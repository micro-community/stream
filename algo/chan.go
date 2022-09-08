/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-06 11:05:32
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-08 21:47:37
 * @FilePath: \stream\algo\chan.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package algo

import (
	"math"
	"sync/atomic"
)

// SafeChan for safety channel，可以防止close后被写入的问题
type SafeChan[T any] struct {
	C       chan T
	senders int32 // senders counts
	len     int
}

// Init a n length chan
func (sc *SafeChan[T]) Init(n int) {
	sc.C = make(chan T, n)
	sc.len = n
}

// Close chan safely senders
func (sc *SafeChan[T]) Close() bool {
	if atomic.CompareAndSwapInt32(&sc.senders, 0, math.MinInt32) {
		close(sc.C)
		return true
	}
	return false
}

func (sc *SafeChan[T]) Send(v T) bool {
	// senders增加后为正数说明有channel未被关闭，可以发送数据
	if atomic.AddInt32(&sc.senders, 1) > 0 {
		sc.C <- v
		//发送后计数减1
		atomic.AddInt32(&sc.senders, -1)
		return true
	}
	return false
}

func (sc *SafeChan[T]) IsClosed() bool {
	return atomic.LoadInt32(&sc.senders) < 0
}

func (sc *SafeChan[T]) IsEmpty() bool {
	return atomic.LoadInt32(&sc.senders) == 0
}

// IsFull return sender exist
func (sc *SafeChan[T]) IsFull() bool {
	return atomic.LoadInt32(&sc.senders) > 0
}

func (sc *SafeChan[T]) ChanCount() int {
	return sc.len
}
