package flow_dsl

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"go.temporal.io/sdk/activity"
	"strings"
	"time"
	"workflow-service/redis_op"
)

var START_PORT = "FlowLoop@begin"
var END_PORT = "FlowLoop@end"

func (d *FlowDslActivities) FlowLoop(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {

	if params.Type == PARAM_TYPE_STATE {
		result.Valid = true
		result.Data = params.WorkflowPayload
		result.Type = params.Type
		return
	}

	logger := activity.GetLogger(ctx)
	result.Data = params.WorkflowPayload
	matchPorts := result.Data.MatchPort
	targetPort := params.OutgoingBusiness["FlowLoop@output"]
	result.Data.MatchPort = targetPort
	result.Type = PARAM_TYPE_EVENT
	_, exist := matchPorts[END_PORT]
	if exist {
		result.Valid = false
		logger.Debug("flow loop end")
		return
	}
	_, exist1 := matchPorts[START_PORT]
	workflowId := activity.GetInfo(ctx).WorkflowExecution.ID
	workflowName := strings.Split(workflowId, "_")[0]
	redisKey := GetRedisPrefix(workflowName) + params.CurrentId
	if exist1 && (params.Type == PARAM_TYPE_TIME_INTERVAL || params.Type == PARAM_TYPE_EVENT) {
		var cacheData SecondInterval
		err1 := redis_op.Rdb.Get(ctx, redisKey).Scan(&cacheData)
		delaySecond := 60
		val, exists := params.Properties["delay"]
		result.Msg = "delaySecond=" + cast.ToString(delaySecond) + " val=" + cast.ToString(val)
		if exists {
			delaySecond = cast.ToInt(val)
			//logger.Debug("read property delaySecond:" + cast.ToString(delaySecond))
		}
		switch {
		case err1 == redis.Nil:
			entity := &SecondInterval{
				Second: time.Now().Unix(),
			}
			_, err2 := redis_op.Rdb.Set(ctx, redisKey, entity, time.Duration(delaySecond*2)*time.Second).Result()
			if err2 != nil {
				err = errors.Wrap(err2, "set cache data error")
				return result, err
			}
			logger.Debug("has no cache data, create new one")
			result.Valid = false
			return
		case err1 != nil:
			err = errors.Wrap(err1, "get cache data error")
			return result, err

		}
		buf, _ := json.Marshal(params.WorkflowPayload.Payload)
		var eventInterval SecondInterval
		err = json.Unmarshal(buf, &eventInterval)
		if err != nil {
			err = errors.Wrap(err, "unmarshal error")
			logger.Error("unmarshal error")
			return result, err
		}
		cacheBuf, _ := json.Marshal(cacheData)
		logger.Debug("SecondInterval " + string(buf) + " cacheData=" + string(cacheBuf))
		if eventInterval.Second-cacheData.Second >= int64(delaySecond) {
			logger.Debug("has cache data, and interval is over")
			_, err2 := redis_op.Rdb.Set(ctx, redisKey, &eventInterval, time.Duration(delaySecond*2)*time.Second).Result()
			if err2 != nil {
				err = errors.Wrap(err2, "set cache data error")
				return result, err
			}
			result.Valid = true
		}

	} else {
		s, _ := json.Marshal(matchPorts)
		logger.Debug("matchPorts=" + string(s) + " event type=" + params.Type)
	}
	return
}
