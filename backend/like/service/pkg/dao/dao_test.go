package dao

import (
	"context"
	"fmt"
	"testing"
)

func TestDao(t *testing.T) {
	m, err := NewMysql()
	if err != nil {
		fmt.Println(err.Error())
	}
	ctx := context.Background()
	d := &dao{mysql: m}
	err = d.AddLikeState(ctx, 1, 1, 1)
	if err != nil {
		fmt.Println(err.Error())
	}
}
