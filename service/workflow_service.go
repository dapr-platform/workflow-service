package service

import (
	"context"
	"encoding/json"
	"github.com/dapr-platform/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"sync"
	"time"
	"workflow-service/config"
	"workflow-service/flow_dsl"
	"workflow-service/model"
	"workflow-service/redis_op"
)

var temporalClient client.Client
var STATUS_ENABLE = 1
var STATUS_DISABLE = 0

var RunningWorkflow = sync.Map{}
var DeviceFilterMap = sync.Map{}
var SecondIntervalMap = sync.Map{}

var cacheLatestDeviceMap = sync.Map{}

type TransLogger struct {
	Logger *logrus.Logger
}

func (t *TransLogger) Debug(msg string, keyvals ...interface{}) {
	t.Logger.Debug(msg, keyvals)
}
func (t *TransLogger) Info(msg string, keyvals ...interface{}) {
	t.Logger.Info(msg, keyvals)
}
func (t *TransLogger) Warn(msg string, keyvals ...interface{}) {
	t.Logger.Warn(msg, keyvals)
}
func (t *TransLogger) Error(msg string, keyvals ...interface{}) {
	t.Logger.Error(msg, keyvals)
}
func init() {
	temporalClient = GetClient()
	go refreshWorkflow(context.Background())
	go loopCheckWorkflow(context.Background())

}

// 启动时，需要将所有workflow停止、启动一次。因为temporal 有non-deterministic error has caused the Workflow Task to fail.还不知道怎么解决
func refreshWorkflow(ctx context.Context) {
	time.Sleep(time.Second * 5)
	startSecondIntervalLoop()
	workflows, err := common.DbQuery[model.Workflow](ctx, common.GetDaprClient(), model.WorkflowTableInfo.Name, "")
	if err != nil {
		common.Logger.Errorln("refreshWorkflow error", err)
		return
	}
	for _, workflow := range workflows {

		id := workflow.Name + "_" + workflow.ID
		err = terminateWorkflowByWorkflowId(ctx, id)
		if err != nil {
			common.Logger.Warningln("terminateWorkflowByWorkflowId error", err)
		}
		if workflow.Status == int32(STATUS_ENABLE) {
			err = StartFlowDslWorkflow(ctx, &workflow)
			if err != nil {
				err = errors.Wrap(err, "StartFlowDslWorkflow error")
				common.Logger.Errorln("StartFlowDslWorkflow error", err)
				continue
			}

		}
	}
}

// 定时检查规则，如果没有被启动则启动，否则停止
func loopCheckWorkflow(ctx context.Context) {
	for {
		time.Sleep(time.Minute)
		workflows, err := common.DbQuery[model.Workflow](ctx, common.GetDaprClient(), model.WorkflowTableInfo.Name, "")
		if err != nil {
			common.Logger.Errorln("refreshWorkflow error", err)
			return
		}
		for _, workflow := range workflows {

			id := workflow.Name + "_" + workflow.ID

			if workflow.Status == int32(STATUS_ENABLE) {
				upTime, exists := RunningWorkflow.Load(id)
				if !exists {
					err = StartFlowDslWorkflow(ctx, &workflow)
					if err != nil {
						err = errors.Wrap(err, "StartFlowDslWorkflow error")
						common.Logger.Errorln("StartFlowDslWorkflow error", err)
						continue
					}
					common.Logger.Debug("loopCheckWorkflow start ", id)
				} else {
					if cast.ToTime(upTime).Unix() != cast.ToTime(workflow.UpdatedTime).Unix() { //更新了， 先停止再启动
						err = terminateWorkflowByWorkflowId(ctx, id)
						if err != nil {
							common.Logger.Warningln("update terminateWorkflowByWorkflowId error", err)
						}
						err = StartFlowDslWorkflow(ctx, &workflow)
						if err != nil {
							err = errors.Wrap(err, "update StartFlowDslWorkflow error")
							common.Logger.Errorln("update StartFlowDslWorkflow error", err)
							continue
						}
					}

				}
			} else {
				_, exists := RunningWorkflow.Load(id)
				if exists {
					err = terminateWorkflowByWorkflowId(ctx, id)
					if err != nil {
						common.Logger.Warningln("terminateWorkflowByWorkflowId error", err)
					}
					common.Logger.Debug("loopCheckWorkflow terminateWorkflowByWorkflowId ", id)
				}
			}
		}

	}
}

func GetClient() client.Client {
	var err error
	for {
		temporalClient, err = client.Dial(client.Options{
			HostPort: config.TEMPORAL_HOST_PORT,
			Logger:   &TransLogger{Logger: common.Logger},
		})
		if err != nil {
			common.Logger.Error("Unable to create client", err)
			time.Sleep(time.Second * 5)
		} else {
			break
		}
	}

	return temporalClient

}

