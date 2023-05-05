package containerutil

import (
	"reflect"
	"testing"
)

func TestIntersection(t *testing.T) {
	s1 := []uint64{1, 3, 5, 7, 9, 10, 11, 640}
	s2 := []uint64{2, 9, 6, 8, 3, 11, 640}
	expected := []uint64{3, 9, 11, 640}
	res := IntersectMap(s1, s2)
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", res, expected)
	}
}

func TestIntersectSorted(t *testing.T) {
	s1 := []uint64{1, 3, 5, 7, 9, 640}
	s2 := []uint64{2, 3, 6, 8, 9, 640}
	expected := []uint64{3, 9, 640}
	res := IntersectSorted(s1, s2)
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", res, expected)
	}
	s3 := []uint64{2, 9, 6, 8, 3}
	expected = []uint64{9} // s3 未排序会漏掉元素
	res = IntersectSorted(s1, s3)
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", res, expected)
	}
}

func TestHasIntersection(t *testing.T) {
	testCases := []struct {
		a          []uint64
		b          []uint64
		wantResult bool
	}{
		{[]uint64{1, 2, 3}, []uint64{4, 5, 6}, false},
		{[]uint64{1, 2, 3}, []uint64{2, 4, 6}, true},
		{[]uint64{}, []uint64{1, 2, 3}, false},
		{[]uint64{}, []uint64{}, false},
		{[]uint64{2, 1}, []uint64{2, 3}, true},
		{[]uint64{2, 1, 10}, []uint64{10, 3}, true},
	}

	for _, tc := range testCases {
		gotResult := HasIntersection(tc.a, tc.b)
		if gotResult != tc.wantResult {
			t.Errorf("hasIntersection(%v, %v) = %v, want %v", tc.a, tc.b, gotResult, tc.wantResult)
		}
	}
}

// BenchmarkIntersectMap-8   	 1914442	       625.3 ns/op	     299 B/op	       5 allocs/op
func BenchmarkIntersectMap(b *testing.B) {
	a := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bb := []uint64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntersectMap(a, bb)
	}
}

// BenchmarkIntersectSorted-8   	 9375063	       120.3 ns/op	     120 B/op	       4 allocs/op
func BenchmarkIntersectSorted(b *testing.B) {
	a := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bb := []uint64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntersectSorted(a, bb)
	}
}
