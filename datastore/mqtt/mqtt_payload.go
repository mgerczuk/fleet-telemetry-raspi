package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/teslamotors/fleet-telemetry/protos"
	"github.com/teslamotors/fleet-telemetry/telemetry"
)

func (p *Producer) process(rec *telemetry.Record, topLevel string, retain bool, obj any) ([]pahomqtt.Token, error) {
	mqttTopicName := fmt.Sprintf("%s/%s/%s", p.config.TopicBase, rec.Vin, topLevel)
	jsonValue, err := json.Marshal(obj)
	if err != nil {
		return []pahomqtt.Token{}, fmt.Errorf("failed to marshal JSON for MQTT topic %s: %v", mqttTopicName, err)
	}
	p.updateMetrics(rec.TxType, len(jsonValue))
	return []pahomqtt.Token{p.client.Publish(mqttTopicName, p.config.QoS, retain, jsonValue)}, nil
}

func (p *Producer) processVehicleFields(rec *telemetry.Record, payload *protos.Payload) ([]pahomqtt.Token, error) {
	return p.process(rec, //
		"v", //
		false,
		map[string]interface{}{
			"created_at": payload.CreatedAt.AsTime(),
			"vin":        payload.Vin,
			"data":       p.dataToMqtt(payload.Data),
		})
}

func (p *Producer) processVehicleAlerts(rec *telemetry.Record, payload *protos.VehicleAlerts) ([]pahomqtt.Token, error) {
	return p.process(rec, //
		"alerts", //
		false,
		map[string]interface{}{
			"created_at": payload.CreatedAt.AsTime(),
			"vin":        payload.Vin,
			"alerts":     vehicleAlertsToMqtt(payload.Alerts),
		})
}

func (p *Producer) processVehicleErrors(rec *telemetry.Record, payload *protos.VehicleErrors) ([]pahomqtt.Token, error) {
	return p.process(rec, //
		"errors", //
		false,
		map[string]interface{}{
			"created_at": payload.CreatedAt.AsTime(),
			"vin":        payload.Vin,
			"errors":     vehicleErrorsToMqtt(payload.Errors),
		})
}

func (p *Producer) processVehicleConnectivity(rec *telemetry.Record, payload *protos.VehicleConnectivity) ([]pahomqtt.Token, error) {
	return p.process(rec, //
		"connectivity", //
		true,
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

func (p *Producer) dataToMqtt(data []*protos.Datum) []Datum {
	result := make([]Datum, len(data))
	for i, d := range data {
		result[i].Key = d.Key.String()
		result[i].Value = d.Value.Value
	}
	return result
}
