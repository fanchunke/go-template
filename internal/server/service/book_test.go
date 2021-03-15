package service

import (
	"context"
	"errors"
	"go-template/internal/server/model"
	"go-template/internal/server/repository"
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
)

var books = []map[string]interface{}{
	{
		"id":   "1",
		"name": "A",
	},
	{
		"id":   "2",
		"name": "B",
	},
}

type mockBookRepo struct{}

func (r *mockBookRepo) Get(bookID string) (*model.Book, error) {
	var book *model.Book
	for _, item := range books {
		if item["id"] == bookID {
			if err := mapstructure.Decode(item, &book); err != nil {
				return nil, err
			}
			return book, nil
		}
	}
	return book, errors.New("Not Found")
}

func Test_bookService_Get(t *testing.T) {
	type fields struct {
		repo repository.BookRepo
	}
	type args struct {
		bookID string
	}
	type test struct {
		name    string
		fields  fields
		args    args
		want    *model.Book
		wantErr bool
	}
	repo := &mockBookRepo{}
	tests := []test{
		{
			name:    "Get 1",
			fields:  fields{repo: repo},
			args:    args{bookID: "1"},
			want:    &model.Book{ID: "1", Name: "A"},
			wantErr: false,
		},
		{
			name:    "Get 2",
			fields:  fields{repo: repo},
			args:    args{bookID: "2"},
			want:    &model.Book{ID: "2", Name: "B"},
			wantErr: false,
		},
		{
			name:    "Get 3",
			fields:  fields{repo: repo},
			args:    args{bookID: "3"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &bookService{
				repo: tt.fields.repo,
			}
			got, err := s.Get(context.TODO(), tt.args.bookID)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
