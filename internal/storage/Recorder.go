package storage

type Recorder interface {
	// Record stores a measurement
	Record(m *Measurement) error

	// Initialize initializes resources for recording
	Initialize() error

	// Close closes the recorder
	Close() error
}
