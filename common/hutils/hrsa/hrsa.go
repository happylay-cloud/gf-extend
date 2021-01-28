package hrsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/gogf/gf/frame/g"
)

// GenRsaKey 产生RSA密钥对
func GenRsaKey(bits ...int) (prvkey, pubkey []byte, err error) {
	// 默认长度
	tmpBits := 1024
	if len(bits) > 0 {
		tmpBits = bits[0]
	}
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, tmpBits)
	if err != nil {
		return nil, nil, err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey = pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = pem.EncodeToMemory(block)

	return prvkey, pubkey, nil
}

// RsaEncryptBlock 公钥加密-分段
func RsaEncryptBlock(src, publicKeyByte []byte) (bytesEncrypt []byte, err error) {
	// 解密pem格式的公钥
	block, _ := pem.Decode(publicKeyByte)
	if block == nil {
		return nil, errors.New("获取公钥失败")
	}
	// 解析公钥
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := publicKey.(*rsa.PublicKey)

	keySize, srcSize := pub.Size(), len(src)

	g.Log().Line(false).Debug("密钥长度：", keySize, "\t明文长度：\t", srcSize)

	// 单次加密的长度需要减掉padding的长度，PKCS1为11
	offSet, once := 0, keySize-11
	buffer := bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + once
		if endIndex > srcSize {
			endIndex = srcSize
		}
		// 加密一部分
		bytesOnce, err := rsa.EncryptPKCS1v15(rand.Reader, pub, src[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	bytesEncrypt = buffer.Bytes()
	return
}

// RsaDecryptBlock 私钥解密-分段
func RsaDecryptBlock(src, privateKeyByte []byte) (bytesDecrypt []byte, err error) {
	// 获取私钥
	block, _ := pem.Decode(privateKeyByte)
	if block == nil {
		return nil, errors.New("获取私钥失败")
	}
	// 解析PKCS1格式的私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	keySize, srcSize := privateKey.Size(), len(src)
	g.Log().Line(false).Debug("密钥长度：", keySize, "\t密文长度：\t", srcSize)

	var offSet = 0
	var buffer = bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + keySize
		if endIndex > srcSize {
			endIndex = srcSize
		}
		// 解密一部分
		bytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, src[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	bytesDecrypt = buffer.Bytes()
	return bytesDecrypt, nil
}
