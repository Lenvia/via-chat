package user_service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"via-chat/models"
	"via-chat/services/helper"
	"via-chat/services/session"
	"via-chat/services/validator"
)

func Login(c *gin.Context) {
	var u validator.User

	if err := c.ShouldBind(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5000, "msg": err.Error()})
		return
	}

	username := u.Username
	pwd := u.Password
	avatarId := u.AvatarId

	user := models.FindUserByField("username", username)
	userInfo := user
	md5Pwd := helper.Md5Encrypt(pwd)

	if userInfo.ID > 0 {
		// json 用户存在
		// 验证密码
		if userInfo.Password != md5Pwd {
			c.JSON(http.StatusOK, gin.H{
				"code": 5000,
				"msg":  "密码错误",
			})
			return
		}

		models.SaveAvatarId(avatarId, user)

	} else {
		// 新用户
		userInfo = models.AddUser(map[string]interface{}{
			"username":  username,
			"password":  md5Pwd,
			"avatar_id": avatarId,
		})
	}

	if userInfo.ID > 0 {
		// 登录成功，将用户信息保存到会话中
		session.SaveAuthSession(c, strconv.Itoa(int(userInfo.ID)))
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		// 到这里session是有值的
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 5001,
			"msg":  "系统错误",
		})
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