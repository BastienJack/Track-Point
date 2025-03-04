package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	httphandler "commerce/http/handler"
	"commerce/pkg/viper"
)

var (
	apiConfig     = viper.Init("api")
	apiServerAddr = fmt.Sprintf("%s:%d", apiConfig.Viper.GetString("server.host"), apiConfig.Viper.GetInt("server.port"))

	curDir, _    = os.Getwd()
	responsePath = filepath.Join(curDir, "/response/web")
)

func registerGroup(server *gin.Engine) {
	commerce := server.Group("/commerce")
	{
		commerce.POST("/register", httphandler.Register)
		commerce.POST("/login", httphandler.Login)

		commerce.Group("/track-point")
		{
			commerce.GET("/get-common-params", httphandler.GetCommonParams)
			commerce.POST("/query-events", httphandler.QueryEvent)
			commerce.POST("/delete-event", httphandler.DeleteEvent)
			commerce.POST("/add-common-params", httphandler.AddCommonParams)
			commerce.POST("/send-event", httphandler.SendEvent)
		}
	}
}

func main() {
	ginServer := gin.Default()

	registerGroup(ginServer)

	err := ginServer.Run(apiServerAddr)
	if err != nil {
		fmt.Print(err)
	}
}
