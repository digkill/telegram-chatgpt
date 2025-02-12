package domains

import (
	"context"
)

type ActionResultType string

const (
	ActionResultText     ActionResultType = "text"
	ActionResultImageB64 ActionResultType = "image_b64"
)

type ActionResult struct {
	Result string
	Type   ActionResultType
}

type ActionInfo struct {
	Message *Message
	Result  *ActionResult
}

func (a *ActionInfo) ExistsResult() bool {
	return a.Result != nil
}

func (a *ActionInfo) GetText() string {
	if msg := a.Message; msg != nil {
		return msg.Content
	}

	return ""
}

type Action interface {
	Handle(ctx context.Context, actionInfo *ActionInfo) (res *ActionResult, err error)
}
