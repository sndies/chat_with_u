package service

import (
	"context"
	"github.com/sndies/chat_with_u/middleware/cache"
	"testing"
)

func TestAddKefuAccount(t *testing.T) {
	type args struct {
		ctx         context.Context
		accessToken string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "测试生成客服账号",
			args: args{
				ctx:         context.Background(),
				accessToken: "69_IKVTpLBu4-MwrxQx7idXYOb4LePysPTAGJCWOFrz0sc7bA99DeZfvtvWXNtgGPlkiTR2OWEWpV2Nnu-E6xmxAA0uqC7NsRVUuYj93hXfUnwStt1rUxGBihORZ6cRVXhAIABFH",
			},
		},
	}
	for _, tt := range tests {
		cache.Init()
		t.Run(tt.name, func(t *testing.T) {
			AddKefuAccount(tt.args.ctx)
		})
	}
}
