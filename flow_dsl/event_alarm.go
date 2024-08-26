package flow_dsl

import (
	"context"
	"github.com/dapr-platform/common"
	"github.com/spf13/cast"
	"go.temporal.io/sdk/activity"
	"strings"
	"time"
	"workflow-service/eventpub"
)

func (d *FlowDslActivities) EventAlarm(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {
	workflowId := activity.GetInfo(ctx).WorkflowExecution.ID
	workflowName := strings.Split(workflowId, "_")[0]

	if params.Type == PARAM_TYPE_EVENT || params.Type == PARAM_TYPE_EVENT_AND_STATE {
		logger := activity.GetLogger(ctx)
		logger.Debug("EventAlarm Process")
		device := cast.ToStringMap(params.WorkflowPayload.Payload)
		id := cast.ToString(device["id"])
		if id == "" {
			result.Msg += "id is empty"
			logger.Error(result.Msg)
			return
		}
		identifier := cast.ToString(device["identifier"])
		level := cast.ToInt(params.Properties["alarmLevel"])
		if level == 0 {
			level = 4
			result.Msg += "alarmLevel is 0, set to 4"
			logger.Error(result.Msg)
		}
		alarmTitle := cast.ToString(params.Properties["alarmTitle"])
		if alarmTitle == "" {
			alarmTitle = workflowName
		} else {
			alarmTitle = replaceVariable(alarmTitle, cast.ToStringMap(params.WorkflowPayload.Payload))
		}
		alarmText := cast.ToString(params.Properties["alarmContent"])
		if alarmText == "" {
			alarmText = workflowName
		} else {
			alarmTitle = replaceVariable(alarmTitle, cast.ToStringMap(params.WorkflowPayload.Payload))
		}
		eventpub.ConstructAndSendEvent(ctx, 1, alarmTitle, alarmText, common.EventStatusActive, int32(level), time.Now(), id, identifier, "")
	}

	return
}

func replaceVariable(str string, device map[string]any) string {
	for k, v := range device {
		str = strings.ReplaceAll(str, "${"+k+"}", cast.ToString(v))
	}
	return str
}
