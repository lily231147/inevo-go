package main

import (
	"jcyh/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	// initialize a route using Default
	router := gin.Default()

	router.Use(CrosHandler())
	// map the request to handler
	// 查询物料信息
	router.GET("/data/bom", controller.BomInfo)
	// 查询原材料信息
	router.GET("/data/material", controller.MaterialInfo)
	// 查询销售信息
	router.GET("/data/sail_info", controller.SailInfo)
	// 查询订单|故障订单
	router.GET("/data/order_search", controller.OrderSearch)
	// 查询供应商和同批次订单
	router.GET("/data/vendor_and_batch", controller.VendorsAndBatch)
	// 查询供应商画像
	router.GET("/data/vendor_avatar", controller.VendorAvatar)
	// 查询所有供应商画像
	router.GET("/data/vendor_all", controller.VendorAll)
	// 查询故障图片
	router.GET("/data/fault_img", controller.FaultImg)
	// 查询指定订单的故障图片
	router.GET("/data/fault_img_by_code", controller.FaultImgByCode)
	// attach the route to a http.Server and start the server
	router.Run("0.0.0.0:8080")
}

func CrosHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		context.Header("Access-Control-Max-Age", "172800")
		context.Header("Access-Control-Allow-Credentials", "false")
		context.Set("content-type", "application/json") //设置返回格式是json
		//处理请求
		context.Next()
	}
}
