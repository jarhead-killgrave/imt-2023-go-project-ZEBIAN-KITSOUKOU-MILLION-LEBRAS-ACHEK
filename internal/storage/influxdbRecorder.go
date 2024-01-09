package storage

import (
	"context"
	"github.com/influxdata/influxdb-client-go/v2"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"strconv"
	"sync"
)

type InfluxDBRecorder struct {
	mu     sync.Mutex
	client influxdb2.Client
	bucket string
	org    string
}

type InfluxDBSettings struct {
	URL          string
	Token        string
	Bucket       string
	Organization string
}

func NewInfluxDBRecorder(settings InfluxDBSettings) (*InfluxDBRecorder, error) {
	client := influxdb2.NewClient(settings.URL, settings.Token)
	return &InfluxDBRecorder{
		mu:     sync.Mutex{},
		client: client,
		bucket: settings.Bucket,
		org:    settings.Organization,
	}, nil
}

func (r *InfluxDBRecorder) RecordOnContext(ctx context.Context, m *sensor.Measurement) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	writeAPI := r.client.WriteAPIBlocking(r.org, r.bucket)

	p := influxdb2.NewPointWithMeasurement(string(m.Type)).
		AddTag("airport", m.AirportID).
		AddTag("sensor", strconv.FormatInt(m.SensorID, 10)).
		AddField("value", m.Value).
		SetTime(m.Timestamp)

	if err := writeAPI.WritePoint(ctx, p); err != nil {
		logutil.Error("Failed to write point: %v", err)
		return err
	}

	return nil
}

// Record stores a measurement
func (r *InfluxDBRecorder) Record(m *sensor.Measurement) error {
	return r.RecordOnContext(context.Background(), m)
}

func (r *InfluxDBRecorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.client.Close()

	return nil
}