package flow_dsl

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"go.temporal.io/sdk/activity"
	"strings"
	"time"
	"workflow-service/redis_op"
)

// 延时。上游传过来的事件，延时后，再想下游发出。
func (d *FlowDslActivities) TimeDelay(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {
	logger := activity.GetLogger(ctx)
	logger.Debug("TimeDelay begin")
	workflowId := activity.GetInfo(ctx).WorkflowExecution.ID
	workflowName := strings.Split(workflowId, "_")[0]
	zsetKey := GetRedisPrefix(workflowName) + "TimeDelay:" + params.CurrentId
	delaySeconds := 60
	delayVal, exist := params.Properties["delay"]

	if exist {
		delaySeconds = cast.ToInt(delayVal)
	} else {
		logger.Error("delay property null")
		result.Msg += "delay property null]\n"
	}
	logger.Debug("params.Type=" + params.Type)
	//必须带有id, 通过redis缓存payload和时间。
	if params.Type == PARAM_TYPE_EVENT || params.Type == PARAM_TYPE_EVENT_AND_STATE {
		id := cast.ToString(params.WorkflowPayload.Payload.(map[string]interface{})["id"])
		if id == "" {
			result.Msg += "id property null\n"
			logger.Debug("id property is null")
			return
		}
		objKey := GetRedisPrefix(workflowName) + "timeDelay:" + params.CurrentId + ":obj:" + id + ":" + cast.ToString(time.Now().Unix())
		logger.Debug("zset key=" + zsetKey + " objkey=" + objKey + " score=" + cast.ToString(time.Now().Unix()))
		zval, err1 := redis_op.Rdb.ZAdd(ctx, zsetKey, redis.Z{cast.ToFloat64(time.Now().Unix()), objKey}).Result()
		if err1 != nil {
			logger.Error("timedelay zadd error " + err1.Error())
			result.Msg += "timedelay zadd error " + err1.Error() + "\n"
			return
		}
		logger.Debug("zval=" + cast.ToString(zval))
		valStr, _ := json.Marshal(params.WorkflowPayload)
		_, err2 := redis_op.Rdb.Set(ctx, objKey, valStr, 0).Result()
		if err2 != nil {
			logger.Error("timedelay set error " + err2.Error())
			result.Msg += "timedelay set error " + err2.Error() + "\n"
			return
		}
		return
	}

	if params.Type == PARAM_TYPE_TIME_INTERVAL {

		maxScore := cast.ToFloat64(time.Now().Unix()) - float64(delaySeconds)
		keys, err1 := redis_op.Rdb.ZRangeByScore(ctx, zsetKey, &redis.ZRangeBy{Min: cast.ToString(0), Max: cast.ToString(maxScore)}).Result()
		if err1 != nil {
			logger.Error("timedelay zrange error " + err1.Error())
			result.Msg += "timedelay zrange error " + err1.Error() + "\n"
			return
		}
		if len(keys) == 0 {
			logger.Debug("timedelay zrange empty")
			result.Msg += "timedelay zrange empty\n"
			return
		}
		resultData := make([]*WorkFlowPayload, 0)
		result.Valid = true

		result.Type = PARAM_TYPE_EVENT_AND_STATE
		result.DataType = RESULT_DATA_TYPE_ARRAY

		for _, key := range keys {
			val, err2 := redis_op.Rdb.Get(ctx, key).Result()
			if err2 == redis.Nil {
				logger.Error("timedelay get nil " + err2.Error())
				result.Msg += "timedelay get " + key + " nil " + err2.Error() + "\n"
				continue
			}
			var payload WorkFlowPayload
			err3 := json.Unmarshal([]byte(val), &payload)
			if err3 != nil {
				logger.Error("timedelay json unmarshal error " + err3.Error())
				result.Msg += "timedelay json unmarshal error " + err3.Error() + "\n"
				return
			}
			redis_op.Rdb.Del(ctx, key)
			payload.MatchPort = params.OutgoingBusiness["TimeDelay@output"]
			resultData = append(resultData, &payload)
		}
		result.Datas = resultData
		zremVal, err1 := redis_op.Rdb.ZRemRangeByScore(ctx, zsetKey, "0", cast.ToString(maxScore)).Result()
		if err1 != nil {
			logger.Error("timedelay zremrangebyscore error " + err1.Error())
			result.Msg += "timedelay zremrangebyscore error " + err1.Error() + "\n"
		}
		if zremVal > 0 {
			logger.Debug("zremrangebyscore success " + cast.ToString(zremVal))
		}
	}
	return
}
