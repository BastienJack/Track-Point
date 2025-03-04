package main

import (
	"commerce/pkg/zap"
	"commerce/rpc/trackpointservice/rpchandler"
)

var (
	logger = zap.InitLogger()
)

func main() {
	trackPointServiceErr := rpchandler.StartTrackPointService()
	if trackPointServiceErr != nil {
		logger.Fatalf("User userservice stopped with error: %v", trackPointServiceErr.Error())
	}
}
