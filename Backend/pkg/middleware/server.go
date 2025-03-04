package middleware

import (
	"context"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

// let compiler do the type check
var _ endpoint.Middleware = ServerMiddleware

// ServerMiddleware server middleware print client address during rpc.
func ServerMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		ri := rpcinfo.GetRPCInfo(ctx)

		// record client information
		logger.Infof("client address: %v", ri.From().Address())

		// rpc
		if err = next(ctx, req, resp); err != nil {
			return err
		}

		return nil
	}
}
