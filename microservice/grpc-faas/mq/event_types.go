// package mq

// type EventTypes string

// const (
// 	ResourceCreate EventTypes = "resource_create"
// 	ResourceDelete EventTypes = "resource_delete"
// 	ResourceUpdate EventTypes = "resource_update"
// )

// type BasicEvent struct {
// 	EventId    string     `json:"event_id"`
// 	EventType  EventTypes `json:"event_type"`
// 	OccurredAt string     `json:"occurred_at"`
// }

// type CreateEvent struct {
// 	BasicEvent     `json:",inline"`
// 	UserId         string `json:"user_id"`
// 	TargetResource string `json:"resource_type"`
// 	StrategyName   string `json:"strategy_name,omitempty"`
// }

// type DeleteEvent struct {
// 	BasicEvent     `json:",inline"`
// 	UserId         string `json:"user_id"`
// 	TargetResource string `json:"resource_type"`
// 	StrategyName   string `json:"strategy_name,omitempty"`
// }

// type UpdateEvent struct {
// 	BasicEvent     `json:",inline"`
// 	UserId         string                 `json:"user_id"`
// 	TargetResource string                 `json:"resource_type"`
// 	StrategyName   string                 `json:"strategy_name,omitempty"`
// 	UpdateInfo     map[string]interface{} `json:"update_info"`
// }

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
