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
	sender := event.Event.Sender
	openId := sender.SenderId.OpenId
	chatId := event.Event.Message.ChatId
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

	if question == "/clear" || question == "清除记忆" {
		p.userCache.Clear(*openId)
		sendMsg(ctx, "🤖️：AI机器人已清除记忆", chatId)
		return nil
	}

	prompt := p.userCache.Get(*openId)
	prompt = fmt.Sprintf("%s\nQ:%s\nA:", prompt, question)
	completions, err := services.Completions(prompt)
	if err != nil {
		sendMsg(ctx, fmt.Sprintf("🤖️：AI机器人摆烂了，请稍后再试～\n错误信息: %v", err), chatId)
		return nil
	}
	p.userCache.Set(*openId, question, completions)
	err = sendMsg(ctx, completions, chatId)
	if err != nil {
		sendMsg(ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), chatId)
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
