package hjsoup

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"github.com/happylay-cloud/gf-extend/common/hcache"
	"github.com/happylay-cloud/gf-extend/common/hutils/hctx"
	"github.com/happylay-cloud/gf-extend/common/hutils/hstr"

	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DefaultProductCodeCacheKey 缓存键
var defaultProductCodeCacheKey = "OPEN_PRODUCT_CODE:"

// DefaultProductCacheBucket 默认缓存桶
var defaultProductCacheBucket = "DEFAULT_DB0"

// ProductCodeDto 商品条码信息
type ProductCodeDto struct {
	ProductCode      string   `json:"product_code"`       // 商品条码
	ProductCodeImage string   `json:"product_code_image"` // 商品条形码图片
	CompanyName      string   `json:"company_name"`       // 企业名称
	CompanyAddress   string   `json:"company_address"`    // 企业注册地址
	ProductName      string   `json:"product_name"`       // 产品名称
	ProductCategory  string   `json:"product_category"`   // 产品分类
	Brand            string   `json:"brand"`              // 品牌
	ProductSpec      string   `json:"product_spec"`       // 产品规格
	StandardNo       string   `json:"standard_no"`        // 标准号
	StandardName     string   `json:"standard_name"`      // 标准名称
	ProductExp       string   `json:"product_exp"`        // 保质期
	UpMarketTime     string   `json:"up_market_time"`     // 上市日期
	DownMarketTime   string   `json:"down_market_time"`   // 下市日期
	ProductImageList []string `json:"product_image_list"` // 图片列表
}

// SearchByProductCode 根据商品条码查询商品信息，警告：此方法仅供学习参考，禁止用于商业
//	@productCode	商品条码
//	@debug			是否开启debug
func SearchByProductCode(productCode string, debug bool) (*ProductCodeDto, error) {
	// 1.获取验证码信息
	response, err := g.Client().
		Timeout(20 * time.Second).
		Header(map[string]string{
			"Host":       "www.chinatrace.org",
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		}).
		Get("http://www.chinatrace.org/trace/verification/image?_=" + strconv.Itoa(time.Now().Second()*1000))
	if err != nil {
		return nil, err
	}

	cookies := response.Cookies()

	// 会话sessionId
	sessionId := cookies[0].Value

	g.Log().Line(false).Debug("当前会话：", sessionId)

	// 2.查询y值
	y := gjson.New(response.ReadAllString()).GetString("y")

	g.Log().Line(false).Info("验证码y值：", y)

	// 定义有效x值
	validX := ""

	// 定义验证码
	doorCode := ""

	// ******************************** 多任务处理-开始 ********************************

	// 定义任务
	doManyTask := hctx.DoManyTask{
		Count:   350,
		Timeout: 20 * time.Second,
		Debug:   true,
	}

	// 定义返回值
	type testTaskValue struct {
		ValidX   string
		DoorCode string
	}

	// 执行任务
	successOne, err := doManyTask.DoTaskSuccessOne(nil, func(do *hctx.DoManyTask, ctx context.Context, channel chan interface{}, wg *sync.WaitGroup, index int, params interface{}) {

		if debug {
			g.Log().Line(false).Debug("任务执行中...，序号：", index)
		}

		// ************************ 业务处理 ************************

		// 封装数据
		taskValue := testTaskValue{}

		var isNeedReturn bool

		// 校验验证码
		respDoorCode, respX, err := validQrCode(index, y, sessionId, debug)
		if err == nil {
			taskValue.DoorCode = respDoorCode
			taskValue.ValidX = respX
			isNeedReturn = true
		} else {
			isNeedReturn = false
		}

		// ************************ 返回数据 ************************

		// 获取返回结果，必须执行
		do.WaitDataReturn(isNeedReturn, ctx, channel, wg, index, taskValue)

	})

	if err != nil {
		return nil, err
	}

	// 获取返回值
	resp := successOne.(testTaskValue)

	// 解析返回值validX
	validX = resp.ValidX

	// 解析验证码
	doorCode = resp.DoorCode

	x, err := strconv.Atoi(validX)
	if err != nil {
		return nil, err
	}

	// 必须重新验证一次（多次验证只有最后一次生效）
	doorCode, validX, err = validQrCode(x, y, sessionId, true)

	// ******************************** 多任务处理-结束 ********************************

	if gstr.LenRune(doorCode) == 0 {
		return nil, errors.New("查询超时")
	}

	// 此处传参是正确的
	formData1 := map[string]string{
		"productCode":  productCode,
		"batchNo":      "",
		"productCode1": "",
		"traceCode":    "",
		"doorCode":     doorCode,
		"validX":       validX,
	}

	g.Log().Line(false).Info("请求参数：", formData1)

	if debug {
		// 警告：不能按照这种方式传参->使用FormPost方法，此处cookie会丢失
		formData2 := url.Values{
			"productCode":  {productCode},
			"batchNo":      {""},
			"productCode1": {""},
			"traceCode":    {""},
			"doorCode":     {doorCode},
			"validX":       {validX},
		}
		g.Log().Line(false).Info("格式化参数：", formData2.Encode())
	}

	formResp, err := g.Client().Timeout(20*time.Second).
		SetCookieMap(map[string]string{
			"JSESSIONID": sessionId,
		}).
		Header(map[string]string{
			"Host":       "www.chinatrace.org",
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		}).
		ContentType("application/x-www-form-urlencoded").
		Post("http://www.chinatrace.org/trace/door/controller/SearchController/searchByProductCode.do", formData1)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(formResp.Body)
	if err != nil {
		return nil, err
	}

	if debug {
		g.Log().Line(false).Info("请求Cookie：", formResp.Request.Cookies())
		g.Log().Line(false).Info("原始返回值：", string(body))
	}

	// 商品条码信息
	productCodeDto := ProductCodeDto{}

	// 设置商品条码
	productCodeDto.ProductCode = productCode

	// 处理异常信息，返回json而不是html
	if gstr.LenRune(string(body)) < 500 {
		return nil, errors.New("请求频繁，稍后再试")
	}

	// 4.解析数据
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	// ID选择器，处理企业信息
	dom.Find("#company").Each(func(i int, selection *goquery.Selection) {
		selection.Find("p").Each(func(i int, selection *goquery.Selection) {
			switch i {
			case 0:
				// 设置企业名称
				productCodeDto.CompanyName = selection.Find("span").Text()
			case 1:
				// 设置企业注册地址
				productCodeDto.CompanyAddress = selection.Find("span").Text()
			}
		})
	})

	productDiv := dom.Find("#product")

	// class选择器
	table := productDiv.Find(".table")

	table.Find("tr").Each(func(t1 int, selection1 *goquery.Selection) {

		// 处理tr单元格
		selection1.Find("td").Each(func(t2 int, selection2 *goquery.Selection) {
			if t2 == 1 {
				switch t1 {
				case 1:
					// 设置产品名称
					productCodeDto.ProductName = hstr.TrimBlank(selection2.Text())
				case 2:
					// 设置产品分类
					productCodeDto.ProductCategory = hstr.TrimBlank(selection2.Text())
				case 3:
					// 设置品牌
					productCodeDto.Brand = hstr.TrimBlank(selection2.Text())
				case 4:
					// 设置商品规格
					productCodeDto.ProductSpec = hstr.TrimBlank(selection2.Text())
				case 5:
					// 设置标准号
					productCodeDto.StandardNo = hstr.TrimBlank(selection2.Text())
				case 6:
					// 设置标准名称
					productCodeDto.StandardName = hstr.TrimBlank(selection2.Text())
				case 7:
					// 设置保质期
					productCodeDto.ProductExp = hstr.TrimBlank(selection2.Text())
				case 8:
					// 设置上市日期
					productCodeDto.UpMarketTime = hstr.TrimBlank(selection2.Text())
				case 9:
					// 设置下市日期
					productCodeDto.DownMarketTime = hstr.TrimBlank(selection2.Text())
				}
			}

		})
	})

	// 实例化图片切片
	imgList := make([]string, 0)

	// 元素选择器，处理图片列表
	table.Find("img").Each(func(i int, selection *goquery.Selection) {
		src, _ := selection.Attr("src")
		if i == 0 {
			// 设置条形码图片
			productCodeDto.ProductCodeImage = src
		} else {
			imgList = append(imgList, src)
		}

	})
	// 设置图片列表
	productCodeDto.ProductImageList = imgList

	if gstr.LenRune(productCodeDto.ProductName) == 0 {
		return nil, errors.New("未知商品")
	}

	return &productCodeDto, nil
}

