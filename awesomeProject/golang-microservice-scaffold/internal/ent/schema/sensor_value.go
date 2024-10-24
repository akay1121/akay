//go:generate ent generate .

// Package schema defines the SQL table schema employed by the service. Database manipulation related code
// is implemented by the ent generator, and you should never edit the generated code in the ent package.
//
// AGAIN: This package only contains the definition of the schema, not the implementation!
package schema

import (
	"ariga.io/atlas/sql/schema"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// SensorValue
type SensorValue struct {
	schema.Schema
}

// Fields of the SensorValue.
func (SensorValue) Fields() []ent.Field {
	return []ent.Field{
		field.Float("value").
			Optional(), // 传感器值
		field.Time("timestamp").
			Default(time.Now), // 时间戳
	}
}

// Edges of the SensorValue.
func (SensorValue) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sensor", Sensor.Type).
			Ref("values").
			Required().
			Unique(),
	}
}
