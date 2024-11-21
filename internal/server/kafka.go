package server

import (
	"context"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/mq/kafka"
	"github.com/go-kratos/kratos-layout/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	kftransport "github.com/tx7do/kratos-transport/transport/kafka"
)

func NewKafkaConsumerServer(msgConf *conf.Kafka, greeterService *service.GreeterService) (*kftransport.Server, func()) {
	// Notice 没有用到pem～～～
	//certBytes, err := ioutil.ReadFile(c.Kafka.Pem)
	//if err != nil {
	//	panic(fmt.Sprintf("kafka client read cert file failed %v", err))
	//}
	//clientCertPool := x509.NewCertPool()
	//ok := clientCertPool.AppendCertsFromPEM(certBytes)
	//if !ok {
	//	panic("kafka client failed to parse root certificate")
	//}

	ctx := context.Background()
	kafkaSrv := kftransport.NewServer(
		kftransport.WithAddress(msgConf.Server.Addr), // addr是个列表
		kftransport.WithCodec("json"),
		// Notice 没有用到TLS～～～
		//kafka.WithTLSConfig(&tls.Config{
		//	RootCAs:            clientCertPool,
		//	InsecureSkipVerify: true,
		//}),
		//kafka.WithPlainMechanism(c.Kafka.Username, c.Kafka.Password),
		//kafka.WithBrokerOptions(hello_broker.OptionContextWithValue("max.poll.records", 5000)),
	)

	// 注册一个订阅者，从kafka中消费数据
	err := kafkaSrv.RegisterSubscriber(
		ctx,
		msgConf.GetTopic(), // topic
		msgConf.GetGroup(), // group
		// Notice 手动ACK消息，将这个参数设置为true
		false,
		kafka.RegisterGreeterHandler(greeterService.HandleReceiveGreeterData),
		kafka.GreeterDataCreator,
	)
	if err != nil {
		panic(err)
	}

	// start
	if err = kafkaSrv.Start(ctx); err != nil {
		panic(err)
	}

	return kafkaSrv, func() {
		if err = kafkaSrv.Stop(ctx); err != nil {
			log.Errorf("kafka expected nil got %v", err)
		}
	}
}