func ChangeWorkflowStatus(ctx context.Context, workflowModel model.Workflow) (result model.Workflow, err error) {
	existsWorkflow, err := common.DbGetOne[model.Workflow](ctx, common.GetDaprClient(), model.WorkflowTableInfo.Name, model.Workflow_FIELD_NAME_id+"="+workflowModel.ID)
	if err != nil {
		err = errors.Wrap(err, "db get one error")
		return
	}
	if existsWorkflow == nil {
		err = errors.New("can't find workflow by " + workflowModel.ID)
		return
	}
	existsWorkflow.Status = workflowModel.Status

	if existsWorkflow.Status == int32(STATUS_ENABLE) {
		err = StartFlowDslWorkflow(ctx, existsWorkflow)
		if err != nil {
			err = errors.Wrap(err, "StartFlowDslWorkflow error")
			return
		}
	} else {
		err = disableOneWorkflowByModel(ctx, existsWorkflow)
		if err != nil {
			err = errors.Wrap(err, "disableOneWorkflowByModel error")
			return
		}
	}

	err = common.DbUpsert[model.Workflow](ctx, common.GetDaprClient(), *existsWorkflow, model.WorkflowTableInfo.Name, model.Workflow_FIELD_NAME_id)
	if err == nil {
		common.PublishDbUpsertMessage(ctx, common.GetDaprClient(), model.WorkflowTableInfo.Name, model.Workflow_FIELD_NAME_id, "", false, *existsWorkflow)
	}
	result = workflowModel
	return
}
func disableOneWorkflowByModel(ctx context.Context, workflowModel *model.Workflow) (err error) {
	parser := flow_dsl.JsonParser{}
	workflow, err := parser.Parse(workflowModel.Content)
	if err != nil {
		err = errors.Wrap(err, "解析dsl内容错误")
		return
	}
	buf, _ := json.Marshal(workflow)
	common.Logger.Debug("disableOneWorkflowByModel workflow json:", string(buf))
	if workflowModel.WorkflowID != "" {
		terminateWorkflowByWorkflowId(ctx, workflowModel.WorkflowID)
	}
	err = common.DbUpsert[model.Workflow](ctx, common.GetDaprClient(), *workflowModel, model.WorkflowTableInfo.Name, model.Workflow_FIELD_NAME_id)
	if err != nil {
		err = errors.Wrap(err, "db upsert error")
	}
	return
}
func terminateWorkflowByWorkflowId(ctx context.Context, workflowId string) (err error) {
	if common.RUNNING_MODE == common.RUNNING_MODE_CENTER {
		return
	}
	common.Logger.Debug("terminateWorkflowByWorkflowId workflowId:", workflowId)
	m := map[string]interface{}{}
	m["status"] = false
	err = temporalClient.TerminateWorkflow(ctx, workflowId, "", "restart")
	common.Logger.Debug("terminateWorkflowByWorkflowId workflow ", workflowId)
	RunningWorkflow.Delete(workflowId)
	DeviceFilterMap.Delete(workflowId)
	SecondIntervalMap.Delete(workflowId)
	return
}
func disableWorkflowByWorkflowId(ctx context.Context, workflowId string) (err error) {
	if common.RUNNING_MODE == common.RUNNING_MODE_CENTER {
		return
	}
	common.Logger.Debug("disableWorkflowByWorkflowId workflowId:", workflowId)
	m := map[string]interface{}{}
	m["status"] = false
	err = temporalClient.SignalWorkflow(ctx, workflowId, "", "Status", m)
	common.Logger.Debug("disableOneWorkflowByModel workflow signal:", workflowId, "", "Status", m)
	return
}
func checkOneWorkflowByModel(ctx context.Context, workflowModel model.Workflow) (err error) {
	parser := flow_dsl.JsonParser{}
	workflow, err := parser.Parse(workflowModel.Content)
	if err != nil {
		err = errors.Wrap(err, "解析dsl内容错误")
		return
	}

	buf, _ := json.Marshal(workflow)
	common.Logger.Debug("workflow json:", string(buf))

	return
}
func StartFlowDslWorkflow(ctx context.Context, workflowModel *model.Workflow) (err error) {
	if common.RUNNING_MODE == common.RUNNING_MODE_CENTER {
		return
	}
	common.Logger.Debug("StartFlowDslWorkflow workflowModel:", workflowModel)
	parser := flow_dsl.JsonParser{}
	workflow, err := parser.Parse(workflowModel.Content)
	if err != nil {
		err = errors.Wrap(err, "解析dsl内容错误")
		return
	}
	wfjson, _ := json.Marshal(workflow)
	common.Logger.Debug("workflow json:", string(wfjson))
	/*
		err = temporalClient.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
			TaskQueue: "flowdsl",
			Operation: &client.BuildIDOpAddNewIDInNewDefaultSet{
				BuildID: "1.0",
			},
		})
		if err != nil {
			common.Logger.Errorln("UpdateWorkerBuildIdCompatibility error", err)
		}

	*/
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowModel.Name + "_" + workflowModel.ID,
		TaskQueue: "flowdsl",
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 1.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
	}
	workflow.UseLocalActivity = true
	workflow.Running = true
	we, err := temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, flow_dsl.FlowDslWorkflow, &workflow)
	if err != nil {
		return err
	}
	common.Logger.Info("Started workflow ", " WorkflowID ", we.GetID(), " RunID ", we.GetRunID())
	workflowModel.WorkflowID = we.GetID()
	workflowModel.RunningID = we.GetRunID()
	err = common.DbUpsert[model.Workflow](ctx, common.GetDaprClient(), *workflowModel, model.WorkflowTableInfo.Name, model.Workflow_FIELD_NAME_id)
	if err != nil {
		err = errors.Wrap(err, "db upsert error")
		temporalClient.CancelWorkflow(ctx, we.GetID(), we.GetRunID())
		return
	}
	RunningWorkflow.Store(we.GetID(), workflowModel.UpdatedTime)
	filters := flow_dsl.ParseWorkflowDeviceFilter(&workflow)
	DeviceFilterMap.Store(we.GetID(), filters)
	usingSecondInterval := flow_dsl.IsWorkflowUsingSecondInterval(&workflow)
	SecondIntervalMap.Store(we.GetID(), usingSecondInterval)

	return
}
func DeleteWorkflow(ctx context.Context, id string) (err error) {

	//删除工作流可能用到的redis keys
	workflow, err := common.DbGetOne[model.Workflow](ctx, common.GetDaprClient(), model.WorkflowTableInfo.Name, model.Workflow_FIELD_NAME_id+"="+id)
	if err != nil {
		common.Logger.Error("DbGetOne error", err)
	} else {
		redisKeyPrefix := flow_dsl.GetRedisPrefix(workflow.Name)
		keys, err1 := redis_op.Rdb.Keys(ctx, redisKeyPrefix+"*").Result()
		if err1 != nil {
			common.Logger.Error("keys error " + err1.Error())
		} else {
			for _, key := range keys {
				_ = redis_op.Rdb.Del(ctx, key).Err()
			}
		}
	}
	return common.DbDelete(ctx, common.GetDaprClient(), model.WorkflowTableInfo.Name, model.Workflow_FIELD_NAME_id, id)
}
func SaveWorkflow(ctx context.Context, workflowModel model.Workflow) (result model.Workflow, err error) {
	err = checkOneWorkflowByModel(ctx, workflowModel)
	if err != nil {
		err = errors.Wrap(err, "checkOneWorkflowByModel error")
		//return
	}
	err = common.DbUpsert[model.Workflow](ctx, common.GetDaprClient(), workflowModel, model.WorkflowTableInfo.Name, model.Workflow_FIELD_NAME_id)
	result = workflowModel
	return
}

