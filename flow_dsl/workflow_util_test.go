package flow_dsl

import (
	"os"
	"testing"
)

func TestFilter(t *testing.T) {
	jsonBytes, err := os.ReadFile("test.json")
	if err != nil {
		t.Error(err)
	}

	parser := &JsonParser{}
	workflow, err := parser.Parse(string(jsonBytes))
	if err != nil {
		t.Error(err)
	}
	filters := ParseWorkflowDeviceFilter(&workflow)
	t.Log(len(filters))
	t.Log(filters)
	device := make(map[string]interface{})
	device["id"] = "8f843d79cff85d83a1547f5c6a42f9af"
	for _, filter := range filters {
		t.Log(filter.Filter(device))
	}
}

func TestIsWorkflowUsingSecondInterval(t *testing.T) {
	jsonBytes, err := os.ReadFile("test1.json")
	if err != nil {
		t.Error(err)
	}

	parser := &JsonParser{}
	workflow, err := parser.Parse(string(jsonBytes))
	if err != nil {
		t.Error(err)
	}

	t.Log(IsWorkflowUsingSecondInterval(&workflow))

}
