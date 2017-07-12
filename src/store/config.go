package store

// Config represents the configuration for all savers.
type Config struct {
	Master        bool   // Whether node is master.
	TestMode      bool   // Whether node is in test mode.
	MemoryMode    bool   // Whether to use local storage in runtime.
	InternalState bool   // Whether to enable internal state.
	ConfigDir     string // Configuration directory.
	CXOPort       int    // CXO listening port.
	CXORPCPort    int    // CXO RPC listening port.
}
