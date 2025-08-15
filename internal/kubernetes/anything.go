package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// -------------------- Workload: Deployment --------------------

func (op *DefaultK8sOperator) GetDeployment(ctx context.Context, cluster, namespace, name string) (*appsv1.Deployment, error) {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return nil, err
	}
	return cs.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (op *DefaultK8sOperator) ListDeployments(ctx context.Context, cluster, namespace, labelSelector string) ([]appsv1.Deployment, error) {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return nil, err
	}
	list, err := cs.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (op *DefaultK8sOperator) DeleteDeployment(ctx context.Context, cluster, namespace, name string, graceSeconds int64) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	return cs.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{
		GracePeriodSeconds: &graceSeconds,
	})
}

func (op *DefaultK8sOperator) ScaleDeployment(ctx context.Context, cluster, namespace, name string, replicas int32) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	return retryOnConflict(func() error {
		dep, err := cs.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		dep.Spec.Replicas = &replicas
		_, err = cs.AppsV1().Deployments(namespace).Update(ctx, dep, metav1.UpdateOptions{})
		return err
	})
}

func (op *DefaultK8sOperator) RolloutRestartDeployment(ctx context.Context, cluster, namespace, name string) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	return retryOnConflict(func() error {
		dep, err := cs.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if dep.Spec.Template.Annotations == nil {
			dep.Spec.Template.Annotations = map[string]string{}
		}
		dep.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
		_, err = cs.AppsV1().Deployments(namespace).Update(ctx, dep, metav1.UpdateOptions{})
		return err
	})
}

// -------------------- Workload: StatefulSet --------------------

func (op *DefaultK8sOperator) GetStatefulSet(ctx context.Context, cluster, namespace, name string) (*appsv1.StatefulSet, error) {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return nil, err
	}
	return cs.AppsV1().StatefulSets(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (op *DefaultK8sOperator) ScaleStatefulSet(ctx context.Context, cluster, namespace, name string, replicas int32) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	return retryOnConflict(func() error {
		sts, err := cs.AppsV1().StatefulSets(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		sts.Spec.Replicas = &replicas
		_, err = cs.AppsV1().StatefulSets(namespace).Update(ctx, sts, metav1.UpdateOptions{})
		return err
	})
}

func (op *DefaultK8sOperator) RolloutRestartStatefulSet(ctx context.Context, cluster, namespace, name string) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	return retryOnConflict(func() error {
		sts, err := cs.AppsV1().StatefulSets(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if sts.Spec.Template.Annotations == nil {
			sts.Spec.Template.Annotations = map[string]string{}
		}
		sts.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
		_, err = cs.AppsV1().StatefulSets(namespace).Update(ctx, sts, metav1.UpdateOptions{})
		return err
	})
}

// -------------------- Pod: 列表/删除/日志/Exec --------------------

func (op *DefaultK8sOperator) ListPods(ctx context.Context, cluster, namespace, labelSelector string) ([]corev1.Pod, error) {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return nil, err
	}
	list, err := cs.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (op *DefaultK8sOperator) ListPodsByDeployment(ctx context.Context, cluster, namespace, name string) (*corev1.PodList, error) {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return nil, err
	}

	// 首先获取 Deployment
	deployment, err := cs.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// 获取 Deployment 的标签选择器
	labelSelector := metav1.FormatLabelSelector(deployment.Spec.Selector)

	// 使用标签选择器列出所有相关的 Pods
	return cs.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
}

func (op *DefaultK8sOperator) DeletePod(ctx context.Context, cluster, namespace, name string, graceSeconds int64) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	return cs.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{
		GracePeriodSeconds: &graceSeconds,
	})
}

func (op *DefaultK8sOperator) GetPodLogs(ctx context.Context, cluster, namespace, pod, container string, tailLines, sinceSeconds int64) (string, error) {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return "", err
	}
	opts := &corev1.PodLogOptions{
		Container:    container,
		TailLines:    ptrInt64(tailLines),
		SinceSeconds: ptrInt64(sinceSeconds),
	}
	req := cs.CoreV1().Pods(namespace).GetLogs(pod, opts)
	stream, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer func() { _ = stream.Close() }()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, stream); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (op *DefaultK8sOperator) ExecPod(ctx context.Context, cluster, namespace, pod, container string, cmd []string, tty bool) (string, error) {
	cs, cfg, err := op.Client(cluster)
	if err != nil {
		return "", err
	}
	req := cs.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   cmd,
			Stdout:    true,
			Stderr:    true,
			TTY:       tty,
		}, k8sscheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		return "", err
	}
	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    tty,
	})
	out := stdout.String()
	if s := stderr.String(); s != "" {
		out = out + "\n" + s
	}
	return out, err
}

// -------------------- Node: Cordon / Uncordon / Drain --------------------

func (op *DefaultK8sOperator) CordonNode(ctx context.Context, cluster, node string) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	return retryOnConflict(func() error {
		n, err := cs.CoreV1().Nodes().Get(ctx, node, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if n.Spec.Unschedulable {
			return nil
		}
		n.Spec.Unschedulable = true
		_, err = cs.CoreV1().Nodes().Update(ctx, n, metav1.UpdateOptions{})
		return err
	})
}

