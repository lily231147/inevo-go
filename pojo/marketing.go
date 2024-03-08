package pojo

// MomOrder 订单，由订单编码MoCode表示
type MomOrder struct {
	Code         string
	OrderDetails []MomOrderDetail
}

// SailInfo 销售信息
type SailInfo struct {
	Qty     int
	Money   int
	Profit  int
	Address string
}

type Date2SailInfo map[string][]SailInfo

// MoSearchInfo 查找得到的订单信息
type MoSearchInfo struct {
	Code  string
	Place string
	Date  string
	Name  string
	Money float32
}

type MaterialInfo struct {
	ComponentId   string
	ComponentName string
	Qty           int
	Price         int
}
