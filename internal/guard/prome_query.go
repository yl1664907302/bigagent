package guard

import (
	"context"
	"fmt"
	papi "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/klog/v2"
	"time"
)

func (gr *Guarder) PromInstantQuery(ql string, apiAddr string, twSec int) (model.Vector, error) {

	//apiAddr := gr.Cg.K8sComputePromAddrMap[clusterName]
	//if apiAddr == "" {
	//	// 打印错误和返回都需要用
	//	msg := fmt.Sprintf("PromInstantQuery.getapiAddr.K8sComputePromAddrMap.empty[clusterName:%v][ql:%v]", clusterName, ql)
	//	err := fmt.Errorf(msg)
	//	klog.Errorf(msg)
	//	return nil, err
	//}

	// 使用prometheus sdk 查询
	// new pclient
	client, err := papi.NewClient(papi.Config{
		Address: apiAddr,
	})
	if err != nil {

		klog.Errorf("[queryOneMetric.Error.creating.client][err:%v]", err)
		return nil, err
	}
	// 用 client 构造查询api对象
	v1Api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(twSec)*time.Second)
	defer cancel()

	result, warnings, err := v1Api.Query(ctx, ql, time.Now())
	if err != nil {
		klog.Errorf("[queryOneMetric.querying.Prometheus.Error][url:%+v][ql][%v][err:%v]", apiAddr, ql, err)
		return nil, err
	}
	if len(warnings) > 0 {
		msg := fmt.Sprintf("[queryOneMetric.querying.Prometheus.warnings][url:%+v][ql:%v][warnings:%v]", apiAddr, ql, warnings)
		klog.Errorf(msg)
		return nil, fmt.Errorf(msg)
	}
	vec, ok := result.(model.Vector)
	if !ok {
		klog.Errorf("[queryOneMetric.querying.Prometheus.result.to.vector.error][result:%+v]", result)
		return nil, err
	}
	return vec, nil
}

func (gr *Guarder) TestInstantQuery() {
	//vecs, err := gr.PromInstantQuery("node_os_version", "http://192.168.0.100:8091", 2)
	vecs, err := gr.PromInstantQuery("node_os_version", "cpu-compute-03", 2)
	fmt.Println(vecs, err)
	for _, vec := range vecs {
		vec := vec
		fmt.Println(vec.Value)
		fmt.Println(vec.Metric)
		fmt.Println(vec.Timestamp)
	}

}
