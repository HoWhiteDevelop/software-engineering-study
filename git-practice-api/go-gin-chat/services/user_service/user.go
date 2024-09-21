package user_service

import (
	"git-practice-api/go-gin-chat/models"
	"git-practice-api/go-gin-chat/result"
	"git-practice-api/go-gin-chat/services/helper"
	"git-practice-api/go-gin-chat/services/session"
	"git-practice-api/go-gin-chat/services/validator"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Login(c *gin.Context) {

	var u validator.User

	if err := c.ShouldBind(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5000, "msg": err.Error()})
		return
	}

	user := models.FindUserByField("username", u.Username)
	userInfo := user
	md5Pwd := helper.Md5Encrypt(u.Password)

	if userInfo.ID > 0 {
		// json 用户存在
		// 验证密码
		if userInfo.Password != md5Pwd {
			result.Failture(result.APIcode.PasswordError, result.APIcode.GetMessage(result.APIcode.PasswordError), c, nil)
			return
		}

		models.SaveAvatarId(u.AvatarId, user)

	} else {
		// 新用户
		userInfo = models.AddUser(map[string]interface{}{
			"username":  u.Username,
			"password":  md5Pwd,
			"avatar_id": u.AvatarId,
		})
	}

	if userInfo.ID > 0 {
		session.SaveAuthSession(c, string(strconv.Itoa(int(userInfo.ID))))
		result.Success("注册成功", c)
		return
	} else {
		result.Failture(result.APIcode.SystemError, result.APIcode.GetMessage(result.APIcode.SystemError), c, nil)
		return
	}
}

func GetUserInfo(c *gin.Context) map[string]interface{} {
	return session.GetSessionUserInfo(c)
}

func Logout(c *gin.Context) {
	session.ClearAuthSession(c)
	c.Redirect(http.StatusFound, "/")
	return
}
