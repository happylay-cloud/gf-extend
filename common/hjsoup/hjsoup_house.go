package hjsoup

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"

	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

// RecordHousePriceDTO 商品住宅明码标价信息
type RecordHousePriceDTO struct {
	HrefId        string `json:"href_id"`        // 跳转详情ID
	RecordNo      string `json:"record_no"`      // 备案号
	EstateName    string `json:"estate_name"`    // 楼盘名称
	HouseNumbers  string `json:"house_numbers"`  // 楼号
	HouseArea     string `json:"house_area"`     // 建筑面积
	HouseNum      string `json:"house_num"`      // 套数
	AvgPrice      string `json:"avg_price"`      // 均价
	ReleaseTime   string `json:"release_time"`   // 发布日期
	ReleaseSource string `json:"release_source"` // 来源
}

// RecordHouseDetailDTO 户型详情信息
type RecordHouseDetailDTO struct {
	RecordNo      string `json:"record_no"`      // 备案号
	HouseNumber   string `json:"house_number"`   // 楼号
	RoomNumber    string `json:"room_number"`    // 房号
	HouseType     string `json:"house_type"`     // 户型
	HouseArea     string `json:"house_area"`     // 建筑面积
	ShareArea     string `json:"share_area"`     // 公摊面积
	IndoorArea    string `json:"indoor_area"`    // 套内面积
	AvgPrice      string `json:"avg_price"`      // 备案单价
	TotalPrince   string `json:"total_prince"`   // 备案总价
	BuildProperty string `json:"build_property"` // 楼盘属性
	DecorateState string `json:"decorate_state"` // 装修属性
	Remark        string `json:"remark"`         // 备注
}

// RecordEstateDetailDTO 楼盘详情信息
type RecordEstateDetailDTO struct {
	RecordNo           string `json:"record_no"`           // 备案号
	EstateName         string `json:"estate_name"`         // 楼盘名称
	EstatePlace        string `json:"estate_place"`        // 坐落位置
	EstateArea         string `json:"estate_area"`         // 所在区域
	EnterpriseName     string `json:"enterprise_name"`     // 开发企业
	PropertyCategories string `json:"property_categories"` // 物业类别
	PropertyCompany    string `json:"property_company"`    // 物业公司
	ProjectInfo        string `json:"project_info"`        // 项目信息
	TrafficInfo        string `json:"traffic_info"`        // 交通信息
	BuildType          string `json:"build_type"`          // 建筑类型
	DesignCompany      string `json:"design_company"`      // 设计单位
	AroundSupport      string `json:"around_support"`      // 周边配套
}

// GetHeFeiFangJiaRecordViewState 获取访问状态，警告：此方法仅供学习参考，禁止用于商业
func GetHeFeiFangJiaRecordViewState() (string, error) {

	// 获取房价信息
	response, err := g.Client().Timeout(20 * time.Second).
		Header(map[string]string{
			"Host":       "drc.hefei.gov.cn",
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		}).
		ContentType("application/x-www-form-urlencoded").
		Post("http://220.178.124.94:8010/fangjia/ws/DefaultList.aspx")
	if err != nil {
		return "", err
	}

	// 解析响应体
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	// 解析数据
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	// 获取客户端状态
	viewState, _ := dom.Find("#__VIEWSTATE").First().Attr("value")
	return viewState, err
}

