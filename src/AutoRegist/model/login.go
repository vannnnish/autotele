/**
 * Created by vannnnish on 2018/6/4.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package model

import (
	"github.com/cjongseok/mtproto"
	"log"
	"fmt"
	"github.com/cjongseok/slog"
	"strings"
	"encoding/json"
	"context"
	"AutoRegist/receiver"
	"sync"
	"math/rand"
	"time"
)

var addr = "149.154.167.50:443"
var subscribers []*subscriber
var Users []*User

func Login(apiId int32, apiHash, appVersion, deviceModel, systemVersion, language, key string, phoneNumber string) (mtproto.RPCaller, error) {
	var manager *mtproto.Manager
	var caller mtproto.RPCaller
	var mconn *mtproto.Conn
	config, err := mtproto.NewConfiguration(apiId, apiHash, appVersion, deviceModel, systemVersion, language, 0, 0, key)
	if err != nil {
		fmt.Println("err:", err)
		return caller, err
	}
	// Connect by phone number

	log.Println("MAIN: load authentication")
	manager, err = mtproto.NewManager(config)
	if err != nil {
		fmt.Println("err:", err)
		return caller, err
	}
	mconn, err = manager.LoadAuthentication(phoneNumber)
	if err != nil {
		fmt.Println("err:", err)
		return caller, err
	}

	caller = mtproto.RPCaller{mconn}
	typeMessageChats, err := caller.MessagesGetAllChats(context.Background(), &mtproto.ReqMessagesGetAllChats{
		ExceptIds: []int32{},
	})
	var chatInfo ChatInfo
	for _, chat := range typeMessageChats.GetMessagesChats().Chats {
		if chat.GetChannel().Username == "testgroupofmine" {
			chatInfo.Id = chat.GetChannel().Id
			chatInfo.AccessHash = chat.GetChannel().AccessHash
			chatInfo.Title = chat.GetChannel().Title
			chatInfo.Name = chat.GetChannel().Username
		}
		fmt.Println("id:", chat.GetChannel().Id)
		fmt.Println("accessHash:", chat.GetChannel().AccessHash)
		fmt.Println("name:", chat.GetChannel().Username)
		fmt.Println("Title:", chat.GetChannel().Title)
	}
	subscriber := newSubscriber(mconn, chatInfo)
	subscribers = append(subscribers, subscriber)
	// mconn.AddUpdateCallback(subscriber)
	fmt.Println("看看这个chatInfo:", chatInfo)
	return caller, nil
}

func (s *subscriber) OnUpdate(u mtproto.Update) {
	messageRow := slog.StringifyIndent(u, "  ")
	messageArr := strings.Split(messageRow, "mtproto.PredUpdates:")
	fmt.Println("截取:", messageArr, len(messageArr))
	var updates Updates
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
	}
}

func sendMessage() {
	for {
		randNumber := rand.Intn(3)
		select {
		case message := <-messageChannel:
			fmt.Println("看看这个随机数:", randNumber)
			caller := mtproto.RPCaller{subscribers[randNumber].mconn}
			receiver.SendMessge(subscribers[randNumber].chatInfo.Id, subscribers[randNumber].chatInfo.AccessHash, message, caller)
		case <-time.Tick(45 * time.Second):
			messageSyncMap.DeleteData()
		}
	}
}

var messageChannel = make(chan string, 20)

var messageSyncMap = MessageSyncMap{
	syncMap:     sync.Map{},
	isDeleteing: false,
}

func (m MessageSyncMap) DeleteData() {
	m.isDeleteing = true
	m.syncMap.Range(func(key, value interface{}) bool {
		m.syncMap.Delete(key)
		return true
	})
	m.isDeleteing = false
}

type MessageSyncMap struct {
	syncMap     sync.Map
	isDeleteing bool
}

type MessageData struct {
	Message string
}

type Updates struct {
	Updates []struct {
		Value struct {
			UpdateNewChannelMessage struct {
				Message struct {
					Value struct {
						Message struct {
							FromId int
							Id     int
							ToId struct {
								Value struct {
									PeerChannel struct {
										ChannelId int
									}
								}
							}
							Message string
						}
					}
				}
			}
		}
	}
	Chats []struct {
		Value struct {
			Channel struct {
				Id         int32
				AccessHash int64
			}
		}
	}
}

type ChatInfo struct {
	Id         int32
	AccessHash int64
	Name       string
	Title      string
}
