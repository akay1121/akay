//go:generate ent generate .

// Package schema defines the SQL table schema employed by the service. Database manipulation related code
// is implemented by the ent generator, and you should never edit the generated code in the ent package.
//
// AGAIN: This package only contains the definition of the schema, not the implementation!
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

type Sensor struct {
	ent.Schema
}

func (Sensor *Sensor) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Unique(), // id primery key
		field.String("sensor_type"). // sensor-type
			NotEmpty(),
		field.String("identifier"). // 类型识别号
			NotEmpty(),
		field.Time("last_updated"). // latest_time
			Default(time.Now),
	}

}
func (Sensor) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("values", SensorValue.Type)} //one sensor could have multiple values
	// can add more edge
}
