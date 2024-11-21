package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos-layout/internal/mq/kafka"
	"github.com/tx7do/kratos-transport/broker"

	v1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	"github.com/go-kratos/kratos-layout/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// 从kafka中接收数据并处理
func (s *GreeterService) HandleReceiveGreeterData(ctx context.Context, topic string, headers broker.Headers, msg *kafka.GreeterData) error {

	fmt.Printf("Topic %s,  Msg:%v", topic, msg)

	// Notice 这里面可以做一些业务逻辑～
	return s.uc.AutoHandleGreeterData(ctx, msg)
}
