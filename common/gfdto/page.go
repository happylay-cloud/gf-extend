package gfdto

import "github.com/gogf/gf/util/gconv"

// 通用分页JSON数据结构
type page struct {
	PageNum    int         `json:"pageNum"`    // 当前页码
	PageSize   int         `json:"pageSize"`   // 每页记录数
	TotalCount int         `json:"totalCount"` // 数据总数
	List       interface{} `json:"list"`       // 数据集合
}

// PageData 封装分页数据
func PageData(pageNum, pageSize, totalCount int, list interface{}) *page {
	return &page{
		PageNum:    pageNum,
		PageSize:   pageSize,
		TotalCount: totalCount,
		List:       list,
	}
}

// NewPage 获取分页对象
func NewPage(pageNum, pageSize int) *page {
	return &page{
		PageNum:    pageNum,
		PageSize:   pageSize,
		TotalCount: 0,
		List:       nil,
	}
}

// LimitPage 获取原生sql分页参数及sql语句
func (p *page) LimitPage() (page, pageSize int, sql string) {

	// 数据处理
	if p.PageSize <= 0 {
		p.PageSize = 1
	}
	if p.PageNum <= 0 {
		p.PageNum = 10
	}
	// 计算分页数据
	page = (p.PageNum - 1) * p.PageSize
	// 每页大小
	pageSize = p.PageSize
	// 分页sql语句
	sql = " LIMIT " + gconv.String(page) + "," + gconv.String(pageSize)

	return page, pageSize, sql
}
