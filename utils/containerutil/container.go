package containerutil

func HasIntersection[E comparable](a, b []E) bool {
	aMap := make(map[E]bool)

	for i := 0; i < len(a); i++ {
		aMap[a[i]] = true
	}

	for i := 0; i < len(b); i++ {
		if aMap[b[i]] {
			return true
		}
	}

	return false
}

// 并集
func Union[E comparable](a, b []E) []E {
	m := make(map[E]struct{})
	var union []E

	for _, item := range a {
		m[item] = struct{}{}
		union = append(union, item)
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			union = append(union, item)
			m[item] = struct{}{}
		}
	}

	return union
}

// 交集 通用, 返回顺序与 a 保持一致.
func IntersectMap(a, b []uint64) []uint64 {
	m := make(map[uint64]bool)
	var intersection []uint64

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if m[item] {
			intersection = append(intersection, item)
			m[item] = false
		}
	}

	return intersection
}

// 交集 用于已排好序的元素集合.
func IntersectSorted(a, b []uint64) []uint64 {
	var result []uint64

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				result = append(result, a[i])
				break
			} else if a[i] < b[j] {
				break
			}
		}
	}

	return result
}

// DiffArray 求A再B中的差集.
// A[1,2] B[2,4] -> [1].
func DiffArray[E comparable](a, b []E) []E {
	var diffArray []E
	temp := map[E]struct{}{}

	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			diffArray = append(diffArray, val)
		}
	}

	return diffArray
}
