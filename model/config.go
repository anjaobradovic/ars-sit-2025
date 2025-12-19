package model

type Config struct {
	ID         string            `json:"id"` // UUID
	Name       string            `json:"name"`
	Version    string            `json:"version"` // npr. "v1", "v2"
	Parameters map[string]string `json:"parameters"`
}

type ConfigRepository interface {
	Add(config Config) error                                     // dodaje novu verziju
	GetByIDAndVersion(id string, version string) (Config, error) // dohvat po UUID + verzija
	DeleteByIDAndVersion(id string, version string) error        // brisanje po UUID + verzija
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

type IdempotencyStatus string

const (
	StatusInProgress IdempotencyStatus = "in_progress"
	StatusCompleted  IdempotencyStatus = "completed"
)

type IdempotencyRecord struct {
	Status     IdempotencyStatus `json:"status"`
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body"`
}
