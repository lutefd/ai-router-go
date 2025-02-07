package idgen

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(workerID int64) error {
	var err error
	node, err = snowflake.NewNode(workerID)
	return err
}

func Generate() string {
	return fmt.Sprintf("chat_%d", node.Generate())
}
