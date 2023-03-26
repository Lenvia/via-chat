package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"via-chat/models"
)

// EnableCookieSession 函数用于启用Cookie Session支持
// 它创建了一个名为 "via-chat" 的Cookie Session存储器，其中包含应用程序的 cookie_key
func EnableCookieSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte(viper.GetString(`app.cookie_key`)))
	return sessions.Sessions("via-chat", store)
}

// SaveAuthSession 注册和登陆时都需要保存 session 信息
func SaveAuthSession(c *gin.Context, info interface{}) {
	// 在请求上下文中创建默认Session存储器的实例
	session := sessions.Default(c)
	session.Set("uid", info)
	// c.SetCookie("user_id",string(info.(map[string]interface{})["b"].(uint)), 1000, "/", "localhost", false, true)
	session.Save() // 将session的更改保存到底层存储
}

// GetSessionUserInfo 函数用于从Session中检索用户信息，以便进行验证或处理用户请求
func GetSessionUserInfo(c *gin.Context) map[string]interface{} {
	session := sessions.Default(c)

	uid := session.Get("uid")

	data := make(map[string]interface{})
	// 如果Session中存在用户ID
	if uid != nil { // 使用此ID检索用户信息，例如：ID，用户名和头像编号
		user := models.FindUserByField("id", uid.(string))
		data["uid"] = user.ID
		data["username"] = user.Username
		data["avatar_id"] = user.AvatarId
	}
	return data
}

// ClearAuthSession 函数用于清除所有与Session相关的数据，包括login和logout过程中所需的信息
func ClearAuthSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

// HasSession 函数用于检查请求中是否包含有效的Session信息
func HasSession(c *gin.Context) bool {
	session := sessions.Default(c)
	if sessionValue := session.Get("uid"); sessionValue == nil {
		return false
	}
	return true
}

// AuthSessionMiddle 函数是用于检查用户是否已经进行身份验证的gin中间件。
// 它检查请求中是否存在有效的Session信息。如果没有Session信息，则将用户重定向到首页。
// 如果Session信息有效，则将Session中的值设置为uid，然后执行链中下一个处理程序。
func AuthSessionMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionValue := session.Get("uid")
		if sessionValue == nil {
			c.Redirect(http.StatusFound, "/")
			return
		}

		uidInt, _ := strconv.Atoi(sessionValue.(string))

		if uidInt <= 0 {
			c.Redirect(http.StatusFound, "/")
			return
		}

		// 将Session值设置为请求上下文中的变量
		c.Set("uid", sessionValue)

		c.Next()
		return
	}
}
