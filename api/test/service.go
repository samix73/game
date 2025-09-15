//go:generate rpc-gen -output ./../generated

package test

import "context"

type Service interface {
	Test(ctx context.Context, req Request) (*Response, error)
}

type Request struct {
	Data string
}

type Response struct {
	Result string
}
