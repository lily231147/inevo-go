package dao

import (
	"encoding/json"
	"jcyh/pojo"
	"jcyh/util"

	"github.com/go-redis/redis/v8"
)

// GetAllBomByOrder
// 根据输入的订单查询所有子件及其结构
// 返回母件与子件的结构map（pojo.S）以及物料id到物料name的map（pojo.N）
func GetAllBomByOrder(MomOrder pojo.MomOrder) (s pojo.S) {
	//以母件id为key，以子件信息为value
	s = make(pojo.S)
	//读取所有的根母件作为键，将对应的子件作为值，得到map
	for _, OrderDetail := range MomOrder.OrderDetails {
		searchComponent(s, OrderDetail.ProductId)
	}
	return
}

// SearchComponent
// 根据母件id查找所有子件信息，然后再递归的查找下一层子件信息
func searchComponent(s pojo.S, parentId string) {
	var cs []pojo.Component
	if child, err := util.Redis.Get(util.Ctx, parentId).Result(); err != redis.Nil {
		json.Unmarshal([]byte(child), &cs)
	} else {
		util.Db.Table("bom_struct").Select([]string{"ComponentId", "ComponentName", "Qty"}).Where("parent_id = ?", parentId).Find(&cs)
		//转换为json
		jcs, _ := json.Marshal(cs)
		util.Redis.Set(util.Ctx, parentId, jcs, 0)
	}
	if len(cs) == 0 {
		s[parentId] = nil
		return
	}
	s[parentId] = cs
	for _, c := range cs {
		searchComponent(s, c.ComponentId)
	}
}
