package flow_dsl

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) TestRun() {
	env := s.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(FlowDslWorkflow)
	env.RegisterActivity(&FlowDslActivities{})

	/*
		env.RegisterDelayedCallback(func() {
			device := map[string]any{"deviceId": "device1"}
			env.SignalWorkflow("Device", device)
		}, time.Second)
		env.RegisterDelayedCallback(func() {
			device := map[string]any{"deviceId": "device2"}
			env.SignalWorkflow("Device", device)
		}, time.Second*2)

	*/
	jsonBytes, err := os.ReadFile("dc.json")
	if err != nil {
		s.Error(err)
	}

	parser := &JsonParser{}
	workflow, err := parser.Parse(string(jsonBytes))
	if err != nil {
		s.Error(err)
	}
	workflow.Running = true
	workflow.UseLocalActivity = true
	go env.ExecuteWorkflow(FlowDslWorkflow, &workflow)
	fmt.Println("after execute workflow")
	//s.True(env.IsWorkflowCompleted())
	//s.NoError(env.GetWorkflowError())

	time.Sleep(time.Second * 60)
	fmt.Println("after sleep")
}

func TestStr(t *testing.T) {
	t1 := time.Now().Format("15:04:05")
	t.Log(t1 > "08:00:00")
	t.Log(t1 < "18:00:00")

}