func (op *DefaultK8sOperator) UncordonNode(ctx context.Context, cluster, node string) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	return retryOnConflict(func() error {
		n, err := cs.CoreV1().Nodes().Get(ctx, node, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if !n.Spec.Unschedulable {
			return nil
		}
		n.Spec.Unschedulable = false
		_, err = cs.CoreV1().Nodes().Update(ctx, n, metav1.UpdateOptions{})
		return err
	})
}

func (op *DefaultK8sOperator) DrainNode(ctx context.Context, cluster, node string, force, ignoreDaemonSets, deleteEmptyDirData bool) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	// 1. cordon
	if err := op.CordonNode(ctx, cluster, node); err != nil {
		return err
	}
	// 2. 列出该节点所有 pod
	pods, err := cs.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + node,
	})
	if err != nil {
		return err
	}
	// 3. 逐个 evict（跳过 DS/mirror，保守处理本地存储）
	for _, p := range pods.Items {
		if p.Namespace == "kube-system" && isMirrorPod(&p) {
			continue
		}
		if !force && hasLocalStorage(&p) && !deleteEmptyDirData {
			continue
		}
		if !ignoreDaemonSets && controlledBy(&p, "DaemonSet") {
			continue
		}
		ev := &policyv1.Eviction{
			ObjectMeta: metav1.ObjectMeta{Name: p.Name, Namespace: p.Namespace},
			DeleteOptions: &metav1.DeleteOptions{
				GracePeriodSeconds: ptrInt64(30),
			},
		}
		err := cs.PolicyV1().Evictions(p.Namespace).Evict(ctx, ev)
		if apierrors.IsTooManyRequests(err) {
			time.Sleep(2 * time.Second)
			err = cs.PolicyV1().Evictions(p.Namespace).Evict(ctx, ev)
		}
		if err != nil && !apierrors.IsNotFound(err) {
			return fmt.Errorf("evict pod %s/%s: %w", p.Namespace, p.Name, err)
		}
	}
	return nil
}

// -------------------- ConfigMap / Secret --------------------

func (op *DefaultK8sOperator) GetConfigMap(ctx context.Context, cluster, namespace, name string) (*corev1.ConfigMap, error) {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return nil, err
	}
	return cs.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (op *DefaultK8sOperator) CreateOrUpdateConfigMap(ctx context.Context, cluster string, cm *corev1.ConfigMap) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	_, err = cs.CoreV1().ConfigMaps(cm.Namespace).Create(ctx, cm, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		_, err = cs.CoreV1().ConfigMaps(cm.Namespace).Update(ctx, cm, metav1.UpdateOptions{})
	}
	return err
}

func (op *DefaultK8sOperator) GetSecret(ctx context.Context, cluster, namespace, name string) (*corev1.Secret, error) {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return nil, err
	}
	return cs.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (op *DefaultK8sOperator) CreateOrUpdateSecret(ctx context.Context, cluster string, s *corev1.Secret) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	_, err = cs.CoreV1().Secrets(s.Namespace).Create(ctx, s, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		_, err = cs.CoreV1().Secrets(s.Namespace).Update(ctx, s, metav1.UpdateOptions{})
	}
	return err
}

// -------------------- Namespace / Service --------------------

func (op *DefaultK8sOperator) EnsureNamespace(ctx context.Context, cluster, name string, labels map[string]string) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name, Labels: labels}}
	_, err = cs.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		_, err = cs.CoreV1().Namespaces().Update(ctx, ns, metav1.UpdateOptions{})
	}
	return err
}

func (op *DefaultK8sOperator) DeleteNamespace(ctx context.Context, cluster, name string) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	return cs.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
}

func (op *DefaultK8sOperator) EnsureServiceClusterIP(ctx context.Context, cluster string, svc *corev1.Service) error {
	cs, _, err := op.Client(cluster)
	if err != nil {
		return err
	}
	if svc.Spec.Type == "" {
		svc.Spec.Type = corev1.ServiceTypeClusterIP
	}
	_, err = cs.CoreV1().Services(svc.Namespace).Create(ctx, svc, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		exist, getErr := cs.CoreV1().Services(svc.Namespace).Get(ctx, svc.Name, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}
		svc.Spec.ClusterIP = exist.Spec.ClusterIP
		_, err = cs.CoreV1().Services(svc.Namespace).Update(ctx, svc, metav1.UpdateOptions{})
	}
	return err
}

// -------------------- Helpers --------------------

func retryOnConflict(fn func() error) error {
	for i := 0; i < 5; i++ {
		if err := fn(); err != nil {
			if apierrors.IsConflict(err) {
				time.Sleep(time.Duration(200*(i+1)) * time.Millisecond)
				continue
			}
			return err
		}
		return nil
	}
	return fmt.Errorf("conflict retries exceeded")
}

func ptrInt64(v int64) *int64 { return &v }

func isMirrorPod(p *corev1.Pod) bool {
	_, ok := p.Annotations["kubernetes.io/config.mirror"]
	return ok
}

func hasLocalStorage(p *corev1.Pod) bool {
	for _, v := range p.Spec.Volumes {
		if v.EmptyDir != nil {
			return true
		}
	}
	return false
}

func controlledBy(p *corev1.Pod, kind string) bool {
	for _, o := range p.OwnerReferences {
		if o.Controller != nil && *o.Controller && strings.EqualFold(o.Kind, kind) {
			return true
		}
	}
	return false
}
