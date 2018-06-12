/**
 * Created by vannnnish on 2018/6/4.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package receiver

import (
	"github.com/cjongseok/mtproto"
	"github.com/cjongseok/slog"
	"context"
	"fmt"
	"strings"
	"encoding/json"
	"math/rand"
)

func ReveiveMessage(caller mtproto.RPCaller) {
	emptyPeer := &mtproto.TypeInputPeer{&mtproto.TypeInputPeer_InputPeerEmpty{&mtproto.PredInputPeerEmpty{}}}
	resp, _ := caller.MessagesGetDialogs(context.Background(), &mtproto.ReqMessagesGetDialogs{
		OffsetDate: 0, OffsetId: 0, OffsetPeer: emptyPeer, Limit: 7,
	})
	switch dialogs := resp.GetValue().(type) {
	case *mtproto.TypeMessagesDialogs_MessagesDialogs:
		fmt.Println("什么消息")
		fmt.Println(slog.StringifyIndent(dialogs.MessagesDialogs, "  "))
		messageRow := slog.StringifyIndent(dialogs.MessagesDialogs, "  ")
		messageArr := strings.Split(messageRow, "PredMessagesDialogs:")
		var returnData ReturnMessage
		if len(messageArr) == 2 {
			err := json.Unmarshal([]byte(strings.TrimSpace(messageArr[1])), &returnData)
			if err != nil {
				fmt.Println("err:", err)
			} else {
				fmt.Println(returnData)
				for _, message := range returnData.Messages {
					if message.Value.Message.ToId.Value.PeerChannel.ChannelId == 0 {
					} else {
						fmt.Println("消息:", message.Value.Message.Message)
					}
				}
			}
		}
	case *mtproto.TypeMessagesDialogs_MessagesDialogsSlice:
		fmt.Println("消息切片")
		messageRow := slog.StringifyIndent(dialogs.MessagesDialogsSlice, "  ")
		fmt.Println(slog.StringifyIndent(dialogs.MessagesDialogsSlice, "  "))
		messageArr := strings.Split(messageRow, "PredMessagesDialogsSlice:")
		var returnData ReturnMessage
		if len(messageArr) == 2 {
			err := json.Unmarshal([]byte(strings.TrimSpace(messageArr[1])), &returnData)
			if err != nil {
				fmt.Println("err:", err)
			} else {
				fmt.Println(returnData)
			}
		}

	}
}

type ReturnMessage struct {
	Messages []struct {
		Value struct {
			Message struct {
				Message string
				FromId  int
				ToId struct {
					Value struct {
						PeerChannel struct {
							ChannelId int
						}
					}
				}
			}
		}
	}
}

func SendMessge(chanId int32, chanHash int64, message string, caller mtproto.RPCaller) {

	peer := &mtproto.TypeInputPeer{&mtproto.TypeInputPeer_InputPeerChannel{
		&mtproto.PredInputPeerChannel{
			int32(chanId), int64(chanHash),
		}}}
	resp, _ := caller.MessagesSendMessage(context.Background(), &mtproto.ReqMessagesSendMessage{
		Peer:     peer,
		Message:  message,
		RandomId: rand.Int63(),
	})
	fmt.Println("send response:", slog.StringifyIndent(resp, "  "))
}
