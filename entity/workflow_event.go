package entity

type WorkflowEvent struct {
	Type         string `json:"type"`
	WorkflowName string `json:"workflow_name"`
}
type DeviceActionEvent struct {
	WorkflowEvent
	PropertyIdentifier string   `json:"property_identifier"`
	Value              any      `json:"value"`
	MatchTags          []string `json:"match_tags"`
	DataType           string   `json:"data_type"`
	SourceDevice       any      `json:"source_device"`
	DeviceIds          []string `json:"device_ids"`
	Tags               []string `json:"tags"`
	ProductId          string   `json:"product_id"`
}
