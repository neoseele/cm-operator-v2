package custommetric

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	cmv1alpha1 "github.com/neoseele/cm-operator/pkg/apis/cm/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_custommetric")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

func int32Ptr(i int32) *int32 { return &i }

// Add creates a new CustomMetric Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCustomMetric{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("custommetric-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CustomMetric
	err = c.Watch(&source.Kind{Type: &cmv1alpha1.CustomMetric{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner CustomMetric
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cmv1alpha1.CustomMetric{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCustomMetric implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCustomMetric{}

// ReconcileCustomMetric reconciles a CustomMetric object
type ReconcileCustomMetric struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CustomMetric object and makes changes based on the state read
// and what is in the CustomMetric.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCustomMetric) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Name", request.Name)
	reqLogger.Info("Reconciling CustomMetric")

	// Fetch the CustomMetric instance
	instance := &cmv1alpha1.CustomMetric{}
	if err := r.client.Get(context.TODO(), request.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// ClusterRole
	if err := r.reconcileClusterRole(instance, reqLogger); err != nil {
		return reconcile.Result{}, err
	}

	// ClusterRoleBinding
	if err := r.reconcileClusterRoleBinding(instance, reqLogger); err != nil {
		return reconcile.Result{}, err
	}

	// ConfigMap
	if err := r.reconcileConfigMap(instance, reqLogger); err != nil {
		return reconcile.Result{}, err
	}

	// Deployment
	if err := r.reconcileDeployment(instance, reqLogger); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileCustomMetric) reconcileClusterRole(cr *cmv1alpha1.CustomMetric, reqLogger logr.Logger) error {
	// Define a new object
	clusterRole := newClusterRole(cr)

	// Set CustomMetric instance as the owner and controller
	// if err := controllerutil.SetControllerReference(cr, clusterRole, r.scheme); err != nil {
	// 	return err
	// }

	// Check if this object already exists
	found := &rbacv1.ClusterRole{}
	err := r.client.Get(context.TODO(), client.ObjectKey{Name: clusterRole.Name}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new clusterRole name", "clusterRole.Name", clusterRole.Name)
		err = r.client.Create(context.TODO(), clusterRole)
		if err != nil {
			return err
		}

		// object created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}

	// object already exists - don't requeue
	reqLogger.Info("Skip reconcile: ConfigMap already exists", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
	return nil
}

func (r *ReconcileCustomMetric) reconcileClusterRoleBinding(cr *cmv1alpha1.CustomMetric, reqLogger logr.Logger) error {
	// Define a new object
	clusterRoleBinding := newClusterRoleBinding(cr)

	// Set CustomMetric instance as the owner and controller
	// if err := controllerutil.SetControllerReference(cr, clusterRoleBinding, r.scheme); err != nil {
	// 	return err
	// }

	// Check if this object already exists
	found := &rbacv1.ClusterRoleBinding{}
	err := r.client.Get(context.TODO(), client.ObjectKey{Name: clusterRoleBinding.Name}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new clusterRoleBinding", "clusterRoleBinding.Name", clusterRoleBinding.Name)
		err = r.client.Create(context.TODO(), clusterRoleBinding)
		if err != nil {
			return err
		}

		// object created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}

	// object already exists - don't requeue
	reqLogger.Info("Skip reconcile: ConfigMap already exists", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
	return nil
}

func (r *ReconcileCustomMetric) reconcileConfigMap(cr *cmv1alpha1.CustomMetric, reqLogger logr.Logger) error {
	// Define a new object
	configmap := newConfigMap(cr)

	// Set CustomMetric instance as the owner and controller
	if err := controllerutil.SetControllerReference(cr, configmap, r.scheme); err != nil {
		return err
	}

	// Check if this object already exists
	found := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: configmap.Name, Namespace: configmap.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ConfigMap", "ConfigMap.Namespace", configmap.Namespace, "ConfigMap.Name", configmap.Name)
		err = r.client.Create(context.TODO(), configmap)
		if err != nil {
			return err
		}

		// object created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}

	// object already exists - don't requeue
	reqLogger.Info("Skip reconcile: ConfigMap already exists", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
	return nil
}

func (r *ReconcileCustomMetric) reconcileDeployment(cr *cmv1alpha1.CustomMetric, reqLogger logr.Logger) error {
	// Define a new object
	deployment := newDeployment(cr)

	// Set CustomMetric instance as the owner and controller
	if err := controllerutil.SetControllerReference(cr, deployment, r.scheme); err != nil {
		return err
	}

	// Check if this object already exists
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			return err
		}

		// object created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}

	// object already exists - don't requeue
	reqLogger.Info("Skip reconcile: Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	return nil
}

func newClusterRole(cr *cmv1alpha1.CustomMetric) *rbacv1.ClusterRole {
	resourceName := cr.Name + "-prometheus"

	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, cmv1alpha1.SchemeGroupVersion.WithKind("CustomMetric")),
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"nodes",
					"nodes/proxy",
					"services",
					"endpoints",
					"pods",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
			},
			{
				APIGroups: []string{
					"extentions",
				},
				Resources: []string{
					"ingresses",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
			},
			{
				NonResourceURLs: []string{
					"/metrics",
				},
				Verbs: []string{
					"get",
				},
			},
		},
	}
}

func newClusterRoleBinding(cr *cmv1alpha1.CustomMetric) *rbacv1.ClusterRoleBinding {
	resourceName := cr.Name + "-prometheus"

	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, cmv1alpha1.SchemeGroupVersion.WithKind("CustomMetric")),
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     resourceName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: cr.Namespace,
			},
		},
	}
}

func newConfigMap(cr *cmv1alpha1.CustomMetric) *corev1.ConfigMap {
	resourceName := cr.Name + "-prometheus"

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: cr.Namespace,
			// OwnerReferences: []metav1.OwnerReference{
			// 	*metav1.NewControllerRef(cr, cmv1alpha1.SchemeGroupVersion.WithKind("CustomMetric")),
			// },
		},
		Data: map[string]string{
			"prometheus.yml": `
scrape_configs:

- job_name: 'kubernetes-pods'

  kubernetes_sd_configs:
    - role: pod

  relabel_configs:
    - source_labels: [__meta_kubernetes_pod_annotation_cm_example_com_scrape]
      action: keep
      regex: true
    - source_labels:
      - __meta_kubernetes_pod_annotationpresent_cm_example_com_path
      - __meta_kubernetes_pod_annotation_cm_example_com_path
      action: replace
      target_label: __metrics_path__
      regex: true;(.+)
    - source_labels: [__address__, __meta_kubernetes_pod_annotation_cm_example_com_port]
      action: replace
      regex: ([^:]+)(?::\d+)?;(\d+)
      replacement: $1:$2
      target_label: __address__
    - action: labelmap
      regex: __meta_kubernetes_pod_label_(.+)
    - source_labels: [__meta_kubernetes_namespace]
      action: replace
      target_label: kubernetes_namespace
    - source_labels: [__meta_kubernetes_pod_name]
      action: replace
      target_label: kubernetes_pod_name

- job_name: 'kubernetes-nodes-cadvisor'

  metrics_path: /metrics/cadvisor

  kubernetes_sd_configs:
    - role: node

  relabel_configs:
    - source_labels: [__meta_kubernetes_node_annotation_cm_example_com_scrape]
      action: keep
      regex: true
    - source_labels: [__address__]
      action: replace
      regex: ([^:]+)(?::\d+)?
      replacement: $1:10255
      target_label: __address__
    - action: labelmap
      regex: __meta_kubernetes_node_label_(.+)
      `,
		},
	}
}

