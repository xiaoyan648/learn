package feishu

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	// httpx ("github.com/go-leo/gox/netx/httpx")

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// SDK 使用文档：https://github.com/larksuite/oapi-sdk-go/tree/v3_main
// 开发者复制该Demo后，需要修改Demo里面的"appID", "appSecret"为自己应用的appId,appSecret.
func CardSend() {
	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	client := lark.NewClient("cli_a5c9a414df7d500e", "BiuygCC61BNjzQYQ6mfwRddNpkzqVGmZ", lark.WithEnableTokenCache(true))

	// 创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`chat_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(`oc_d6cfe38ce6ffbdd4684ba6358c26e2bc`).
			MsgType(`interactive`).
			Content(`{"config": {"wide_screen_mode": true},"header": {"title": {"tag": "plain_text","content": "测试卡片"}},"elements": [{"tag": "div","fields": [{"is_short": false,"text": {"tag": "lark_md","content": ""}},{"is_short": false,"text": {"tag": "lark_md","content": "**时间：**\n2020-4-8 至 2020-4-10（共3天）"}},{"is_short": false,"text": {"tag": "lark_md","content": ""}},{"is_short": true,"text": {"tag": "lark_md","content": "**备注**\n测试功能"}}]},{"tag": "hr"},{"tag": "action","layout": "bisected","actions": [{"tag": "button","text": {"tag": "plain_text","content": "批准"},"type": "primary","value": {"chosen": "approve"}},{"tag": "button","text": {"tag": "plain_text","content": "拒绝"},"type": "primary","value": {"chosen": "decline"}}]}]}`).
			Build()).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Im.Message.Create(context.Background(), req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
}

func CardSendWebHook() {
	url := "https://open.feishu.cn/open-apis/bot/v2/hook/12a23fbf-dccf-46de-9ec6-a7fce670f229"
	req, err := http.NewRequest("POST", url, strings.NewReader(`{"config": {"wide_screen_mode": true},"header": {"title": {"tag": "plain_text","content": "测试卡片"}},"elements": [{"tag": "div","fields": [{"is_short": false,"text": {"tag": "lark_md","content": ""}},{"is_short": false,"text": {"tag": "lark_md","content": "**时间：**\n2020-4-8 至 2020-4-10（共3天）"}},{"is_short": false,"text": {"tag": "lark_md","content": ""}},{"is_short": true,"text": {"tag": "lark_md","content": "**备注**\n测试功能"}}]},{"tag": "hr"},{"tag": "action","layout": "bisected","actions": [{"tag": "button","text": {"tag": "plain_text","content": "批准"},"type": "primary","value": {"chosen": "approve"}},{"tag": "button","text": {"tag": "plain_text","content": "拒绝"},"type": "primary","value": {"chosen": "decline"}}]}]}`))
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
}
