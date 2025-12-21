package handlers

import "github.com/anjaobradovic/ars-sit-2025/model"

// -------------------- CONFIGS --------------------

// swagger:parameters createConfiguration
type createConfigurationParams struct {
	// in: body
	// required: true
	Body model.Config `json:"body"`
}

// swagger:parameters getConfigurationByNameAndVersion deleteConfigurationByNameAndVersion
type configPathParams struct {
	// in: path
	// required: true
	Name string `json:"name"`

	// in: path
	// required: true
	Version string `json:"version"`
}

// -------------------- GROUPS --------------------

// swagger:parameters createGroup
type createGroupParams struct {
	// in: body
	// required: true
	Body model.ConfigurationGroup `json:"body"`
}

// swagger:parameters getGroup deleteGroup addConfig removeConfig getConfigsByLabels
type groupPathParams struct {
	// in: path
	// required: true
	Name string `json:"name"`

	// in: path
	// required: true
	Version string `json:"version"`
}

// swagger:parameters addConfig
type addConfigParams struct {
	groupPathParams

	// in: body
	// required: true
	Body model.LabeledConfiguration `json:"body"`
}

// swagger:parameters removeConfig
type removeConfigParams struct {
	groupPathParams

	// in: body
	// required: true
	Body struct {
		// ID of labeled configuration to remove
		ConfigID string `json:"configId"`
	} `json:"body"`
}

// swagger:parameters getConfigsByLabels
type getConfigsByLabelsParams struct {
	groupPathParams

	// Labels filter in format key:value;key2:value2 (all must match). For example env:prod;region:eu
	// in: query
	// required: false
	Labels string `json:"labels"`
}
