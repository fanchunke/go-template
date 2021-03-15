package cache

import (
	"go-template/internal/server/model"
	"reflect"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
)

func Test_userCache_Get(t *testing.T) {

	// miniredis for unittest
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	s.HSet("1", "id", "1", "name", "A")
	s.HSet("2", "id", "2", "name", "B")
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", s.Addr())
		},
	}

	type fields struct {
		pool *redis.Pool
	}
	type args struct {
		userID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.UserCache
		wantErr bool
	}{
		{
			name:    "Get 1",
			fields:  fields{pool: pool},
			args:    args{userID: "1"},
			want:    &model.UserCache{ID: "1", Name: "A"},
			wantErr: false,
		},
		{
			name:    "Get 2",
			fields:  fields{pool: pool},
			args:    args{userID: "2"},
			want:    &model.UserCache{ID: "2", Name: "B"},
			wantErr: false,
		},
		{
			name:    "Get 3",
			fields:  fields{pool: pool},
			args:    args{userID: "3"},
			want:    &model.UserCache{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &userCache{
				pool: tt.fields.pool,
			}
			got, err := c.Get(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("userCache.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userCache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userCache_Set(t *testing.T) {

	// miniredis for unittest
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", s.Addr())
		},
	}

	type fields struct {
		pool *redis.Pool
	}
	type args struct {
		userID string
		user   *model.UserCache
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Set 1",
			fields:  fields{pool: pool},
			args:    args{userID: "1", user: &model.UserCache{ID: "1", Name: "A"}},
			wantErr: false,
		},
		{
			name:    "Set 2",
			fields:  fields{pool: pool},
			args:    args{userID: "2", user: &model.UserCache{ID: "2", Name: "B"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &userCache{
				pool: tt.fields.pool,
			}
			if err := c.Set(tt.args.userID, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userCache.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			if result := s.HGet(tt.args.userID, "id"); result != tt.args.user.ID {
				t.Errorf("userCache.Get got = %v, want %v", result, tt.args.user.ID)
			}
			if result := s.HGet(tt.args.userID, "name"); result != tt.args.user.Name {
				t.Errorf("userCache.Get got = %v, want %v", result, tt.args.user.Name)
			}
		})
	}
}