func newDeployment(cr *cmv1alpha1.CustomMetric) *appsv1.Deployment {
	resourceName := cr.Name + "-prometheus"

	labels := map[string]string{
		"app":        "prometheus-server",
		"controller": cr.Name,
	}

	sidecarArgs := []string{
		fmt.Sprintf("--stackdriver.project-id=%s", cr.Spec.Project),
		fmt.Sprintf("--stackdriver.kubernetes.cluster-name=%s", cr.Spec.Cluster),
		fmt.Sprintf("--stackdriver.kubernetes.location=%s", cr.Spec.Location),
		"--prometheus.wal-directory=/prometheus/wal",
		"--log.level=debug",
	}

	for _, m := range cr.Spec.Metrics {
		sidecarArgs = append(sidecarArgs, fmt.Sprintf("--include={__name__=~\"%s\"}", m))
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: cr.Namespace,
			// OwnerReferences: []metav1.OwnerReference{
			// 	*metav1.NewControllerRef(cr, cmv1alpha1.SchemeGroupVersion.WithKind("CustomMetric")),
			// },
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "prometheus",
							Image: "prom/prometheus:v2.6.1",
							Args: []string{
								"--config.file=/etc/prometheus/prometheus.yml",
								"--storage.tsdb.path=/prometheus/",
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9090,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "prometheus-config-volume",
									MountPath: "/etc/prometheus/",
								},
								{
									Name:      "prometheus-storage-volume",
									MountPath: "/prometheus/",
								},
							},
						},
						{
							Name:  "sidecar",
							Image: "gcr.io/stackdriver-prometheus/stackdriver-prometheus-sidecar:0.8.0",
							Args:  sidecarArgs,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9091,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "prometheus-storage-volume",
									MountPath: "/prometheus/",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "prometheus-config-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: resourceName,
									},
									DefaultMode: int32Ptr(420),
								},
							},
						},
						{
							Name: "prometheus-storage-volume", // default to emptyDir
						},
					},
				},
			},
		},
	}
}
