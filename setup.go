/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by Kubeform. DO NOT EDIT.

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gobuffalo/flect"
	vultr "github.com/vultr/terraform-provider-vultr/vultr"
	auditlib "go.bytebuilders.dev/audit/lib"
	arv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	informers "k8s.io/apiextensions-apiserver/pkg/client/informers/externalversions"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	admissionregistrationv1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	barev1alpha1 "kubeform.dev/provider-vultr-api/apis/bare/v1alpha1"
	blockv1alpha1 "kubeform.dev/provider-vultr-api/apis/block/v1alpha1"
	dnsv1alpha1 "kubeform.dev/provider-vultr-api/apis/dns/v1alpha1"
	firewallv1alpha1 "kubeform.dev/provider-vultr-api/apis/firewall/v1alpha1"
	instancev1alpha1 "kubeform.dev/provider-vultr-api/apis/instance/v1alpha1"
	isov1alpha1 "kubeform.dev/provider-vultr-api/apis/iso/v1alpha1"
	kubernetesv1alpha1 "kubeform.dev/provider-vultr-api/apis/kubernetes/v1alpha1"
	loadv1alpha1 "kubeform.dev/provider-vultr-api/apis/load/v1alpha1"
	objectv1alpha1 "kubeform.dev/provider-vultr-api/apis/object/v1alpha1"
	privatev1alpha1 "kubeform.dev/provider-vultr-api/apis/private/v1alpha1"
	reservedv1alpha1 "kubeform.dev/provider-vultr-api/apis/reserved/v1alpha1"
	reversev1alpha1 "kubeform.dev/provider-vultr-api/apis/reverse/v1alpha1"
	snapshotv1alpha1 "kubeform.dev/provider-vultr-api/apis/snapshot/v1alpha1"
	sshv1alpha1 "kubeform.dev/provider-vultr-api/apis/ssh/v1alpha1"
	startupv1alpha1 "kubeform.dev/provider-vultr-api/apis/startup/v1alpha1"
	userv1alpha1 "kubeform.dev/provider-vultr-api/apis/user/v1alpha1"
	controllersbare "kubeform.dev/provider-vultr-controller/controllers/bare"
	controllersblock "kubeform.dev/provider-vultr-controller/controllers/block"
	controllersdns "kubeform.dev/provider-vultr-controller/controllers/dns"
	controllersfirewall "kubeform.dev/provider-vultr-controller/controllers/firewall"
	controllersinstance "kubeform.dev/provider-vultr-controller/controllers/instance"
	controllersiso "kubeform.dev/provider-vultr-controller/controllers/iso"
	controllerskubernetes "kubeform.dev/provider-vultr-controller/controllers/kubernetes"
	controllersload "kubeform.dev/provider-vultr-controller/controllers/load"
	controllersobject "kubeform.dev/provider-vultr-controller/controllers/object"
	controllersprivate "kubeform.dev/provider-vultr-controller/controllers/private"
	controllersreserved "kubeform.dev/provider-vultr-controller/controllers/reserved"
	controllersreverse "kubeform.dev/provider-vultr-controller/controllers/reverse"
	controllerssnapshot "kubeform.dev/provider-vultr-controller/controllers/snapshot"
	controllersssh "kubeform.dev/provider-vultr-controller/controllers/ssh"
	controllersstartup "kubeform.dev/provider-vultr-controller/controllers/startup"
	controllersuser "kubeform.dev/provider-vultr-controller/controllers/user"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var _provider = vultr.Provider()

var runningControllers = struct {
	sync.RWMutex
	mp map[schema.GroupVersionKind]bool
}{mp: make(map[schema.GroupVersionKind]bool)}

