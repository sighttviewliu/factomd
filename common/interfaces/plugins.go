package interfaces

// IManagerPlugin is the interface we are exposing as a plugin. It is
// not directly a manager interface, as we have to handle goroutines
// in the plugin
type IManagerController interface {
	// Manager functions extended
	RetrieveDBStateByHeight(height uint32) error
	UploadDBStateBytes(data []byte, sign bool) error
	RequestMoreUploads() int

	// Control function
	IsBufferEmpty() bool
	FetchFromBuffer() []byte
	SetSigningKey(sec []byte) error

	// Plugin Control
	Alive() error
}