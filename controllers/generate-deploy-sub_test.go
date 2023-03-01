package controllers

import (
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	devopsAppsV1Beta1 "gitlab.myshuju.top/heshiying/devops/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

func Test_parseTemplate(t *testing.T) {
	type args struct {
		tmplName string
		md       *devopsAppsV1Beta1.Deploy
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTemplate(tt.args.tmplName, tt.args.md); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDeployment(t *testing.T) {
	type args struct {
		dp *devopsAppsV1Beta1.Deploy
	}
	bytes, err := os.ReadFile(filepath.Join("testdata", "deploy-ingress.yaml"))
	if err != nil {
		return
	}
	db := devopsAppsV1Beta1.Deploy{}
	yaml.Unmarshal(bytes, &db)
	tests := []struct {
		name    string
		args    args
		want    *appsv1.Deployment
		wantErr bool
	}{
		{
			name: "hh",
			args: args{dp: &db},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDeployment(tt.args.dp)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	type args struct {
		dp *devopsAppsV1Beta1.Deploy
	}
	//bytes, err := os.ReadFile(filepath.Join("/Users/john/study/code/go/k8s_learn/heroku/controllers/testdata", "deploy-ingress.yaml"))
	bytes, err := os.ReadFile(filepath.Join("/Users/john/study/code/go/k8s_learn/heroku/controllers/testdata", "deploy-nodeport.yaml"))
	if err != nil {
		return
	}
	db := devopsAppsV1Beta1.Deploy{}
	yaml.Unmarshal(bytes, &db)
	tests := []struct {
		name    string
		args    args
		want    *corev1.Service
		wantErr bool
	}{
		{
			args: args{
				dp: &db,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewService(tt.args.dp)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIngress(t *testing.T) {
	type args struct {
		dp *devopsAppsV1Beta1.Deploy
	}
	tests := []struct {
		name    string
		args    args
		want    *networkingv1.Ingress
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewIngress(tt.args.dp)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIngress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIngress() = %v, want %v", got, tt.want)
			}
		})
	}
}
