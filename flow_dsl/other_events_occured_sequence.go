package flow_dsl

import (
	"context"
	"encoding/json"
	"go.temporal.io/sdk/activity"
)

func (d *FlowDslActivities) OtherEventsOccuredSequence(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	logger := activity.GetLogger(ctx)
	data, _ := json.Marshal(params)
	logger.Debug("process " + name + " param :" + string(data))
	return
}
