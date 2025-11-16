package utils

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
	once sync.Once
)

// InitSnowflake 初始化雪花算法节点
// nodeID 是 0-1023 之间的整数，用于区分不同的服务实例
func InitSnowflake(nodeID int64) (err error) {
	once.Do(func() {
		node, err = snowflake.NewNode(nodeID)
		if err != nil {
			// 使用 fmt.Errorf 来包装错误，提供更多上下文
			err = fmt.Errorf("failed to create snowflake node: %w", err)
		}
	})
	return err
}

// GenID 生成一个新的雪花 ID
// 返回 int64 类型的 ID
func GenID() int64 {
	if node == nil {
		// 如果节点未初始化，这是一个严重错误，应该直接 panic
		// 这样可以及早发现配置问题
		panic("snowflake node not initialized")
	}
	return node.Generate().Int64()
}

