package flow_dsl

import (
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
	"workflow-service/redis_op"
)

var NEED_TO_CONTINUE_AS_NEW_NUM = 20000

type (
	Workflow struct {
		workflow.Context  `json:"-"`
		RuleId            string         `json:"rule_id"`
		UseLocalActivity  bool           `json:"use_local_activity"`
		Params            WorkflowParams `json:"params"`
		Running           bool           `json:"running"`
		NeedContinueAsNew bool           `json:"need_continue_as_new"`
		Statements        []*Statement   `json:"statements"`
		StopChan          chan bool      `json:"-"`
	}

	Statement struct {
		Id               string
		Type             string //类型
		Properties       Properties
		IncomingBusiness map[string]map[string]string //key incoming business,value  为 self business
		OutgoingBusiness map[string]map[string]string //key self business,value  为 target business
		Activity         *ActivityInvocation
		NextStatements   []*Statement
	}
	ActivityInvocation struct {
		Id     string
		Name   string
		param  WorkflowParams
		Result WorkflowResult
	}

	WorkflowParams struct {
		WorkflowPayload  *WorkFlowPayload
		CurrentId        string
		Properties       Properties
		IncomingBusiness map[string]map[string]string //key 为incoming node id, value 为business
		OutgoingBusiness map[string]map[string]string
		//	SignalType       string //有时间流和device对象流. 不同的卡片，需要处理不同的流，对于不需要处理的，就直接返回true。处理后，根据需要转成对应的流
		Type string //event or state
		//PreNodeResult    map[string]WorkflowResult
	}
	WorkFlowPayload struct {
		Payload    any
		MatchPort  map[string]string //default output
		MatchKey   string            `json:"match_key"`
		MatchValue any               `json:"match_value"`
	}
	WorkflowResult struct {
		Id       string
		Valid    bool   //节点结果是否有效
		Type     string //结果类型 event or state
		DataType string //节点结果类型,Object or Array
		Data     *WorkFlowPayload
		Datas    []*WorkFlowPayload
		Msg      string //for debug
	}

	Properties map[string]any

	executable interface {
		execute(ctx workflow.Context, useLocalActivity bool, bindings WorkflowParams) (result WorkflowResult, err error)
	}
)

func NewWorkflowParams() WorkflowParams {
	return WorkflowParams{
		IncomingBusiness: make(map[string]map[string]string, 0),
		OutgoingBusiness: make(map[string]map[string]string, 0),
		//PreNodeResult:    make(map[string]WorkflowResult, 0),
	}
}

// 每个workflow一个时间发生器
func (w *Workflow) StartSecondInterval(ctx workflow.Context) {
	//ticker := time.NewTimer(time.Second)
	//for range ticker.C {
	//	if !w.Running {
	//		ticker.Stop()
	//		return
	//	}
	//	workflow.Go(ctx, w.secondIntervalFunc)
	//}

	selector := workflow.NewSelector(ctx)
	timerSleepMills := int64(1000)
	redis_op.Rdb.Set(context.Background(), w.RuleId+"_second_interval", cast.ToString(time.Now().UnixMilli()), time.Second*10)
	for {

		timerFuture := workflow.NewTimer(ctx, time.Millisecond*time.Duration(timerSleepMills))
		selector.AddFuture(timerFuture, func(f workflow.Future) {

			var data = SecondInterval{
				Second: time.Now().Unix(),
			}
			params := NewWorkflowParams()
			params.Type = PARAM_TYPE_TIME_INTERVAL
			params.WorkflowPayload = &WorkFlowPayload{
				Payload: data,
			}

			for _, statement := range w.Statements {
				///logger.Debug("statement type", statement.Type, " id=", statement.Id)

				_, _ = statement.execute(ctx, w.UseLocalActivity, params)

			}
		})
		//logger.Debug("before select")
		if workflow.GetInfo(ctx).GetCurrentHistoryLength() > NEED_TO_CONTINUE_AS_NEW_NUM {
			w.NeedContinueAsNew = true
			return
		}
		selector.Select(ctx)
		//logger.Debug("after select")
		beginTime := cast.ToInt64(redis_op.Rdb.Get(context.Background(), w.RuleId+"_second_interval").Val())
		if time.Now().UnixMilli()-beginTime > 1000 {
			timerSleepMills = 2000 - (time.Now().UnixMilli() - beginTime)
		} else {
			timerSleepMills = 1000
		}
		redis_op.Rdb.Set(context.Background(), w.RuleId+"_second_interval", cast.ToString(beginTime+1000), time.Second*10)

	}

}

/*
	func (w *Workflow) startSubscribe(ctx workflow.Context) {
		logger := workflow.GetLogger(ctx)
		eventsub.RegisterDeviceEventHandler(w.RuleId, ctx, w.handleDeviceEvent)
		logger.Debug(w.RuleId + " start subscribe")
	}
*/

