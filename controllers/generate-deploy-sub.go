package controllers

import (
	"bytes"
	"fmt"
	devopsAppsV1Beta1 "gitlab.myshuju.top/heshiying/devops/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"text/template"
)

func parseTemplate(tmplName string, md *devopsAppsV1Beta1.Deploy) []byte {
	tmpl, err := template.ParseFiles(fmt.Sprintf("controllers/templates/%s", tmplName))
	//tmpl, err := template.ParseFiles(fmt.Sprintf("templates/%s", tmplName))
	if err != nil {
		panic(err)
	}
	buffer := new(bytes.Buffer)
	err = tmpl.Execute(buffer, md)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}
func NewDeployment(dp *devopsAppsV1Beta1.Deploy) (*appsv1.Deployment, error) {
	deployment := appsv1.Deployment{}
	err := yaml.Unmarshal(parseTemplate("deployment.tmpl", dp), &deployment)
	if err != nil {
		return &deployment, err
	}
	return &deployment, nil
}

func NewService(dp *devopsAppsV1Beta1.Deploy) (*corev1.Service, error) {
	service := corev1.Service{}
	err := yaml.Unmarshal(parseTemplate("service.tmpl", dp), &service)
	if err != nil {
		return &service, err
	}
	return &service, nil
}
func NewIngress(dp *devopsAppsV1Beta1.Deploy) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	err := yaml.Unmarshal(parseTemplate("ingress.tmpl", dp), &ingress)
	if err != nil {
		return &ingress, err
	}
	return &ingress, nil
}
