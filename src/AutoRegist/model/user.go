/**
 * Created by vannnnish on 2018/6/4.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package model

import (
	"github.com/cjongseok/mtproto"
)

type User struct {
	SessionName string
	UserName    string
	IsLogin     bool
	UserId      int
	scriber     *subscriber
	Message     chan string
}

func NewUser(sessionName, userName string, scriber *subscriber) *User {
	return &User{
		SessionName: sessionName,
		UserName:    userName,
		IsLogin:     true,
		scriber:     scriber,
		Message:     make(chan string),
	}
}

type subscriber struct {
	chatInfo ChatInfo
	mconn    *mtproto.Conn
}

func newSubscriber(mconn *mtproto.Conn, chatInfo ChatInfo) *subscriber {
	s := new(subscriber)
	s.chatInfo = chatInfo
	s.mconn = mconn
	return s
}
