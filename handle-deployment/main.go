package main

import (
	"context"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal("Get config file err", err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Get configSet err", err)
	}

	// 获取deployment clientset
	dpClient := clientSet.AppsV1().Deployments(corev1.NamespaceDefault)
	log.Println("Create deploymentClient...")

	// 调用创建deployment的方法
	log.Println("Create deployment...")
	if err := createDeployment(dpClient); err != nil {
		log.Fatal(err)
	}
	<-time.Tick(120 * time.Second)

	// 调用更新deploymen的方法
	log.Println("Update deployment...")
	if err := updateDeployment(dpClient); err != nil {
		log.Fatal(err)
	}
	<-time.Tick(120 * time.Second)

	// 调用删除deployment的方法
	log.Println("Delete deployment...")
	if err := deleteDeployment(dpClient); err != nil {
		log.Fatal(err)
	}
	<-time.Tick(120 * time.Second)
}

// 创建deployment的方法
func createDeployment(dpClient v1.DeploymentInterface) error {
	replicas := int32(3)
	newDP := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deploy",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.17",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := dpClient.Create(context.TODO(), newDP, metav1.CreateOptions{})
	return err
}

// 更新deployment的方法
func updateDeployment(dpClient v1.DeploymentInterface) error {
	dp, err := dpClient.Get(context.TODO(), "nginx-deploy", metav1.GetOptions{})
	if err != nil {
		return err
	}
	dp.Spec.Template.Spec.Containers[0].Image = "nginx:1.18"

	return retry.RetryOnConflict(
		retry.DefaultRetry, func() error {
			_, err = dpClient.Update(context.TODO(), dp, metav1.UpdateOptions{})
			return err
		},
	)
}

//删除deployment的方法
func deleteDeployment(dpClient v1.DeploymentInterface) error {
	deletePolicy := metav1.DeletePropagationForeground
	return dpClient.Delete(
		context.TODO(), "nginx-deploy", metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		},
	)
}
