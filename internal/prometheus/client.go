package prom

import (
	"context"
	"fmt"
	"log"
	"time"

	papi "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// PromClient 是 Prometheus 客户端封装
type PromClient struct {
	Address string
	client  papi.Client
	v1api   v1.API
	timeout time.Duration
}

// NewPromClient 创建一个 Prometheus 客户端
func NewPromClient(addr string, timeoutSec int) (*PromClient, error) {
	client, err := papi.NewClient(papi.Config{
		Address: addr,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus client: %w", err)
	}

	return &PromClient{
		Address: addr,
		client:  client,
		v1api:   v1.NewAPI(client),
		timeout: time.Duration(timeoutSec) * time.Second,
	}, nil
}

// InstantQuery 执行即时查询（单点查询）
func (pc *PromClient) InstantQuery(ql string) (model.Vector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pc.timeout)
	defer cancel()

	log.Printf("[PromQL] running query: %s", ql)

	// 执行查询
	result, warnings, err := pc.v1api.Query(ctx, ql, time.Now())
	if err != nil {
		return nil, fmt.Errorf("prometheus query error: %w", err)
	}
	if len(warnings) > 0 {
		log.Printf("[PromQL] warnings: %v", warnings)
	}

	// 打印返回值类型和内容
	log.Printf("[PromQL] raw result type: %T", result)
	log.Printf("[PromQL] raw result: %#v", result)

	// 根据类型断言
	switch val := result.(type) {
	case model.Vector:
		if len(val) == 0 {
			log.Println("[PromQL] query result is empty (vector len=0)")
		}
		return val, nil
	case *model.Scalar:
		log.Printf("[PromQL] got scalar value: %v at %v", val.Value, val.Timestamp)
		// 把 scalar 转成 vector，方便后续处理
		vec := model.Vector{
			&model.Sample{
				Metric:    model.Metric{"__name__": "scalar"},
				Value:     val.Value,
				Timestamp: val.Timestamp,
			},
		}
		return vec, nil
	default:
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}
}

// RangeQuery 执行区间查询（时间序列范围）
func (pc *PromClient) RangeQuery(ql string, start, end time.Time, step time.Duration) (model.Matrix, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pc.timeout)
	defer cancel()

	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}
	result, warnings, err := pc.v1api.QueryRange(ctx, ql, r)
	if err != nil {
		return nil, fmt.Errorf("prometheus range query error: %w", err)
	}
	if len(warnings) > 0 {
		return nil, fmt.Errorf("prometheus range query warnings: %v", warnings)
	}

	matrix, ok := result.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}
	return matrix, nil
}
