package hdto

import (
	"math"

	"github.com/gogf/gf/util/gconv"
)

// Page 通用分页JSON数据结构
type Page struct {
	PageNum    int         `json:"pageNum"`    // 当前页码
	PageSize   int         `json:"pageSize"`   // 每页记录数
	TotalCount int         `json:"totalCount"` // 数据总数
	PageCount  int         `json:"pageCount"`  // 总分页数
	List       interface{} `json:"list"`       // 数据集合
}

// PageData 封装分页数据
func PageData(pageNum, pageSize, totalCount int, list interface{}) *Page {
	return &Page{
		PageNum:    pageNum,
		PageSize:   pageSize,
		TotalCount: totalCount,
		List:       list,
	}
}

// NewPage 获取分页对象
//	@pageNum			当前页码
//	@pageSize			每页记录数
//	@handlerErrData 	是否处理异常数据，可选参数，bool类型，false 不处理，true 处理。默认true处理。
func NewPage(pageNum, pageSize int, handlerErrData ...bool) *Page {

	// 默认处理异常数据
	if len(handlerErrData) == 0 || (len(handlerErrData) > 0 && handlerErrData[0]) {
		// 异常数据处理
		if pageNum <= 0 {
			pageNum = 1
		}

		if pageSize <= 0 {
			pageSize = 10
		}
	}

	return &Page{
		PageNum:    pageNum,
		PageSize:   pageSize,
		TotalCount: 0,
		List:       nil,
	}
}

// LimitPage 获取原生sql分页参数及sql语句
func (p *Page) LimitPage() (page, pageSize int, sql string) {

	// 计算分页数据
	page = (p.PageNum - 1) * p.PageSize
	// 每页大小
	pageSize = p.PageSize
	// 分页sql语句
	sql = " LIMIT " + gconv.String(page) + "," + gconv.String(pageSize)

	return page, pageSize, sql
}

// RedisPage 获取Redis分页
func (p *Page) RedisPage() (start, stop int) {

	// 计算分页数据-对应redis start
	start = (p.PageNum - 1) * p.PageSize
	// 每页大小-对应redis stop
	stop = start + p.PageSize - 1

	return start, stop
}

// CalPageCount 计算分页总数
func (p *Page) CalPageCount(totalCount int) (pageCount int) {

	// 计算分页总数
	pageCount = int(math.Ceil(float64(totalCount) / float64(p.PageSize)))

	return pageCount
}