func watchCRD(ctx context.Context, crdClient *clientset.Clientset, vwcClient *admissionregistrationv1.AdmissionregistrationV1Client, stopCh <-chan struct{}, mgr manager.Manager, auditor *auditlib.EventPublisher, restrictToNamespace string) error {
	informerFactory := informers.NewSharedInformerFactory(crdClient, time.Second*30)
	i := informerFactory.Apiextensions().V1().CustomResourceDefinitions().Informer()
	l := informerFactory.Apiextensions().V1().CustomResourceDefinitions().Lister()

	i.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			var key string
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				klog.Error(err)
				return
			}

			_, name, err := cache.SplitMetaNamespaceKey(key)
			if err != nil {
				klog.Error(err)
				return
			}

			crd, err := l.Get(name)
			if err != nil {
				klog.Error(err)
				return
			}
			if strings.Contains(crd.Spec.Group, "vultr.kubeform.com") {
				gvk := schema.GroupVersionKind{
					Group:   crd.Spec.Group,
					Version: crd.Spec.Versions[0].Name,
					Kind:    crd.Spec.Names.Kind,
				}

				// check whether this gvk came before, if no then start the controller
				runningControllers.RLock()
				_, ok := runningControllers.mp[gvk]
				runningControllers.RUnlock()

				if !ok {
					runningControllers.Lock()
					runningControllers.mp[gvk] = true
					runningControllers.Unlock()

					if enableValidatingWebhook {
						// add dynamic ValidatingWebhookConfiguration

						// create empty VWC if the group has come for the first time
						err := createEmptyVWC(vwcClient, gvk)
						if err != nil {
							klog.Error(err)
							return
						}

						// update
						err = updateVWC(vwcClient, gvk)
						if err != nil {
							klog.Error(err)
							return
						}

						err = SetupWebhook(mgr, gvk)
						if err != nil {
							setupLog.Error(err, "unable to enable webhook")
							os.Exit(1)
						}
					}

					err = SetupManager(ctx, mgr, gvk, auditor, restrictToNamespace)
					if err != nil {
						setupLog.Error(err, "unable to start manager")
						os.Exit(1)
					}
				}
			}
		},
	})

	informerFactory.Start(stopCh)

	return nil
}

