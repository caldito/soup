package deployment

import (
	// from this repo
	"github.com/caldito/soup/pkg/k8s"

	// from other repos
	"context"
	"fmt"
	"k8s.io/client-go/rest"
)

func Deploy(namespace string, manifests []string, cloneLocation string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("Error getting cluster config")
		panic(err)
	}
	ctx := context.TODO()
	err = k8s.DeclareNamespaceSSA(ctx, config, namespace)
	if err != nil {
		fmt.Println("Error preparing namespace " + namespace)
		panic(err)
	}
	for _, manifest := range manifests {
		err = k8s.DoSSA(ctx, config, namespace, cloneLocation+"/"+manifest)
		if err != nil {
			fmt.Println("Error deploying the manifest " + manifest)
			panic(err)
		}
	}
	return nil
}