/**
 * Created by vannnnish on 2018/6/4.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package main

import (
	"fmt"
	"os"
	"github.com/cjongseok/slog"
	"net"
	"AutoRegist/model"
	"AutoRegist/constant"
)

const (
	defaultNewKeyFile = ".mtproto"
	appId             = 219949
	appHash           = "54a783ba8153690890d8436aae3b27cf"
	appVersion        = "0.0.1"
	deviceModel       = ""
	systemVersion     = ""
	language          = ""
	telegramAddress   = "149.154.167.50:443"
)

func isServerEndpoint(addr string) (err error) {
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", addr)
	if err == nil && tcpAddr.IP.To4() == nil && tcpAddr.IP.To16() == nil {
		err = fmt.Errorf("invalid ip address")
	}
	return
}

func main() {
	// set up logging
	logf, err := os.OpenFile("ss.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer logf.Close()
	slog.SetLogOutput(logf)

	for _, phone := range constant.PhoneArray {
		// phone := "3044608479"
		key := fmt.Sprintf("./session/+1%s.mtproto", phone)
		go model.Login(appId, appHash, appVersion, deviceModel, systemVersion, language, key, phone)
	}
	// 接收group信息

	blockChan := make(chan bool)
	blockChan <- true
}
