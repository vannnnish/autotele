/**
 * Created by vannnnish on 2018/6/4.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package module

import (
	"github.com/cjongseok/mtproto"
	"github.com/cjongseok/slog"
	"strings"
	"fmt"
)

type User struct {
	SessionName string
	UserName    string
	IsLogin     bool
	UserId      int
	scriber     *subscriber
	Message     chan string
	Chats       []ChatInfo
}

func NewUser(sessionName, userName string, scriber *subscriber, chatInfos []ChatInfo) *User {
	return &User{
		SessionName: sessionName,
		UserName:    userName,
		IsLogin:     true,
		scriber:     scriber,
		Message:     make(chan string),
		Chats:       chatInfos,
	}
}

type subscriber struct {
	caller mtproto.RPCaller
}

func newSubscriber(caller mtproto.RPCaller) *subscriber {
	s := new(subscriber)
	s.caller = caller
	return s
}

func (s *subscriber) OnUpdate(u mtproto.Update) {
	messageRow := slog.StringifyIndent(u, "  ")
	messageArr := strings.Split(messageRow, "mtproto.PredUpdates:")
	fmt.Println("截取:", messageArr, len(messageArr))
	/*	var updates Updates
		 if len(messageArr) == 2 {
			err := json.Unmarshal([]byte(strings.TrimSpace(messageArr[1])), &updates)
			if err != nil {
				fmt.Println("出错:", err)
			} else {
				fmt.Println("结果", updates)
				// caller := mtproto.RPCaller{s.mconn}
				for _, update := range updates.Updates {
					if s.chatInfo.Id != 0 && update.Value.UpdateNewChannelMessage.Message.Value.Message.ToId.Value.PeerChannel.ChannelId == 1147502513 && update.Value.UpdateNewChannelMessage.Message.Value.Message.Message != "" {
						// fmt.Println("消息", update.Value.UpdateNewChannelMessage.Message.Value.Message.Message)
						if !messageSyncMap.isDeleteing {
							// 如果已经存在那么就直接返回
							_, ok := messageSyncMap.syncMap.Load(update.Value.UpdateNewChannelMessage.Message.Value.Message.Id)
							if ok {
								return
							} else {
								messageSyncMap.syncMap.Store(update.Value.UpdateNewChannelMessage.Message.Value.Message.Id, true)
								messageChannel <- update.Value.UpdateNewChannelMessage.Message.Value.Message.Message
							}
						}
						(&sync.Once{}).Do(func() {
							go sendMessage()
						})
						//	messageSyncMap.Store(update.Value.UpdateNewChannelMessage.Message.Value.Message.Message, true)

					}
				}
			}
		}*/
}
