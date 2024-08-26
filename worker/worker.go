package worker

import (
	"github.com/dapr-platform/common"
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"time"
	"workflow-service/config"
	"workflow-service/flow_dsl"
)

func init() {
	go newWorker() //启动多个？

}

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

func newWorker() {
	var c client.Client
	var err error
	for {
		c, err = client.Dial(client.Options{
			HostPort: config.TEMPORAL_HOST_PORT,
			Logger:   &TransLogger{Logger: common.Logger},
		})
		if err != nil {
			common.Logger.Fatalln("Unable to create client", err)
			time.Sleep(time.Second * 5)
		} else {
			break
		}
	}

	w := worker.New(c, "flowdsl", worker.Options{})
	w.RegisterWorkflow(flow_dsl.FlowDslWorkflow)
	w.RegisterActivity(&flow_dsl.FlowDslActivities{})

	err = w.Start()
	if err != nil {
		common.Logger.Fatalln("Unable to start worker", err)
	}

}
