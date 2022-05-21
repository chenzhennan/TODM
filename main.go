package main

import (
	"TODM/spider"
	"context"
	"fmt"
	bot "github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"strings"
	"time"
)

const AppID uint64 = 102006691
const BotToken = "MLOp3gg6nzUUOemje5z7i9zyhcFxITRX"

var processor Processor

func main() {
	botToken := token.BotToken(AppID, BotToken)
	api := bot.NewOpenAPI(botToken).WithTimeout(3 * time.Second)
	ctx := context.Background()

	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Printf("%+v,\n err:%v \n", ws, err)

	processor = Processor{
		api:    api,
		spider: spider.NewSpider("", ""),
	}

	intent := websocket.RegisterHandlers(
		ATMessageEventHandler(),
	)

	// 启动 session manager 进行 ws 连接的管理，如果接口返回需要启动多个 shard 的连接，这里也会自动启动多个
	bot.NewSessionManager().Start(ws, botToken, &intent)
}

// ATMessageEventHandler 实现处理 at 消息的回调
func ATMessageEventHandler() event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		fmt.Println("实现处理 at 消息的回调")
		input := strings.ToLower(message.ETLInput(data.Content))
		return processor.ProcessATMessage(input, data)
	}
}