func startSecondIntervalLoop() {
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			var data = flow_dsl.SecondInterval{
				Second: time.Now().Unix(),
			}
			SecondIntervalMap.Range(func(key, value any) bool {
				if cast.ToBool(value) {
					temporalClient.SignalWorkflow(context.Background(), key.(string), "", "SecondInterval", data)
				}
				return true
			})
		}
		ticker.Stop()
	}()
}
func SendDeviceDataSignal(device map[string]any) {
	common.Logger.Debug("SendDeviceDataSignal device:", device)
	/*
		old, exists := cacheLatestDeviceMap.Load(device["id"])
		if exists {
			needSend := false
			for k, v := range device {
				if cast.ToString(v) != cast.ToString(old.(map[string]any)[k]) {
					needSend = true
					break
				}
			}
			if !needSend {
				return
			}
		}
		cacheLatestDeviceMap.Store(device["id"], device)

	*/
	DeviceFilterMap.Range(func(key, value any) bool {
		common.Logger.Debug("device signal " + key.(string))

		deviceFilter := value.([]flow_dsl.DeviceFilter)
		shouldSend := false
		for _, filter := range deviceFilter {
			if filter.Filter(device) {
				shouldSend = true
				break
			}
		}
		if shouldSend {
			err := temporalClient.SignalWorkflow(context.Background(), key.(string), "", "Device", device)
			if err != nil {
				common.Logger.Error("device signal error", err)
			}
		}

		return true
	})
	return
}
