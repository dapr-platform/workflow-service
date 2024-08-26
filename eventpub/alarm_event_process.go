package eventpub

import (
	"context"
	"github.com/dapr-platform/common"
	"log"
	"strings"
	"time"
)

func ConstructAndSendEvent(ctx context.Context, eventType int32, eventTitle string, eventDescription string, eventStatus int32, eventLevel int32, eventTime time.Time, objectID string, objectName string, location string) {
	dn := objectID + "_" + eventTitle
	dn = strings.Replace(dn, " ", "_", -1)
	event := common.Event{
		Dn:          dn,
		Title:       eventTitle,
		Type:        eventType,
		Description: eventDescription,
		Status:      eventStatus,
		Level:       eventLevel,
		EventTime:   common.LocalTime(eventTime),
		CreateAt:    common.LocalTime(time.Now()),
		UpdatedAt:   common.LocalTime(time.Now()),
		ObjectID:    objectID,
		ObjectName:  objectName,
		Location:    location,
	}
	err := publishEvent(ctx, &event)
	if err != nil {
		log.Println("publishInternalEvent err:", err)
	}
}

func publishEvent(ctx context.Context, event *common.Event) error {

	return common.GetDaprClient().PublishEvent(ctx, common.PUBSUB_NAME, common.EVENT_TOPIC_NAME, event)
}
