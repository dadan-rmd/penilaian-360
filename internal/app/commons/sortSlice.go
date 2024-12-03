package commons

import (
	"sort"

	"github.com/spf13/cast"
)

func SortSlice(arr []string) []string {
	arr1 := arr
	sort.Slice(arr1, func(i1, i2 int) bool {
		numA := cast.ToInt(arr1[i1])
		numB := cast.ToInt(arr1[i2])
		return numA < numB
	})
	return arr1
}
