package flow_dsl

import (
	"context"
	"go.temporal.io/sdk/activity"
	"strings"
	"workflow-service/eventpub"
)

func (d *FlowDslActivities) DeviceDataCollection(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {
	if params.Type == PARAM_TYPE_EVENT {
		logger := activity.GetLogger(ctx)
		workflowId := activity.GetInfo(ctx).WorkflowExecution.ID
		workflowName := strings.Split(workflowId, "_")[0]
		err = eventpub.PublishInternalMessage(ctx, "device_collection", params.Properties)
		result.Msg = "publish device_collection message, and finished"
		logger.Debug(workflowName + " DeviceDataCollection " + params.CurrentId + " send msg end")
	}

	return
}
