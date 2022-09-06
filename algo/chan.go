/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-06 11:05:32
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-06 11:06:06
 * @FilePath: \stream\algo\chan.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package algo

// SafeChan for safety channel，可以防止close后被写入的问题
type SafeChan[T any] struct {
	C       chan T
	senders int32 // senders counts
}