// ListHeFeiFangJiaRecordPage 获取商品住宅明码标价分页数据
// @viewState 客户端状态
// @pageNum 页码
func ListHeFeiFangJiaRecordPage(viewState string, pageNum int) ([]*RecordHousePriceDTO, error) {

	// 1-1、定义分页参数
	formData := map[string]string{
		"__VIEWSTATE":     viewState,
		"__EVENTTARGET":   "AspNetPager1",
		"__EVENTARGUMENT": strconv.Itoa(pageNum),
	}

	// 1-2、获取房价信息
	response, err := g.Client().Timeout(20*time.Second).
		Header(map[string]string{
			"Host":       "drc.hefei.gov.cn",
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		}).
		ContentType("application/x-www-form-urlencoded").
		Post("http://220.178.124.94:8010/fangjia/ws/DefaultList.aspx", formData)
	if err != nil {
		return nil, err
	}

	// 1-3、解析响应体
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// 2-1.解析数据
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	// 定义数据切片
	var list []*RecordHousePriceDTO
	dom.Find("tr").Each(func(t1 int, s1 *goquery.Selection) {
		if t1 >= 2 {
			recordHousePriceDTO := RecordHousePriceDTO{}
			s1.Find("td").Each(func(t2 int, s2 *goquery.Selection) {
				switch t2 {
				case 0:
					// 获取跳转链接
					href, find := s2.Find("a").Attr("href")
					if find {
						hrefId := gstr.Replace(strings.TrimSpace(href), "Detail2.aspx?Id=", "")
						// 跳转详情ID
						recordHousePriceDTO.HrefId = hrefId
					}

					// 备案号
					recordHousePriceDTO.RecordNo = strings.TrimSpace(s2.Text())
				case 1:
					// 楼盘名称
					recordHousePriceDTO.EstateName = strings.TrimSpace(s2.Text())
				case 2:
					// 楼号
					recordHousePriceDTO.HouseNumbers = strings.TrimSpace(s2.Text())
				case 3:
					// 建筑面积(㎡)
					recordHousePriceDTO.HouseArea = strings.TrimSpace(
						gstr.Replace(gstr.TrimAll(s2.Text()), ",", ""),
					)
				case 4:
					// 套数
					recordHousePriceDTO.HouseNum = strings.TrimSpace(
						gstr.Replace(gstr.TrimAll(s2.Text()), ",", ""),
					)
				case 5:
					// 均价(元/㎡)，去除特殊空格及逗号
					recordHousePriceDTO.AvgPrice = strings.TrimSpace(
						gstr.Replace(gstr.TrimAll(s2.Text()), ",", ""),
					)
				}
			})

			if !gstr.Contains(recordHousePriceDTO.RecordNo, "首页") {
				// 添加切片
				list = append(list, &recordHousePriceDTO)
			}

		}

	})

	return list, err
}

// GetHeFeiFangJiaDetail 获取楼盘详情分页数据
// @viewState 客户端状态
// @hrefId 跳转详情ID
// @pageNum 页码
func GetHeFeiFangJiaDetail(hrefId string) (*RecordEstateDetailDTO, string, error) {

	// 1-2、获取房价详情
	response, err := g.Client().Timeout(20 * time.Second).
		Header(map[string]string{
			"Host":       "drc.hefei.gov.cn",
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		}).
		ContentType("application/x-www-form-urlencoded").
		Post("http://220.178.124.94:8010/fangjia/ws/Detail2.aspx?Id=" + hrefId)
	if err != nil {
		return nil, "", err
	}

	// 1-3、解析响应体
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	// 2-1.解析数据
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, "", err
	}

	// 获取客户端状态
	viewState, find := dom.Find("#__VIEWSTATE").First().Attr("value")

	if !find {
		return nil, "", errors.New("跳转详情ID无效")
	}

	// 定义详情数据
	recordEstateDetailDTO := RecordEstateDetailDTO{}

	title := dom.Find("#txtTitle").First().Text()
	// 楼盘名称
	recordEstateDetailDTO.EstateName = strings.TrimSpace(title)

	detailMap := map[string]string{}
	dom.Find("#IsTableShow").Find("tr").Each(func(t1 int, s1 *goquery.Selection) {
		var key, value string
		s1.Find("td").Each(func(t2 int, s2 *goquery.Selection) {
			switch t2 % 2 {
			case 0:
				key = gstr.Replace(strings.TrimSpace(s2.Text()), "\n", "")
			case 1:
				value = gstr.Replace(strings.TrimSpace(s2.Text()), "\n", "")
				emptyCount := gstr.Count(value, " ")
				// 去除多余空格，8可调，建议5-8
				if emptyCount > 8 {
					value = gstr.Replace(value, " ", "")
				}
				// 追加数据
				detailMap[key] = value
			}
		})

	})

	// 处理详情数据
	for key, value := range detailMap {
		switch key {
		case "交通状况：":
			recordEstateDetailDTO.TrafficInfo = value
		case "周边配套：":
			recordEstateDetailDTO.AroundSupport = value
		case "坐落位置：":
			recordEstateDetailDTO.EstatePlace = value
		case "建筑类型：":
			recordEstateDetailDTO.BuildType = value
		case "开发企业：":
			recordEstateDetailDTO.EnterpriseName = value
		case "所在区域：":
			recordEstateDetailDTO.EstateArea = value
		case "物业公司：":
			recordEstateDetailDTO.PropertyCompany = value
		case "物业类别：":
			recordEstateDetailDTO.PropertyCategories = value
		case "设计单位：":
			recordEstateDetailDTO.DesignCompany = value
		case "项目信息：":
			recordEstateDetailDTO.ProjectInfo = value
		}
	}

	return &recordEstateDetailDTO, viewState, err
}
