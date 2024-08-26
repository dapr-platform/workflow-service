package eventpub

import (
	"context"
	"github.com/dapr-platform/common"
	"github.com/pkg/errors"
)

func PublishInternalMessage(ctx context.Context, topic string, msg any) error {
	err := common.GetDaprClient().PublishEvent(ctx, common.PUBSUB_NAME, topic, msg)
	if err != nil {
		err = errors.Wrap(err, "PublishInternalMessage")
	}
	return err
}
