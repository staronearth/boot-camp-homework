package week1

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteAt(t *testing.T) {
	testcase := []struct {
		name string

		inputSlice []any
		inputIdx   int

		wantRes []any
		wantErr error
	}{
		{
			name:       "nil",
			inputSlice: nil,
			inputIdx:   13,

			wantRes: nil,
			wantErr: errors.New("slice不能为空或者不能nil"),
		},
		{
			name:       "nil",
			inputSlice: []any{},
			inputIdx:   13,

			wantRes: nil,
			wantErr: errors.New("slice不能为空或者不能nil"),
		},
		{
			name:       "一个元素",
			inputSlice: []any{1},
			inputIdx:   0,

			wantRes: []any{},
			wantErr: nil,
		},
		{
			name:       "错误下标",
			inputSlice: []any{1, 2, 4},
			inputIdx:   13,

			wantRes: nil,
			wantErr: ErrIndexOutOfRange,
		},

		{
			name:       "中间元素",
			inputSlice: []any{1, 2, 4, 10},
			inputIdx:   2,

			wantRes: []any{1, 2, 10},
			wantErr: nil,
		},
		{
			name:       "nil",
			inputSlice: []any{1, 2, 4, 10},
			inputIdx:   3,

			wantRes: []any{1, 2, 4},
			wantErr: nil,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := DeleteAt(tc.inputSlice, tc.inputIdx)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestDeleteAtSlice(t *testing.T) {
	testcase := []struct {
		name string

		inputSlice []any
		inputIdx   int

		wantRes []any
		wantErr error
	}{
		{
			name:       "nil",
			inputSlice: nil,
			inputIdx:   13,

			wantRes: nil,
			wantErr: errors.New("slice不能为空或者不能nil"),
		},
		{
			name:       "nil",
			inputSlice: []any{},
			inputIdx:   13,

			wantRes: nil,
			wantErr: errors.New("slice不能为空或者不能nil"),
		},
		{
			name:       "一个元素",
			inputSlice: []any{1},
			inputIdx:   0,

			wantRes: []any{},
			wantErr: nil,
		},
		{
			name:       "错误下标",
			inputSlice: []any{1, 2, 4},
			inputIdx:   13,

			wantRes: nil,
			wantErr: ErrIndexOutOfRange,
		},

		{
			name:       "中间元素",
			inputSlice: []any{1, 2, 4, 10},
			inputIdx:   2,

			wantRes: []any{1, 2, 10},
			wantErr: nil,
		},
		{
			name:       "删除最后一个元素",
			inputSlice: []any{1, 2, 4, 10},
			inputIdx:   3,

			wantRes: []any{1, 2, 4},
			wantErr: nil,
		},
		{
			name:       "删除第一个元素",
			inputSlice: []any{1, 2, 4, 10},
			inputIdx:   0,

			wantRes: []any{2, 4, 10},
			wantErr: nil,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := DeleteAtSlice(tc.inputSlice, tc.inputIdx)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
