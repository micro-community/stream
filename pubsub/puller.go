/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-06 10:14:25
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-06 10:38:46
 * @FilePath: \stream\pubsub\puller.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package pubsub

type IPull interface {
	IPublish
	Connect() error
	Pull()
	Reconnect() bool
	init(streamPathUrl string, url string, option *PullOption)
}

//Puller for remote pull stream
type Puller struct {
	Client[PullOption]
}
