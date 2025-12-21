package model

// Config represents a single configuration item
// swagger:model Config
type Config struct {
	// The unique identifier for the configuration
	// example: config-123
	ID string `json:"id"`

	// The name of the configuration
	// example: Database Configuration
	Name string `json:"name"`

	// The version of the configuration
	// example: v1.0.0
	Version string `json:"version"`

	// Key-value pairs representing configuration parameters
	// example: {"host": "localhost", "port": "5432"}
	Parameters map[string]string `json:"parameters"`
}

// ConfigurationGroup represents a group of related configurations
// swagger:model ConfigurationGroup
type ConfigurationGroup struct {
	// The unique identifier for the configuration group
	// example: group-456
	Id string `json:"id"`

	// The name of the configuration group
	// example: Database Configurations
	Name string `json:"name"`

	// The version of the configuration group
	// example: v2.1.0
	Version string `json:"version"`

	// Array of labeled configurations in this group
	Configurations []*LabeledConfiguration `json:"configurations"`
}

// LabeledConfiguration represents a configuration with associated labels
// swagger:model LabeledConfiguration
type LabeledConfiguration struct {
	// The unique identifier for the labeled configuration
	// example: labeled-config-789
	Id string `json:"id"`

	// The configuration object
	Configuration *Config `json:"configuration"`

	// Key-value pairs representing labels for this configuration
	// example: {"environment": "production", "region": "us-east-1"}
	Labels map[string]string `json:"labels"`
}

// IdempotencyStatus represents the status of an idempotent operation
// swagger:model IdempotencyStatus
type IdempotencyStatus string

const (
	StatusInProgress IdempotencyStatus = "in_progress"
	StatusCompleted  IdempotencyStatus = "completed"
)

// IdempotencyRecord stores information about an idempotent request
// swagger:model IdempotencyRecord
type IdempotencyRecord struct {
	// Status of the request
	// example: completed
	// swagger:enum
	Status IdempotencyStatus `json:"status"`

	// HTTP status code returned for the request
	// example: 200
	StatusCode int `json:"statusCode"`

	// Body of the response
	// example: {"id":"config-123","name":"DB Config","version":"v1.0.0"}
	Body string `json:"body"`
}

// ErrorResponse represents a standard error
// swagger:model ErrorResponse
type ErrorResponse struct {
	// Error message
	// example: configuration not found
	Message string `json:"message"`
}

// NoContentResponse represents an empty response
// swagger:model NoContentResponse
type NoContentResponse struct{}
