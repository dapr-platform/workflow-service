package flow_dsl

import (
	"context"
	"encoding/json"
)

func (d *FlowDslActivities) EventDebug(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {
	paramsJson, _ := json.Marshal(params)
	result.Msg = string(paramsJson)
	return
}
