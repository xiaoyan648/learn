package main

import (
	"fmt"
	"time"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
)

const ruleFile = `
{
	"ruleChain": {
	  "id": "rule01",
	  "name": "测试规则链",
	  "root": true
	},
	"metadata": {
	  "nodes": [
		{
			"id":"s1",
			"type": "jsFilter",
			"name": "过滤",
			"debugMode": true,
			"configuration": {
			  "jsScript": "return msg.temperature>10;"
			}
		},
		{
		  "id": "s2",
		  "type": "restApiCall",
		  "name": "获取baidu信息",
		  "debugMode": true,
		  "configuration": {
			"restEndpointUrlPattern": "https://www.baidu.com/",
			"requestMethod": "GET",
			"maxParallelRequestsCount": 200
		  }
		},
		{
			"id": "s3",
			"type": "restApiCall",
			"name": "获取baidu信息",
			"debugMode": true,
			"configuration": {
			  "restEndpointUrlPattern": "https://www.baidu.com/",
			  "requestMethod": "GET",
			  "maxParallelRequestsCount": 200
			}
		}
	  ],
	  "connections": [
		  {
			"fromId": "s1",
			"toId": "s2",
			"type": "True"
		  },
		  {
			"fromId": "s1",
			"toId": "s3",
			"type": "True"
		  }
	],
	  "ruleChainConnections": null
	}
  }  
`

func main() {
	// 创建一个规则引擎实例，每个规则引擎实例有且只有一个根规则链
	ruleEngine, err := rulego.New("rule01", []byte(ruleFile))
	if err != nil {
		panic(err)
	}

	// 定义消息元数据
	metaData := types.NewMetadata()
	metaData.PutValue("productType", "test01")
	// 定义消息内容和消息类型
	msg := types.NewMsg(0, "TELEMETRY_MSG", types.JSON, metaData, "{\"temperature\":15}")

	// 把消息交给规则引擎处理
	ruleEngine.OnMsgWithEndFunc(msg, func(msg types.RuleMsg, err error) {
		fmt.Printf("处理完毕 resp:%+v", msg)
	})

	time.Sleep(2 * time.Second)
}
