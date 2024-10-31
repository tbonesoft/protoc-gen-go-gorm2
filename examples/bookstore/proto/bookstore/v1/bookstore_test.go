package bookstore_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/lib/pq"

	bookstore "github.com/tbonesoft/protoc-gen-go-gorm2/gen/bookstore/v1"
)

func TestBookORM_TableName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
		{
			name: `Custom table name "my_book"`,
			want: "my_book",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &bookstore.BookORM{}
			if got := x.TableName(); got != tt.want {
				t.Errorf("BookORM.TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookORM_ToPB(t *testing.T) {
	sampleBigIntId := uint64(1848305773055053824)
	sampleInt32 := int32(1345)

	type fields struct {
		Id      uint64
		Title   string
		TitleCn *int32
		Tags    pq.StringArray
		Note    []byte
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRs  *bookstore.Book
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test copy BigInt/uint64 pq.StringArray/[]string and []byte data type values",
			fields: fields{
				Id:      sampleBigIntId,
				Title:   "_title",
				TitleCn: &sampleInt32,
				Tags:    pq.StringArray{"tag_1", "tag_b"},
				Note:    []byte(`a note`),
			},
			args: args{
				ctx: context.TODO(),
			},
			wantRs: &bookstore.Book{
				Id:      sampleBigIntId,
				Title:   "_title",
				TitleCn: &sampleInt32,
				Tags:    []string{"tag_1", "tag_b"},
				Note:    []byte(`a note`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &bookstore.BookORM{
				Id:      tt.fields.Id,
				Title:   tt.fields.Title,
				TitleCn: tt.fields.TitleCn,
				Tags:    tt.fields.Tags,
				Note:    tt.fields.Note,
			}
			gotRs, err := x.ToPB(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("BookORM.ToPB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRs, tt.wantRs) {
				t.Errorf("BookORM.ToPB() = %v, want %v", gotRs, tt.wantRs)
			}
		})
	}
}

func TestBook_ToORM(t *testing.T) {
	sampleBigIntId := uint64(1848305773055053824)
	sampleInt32 := int32(1345)

	type fields struct {
		Id      uint64
		Title   string
		TitleCn *int32
		Tags    []string
		Note    []byte
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRs  *bookstore.BookORM
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test copy BigInt/uint64 pq.StringArray/[]string and []byte data type values",
			fields: fields{
				Id:      sampleBigIntId,
				Title:   "_title",
				TitleCn: &sampleInt32,
				Tags:    []string{"tag_1", "tag_b"},
				Note:    []byte(`a note`),
			},
			args: args{
				ctx: context.TODO(),
			},
			wantRs: &bookstore.BookORM{
				Id:      sampleBigIntId,
				Title:   "_title",
				TitleCn: &sampleInt32,
				Tags:    pq.StringArray{"tag_1", "tag_b"},
				Note:    []byte(`a note`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &bookstore.Book{
				Id:      tt.fields.Id,
				Title:   tt.fields.Title,
				TitleCn: tt.fields.TitleCn,
				Tags:    tt.fields.Tags,
				Note:    tt.fields.Note,
			}
			gotRs, err := x.ToORM(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Book.ToORM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRs, tt.wantRs) {
				t.Errorf("Book.ToORM() = %v, want %v", gotRs, tt.wantRs)
			}
		})
	}
}
