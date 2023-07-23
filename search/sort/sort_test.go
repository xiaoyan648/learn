package sort

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	items := []int{0, 1, 2, 3, 4, 5}
	sort.Slice(items, func(i, j int) bool {
		if items[j] == 0 {
			return true
		}
		return items[i] < items[j]
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5, 0}, items)

	items = []int{0, 1, 2, 3, 4, 5, 6}
	sort.Slice(items, func(i, j int) bool {
		if items[i] == 0 { //
			return false
		}
		if items[j] == 0 {
			return true
		}
		return items[i] < items[j]
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 0}, items)

	items = []int{0, 1, 2, 3, 4, 5, 6}
	sort.SliceStable(items, func(i, j int) bool {
		if items[j] == 0 {
			return true
		}
		return items[i] < items[j]
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 0}, items)
}

func Test(t *testing.T) {
	now := time.Now()
	utcTime := now.UTC()
	fmt.Printf(" now: %v, utc: %v \n", now, utcTime)
	fmt.Printf(" now: %v, utc: %v \n", now.Unix(), utcTime.Unix())
}
