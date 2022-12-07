package main

import (
	"fmt"

	"github.com/aZ4ziL/drug-api/handlers"
	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("server is running...")
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	userGroup := r.Group("/v1/user")
	userGroup.POST("/sign-up", handlers.NewUserHandlerV1().SignUp())
	userGroup.POST("/get-token", handlers.NewUserHandlerV1().GetToken())

	drugGroup := r.Group("/v1/drug")
	drugGroup.Use(handlers.AuthenticationMiddleware())
	drugGroup.GET("", handlers.NewDrugHandlerV1().Index())
	drugGroup.POST("", handlers.NewDrugHandlerV1().Add())
	drugGroup.PUT("", handlers.NewDrugHandlerV1().Edit())
	drugGroup.DELETE("", handlers.NewDrugHandlerV1().Delete())

	r.Run(":8000")
}
