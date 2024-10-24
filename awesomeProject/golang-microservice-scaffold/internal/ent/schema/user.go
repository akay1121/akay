//go:generate ent generate .

// Package schema defines the SQL table schema employed by the service. Database manipulation related code
// is implemented by the ent generator, and you should never edit the generated code in the ent package.
//
// AGAIN: This package only contains the definition of the schema, not the implementation!
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			Comment("Unique identifier"),
		field.Int64("parent_id").
			Default(-1).
			Comment("Identifier of the parent user group"),
		field.Int16("type").
			Default(0).
			Comment("User type (0: normal user, 1: user group)"),
		field.String("name").
			MaxLen(64).
			NotEmpty().
			Comment("Username"),
		field.String("password").
			MaxLen(64).
			Sensitive().
			NotEmpty().
			Comment("Hashed password with salt"),
		field.Bytes("salt").
			MaxLen(64).
			NotEmpty().
			Comment("Salt value to ensure security"),
		field.String("two_fa_method").
			Default("disabled").
			MaxLen(16).
			NotEmpty().
			Comment("2FA validation method"),
		field.String("secret").
			Default("").
			MaxLen(128).
			NotEmpty().
			Comment("Secret value of 2FA or just an email address"),
		field.Bool("deleted").
			Default(false).
			Comment(""),
		field.Bytes("login_ip").
			MaxLen(16).
			NotEmpty().
			Comment("IP address of the last login time"),
		field.Time("last_login").
			Comment("Date of the last login time"),
		field.String("nickname").
			Default("").
			MaxLen(64).
			Comment("Optional nickname"),
		field.String("email").
			Default("").
			MaxLen(64).
			Comment("Email address"),
		field.String("phone_number").
			Default("").
			MaxLen(15).
			Comment("Phone number"),
		field.String("avatar").
			Default("").
			MaxLen(128).
			Comment("Relative path of the avatar image"),
		field.Int8("gender").
			Default(0).
			Comment("0: Unknown, 1: Male, 2: Female"),
		field.Time("create_time").
			Default(time.Now).
			Immutable().
			Comment("Creation time for audit purposes"),
		field.Time("last_update").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Last update time of this record for audit purposes"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", User.Type).
			From("parent").
			Field("parent_id").
			Required().
			Unique(),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "password", "salt", "two_fa_method", "secret").
			StorageKey("idx_user_login"),
		index.Fields("parent_id").
			StorageKey("idx_user_tree"),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{
			Table:     "sys_users",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
			Options:   "ENGINE = InnoDB",
		},
		schema.Comment("User information sheet"),
	}
}
