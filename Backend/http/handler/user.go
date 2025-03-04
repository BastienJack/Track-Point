package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	userrpc "commerce/idl/user/kitex_gen/user"
	"commerce/message"
	rpcclient "commerce/rpc/client"
)

func Register(ctx *gin.Context) {
	// get raw http json data
	body, err := ctx.GetRawData()
	if err != nil {
		panic(err)
	}

	// unmarshall raw http json data
	var regHttpReq message.RegisterRequest
	if err := json.Unmarshal(body, &regHttpReq); err != nil {
		ctx.JSON(http.StatusBadRequest, message.RegisterResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  "Unmarshall json data error.",
			},
		})
		return
	}

	// length check
	if len(regHttpReq.Username) == 0 || len(regHttpReq.Password) == 0 {
		ctx.JSON(http.StatusBadRequest, message.RegisterResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  "Username or password should not be empty.",
			},
		})
	}

	// confirm password error
	if regHttpReq.Password != regHttpReq.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, message.RegisterResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  "Confirm password not identical.",
			},
		})
	}

	// register RPC request data
	regRpcReq := &userrpc.RegisterRequest{
		Username:        regHttpReq.Username,
		Password:        regHttpReq.Password,
		ConfirmPassword: regHttpReq.ConfirmPassword,
	}

	// register RPC call
	regRpcRes, err := rpcclient.Register(ctx, regRpcReq)

	// register RPC error response
	if err != nil || regRpcRes.StatusCode == -1 {
		ctx.JSON(http.StatusBadRequest, message.RegisterResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  regRpcRes.StatusMsg,
			},
		})
		return
	}

	// register RPC success response
	ctx.JSON(http.StatusOK, message.RegisterResponse{
		Base: message.Base{
			StatusCode: 0,
			StatusMsg:  regRpcRes.StatusMsg,
		},
		UserID: regRpcRes.UserId,
	})
}

func Login(ctx *gin.Context) {
	// get raw http json data
	body, err := ctx.GetRawData()
	if err != nil {
		panic(err)
	}

	// unmarshall raw http json data
	var loginHttpReq message.LoginRequest
	if err := json.Unmarshal(body, &loginHttpReq); err != nil {
		ctx.JSON(http.StatusBadRequest, message.LoginResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  "Unmarshall json data error.",
			},
		})
		return
	}

	// login RPC request data
	loginRpcReq := &userrpc.LoginRequest{
		Username: loginHttpReq.Username,
		Password: loginHttpReq.Password,
	}

	// login RPC call
	loginRpcRes, err := rpcclient.Login(ctx, loginRpcReq)

	// login RPC error response
	if err != nil || loginRpcRes.StatusCode == -1 {
		ctx.JSON(http.StatusBadRequest, message.LoginResponse{
			Base: message.Base{
				StatusCode: -1,
				StatusMsg:  loginRpcRes.StatusMsg,
			},
		})
		return
	}

	// login RPC success response
	ctx.JSON(http.StatusOK, message.LoginResponse{
		Base: message.Base{
			StatusCode: 0,
			StatusMsg:  loginRpcRes.StatusMsg,
		},
		UserID: loginRpcRes.UserId,
	})
}
