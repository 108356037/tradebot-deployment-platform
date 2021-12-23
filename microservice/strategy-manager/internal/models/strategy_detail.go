package models

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

// given a specific user strategy, returned its current resource request as (cpu(*resource.Quantity), mem(*resource.Quantity))
func GetStrategyRequest(user, strategy string) (*resource.Quantity, *resource.Quantity) {
	strategyDoc, err := GetSingleStrategy(user, strategy)
	if err != nil {
		return nil, nil
	}

	cpuUsed := resource.MustParse(strategyDoc.CpuRequest)
	memUsed := resource.MustParse(strategyDoc.MemRequest)
	return &cpuUsed, &memUsed
}

// given a specific user strategy, returned its current resource limit as (cpu(*resource.Quantity), mem(*resource.Quantity))
func GetStrategyLimit(user, strategy string) (*resource.Quantity, *resource.Quantity) {
	strategyDoc, err := GetSingleStrategy(user, strategy)
	if err != nil {
		return nil, nil
	}

	cpuUsed := resource.MustParse(strategyDoc.CpuLimit)
	memUsed := resource.MustParse(strategyDoc.MemLimit)
	return &cpuUsed, &memUsed
}
