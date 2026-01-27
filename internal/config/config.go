package config

// Config holds the complete resolved application configuration
type Config struct {
	Metrics  []MetricConfig
	Export   ExportConfig
	Settings SettingsConfig
}
