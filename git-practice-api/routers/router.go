package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// dev 开发用 避免修改静态资源需要重启服务
	router.StaticFS("/static", http.Dir("static"))
	return router
}
