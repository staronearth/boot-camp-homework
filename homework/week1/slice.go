package week1

import "errors"

var ErrIndexOutOfRange = errors.New("下标超出范围")

// DeleteAt 删除指定位置的元素
// 如果下标不是合法的下标，返回 ErrIndexOutOfRange
func DeleteAt[T any](s []T, idx int) ([]T, error) {
	if len(s) == 0 {
		return nil, errors.New("slice不能为空或者不能nil")
	}
	//panic("implement me")
	if idx < 0 || idx >= len(s) {
		return nil, ErrIndexOutOfRange
	}
	var res []T = make([]T, 0)
	for i := 0; i < len(s); i++ {
		if i == idx {
			continue
		}
		res = append(res, s[i])
	}
	return res, nil
}

func DeleteAtSlice[T any](s []T, idx int) ([]T, error) {
	if len(s) == 0 {
		return nil, errors.New("slice不能为空或者不能nil")
	}
	//panic("implement me")
	if idx < 0 || idx >= len(s) {
		return nil, ErrIndexOutOfRange
	}
	s = append(s[:idx], s[idx+1:]...)
	s = s[:len(s)]
	return s, nil
}
