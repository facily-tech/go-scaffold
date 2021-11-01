package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Quote holds the schema definition for the Quote entity.
type Quote struct {
	ent.Schema
}

// Fields of the Quote.
func (Quote) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("content"),
	}
}

// Edges of the Quote.
func (Quote) Edges() []ent.Edge {
	return nil
}
