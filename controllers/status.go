package controllers

import (
	"context"
	"github.com/go-logr/logr"
	devopsv1beta1 "gitlab.myshuju.top/heshiying/devops/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *DeployReconciler) updateStatus(ctx context.Context, deploy *devopsv1beta1.Deploy, deployment *appsv1.Deployment) error {
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		logr.Logger{}.Error(err, "Error retrieving Deployment labels")
		return err
	}

	deploy.Status.Selector = selector.String()

	if deployment.Spec.Replicas != nil {
		deploy.Status.Replicas = deployment.Status.ReadyReplicas
	}
	return r.Client.Status().Update(ctx, deploy)
}
