package go_ws

import (
	"encoding/json"
	"git-practice-api/go-gin-chat/models"
	"git-practice-api/go-gin-chat/services/helper"
	"git-practice-api/go-gin-chat/services/safe"
	"git-practice-api/go-gin-chat/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jianfengye/collection"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 客户端连接详情
type wsClients struct {
	Conn *websocket.Conn `json:"conn"`

	RemoteAddr string `json:"remote_addr"`

	Uid string `json:"uid"`

	Username string `json:"username"`

	RoomId string `json:"room_id"`

	AvatarId string `json:"avatar_id"`
}

type msgData struct {
	Uid      string        `json:"uid"`
	Username string        `json:"username"`
	AvatarId string        `json:"avatar_id"`
	ToUid    string        `json:"to_uid"`
	Content  string        `json:"content"`
	ImageUrl string        `json:"image_url"`
	RoomId   string        `json:"room_id"`
	Count    int           `json:"count"`
	List     []interface{} `json:"list"`
	Time     int64         `json:"time"`
}

// client & serve 的消息体
type msg struct {
	Status int             `json:"status"`
	Data   msgData         `json:"data"`
	Conn   *websocket.Conn `json:"conn"`
}

type pingStorage struct {
	Conn       *websocket.Conn `json:"conn"`
	RemoteAddr string          `json:"remote_addr"`
	Time       int64           `json:"time"`
}

// 变量定义初始化
var (
	wsUpgrader = websocket.Upgrader{}

	clientMsg = msg{}

	mutex = sync.Mutex{}

	//rooms = [roomCount + 1][]wsClients{}
	//值是存储在该房间内的 WebSocket 客户端集合（[]interface{}）。
	//这里的 interface{} 代表的是一个可以存储任何类型的空接口，
	//但实际存储的是 wsClients 结构体的实例。
	rooms = make(map[int][]interface{})
	//用于传递新加入房间的客户端连接。
	//它接收的是 wsClients 结构体，
	//这个结构体包含了客户端连接的信息（如 Conn、RemoteAddr、Uid 等
	enterRooms = make(chan wsClients)
	//用于传递服务器和客户端之间的消息（类型为 msg 结构体）。
	//当服务器或客户端发消息时，消息会通过该通道在不同的 goroutine 之间传递。
	sMsg = make(chan msg)
	//这是一个 chan 通道，用于传递需要断开连接的客户端（类型为 *websocket.Conn）。
	//当客户端连接断开或超时时，连接会通过该通道传递，并触发相关的断开处理。
	offline = make(chan *websocket.Conn)
	//这是一个带有缓冲区大小为 1 的 chan 通道，
	//通常用于在并发操作中起到同步或通知的作用。
	//它在一些操作（如消息通知、写入）之间起到阻塞或控制的作用，避免多个 goroutine 并发修改同一资源。
	//在写消息时，它会阻塞其他 goroutine，
	//确保同一时间只有一个消息被发送，
	//避免 concurrent write to websocket connection 错误。
	chNotify = make(chan int, 1)
	//pingMap 是一个 []interface{} 切片，用于存储心跳检测的客户端连接信息。
	//每个连接的信息被存储为 pingStorage 结构体，记录了连接的 WebSocket 对象、远程地址和最近一次心跳的时间。
	//它用于处理 WebSocket 心跳机制，定期检查并清理超时的客户端连接。
	pingMap []interface{}

	clientMsgLock = sync.Mutex{}
	clientMsgData = clientMsg // 临时存储 clientMsg 数据
)

// 定义消息类型
const msgTypeOnline = 1        // 上线
const msgTypeOffline = 2       // 离线
const msgTypeSend = 3          // 消息发送
const msgTypeGetOnlineUser = 4 // 获取用户列表
const msgTypePrivateChat = 5   // 私聊

const roomCount = 6 // 房间总数

type GoServe struct {
	ws.ServeInterface
}

func (goServe *GoServe) RunWs(gin *gin.Context) {
	// 使用 channel goroutine
	Run(gin)
}

func (goServe *GoServe) GetOnlineUserCount() int {
	return GetOnlineUserCount()
}

func (goServe *GoServe) GetOnlineRoomUserCount(roomId int) int {
	return GetOnlineRoomUserCount(roomId)
}

// 建立连接开启读写的协程
func Run(gin *gin.Context) {

	// @see https://github.com/gorilla/websocket/issues/523
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	//建立websocket连接
	c, _ := wsUpgrader.Upgrade(gin.Writer, gin.Request, nil)

	defer c.Close()
	done := make(chan struct{})

	go read(c, done)
	go write(done)

	//for {
	//	select {
	//	case <-done:
	//		return
	//	}
	//}
	//类似于for{}用于阻塞
	select {}

}

// HandelOfflineCoon 定时任务清理没有心跳的连接
func HandelOfflineCoon() {
	//将 pingMap 转换成一个 ObjCollection 集合，这样就可以对其进行集合操作，如筛选、过滤、追加等
	objColl := collection.NewObjCollection(pingMap)
	//safe.Safety.Do(...) 是一个并发安全的操作。它确保在对集合 objColl 进行操作时，不会有并发问题。
	retColl := safe.Safety.Do(func() interface{} {
		//从集合中删除符合某些条件的对象
		return objColl.Reject(func(obj interface{}, index int) bool {
			nowTime := time.Now().Unix()
			timeDiff := nowTime - obj.(pingStorage).Time
			// log.Println("timeDiff", nowTime, obj.(pingStorage).Time, timeDiff)
			if timeDiff > 60 { // 超过 60s 没有心跳 主动断开连接
				offline <- obj.(pingStorage).Conn
				return true
			}
			return false
		})
	}).(collection.ICollection)

	interfaces, _ := retColl.ToInterfaces()

	pingMap = interfaces
}

func appendPing(c *websocket.Conn) {
	objColl := collection.NewObjCollection(pingMap)

	// 先删除相同的
	retColl := safe.Safety.Do(func() interface{} {
		return objColl.Reject(func(obj interface{}, index int) bool {
			if obj.(pingStorage).RemoteAddr == c.RemoteAddr().String() {
				return true
			}
			return false
		})
	}).(collection.ICollection)

	// 再追加
	retColl = safe.Safety.Do(func() interface{} {
		return retColl.Append(pingStorage{
			Conn:       c,
			RemoteAddr: c.RemoteAddr().String(),
			Time:       time.Now().Unix(),
		})
	}).(collection.ICollection)

	interfaces, _ := retColl.ToInterfaces()

	pingMap = interfaces

}

func read(c *websocket.Conn, done chan<- struct{}) {

	defer func() {
		//捕获read抛出的panic
		if err := recover(); err != nil {
			log.Println("read发生错误", err)
			//panic(nil)
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		//log.Println("client message", string(message), c.RemoteAddr())
		if err != nil { // 离线通知
			offline <- c
			log.Println("ReadMessage error1", err)
			c.Close()
			close(done)
			return
		}

		// 处理心跳响应 , heartbeat为与客户端约定的值
		if string(message) == `heartbeat` {
			appendPing(c)
			chNotify <- 1
			// log.Println("heartbeat pingMap：", pingMap)
			c.WriteMessage(websocket.TextMessage, []byte(`{"status":0,"data":"heartbeat ok"}`))
			<-chNotify
			continue
		}

		json.Unmarshal(message, &clientMsgData)

		clientMsgLock.Lock()
		clientMsg = clientMsgData
		clientMsgLock.Unlock()

		//fmt.Println("来自客户端的消息", clientMsg, c.RemoteAddr())
		if clientMsg.Data.Uid != "" {
			if clientMsg.Status == msgTypeOnline { // 进入房间，建立连接
				roomId, _ := getRoomId()

				enterRooms <- wsClients{
					Conn:       c,
					RemoteAddr: c.RemoteAddr().String(),
					Uid:        clientMsg.Data.Uid,
					Username:   clientMsg.Data.Username,
					RoomId:     roomId,
					AvatarId:   clientMsg.Data.AvatarId,
				}
			}

			_, serveMsg := formatServeMsgStr(clientMsg.Status, c)
			sMsg <- serveMsg
		}
	}
}

func write(done <-chan struct{}) {

	defer func() {
		//捕获write抛出的panic
		if err := recover(); err != nil {
			log.Println("write发生错误", err)
			//panic(err)
		}
	}()

	for {
		select {
		case <-done: // 当 done 通道关闭时，退出 write 函数
			return
		case r := <-enterRooms:
			handleConnClients(r.Conn)
		case cl := <-sMsg:
			serveMsgStr, _ := json.Marshal(cl)
			switch cl.Status {
			case msgTypeOnline, msgTypeSend:
				notify(cl.Conn, string(serveMsgStr))
			case msgTypeGetOnlineUser:
				chNotify <- 1
				cl.Conn.WriteMessage(websocket.TextMessage, serveMsgStr)
				<-chNotify
			case msgTypePrivateChat:
				chNotify <- 1
				toC := findToUserCoonClient()
				if toC != nil {
					toC.(wsClients).Conn.WriteMessage(websocket.TextMessage, serveMsgStr)
				}
				<-chNotify
			}
		case o := <-offline:
			disconnect(o)
		}
	}
}

func handleConnClients(c *websocket.Conn) {
	roomId, roomIdInt := getRoomId()

	objColl := collection.NewObjCollection(rooms[roomIdInt])
	//过滤移除相同 UID 的用户
	retColl := safe.Safety.Do(func() interface{} {
		return objColl.Reject(func(item interface{}, key int) bool {
			if item.(wsClients).Uid == clientMsg.Data.Uid {
				chNotify <- 1
				item.(wsClients).Conn.WriteMessage(websocket.TextMessage, []byte(`{"status":-1,"data":[]}`))
				<-chNotify
				return true
			}
			return false
		})
	}).(collection.ICollection)
	//将新的用户信息追加到房间
	retColl = safe.Safety.Do(func() interface{} {
		return retColl.Append(wsClients{
			Conn:       c,
			RemoteAddr: c.RemoteAddr().String(),
			Uid:        clientMsg.Data.Uid,
			Username:   clientMsg.Data.Username,
			RoomId:     roomId,
			AvatarId:   clientMsg.Data.AvatarId,
		})
	}).(collection.ICollection)
	//更新房间中的用户集合
	interfaces, _ := retColl.ToInterfaces()

	rooms[roomIdInt] = interfaces
}

// 获取私聊的用户连接
func findToUserCoonClient() interface{} {
	_, roomIdInt := getRoomId()

	toUserUid := clientMsg.Data.ToUid
	assignRoom := rooms[roomIdInt]
	for _, c := range assignRoom {
		stringUid := c.(wsClients).Uid
		if stringUid == toUserUid {
			return c
		}
	}

	return nil
}

// 统一消息发放
func notify(conn *websocket.Conn, msg string) {
	chNotify <- 1 // 利用channel阻塞 避免并发去对同一个连接发送消息出现panic: concurrent write to websocket connection这样的异常
	_, roomIdInt := getRoomId()
	assignRoom := rooms[roomIdInt]
	for _, con := range assignRoom {
		if con.(wsClients).RemoteAddr != conn.RemoteAddr().String() {
			con.(wsClients).Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
	<-chNotify
}

// 离线通知
func disconnect(conn *websocket.Conn) {
	_, roomIdInt := getRoomId()

	objColl := collection.NewObjCollection(rooms[roomIdInt])

	retColl := safe.Safety.Do(func() interface{} {
		return objColl.Reject(func(item interface{}, key int) bool {
			if item.(wsClients).RemoteAddr == conn.RemoteAddr().String() {

				data := msgData{
					Username: item.(wsClients).Username,
					Uid:      item.(wsClients).Uid,
					Time:     time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()
				}

				jsonStrServeMsg := msg{
					Status: msgTypeOffline,
					Data:   data,
				}
				serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

				disMsg := string(serveMsgStr)

				item.(wsClients).Conn.Close()

				notify(conn, disMsg)

				return true
			}
			return false
		})
	}).(collection.ICollection)

	interfaces, _ := retColl.ToInterfaces()
	rooms[roomIdInt] = interfaces
}

// 格式化传送给客户端的消息数据
func formatServeMsgStr(status int, conn *websocket.Conn) ([]byte, msg) {

	roomId, roomIdInt := getRoomId()

	//log.Println(reflect.TypeOf(var))

	data := msgData{
		Username: clientMsg.Data.Username,
		Uid:      clientMsg.Data.Uid,
		RoomId:   roomId,
		Time:     time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()
	}

	if status == msgTypeSend || status == msgTypePrivateChat {
		data.AvatarId = clientMsg.Data.AvatarId
		content := clientMsg.Data.Content

		data.Content = content
		if helper.MbStrLen(content) > 800 {
			// 直接截断
			data.Content = string([]rune(content)[:800])
		}

		toUidStr := clientMsg.Data.ToUid
		toUid, _ := strconv.Atoi(toUidStr)

		// 保存消息
		stringUid := data.Uid
		intUid, _ := strconv.Atoi(stringUid)

		if clientMsg.Data.ImageUrl != "" {
			// 存在图片
			models.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"content":    data.Content,
				"room_id":    data.RoomId,
				"image_url":  clientMsg.Data.ImageUrl,
			})
		} else {
			models.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"content":    data.Content,
				"room_id":    data.RoomId,
			})
		}

	}

	if status == msgTypeGetOnlineUser {
		ro := rooms[roomIdInt]
		data.Count = len(ro)
		data.List = ro
	}

	jsonStrServeMsg := msg{
		Status: status,
		Data:   data,
		Conn:   conn,
	}
	serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

	return serveMsgStr, jsonStrServeMsg
}

func getRoomId() (string, int) {
	roomId := clientMsg.Data.RoomId

	roomIdInt, _ := strconv.Atoi(roomId)
	return roomId, roomIdInt
}

// =======================对外方法=====================================

func GetOnlineUserCount() int {
	num := 0
	for i := 1; i <= roomCount; i++ {
		num = num + GetOnlineRoomUserCount(i)
	}
	return num
}

func GetOnlineRoomUserCount(roomId int) int {
	return len(rooms[roomId])
}
