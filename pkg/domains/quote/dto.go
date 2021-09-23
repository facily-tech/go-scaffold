package quote

func (q Quote) toJSONResponse() JSONResponse {
	return JSONResponse{
		ID:      &q.ID,
		Content: q.Content,
	}
}
