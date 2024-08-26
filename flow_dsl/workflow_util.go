package flow_dsl

import (
	"github.com/spf13/cast"
	"strings"
)

var UsingSecondIntervalStatmentTypes = map[string]bool{
	"TimeDelay":        true,
	"TimeFixed":        true,
	"TimeSpan":         true,
	"TimeStateKeeping": true,
	"FlowLoop":         true,
	"FlowWhenIfThen":   true,
}

func GetRedisPrefix(workflowName string) string {
	return "workflow:" + workflowName + ":"
}
func IsWorkflowUsingSecondInterval(workflow *Workflow) bool {
	for _, stat := range workflow.Statements {
		if isStatmentUsingSecondInterval(stat) {
			return true
		}

	}
	return false
}

func isStatmentUsingSecondInterval(stat *Statement) bool {
	if _, exist := UsingSecondIntervalStatmentTypes[stat.Type]; exist {
		return true
	} else {
		for _, child := range stat.NextStatements {
			if isStatmentUsingSecondInterval(child) {
				return true
			}
		}
	}
	return false
}

func ParseWorkflowDeviceFilter(workflow *Workflow) (filters []DeviceFilter) {
	filters = make([]DeviceFilter, 0)
	for _, stat := range workflow.Statements {
		if stat.Type == "DeviceEventOrStateChange" {
			tags := strings.Split(cast.ToString(stat.Properties["tags"]), ",")

			ids := make([]string, 0)
			devices := cast.ToSlice(stat.Properties["devices"])
			for _, d := range devices {
				device := cast.ToStringMap(d)
				ids = append(ids, cast.ToString(device["id"]))
			}
			filter := DeviceFilter{
				Type:      cast.ToString(stat.Properties["dataType"]),
				Tags:      tags,
				ProductId: cast.ToString(stat.Properties["product"]),
				IDs:       ids,
			}
			filters = append(filters, filter)
		}
	}
	return
}
