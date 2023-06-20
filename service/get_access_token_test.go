package service

import (
	"context"
	"github.com/sndies/chat_with_u/middleware/cache"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "测试获取access_token",
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.Init()
			got, err := GetAccessToken(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("GetAccessToken() = %v", got)
			}
		})
	}
}
