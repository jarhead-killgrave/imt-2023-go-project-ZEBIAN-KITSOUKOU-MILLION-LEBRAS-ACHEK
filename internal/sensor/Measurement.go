package sensor

import (
	"encoding/json"
	"fmt"
	"time"
)

// Measurement represents the data from a sensor
type Measurement struct {
	SensorID  int64     `json:"sensor_id"`
	AirportID string    `json:"airport_id"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
}

func (m *Measurement) String() string {
	return fmt.Sprintf(
		"SensorID: %d, AirportID: %s, Value: %f, Unit: %s, Timestamp: %s",
		m.SensorID, m.AirportID, m.Value, m.Unit, m.Timestamp,
	)
}

func (m *Measurement) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Measurement) ToCSV(separator rune, timeFormat string) string {
	return fmt.Sprintf(
		"%d%c%s%c%f%c%s%c%s",
		m.SensorID, separator, m.AirportID, separator, m.Value, separator, m.Unit, separator, m.Timestamp.Format(timeFormat),
	)
}

// Static method to get the field names of a measurement
func MeasurementFieldNames(separator rune) string {
	return fmt.Sprintf(
		"sensor_id%cairport_id%cvalue%cunit%ctimestamp",
		separator, separator, separator, separator,
	)
}
