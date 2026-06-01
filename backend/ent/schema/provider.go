package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Provider holds the schema definition for the system-level Provider entity.
// Managed by admin; users can only use providers that are enabled.
type Provider struct {
	ent.Schema
}

// Fields of the Provider.
func (Provider) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Positive(),
		field.String("name").NotEmpty(),          // 显示名称，如 阿里云
		field.String("type").Unique().NotEmpty(), // 内部标识，如 aliyun
		field.String("description").Optional(),   // 简介
		field.Bool("enabled").Default(true),      // 是否启用
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the Provider.
func (Provider) Edges() []ent.Edge {
	return nil
}
