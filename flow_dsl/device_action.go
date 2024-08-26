package flow_dsl

import (
	"context"
	"encoding/json"
	"github.com/spf13/cast"
	"go.temporal.io/sdk/activity"
	"strings"
	"workflow-service/entity"
	"workflow-service/eventpub"
)

func (d *FlowDslActivities) DeviceAction(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {
	logger := activity.GetLogger(ctx)
	workflowId := activity.GetInfo(ctx).WorkflowExecution.ID
	workflowName := strings.Split(workflowId, "_")[0]

	if params.Type != PARAM_TYPE_EVENT && params.Type != PARAM_TYPE_EVENT_AND_STATE {
		result.Valid = true
		result.Data = params.WorkflowPayload
		result.Type = params.Type
		result.Msg = "params.Type != PARAM_TYPE_EVENT && params.Type != PARAM_TYPE_EVENT_AND_STATE,skip"

		return
	}
	logger.Debug(workflowName + " DeviceAction start " + params.CurrentId)
	_, exist := params.WorkflowPayload.MatchPort["DeviceAction@input"]
	if !exist {
		str, _ := json.Marshal(params.WorkflowPayload.MatchPort)
		logger.Debug("DeviceAction matchPort " + string(str) + " not contains DeviceAction@input,skip")
		logger.Debug(workflowName + " DeviceAction " + params.CurrentId + "not match port end")
		return
	}

	sourceDevice := params.WorkflowPayload.Payload
	if sourceDevice != nil {
		logger.Debug("DeviceAction payloadId=" + cast.ToString(params.WorkflowPayload.Payload.(map[string]any)["identifier"]))
	}
	logger.Debug(workflowName+" DeviceAction properties = ", params.Properties)
	property := params.Properties
	deviceActionEvent := entity.DeviceActionEvent{
		DataType:     property["dataType"].(string),
		SourceDevice: sourceDevice,
	}
	if cast.ToString(property["matchTag"]) != "" {
		deviceActionEvent.MatchTags = strings.Split(property["matchTag"].(string), ",")
	} else {
		deviceActionEvent.MatchTags = make([]string, 0)
	}
	deviceActionEvent.WorkflowEvent.Type = "DeviceAction"
	deviceActionEvent.WorkflowEvent.WorkflowName = workflowName
	deviceActionEvent.PropertyIdentifier = cast.ToString(property["key"])
	deviceActionEvent.Value = property["value"]
	switch property["dataType"] {
	case "tag":
		tags := property["tags"].(string)
		tagArr := strings.Split(tags, ",")
		deviceActionEvent.Tags = tagArr

	case "device":
		devices := property["devices"].([]any)
		if devices != nil {
			deviceIds := make([]string, 0)
			for _, d := range devices {
				device := cast.ToStringMap(d)
				deviceIds = append(deviceIds, device["id"].(string))
			}
			deviceActionEvent.DeviceIds = deviceIds
		}

	case "product":
		productId := property["product"]
		if productId == "" {
			logger.Error("productId is empty")
			return
		}
		deviceActionEvent.ProductId = productId.(string)
	}
	err = eventpub.PublishInternalMessage(context.Background(), "workflow", deviceActionEvent)
	logger.Debug(workflowName + " DeviceAction " + params.CurrentId + "send msg and end")
	return
}
