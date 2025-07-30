package osquery

import (
	"context"
	"github.com/osquery/osquery-go"
	"reflect"
	"testing"
	"time"
)

func TestNewQuerier(t *testing.T) {
	type args struct {
		socketPath string
		timeout    time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *Querier
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQuerier(tt.args.socketPath, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewQuerier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQuerier() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuerier_Close(t *testing.T) {
	type fields struct {
		client *osquery.ExtensionManagerClient
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Querier{
				client: tt.fields.client,
			}
			q.Close()
		})
	}
}

func TestQuerier_Query(t *testing.T) {
	type fields struct {
		client *osquery.ExtensionManagerClient
	}
	type args struct {
		ctx context.Context
		sql string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Querier{
				client: tt.fields.client,
			}
			got, err := q.Query(tt.args.ctx, tt.args.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query() got = %v, want %v", got, tt.want)
			}
		})
	}
}
