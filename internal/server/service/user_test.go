package service

import (
	"context"
	"errors"
	"go-template/internal/server/cache"
	"go-template/internal/server/model"
	"go-template/internal/server/repository"
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

var users = []map[string]interface{}{
	{
		"id":   "1",
		"name": "A",
	},
	{
		"id":   "2",
		"name": "B",
	},
}

func get(userID string) map[string]interface{} {
	for _, item := range users {
		if item["id"] == userID {
			return item
		}
	}
	return nil
}

type mockUserRepo struct{}

func (r *mockUserRepo) Get(userID string) (*model.User, error) {
	var user *model.User
	v := get(userID)
	if v != nil {
		if err := mapstructure.Decode(v, &user); err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, errors.New("Not Found")
}

type mockUserCache struct{}

func (c *mockUserCache) Get(userID string) (*model.UserCache, error) {
	var userCache *model.UserCache
	v := get(userID)
	if v != nil {
		if err := mapstructure.Decode(v, &userCache); err != nil {
			return nil, err
		}
		return userCache, nil
	}
	return nil, errors.New("Not Found")
}

func (c *mockUserCache) Set(userID string, user *model.UserCache) error {
	panic("not implemented") // TODO: Implement
}

func Test_userService_Get(t *testing.T) {
	type fields struct {
		repo   repository.UserRepo
		cache  cache.UserCache
		logger *zap.Logger
	}
	type args struct {
		userID string
	}
	type test struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}
	f := fields{repo: &mockUserRepo{}, cache: &mockUserCache{}, logger: zap.NewNop()}
	tests := []test{
		{
			name:    "Get 1",
			fields:  f,
			args:    args{userID: "1"},
			want:    &model.User{ID: "1", Name: "A"},
			wantErr: false,
		},
		{
			name:    "Get 2",
			fields:  f,
			args:    args{userID: "2"},
			want:    &model.User{ID: "2", Name: "B"},
			wantErr: false,
		},
		{
			name:    "Get 3",
			fields:  f,
			args:    args{userID: "3"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserService(tt.fields.repo, tt.fields.cache, tt.fields.logger)
			got, err := s.Get(context.TODO(), tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
