package memory

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func newLogger() zerolog.Logger {
	out := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05 -0700",
		NoColor:    true,
	}

	l := zerolog.New(out)

	return l.Level(zerolog.InfoLevel).With().Timestamp().Logger()
}

func TestFileStore_Get(t *testing.T) {
	type fields struct {
		memoryStore   StoreMemory
		storeRestore  bool
		storeInterval time.Duration
		storePath     string
	}
	type args struct {
		ctx context.Context
	}

	data := make(map[string]Type)
	gauge := 1.0
	counter := int64(1)
	data["test"] = Type{
		Gauge:   &gauge,
		Counter: &counter,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]Type
		wantErr bool
	}{
		{
			name: "positive",
			fields: fields{
				memoryStore: StoreMemory{
					metric: data,
					log:    newLogger(),
				},
				storeRestore:  false,
				storeInterval: 0,
				storePath:     "",
			},
			args: args{
				ctx: context.Background(),
			},
			want:    data,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileStore{
				memoryStore:   tt.fields.memoryStore,
				storeRestore:  tt.fields.storeRestore,
				storeInterval: tt.fields.storeInterval,
				storePath:     tt.fields.storePath,
			}
			got, err := f.Get(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_GetOne(t *testing.T) {
	type fields struct {
		memoryStore   StoreMemory
		storeRestore  bool
		storeInterval time.Duration
		storePath     string
	}
	type args struct {
		ctx  context.Context
		name string
	}

	data := make(map[string]Type)
	gauge := 1.0
	counter := int64(1)
	data["test"] = Type{
		Gauge:   &gauge,
		Counter: &counter,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Type
		wantErr bool
	}{
		{
			name: "positive",
			fields: fields{
				memoryStore: StoreMemory{
					metric: data,
					log:    newLogger(),
				},
				storeRestore:  false,
				storeInterval: 0,
				storePath:     "",
			},
			args: args{
				ctx:  context.Background(),
				name: "test",
			},
			want: Type{
				Gauge:   &gauge,
				Counter: &counter,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileStore{
				memoryStore:   tt.fields.memoryStore,
				storeRestore:  tt.fields.storeRestore,
				storeInterval: tt.fields.storeInterval,
				storePath:     tt.fields.storePath,
			}
			got, err := f.GetOne(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_Ping(t *testing.T) {
	type fields struct {
		memoryStore   StoreMemory
		storeRestore  bool
		storeInterval time.Duration
		storePath     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "positive",
			fields: fields{
				memoryStore: StoreMemory{
					metric: make(map[string]Type),
					log:    newLogger(),
				},
				storeRestore:  false,
				storeInterval: 0,
				storePath:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileStore{
				memoryStore:   tt.fields.memoryStore,
				storeRestore:  tt.fields.storeRestore,
				storeInterval: tt.fields.storeInterval,
				storePath:     tt.fields.storePath,
			}
			if err := f.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileStore_Update(t *testing.T) {
	type fields struct {
		memoryStore   StoreMemory
		storeRestore  bool
		storeInterval time.Duration
		storePath     string
	}
	type args struct {
		ctx  context.Context
		name string
		data map[string]Type
	}

	data := make(map[string]Type)
	gauge := 1.0
	counter := int64(1)
	data["test"] = Type{
		Gauge:   &gauge,
		Counter: &counter,
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive",
			fields: fields{
				memoryStore: StoreMemory{
					metric: data,
					log:    newLogger(),
				},
				storeRestore:  false,
				storeInterval: 0,
				storePath:     "./test.json",
			},
			args: args{
				ctx:  context.Background(),
				name: "test",
				data: data,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileStore{
				memoryStore:   tt.fields.memoryStore,
				storeRestore:  tt.fields.storeRestore,
				storeInterval: tt.fields.storeInterval,
				storePath:     tt.fields.storePath,
			}
			f.Update(tt.args.ctx, tt.args.name, tt.args.data)
		})
	}
	if err := os.Remove("./test.json"); err != nil {
		log.Err(err).Msg("os.Remove")
	}
}

func TestFileStore_Updates(t *testing.T) {
	type fields struct {
		memoryStore   StoreMemory
		storeRestore  bool
		storeInterval time.Duration
		storePath     string
	}
	type args struct {
		ctx context.Context
		req []Metric
	}

	data := make(map[string]Type)
	gauge := 1.0
	counter := int64(1)
	data["test"] = Type{
		Gauge:   &gauge,
		Counter: &counter,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "positive",
			fields: fields{
				memoryStore: StoreMemory{
					metric: data,
					log:    newLogger(),
				},
				storeRestore:  false,
				storeInterval: 0,
				storePath:     "./test.json",
			},
			args: args{
				ctx: context.Background(),
				req: []Metric{
					{
						Name:  "test",
						Value: &gauge,
						Delta: &counter,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileStore{
				memoryStore:   tt.fields.memoryStore,
				storeRestore:  tt.fields.storeRestore,
				storeInterval: tt.fields.storeInterval,
				storePath:     tt.fields.storePath,
			}
			if err := f.Updates(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Updates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if err := os.Remove("./test.json"); err != nil {
		log.Err(err).Msg("os.Remove")
	}
}
