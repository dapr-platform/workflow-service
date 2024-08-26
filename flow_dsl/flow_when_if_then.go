package flow_dsl

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.temporal.io/sdk/activity"
	"strings"
	"time"
	"workflow-service/redis_op"
)

type FlowWhenIfThenStatus struct {
	WhenStatus      bool
	WhenPayload     *WorkFlowPayload
	ConditionStatus bool
	OccuredTime     int64
}

var whenPort = "FlowWhenIfThen@when"
var conditionPort = "FlowWhenIfThen@condition"
var thenPort = "FlowWhenIfThen@then"
var elsePort = "FlowWhenIfThen@else"
var expiredSeconds = 5

func (d *FlowDslActivities) FlowWhenIfThen(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {
	logger := activity.GetLogger(ctx)
	workflowId := activity.GetInfo(ctx).WorkflowExecution.ID
	workflowName := strings.Split(workflowId, "_")[0]
	redisKey := GetRedisPrefix(workflowName) + "FlowWhenIfThen:" + params.CurrentId
	if params.Type == PARAM_TYPE_EVENT || params.Type == PARAM_TYPE_EVENT_AND_STATE {

		status, err1 := redis_op.GetRedisVal[FlowWhenIfThenStatus](ctx, redisKey)
		if err1 == redis.Nil {
			status = FlowWhenIfThenStatus{
				WhenStatus:      false,
				ConditionStatus: false,
				OccuredTime:     time.Now().Unix(),
			}
		}
		for k, _ := range params.WorkflowPayload.MatchPort {
			if k == whenPort {
				status.WhenStatus = true
				status.WhenPayload = params.WorkflowPayload
			}
			if k == conditionPort {
				status.ConditionStatus = true
			}
		}
		err2 := redis_op.SetRedisVal(ctx, redisKey, status, expiredSeconds*2)
		if err2 != nil {
			logger.Error("FlowWhenIfThen SetRedisVal err:", err2)
			return
		}
		return
	}
	if params.Type == PARAM_TYPE_TIME_INTERVAL { //condition的事件可能为false，即不会触发到这里，所以根据超时判断。
		status, err3 := redis_op.GetRedisVal[FlowWhenIfThenStatus](ctx, redisKey)
		if err3 == redis.Nil {
			return
		}
		if status.OccuredTime+int64(expiredSeconds) > time.Now().Unix() {
			if status.WhenStatus && status.ConditionStatus {
				result.Valid = true
				status.WhenPayload.MatchPort = params.OutgoingBusiness[thenPort]
				result.Data = status.WhenPayload
				result.Type = PARAM_TYPE_EVENT_AND_STATE
				result.DataType = RESULT_DATA_TYPE_OBJECT

			}
			if status.WhenStatus && !status.ConditionStatus {
				result.Valid = true
				status.WhenPayload.MatchPort = params.OutgoingBusiness[elsePort]
				result.Data = status.WhenPayload
				result.Type = PARAM_TYPE_EVENT_AND_STATE
				result.DataType = RESULT_DATA_TYPE_OBJECT
			}
		}

	}

	return
}
