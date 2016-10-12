package database

//分页参数结构体
type Page struct {
	Pages   int `json:"pages"`   //页码数量
	Totle   int `json:"totle"`   //总数据
	Limit   int `json:"limit"`   //每页数量
	Current int `json:"current"` //当前页
}
