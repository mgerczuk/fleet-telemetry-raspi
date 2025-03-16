package mqtt

import (
	"fmt"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/teslamotors/fleet-telemetry/telemetry"
	"google.golang.org/protobuf/encoding/protojson"
)

func (p *Producer) process(rec *telemetry.Record, topLevel string, retain bool) ([]pahomqtt.Token, error) {
	mqttTopicName := fmt.Sprintf("%s/%s/%s", p.config.TopicBase, rec.Vin, topLevel)
	jsonValue, err := protojson.Marshal(rec.GetProtoMessage())
	if err != nil {
		return []pahomqtt.Token{}, fmt.Errorf("failed to marshal JSON for MQTT topic %s: %v", mqttTopicName, err)
	}
	p.updateMetrics(rec.TxType, len(jsonValue))
	return []pahomqtt.Token{p.client.Publish(mqttTopicName, p.config.QoS, retain, jsonValue)}, nil
}

func (p *Producer) processVehicleFields(rec *telemetry.Record) ([]pahomqtt.Token, error) {
	return p.process(rec, "v", false)
}

func (p *Producer) processVehicleAlerts(rec *telemetry.Record) ([]pahomqtt.Token, error) {
	return p.process(rec, "alerts", false)
}

func (p *Producer) processVehicleErrors(rec *telemetry.Record) ([]pahomqtt.Token, error) {
	return p.process(rec, "errors", false)
}

func (p *Producer) processVehicleConnectivity(rec *telemetry.Record) ([]pahomqtt.Token, error) {
	return p.process(rec, "connectivity", true)
}
