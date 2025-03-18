package mqtt_test

import (
	"context"
	"fmt"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus/hooks/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/teslamotors/fleet-telemetry/messages"
	"github.com/teslamotors/fleet-telemetry/protos"
	"github.com/teslamotors/fleet-telemetry/server/airbrake"

	"github.com/teslamotors/fleet-telemetry/datastore/mqtt"
	logrus "github.com/teslamotors/fleet-telemetry/logger"
	"github.com/teslamotors/fleet-telemetry/metrics"
	"github.com/teslamotors/fleet-telemetry/telemetry"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MockMQTTClient struct {
	ConnectFunc           func() pahomqtt.Token
	PublishFunc           func(topic string, qos byte, retained bool, payload interface{}) pahomqtt.Token
	DisconnectFunc        func(quiesce uint)
	IsConnectedFunc       func() bool
	IsConnectionOpenFunc  func() bool
	SubscribeFunc         func(topic string, qos byte, callback pahomqtt.MessageHandler) pahomqtt.Token
	SubscribeMultipleFunc func(filters map[string]byte, callback pahomqtt.MessageHandler) pahomqtt.Token
	UnsubscribeFunc       func(topics ...string) pahomqtt.Token
	AddRouteFunc          func(topic string, callback pahomqtt.MessageHandler)
	OptionsReaderFunc     func() pahomqtt.ClientOptionsReader
}

func (m *MockMQTTClient) Connect() pahomqtt.Token {
	return m.ConnectFunc()
}

func (m *MockMQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) pahomqtt.Token {
	return m.PublishFunc(topic, qos, retained, payload)
}

func (m *MockMQTTClient) Disconnect(quiesce uint) {
	m.DisconnectFunc(quiesce)
}

func (m *MockMQTTClient) IsConnected() bool {
	return m.IsConnectedFunc()
}

func (m *MockMQTTClient) IsConnectionOpen() bool {
	return m.IsConnectionOpenFunc()
}

func (m *MockMQTTClient) Subscribe(topic string, qos byte, callback pahomqtt.MessageHandler) pahomqtt.Token {
	return m.SubscribeFunc(topic, qos, callback)
}

func (m *MockMQTTClient) SubscribeMultiple(filters map[string]byte, callback pahomqtt.MessageHandler) pahomqtt.Token {
	return m.SubscribeMultipleFunc(filters, callback)
}

func (m *MockMQTTClient) Unsubscribe(topics ...string) pahomqtt.Token {
	return m.UnsubscribeFunc(topics...)
}

func (m *MockMQTTClient) AddRoute(topic string, callback pahomqtt.MessageHandler) {
	m.AddRouteFunc(topic, callback)
}

func (m *MockMQTTClient) OptionsReader() pahomqtt.ClientOptionsReader {
	return m.OptionsReaderFunc()
}

type MockToken struct {
	WaitFunc        func() bool
	WaitTimeoutFunc func(time.Duration) bool
	DoneFunc        func() <-chan struct{}
	ErrorFunc       func() error
}

func (m *MockToken) Wait() bool {
	return m.WaitFunc()
}

func (m *MockToken) WaitTimeout(d time.Duration) bool {
	return m.WaitTimeoutFunc(d)
}

func (m *MockToken) Done() <-chan struct{} {
	return m.DoneFunc()
}

func (m *MockToken) Error() error {
	return m.ErrorFunc()
}

var publishedTopics = make(map[string][]byte)

func resetPublishedTopics() {
	publishedTopics = make(map[string][]byte)
}

func mockPahoNewClient(_ *pahomqtt.ClientOptions) pahomqtt.Client {
	return &MockMQTTClient{

		ConnectFunc: func() pahomqtt.Token {
			return &MockToken{
				WaitFunc:  func() bool { return true },
				ErrorFunc: func() error { return nil },
			}
		},
		IsConnectedFunc: func() bool {
			return true
		},
		PublishFunc: func(topic string, _ byte, _ bool, payload interface{}) pahomqtt.Token {
			publishedTopics[topic] = payload.([]byte)
			return &MockToken{
				WaitTimeoutFunc: func(_ time.Duration) bool { return true },
				WaitFunc:        func() bool { return true },
				ErrorFunc:       func() error { return nil },
			}
		},
	}
}

var _ = Describe("MQTTProducer", func() {
	var (
		mockLogger        *logrus.Logger
		mockCollector     metrics.MetricCollector
		mockConfig        *mqtt.Config
		mockAirbrake      *airbrake.Handler
		originalNewClient func(*pahomqtt.ClientOptions) pahomqtt.Client
		loggerHook        *test.Hook
		serializer        *telemetry.BinarySerializer
	)

	BeforeEach(func() {
		resetPublishedTopics()
		originalNewClient = mqtt.PahoNewClient
		mqtt.PahoNewClient = mockPahoNewClient

		mockLogger, loggerHook = logrus.NoOpLogger()
		mockCollector = metrics.NewCollector(nil, mockLogger)
		mockAirbrake = airbrake.NewAirbrakeHandler(nil)
		mockConfig = &mqtt.Config{
			Broker:    "tcp://localhost:1883",
			ClientID:  "test-client",
			Username:  "testuser",
			Password:  "testpass",
			TopicBase: "test/topic",
			QoS:       1,
			Retained:  false,
		}

		serializer = telemetry.NewBinarySerializer(
			&telemetry.RequestIdentity{
				DeviceID: "TEST123",
				SenderID: "vehicle_device.TEST123",
			},
			map[string][]telemetry.Producer{},
			mockLogger,
		)

	})

	AfterEach(func() {
		mqtt.PahoNewClient = originalNewClient
	})

	Describe("Produce", func() {
		It("should publish MQTT messages for each field in the payload", func() {
			producer, err := mqtt.NewProducer(
				context.Background(),
				mockConfig,
				mockCollector,
				"test_namespace",
				mockAirbrake,
				nil,
				nil,
				mockLogger,
			)
			Expect(err).NotTo(HaveOccurred())

			createdAt := timestamppb.Now()

			payload := &protos.Payload{
				Vin: "TEST123",
				Data: []*protos.Datum{
					{
						Key: protos.Field_VehicleName,
						Value: &protos.Value{
							Value: &protos.Value_StringValue{StringValue: "My Tesla"},
						},
					},
					{
						Key: protos.Field_TimeToFullCharge,
						Value: &protos.Value{
							Value: &protos.Value_Invalid{Invalid: true},
						},
					},
					{
						Key: protos.Field_Location,
						Value: &protos.Value{
							Value: &protos.Value_LocationValue{
								LocationValue: &protos.LocationValue{
									Latitude:  37.7749,
									Longitude: -122.4194,
								},
							},
						},
					},
					{
						Key: protos.Field_BatteryLevel,
						Value: &protos.Value{
							Value: &protos.Value_FloatValue{FloatValue: 75.5},
						},
					},
				},
				CreatedAt: createdAt,
			}

			payloadBytes, err := proto.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())

			// Create stream message
			message := messages.StreamMessage{
				TXID:         []byte("1234"),
				SenderID:     []byte("vehicle_device.TEST123"),
				MessageTopic: []byte("V"),
				Payload:      payloadBytes,
			}
			msgBytes, err := message.ToBytes()
			Expect(err).NotTo(HaveOccurred())

			// Create record properly using NewRecord
			record, err := telemetry.NewRecord(serializer, msgBytes, "1", true)
			Expect(err).NotTo(HaveOccurred())

			producer.Produce(record)

			Expect(publishedTopics).To(HaveLen(1))

			topic := "test/topic/TEST123/v"

			jsonValue := `{"data":[` +
				`{"key":"VehicleName","value":{"stringValue":"My Tesla"}},` +
				`{"key":"TimeToFullCharge","value":{"invalid":true}},` +
				`{"key":"Location","value":{"locationValue":{"latitude":37.7749,"longitude":-122.4194}}},` +
				`{"key":"BatteryLevel","value":{"floatValue":75.5}}` +
				`],` +
				`"createdAt":"` + createdAt.AsTime().Format(time.RFC3339Nano) + `",` +
				`"vin":"TEST123"}`

			Expect(publishedTopics).To(HaveKey(topic))
			Expect(string(publishedTopics[topic])).To(Equal(jsonValue))
		})

		It("should publish MQTT messages for vehicle alerts", func() {

			producer, err := mqtt.NewProducer(
				context.Background(),
				mockConfig,
				mockCollector,
				"test_namespace",
				mockAirbrake,
				nil,
				nil,
				mockLogger,
			)
			Expect(err).NotTo(HaveOccurred())

			createdAt := timestamppb.Now()
			alerts := &protos.VehicleAlerts{
				Vin: "TEST123",
				Alerts: []*protos.VehicleAlert{
					{
						Name:      "TestAlert1",
						StartedAt: createdAt,
						EndedAt:   nil,
						Audiences: []protos.Audience{protos.Audience_Customer, protos.Audience_Service},
					},
					{
						Name:      "TestAlert2",
						StartedAt: createdAt,
						EndedAt:   createdAt,
						Audiences: []protos.Audience{protos.Audience_ServiceFix},
					},
				},
				CreatedAt: createdAt,
			}

			alertsBytes, err := proto.Marshal(alerts)
			Expect(err).NotTo(HaveOccurred())

			// Create stream message
			message := messages.StreamMessage{
				TXID:         []byte("1234"),
				SenderID:     []byte("vehicle_device.TEST123"),
				MessageTopic: []byte("alerts"),
				Payload:      alertsBytes,
			}
			msgBytes, err := message.ToBytes()
			Expect(err).NotTo(HaveOccurred())

			// Create record properly using NewRecord
			record, err := telemetry.NewRecord(serializer, msgBytes, "1", true)
			Expect(err).NotTo(HaveOccurred())

			producer.Produce(record)

			Expect(publishedTopics).To(HaveLen(1))

			topic := "test/topic/TEST123/alerts"

			createdAtAsTime := createdAt.AsTime()

			jsonValue := `{"alerts":[` +
				`{"name":"TestAlert1","audiences":["Customer","Service"],"startedAt":"` + createdAtAsTime.Format(time.RFC3339Nano) + `"},` +
				`{"name":"TestAlert2","audiences":["ServiceFix"],"startedAt":"` + createdAtAsTime.Format(time.RFC3339Nano) + `","endedAt":"` + createdAtAsTime.Format(time.RFC3339Nano) + `"}` +
				`],` +
				`"createdAt":"` + createdAtAsTime.Format(time.RFC3339Nano) + `",` +
				`"vin":"TEST123"}`

			Expect(publishedTopics).To(HaveKey(topic))
			Expect(string(publishedTopics[topic])).To(Equal(jsonValue))
		})

		It("should publish MQTT messages for vehicle errors", func() {
			producer, err := mqtt.NewProducer(
				context.Background(),
				mockConfig,
				mockCollector,
				"test_namespace",
				nil,
				nil,
				nil,
				mockLogger,
			)
			Expect(err).NotTo(HaveOccurred())

			createdAt := timestamppb.Now()
			vehicleErrors := &protos.VehicleErrors{
				Vin: "TEST123",
				Errors: []*protos.VehicleError{
					{
						Name:      "TestError1",
						Body:      "This is a test error",
						Tags:      map[string]string{"tag1": "value1", "tag2": "value2"},
						CreatedAt: createdAt,
					},
					{
						Name:      "TestError2",
						Body:      "This is another test error",
						Tags:      map[string]string{"tagA": "valueA"},
						CreatedAt: createdAt,
					},
				},
				CreatedAt: createdAt,
			}

			errorsBytes, err := proto.Marshal(vehicleErrors)
			Expect(err).NotTo(HaveOccurred())

			message := messages.StreamMessage{
				TXID:         []byte("1234"),
				SenderID:     []byte("vehicle_device.TEST123"),
				MessageTopic: []byte("errors"),
				Payload:      errorsBytes,
			}
			msgBytes, err := message.ToBytes()
			Expect(err).NotTo(HaveOccurred())

			// Create record properly using NewRecord
			record, err := telemetry.NewRecord(serializer, msgBytes, "1", true)
			Expect(err).NotTo(HaveOccurred())

			producer.Produce(record)

			Expect(publishedTopics).To(HaveLen(1))

			topic := "test/topic/TEST123/errors"

			createdAtAsTime := createdAt.AsTime()

			jsonValue := `{"errors":[` +
				`{"createdAt":"` + createdAtAsTime.Format(time.RFC3339Nano) + `","name":"TestError1","tags":{"tag1":"value1","tag2":"value2"},"body":"This is a test error"},` +
				`{"createdAt":"` + createdAtAsTime.Format(time.RFC3339Nano) + `","name":"TestError2","tags":{"tagA":"valueA"},"body":"This is another test error"}` +
				`],` +
				`"createdAt":"` + createdAtAsTime.Format(time.RFC3339Nano) + `",` +
				`"vin":"TEST123"}`

			fmt.Println(string(publishedTopics[topic]))
			fmt.Println(jsonValue)
			Expect(publishedTopics).To(HaveKey(topic))
			Expect(string(publishedTopics[topic])).To(Equal(jsonValue))
		})

		It("should publish MQTT messages for vehicle connectivity", func() {
			producer, err := mqtt.NewProducer(
				context.Background(),
				mockConfig,
				mockCollector,
				"test_namespace",
				nil,
				nil,
				nil,
				mockLogger,
			)
			Expect(err).NotTo(HaveOccurred())

			createdAt := timestamppb.Now()
			connectivity := &protos.VehicleConnectivity{
				Vin:              "TEST123",
				ConnectionId:     "connid",
				Status:           protos.ConnectivityEvent_DISCONNECTED,
				CreatedAt:        createdAt,
				NetworkInterface: "xyz",
			}

			errorsBytes, err := proto.Marshal(connectivity)
			Expect(err).NotTo(HaveOccurred())

			message := messages.StreamMessage{
				TXID:         []byte("1234"),
				SenderID:     []byte("vehicle_device.TEST123"),
				MessageTopic: []byte("connectivity"),
				Payload:      errorsBytes,
			}
			msgBytes, err := message.ToBytes()
			Expect(err).NotTo(HaveOccurred())

			// Create record properly using NewRecord
			record, err := telemetry.NewRecord(serializer, msgBytes, "1", true)
			Expect(err).NotTo(HaveOccurred())

			producer.Produce(record)

			Expect(publishedTopics).To(HaveLen(1))

			topic := "test/topic/TEST123/connectivity"

			createdAtAsTime := createdAt.AsTime()

			jsonValue := `{"vin":"TEST123","connectionId":"connid",` +
				`"status":"DISCONNECTED",` +
				`"createdAt":"` + createdAtAsTime.Format(time.RFC3339Nano) + `",` +
				`"networkInterface":"xyz"}`

			Expect(publishedTopics).To(HaveKey(topic))
			Expect(string(publishedTopics[topic])).To(Equal(jsonValue))
		})

		It("should handle timeouts when publishing MQTT messages", func() {
			// Mock a slow publish function that always times out
			mqtt.PahoNewClient = func(o *pahomqtt.ClientOptions) pahomqtt.Client {
				return &MockMQTTClient{
					ConnectFunc: func() pahomqtt.Token {
						return &MockToken{
							WaitFunc:  func() bool { return true },
							ErrorFunc: func() error { return nil },
						}
					},
					IsConnectedFunc: func() bool {
						return true
					},
					PublishFunc: func(topic string, qos byte, retained bool, payload interface{}) pahomqtt.Token {
						return &MockToken{
							WaitTimeoutFunc: func(d time.Duration) bool { return false },
							WaitFunc:        func() bool { return false },
							ErrorFunc:       func() error { return pahomqtt.TimedOut },
						}
					},
				}
			}

			producer, err := mqtt.NewProducer(
				context.Background(),
				mockConfig,
				mockCollector,
				"test_namespace",
				mockAirbrake,
				nil,
				nil,
				mockLogger,
			)
			Expect(err).NotTo(HaveOccurred())

			payload := &protos.Payload{
				Vin: "TEST123",
				Data: []*protos.Datum{
					{
						Key: protos.Field_VehicleName,
						Value: &protos.Value{
							Value: &protos.Value_StringValue{StringValue: "My Tesla"},
						},
					},
				},
				CreatedAt: timestamppb.Now(),
			}

			payloadBytes, err := proto.Marshal(payload)
			Expect(err).NotTo(HaveOccurred())

			message := messages.StreamMessage{
				TXID:         []byte("1234"),
				SenderID:     []byte("vehicle_device.TEST123"),
				MessageTopic: []byte("V"),
				Payload:      payloadBytes,
			}
			msgBytes, err := message.ToBytes()
			Expect(err).NotTo(HaveOccurred())

			// Create record properly using NewRecord
			record, err := telemetry.NewRecord(serializer, msgBytes, "1", true)
			Expect(err).NotTo(HaveOccurred())

			producer.Produce(record)

			// Check that an error was logged
			Expect(loggerHook.LastEntry().Message).To(Equal("mqtt_publish_error"))

		})
	})
})