func (w *Workflow) Start(ctx workflow.Context) {
	//w.startSubscribe(ctx)
	logger := workflow.GetLogger(ctx)
	logger.Debug("workflow begin listen")
	selector := workflow.NewSelector(ctx)
	count := 0
	selector.AddReceive(workflow.GetSignalChannel(ctx, "Status"), func(c workflow.ReceiveChannel, more bool) {
		var statusMap map[string]bool
		c.Receive(ctx, &statusMap)
		w.Running = statusMap["status"]
		count++
		logger.Info("Status signal Received", statusMap)
	})
	selector.AddReceive(workflow.GetSignalChannel(ctx, "SecondInterval"), func(c workflow.ReceiveChannel, more bool) {
		var data SecondInterval
		c.Receive(ctx, &data)
		count++
		params := NewWorkflowParams()
		params.Type = PARAM_TYPE_TIME_INTERVAL
		params.WorkflowPayload = &WorkFlowPayload{
			Payload: data,
		}

		for _, statement := range w.Statements {
			///logger.Debug("statement type", statement.Type, " id=", statement.Id)

			_, _ = statement.execute(ctx, w.UseLocalActivity, params)

		}

		//logger.Info("Device signal Received", device)
	})

	selector.AddReceive(workflow.GetSignalChannel(ctx, "Device"), func(c workflow.ReceiveChannel, more bool) {
		var device map[string]any
		c.Receive(ctx, &device)
		count++
		for _, statement := range w.Statements {
			params := NewWorkflowParams()
			params.WorkflowPayload = &WorkFlowPayload{
				Payload: device,
			}
			params.Type = PARAM_TYPE_DEVICE_OBJECT
			_, err := statement.execute(ctx, w.UseLocalActivity, params)
			if err != nil {
				continue
			}
		}

		//logger.Info("Device signal Received", device)
	})
	for {
		if workflow.GetInfo(ctx).GetCurrentHistoryLength() > NEED_TO_CONTINUE_AS_NEW_NUM {
			logger.Debug("count=" + cast.ToString(count))
			w.NeedContinueAsNew = true
			return
		}
		//logger.Debug("before select " + cast.ToString(workflow.GetInfo(ctx).GetCurrentHistoryLength()) + " count=" + cast.ToString(count))
		selector.Select(ctx)
		//logger.Debug("after select")
	}
}

func FlowDslWorkflow(ctx workflow.Context, dslWorkflow *Workflow) ([]byte, error) {

	if dslWorkflow.UseLocalActivity {
		aol := workflow.LocalActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    time.Second,
				BackoffCoefficient: 1.0,
				MaximumInterval:    time.Minute,
				MaximumAttempts:    5,
			},
		}
		ctx = workflow.WithLocalActivityOptions(ctx, aol)
	} else {
		ao := workflow.ActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    time.Second,
				BackoffCoefficient: 1.0,
				MaximumInterval:    time.Minute,
				MaximumAttempts:    5,
			},
		}
		ctx = workflow.WithActivityOptions(ctx, ao)
	}
	dslWorkflow.Context = ctx
	dslWorkflow.Running = true
	logger := workflow.GetLogger(ctx)

	//TODO check workflow if need to start second interval
	//workflow.Go(dslWorkflow, dslWorkflow.StartSecondInterval)
	workflow.Go(dslWorkflow, dslWorkflow.Start)
	//dslWorkflow.Start(ctx)
	workflow.Await(ctx, func() bool {
		return dslWorkflow.NeedContinueAsNew || !dslWorkflow.Running
	})
	logger.Info("DSL Workflow completed.")
	dslWorkflow.NeedContinueAsNew = false

	if !dslWorkflow.Running {
		//eventsub.UnRegisterDeviceEventHandler(dslWorkflow.RuleId)
		return nil, nil
	}
	return nil, workflow.NewContinueAsNewError(dslWorkflow, FlowDslWorkflow, dslWorkflow)
}

func (b *Statement) execute(ctx workflow.Context, useLocalActivity bool, params WorkflowParams) (result WorkflowResult, err error) {
	logger := workflow.GetLogger(ctx)
	//logger.Debug("execute statement", b.Type, b.Id, " useLocalActivity=", useLocalActivity)
	if b.Activity == nil {
		logger.Error(b.Id + " activity is nil")
		return
	}
	params.CurrentId = b.Id
	params.Properties = b.Properties
	params.IncomingBusiness = b.IncomingBusiness
	params.OutgoingBusiness = b.OutgoingBusiness
	result, err = b.Activity.execute(ctx, useLocalActivity, params)
	if err != nil {
		logger.Error(b.Type + " execute error " + err.Error())
		err = errors.Wrap(err, b.Type+" execute error")
		return
	}
	if result.Valid {
		params.Type = result.Type
		for _, s := range b.NextStatements {
			if result.DataType == "" || result.DataType == RESULT_DATA_TYPE_OBJECT {
				params.WorkflowPayload = result.Data

				_, err := s.execute(ctx, useLocalActivity, params)
				if err != nil {
					return result, err
				}
			} else {
				for _, v := range result.Datas {
					params.WorkflowPayload = v
					_, err := s.execute(ctx, useLocalActivity, params)
					if err != nil {
						return result, err
					}
				}

			}

		}
	}

	/*
		childCtx, cancelHandler := workflow.WithCancel(ctx)
		selector := workflow.NewSelector(ctx)
		var activityErr error
		for _, s := range b.NextStatements {
			logger.Debug("next statement", s.Type, s.Id)

			params.CurrentId = s.Id
			f := executeAsync(s, childCtx, useLocalActivity, params)
			selector.AddFuture(f, func(f workflow.Future) {
				err := f.Get(ctx, nil)
				if err != nil {
					// cancel all pending activities
					cancelHandler()
					activityErr = err
				}
			})
		}

		for i := 0; i < len(b.NextStatements); i++ {
			selector.Select(ctx) // this will wait for one branch
			if activityErr != nil {
				return result, activityErr
			}
		}
	*/
	return result, nil
}

func (a ActivityInvocation) execute(ctx workflow.Context, useLocalActivity bool, params WorkflowParams) (result WorkflowResult, err error) {

	if useLocalActivity {
		err = workflow.ExecuteLocalActivity(ctx, a.Name, params).Get(ctx, &result)
	} else {
		err = workflow.ExecuteActivity(ctx, a.Name, params).Get(ctx, &result)
	}
	return
}

func executeAsync(exe executable, ctx workflow.Context, useLocalActivity bool, params WorkflowParams) workflow.Future {
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		ret, err := exe.execute(ctx, useLocalActivity, params)
		settable.Set(ret, err)
	})
	return future
}
