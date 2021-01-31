package hgrpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/happylay-cloud/gf-extend/common/hgrpc/test/message"
	"google.golang.org/grpc"
)

// TestServer 服务端
func TestServer(t *testing.T) {
	// 建立服务端监听
	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err.Error())
	}
	// 建立GRPC服务
	rpcServer := grpc.NewServer()
	// 服务注册
	message.RegisterOrderServiceServer(rpcServer, new(OrderServiceImpl))
	err = rpcServer.Serve(lis)
	if err != nil {
		panic(err.Error())
	}
}

// TestClient 客户端
func TestClient(t *testing.T) {

	// 建立RPC通道
	client, err := grpc.Dial("localhost:8090", grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	// 封装请求参数
	orderRequest := &message.OrderRequest{OrderId: "201907300001", TimeStamp: time.Now().Unix()}

	// 创建rpc对象
	orderServiceClient := message.NewOrderServiceClient(client)

	// 调用服务
	orderInfo, _ := orderServiceClient.GetOrderInfo(ctx, orderRequest)
	if orderInfo != nil {
		fmt.Println(orderInfo.GetOrderId())
		fmt.Println(orderInfo.GetOrderName())
		fmt.Println(orderInfo.GetOrderStatus())
	}

	// 调用服务
	orderInfosClient, _ := orderServiceClient.GetOrderInfos(context.TODO(), orderRequest)
	for {
		orderInfo, err := orderInfosClient.Recv()
		if err == io.EOF {
			fmt.Println("读取结束")
			return
		}
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("读取到的信息：", orderInfo)
	}
}

type OrderServiceImpl struct{}

// 具体的方法实现
func (os *OrderServiceImpl) GetOrderInfo(ctx context.Context, request *message.OrderRequest) (*message.OrderInfo, error) {
	orderMap := map[string]message.OrderInfo{
		"201907300001": {OrderId: "201907300001", OrderName: "衣服", OrderStatus: "已付款"},
		"201907310001": {OrderId: "201907310001", OrderName: "零食", OrderStatus: "已付款"},
		"201907310002": {OrderId: "201907310002", OrderName: "食品", OrderStatus: "未付款"},
	}

	var response *message.OrderInfo
	current := time.Now().Unix()
	if request.TimeStamp > current {
		response = &message.OrderInfo{OrderId: "0", OrderName: "", OrderStatus: "订单信息异常"}
	} else {
		result := orderMap[request.OrderId]
		if result.OrderId != "" {
			fmt.Println("请求结果", result)
			return &result, nil
		} else {
			return nil, errors.New("服务器内部错误")
		}
	}
	return response, nil
}

// 获取订单信息
func (os *OrderServiceImpl) GetOrderInfos(request *message.OrderRequest, stream message.OrderService_GetOrderInfosServer) error {
	fmt.Println(" 服务端流RPC模式")

	orderMap := map[string]message.OrderInfo{
		"201907300001": message.OrderInfo{OrderId: "201907300001", OrderName: "衣服", OrderStatus: "已付款"},
		"201907310001": message.OrderInfo{OrderId: "201907310001", OrderName: "零食", OrderStatus: "已付款"},
		"201907310002": message.OrderInfo{OrderId: "201907310002", OrderName: "食品", OrderStatus: "未付款"},
	}
	for id, info := range orderMap {
		if time.Now().Unix() >= request.TimeStamp {
			fmt.Println("订单序列号ID：", id)
			fmt.Println("订单详情：", info)
			// 通过流模式发送给客户端
			_ = stream.Send(&info)
		}
	}
	return nil
}
