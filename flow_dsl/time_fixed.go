package flow_dsl

import (
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.temporal.io/sdk/activity"
	"strings"
	"time"
	"workflow-service/holiday"
)

// 定时
func (d *FlowDslActivities) TimeFixed(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {
	if params.Type == PARAM_TYPE_TIME_INTERVAL {
		logger := activity.GetLogger(ctx)
		valid, err1 := timeFixedCheckTimeValid(params)
		if err1 != nil {
			err = errors.Wrap(err1, "TimeFixed")
			result.Msg = err1.Error()
			logger.Error(err.Error())
			return
		}
		if valid {
			result.Msg = "pass"
			result.Valid = true
			result.Data = params.WorkflowPayload
			targetPort := params.OutgoingBusiness["TimeFixed@output"]
			result.Data.MatchPort = targetPort
			result.Type = PARAM_TYPE_EVENT
		} else {
			result.Msg = "fail"
		}
	} else {
		result.Msg = "params type " + params.Type + " not process"
	}

	return
}

func timeFixedCheckTimeValid(params WorkflowParams) (ret bool, err error) {
	timeType := cast.ToString(params.Properties["timeType"])
	if timeType == "" {
		err = errors.New("timeType error")
		return
	}
	timeStr := params.Properties["time"]
	if time.Now().Format("15:04:05") == timeStr {
		switch timeType {
		case "1": //每天
			ret = true
		case "2": //工作日
			if holiday.IsWorkday(time.Now()) {
				ret = true
			}
		case "3": //节假日
			if holiday.IsHoliday(time.Now()) {
				ret = true
			}
		case "4": //自定义
			week := cast.ToString(params.Properties["week"])
			if week != "" {
				weekArr := strings.Split(week, ",")
				weekNow := time.Now().Weekday().String()
				for _, v := range weekArr {
					if v == weekNow {
						ret = true
					}
				}
			}

		}
	}
	return
}
