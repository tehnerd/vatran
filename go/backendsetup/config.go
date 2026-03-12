package backendsetup

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// katranConfig is a minimal representation of the katran YAML config,
// used only to extract VIP addresses. This avoids importing the full
// server package which has CGO dependencies.
type katranConfig struct {
	VIPs []vipEntry `yaml:"vips"`
}

// vipEntry represents a single VIP entry in the config file.
type vipEntry struct {
	Address string `yaml:"address"`
}

// extractVIPAddresses parses YAML config data and returns the unique
// VIP addresses from the vips section.
//
// Parameters:
//   - data: raw YAML config file contents.
//
// Returns deduplicated VIP addresses or an error if parsing fails.
func extractVIPAddresses(data []byte) ([]string, error) {
	var cfg katranConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	vips := make([]string, 0, len(cfg.VIPs))
	for _, v := range cfg.VIPs {
		if v.Address != "" {
			vips = append(vips, v.Address)
		}
	}
	return DeduplicateVIPs(vips), nil
}
