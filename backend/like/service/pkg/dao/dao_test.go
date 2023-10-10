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
	like, err := d.GetLikeState(ctx, 20, 1)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(like)
}
