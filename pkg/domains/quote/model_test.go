package quote

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewQuote(t *testing.T) {
	type args struct {
		id      *uuid.UUID
		content string
	}
	tests := []struct {
		name string
		args args
		want Quote
		err  error
	}{
		{
			name: "success, create a new quote without id",
			args: args{id: nil, content: testQuote.Content},
			want: testQuote,
			err:  nil,
		},
		{
			name: "success, create a new quote with id",
			args: args{id: &testQuote.ID, content: testQuote.Content},
			want: testQuote,
			err:  nil,
		},
		{
			name: "fail, try to create quote",
			args: args{id: nil, content: ""},
			want: Quote{},
			err:  ErrNew,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQuote(tt.args.id, tt.args.content)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.want.Content, got.Content)
		})
	}
}
