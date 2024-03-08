package dao

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"jcyh/pojo"
	"jcyh/util"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// GetOrder 根据订单编码查找订单信息，订单可能包括多个子订单
// 子订单：MomOrderDetails{ModId, Qty,  RootId,  RootName}
// ModId：子订单id
// Qty：子订单产品数量
// RootId : 子订单产品id
// RootName ： 子订单产品名称
func GetOrder(Code string, db *gorm.DB) (momOrder pojo.MomOrder) {
	momOrder = pojo.MomOrder{Code: Code}
	db.Table("order_").Select([]string{"ModId", "qty", "ProductId", "ProductName"}).Where("code = ?", Code).Find(&(momOrder.OrderDetails))
	return
}

// 获得订单的价格等等信息
func GetOrderInfo(Code string) (orderInfo pojo.MoSearchInfo) {
	util.Db.Table("order_search").Where("code='" + Code + "'").Find(&orderInfo)
	return
}

// 获得物料的价格等等信息
func GetMaterialInfo() (materialInfos []pojo.MaterialInfo) {
	util.Db.Table("material").Limit(100).Find(&materialInfos)
	//util.Db.Table("bom_struct a").Select([]string{"min(a.component_id) component_id", "a.component_name", "cast(sum(qty) as signed) as qty",
	//	"cast(cast(conv(Right(md5(a.component_name),3),16,10) as signed)/(cast(sum(qty) as signed)+10) as signed)+cast(conv(Right(md5(a.component_name),2),16,10) as signed) as price"}).
	//	Where("not exists(select * from bom_struct b where a.component_id=b.parent_id)").Group("a.component_name").Limit(100).Find(&materialInfos)
	return
}

// GetSailInfo
// 以指定时间单位（月| 季度）和产品加载销售信息，name是产品名，如果为all的话是查所有的
func GetSailInfo(name string) (dateRange pojo.DateRange, info [][]pojo.SailInfo) {
	//读取季度范围
	if dr, err := util.Redis.Get(util.Ctx, "sail_"+name).Result(); err != redis.Nil {
		json.Unmarshal([]byte(dr), &dateRange)
	} else {
		util.Db.Table("sail_info").Select([]string{"min(quarter) min", "max(quarter) max"}).Find(&dateRange)
		//转化为json
		jsi, _ := json.Marshal(dateRange)
		util.Redis.Set(util.Ctx, "sail_"+name, jsi, 0)
	}
	//指定日期和产品大类名称，将查询得到结果按地区分组
	for date := dateRange.Min; strings.Compare(date, dateRange.Max) != 1; date = dateAdd(date, "quarter") {
		var sailInfo []pojo.SailInfo
		if redisCache, err := util.Redis.Get(util.Ctx, "sail_"+date+"_"+name).Result(); err != redis.Nil {
			json.Unmarshal([]byte(redisCache), &sailInfo)
		} else {
			util.Db.Table("sail").Select([]string{"sum(qty) qty", "sum(money) money", "sum(profit) profit", "Address address"}).
				Where("quarter=? and name like '%"+name+"%'", date).Group("address").Find(&sailInfo)
			//转化为json
			redisCache, _ := json.Marshal(sailInfo)
			util.Redis.Set(util.Ctx, "sail_"+date+"_"+name, redisCache, 0)
		}
		info = append(info, sailInfo)
	}
	return
}

// GetMoCodeSearch
// 根据输入的信息去查找候选订单号
func GetOrderSearchResult(code string, place string, startDate, endDate string, name string, minMoney string, maxMoney string, isFault bool) (moSearchInfo []pojo.MoSearchInfo) {
	condition := "1=1"
	if isFault {
		condition = "is_fault=1"
	}
	if place != "" {
		condition = condition + " and place like '%" + place + "%'"
	}
	if code != "" {
		condition = condition + " and code='" + code + "'"
		util.Db.Table("order_search").Where(condition).Find(&moSearchInfo)
		return
	}
	if startDate != "" {
		condition = condition + " and date > '" + startDate + "'"
	}
	if endDate != "" {
		condition = condition + " and date < '" + endDate + "'"
	}
	if name != "" {
		condition = condition + " and name like '%" + name + "%'"
	}
	if minMoney != "" {
		condition = condition + " and money > '" + minMoney + "'"
	}
	if maxMoney != "" {
		condition = condition + " and money < '" + maxMoney + "'"
	}
	util.Db.Table("order_search").Where(condition).Find(&moSearchInfo)
	return
}

// GetMomOrderByBatch
// 根据mocode获得同批次的其他订单
func GetMomOrderByBatch(moCode string) (moSearchInfo []pojo.MoSearchInfo) {
	if bi, err := util.Redis.Get(util.Ctx, "batch"+moCode).Result(); err != redis.Nil {
		json.Unmarshal([]byte(bi), &moSearchInfo)
	} else {
		sql := "batch_id in (select batch_id from order_search a where code=?)"
		util.Db.Table("order_search").Select([]string{"code", "place", "date", "name", "money"}).Where(sql, moCode).Find(&moSearchInfo)
		//转换为json
		jvs, _ := json.Marshal(moSearchInfo)
		util.Redis.Set(util.Ctx, "batch"+moCode, jvs, 0)
	}
	return
}

// 日期字符串加法实现
func dateAdd(current string, period string) (next string) {
	var max int
	if period == "month" {
		max = 12
	} else {
		max = 4
	}
	strs := strings.Split(current, "-")
	monthOrQuarter, _ := strconv.Atoi(strs[1])
	if max > monthOrQuarter {
		strs[1] = strconv.Itoa(monthOrQuarter + 1)
	} else {
		year, _ := strconv.Atoi(strs[0])
		strs[0] = strconv.Itoa(year + 1)
		strs[1] = "1"
	}
	next = strs[0] + "-" + strs[1]
	return
}

// GetFaultImgPaths
// 根据关键词查询相关故障图片的路径
func GetFaultImgPaths(key string) (paths []string) {
	if redisCache, err := util.Redis.Get(util.Ctx, "fault_img_"+key).Result(); err != redis.Nil {
		json.Unmarshal([]byte(redisCache), &paths)
	} else {
		util.Db.Table("fault_img").Select("path").Where("keyword like '%" + key + "%'").Find(&paths)
		//转换为json
		redisCache, _ := json.Marshal(paths)
		util.Redis.Set(util.Ctx, "fault_img_"+key, redisCache, 0)
	}
	return
}

// GetFaultImgPaths
// 根据关键词查询相关故障图片的路径
func GetFaultImgPathsByCode(code string) (path string) {
	var paths []string
	if redisCache, err := util.Redis.Get(util.Ctx, "fault_img_by_code"+code).Result(); err != redis.Nil {
		json.Unmarshal([]byte(redisCache), &path)
	} else {
		util.Db.Table("fault_img").Select("path").Find(&paths)
		hash := sha256.Sum256([]byte(code))
		num, _ := strconv.ParseInt(hex.EncodeToString(hash[:])[:8], 16, 0)
		path = paths[num%int64(len(paths))]
		//转换为json
		redisCache, _ := json.Marshal(path)
		util.Redis.Set(util.Ctx, "fault_img_by_code"+code, redisCache, 0)
	}
	return
}
