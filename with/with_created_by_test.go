package with

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatedBy_UpdatesCreatedBy(t *testing.T) {
	type fields struct {
		CreatedBy string
	}
	tests := []struct {
		name   string
		fields fields
		want   []dal.Update
	}{
		{
			name: "ok",
			fields: fields{
				CreatedBy: "u1",
			},
			want: []dal.Update{
				{Field: "createdBy", Value: "u1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &CreatedByField{
				CreatedBy: tt.fields.CreatedBy,
			}
			assert.Equalf(t, tt.want, v.UpdatesCreatedBy(), "UpdatesCreatedBy()")
		})
	}
}

func TestCreatedBy_Validate(t *testing.T) {
	type fields struct {
		CreatedBy string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			fields: fields{
				CreatedBy: "u1",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err, i...)
			},
		},
		{
			name: "missing",
			fields: fields{
				CreatedBy: "",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &CreatedByField{
				CreatedBy: tt.fields.CreatedBy,
			}
			tt.wantErr(t, v.Validate(), fmt.Sprintf("{CreatedByField=%s}Validate()", v.CreatedBy))
		})
	}
}
