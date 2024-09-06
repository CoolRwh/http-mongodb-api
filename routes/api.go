package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"http-mongodb-api/app/controller"
	"net/http"
)

var Router *gin.Engine
var Api = new(apiRouter)

type apiRouter struct{}

func (apiRouter *apiRouter) InitRouter() {
	Router = gin.Default()
	//_ = Router.SetTrustedProxies([]string{"127.0.0.1"})
	//跨域中间件
	Router.Use(cors())

	Router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, "ok")
		return
	})

	api := Router.Group("api/v1")
	{
		mongodb := api.Group("mongodb")
		{
			mongodb.POST("cs", controller.Mongodb.Cs)
			mongodb.POST("find", controller.Mongodb.Find)
			mongodb.POST("fineOne", controller.Mongodb.FindOne)
			mongodb.POST("installMany", controller.Mongodb.InsertMany)
			mongodb.POST("installOne", controller.Mongodb.InsertOne)
			mongodb.POST("updateById", controller.Mongodb.UpdateById)
			mongodb.POST("updateOne", controller.Mongodb.UpdateOne)
			mongodb.POST("deleteById", controller.Mongodb.DeleteById)
			mongodb.POST("deleteMany", controller.Mongodb.DeleteMany)
		}
	}
	if err := Router.Run(viper.GetString("server.addr")); err != nil {
		fmt.Printf("error:" + err.Error())
	}
}
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
