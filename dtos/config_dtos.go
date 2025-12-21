package dtos

// ConfigurationGroupConfigurationDto represents a configuration reference in a group
// swagger:model ConfigurationGroupConfigurationDto
type ConfigurationGroupConfigurationDto struct {
	// The ID of the configuration to include in the group
	// example: config-123
	Id string `json:"id"`

	// Labels to apply to this configuration in the group context
	// example: {"environment":"production","region":"us-east-1"}
	Labels map[string]string `json:"labels"`
}

// CreateConfigurationDto represents the request body for creating a new configuration
// swagger:model CreateConfigurationDto
type CreateConfigurationDto struct {
	// The name of the configuration
	// example: database-config
	Name string `json:"name"`

	// The version of the configuration
	// example: v1.0
	Version string `json:"version"`

	// Key-value pairs representing configuration parameters
	// example: {"db.host":"localhost","db.port":"5432"}
	Parameters map[string]string `json:"parameters"`
}

// ConfigurationGroupDto represents the request/response body for configuration groups
// swagger:model ConfigurationGroupDto
type ConfigurationGroupDto struct {
	// The name of the configuration group
	// example: backend-group
	Name string `json:"name"`

	// The version of the configuration group
	// example: v1
	Version string `json:"version"`

	// List of configurations to include in this group
	// minItems: 1
	ConfigurationList []*ConfigurationGroupConfigurationDto `json:"configuration_list"`
}
