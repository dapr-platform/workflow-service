package dataquery

import (
	"context"
	"encoding/json"
	"github.com/dapr-platform/common"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"strings"
	"workflow-service/entity"
)

var STATE_DEVICE_MIRROR_KEY = "device_mirror:"

func QueryDeviceStatesByProductIdAndTags(ctx context.Context, productId string, tags []string, property string) (result []map[string]any, err error) {
	result = make([]map[string]any, 0)
	tagArr := strings.Join(tags, ",")
	items, err := queryDeviceByTagsAndProductIdAndId(ctx, tagArr, "", productId, "enabled=1")
	if err != nil {
		return
	}
	for _, item := range items {

		d, err1 := queryDeviceMirrorProperty(ctx, cast.ToString(item["identifier"]), property)
		if err1 != nil {
			err = errors.Wrap(err1, "QueryDeviceStatesByProductIdAndTags invoke method failed")
			return
		}
		if d != nil {
			d["tags"] = item["tags"]
			d["id"] = item["id"]
			result = append(result, d)
		}

	}
	return
}

func QueryDeviceStatesByDeviceIds(ctx context.Context, deviceIds []string, property string) (result []map[string]any, err error) {
	result = make([]map[string]any, 0)
	for _, deviceId := range deviceIds {
		items, err1 := queryDeviceByTagsAndProductIdAndId(ctx, "", deviceId, "", "enabled=1")
		if err1 != nil {
			err = errors.Wrap(err1, "QueryDeviceStatesByDeviceIds invoke method failed")
			return
		}
		for _, item := range items {
			d, err1 := queryDeviceMirrorProperty(ctx, cast.ToString(item["identifier"]), property)
			if err1 != nil {
				err = errors.Wrap(err1, "QueryDeviceStatesByProductIdAndTags invoke method failed")
				return
			}
			if d != nil {
				d["tags"] = item["tags"]
				d["id"] = item["id"]
				result = append(result, d)
			}

		}

	}
	return
}
func QueryDeviceStatesByProductId(ctx context.Context, productId string, property string) (result []map[string]any, err error) {
	result = make([]map[string]any, 0)
	items, err1 := queryDeviceByTagsAndProductIdAndId(ctx, "", "", productId, "enabled=1")
	if err1 != nil {
		err = errors.Wrap(err1, "QueryDeviceStatesByDeviceIds invoke method failed")
		return
	}
	for _, item := range items {
		d, err1 := queryDeviceMirrorProperty(ctx, cast.ToString(item["identifier"]), property)
		if err1 != nil {
			err = errors.Wrap(err1, "QueryDeviceStatesByProductIdAndTags invoke method failed")
			return
		}
		if d != nil {
			d["tags"] = item["tags"]
			d["id"] = item["id"]
			result = append(result, d)
		} else {
			common.Logger.Debug("QueryDeviceStatesByProductId " + productId + " " + cast.ToString(item["identifier"]) + " device mirror is empty")
		}
	}

	return
}

func queryDeviceByTagsAndProductIdAndId(ctx context.Context, tags string, id, productId string, whereString string) (data []map[string]any, err error) {
	selectStr := "id,identifier,enabled,tags"
	fromStr := "v_device_with_tag"
	whereStr := ""
	if whereString != "" {
		whereStr = whereString + " and "
	}

	if tags != "" {
		arr := strings.Split(tags, ",")
		for _, s := range arr {
			if s != "" {
				whereStr += " '" + s + "'=any(tags)" + " and"
			}

		}
	}

	if id != "" {
		whereStr += " id='" + id + "'" + " and"
	}

	if productId != "" {
		whereStr += " product_id='" + productId + "' and"
	}

	if whereStr != "" {
		whereStr = whereStr[:strings.LastIndex(whereStr, " and")]
	} else {
		whereStr = "1=1"
	}

	data, err = common.CustomSql[map[string]any](ctx, common.GetDaprClient(), selectStr, fromStr, whereStr)
	if err != nil {
		err = errors.Wrap(err, "select data error")
		return
	}

	return
}

func queryDeviceMirrorProperty(ctx context.Context, deviceIdentifier string, property string) (data map[string]any, err error) {

	buf, err := common.GetInStateStore(ctx, common.GetDaprClient(), common.GLOBAL_STATESTOR_NAME, STATE_DEVICE_MIRROR_KEY+deviceIdentifier)
	if err != nil {
		common.Logger.Error(err.Error())
		err = errors.Wrap(err, "GetDeviceMirror")
	} else {
		if len(buf) == 0 {
			common.Logger.Error(deviceIdentifier + " GetDeviceMirror from StateStore: empty")
			return
		} else {
			deviceMirror := new(entity.DeviceMirror)
			err = json.Unmarshal(buf, deviceMirror)
			if err != nil {
				common.Logger.Error(err.Error())
				err = errors.Wrap(err, "GetDeviceMirror")
				return
			}
			data = make(map[string]any)
			data["identifier"] = deviceIdentifier
			val, exist := deviceMirror.State.Desired[property]
			if exist {
				data[property] = val
			} else {
				val, exist = deviceMirror.State.Reported[property]
				if exist {
					data[property] = val
				}
			}
			return
		}

	}

	return
}
