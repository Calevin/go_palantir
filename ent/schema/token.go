package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Token holds the schema definition for the Token entity.
type Token struct {
	ent.Schema
}

// Fields of the Token.
func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.Int("line").
			Comment("Número de línea donde se encontró el token"),
		field.Int("order").
			Comment("Orden del token en la línea"),
		field.String("token").
			Comment("Contenido del token"),
	}
}

// Edges of the Token.
func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		// Define la relación inversa: Cada Token pertenece a un File.
		edge.From("file", File.Type).
			Ref("tokens").
			Unique(),
	}
}
