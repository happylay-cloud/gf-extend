package hgrpc

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/frame/g"
	"github.com/golang/protobuf/proto"
	"github.com/happylay-cloud/gf-extend/common/hgrpc/test/pb"
)

func TestHgrpc(t *testing.T) {

	person := &pb.Person{
		Name:   "happylay",
		Age:    16,
		Emails: []string{"xxx@qq.com", "xxx@qq.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "181xxxxxxxx",
				Type:   pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "182xxxxxxxx",
				Type:   pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "183xxxxxxxx",
				Type:   pb.PhoneType_WORK,
			},
		},
	}

	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("序列化异常：", err)
	}

	newData := &pb.Person{}
	err = proto.Unmarshal(data, newData)
	if err != nil {
		fmt.Println("反序列化异常：", err)
	}

	g.Dump(newData)

}
