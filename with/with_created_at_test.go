package with

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreatedOn_Validate(t *testing.T) {
	type fields struct {
		CreatedAt string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			fields: fields{
				CreatedAt: "2020-12-31",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err, i...)
			},
		},
		{
			name: "missing",
			fields: fields{
				CreatedAt: "",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i...)
			},
		},
		{
			name: "invalid_letters",
			fields: fields{
				CreatedAt: "abc",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i...)
			},
		},
		{
			name: "invalid_no_dashes",
			fields: fields{
				CreatedAt: "20201231",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i...)
			},
		},
		{
			name: "invalid_slashes",
			fields: fields{
				CreatedAt: "2020/12/31",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i...)
			},
		},
		{
			name: "invalid_us",
			fields: fields{
				CreatedAt: "31/12/2020",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &CreatedAtField{}
			if tt.fields.CreatedAt != "" {
				v.CreatedAt, _ = time.Parse(time.DateOnly, tt.fields.CreatedAt)
			}
			tt.wantErr(t, v.Validate(), fmt.Sprintf("{CreatedAtField=%s}.Validate()", v.CreatedAt))
		})
	}
}
