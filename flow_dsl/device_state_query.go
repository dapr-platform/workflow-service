package flow_dsl

import (
	"context"
	"github.com/spf13/cast"
	"go.temporal.io/sdk/activity"
	"strings"
	"workflow-service/dataquery"
)

var DeviceStateQuery_PORT_TRUE = "DeviceStateQuery@true"
var DeviceStateQuery_PORT_FALSE = "DeviceStateQuery@false"

func (d *FlowDslActivities) DeviceStateQuery(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {
	logger := activity.GetLogger(ctx)
	if params.Type != PARAM_TYPE_EVENT && params.Type != PARAM_TYPE_EVENT_AND_STATE {
		result.Valid = true
		result.Data = params.WorkflowPayload
		result.Type = params.Type
		logger.Debug("params.Type != PARAM_TYPE_EVENT || params.Type !=PARAM_TYPE_EVENT_AND_STATE, ignore " + params.Type)
		result.Msg = "params.Type != PARAM_TYPE_EVENT || params.Type !=PARAM_TYPE_EVENT_AND_STATE, ignore "
		return
	}

	logger.Debug("DeviceStateQuery properties = ", params.Properties)
	logger.Debug("DeviceStateQuery outgoingports = ", params.OutgoingBusiness)
	property := params.Properties
	resultData := make([]*WorkFlowPayload, 0)
	result.Datas = resultData
	result.Type = PARAM_TYPE_EVENT_AND_STATE
	result.DataType = RESULT_DATA_TYPE_ARRAY
	switch property["dataType"] {
	case "tag":
		tags := property["tags"].(string)
		tagArr := strings.Split(tags, ",")
		devices, err1 := dataquery.QueryDeviceStatesByProductIdAndTags(ctx, cast.ToString(property["product"]), tagArr, cast.ToString(property["key"]))
		if err1 != nil {
			logger.Error("QueryDeviceStatesByProductIdAndTags err=", err1)
			err = err1
			return
		} else {
			for _, device := range devices {
				match, err2 := CheckDeviceDataMatchProperty(property, device)
				if err2 != nil {
					logger.Error("CheckDeviceDataMatchProperty err=", err2)
					err = err2
					return
				}
				if match {

					resultData = append(resultData, &WorkFlowPayload{Payload: device, MatchPort: params.OutgoingBusiness[DeviceStateQuery_PORT_TRUE], MatchKey: property["key"].(string), MatchValue: device[property["key"].(string)]})
				} else {
					resultData = append(resultData, &WorkFlowPayload{Payload: device, MatchPort: params.OutgoingBusiness[DeviceStateQuery_PORT_FALSE], MatchKey: property["key"].(string), MatchValue: device[property["key"].(string)]})
				}
			}
		}

	case "device":
		devices := property["devices"].([]any)
		if devices != nil {
			deviceIds := make([]string, 0)
			for _, device := range devices {
				d := cast.ToStringMap(device)
				deviceIds = append(deviceIds, d["id"].(string))
			}
			devices, err1 := dataquery.QueryDeviceStatesByDeviceIds(ctx, deviceIds, cast.ToString(property["key"]))
			if err1 != nil {
				logger.Error("QueryDeviceStatesByDeviceIds err=", err1)
				err = err1
				return
			}
			for _, device := range devices {
				match, err2 := CheckDeviceDataMatchProperty(property, device)
				if err2 != nil {
					logger.Error("CheckDeviceDataMatchProperty err=", err2)
					err = err2
					continue
				}
				if match {
					resultData = append(resultData, &WorkFlowPayload{Payload: device, MatchPort: params.OutgoingBusiness[DeviceStateQuery_PORT_TRUE], MatchKey: property["key"].(string), MatchValue: device[property["key"].(string)]})
				} else {
					resultData = append(resultData, &WorkFlowPayload{Payload: device, MatchPort: params.OutgoingBusiness[DeviceStateQuery_PORT_FALSE], MatchKey: property["key"].(string), MatchValue: device[property["key"].(string)]})
				}
			}
		}

	case "product":
		productId := property["product"]
		if productId == "" {
			logger.Error("productId is empty")
			return
		}
		logger.Debug("productId=" + cast.ToString(productId))
		devices, err1 := dataquery.QueryDeviceStatesByProductId(ctx, cast.ToString(productId), cast.ToString(property["key"]))
		if err1 != nil {
			logger.Error("QueryDeviceStatesByDeviceIds err=", err1)
			err = err1
			return
		}
		logger.Debug("device length=" + cast.ToString(len(devices)))
		for _, device := range devices {
			match, err2 := CheckDeviceDataMatchProperty(property, device)
			if err2 != nil {
				logger.Error("CheckDeviceDataMatchProperty err=", err2)
				continue
			}
			if match {
				logger.Debug("match device = ", device)
				logger.Debug("outgoings = ", params.OutgoingBusiness)
				resultData = append(resultData, &WorkFlowPayload{Payload: device, MatchPort: params.OutgoingBusiness[DeviceStateQuery_PORT_TRUE], MatchKey: property["key"].(string), MatchValue: device[property["key"].(string)]})
			} else {
				logger.Debug("not match device = ", device)
				resultData = append(resultData, &WorkFlowPayload{Payload: device, MatchPort: params.OutgoingBusiness[DeviceStateQuery_PORT_FALSE], MatchKey: property["key"].(string), MatchValue: device[property["key"].(string)]})
			}
		}
	}
	result.Datas = resultData
	if len(resultData) > 0 {
		result.Valid = true
	} else {
		result.Valid = false
	}

	return
}
