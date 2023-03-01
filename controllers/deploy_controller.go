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

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	devopsAppsV1Beta1 "gitlab.myshuju.top/heshiying/devops/api/v1beta1"
)

// DeployReconciler reconciles a Deploy object
type DeployReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.myshuju.top,resources=deploys,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.myshuju.top,resources=deploys/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.myshuju.top,resources=deploys/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Deploy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *DeployReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx, "Deploy", req.NamespacedName)
	logger.Info("Deploy change")

	// TODO(user): your logic here
	// 1处理资源
	// 获取资源对象,
	deploy := devopsAppsV1Beta1.Deploy{}
	if err := r.Client.Get(ctx, req.NamespacedName, &deploy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// 创建一个copy，避免在后面操作对象，污染缓存
	deployCopy := deploy.DeepCopy()

	// 2.处理deployment，如果存在更新，如果不存在则创建
	// 获取deployment对象
	deployment := appsv1.Deployment{}
	if err := r.Client.Get(ctx, req.NamespacedName, &deployment); err != nil {
		if errors.IsNotFound(err) {
			err = r.createDeployment(ctx, req, deployCopy, logger)
			if err != nil {
				return ctrl.Result{}, err
			}
		} else {
			logger.Error(err, "get Deployment err")
			return ctrl.Result{}, err
		}

	} else {

		err = r.updateDeployment(ctx, req, deployCopy, deployment, logger)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	// 3.处理service
	// >>如果mode是ingress,如果不存在，则创建，如果存在更新普通service
	// >>如果mode是nodePort,如果不存在，则创建，如果存在更新nodePort
	service := corev1.Service{}
	if err := r.Client.Get(ctx, req.NamespacedName, &service); err != nil {
		if errors.IsNotFound(err) {
			if err := r.createService(ctx, req, deployCopy, logger); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			logger.Error(err, "get service err")
			return ctrl.Result{}, err
		}
	} else {
		if err := r.updateService(ctx, req, deployCopy, deployCopy.Spec.Expose.Mode, service, logger); err != nil {
			return ctrl.Result{}, err
		}
	}

	// 4.处理ingress
	// >>如果mode是ingress,如果不存在，创建ingress，如果存在更新ingress
	// >>如果mode是nodePort,如果存在，则删除。
	ingress := networkingv1.Ingress{}
	if err := r.Client.Get(ctx, req.NamespacedName, &ingress); err != nil {
		if errors.IsNotFound(err) {
			if err = r.ingressNotExistDeal(ctx, req.NamespacedName, deployCopy, deployCopy.Spec.Expose.Mode, logger); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			if err = r.ingressExistDeal(ctx, &ingress, deployCopy, deployCopy.Spec.Expose.Mode, logger); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeployReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devopsAppsV1Beta1.Deploy{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}

func (r *DeployReconciler) createDeployment(ctx context.Context, req ctrl.Request, deploy *devopsAppsV1Beta1.Deploy, logger logr.Logger) error {
	deployment, err := NewDeployment(deploy)
	if err != nil {
		return err
	}
	err = controllerutil.SetControllerReference(deploy, deployment, r.Scheme)
	if err != nil {
		return err
	}
	return r.Client.Create(ctx, deployment)

}

func (r *DeployReconciler) updateDeployment(ctx context.Context, req ctrl.Request, deploy *devopsAppsV1Beta1.Deploy, oldDeployment appsv1.Deployment, logger logr.Logger) error {
	deployment, err := NewDeployment(deploy)
	if err != nil {
		return err
	}
	err = r.Client.Update(ctx, deployment, client.DryRunAll)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(deployment.Spec, oldDeployment.Spec) {
		return nil
	}
	err = controllerutil.SetControllerReference(deploy, deployment, r.Scheme)
	if err != nil {
		return err
	}

	return r.Client.Update(ctx, deployment)
}

func (r *DeployReconciler) createService(ctx context.Context, req ctrl.Request, deploy *devopsAppsV1Beta1.Deploy, logger logr.Logger) error {
	service, err := NewService(deploy)
	if err != nil {
		return err
	}
	err = controllerutil.SetControllerReference(deploy, service, r.Scheme)
	if err != nil {
		return err
	}
	return r.Client.Create(ctx, service)
}

func (r *DeployReconciler) updateService(ctx context.Context, req ctrl.Request, deploy *devopsAppsV1Beta1.Deploy, mode devopsAppsV1Beta1.ExposeMode, oldService corev1.Service, logger logr.Logger) error {
	service, err := NewService(deploy)
	if err != nil {
		return err
	}
	err = r.Client.Update(ctx, service, client.DryRunAll)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(service.Spec, oldService.Spec) {
		return nil
	}
	err = controllerutil.SetControllerReference(deploy, service, r.Scheme)
	if err != nil {
		return err
	}
	return r.Client.Update(ctx, service)
}

func (r *DeployReconciler) ingressNotExistDeal(ctx context.Context, name types.NamespacedName, deploy *devopsAppsV1Beta1.Deploy, mode devopsAppsV1Beta1.ExposeMode, logger logr.Logger) error {
	if mode == devopsAppsV1Beta1.ExposeModeNodePort {
		return nil
	}
	if mode == devopsAppsV1Beta1.ExposeModeIngress {
		ingress, err := NewIngress(deploy)
		if err != nil {
			return err
		}
		err = controllerutil.SetControllerReference(deploy, ingress, r.Scheme)
		if err != nil {
			return err
		}
		return r.Client.Create(ctx, ingress)
	}
	return nil
}

func (r *DeployReconciler) ingressExistDeal(ctx context.Context, ingress *networkingv1.Ingress, deploy *devopsAppsV1Beta1.Deploy, mode devopsAppsV1Beta1.ExposeMode, logger logr.Logger) error {
	// 如果是nodePort就删除原有的ingress
	if mode == devopsAppsV1Beta1.ExposeModeNodePort {
		r.Client.DeleteAllOf(ctx, ingress)
		return nil
	}
	if mode == devopsAppsV1Beta1.ExposeModeIngress {
		newIngress, err := NewIngress(deploy)
		if err != nil {
			return err
		}
		err = r.Client.Update(ctx, newIngress, client.DryRunAll)
		if err != nil {
			return err
		}
		if reflect.DeepEqual(ingress.Spec, newIngress.Spec) {
			return nil
		}
		err = controllerutil.SetControllerReference(deploy, newIngress, r.Scheme)
		if err != nil {
			return err
		}
		return r.Client.Update(ctx, newIngress)
	}

	return nil
}
