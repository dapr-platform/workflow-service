package eventsub

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
	"log"
	"workflow-service/service"
)

/*
var subscribeMap = sync.Map{}
var subScribeCtxMap = sync.Map{}

func RegisterDeviceEventHandler(id string, ctx workflow.Context, f func(ctx workflow.Context, event map[string]any)) {
	subscribeMap.Store(id, f)
	subScribeCtxMap.Store(id, ctx)
}
func UnRegisterDeviceEventHandler(id string) {
	subscribeMap.Delete(id)
	subScribeCtxMap.Delete(id)
}

*/

func NewDeviceEventHandler(server common.Service, eventPub string, eventTopic string) {

	var sub = &common.Subscription{
		PubsubName: eventPub,
		Topic:      eventTopic,
		Route:      "/DeviceEventHandler",
	}

	err := server.AddTopicEventHandler(sub, deviceEventHandler)

	if err != nil {
		log.Println("sub error ", err)
		panic(err)
	}
}

func deviceEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {

	var event = make(map[string]any, 0)
	err = json.Unmarshal(e.RawData, &event)
	if err != nil {
		log.Println("eventsub - unmarshal error: ", err)
	} else {
		service.SendDeviceDataSignal(event)
		/*
			subscribeMap.Range(func(key, value interface{}) bool {
				wctx, _ := subScribeCtxMap.Load(key)
				value.(func(workflow.Context, map[string]any))(wctx.(workflow.Context), event)
				return true
			})*/

	}

	return false, nil
}
