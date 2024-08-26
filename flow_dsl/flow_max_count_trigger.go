package flow_dsl

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"go.temporal.io/sdk/activity"
	"strings"
	"time"
	"workflow-service/redis_op"
)

var PORT_COUNT = "FlowMaxCountTrigger@count"
var PORT_ZERO = "FlowMaxCountTrigger@zero"
var PORT_OUTPUT = "FlowMaxCountTrigger@output"

func (d *FlowDslActivities) FlowMaxCountTrigger(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {
	logger := activity.GetLogger(ctx)
	logger.Debug("FlowMaxCountTrigger begin")
	workflowId := activity.GetInfo(ctx).WorkflowExecution.ID
	workflowName := strings.Split(workflowId, "_")[0]
	logger.Debug(workflowName + " FlowMaxCountTrigger " + params.CurrentId + " start")
	redisKey := GetRedisPrefix(workflowName) + "FlowMaxCountTrigger:" + params.CurrentId
	countLimit := 1
	val, exist := params.Properties["count"]

	if exist {
		countLimit = cast.ToInt(val)
	} else {
		logger.Error(workflowName + " FlowMaxCountTrigger " + "count property null")
		result.Msg += "count property null]\n"
	}
	if params.Type == PARAM_TYPE_EVENT || params.Type == PARAM_TYPE_EVENT_AND_STATE {
		countLogic := false
		for k, _ := range params.WorkflowPayload.MatchPort {
			if k == PORT_COUNT {
				countLogic = true
			}
		}
		logger.Debug(workflowName + " FlowMaxCountTrigger " + "countLogic " + cast.ToString(countLogic))
		id := cast.ToString(params.WorkflowPayload.Payload.(map[string]interface{})["id"])
		if id == "" {
			result.Msg += "id property null\n"
			logger.Debug(workflowName + " FlowMaxCountTrigger " + "id property is null")

		}
		cacheKey := redisKey + ":" + id
		if countLogic {
			countVal := 0
			cacheVal, err2 := redis_op.Rdb.Get(ctx, cacheKey).Result()
			if err2 != nil && err2 != redis.Nil {
				logger.Error(workflowName + " FlowMaxCountTrigger " + "redis get error " + err2.Error())
			}
			if cacheVal != "" {
				countVal = cast.ToInt(cacheVal) + 1
			} else {
				countVal = 1
			}
			logger.Debug(workflowName + " FlowMaxCountTrigger " + "countVal " + cast.ToString(countVal) + " countLimit " + cast.ToString(countLimit))
			if countVal <= countLimit {
				result.Valid = true
				result.Data = params.WorkflowPayload
				result.Type = params.Type
				result.Data.MatchPort = params.OutgoingBusiness[PORT_OUTPUT]

			} else {
				result.Valid = false
			}
			redis_op.Rdb.Set(ctx, cacheKey, countVal, time.Hour*24)
		} else {
			redis_op.Rdb.Del(ctx, cacheKey)
		}

		return
	}

	return
}
