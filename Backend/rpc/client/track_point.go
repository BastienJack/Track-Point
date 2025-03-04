package rpcclient

import (
	trackpoint "commerce/idl/track_point/kitex_gen/track-point"
	"commerce/idl/track_point/kitex_gen/track-point/trackpointservice"
	"commerce/pkg/etcd"
	"commerce/pkg/middleware"
	"commerce/pkg/viper"

	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"time"
)

var (
	eventConfig             = viper.Init("event")
	trackPointServiceClient trackpointservice.Client
)

func InitTrackPoint(config *viper.Config) {
	// get etcd address
	etcdAddr := fmt.Sprintf("%s:%d",
		eventConfig.Viper.GetString("etcd.host"),
		eventConfig.Viper.GetInt("etcd.port"))

	// create an etcd resolver
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}

	// create trackpointservice
	serviceName := eventConfig.Viper.GetString("server.name")
	c, err := trackpointservice.NewClient(
		serviceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(1),
		client.WithRPCTimeout(30*time.Second),
		client.WithConnectTimeout(30000*time.Millisecond),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithResolver(r),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)

	if err != nil {
		panic(fmt.Sprintf("InitUser err %v", err))
	}
	trackPointServiceClient = c
}

func SendEvent(ctx context.Context, req *trackpoint.SendEventRequest) (*trackpoint.SendEventResponse, error) {
	return trackPointServiceClient.SendEvent(ctx, req)
}

func QueryEvent(ctx context.Context, req *trackpoint.QueryEventRequest) (*trackpoint.QueryEventResponse, error) {
	return trackPointServiceClient.QueryEvent(ctx, req)
}

func DeleteEvent(ctx context.Context, req *trackpoint.DeleteEventRequest) (*trackpoint.DeleteEventResponse, error) {
	return trackPointServiceClient.DeleteEvent(ctx, req)
}