func createEmptyVWC(vwcClient *admissionregistrationv1.AdmissionregistrationV1Client, gvk schema.GroupVersionKind) error {
	vwcName := strings.ReplaceAll(strings.ToLower(gvk.Group), ".", "-")
	_, err := vwcClient.ValidatingWebhookConfigurations().Get(context.TODO(), vwcName, metav1.GetOptions{})
	if err == nil || !(errors.IsNotFound(err)) {
		return err
	}

	emptyVWC := &arv1.ValidatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ValidatingWebhookConfiguration",
			APIVersion: "admissionregistration.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ReplaceAll(strings.ToLower(gvk.Group), ".", "-"),
			Labels: map[string]string{
				"app.kubernetes.io/instance": "vultr.kubeform.com",
				"app.kubernetes.io/part-of":  "kubeform.com",
			},
		},
	}
	_, err = vwcClient.ValidatingWebhookConfigurations().Create(context.TODO(), emptyVWC, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func updateVWC(vwcClient *admissionregistrationv1.AdmissionregistrationV1Client, gvk schema.GroupVersionKind) error {
	vwcName := strings.ReplaceAll(strings.ToLower(gvk.Group), ".", "-")
	vwc, err := vwcClient.ValidatingWebhookConfigurations().Get(context.TODO(), vwcName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	path := "/validate-" + strings.ReplaceAll(strings.ToLower(gvk.Group), ".", "-") + "-v1alpha1-" + strings.ToLower(gvk.Kind)
	fail := arv1.Fail
	sideEffects := arv1.SideEffectClassNone
	admissionReviewVersions := []string{"v1beta1"}

	rules := []arv1.RuleWithOperations{
		{
			Operations: []arv1.OperationType{
				arv1.Delete,
				arv1.Update,
			},
			Rule: arv1.Rule{
				APIGroups:   []string{strings.ToLower(gvk.Group)},
				APIVersions: []string{gvk.Version},
				Resources:   []string{strings.ToLower(flect.Pluralize(gvk.Kind))},
			},
		},
	}

	data, err := ioutil.ReadFile("/tmp/k8s-webhook-server/serving-certs/ca.crt")
	if err != nil {
		return err
	}

	name := strings.ToLower(gvk.Kind) + "." + gvk.Group
	for _, webhook := range vwc.Webhooks {
		if webhook.Name == name {
			return nil
		}
	}

	newWebhook := arv1.ValidatingWebhook{
		Name: name,
		ClientConfig: arv1.WebhookClientConfig{
			Service: &arv1.ServiceReference{
				Namespace: webhookNamespace,
				Name:      webhookName,
				Path:      &path,
			},
			CABundle: data,
		},
		Rules:                   rules,
		FailurePolicy:           &fail,
		SideEffects:             &sideEffects,
		AdmissionReviewVersions: admissionReviewVersions,
	}

	vwc.Webhooks = append(vwc.Webhooks, newWebhook)

	_, err = vwcClient.ValidatingWebhookConfigurations().Update(context.TODO(), vwc, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func SetupManager(ctx context.Context, mgr manager.Manager, gvk schema.GroupVersionKind, auditor *auditlib.EventPublisher, restrictToNamespace string) error {
	switch gvk {
	case schema.GroupVersionKind{
		Group:   "bare.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "MetalServer",
	}:
		if err := (&controllersbare.MetalServerReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("MetalServer"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_bare_metal_server"],
			TypeName: "vultr_bare_metal_server",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "MetalServer")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "block.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Storage",
	}:
		if err := (&controllersblock.StorageReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Storage"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_block_storage"],
			TypeName: "vultr_block_storage",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Storage")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "dns.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Domain",
	}:
		if err := (&controllersdns.DomainReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Domain"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_dns_domain"],
			TypeName: "vultr_dns_domain",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Domain")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "dns.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Record",
	}:
		if err := (&controllersdns.RecordReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Record"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_dns_record"],
			TypeName: "vultr_dns_record",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Record")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "firewall.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Group",
	}:
		if err := (&controllersfirewall.GroupReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Group"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_firewall_group"],
			TypeName: "vultr_firewall_group",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Group")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "firewall.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Rule",
	}:
		if err := (&controllersfirewall.RuleReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Rule"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_firewall_rule"],
			TypeName: "vultr_firewall_rule",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Rule")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "instance.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Instance",
	}:
		if err := (&controllersinstance.InstanceReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Instance"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_instance"],
			TypeName: "vultr_instance",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Instance")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "instance.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Ipv4",
	}:
		if err := (&controllersinstance.Ipv4Reconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Ipv4"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_instance_ipv4"],
			TypeName: "vultr_instance_ipv4",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Ipv4")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "iso.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Private",
	}:
		if err := (&controllersiso.PrivateReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Private"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_iso_private"],
			TypeName: "vultr_iso_private",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Private")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "kubernetes.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Kubernetes",
	}:
		if err := (&controllerskubernetes.KubernetesReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Kubernetes"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_kubernetes"],
			TypeName: "vultr_kubernetes",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Kubernetes")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "kubernetes.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "NodePools",
	}:
		if err := (&controllerskubernetes.NodePoolsReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("NodePools"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_kubernetes_node_pools"],
			TypeName: "vultr_kubernetes_node_pools",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "NodePools")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "load.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Balancer",
	}:
		if err := (&controllersload.BalancerReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Balancer"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_load_balancer"],
			TypeName: "vultr_load_balancer",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Balancer")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "object.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Storage",
	}:
		if err := (&controllersobject.StorageReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Storage"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_object_storage"],
			TypeName: "vultr_object_storage",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Storage")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "private.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Network",
	}:
		if err := (&controllersprivate.NetworkReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Network"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_private_network"],
			TypeName: "vultr_private_network",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Network")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "reserved.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Ip",
	}:
		if err := (&controllersreserved.IpReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Ip"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_reserved_ip"],
			TypeName: "vultr_reserved_ip",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Ip")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "reverse.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Ipv4",
	}:
		if err := (&controllersreverse.Ipv4Reconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Ipv4"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_reverse_ipv4"],
			TypeName: "vultr_reverse_ipv4",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Ipv4")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "reverse.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Ipv6",
	}:
		if err := (&controllersreverse.Ipv6Reconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Ipv6"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_reverse_ipv6"],
			TypeName: "vultr_reverse_ipv6",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Ipv6")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "snapshot.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Snapshot",
	}:
		if err := (&controllerssnapshot.SnapshotReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Snapshot"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_snapshot"],
			TypeName: "vultr_snapshot",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Snapshot")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "snapshot.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "FromURL",
	}:
		if err := (&controllerssnapshot.FromURLReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("FromURL"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_snapshot_from_url"],
			TypeName: "vultr_snapshot_from_url",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "FromURL")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "ssh.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Key",
	}:
		if err := (&controllersssh.KeyReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Key"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_ssh_key"],
			TypeName: "vultr_ssh_key",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Key")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "startup.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Script",
	}:
		if err := (&controllersstartup.ScriptReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("Script"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_startup_script"],
			TypeName: "vultr_startup_script",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Script")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "user.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "User",
	}:
		if err := (&controllersuser.UserReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("User"),
			Scheme:   mgr.GetScheme(),
			Gvk:      gvk,
			Provider: _provider,
			Resource: _provider.ResourcesMap["vultr_user"],
			TypeName: "vultr_user",
		}).SetupWithManager(ctx, mgr, auditor, restrictToNamespace); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "User")
			return err
		}

	default:
		return fmt.Errorf("Invalid CRD")
	}

	return nil
}

