package dao

import (
	"encoding/json"
	"jcyh/pojo"
	"jcyh/util"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// GetVendorsByOrderCode
// 根据订单编码获取供应商信息
func GetVendorsByOrderCode(code string) (vendors []pojo.Vendor) {

	if vs, err := util.Redis.Get(util.Ctx, "venbymo"+code).Result(); err != redis.Nil {
		json.Unmarshal([]byte(vs), &vendors)
	} else {
		var productIds []struct{ ProductId int }
		util.Db.Table("order_").Select("product_id").Where("code='" + code + "'").Find(&productIds)

		var t []pojo.Vendor
		for _, item := range productIds {
			util.Db.Raw("call get_vendors(" + strconv.Itoa(item.ProductId) + ")").Find(&t)
			vendors = append(vendors, t...)
		}
		//转换为json
		jvs, _ := json.Marshal(vendors)
		util.Redis.Set(util.Ctx, "venbymo"+code, jvs, 0)
	}
	return
}

// GetVendorAvatarByCode
// 根据供应商id拿到供应商画像
func GetVendorAvatarByCode(code string) (vendorAvatar pojo.VendorAvatar) {
	if redisCache, err := util.Redis.Get(util.Ctx, "vendor_avatar_"+code).Result(); err != redis.Nil {
		json.Unmarshal([]byte(redisCache), &vendorAvatar)
	} else {
		util.Db.Table("vendor").Where("code='" + code + "'").First(&vendorAvatar)
		//转换为json
		redisCache, _ := json.Marshal(vendorAvatar)
		util.Redis.Set(util.Ctx, "vendor_avatar_"+code, redisCache, 0)
	}
	return
}

// GetVendorAvatarByBomId
// 根据原材料id拿到供应商画像
func GetVendorAvatarByBomId(id string) (vendorAvatar pojo.VendorAvatar) {
	var code string
	util.Db.Table("vendor_bom").Select("code").Where("component_id=" + id).Find(&code)
	return GetVendorAvatarByCode(code)
}

// GetVendorAll
// 根据原材料id拿到供应商画像
func GetVendorAll() (vendorAvatars []pojo.VendorAvatar) {
	util.Db.Table("vendor").Find(&vendorAvatars)
	return
}
