package quote

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Quote struct {
	ID      uuid.UUID
	Content string
}

func NewQuote(id *uuid.UUID, content string) (Quote, error) {
	if content == "" {
		return Quote{}, errors.Wrap(ErrNew, "empty code not allowed")
	}

	if id != nil {
		return Quote{ID: *id, Content: content}, nil
	}

	iid, err := uuid.NewRandom()
	if err != nil {
		return Quote{}, errors.Wrap(ErrNew, err.Error())
	}

	return Quote{ID: iid, Content: content}, nil
}
