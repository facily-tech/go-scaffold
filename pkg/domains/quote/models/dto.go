package models

import "github.com/google/uuid"

type JSONResponse struct {
	ID      *uuid.UUID `json:"id,omitempty"`
	Content string     `json:"content,omitempty"`
	Error   string     `json:"error,omitempty"`
}

func (q Quote) ToJSONResponse() JSONResponse {
	return JSONResponse{
		ID:      &q.ID,
		Content: q.Content,
	}
}
