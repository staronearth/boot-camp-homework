package main

import (
	"boot-camp-homework/homework/week2/webhook/internal/repository"
	dao "boot-camp-homework/homework/week2/webhook/internal/repository/dao"
	"boot-camp-homework/homework/week2/webhook/internal/service"
	"boot-camp-homework/homework/week2/webhook/internal/web"
	"boot-camp-homework/homework/week2/webhook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	u := initUser(db)
	server := initWebServer()
	u.RegisterRoutes(server)
	server.Run(":8090")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(func(context *gin.Context) {
		println("这是第一个middleware")
	})
	server.Use(func(context *gin.Context) {
		println("这是第二个middleware")
	})
	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods: []string{"PUT", "PATCH", "OPTION", "POST"},
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"x-jwt-token"},
		//是否允许携带cookie之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "staronearth")
		},
		MaxAge: 12 * time.Minute,
	}))
	//步骤1,放在cookie中
	//store := cookie.NewStore([]byte("secret"))
	//步骤2,放在内存中
	//store := memstore.NewStore([]byte("Mq.8f2@UErXH1EE4a`h{|:bnf#2>Xbi6"), []byte("g7tt&Z:0yz{9,6q_m4zbVf95+=YF/>FG"))
	//公司中的
	//store, err := redis.NewStore(16, "tcp", "172.17.32.110:32381", "", []byte("Mq.8f2@UErXH1EE4a`h{|:bnf#2>Xbi6"), []byte("g7tt&Z:0yz{9,6q_m4zbVf95+=YF/>FG"))
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("Mq.8f2@UErXH1EE4a`h{|:bnf#2>Xbi6"), []byte("g7tt&Z:0yz{9,6q_m4zbVf95+=YF/>FG"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().IgnorePaths("/users/login").IgnorePaths("/users/signup").Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:example@tcp(localhost:3306)/webook"), &gorm.Config{})
	//db, err := gorm.Open(mysql.Open("root:devops@tcp(172.17.32.110:32658)/webook"), &gorm.Config{})
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic相当于整个goroutine结束
		// 一旦初始化过程出错,应用就不要启动了
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
