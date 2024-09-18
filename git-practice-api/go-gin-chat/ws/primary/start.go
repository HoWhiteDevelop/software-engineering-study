package primary

import (
	"git-practice-api/go-gin-chat/ws"
	"git-practice-api/go-gin-chat/ws/go_ws"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 定义 serve 的映射关系
var serveMap = map[string]ws.ServeInterface{
	"Serve":   &ws.Serve{},
	"GoServe": &go_ws.GoServe{},
}

func Create() ws.ServeInterface {
	// GoServe or Serve
	_type := viper.GetString("app.serve_type")
	return serveMap[_type]
}

func Start(gin *gin.Context) {
	Create().RunWs(gin)
}

func OnlineUserCount() int {
	return Create().GetOnlineUserCount()
}

func OnlineRoomUserCount(roomId int) int {
	return Create().GetOnlineRoomUserCount(roomId)
}
