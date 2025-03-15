package config

const TestConfig = `{
	"host": "127.0.0.1",
	"port": 443,
	"status_port": 8080,
	"log_level": "info",
	"json_log_enable": true,
	"namespace": "tesla_telemetry",
	"reliable_ack_sources": {
		"V": "mqtt"
	},
	"mqtt": {
		"broker": "mqtt:1883",
		"client_id": "client-1",
		"topic_base": "telemetry",
		"qos": 1,
		"retained": false,
		"connect_timeout_ms": 30000,
		"publish_timeout_ms": 1000 
	},
	"monitoring": {
		"prometheus_metrics_port": 9090,
		"profiler_port": 4269,
		"profiling_path": "/tmp/fleet-telemetry/profile"
	},
	"rate_limit": {
		"enabled": true,
		"message_interval_time": 30,
		"message_limit": 1000
	},
	"records": {
		"V": ["mqtt"]
	},
	"tls": {
		"ca_file": "tesla.ca",
		"server_cert": "your_own_cert.crt",
		"server_key": "your_own_key.key"
	}
}
`

const TestSmallConfig = `
{
	"host": "127.0.0.1",
	"port": 443,
	"status_port": 8080,
	"namespace": "tesla_telemetry",
	"mqtt": {
		"broker": "mqtt:1883",
		"client_id": "client-1",
		"topic_base": "telemetry",
		"qos": 1,
		"retained": false,
		"connect_timeout_ms": 30000,
		"publish_timeout_ms": 1000
	},
	"records": {
		"V": ["mqtt"]
	},
	"tls": {
		"ca_file": "tesla.ca",
		"server_cert": "your_own_cert.crt",
		"server_key": "your_own_key.key"
	}
}
`

const BadVinsConfig = `
{
	"host": "127.0.0.1",
	"port": 443,
	"status_port": 8080,
	"namespace": "tesla_telemetry",
	"kafka": {
		"bootstrap.servers": "some.broker1:9093,some.broker1:9093",
		"ssl.ca.location": "kafka.ca",
		"ssl.certificate.location": "kafka.crt",
		"ssl.key.location": "kafka.key",
		"queue.buffering.max.messages": 1000000
	},
	"records": {
		"V": ["kafka"]
	},
	"vins_signal_tracking_enabled": ["vin1", "vin2", "vin3"],
	"tls": {
		"ca_file": "tesla.ca",
		"server_cert": "your_own_cert.crt",
		"server_key": "your_own_key.key"
	}
}
`
const TestBadReliableAckConfig = `
{
	"host": "127.0.0.1",
	"port": 443,
	"status_port": 8080,
	"namespace": "tesla_telemetry",
	"reliable_ack_sources": {
		"V": "pubsub"
	},
	"kafka": {
		"bootstrap.servers": "some.broker1:9093,some.broker1:9093",
		"ssl.ca.location": "kafka.ca",
		"ssl.certificate.location": "kafka.crt",
		"ssl.key.location": "kafka.key",
		"queue.buffering.max.messages": 1000000
	},
	"records": {
		"V": ["kafka"]
	},
	"tls": {
		"ca_file": "tesla.ca",
		"server_cert": "your_own_cert.crt",
		"server_key": "your_own_key.key"
	}
}
`

const TestLoggerAsReliableAckConfig = `
{
	"host": "127.0.0.1",
	"port": 443,
	"status_port": 8080,
	"namespace": "tesla_telemetry",
	"reliable_ack_sources": {
		"V": "logger"
	},
	"kafka": {
		"bootstrap.servers": "some.broker1:9093,some.broker1:9093",
		"ssl.ca.location": "kafka.ca",
		"ssl.certificate.location": "kafka.crt",
		"ssl.key.location": "kafka.key",
		"queue.buffering.max.messages": 1000000
	},
	"records": {
		"V": ["kafka", "logger"]
	},
	"tls": {
		"ca_file": "tesla.ca",
		"server_cert": "your_own_cert.crt",
		"server_key": "your_own_key.key"
	}
}
`

const TestUnusedTxTypeAsReliableAckConfig = `
{
	"host": "127.0.0.1",
	"port": 443,
	"status_port": 8080,
	"namespace": "tesla_telemetry",
	"reliable_ack_sources": {
		"error": "kafka"
	},
	"kafka": {
		"bootstrap.servers": "some.broker1:9093,some.broker1:9093",
		"ssl.ca.location": "kafka.ca",
		"ssl.certificate.location": "kafka.crt",
		"ssl.key.location": "kafka.key",
		"queue.buffering.max.messages": 1000000
	},
	"records": {
		"V": ["kafka", "logger"]
	},
	"tls": {
		"ca_file": "tesla.ca",
		"server_cert": "your_own_cert.crt",
		"server_key": "your_own_key.key"
	}
}
`

const TestPubsubConfig = `
{
	"host": "127.0.0.1",
	"port": 443,
	"status_port": 8080,
	"pubsub": {
        "gcp_project_id": "some-project-id",
		"reliable_ack": "true"
    },
	"records": {
		"V": ["pubsub"]
	}
}
`

const BadTopicConfig = `
{
	"host": "127.0.0.1",
	"port": "",
}`

const TestZMQConfig = `
{
  "host": "127.0.0.1",
  "port": 443,
  "status_port": 8080,
  "zmq": {
    "addr": "tcp://127.0.0.1:5288"
  },
  "records": {
    "V": ["zmq"]
  }
}
`

const TestTransmitDecodedRecords = `
{
	"host": "127.0.0.1",
	"port": 443,
	"status_port": 8080,
	"transmit_decoded_records": true,
	"records": {
		"V": ["logger"]
	}
}
`

const TestVinsToTrackConfig = `
{
	"host": "127.0.0.1",
	"port": 443,
	"status_port": 8080,
	"transmit_decoded_records": true,
	"records": {
		"V": ["logger"]
	},
	"vins_signal_tracking_enabled": ["v1", "v2"]
}
`

const TestAirbrakeConfig = `
{
	"host": "127.0.0.1",
	"port": 443,
	"status_port": 8080,
	"namespace": "tesla_telemetry",
	"kafka": {
		"bootstrap.servers": "some.broker1:9093,some.broker1:9093",
		"ssl.ca.location": "kafka.ca",
		"ssl.certificate.location": "kafka.crt",
		"ssl.key.location": "kafka.key",
		"queue.buffering.max.messages": 1000000
	},
	"records": {
		"V": ["kafka"]
	},
	"tls": {
		"ca_file": "tesla.ca",
		"server_cert": "your_own_cert.crt",
		"server_key": "your_own_key.key"
	},
	"airbrake": {
        "project_id": 1,
        "project_key": "test1",
        "environment": "integration",
        "host": "http://errbit-test.example.com"
    }
}
`

const TestBadTxTypeReliableAckConfig = `
{
	"host": "127.0.0.1",
	"port": 443,
	"status_port": 8080,
	"namespace": "tesla_telemetry",
	"reliable_ack_sources": {
		"connectivity": "kafka"
	},
	"kafka": {
		"bootstrap.servers": "some.broker1:9093,some.broker1:9093",
		"ssl.ca.location": "kafka.ca",
		"ssl.certificate.location": "kafka.crt",
		"ssl.key.location": "kafka.key",
		"queue.buffering.max.messages": 1000000
	},
	"records": {
		"V": ["kafka"],
		"connectivity": ["kafka"]
	},
	"tls": {
		"server_cert": "your_own_cert.crt",
		"server_key": "your_own_key.key"
	}
}
`
