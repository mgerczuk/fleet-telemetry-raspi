package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/teslamotors/fleet-telemetry/protos"
	"github.com/teslamotors/fleet-telemetry/telemetry"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (p *MQTTProducer) process(rec *telemetry.Record, topLevel string, obj any) ([]pahomqtt.Token, error) {
	mqttTopicName := fmt.Sprintf("%s/%s/%s", p.config.TopicBase, rec.Vin, topLevel)
	jsonValue, err := json.Marshal(obj)
	if err != nil {
		return []pahomqtt.Token{}, fmt.Errorf("failed to marshal JSON for MQTT topic %s: %v", mqttTopicName, err)
	}
	p.updateMetrics(rec.TxType, len(jsonValue))
	return []pahomqtt.Token{p.client.Publish(mqttTopicName, p.config.QoS, p.config.Retained, jsonValue)}, nil
}

func (p *MQTTProducer) processVehicleFields(rec *telemetry.Record, payload *protos.Payload) ([]pahomqtt.Token, error) {
	return p.process(rec, //
		"v", //
		map[string]interface{}{
			"created_at": payload.CreatedAt.AsTime(),
			"vin":        payload.Vin,
			"data":       p.dataToMqtt(payload.Data),
		})
}

func (p *MQTTProducer) processVehicleAlerts(rec *telemetry.Record, payload *protos.VehicleAlerts) ([]pahomqtt.Token, error) {
	return p.process(rec, //
		"alerts", //
		map[string]interface{}{
			"created_at": payload.CreatedAt.AsTime(),
			"vin":        payload.Vin,
			"alerts":     vehicleAlertsToMqtt(payload.Alerts),
		})
}

func (p *MQTTProducer) processVehicleErrors(rec *telemetry.Record, payload *protos.VehicleErrors) ([]pahomqtt.Token, error) {
	return p.process(rec, //
		"errors", //
		map[string]interface{}{
			"created_at": payload.CreatedAt.AsTime(),
			"vin":        payload.Vin,
			"errors":     vehicleErrorsToMqtt(payload.Errors),
		})
}

func (p *MQTTProducer) processVehicleConnectivity(rec *telemetry.Record, payload *protos.VehicleConnectivity) ([]pahomqtt.Token, error) {
	return p.process(rec, //
		"connectivity", //
		map[string]interface{}{
			"vin":               payload.Vin,
			"connection_id":     payload.GetConnectionId(),
			"status":            payload.GetStatus().String(),
			"created_at":        payload.GetCreatedAt().AsTime(),
			"network_interface": payload.NetworkInterface,
		})
}

type VehicleAlert struct {
	Name      string            `json:"name,omitempty"`
	Audiences []protos.Audience `json:"audiences,omitempty"`
	StartedAt *time.Time        `json:"started_at,omitempty"`
	EndedAt   *time.Time        `json:"ended_at,omitempty"`
}

func vehicleAlertsToMqtt(alerts []*protos.VehicleAlert) []VehicleAlert {
	result := make([]VehicleAlert, len(alerts))
	for i, d := range alerts {
		result[i].Name = d.Name
		result[i].Audiences = d.Audiences
		if d.StartedAt != nil {
			time := d.StartedAt.AsTime()
			result[i].StartedAt = &time
		}
		if d.EndedAt != nil {
			time := d.EndedAt.AsTime()
			result[i].EndedAt = &time
		}
	}
	return result
}

