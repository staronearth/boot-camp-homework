package web

import "github.com/gin-gonic/gin"

func RegisterRoutes() *gin.Engine {
	server := gin.Default()
	RegisterUserRoutes(server)
	return server
}

func RegisterUserRoutes(server *gin.Engine) {
	u := &UserHandler{}
	server.POST("/users/signup", u.Signup)
	//这是REST风格
	//server.PUT("/users", u.Signup)

	server.POST("/users/login", u.Login)

	server.POST("/users/edit", u.Edit)
	//REST风格
	//server.POST("/users/:id", u.Edit)
	server.GET("/users/profile", u.Profile)
	//REST风格
	//server.GET("/users/:id", u.Profile)
}
