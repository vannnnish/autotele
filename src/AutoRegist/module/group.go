/**
 * Created by vannnnish on 2018/6/12.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package module

import (
	"github.com/cjongseok/mtproto"
	"context"
	"fmt"
)

type ChatInfo struct {
	Id         int32
	AccessHash int64
	Name       string
	Title      string
}

// 获取群组消息
func GetGroupInfo(caller mtproto.RPCaller) ([]ChatInfo, error) {
	var chatInfos []ChatInfo

	typeMessageChats, err := caller.MessagesGetAllChats(context.Background(), &mtproto.ReqMessagesGetAllChats{
		ExceptIds: []int32{},
	})
	if err != nil {
		return chatInfos, err
	}
	fmt.Println("这个TypeMessageChats", typeMessageChats.GetMessagesChats())
	for _, chat := range typeMessageChats.GetMessagesChats().Chats {
		var chatInfo ChatInfo
		chatInfo.Id = chat.GetChannel().Id
		chatInfo.AccessHash = chat.GetChannel().AccessHash
		chatInfo.Title = chat.GetChannel().Title
		chatInfo.Name = chat.GetChannel().Username
		chatInfos = append(chatInfos, chatInfo)
	}
	return chatInfos, nil
}

// 获取当前用户信息
func GetUserSelf(caller mtproto.RPCaller) (*mtproto.PredUser, error) {
	var predUser *mtproto.PredUser
	typeUserFull, err := caller.UsersGetFullUser(context.Background(), &mtproto.ReqUsersGetFullUser{
		Id: &mtproto.TypeInputUser{
			Value: &mtproto.TypeInputUser_InputUserSelf{},
		},
	})
	if err != nil {
		fmt.Println("err,", err)
		return predUser, err
	}
	predUser = typeUserFull.GetValue().User.GetUser()
	return predUser, nil
}
