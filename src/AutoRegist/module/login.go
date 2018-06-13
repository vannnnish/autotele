/**
 * Created by vannnnish on 2018/6/4.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package module

import (
	"github.com/cjongseok/mtproto"
	"fmt"
	"sync"
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
	// 获取组信息
	chatInfos, err := GetGroupInfo(caller)
	if err != nil {
		return caller, err
	}
	fmt.Println("组信息:", chatInfos)
	// 获取用户信息
	userInfo, err := GetUserSelf(caller)
	if err != nil {
		return caller, err
	}
	subscriber := newSubscriber(caller)
	user := NewUser(phoneNumber, userInfo.FirstName, subscriber, chatInfos)
	SendMessageToGroup(user, "看看这个")
	Users = append(Users, user)
	// subscribers = append(subscribers, subscriber)
	// SendMessageBack(mconn, subscriber)
	return caller, nil
}

/*func sendMessage() {
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
}*/

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
