package pojo

// MomOrderDetail 订单明细，一个订单可以有多个订单明细，每个订单明细对应部分订单需求
type MomOrderDetail struct {
	ModId       string
	Qty         float32
	ProductId   string
	ProductName string
}

type DateRange struct {
	Min string
	Max string
}
