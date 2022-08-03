package main

import (
	"context"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/linode/linodego"
)

type Environment struct {
	Debug               bool     `default:"false"`
	LinodeToken         string   `required:"true" split_words:"true"`
	InstanceIDs         []int    `required:"true" envconfig:"INSTANCE_IDS"`
	AdditionalAllowList []string `default:"192.168.128.0/32" split_words:"true"`
	Kubeconfig          string
}

func main() {
	var env Environment
	err := envconfig.Process("", &env)
	if err != nil {
		fmt.Printf("failed to get env vars: %s\n", err.Error())
		return
	}
	if env.Kubeconfig != "" {
		var err error
		env.Kubeconfig, err = cleanKubeconfigPath(env.Kubeconfig)
		if err != nil {
			fmt.Printf("failed to get env vars: %s\n", err.Error())
			return
		}
	}
	if env.Debug {
		fmt.Printf("loaded environment: %v\n", env)
	}

	linode := linodeClient(env.LinodeToken, env.Debug)
	k8s, err := kubernetesClient(env.Kubeconfig)
	if err != nil {
		fmt.Printf("failed to create kubernetes client: %s\n", err.Error())
		return
	}

	externalIPs, err := listExternalIPs(k8s)
	if err != nil {
		fmt.Printf("failed to list ExternalIPs: %s\n", err.Error())
		return
	}
	targetAllowList := append(env.AdditionalAllowList, externalIPs...)
	if len(targetAllowList) == 0 {
		fmt.Print("No ExternalIPs found and no additionalAllowList entries set. Got nothing to do, exiting.")
		return
	}
	fmt.Printf("targetAllowList: %v\n", targetAllowList)

	for _, id := range env.InstanceIDs {
		mdb, err := linode.GetMongoDatabase(context.Background(), id)
		if err != nil {
			fmt.Printf("failed to fetch MongoDB (#%d): %s", id, err.Error())
			continue
		}
		fmt.Printf("current MongoDB #%d allowList: %v\n", id, mdb.AllowList)

		if !equalAllowLists(mdb.AllowList, targetAllowList) {
			fmt.Printf("updating allowList of #%d to %v\n", id, targetAllowList)
			_, err = linode.UpdateMongoDatabase(context.Background(), id, linodego.MongoUpdateOptions{AllowList: &targetAllowList})
			if err != nil {
				fmt.Printf("failed to update MongoDB #%d allowList: %s", id, err.Error())
			}
		} else {
			fmt.Printf("allowList of MongoDB #%d is looking good, nothing to do\n", id)
		}
	}

	fmt.Println("done")
}
