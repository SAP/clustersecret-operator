/*
Copyright (c) 2023 SAP SE

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

package controller

import (
	"context"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes" // could also be aliased 'kubeclients' but we keep it as 'kubernetes' since most people do
	kubescheme "k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	kubecorev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	corev1alpha1 "github.com/sap/clustersecret-operator/pkg/apis/core.cs.sap.com/v1alpha1"
	coreclients "github.com/sap/clustersecret-operator/pkg/client/clientset/versioned"
	corescheme "github.com/sap/clustersecret-operator/pkg/client/clientset/versioned/scheme"
	coreinformers "github.com/sap/clustersecret-operator/pkg/client/informers/externalversions"
	corev1alpha1listers "github.com/sap/clustersecret-operator/pkg/client/listers/core.cs.sap.com/v1alpha1"
)

const (
	ControllerName = "clustersecret-operator.cs.sap.com"
)

type Controller struct {
	ctx                   context.Context                         // controller context; controller will terminate when context is cancelled
	kubeclient            kubernetes.Interface                    // kubernetes client; use client interface, so we can mock it (e.g. with the fake client)
	coreclient            coreclients.Interface                   // core client; use client interface, so we can mock it (e.g. with the fake client)
	kubeinformerFactory   kubeinformers.SharedInformerFactory     // kubernetes informer factory
	coreinformerFactory   coreinformers.SharedInformerFactory     // core informer factory
	namespaceInformer     cache.SharedIndexInformer               // namespace informer
	secretInformer        cache.SharedIndexInformer               // secret informer
	clusterSecretInformer cache.SharedIndexInformer               // clustersecret informer
	namespaceLister       kubecorev1listers.NamespaceLister       // namespace lister
	secretLister          kubecorev1listers.SecretLister          // secret lister
	clusterSecretLister   corev1alpha1listers.ClusterSecretLister // clustersecret lister
	eventRecorder         record.EventRecorder                    // event recorder
	workqueue             workqueue.RateLimitingInterface         // workqueue
	numWorkers            int                                     // number of worker routines
	wgWorkers             sync.WaitGroup                          // wait group to be able to work for workers to complete
	synchronizer          Synchronizer                            // cache synchronizer
}

type workqueueItem struct {
	key  int
	name string
}

const (
	workqueueItemKeyNamespace = iota
	workqueueItemKeyClusterSecret
)

func NewController(ctx context.Context, kubeclient kubernetes.Interface, coreclient coreclients.Interface, synchronizer Synchronizer) *Controller {
	// kubernetes client (for namespaces, secrets)
	kubeinformerFactory := kubeinformers.NewSharedInformerFactory(kubeclient, 300*time.Second)
	nsInformer := kubeinformerFactory.Core().V1().Namespaces()
	scInformer := kubeinformerFactory.Core().V1().Secrets()
	// attention: important to create informer and lister before starting the factory !!!
	namespaceInformer := nsInformer.Informer()
	namespaceLister := nsInformer.Lister()
	secretInformer := scInformer.Informer()
	secretLister := scInformer.Lister()

	// core client (for our custom resources, i.e. for clustersecrets)
	coreinformerFactory := coreinformers.NewSharedInformerFactory(coreclient, 300*time.Second)
	csInformer := coreinformerFactory.Core().V1alpha1().ClusterSecrets()
	// attention: important to create informer and lister before starting the factory !!!
	clusterSecretInformer := csInformer.Informer()
	clusterSecretLister := csInformer.Lister()

	// setup event recorder
	scheme := runtime.NewScheme()
	kubescheme.AddToScheme(scheme)
	corescheme.AddToScheme(scheme)
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.V(3).Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclient.CoreV1().Events("")})
	eventRecorder := eventBroadcaster.NewRecorder(scheme, corev1.EventSource{Component: ControllerName})

	// setup workqueue
	workqueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "")
	go func() {
		<-ctx.Done()
		klog.V(1).Info("shutting down work queue")
		workqueue.ShutDown()
	}()

	// init synchronizer
	if synchronizer != nil {
		informers := make(map[schema.GroupVersionKind]cache.SharedIndexInformer)
		// todo: use constants
		informers[schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}] = namespaceInformer
		informers[schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}] = secretInformer
		informers[corev1alpha1.ClusterSecretGroupVersionKind] = clusterSecretInformer
		synchronizer.Init(informers)
	}

	return &Controller{
		ctx:                   ctx,
		kubeclient:            kubeclient,
		coreclient:            coreclient,
		kubeinformerFactory:   kubeinformerFactory,
		coreinformerFactory:   coreinformerFactory,
		namespaceInformer:     namespaceInformer,
		secretInformer:        secretInformer,
		clusterSecretInformer: clusterSecretInformer,
		namespaceLister:       namespaceLister,
		secretLister:          secretLister,
		clusterSecretLister:   clusterSecretLister,
		eventRecorder:         eventRecorder,
		workqueue:             workqueue,
		numWorkers:            3, // todo: make configurable
		synchronizer:          synchronizer,
	}
}

// this method should not be called more than once on the same receiver; todo: safeguard with some lock
func (c *Controller) Start() {
	c.startEventHandlers()
	c.startWorkers()
	c.startInformers()
}

func (c *Controller) Wait() {
	<-c.ctx.Done()
	c.wgWorkers.Wait()
}

func (c *Controller) startInformers() {
	c.kubeinformerFactory.Start(c.ctx.Done())
	for _, ok := range c.kubeinformerFactory.WaitForCacheSync(c.ctx.Done()) {
		if !ok {
			klog.Fatal("error waiting for informer caches to sync")
		}
	}
	c.coreinformerFactory.Start(c.ctx.Done())
	for _, ok := range c.coreinformerFactory.WaitForCacheSync(c.ctx.Done()) {
		if !ok {
			klog.Fatal("error waiting for informer caches to sync")
		}
	}
}

func (c *Controller) startEventHandlers() {
	c.namespaceInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(new interface{}) {
				c.enqueueNamespace("ADD", new)
			},
			UpdateFunc: func(old, new interface{}) {
				c.enqueueNamespace("UPDATE", new)
			},
			// ignore deletions (managed secrets in the deleted namespace will anyway be deleted)
			// DeleteFunc: c.enqueueNamespace,
		},
	)
	c.clusterSecretInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(new interface{}) {
				c.enqueueClusterSecret("ADD", new)
			},
			UpdateFunc: func(old, new interface{}) {
				oldClusterSecret, ok := old.(*corev1alpha1.ClusterSecret)
				if !ok {
					panic("this cannot happen")
				}
				newClusterSecret, ok := new.(*corev1alpha1.ClusterSecret)
				if !ok {
					panic("this cannot happen")
				}
				if oldClusterSecret.Generation != newClusterSecret.Generation {
					c.enqueueClusterSecret("UPDATE", new)
				}
			},
			DeleteFunc: func(old interface{}) {
				c.enqueueClusterSecret("DELETE", old)
			},
		},
	)
}

func (c *Controller) startWorkers() {
	// spawn worker routines
	for i := 0; i < c.numWorkers; i++ {
		c.wgWorkers.Add(1)
		go func(i int) {
			defer c.wgWorkers.Done()
			klog.V(1).Infof("worker %d starting", i)
			for {
				// get object from queue; will block as long as queue is empty
				// note: due to the implementation of the workqueue it is guaranteed that an item cannot be processed by more than one worker at the same time
				obj, shutdown := c.workqueue.Get()
				if shutdown {
					klog.V(1).Infof("worker %d exiting", i)
					return
				}
				// cast to workqueueItem (we know that there cannot be anything different in the queue)
				item, ok := obj.(workqueueItem)
				if !ok {
					panic("this cannot happen")
				}
				// process item; use an anonymous func to easily ensure calling Done() for the item
				func(item workqueueItem) {
					defer c.workqueue.Done(item)
					switch item.key {
					case workqueueItemKeyNamespace:
						if err := c.reconcileNamespace(item.name); err != nil {
							c.workqueue.AddRateLimited(item)
							klog.Errorf("error reconciling namespace %s: %s (requeuing)", item.name, err)
							return
						}
						c.workqueue.Forget(item)
						klog.V(2).Infof("successfully reconciled namespace %s", item.name)
					case workqueueItemKeyClusterSecret:
						if err := c.reconcileClusterSecret(item.name); err != nil {
							c.workqueue.AddRateLimited(item)
							klog.Errorf("error reconciling clustersecret %s: %s (requeuing)", item.name, err)
							return
						}
						c.workqueue.Forget(item)
						klog.V(2).Infof("successfully reconciled clustersecret %s", item.name)
					default:
						panic("this cannot happen")
					}
				}(item)
			}
		}(i)
	}
}
