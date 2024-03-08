package pojo

type Vendor struct {
	Code string `json:"id"`
	Name string `json:"name"`
}

type VendorAvatar struct {
	Code     string
	Name     string
	Turnover float32
	Delivery float32
	Quantity float32
	Service  float32
	Scale    float32
	Region   string
}
