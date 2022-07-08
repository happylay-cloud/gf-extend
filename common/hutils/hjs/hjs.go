package hjs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/dop251/goja"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/crypto/gsha1"
	"github.com/gogf/gf/text/gstr"
)

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
