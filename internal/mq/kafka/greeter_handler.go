package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tx7do/kratos-transport/broker"
)

// 从kafka中解析GreeterData数据
func GreeterDataCreator() broker.Any { return &GreeterData{} }

// service中处理的函数需要跟它保持一致
type GreeterDataHandler func(ctx context.Context, topic string, headers broker.Headers, msg *GreeterData) error

// handle
func RegisterGreeterHandler(fnc GreeterDataHandler) broker.Handler {
	return func(ctx context.Context, event broker.Event) error {
		var msg *GreeterData = nil

		eventMsgBody := event.Message().Body
		fmt.Println("eventMsgBody:>>>>> ", eventMsgBody)
		switch t := eventMsgBody.(type) {
		case []byte:
			fmt.Println("byte............................")
			msg = &GreeterData{}
			if err := json.Unmarshal(t, msg); err != nil {
				return err
			}
		case string:
			fmt.Println("string............................")
			msg = &GreeterData{}
			if err := json.Unmarshal([]byte(t), msg); err != nil {
				return err
			}
		case *GreeterData:
			fmt.Println("相同类型............................", t.Hello)
			msg = t
		default:
			fmt.Println("default............................")
			return fmt.Errorf("unsupported type: %T", t)
		}

		// Notice 使用service中的handler函数做业务处理，把msg传进去～
		fmt.Println("evnet.Message().Headers:>>>>>>>>>>> ", event.Message().Headers)
		if err := fnc(ctx, event.Topic(), event.Message().Headers, msg); err != nil {
			return err
		}

		fmt.Println("msg:>>>>>>>>>>>> ", msg)

		return nil
	}
}
