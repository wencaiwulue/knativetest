package controllers

import "context"

type Action interface {
	Process(ctx context.Context) interface{}
}
