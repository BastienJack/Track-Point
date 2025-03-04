package httphandler

import (
	trackpoint "commerce/idl/track_point/kitex_gen/track-point"
	"commerce/message"
	"commerce/pkg/viper"
	rpcclient "commerce/rpc/client"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	eventConfig = viper.Init("event")

	eventCommonParamsKey = "event_params"
	eventCommonParams    = eventConfig.Viper.Get(eventCommonParamsKey)
)

func GetCommonParams(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, eventCommonParams)
}

func AddCommonParams(ctx *gin.Context) {
	// get raw http json data
	body, err := ctx.GetRawData()
	if err != nil {
		panic(err)
	}

	// get json data
	var param message.AddCommonParamRequest
	err = json.Unmarshal(body, &param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, message.AddCommonParamResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	// try to set key
	eventConfig.Viper.Set(eventCommonParamsKey+"."+param.Key, eventCommonParamsKey+"."+param.Value)

	// save config
	err = eventConfig.Viper.SafeWriteConfig()

	// renew config
	eventConfig = viper.Init("event")
	eventCommonParams = eventConfig.Viper.Get(eventCommonParamsKey)

	// error response
	if err != nil {
		ctx.JSON(http.StatusBadRequest, message.AddCommonParamResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		})

		return
	}

	// success response
	ctx.JSON(http.StatusOK, message.AddCommonParamResponse{
		Base: message.Base{
			StatusCode: 0,
			StatusMsg:  "Add params success.",
		},
	})
}

func SendEvent(ctx *gin.Context) {
	// get raw http json data
	body, err := ctx.GetRawData()
	if err != nil {
		panic(err)
	}

	// get json data
	var event *message.SendEventRequest
	err = json.Unmarshal(body, &event)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, message.AddCommonParamResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	// send event RPC request data
	sendEventRpcReq := &trackpoint.SendEventRequest{
		JsonEventParams: event.Event,
	}

	// send event RPC call
	sendEventRpcRes, err := rpcclient.SendEvent(ctx, sendEventRpcReq)

	// send event RPC error response
	if err != nil || sendEventRpcRes.StatusCode == -1 {
		ctx.JSON(http.StatusBadRequest, message.SendEventResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  sendEventRpcRes.StatusMsg,
			},
		})
		return
	}

	// send event RPC success response
	ctx.JSON(http.StatusOK, message.SendEventResponse{
		Base: message.Base{
			StatusCode: 0,
			StatusMsg:  sendEventRpcRes.StatusMsg,
		},
	})
}

func QueryEvent(ctx *gin.Context) {
	// get raw http json data
	body, err := ctx.GetRawData()
	if err != nil {
		panic(err)
	}

	// get json data
	var queryParam *message.QueryEventRequest
	err = json.Unmarshal(body, &queryParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, message.QueryEventResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
			Events: make([]*trackpoint.Event, 0),
		})
		return
	}

	// query event RPC request data
	queryEventRpcReq := &trackpoint.QueryEventRequest{
		Offset: int32(queryParam.Offset),
		Limit:  int32(queryParam.Limit),
	}

	// query event RPC call
	queryEventRpcRes, err := rpcclient.QueryEvent(ctx, queryEventRpcReq)

	// query event RPC error response
	if err != nil || queryEventRpcRes.StatusCode == -1 {
		ctx.JSON(http.StatusBadRequest, message.QueryEventResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  queryEventRpcRes.StatusMsg,
			},
			Events: queryEventRpcRes.Events,
		})
		return
	}

	// query event RPC success response
	ctx.JSON(http.StatusOK, message.QueryEventResponse{
		Base: message.Base{
			StatusCode: 0,
			StatusMsg:  queryEventRpcRes.StatusMsg,
		},
		Events: queryEventRpcRes.Events,
	})
}

func DeleteEvent(ctx *gin.Context) {
	// get raw http json data
	body, err := ctx.GetRawData()
	if err != nil {
		panic(err)
	}

	// get json data
	var queryParam *message.DeleteEventRequest
	err = json.Unmarshal(body, &queryParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, message.DeleteEventResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	// delete event RPC request data
	deleteEventRpcReq := &trackpoint.DeleteEventRequest{
		EventId: queryParam.EventId,
	}

	// delete event RPC call
	deleteEventRpcRes, err := rpcclient.DeleteEvent(ctx, deleteEventRpcReq)

	// delete event RPC error response
	if err != nil || deleteEventRpcRes.StatusCode == -1 {
		ctx.JSON(http.StatusBadRequest, message.DeleteEventResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  deleteEventRpcRes.StatusMsg,
			},
		})
		return
	}

	// delete event RPC success response
	ctx.JSON(http.StatusOK, message.DeleteEventResponse{
		Base: message.Base{
			StatusCode: 0,
			StatusMsg:  deleteEventRpcRes.StatusMsg,
		},
	})
}
