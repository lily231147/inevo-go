package controller

import (
	"jcyh/dao"
	"jcyh/pojo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VendorsAndBatch(c *gin.Context) {
	code := c.Query("order_code")
	vs := dao.GetVendorsByOrderCode(code)
	bat := dao.GetMomOrderByBatch(code)
	c.IndentedJSON(http.StatusOK, gin.H{"vendors": vs, "batch": bat})
}

func VendorAvatar(c *gin.Context) {
	var avatar pojo.VendorAvatar
	if code := c.DefaultQuery("vencode", ""); code != "" {
		avatar = dao.GetVendorAvatarByCode(code)
	} else {
		id := c.DefaultQuery("bomid", "")
		avatar = dao.GetVendorAvatarByBomId(id)
	}
	c.IndentedJSON(http.StatusOK, avatar)
}

func VendorAll(c *gin.Context) {
	var avatars = dao.GetVendorAll()
	c.IndentedJSON(http.StatusOK, avatars)
}