type VehicleError struct {
	CreatedAt *time.Time        `json:"created_at,omitempty"`
	Name      string            `json:"name,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
	Body      string            `json:"body,omitempty"`
}

func vehicleErrorsToMqtt(vehicleErrors []*protos.VehicleError) []VehicleError {
	result := make([]VehicleError, len(vehicleErrors))
	for i, d := range vehicleErrors {
		if d.CreatedAt != nil {
			time := d.CreatedAt.AsTime()
			result[i].CreatedAt = &time
		}
		result[i].Name = d.Name
		result[i].Tags = d.Tags
		result[i].Body = d.Body
	}
	return result
}

type Datum struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

func (p *MQTTProducer) dataToMqtt(data []*protos.Datum) []Datum {
	result := make([]Datum, len(data))
	for i, d := range data {
		result[i].Key = d.Key.String()
		result[i].Value = getDatumValue(d.Value)
	}
	return result
}

func getDatumValue(value *protos.Value) interface{} {
	// ordered by expected frequency (see payload.go transformValue)
	switch v := value.Value.(type) {
	case *protos.Value_StringValue:
		return v.StringValue
	case *protos.Value_LocationValue:
		return map[string]float64{
			"latitude":  v.LocationValue.Latitude,
			"longitude": v.LocationValue.Longitude,
		}
	case *protos.Value_FloatValue:
		return v.FloatValue
	case *protos.Value_IntValue:
		return v.IntValue
	case *protos.Value_DoubleValue:
		return v.DoubleValue
	case *protos.Value_LongValue:
		return v.LongValue
	case *protos.Value_BooleanValue:
		return v.BooleanValue
	case *protos.Value_Invalid:
		return nil
	case *protos.Value_ShiftStateValue:
		return v.ShiftStateValue.String()
	case *protos.Value_LaneAssistLevelValue:
		return v.LaneAssistLevelValue.String()
	case *protos.Value_ScheduledChargingModeValue:
		return v.ScheduledChargingModeValue.String()
	case *protos.Value_SentryModeStateValue:
		return v.SentryModeStateValue.String()
	case *protos.Value_SpeedAssistLevelValue:
		return v.SpeedAssistLevelValue.String()
	case *protos.Value_BmsStateValue:
		return v.BmsStateValue.String()
	case *protos.Value_BuckleStatusValue:
		return v.BuckleStatusValue.String()
	case *protos.Value_CarTypeValue:
		return v.CarTypeValue.String()
	case *protos.Value_ChargePortValue:
		return v.ChargePortValue.String()
	case *protos.Value_ChargePortLatchValue:
		return v.ChargePortLatchValue.String()
	case *protos.Value_DoorValue:
		return map[string]bool{
			"DriverFront":    v.DoorValue.DriverFront,
			"PassengerFront": v.DoorValue.PassengerFront,
			"DriverRear":     v.DoorValue.DriverRear,
			"PassengerRear":  v.DoorValue.PassengerRear,
			"TrunkFront":     v.DoorValue.TrunkFront,
			"TrunkRear":      v.DoorValue.TrunkRear,
		}
	case *protos.Value_DriveInverterStateValue:
		return v.DriveInverterStateValue.String()
	case *protos.Value_HvilStatusValue:
		return v.HvilStatusValue.String()
	case *protos.Value_WindowStateValue:
		return v.WindowStateValue.String()
	case *protos.Value_SeatFoldPositionValue:
		return v.SeatFoldPositionValue.String()
	case *protos.Value_TractorAirStatusValue:
		return v.TractorAirStatusValue.String()
	case *protos.Value_FollowDistanceValue:
		return v.FollowDistanceValue.String()
	case *protos.Value_ForwardCollisionSensitivityValue:
		return v.ForwardCollisionSensitivityValue.String()
	case *protos.Value_GuestModeMobileAccessValue:
		return v.GuestModeMobileAccessValue.String()
	case *protos.Value_TrailerAirStatusValue:
		return v.TrailerAirStatusValue.String()
	case *protos.Value_TimeValue:
		return fmt.Sprintf("%02d:%02d:%02d", v.TimeValue.Hour, v.TimeValue.Minute, v.TimeValue.Second)
	case *protos.Value_DetailedChargeStateValue:
		return v.DetailedChargeStateValue.String()
	default:
		// Any other value will be processed using the generic getProtoValue function
		return getProtoValue(value, true)
	}
}

// getProtoValue returns the value as an interface{} for consumption in a mqtt payload
func getProtoValue(protoMsg protoreflect.ProtoMessage, returnAsToplevel bool) interface{} {
	if protoMsg == nil || !protoMsg.ProtoReflect().IsValid() {
		return nil
	}

	m := protoMsg.ProtoReflect()
	result := make(map[string]interface{})

	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch fd.Kind() {
		case protoreflect.BoolKind:
			result[string(fd.Name())] = v.Bool()
		case protoreflect.StringKind:
			result[string(fd.Name())] = v.String()
		case protoreflect.Int32Kind, protoreflect.Int64Kind,
			protoreflect.Sint32Kind, protoreflect.Sint64Kind,
			protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
			result[string(fd.Name())] = v.Int()
		case protoreflect.Uint32Kind, protoreflect.Uint64Kind,
			protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
			result[string(fd.Name())] = v.Uint()
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			result[string(fd.Name())] = v.Float()
		case protoreflect.BytesKind:
			result[string(fd.Name())] = v.Bytes()
		case protoreflect.EnumKind:
			if desc := fd.Enum().Values().ByNumber(v.Enum()); desc != nil {
				result[string(fd.Name())] = string(desc.Name())
			} else {
				result[string(fd.Name())] = int32(v.Enum())
			}
		case protoreflect.MessageKind, protoreflect.GroupKind:
			if msg := v.Message(); msg.IsValid() {
				result[string(fd.Name())] = getProtoValue(msg.Interface(), false)
			}
		}
		return true
	})

	// If there's only a single toplevel field, return its value directly
	if returnAsToplevel && len(result) == 1 {
		for _, v := range result {
			return v
		}
	}

	return result
}
