package middleware

import (
	"context"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

// let compiler do the type check
var _ endpoint.Middleware = ClientMiddleware

// ClientMiddleware client middleware print server address, rpc timeout and connection timeout during rpc.
func ClientMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		ri := rpcinfo.GetRPCInfo(ctx)

		// client record  server information
		logger.Infof("server address: %v, rpc timeout: %v, readwrite timeout: %v",
			ri.To().Address(), ri.Config().RPCTimeout(), ri.Config().ConnectTimeout())

		// rpc
		if err = next(ctx, req, resp); err != nil {
			return err
		}

		return nil
	}
}
