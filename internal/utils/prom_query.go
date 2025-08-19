package utils

import (
	"context"
	"fmt"
	"time"

	papi "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// PromInstantQuery 执行一次 Prometheus 即时查询
func PromInstantQuery(ctx context.Context, ql string, apiAddr string, timeoutSeconds int) (model.Vector, error) {
	client, err := papi.NewClient(papi.Config{Address: apiAddr})
	if err != nil {
		return nil, err
	}
	api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	result, warnings, err := api.Query(ctx, ql, time.Now())
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		return nil, fmt.Errorf("prometheus warnings: %v", warnings)
	}
	vec, ok := result.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}
	return vec, nil
}
