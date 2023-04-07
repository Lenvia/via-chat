package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"via-chat/controller"
	"via-chat/middleware"
	"via-chat/services/session"
	"via-chat/static"
	"via-chat/ws/primary"
)

func InitRoute() *gin.Engine {
	//router := gin.Default()
	router := gin.New()
	router.Use(middleware.Cors())

	if viper.GetString(`app.debug_mod`) == "false" {
		// live 模式 打包用，使用嵌入式的静态资源FS
		router.StaticFS("/static", http.FS(static.EmbedStatic))
	} else {
		// dev 开发用 避免修改静态资源需要重启服务，使用本地的静态文件目录
		router.StaticFS("/static", http.Dir("static"))
	}

	// 创建路由分组，并启用 cookie-based 会话
	sr := router.Group("/", session.EnableCookieSession())
	{
		sr.GET("/", controller.Index)

		sr.POST("/login", controller.Login)
		sr.GET("/logout", controller.Logout)
		sr.GET("/ws", primary.Start)

		authorized := sr.Group("/") // 使用 AuthSessionMiddle() 中间件检查用户是否已登录
		authorized.Use(session.AuthSessionMiddle())
		{
			//authorized.GET("/ws", ws.Run)
			authorized.GET("/home", controller.Home)
			authorized.GET("/room/:room_id", controller.Room)
			authorized.GET("/private-chat", controller.PrivateChat)
			authorized.POST("/img-kr-upload", controller.ImgKrUpload)
			authorized.GET("/pagination", controller.Pagination)
		}

	}

	return router
}
