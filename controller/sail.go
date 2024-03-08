package controller

import (
	"jcyh/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SailInfo(c *gin.Context) {
	productName := c.Query("productname")
	dateRange, info := dao.GetSailInfo(productName)
	c.IndentedJSON(http.StatusOK, gin.H{"dateRange": dateRange, "info": info})
}

func FaultImg(c *gin.Context) {
	key := c.Query("key")
	paths := dao.GetFaultImgPaths(key)
	c.IndentedJSON(http.StatusOK, paths)
}

func FaultImgByCode(c *gin.Context) {
	key := c.Query("code")
	paths := dao.GetFaultImgPathsByCode(key)
	c.IndentedJSON(http.StatusOK, paths)
}
