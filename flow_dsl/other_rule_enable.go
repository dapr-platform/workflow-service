package flow_dsl

import (
	"context"
)

func (d *FlowDslActivities) OtherRuleEnable(ctx context.Context, params WorkflowParams) (result WorkflowResult, err error) {
	result.Valid = true
	result.Data = params.WorkflowPayload
	targetPort := params.OutgoingBusiness["OtherRuleEnable@output"]
	result.Data.MatchPort = targetPort
	result.Type = PARAM_TYPE_EVENT
	//logger := activity.GetLogger(ctx)
	//logger.Debug("matchport " + targetPort)
	return
}
