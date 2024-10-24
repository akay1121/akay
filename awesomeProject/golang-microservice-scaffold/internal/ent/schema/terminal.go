package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// Terminal holds the schema definition for the Terminal entity.
type Terminal struct {
	ent.Schema
}

// Fields of the Terminal.
func (Terminal) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("status").Default("offline"),
		field.Int("timeout").Default(60),             // default timeout time is 60s
		field.Time("last_updated").Default(time.Now), // former update time
	}
}
