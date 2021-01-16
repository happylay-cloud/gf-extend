package gfreq

// PageReq 分页请求参数
type PageReq struct {
	PageNum  int `d:"1"  json:"pageNum"`  // 当前页码
	PageSize int `d:"10" json:"pageSize"` // 每页条数
}
