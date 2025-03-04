package rpchandler

import (
	"context"
	"errors"
	"fmt"
	"net"

	"gorm.io/gorm"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"

	"commerce/idl/user/kitex_gen/user"
	"commerce/idl/user/kitex_gen/user/userservice"

	"commerce/storage/db"

	"commerce/pkg/etcd"
	"commerce/pkg/middleware"
	"commerce/pkg/viper"
	"commerce/pkg/zap"
)

var (
	userConfig = viper.Init("user")

	userServiceName = userConfig.Viper.GetString("server.name")
	userServiceAddr = fmt.Sprintf("%s:%d",
		userConfig.Viper.GetString("server.host"),
		userConfig.Viper.GetInt("server.port"))

	userEtcdAddr = fmt.Sprintf("%s:%d",
		userConfig.Viper.GetString("etcd.host"),
		userConfig.Viper.GetInt("etcd.port"))

	serviceLogger = zap.InitLogger()
)

func StartUserService() error {
	// userservice registry
	r, err := etcd.NewEtcdRegistry([]string{userEtcdAddr})
	if err != nil {
		serviceLogger.Fatalf("Server register failed: %v", err)
	}

	// userservice resolver for userservice discovery
	addr, err := net.ResolveTCPAddr("tcp", userServiceAddr)
	if err != nil {
		serviceLogger.Fatalf("Resolver tcp addr failed: %v", err)
	}

	// initialize rpc server
	s := userservice.NewServer(new(UserServiceImpl),
		server.WithServiceAddr(addr),
		server.WithMiddleware(middleware.CommonMiddleware),
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: userServiceName}),
	)

	// run rpc server
	if err := s.Run(); err != nil {
		return err
	}

	return nil
}

// UserServiceImpl implements the last userservice interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	logger := zap.InitLogger()

	// check username conflict
	usr, err := db.GetUserByName(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("Error: select username error. Info: %+v.", err)

		res := &user.RegisterResponse{
			StatusCode: -1,
			StatusMsg:  "Error: server internal error.",
		}

		return res, nil
	} else if usr != nil && usr.Username != "" {
		logger.Errorf("Error: username conflict. Info: user_info=%+v.", usr)

		res := &user.RegisterResponse{
			StatusCode: -1,
			StatusMsg:  "Error: username conflict.",
		}

		return res, nil
	}

	// create user
	usr = &db.User{
		Username: req.Username,
		Password: req.Password,
	}

	if err := db.CreateUser(ctx, usr); err != nil {
		logger.Errorf("Error: create user in database error. Info: %+v.", err)

		res := &user.RegisterResponse{
			StatusCode: -1,
			StatusMsg:  "Error: server internal error.",
		}

		return res, nil
	}

	// query user by name
	usr, err = db.GetUserByName(ctx, req.Username)
	if err != nil || usr == nil {
		logger.Errorf("Error: get user id error. Info: error=%+v, user=%+v.", err, usr)

		res := &user.RegisterResponse{
			StatusCode: -1,
			StatusMsg:  "Error: server internal error.",
		}

		return res, nil
	}

	// register response
	res := &user.RegisterResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		UserId:     int64(usr.ID),
	}

	return res, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	logger := zap.InitLogger()

	// query user by name
	usr, err := db.GetUserByName(ctx, req.Username)
	if err != nil || usr == nil {
		logger.Errorf("Error: get user id error. Info: error=%+v, user=%+v.", err, usr)

		res := &user.LoginResponse{
			StatusCode: -1,
			StatusMsg:  "Error: server internal error.",
		}

		return res, nil
	}

	// login failure response
	if usr.Password != req.Password {
		res := &user.LoginResponse{
			StatusCode: -1,
			StatusMsg:  "Error: password not match.",
		}

		return res, nil
	}

	// login success response
	res := &user.LoginResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		UserId:     int64(usr.ID),
	}

	return res, nil
}
