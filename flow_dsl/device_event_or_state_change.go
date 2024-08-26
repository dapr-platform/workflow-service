package flow_dsl

import (
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.temporal.io/sdk/activity"
	"strings"
)

type DeviceEventOrStateChangeParam struct {
	DataType    string `json:"dataType"`
	Product     string `json:"product"`
	ProductName string `json:"productName"`
	MetaType    string `json:"meta_type"`
	Key         string `json:"key"`
	Op          string `json:"op"`
	Value       string `json:"value"`
	Value2      string `json:"value2"`
}

var OUTPUT_PORT = "DeviceEventOrStateChange@output"

// 根据device singal 处理单个device数据, 发出事件或状态
func (d *FlowDslActivities) DeviceEventOrStateChange(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {
	if params.Type == PARAM_TYPE_TIME_INTERVAL {
		result.Valid = true
		result.Data = params.WorkflowPayload
		result.Type = params.Type
		return
	}
	logger := activity.GetLogger(ctx)
	if params.Type == PARAM_TYPE_DEVICE_OBJECT {

		deviceWithTags := params.WorkflowPayload.Payload.(map[string]any)
		property := params.Properties
		result.DataType = RESULT_DATA_TYPE_OBJECT
		metaType := property["meta_type"].(string)

		result.Data = &WorkFlowPayload{
			Payload:    deviceWithTags,
			MatchKey:   property["key"].(string),
			MatchValue: property["value"],
			MatchPort:  params.OutgoingBusiness[OUTPUT_PORT],
		}
		if metaType == "events" {
			result.Type = PARAM_TYPE_EVENT
		} else {
			result.Type = PARAM_TYPE_EVENT_AND_STATE
		}

		switch property["dataType"] {
		case "tag":
			tags := property["tags"].(string)
			tagArr := strings.Split(tags, ",")
			deviceTag := deviceWithTags["tags"].([]string)

			for _, tag := range tagArr {
				for _, dtag := range deviceTag {
					if tag == dtag {
						match, err1 := CheckDeviceDataMatchProperty(property, deviceWithTags)
						if err1 != nil {
							logger.Error(err1.Error())
							err = errors.Wrap(err1, "CheckDeviceDataMatchProperty")
						} else {
							logger.Debug(dtag+" match:", match)
							result.Valid = match
						}

						return
					}
				}
			}
		case "device":
			devices := property["devices"].([]any)
			for _, d := range devices {
				device := cast.ToStringMap(d)
				if device["id"] == deviceWithTags["id"] {
					match, err1 := CheckDeviceDataMatchProperty(property, deviceWithTags)
					if err1 != nil {
						logger.Error(err1.Error())
						err = errors.Wrap(err1, "CheckDeviceDataMatchProperty")
						result.Msg = err1.Error()
					} else {
						logger.Debug("device "+device["id"].(string)+" match:", match)
						result.Valid = match

					}
					return
				}
			}
		case "product":
			productId := property["product"]
			if productId == "" {
				logger.Error("productId is empty")
				return
			}
			if deviceWithTags["product_id"] == productId {
				match, err1 := CheckDeviceDataMatchProperty(property, deviceWithTags)
				if err1 != nil {
					logger.Error(err1.Error())
					err = errors.Wrap(err1, "CheckDeviceDataMatchProperty")
				} else {
					logger.Debug("product match:", match)
					result.Valid = match
				}

			}

		}

	}

	return
}
