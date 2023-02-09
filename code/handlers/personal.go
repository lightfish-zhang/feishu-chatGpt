package handlers

import (
	"context"
	"fmt"
	"start-feishubot/services"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type PersonalMessageHandler struct {
	userCache services.UserCacheInterface
	msgCache  services.MsgCacheInterface
}

func (p PersonalMessageHandler) handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {

	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	if p.msgCache.IfProcessed(*msgId) {
		fmt.Println("msgId", *msgId, "processed")
		return nil
	}
	p.msgCache.TagProcessed(*msgId)
	question := parseContent(*content)
	if len(question) == 0 {
		fmt.Println("msgId", *msgId, "message.text is empty")
		return nil
	}
	sender := event.Event.Sender
	openId := sender.SenderId.OpenId
	prompt := p.userCache.Get(*openId)
	prompt = fmt.Sprintf("%s\nQ:%s\nA:", prompt, question)
	completions, err := services.Completions(prompt)
	if err != nil {
		sendMsg(ctx, fmt.Sprintf("🤖️：AI机器人摆烂了，请稍后再试～ \n 错误: %v", err), event.Event.Message.ChatId)
		return nil
	}
	p.userCache.Set(*openId, question, completions)
	err = sendMsg(ctx, completions, event.Event.Message.ChatId)
	if err != nil {
		sendMsg(ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～ \n 错误: %v", err), event.Event.Message.ChatId)
		return nil
	}
	return nil

}

var _ MessageHandlerInterface = (*PersonalMessageHandler)(nil)

func NewPersonalMessageHandler() MessageHandlerInterface {
	return &PersonalMessageHandler{
		userCache: services.GetUserCache(),
		msgCache:  services.GetMsgCache(),
	}
}
