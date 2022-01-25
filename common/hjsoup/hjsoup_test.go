package hjsoup

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

func TestHttpClient(t *testing.T) {

	// 1.获取验证码信息
	response, err := g.Client().
		Timeout(20 * time.Second).
		Header(map[string]string{
			"Host":       "www.chinatrace.org",
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		}).
		Get("http://www.chinatrace.org/trace/verification/image?_=" + strconv.Itoa(time.Now().Second()*1000))
	if err != nil {
		panic(err)
	}

	cookies := response.Cookies()

	// 会话sessionId
	sessionId := cookies[0].Value

	fmt.Println(sessionId)

	// 2.查询y值
	y := gjson.New(response.ReadAllString()).GetString("y")
	fmt.Println(y)

	validX := ""

	doorCode := ""
	// 3.解析验证码
	for x := 0; x < 350; x++ {
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
			fmt.Println("执行次数："+strconv.Itoa(x)+"次，", "获取验证码：", doorCode)
			break
		}
	}

	if gstr.LenRune(doorCode) == 0 {
		return
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
		Post("http://www.chinatrace.org/trace/door/controller/SearchController/searchByProductCode.do", map[string]string{
			"productCode":  "6956401264074",
			"batchNo":      "",
			"productCode1": "",
			"traceCode":    "",
			"doorCode":     doorCode,
			"validX":       validX,
		})

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(formResp.Body)
	if err != nil {
		return
	}

	formResp.RawDump()

	fmt.Println(string(body))

}
