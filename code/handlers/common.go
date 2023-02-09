package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"start-feishubot/initialization"
	"strings"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func sendMsg(ctx context.Context, msg string, chatId *string) error {
	fmt.Println("sendMsg", msg, chatId)
	msgB, _ := json.Marshal(strings.TrimSpace(msg))
	msg = string(msgB)
	if len(msg) >= 2 {
		msg = msg[1 : len(msg)-1]
	}
	client := initialization.GetLarkClient()
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()
	fmt.Println("content", content)

	resp, err := client.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(*chatId).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		err = fmt.Errorf("[%d][%s]%s", resp.Code, resp.RequestId(), resp.Msg)
		fmt.Println(err)
		return err
	}
	return nil
}
func msgFilter(msg string) string {
	//replace @到下一个非空的字段 为 ''
	regex := regexp.MustCompile(`@[^ ]*`)
	return regex.ReplaceAllString(msg, "")

}
func parseContent(content string) string {
	//"{\"text\":\"@_user_1  hahaha\"}",
	//only get text content hahaha
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
	}
	text := contentMap["text"].(string)
	return msgFilter(text)
}
