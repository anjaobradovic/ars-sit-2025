package model

type Config struct {
	Id      string            `json:"id"`
	Name    string            `json: "name"`
	Version string            `json: "version"`
	Params  map[string]string `json: "params"`
}

// todo: dodati metode

type ConfigRepository interface {
	// todo: dodati metode
	Add(config Config)
	Get(name string, version int) (Config, error)
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
