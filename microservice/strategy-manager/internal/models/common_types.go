package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StrategyDoc struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID          string             `json:"user_id" bson:"user_id"`
	StrategyName    string             `json:"strategy_name" bson:"strategy_name"`
	CreatedAt       string             `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       string             `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CrontabSchedule string             `json:"schedule,omitempty" bson:"schedule,omitempty"`
	CpuRequest      string             `json:"cpu_request,omitempty" bson:"cpu_request,omitempty"`
	MemRequest      string             `json:"mem_request,omitempty" bson:"mem_request,omitempty"`
	CpuLimit        string             `json:"cpu_limit,omitempty" bson:"cpu_limit,omitempty"`
	MemLimit        string             `json:"mem_limit,omitempty" bson:"mem_limit,omitempty"`
}

type StrategyDocJson struct {
	ID              primitive.ObjectID `json:"_id,omitempty"`
	UserID          string             `json:"user_id"`
	StrategyName    string             `json:"strategy_name"`
	CreatedAt       string             `json:"created_at,omitempty"`
	UpdatedAt       string             `json:"updated_at,omitempty"`
	CrontabSchedule string             `json:"schedule,omitempty"`
	CpuRequest      string             `json:"cpu_request,omitempty"`
	MemRequest      string             `json:"mem_request,omitempty"`
	CpuLimit        string             `json:"cpu_limit,omitempty"`
	MemLimit        string             `json:"mem_limit,omitempty"`
}
