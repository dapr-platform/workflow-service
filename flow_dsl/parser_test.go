package flow_dsl

import (
	"os"
	"testing"
)

func TestJsonParser_Parse(t *testing.T) {
	jsonBytes, err := os.ReadFile("test.json")
	if err != nil {
		t.Error(err.Error())
	}
	parser := &JsonParser{}
	workflow, err := parser.Parse(string(jsonBytes))
	if err != nil {
		t.Error(err.Error())
	}
	if len(workflow.Statements) == 0 {
		t.Error("parse error")
	}

}
