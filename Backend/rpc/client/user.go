package rpcclient

import (
	"commerce/idl/user/kitex_gen/user"
	"commerce/idl/user/kitex_gen/user/userservice"
	"commerce/pkg/etcd"
	"commerce/pkg/middleware"
	"commerce/pkg/viper"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"

	"context"
	"fmt"
	"time"
)

var (
	userConfig = viper.Init("user")
	userClient userservice.Client
)

func InitUser(config *viper.Config) {
	// get etcd address
	etcdAddr := fmt.Sprintf("%s:%d",
		userConfig.Viper.GetString("etcd.host"),
		userConfig.Viper.GetInt("etcd.port"))

	// create an etcd resolver
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}

	// create user userservice
	serviceName := userConfig.Viper.GetString("server.name")
	c, err := userservice.NewClient(
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
	userClient = c
}

func Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	return userClient.Register(ctx, req)
}

func Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	return userClient.Login(ctx, req)
}
