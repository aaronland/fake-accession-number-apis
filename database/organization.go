package database

import (
	"github.com/bwmarrin/snowflake"
)

var snowflake_node *snowflake.Node

func init() {

	node, err := snowflake.NewNode(1)

	if err != nil {
		panic(err)
	}

	snowflake_node = node
}

// NewOrganizationId returns a unique identifier that associate with an organization.
func NewOrganizationId() int64 {
	id := snowflake_node.Generate()
	return id.Int64()
}
