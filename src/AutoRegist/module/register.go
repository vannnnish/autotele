/**
 * Created by vannnnish on 2018/6/12.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package module

import (
	"github.com/cjongseok/mtproto"
	"fmt"
	"log"
)

const (
	defaultNewKeyFile = ".mtproto"

	appVersion      = "0.0.1"
	deviceModel     = ""
	systemVersion   = ""
	language        = ""
	telegramAddress = "149.154.167.50:443"
)

func Register(apiId int32, apiHash, phoneNumber, addr, key string) (*mtproto.Conn, error) {
	var mconn *mtproto.Conn
	config, err := mtproto.NewConfiguration(apiId, apiHash, appVersion, deviceModel, systemVersion, language, 0, 0, key)
	if err != nil {
		return mconn, err
	}
	config.KeyPath = "session/" + phoneNumber + defaultNewKeyFile
	// request to send authentication code to the phone
	var sentCode *mtproto.TypeAuthSentCode
	manager, err := mtproto.NewManager(config)
	if err != nil {
		return mconn, err
	}
	mconn, sentCode, err = manager.NewAuthentication(phoneNumber, addr, false)
	if err != nil {
		return mconn, err
	}

	// sign-in with the code from the user input
	var code string
	fmt.Printf("Enter Code: ")
	fmt.Scanf("%s", &code)
	log.Println("entered code = ", code)
	_, err = mconn.SignIn(phoneNumber, code, sentCode.GetValue().PhoneCodeHash)
	return mconn, nil
}
