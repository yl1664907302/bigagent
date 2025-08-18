package utils

import (
	"bigagent/internal/config/global"
	"bytes"
	"encoding/json"
	"fmt"
	"k8s.io/klog/v2"
	"net/http"
)

const DingAlertTitle = "logoperator钉钉告警"

// 发送钉钉消息的结构体
type dingMsgNew struct {
	Msgtype string  `json:"msgtype"`
	Text    content `json:"text"`
	At      at      `json:"at"`
}

type content struct {
	Content string `json:"content"`
}
type at struct {
	AtMobiles []string `json:"atMobiles"`
}

type T struct {
	At struct {
		AtMobiles []string `json:"atMobiles"`
		AtUserIds []string `json:"atUserIds"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	Msgtype string `json:"msgtype"`
}

func DingDingMsgDirectSend(imc *global.ImDingDingConf, msg string) {

	dm := dingMsgNew{Msgtype: "text"}
	dm.Text.Content = fmt.Sprintf("[%s]\n[%s]",
		imc.Title,
		msg,
	)
	dm.At.AtMobiles = imc.AtMobiles
	bs, err := json.Marshal(dm)
	if err != nil {
		klog.Errorf(
			"[msg:dingding.msg.json.Marshal.error][dm:%v][err:%v]", dm, err)
		return
	}
	res, err := http.Post(imc.BotApiAddr, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		klog.Errorf("[send.dingding.error][err:%v][msg:v]", err, msg)
		return
	}
	if res != nil && res.StatusCode != 200 {
		klog.Errorf("[send.dingding.status_code.error][code:%v][msg:v]", res.StatusCode, dm.Text.Content)
		return
	}
	klog.Infof("[send.dingding.success][code:%v][msg:%v]", res.StatusCode, dm.Text.Content)

}

func ImTestSent() {
	imc := &global.ImDingDingConf{
		BotApiAddr: "https://oapi.dingtalk.com/robot/send?access_token=75f08bf6f2fa40d45bc987608fa3ffa860bc9d8e2cd2b6099a5cc644ba0b3c50",
		Title:      "通知ntp模块处理结果",
		AtMobiles:  []string{""},
	}
	msg := "node01已cordon"
	DingDingMsgDirectSend(imc, msg)

}
