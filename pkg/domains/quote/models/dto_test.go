package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestQuote_toJSONResponse(t *testing.T) {
	type fields struct {
		ID      uuid.UUID
		Content string
	}
	tests := []struct {
		name   string
		fields fields
		want   JSONResponse
	}{
		{
			name:   "success",
			fields: fields(TestQuote),
			want:   JSONResponse{ID: &TestQuote.ID, Content: TestQuote.Content},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := Quote{
				ID:      tt.fields.ID,
				Content: tt.fields.Content,
			}
			got := q.ToJSONResponse()
			assert.Equal(t, tt.want, got)
		})
	}
}
