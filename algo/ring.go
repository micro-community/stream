/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-02 12:47:32
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-02 16:51:15
 * @FilePath: \stream\algo\ring.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package algo

// Ring for anytype
type Ring[T any] struct {
	next, prev *Ring[T]
	Value      T // used by client
}

func (r *Ring[T]) init() *Ring[T] {
	r.next = r
	r.prev = r
	return r
}
