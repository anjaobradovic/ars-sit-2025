package model

type Config struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	Parameters map[string]string `json:"parameters"`
}

type ConfigRepository interface {
	Add(config Config) error
	Get(name string, version string) (Config, error)
}

type ConfigurationGroup struct {
	Id             string                  `json:"id"`
	Name           string                  `json:"name"`
	Version        string                  `json:"version"`
	Configurations []*LabeledConfiguration `json:"configurations"`
}

type LabeledConfiguration struct {
	Id            string            `json:"id"`
	Configuration *Config           `json:"configuration"`
	Labels        map[string]string `json:"labels"`
}
