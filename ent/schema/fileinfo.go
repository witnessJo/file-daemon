package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type FileInfo struct {
	ent.Schema
}

// Fields of the Block.
func (FileInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("node_name").NotEmpty().Comment("Name of the node"),
		field.String("mount_path").NotEmpty().Comment("Mount path of the file"),
		field.JSON("file_list", []string{}).Default([]string{}).Comment("List of files in JSON format"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("Creation time of the record"),
	}
}

// Edges of the Block.
func (FileInfo) Edges() []ent.Edge {
	return []ent.Edge{}
}
