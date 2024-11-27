/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/pflag"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"

	"github.com/sap/clustersecret-operator/internal/controller"

	coreclients "github.com/sap/clustersecret-operator/pkg/client/clientset/versioned"
)

var (
	kubeconfig     string
	leaseNamespace string
	leaseName      string
	leaseId        string
)

func main() {
	// initialize stderr logger
	errlog := log.New(os.Stderr, "", 0)

	// parse flags
	pflag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required/allowed if running out-of-cluster")
	pflag.StringVar(&leaseNamespace, "lease_namespace", "", "Lease namespace. Required if running out-of-cluster; otherwise defaults to controller's namespace")
	pflag.StringVar(&leaseName, "lease_name", "", "Lease name. Required")
	pflag.StringVar(&leaseId, "lease_id", "", "Lease ID. Optional; if unspecified, a unique ID will be generated")

	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.CommandLine.SortFlags = false
	pflag.Parse()

	// check if running in-cluster or out-of-cluster
	inCluster, namespace, err := checkIfRunningInCluster()
	if err != nil {
		klog.Fatalf("error checking whether running in-cluster or out-of-cluster: %s", err)
	}

	// use fallback from environment for certain flags
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}

	// check/default flags
	if inCluster && kubeconfig != "" {
		errlog.Fatal("flag --kubeconfig not allowed when running in-cluster")
	}
	if leaseNamespace == "" {
		leaseNamespace = namespace
	}
	if leaseNamespace == "" {
		errlog.Fatal("flag --lease_namespace empty or not provided; required if running out-of-cluster")
	}
	if leaseName == "" {
		errlog.Fatal("flag --lease_name empty or not provided")
	}
	if leaseId == "" {
		leaseId = uuid.New().String()
	}

	// setup api clients
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		errlog.Fatalf("error building kubeconfig: %s", err)
	}

	kubeclient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("error building kubernetes client: %s", err)
	}

	coreclient, err := coreclients.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("error building core client: %s", err)
	}

	// create main context
	ctx, cancel := context.WithCancel(context.Background())

	// create controller
	controller := controller.NewController(ctx, kubeclient, coreclient, nil)

	// trying to become leader
	leaderelection.RunOrDie(
		ctx,
		leaderelection.LeaderElectionConfig{
			Lock: &resourcelock.LeaseLock{
				LeaseMeta: metav1.ObjectMeta{
					Name:      leaseName,
					Namespace: leaseNamespace,
				},
				Client: kubeclient.CoordinationV1(),
				LockConfig: resourcelock.ResourceLockConfig{
					Identity: leaseId,
				},
			},
			ReleaseOnCancel: false,
			LeaseDuration:   15 * time.Second,
			RenewDeadline:   10 * time.Second,
			RetryPeriod:     2 * time.Second,
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: func(ctx context.Context) {
					klog.Infof("successfully acquired leadership (my id: %s); starting controller", leaseId)
					controller.Start()
				},
				OnStoppedLeading: func() {
					klog.Infof("stopped leading (my id: %s)", leaseId)
					cancel()
					controller.Wait()
				},
				OnNewLeader: func(identity string) {
					if identity == leaseId {
						return
					}
					klog.Infof("observed new leader (my id: %s, new leader id: %s); waiting to become leader", leaseId, identity)
				},
			},
		},
	)

	// exit
	klog.Info("exiting")
}

func checkIfRunningInCluster() (bool, string, error) {
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		// running in-cluster
		if raw, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
			return true, string(raw), nil
		} else {
			return false, "", err
		}
	} else if os.IsNotExist(err) {
		// running out-of-cluster
		return false, "", nil
	} else {
		return false, "", err
	}
}
