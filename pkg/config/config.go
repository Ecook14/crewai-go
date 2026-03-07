package config

// FrameworkConfig holds framework-wide defaults and feature flags.
type FrameworkConfig struct {
	TelemetryEnabled bool
	DefaultTimeout   int
	LoggingLevel     string // "info", "debug", "warn", "error"
}

// Global defaults
var DefaultConfig = FrameworkConfig{
	TelemetryEnabled: false,
	DefaultTimeout:   30,
	LoggingLevel:     "info",
}

// AgentConfig specific overrides for an individual agent
type AgentConfig struct {
	AllowDelegation bool
	MemoryEnabled   bool
	SelfHealing     bool
	MaxIterations   int
}

// CrewConfig specific overrides for a crew
type CrewConfig struct {
	Verbose        bool
	ProcessTimeout int
	MemoryBackend  string // "sqlite", "redis", "chroma"
}