// validQrCode 验证二维码
func validQrCode(x int, y string, sessionId string, debug bool) (doorCode string, validX string, err error) {
	body := g.Client().
		SetCookieMap(map[string]string{
			"JSESSIONID": sessionId,
		}).
		Header(map[string]string{
			"Host":       "www.chinatrace.org",
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		}).
		PostContent("http://www.chinatrace.org/trace/verification/result?x=" + strconv.Itoa(x) + "&y=" + y)

	if gstr.LenRune(body) > 0 {
		doorCode = body
		validX = strconv.Itoa(x)
		if debug {
			g.Log().Line(false).Debug("Valid："+strconv.Itoa(x)+"，", "验证码：", doorCode)
		}
		return doorCode, validX, nil
	}

	return "", "", errors.New("验证无效")
}

// SearchByProductCodeCache 根据商品条码查询商品信息，优先基于本地缓存，警告：此方法仅供学习参考，禁止用于商业
// @productCode	商品条码
func SearchByProductCodeCache(productCode string, debug bool) (productCodeDto *ProductCodeDto, err error) {

	// 查询缓存
	entry, err := hcache.GetCache(defaultProductCacheBucket, []byte(defaultProductCodeCacheKey+productCode))
	if err == nil {

		// 处理空数据
		if gstr.LenRune(string(entry.Value)) == 0 {
			return nil, errors.New("未知商品")
		}

		err = json.Unmarshal(entry.Value, &productCodeDto)
		if err != nil {
			return nil, err
		}

		return productCodeDto, nil
	}

	// 查询数据
	productCodeInfo, err := SearchByProductCode(productCode, debug)

	if err != nil {
		// 数据不存在，添加缓存并设置过期时间
		err = hcache.SetCache(defaultProductCacheBucket, []byte(defaultProductCacheBucket+productCode), []byte(""), 60)
		if err != nil {
			g.Log().Line(false).Error("NutsDb保存缓存失败，异常信息：" + err.Error())
		}
		return productCodeInfo, nil
	}

	if err == nil {
		// 保存数据至缓存
		jsonByte, err := json.Marshal(productCodeInfo)
		if err == nil {
			err = hcache.SetCache(defaultProductCacheBucket, []byte(defaultProductCodeCacheKey+productCode), jsonByte, 0)
			if err != nil {
				g.Log().Line(false).Error("NutsDb保存缓存失败，异常信息：" + err.Error())
			}
			return productCodeInfo, nil
		}
	}

	return productCodeInfo, err
}
