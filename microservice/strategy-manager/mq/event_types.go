package mq

type BasicEvent struct {
	EventId    string `json:"event_id"`
	OccurredAt string `json:"occurred_at"`
}

type ResourceEventTypes string

const (
	ResourceCreate ResourceEventTypes = "resource_create"
	ResourceDelete ResourceEventTypes = "resource_delete"
	ResourceUpdate ResourceEventTypes = "resource_update"
)

type TargetResourceTypes string

const (
	Strategy  TargetResourceTypes = "strategy"
	Namespace TargetResourceTypes = "namespace"
)

type ResourceEvent struct {
	BasicEvent         `json:",inline"`
	ResourceEventType  ResourceEventTypes     `json:"event_type"`
	UserId             string                 `json:"user_id,omitempty"`
	TargetResourceType TargetResourceTypes    `json:"resource_type"`
	ResourceEventInfo  map[string]string      `json:"resource_event_info,omitempty"`
	ResourceUpdateInfo map[string]interface{} `json:"update_info,omitempty"`
}

// type CreateEvent struct {
// 	ResourceEvent  `json:",inline"`
// 	UserId         string               `json:"user_id"`
// 	TargetResource string               `json:"resource_type"`
// 	CreateInfo     ResourceEventInfoMap `json:"create_info,omitempty"`
// }

// type DeleteEvent struct {
// 	BasicEvent     `json:",inline"`
// 	UserId         string            `json:"user_id"`
// 	TargetResource string            `json:"resource_type"`
// 	DeleteInfo     map[string]string `json:"delete_info,omitempty"`
// }

// type UpdateEvent struct {
// 	BasicEvent     `json:",inline"`
// 	UserId         string            `json:"user_id"`
// 	TargetResource string            `json:"resource_type"`
// 	UpdateInfo     map[string]string `json:"update_info,omitempty"`
// }
