package hrsa

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/frame/g"
)

func TestRsa(t *testing.T) {
	prvkey, pubkey, _ := GenRsaKey(2048)
	g.Dump(prvkey, pubkey)

	block, err := RsaEncryptBlock([]byte("RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密RSA加密"), pubkey)
	if err != nil {
		fmt.Println(err)
	}

	g.Dump(gbase64.EncodeToString(block))

	decryptData, err := RsaDecryptBlock(block, prvkey)
	if err != nil {
		fmt.Println(err)
	}
	g.Dump(decryptData)

	sha256SignData, err := RsaSignWithSha256(decryptData, prvkey)
	if err != nil {
		fmt.Println(err)
	}
	g.Dump(gbase64.EncodeToString(sha256SignData))

	withSha256, err := RsaVerySignWithSha256(decryptData, sha256SignData, pubkey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(withSha256)
}
