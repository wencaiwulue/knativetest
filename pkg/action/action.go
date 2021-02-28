package action

import "context"

type Action interface {
	Process(ctx context.Context) interface{}
}
