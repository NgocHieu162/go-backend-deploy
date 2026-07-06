package schema

import (
	softdelete "go-backend/ent/soft-delete"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ChatGroups holds the schema definition for the ChatGroups entity.
type ChatGroups struct {
	ent.Schema
}

// Fields of the ChatGroups.
func (ChatGroups) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").StructTag(`json:"userId"`),
		field.String("Name").Optional().Nillable().StructTag(`json:"name"`),

		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the ChatGroups.
func (ChatGroups) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ChatGroupMembers", ChatGroupMembers.Type),
		edge.To("ChatMessages", ChatMessages.Type),

		edge.From("Users", Users.Type).Ref("ChatGroups").Field("user_id").Unique().Required(),
	}
}

func (ChatGroups) Mixin() []ent.Mixin {
	return []ent.Mixin{
		softdelete.SoftDeleteMixin{},
	}
}
