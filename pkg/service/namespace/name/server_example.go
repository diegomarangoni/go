package name

import (
	"context"
	"fmt"

	name "github.com/diegomarangoni/gomonorepo/pkg/pb/service/namespace/name/v1"
)

type Example struct {
}

func (h Example) HelloWorld(ctx context.Context, req *name.HelloWorldRequest) (*name.HelloWorldResponse, error) {
	return &name.HelloWorldResponse{
		Content: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}
