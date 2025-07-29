package osquery

import (
	"context"
	"fmt"
	"time"

	"github.com/osquery/osquery-go"
)

var (
	OQry *Querier
)

// Querier 定义了一个用于查询 osquery 的结构体。
type Querier struct {
	client *osquery.ExtensionManagerClient
}

// NewQuerier 创建一个新的 Querier 实例。
// 它需要 osquery socket 的路径。
// 在 Windows 上，默认路径是 \\.\pipe\osquery.em
func NewQuerier(socketPath string, timeout time.Duration) (*Querier, error) {
	client, err := osquery.NewClient(socketPath, timeout)
	if err != nil {
		return nil, fmt.Errorf("创建 osquery 客户端失败: %w", err)
	}
	return &Querier{client: client}, nil
}

// Query 对 osquery 执行一个 SQL 查询。
func (q *Querier) Query(ctx context.Context, sql string) ([]map[string]string, error) {
	resp, err := q.client.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("执行查询失败: %w", err)
	}
	if resp.Status.Code != 0 {
		return nil, fmt.Errorf("查询失败，状态码 %d: %s", resp.Status.Code, resp.Status.Message)
	}
	return resp.Response, nil
}

// Close 关闭与 osquery 的连接。
func (q *Querier) Close() {
	q.client.Close()
}
