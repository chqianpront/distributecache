package api

import "github.com/gin-gonic/gin"

func (svr *ApiServ) Create(ctx *gin.Context) {
	addr := ctx.PostForm("addr")
	port := ctx.PostForm("port")

}
