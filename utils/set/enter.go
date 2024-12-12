package set

// Union 并集
func Union[T uint | int | string](slice1, slice2 []T) []T {
	m := make(map[T]int)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// Intersect 求交集
func Intersect[T string | uint | int | uint32](slice1, slice2 []T) []T {
	m := make(map[T]int)
	nn := make([]T, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		time, _ := m[v]
		if time == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}
func Difference[T uint | int | string](slice1, slice2 []T) []T {
	m := make(map[T]int)
	nn := make([]T, 0)
	inter := Intersect(slice1, slice2)

	// 计算交集元素在 slice1 中出现的次数
	for _, v := range inter {
		m[v]++
	}

	// 遍历 slice1，找出不在交集中的元素
	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}

	return nn
}
