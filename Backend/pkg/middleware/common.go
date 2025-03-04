package middleware

import (
	z "commerce/pkg/zap"
	"context"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"go.uber.org/zap"
)

var (
	_      endpoint.Middleware = CommonMiddleware
	logger *zap.SugaredLogger
)

func init() {
	logger = z.InitLogger()
	defer logger.Sync()
}

func CommonMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		ri := rpcinfo.GetRPCInfo(ctx)

		// record real request
		logger.Debugf("real request: %+v", req)

		// record remote userservice information
		logger.Debugf("remote userservice name: %s, remote method: %s",
			ri.To().ServiceName(), ri.To().Method())

		// rpc
		if err := next(ctx, req, resp); err != nil {
			return err
		}

		// record real response
		logger.Infof("real response: %+v", resp)

		return nil
	}
}
