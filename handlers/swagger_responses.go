package handlers

import "github.com/anjaobradovic/ars-sit-2025/model"

// Labeled configurations response
// swagger:response labeledConfigurationsResponse
type labeledConfigurationsResponse struct {
	// in:body
	Body []*model.LabeledConfiguration
}
