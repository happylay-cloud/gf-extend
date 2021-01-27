package hnet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/happylay-cloud/gf-extend/net/htcp/hiface"
)

// Package 数据包
type Package struct {
	//--------------------------消息头Head------------------------------
	// ----版本数据长度2B + 密钥长度2B + 加密数据长度2B + 签名数据长度2B ----
	VersionLength int16 // 版本数据长度，2B
	KeyLength     int16 // 加密密钥长度，2B
	DataLength    int16 // 数据部分长度，2B
	SignLength    int16 // 签名数据长度，2B
	//--------------------------消息体Body---------------------=--------
	Version  []byte // Tcp协议版本
	dataType byte   // 内容类型，0x00：none，0x01：json，0x02：bin，1B
	key      []byte // 加密密钥，Ras加密Aes客户端密钥
	Data     []byte // 加密数据
	// --------------------------签名---------------------------------
	Timestamp int32  // 时间戳，4B
	Signature []byte // 加密数据签名，数据 + 时间进行签名。
	//-----------------------------------------------------------------
}

func (p *Package) GetPkg() []byte {
	pkg, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	return pkg
}

func (p *Package) GetPkgBody() []byte {
	return p.Data
}

// NewAcceptPackage 构建数据包拆包对象
func NewAcceptPackage() *Package {
	return &Package{}
}

// NewPackage 构建数据包封包对象
//  @dataType 0x00：none，0x01：json，0x02：bin
//  @data 数据内容
func NewPackage(dataType byte, data []byte) *Package {

	version := []byte("默认版本号")
	secretData := data
	signature := []byte("默认签名")
	key := []byte("默认加密密钥")

	p := &Package{
		VersionLength: int16(len(version)),
		KeyLength:     int16(len(key)),
		DataLength:    int16(len(secretData)),
		SignLength:    int16(len(signature)),
		Version:       version,
		dataType:      dataType,
		key:           key,
		Data:          secretData,
		Timestamp:     int32(time.Now().Unix()),
		Signature:     signature,
	}

	return p
}

// GetHandlerRouter 获取路由（即消息处理器）
func (p *Package) GetHandlerRouter() string {
	// TODO 保留
	return ""
}

// GetPkgHeadLen 获取数据包消息头长度
func (p *Package) GetPkgHeadLen() int16 {
	return 2 + 2 + 2 + 2
}

// Pack 封包（压缩数据）
func (p *Package) Pack() ([]byte, error) {

	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	// 写入版本数据长度
	if err := binary.Write(dataBuff, binary.BigEndian, &p.VersionLength); err != nil {
		return nil, err
	}
	// 写入加密密钥长度
	if err := binary.Write(dataBuff, binary.BigEndian, &p.KeyLength); err != nil {
		return nil, err
	}
	// 写入加密数据长度
	if err := binary.Write(dataBuff, binary.BigEndian, &p.DataLength); err != nil {
		return nil, err
	}
	// 写入签名长度
	if err := binary.Write(dataBuff, binary.BigEndian, &p.SignLength); err != nil {
		return nil, err
	}
	// 写入版本协议
	if err := binary.Write(dataBuff, binary.BigEndian, &p.Version); err != nil {
		return nil, err
	}
	// 写入内容类型数据
	if err := binary.Write(dataBuff, binary.BigEndian, &p.dataType); err != nil {
		return nil, err
	}
	// 写入加密密钥数据
	if err := binary.Write(dataBuff, binary.BigEndian, &p.key); err != nil {
		return nil, err
	}
	// 写入加密数据
	if err := binary.Write(dataBuff, binary.BigEndian, &p.Data); err != nil {
		return nil, err
	}
	// 写入时间戳
	if err := binary.Write(dataBuff, binary.BigEndian, &p.Timestamp); err != nil {
		return nil, err
	}
	// 写入签名
	if err := binary.Write(dataBuff, binary.BigEndian, &p.Signature); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// UnpackHead 拆包head（解压数据）
func (p *Package) UnpackHead(binaryData []byte) (*Package, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 读取写入版本数据长度
	if err := binary.Read(dataBuff, binary.BigEndian, &p.VersionLength); err != nil {
		return nil, err
	}

	// 读取密钥长度
	if err := binary.Read(dataBuff, binary.BigEndian, &p.KeyLength); err != nil {
		return nil, err
	}

	// 读取数据部分长度
	if err := binary.Read(dataBuff, binary.BigEndian, &p.DataLength); err != nil {
		return nil, err
	}
	// 读取签名部分长度
	if err := binary.Read(dataBuff, binary.BigEndian, &p.SignLength); err != nil {
		return nil, err
	}

	return p, nil
}

// Unpack 拆包body（解压数据）
func (p *Package) UnpackBody(binaryData []byte) (*Package, error) {

	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	// 读取协议版本
	p.Version = make([]byte, p.VersionLength)
	if err := binary.Read(dataBuff, binary.BigEndian, &p.Version); err != nil {
		return nil, err
	}
	// 读取内容类型
	if err := binary.Read(dataBuff, binary.BigEndian, &p.dataType); err != nil {
		return nil, err
	}
	// 读取加密密钥
	p.key = make([]byte, p.KeyLength)
	if err := binary.Read(dataBuff, binary.BigEndian, &p.key); err != nil {
		return nil, err
	}
	// 读取加密数据
	p.Data = make([]byte, p.DataLength)
	if err := binary.Read(dataBuff, binary.BigEndian, &p.Data); err != nil {
		return nil, err
	}
	// 读取时间戳
	if err := binary.Read(dataBuff, binary.BigEndian, &p.Timestamp); err != nil {
		return nil, err
	}
	// 读取签名数据
	p.Signature = make([]byte, p.SignLength)
	if err := binary.Read(dataBuff, binary.BigEndian, &p.Signature); err != nil {
		return nil, err
	}
	return p, nil
}

// Unpack 拆包
func (p *Package) Unpack(conn io.Reader) (hiface.ITcpPkg, error) {

	// 先读出流中的head部分
	headDataBin := make([]byte, p.GetPkgHeadLen())
	// ReadFull 会把msg填充满为止
	_, err := io.ReadFull(conn, headDataBin)
	if err != nil {
		g.Log().Line(false).Error("读取head数据异常", err)
		return p, err
	}
	// 将headDataBin字节流拆包到headData中
	headData, err := p.UnpackHead(headDataBin)
	if err != nil {
		g.Log().Line(false).Error("数据包head拆包异常", err)
		return p, err
	}
	// 数据长度
	bodyLength := headData.VersionLength + headData.KeyLength + headData.DataLength + headData.SignLength + 4 + 1

	// body中还有数据
	if bodyLength >= 5 {
		bodyDataBin := make([]byte, bodyLength)
		// ReadFull 会把msg填充满为止
		_, err := io.ReadFull(conn, bodyDataBin)
		if err != nil {
			g.Log().Line(false).Error("读取body数据异常", err)
			return p, err
		}
		// 将bodyDataBin字节流拆包到bodyData中
		_, err = p.UnpackBody(bodyDataBin)
		if err != nil {
			g.Log().Line(false).Error("数据包body拆包异常", err)
			return p, err
		}
	}

	return p, nil
}
