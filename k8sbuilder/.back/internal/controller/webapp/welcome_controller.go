/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webapp

import (
	"context"

	webappv1 "github.com/futugyousuzu/k8sbuilder/api/webapp/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// WelcomeReconciler reconciles a Welcome object
type WelcomeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.vishel.io,resources=welcomes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.vishel.io,resources=welcomes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.vishel.io,resources=welcomes/finalizers,verbs=update

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=list;watch;get;patch;create;update
//+kubebuilder:rbac:groups=core,resources=services,verbs=list;watch;get;patch;create;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Welcome object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *WelcomeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.WithValues("welcome", req.NamespacedName)
	log.Info("reconciling welcome")

	var welcome webappv1.Welcome
	if err := r.Get(ctx, req.NamespacedName, &welcome); err != nil {
		log.Error(err, "unable to fetch CronJob")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deployment, err := r.createWelcomeDeployment(ctx, welcome)
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Info("create deployment success!")
	svc, err := r.createService(ctx, welcome)
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Info("create service success!")
	applyOpts := []client.PatchOption{client.ForceOwnership, client.
		FieldOwner("welcome_controller")}
	err = r.Patch(ctx, &deployment, client.Apply, applyOpts...)

	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.Patch(ctx, &svc, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WelcomeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Welcome{}).
		Complete(r)
}

func (r *WelcomeReconciler) createWelcomeDeployment(ctx context.Context, welcome webappv1.Welcome) (appsv1.Deployment, error) {
	defOne := int32(1)
	name := welcome.Spec.Name
	if name == "" {
		name = "world"
	}

	depl := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: appsv1.SchemeGroupVersion.String(), Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      welcome.Name,
			Namespace: welcome.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &defOne,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"welcome": welcome.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"welcome": welcome.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "welcome",
							Env: []corev1.EnvVar{
								{Name: "NAME", Value: name},
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 8080, Name: "http", Protocol: "TCP"},
							},
							Image: "sdfcdwefe/operatordemo:v1",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewMilliQuantity(100000, resource.BinarySI),
								},
							},
						},
					},
				},
			},
		},
	}
	if err := r.Create(ctx, &depl); err != nil {
		return depl, err
	}
	return depl, nil
}

func (r *WelcomeReconciler) createService(ctx context.Context, welcome webappv1.Welcome) (corev1.Service, error) {
	svc := corev1.Service{
		TypeMeta: metav1.TypeMeta{APIVersion: corev1.SchemeGroupVersion.String(), Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      welcome.Name,
			Namespace: welcome.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "http", Port: 8080, Protocol: "TCP", TargetPort: intstr.FromString("http")},
			},
			Selector: map[string]string{"welcome": welcome.Name},
			Type:     corev1.ServiceTypeLoadBalancer,
		},
	}
	if err := r.Create(ctx, &svc); err != nil {
		return svc, err
	}
	return svc, nil
}
