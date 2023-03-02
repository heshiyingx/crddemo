package controllers

import (
	"context"
	"github.com/go-logr/logr"
	devopsv1beta1 "gitlab.myshuju.top/heshiying/devops/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *DeployReconciler) updateStatus(ctx context.Context, req ctrl.Request, deployment *appsv1.Deployment) error {
	deploy := devopsv1beta1.Deploy{}
	if err := r.Client.Get(ctx, req.NamespacedName, &deploy); err != nil {
		return client.IgnoreNotFound(err)
	}
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		logr.Logger{}.Error(err, "Error retrieving Deployment labels")
		return err
	}

	deploy.Status.Selector = selector.String()

	if deployment.Spec.Replicas != nil {
		deploy.Status.Replicas = deployment.Status.ReadyReplicas
	}
	return r.Client.Status().Update(ctx, &deploy)
}
