package pool

import (
	"github.com/funny/slab"
)

var (
	slicePool = slab.NewChanPool(
		16,        // The smallest chunk size is 16B.
		64*1024,   // The largest chunk size is 64KB.
		2,         // Power of 2 growth in chunk size.
		1024*1024, // Each slab will be 1MB in size.
	)
)

//RecycleSlice 回收
func RecycleSlice(slice []byte) {
	slicePool.Free(slice)
}

//GetSlice 获取流
func GetSlice(s int) []byte {
	return slicePool.Alloc(s)
}
