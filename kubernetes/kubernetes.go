package kubernetes

import (
	"context"
	"flag"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func GetClient() client.Client {
	flag.Parse()

	cfg := config.GetConfigOrDie()
	clt, err := client.New(cfg, client.Options{})
	if err != nil {
		panic(err.Error())
	}

	return clt
}

func GetPods(ctx context.Context, clt client.Client, namespace string) []string {
	fmt.Println("getpods")

	pod := new(corev1.Pod)

	var err error
	if namespace != "" {
		err = clt.Get(ctx, client.ObjectKey{Namespace: namespace}, pod)
	} else {
		err = clt.Get(ctx, client.ObjectKey{}, pod)
	}
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("# containers: %d\n", len(pod.Spec.Containers))

	pods, err := listPods(ctx, clt)
	if err != nil {
		panic(err.Error())
	}

	return pods
}

func listPods(ctx context.Context, clt client.Client) ([]string, error) {
	podList := new(corev1.PodList)
	selector, err := labels.Parse("k8s-app in (kube-dns, kube-dns-autoscaler)")
	if err != nil {
		return nil, err
	}
	err = clt.List(ctx, podList, client.MatchingLabelsSelector{Selector: selector})
	
	if err != nil {
		return nil, err
	}

	pods := make([]string, 0)
	for _, pod := range podList.Items {
		fmt.Printf("pod: %s.%s\n", pod.Name, pod.Namespace)
		pods = append(pods, pod.Name)
	}
	return pods, nil
}