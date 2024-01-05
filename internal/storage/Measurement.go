package storage

import (
	"encoding/json"
	"fmt"
	"time"
)

// Measurement represents the data from a sensor
type Measurement struct {
	SensorID  int64     `json:"sensor_id"`
	AirportID string    `json:"airport_id"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
}

func (m *Measurement) String() string {
	return fmt.Sprintf(
		"SensorID: %d, AirportID: %s, Type: %s, Value: %f, Unit: %s, Timestamp: %s",
		m.SensorID, m.AirportID, m.Type, m.Value, m.Unit, m.Timestamp.Format(time.RFC3339),
	)
}

func (m *Measurement) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Measurement) FromJSON(data []byte) error {
	return json.Unmarshal(data, m)
}
