package rpcclient

import "commerce/pkg/viper"

// rpc client initialization
func init() {
	// user rpc config
	userRpcConfig := viper.Init("user")
	InitUser(&userRpcConfig)

	// trackpoint rpc config
	trackpointRpcConfig := viper.Init("event")
	InitTrackPoint(&trackpointRpcConfig)
}
