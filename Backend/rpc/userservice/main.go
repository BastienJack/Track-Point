package main

import (
	"commerce/pkg/zap"
	rpchandler "commerce/rpc/userservice/rpchandler"
)

var (
	logger = zap.InitLogger()
)

func main() {
	userServiceErr := rpchandler.StartUserService()
	if userServiceErr != nil {
		logger.Fatalf("User userservice stopped with error: %v", userServiceErr.Error())
	}
}
