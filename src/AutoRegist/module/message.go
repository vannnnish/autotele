/**
 * Created by vannnnish on 2018/6/12.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package module

import (
	"github.com/cjongseok/mtproto"
	"github.com/cjongseok/slog"
	"fmt"
	"context"
	"math/rand"
)

// 自动向其他群组发送消息
func SendMessageBack(mconn *mtproto.Conn, subscriber mtproto.UpdateCallback) {
	mconn.AddUpdateCallback(subscriber)
}

// 用户向群组发送消息
func SendMessageToGroup(user *User, message string, chatNames ...string) {
	if len(chatNames) == 0 {
		for _, chat := range user.Chats {
			peer := &mtproto.TypeInputPeer{&mtproto.TypeInputPeer_InputPeerChannel{
				InputPeerChannel: &mtproto.PredInputPeerChannel{
					ChannelId:  int32(chat.Id),
					AccessHash: int64(chat.AccessHash),
				}}}
			resp, _ := user.scriber.caller.MessagesSendMessage(context.Background(), &mtproto.ReqMessagesSendMessage{
				Peer:     peer,
				Message:  message,
				RandomId: rand.Int63(),
			})
			resp, err := user.scriber.caller.MessagesSendMedia(context.Background(), &mtproto.ReqMessagesSendMedia{
				Peer: peer,
				Media: &mtproto.TypeInputMedia{
					Value: &mtproto.TypeInputMedia_InputMediaPhoto{
						InputMediaPhoto: &mtproto.PredInputMediaPhoto{
							Caption: "什么鬼",
							Id: &mtproto.TypeInputPhoto{
								Value: &mtproto.TypeInputPhoto_InputPhoto{
									InputPhoto: &mtproto.PredInputPhoto{
										Id:         int64(chat.Id),
										AccessHash: int64(chat.AccessHash),
									},
								},
							},
						},
					},
				},
				RandomId: rand.Int63(),
			})
			if err != nil {
				fmt.Println("看看这个错误:", err)
			}
			fmt.Println("send response:", slog.StringifyIndent(resp, "  "))
		}

	}
}
