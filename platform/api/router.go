package api

import "github.com/gin-gonic/gin"

type ApiServ struct{}

func RegisteRouter(r *gin.RouterGroup) *ApiServ {
	svr := &ApiServ{}
	nodes := r.Group("nodes")
	{
		nodes.POST("/create", svr.Create)
	}
	infos := r.Group("infos")
	{
		infos.GET("/", svr.Index)
	}
	return svr
}
