package platform

import (
	"fmt"

	"chen.com/distributecache/config"
	"chen.com/distributecache/platform/api"
	"github.com/gin-gonic/gin"
)

func StartServ() {
	conf := config.GetConfig()
	route := gin.Default()
	v1 := route.Group("/api/v1")
	{
		api.RegisteRouter(v1)
	}
	route.Run(fmt.Sprintf("%s:%d", conf.PlatformAddr, conf.PlatformPort))
}
