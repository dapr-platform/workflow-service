package flow_dsl

import "testing"

func TestTimespan(t *testing.T) {
	params := WorkflowParams{
		Properties: make(map[string]interface{}),
	}
	params.Properties["time"] = []string{"08:00:00", "20:00:00"}
	params.Properties["timeType"] = 4
	params.Properties["week"] = "1,2,3"
	ret, err := timeSpanCheckTimeValid(params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ret)
}
