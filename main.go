package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello，启动成功了！")
	})

	//db := initDB()
	//server := initWebServer()
	//initUserHdl(db, server)

	server.Run(":8080")
}

//func initUserHdl(db *gorm.DB,
//	redisClient redis.Cmdable,
//	codeSvc service.CodeService,
//	server *gin.Engine) {
//	ud := dao.NewUserDAO(db)
//	uc := cache.NewUserCache(redisClient)
//	ur := repository.NewCachedUserRepository(ud, uc)
//	us := service.NewUserService(ur)
//	hdl := web.NewUserHandler(us, codeSvc)
//	hdl.RegisterRoutes(server)
//}
//
//func initDB() *gorm.DB {
//	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
//	if err != nil {
//		panic(err)
//	}
//	err = dao.InitTables(db)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}
//
//func initCodeSvc(redisClient redis.Cmdable) *service.CodeService {
//	cc := cache.NewCodeCache(redisClient)
//	crepo := repository.NewCodeRepository(cc)
//	return service.NewCodeService(crepo, initMemorySms())
//}
//
//func initMemorySms() sms.Service {
//	return localsms.NewService()
//}
//
//func initWebServer() *gin.Engine {
//	server := gin.Default()
//
//	server.Use(cors.New(cors.Config{
//		//AllowOrigins:     []string{"http://localhost:3000"},
//		AllowCredentials: true,
//		AllowHeaders:     []string{"Content-Type", "Authorization"},
//		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
//		AllowOriginFunc: func(origin string) bool {
//			if strings.HasPrefix(origin, "http://localhost") {
//				//if strings.Contains(origin, "localhost") {
//				return true
//			}
//			return strings.Contains(origin, "webook.com")
//		},
//		MaxAge: 6 * time.Hour,
//	}),
//		func(ctx *gin.Context) {
//			println("NOTE: Implement Another Middleware")
//		})
//
//	//redisClient := redis.NewClient(&redis.Options{
//	//	Addr: config.Config.Redis.Addr,
//	//})
//
//	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 1).Build())
//
//	//useSession(server)
//	useJWT(server)
//
//	return server
//}
//
//func useJWT(server *gin.Engine) {
//	login := middleware.LoginJWTMiddlewareBuilder{}
//	server.Use(login.CheckLogin())
//}

//func useSession(server *gin.Engine) {
//	login := &middleware.LoginMiddlewareBuilder{}
//	// 1. userId is saved in Cookie
//	store := cookie.NewStore([]byte("secret"))
//	// or 2. save in memory
//	// store := memstore.NewStore([]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"),
//	//	[]byte("eF1`yQ9>yT1`tH1,sJ0.zD8;mZ9~nC6("))
//	// or 3. save in redis
//	//store, err := redis.NewStore(16, "tcp",
//	//	"localhost:6379", "",
//	//	[]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"),
//	//	[]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgA"))
//	//if err != nil {
//	//	panic(err)
//	//}
//	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
//}
