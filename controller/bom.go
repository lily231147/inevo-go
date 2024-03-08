package controller

import (
	"jcyh/dao"
	"jcyh/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func BomInfo(c *gin.Context) {
	code := c.Query("orderid")
	order := dao.GetOrder(code, util.Db)
	s := dao.GetAllBomByOrder(order)
	info := dao.GetOrderInfo(code)
	c.IndentedJSON(http.StatusOK, gin.H{"order": order, "struct": s, "info": info})
}

func MaterialInfo(c *gin.Context) {
	info := dao.GetMaterialInfo()
	c.IndentedJSON(http.StatusOK, info)
}

func OrderSearch(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	place := c.DefaultQuery("place", "")
	name := c.DefaultQuery("name", "")
	dateRange := c.DefaultQuery("dateRange", "")
	startDate := strings.Split(dateRange, ",")[0]
	endDate := strings.Split(dateRange, ",")[1]
	minMoney := c.DefaultQuery("minMoney", "")
	maxMoney := c.DefaultQuery("maxMoney", "")
	fault := c.DefaultQuery("fault", "false")
	isFault, _ := strconv.ParseBool(fault)
	results := dao.GetOrderSearchResult(code, place, startDate, endDate, name, minMoney, maxMoney, isFault)
	c.IndentedJSON(http.StatusOK, results)
}