func SetupWebhook(mgr manager.Manager, gvk schema.GroupVersionKind) error {
	switch gvk {
	case schema.GroupVersionKind{
		Group:   "bare.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "MetalServer",
	}:
		if err := (&barev1alpha1.MetalServer{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "MetalServer")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "block.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Storage",
	}:
		if err := (&blockv1alpha1.Storage{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Storage")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "dns.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Domain",
	}:
		if err := (&dnsv1alpha1.Domain{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Domain")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "dns.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Record",
	}:
		if err := (&dnsv1alpha1.Record{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Record")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "firewall.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Group",
	}:
		if err := (&firewallv1alpha1.Group{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Group")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "firewall.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Rule",
	}:
		if err := (&firewallv1alpha1.Rule{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Rule")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "instance.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Instance",
	}:
		if err := (&instancev1alpha1.Instance{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Instance")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "instance.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Ipv4",
	}:
		if err := (&instancev1alpha1.Ipv4{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Ipv4")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "iso.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Private",
	}:
		if err := (&isov1alpha1.Private{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Private")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "kubernetes.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Kubernetes",
	}:
		if err := (&kubernetesv1alpha1.Kubernetes{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Kubernetes")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "kubernetes.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "NodePools",
	}:
		if err := (&kubernetesv1alpha1.NodePools{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "NodePools")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "load.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Balancer",
	}:
		if err := (&loadv1alpha1.Balancer{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Balancer")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "object.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Storage",
	}:
		if err := (&objectv1alpha1.Storage{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Storage")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "private.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Network",
	}:
		if err := (&privatev1alpha1.Network{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Network")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "reserved.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Ip",
	}:
		if err := (&reservedv1alpha1.Ip{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Ip")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "reverse.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Ipv4",
	}:
		if err := (&reversev1alpha1.Ipv4{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Ipv4")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "reverse.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Ipv6",
	}:
		if err := (&reversev1alpha1.Ipv6{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Ipv6")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "snapshot.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Snapshot",
	}:
		if err := (&snapshotv1alpha1.Snapshot{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Snapshot")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "snapshot.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "FromURL",
	}:
		if err := (&snapshotv1alpha1.FromURL{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "FromURL")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "ssh.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Key",
	}:
		if err := (&sshv1alpha1.Key{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Key")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "startup.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "Script",
	}:
		if err := (&startupv1alpha1.Script{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Script")
			return err
		}
	case schema.GroupVersionKind{
		Group:   "user.vultr.kubeform.com",
		Version: "v1alpha1",
		Kind:    "User",
	}:
		if err := (&userv1alpha1.User{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "User")
			return err
		}

	default:
		return fmt.Errorf("Invalid Webhook")
	}

	return nil
}
