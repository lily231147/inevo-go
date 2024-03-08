package pojo

// Component 子件信息
type Component struct {
	ComponentId   string
	ComponentName string
	Qty           float32
}

// S 母件到子件
type S map[string][]Component

// T 用于标识物料是否为原材料
type T map[string]string

type BomInfo struct {
	BomCode     string `json:"id"`
	BomName     string `json:"name"`
	MoCode      string `json:"mocode"`
	ProductCode string `json:"pcode"`
	ProductName string `json:"pname"`
}

// Parent 母件信息
//type Parent struct {
//	ParentId   string
//	ParentName string
//}
