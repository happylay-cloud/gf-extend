package hjs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gogf/gf/encoding/gjson"
	"strings"

	"github.com/dop251/goja"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/crypto/gsha1"
	"github.com/gogf/gf/text/gstr"
)

// JsCookieTwoGo 参数
type JsCookieTwoGo struct {
	BtsStr string
	Charts string
	Ct     string
	Ha     string
	Tn     string
	Vt     string
	Wt     string
}

// GetJsCookieOne 获取第一个JsCookie内容，__jsl_clearance
func GetJsCookieOne(script string) (string, error) {
	// 处理原始cookie
	scriptStr := gstr.Replace(script, "<script>document.cookie=", "")
	scriptStr = gstr.Replace(scriptStr, ";location.href=location.pathname+location.search</script>", "")
	// 解析cookie
	result, err := Eval(scriptStr)
	if err != nil {
		return "", err
	}
	return gstr.Split(result.String(), ";")[0], nil
}

// GetJsCookieTwo 获取第二个JsCookie内容，__jsl_clearance
func GetJsCookieTwo(script string) (string, error) {
	// 解析参数
	twoGo := GetJsCookieTwoGo(script)
	// 解析第二次Cookie值
	return GetJsCookie(twoGo.Charts, twoGo.BtsStr, twoGo.Ct, twoGo.Ha, twoGo.Tn)
}

// GetJsCookieTwoGo 获取JsCookie混淆参数
func GetJsCookieTwoGo(script string) *JsCookieTwoGo {
	// 过滤信息
	scriptStr := gstr.Replace(script, "<script>", "")
	scriptStr = gstr.Replace(scriptStr, "})</script>", "")
	itemArrStr := gstr.StrEx(scriptStr, "go({")
	// 提取核心参数
	jsParams := "{" + strings.TrimSpace(itemArrStr) + "}"
	json, _ := gjson.LoadContent(jsParams)
	// 解析参数
	btsArr := json.GetStrings("bts")
	// 加密字符串
	btsStr := btsArr[0] + "," + btsArr[1]
	// 加密字符串
	chars := json.GetString("chars")
	// 校验结果
	ct := json.GetString("ct")
	// 加密方式
	ha := json.GetString("ha")
	// Cookie请求头
	tn := json.GetString("tn")

	vt := json.GetString("vt")
	wt := json.GetString("wt")

	return &JsCookieTwoGo{
		btsStr,
		chars,
		ct,
		ha,
		tn,
		vt,
		wt,
	}
}

// Eval javascript执行函数
func Eval(script string) (goja.Value, error) {
	// 实例化js引擎
	vm := goja.New()
	// 执行脚本
	return vm.RunString(script)
}

// GetJsCookie 获取js-cookie值
func GetJsCookie(chars string, btsStr string, ct string, ha string, tn string) (string, error) {
	bts := gstr.Split(btsStr, ",")
	length := len(chars)
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			// 参与加密值
			value := bts[0] + gstr.SubStr(chars, i, 1) + gstr.SubStr(chars, j, 1) + bts[1]
			// 解析加密算法
			switch ha {
			case "sha1":
				sha1Value := gsha1.Encrypt(value)
				if gstr.Equal(sha1Value, ct) {
					return tn + "=" + value, nil
				}
			case "md5":
				md5Value, _ := gmd5.Encrypt(value)
				if gstr.Equal(md5Value, ct) {
					return tn + "=" + value, nil
				}
			case "sha256":
				sha256.New()
				h := sha256.New()
				h.Write([]byte(value))
				sha256Value := hex.EncodeToString(h.Sum(nil))
				if gstr.Equal(sha256Value, ct) {
					return tn + "=" + value, nil
				}
			default:
				return "", errors.New("参数ha已更新")
			}

		}
	}

	return "", errors.New("无效js参数")
}
