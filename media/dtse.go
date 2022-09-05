/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-05 17:23:16
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-05 17:23:48
 * @FilePath: \stream\media\dtse.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package media

// DTSEstimator is a DTS estimator.
type DTSEstimator struct {
	prevDTS     uint32
	prevPTS     uint32
	prevPrevPTS uint32
	dts         func(uint32) uint32
	delta       uint32
}
