package configresolve

import (
	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/obsbox/internal/configparse"
)

// resolveSettings converts raw settings config to resolved settings config
func resolveSettings(raw *configparse.RawSettingsConfig) (config.SettingsConfig, error) {
	result := config.SettingsConfig{
		Seed: raw.Seed,
		InternalMetrics: config.InternalMetricsConfig{
			Enabled: raw.InternalMetrics.Enabled,
			Format:  config.NamingFormat(raw.InternalMetrics.Format),
		},
	}

	// Validate converted config
	if err := result.Validate(); err != nil {
		return config.SettingsConfig{}, err
	}

	return result, nil
}
