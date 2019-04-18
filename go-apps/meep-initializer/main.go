// Copyright 2017 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-initializer/log"
	"github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"

	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/getlantern/deepcopy"
	"github.com/ghodss/yaml"

	"k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/watch"
	//"k8s.io/apimachinery/pkg/labels" 
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const (
	defaultAnnotation      = "initializer.kubernetes.io/sidecar"
	defaultInitializerName = "sidecar.initializer.kubernetes.io"
	defaultConfigmap       = "meep-initializer"
	defaultNamespace       = "default"
)

var (
	annotation        string
	configmap         string
	initializerName   string
	namespace         string
	requireAnnotation bool
)

type config struct {
	Containers []corev1.Container
	Volumes    []corev1.Volume
}

func getMeSvcName(sourceType string, str string, uniqueName string) (string, string) {

	//2 sourceType
	//"svc", "pod", "statefulset"
	var scenario model.Scenario
	_ = json.Unmarshal([]byte(str), &scenario) 
        // Parse domains
        for _, domain := range scenario.Deployment.Domains {
        	// Parse zones
                for _, zone := range domain.Zones {
                        // Parse Network Locations
                        for _, nl := range zone.NetworkLocations {
                                // Parse Physical locations
                                for _, pl := range nl.PhysicalLocations {
                                        // Parse Processes
                                        for _, proc := range pl.Processes {
                                               	if proc.UserChartLocation != "" {
							//entry from the scenario was found, return service config
							//return svcName and meSvcName
							userChartGroupElement := strings.Split(proc.UserChartGroup, ":")
							if userChartGroupElement[0] == uniqueName {
								return userChartGroupElement[0], userChartGroupElement[1]
							} else {
								//pods use uniqueName which is same as webUi, while the service can have different one 
								if sourceType != "svc" && proc.Name == uniqueName {
									return userChartGroupElement[0], userChartGroupElement[1]
								}
							}
						} else {
							if proc.Name == uniqueName {
								return proc.ServiceConfig.Name, proc.ServiceConfig.MeSvcName
							}
						}
					}
				}
			}
		}
	}
	return "", ""
}

func main() {

	flag.StringVar(&annotation, "annotation", defaultAnnotation, "The annotation to trigger initialization")
	flag.StringVar(&configmap, "configmap", defaultConfigmap, "The sidecar initializer configuration configmap")
	flag.StringVar(&initializerName, "initializer-name", defaultInitializerName, "The initializer name")
	flag.StringVar(&namespace, "namespace", "default", "The configuration namespace")
	flag.BoolVar(&requireAnnotation, "require-annotation", true, "Require annotation for initialization")
	flag.Parse()

	log.Info("Starting the Kubernetes initializer...")
	log.Info("Initializer name set to: ", initializerName)
	log.Info("Initializer annotation set to: ", annotation)
	log.Info("Initializer bool set to: ", requireAnnotation)

	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Load the Envoy Initializer configuration from a Kubernetes ConfigMap.
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(configmap, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	c, err := configmapToConfig(cm)
	if err != nil {
		log.Fatal(err)
	}

	// Watch uninitialized Deployments in all namespaces.
	restClient := clientset.AppsV1beta1().RESTClient()
	restClientCore := clientset.CoreV1().RESTClient()
	watchlist_deployments := cache.NewListWatchFromClient(restClient, "deployments", corev1.NamespaceAll, fields.Everything())
	watchlist_statefulsets := cache.NewListWatchFromClient(restClient, "statefulsets", corev1.NamespaceAll, fields.Everything())
	watchlist_services := cache.NewListWatchFromClient(restClientCore, "services", corev1.NamespaceAll, fields.Everything())



	//connect to scenario DB
	nbOfErrorsAllowed := 5

	// Connect to Redis DB
	for nbOfErrorsAllowed > 0 {
		err = RedisDBConnect()
		if err == nil {
			break
		}
		log.Error("Failed to connect to Redis DB. Error: ", err)
		nbOfErrorsAllowed--
		time.Sleep(2 * time.Second)
	}

	if nbOfErrorsAllowed == 0 {
		log.Fatal(err)
	}

	// Wrap the returned watchlist to workaround the inability to include
	// the `IncludeUninitialized` list option when setting up watch clients.
	includeUninitializedWatchlist_deployments := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.IncludeUninitialized = true
			return watchlist_deployments.List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.IncludeUninitialized = true
			return watchlist_deployments.Watch(options)
		},
	}

        includeUninitializedWatchlist_statefulsets := &cache.ListWatch{
                ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
                        options.IncludeUninitialized = true
                        return watchlist_statefulsets.List(options)
                },
                WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
                        options.IncludeUninitialized = true
                        return watchlist_statefulsets.Watch(options)
                },
        }

        includeUninitializedWatchlist_services := &cache.ListWatch{
                ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
                        options.IncludeUninitialized = true
                        return watchlist_services.List(options)
                },
                WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
                        options.IncludeUninitialized = true
                        return watchlist_services.Watch(options)
                },
        }

	resyncPeriod := 30 * time.Second

	_, controller := cache.NewInformer(includeUninitializedWatchlist_deployments, &v1beta1.Deployment{}, resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				deployment, ok := obj.(*v1beta1.Deployment)
				if ok {
					err := initializeDeployment(deployment, c, clientset)
					if err != nil {
						log.Error(err)
					}
				}
			},
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)

        _, controller2 := cache.NewInformer(includeUninitializedWatchlist_statefulsets, &v1beta1.StatefulSet{}, resyncPeriod,
                cache.ResourceEventHandlerFuncs{
                        AddFunc: func(obj interface{}) {
				statefulset, ok := obj.(*v1beta1.StatefulSet)
				if ok {
                                	err := initializeStatefulSet(statefulset, c, clientset)
	                                if err != nil {
       		                        	log.Error(err)
					}
				}
			},
		},
	)

	stop2 := make(chan struct{})
	go controller2.Run(stop2)

        _, controller3 := cache.NewInformer(includeUninitializedWatchlist_services, &corev1.Service{}, resyncPeriod,
                cache.ResourceEventHandlerFuncs{
                        AddFunc: func(obj interface{}) {
				service, ok := obj.(*corev1.Service)
				if ok {
                                        err := initializeService(service, c, clientset)
                                        if err != nil {
                                                log.Error(err)
                                        }
                                } 
                        },
                },
        )

        stop3 := make(chan struct{})
        go controller3.Run(stop3)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Info("Shutdown signal received, exiting...")
	close(stop)
	close(stop2)
	close(stop3)

}

func initializeService(service *corev1.Service, c *config, clientset *kubernetes.Clientset) error {
        log.Info("service ", service.Name)

        if service.ObjectMeta.GetInitializers() != nil {

                pendingInitializers := service.ObjectMeta.GetInitializers().Pending

                if initializerName == pendingInitializers[0].Name {
                        log.Info("Initializing a new service: ", service.Name)

                        var newService corev1.Service
                        err := deepcopy.Copy(&newService, service)
                        if err != nil {
                                log.Info("Deepcopy error")
				return err
                        }

                        // Remove self from the list of pending Initializers while preserving ordering.
                        if len(pendingInitializers) == 1 {
                                newService.ObjectMeta.Initializers = nil
                        } else {
                                newService.ObjectMeta.Initializers.Pending = append(pendingInitializers[:0], pendingInitializers[1:]...)
                        }


                        //comparing to see if the pod should be updated
                        //looking at the pod name if it contains the name of the scenario
                        //

                        jsonScenarioName := ""
                        scenarioName := ""
                        isPresent := false

                        jsonScenarioFull, err := RedisDBJsonGetEntry("ctrl-engine:active", "")
                        if err == nil {
                                jsonScenarioName, err = RedisDBJsonGetEntry("ctrl-engine:active", "name")
                                if err == nil {
                                        //removing extra character '\'
                                        scenarioName = jsonScenarioName[1:len(jsonScenarioName)-1]
                                        isPresent = strings.Contains(service.Name, scenarioName)
                                } else {
                                        log.Error(err.Error())
                                }
				//no check if service is part of the meep deployed objects, we only check for deployments to add sidecars
				isPresent = true
	                        if isPresent == false {
	                                log.Info("Required criteria not met for pod name; skipping sidecar container injection: ", scenarioName, " not in ", service.Name)
	                                _, err = clientset.CoreV1().Services(service.Namespace).Update(&newService)
	                                if err != nil {
	                                        log.Error(err.Error())
	                                }
	                        } else {
	                                // Modify the Service's template to include new labels
	                                // and configuration volume. Then patch the original service.
       	                        	origSvcLabels := service.ObjectMeta.GetLabels()
	                                newSvcLabels := make(map[string]string)

					//we keep the original labels
	                                for k, v := range origSvcLabels {
	                                        newSvcLabels[k] = v
	                                }

					meepAppName := service.Name//[len(scenarioName)+6:]
					svcName, meSvcName := getMeSvcName("svc", jsonScenarioFull, meepAppName)

					if svcName != "" {
						for k := range newService.Spec.Selector {
						    delete(newService.Spec.Selector, k)
						}
	                                	newService.Spec.Selector["meepSvc"] = svcName
					}
	                                if meSvcName != "" {
	                                        newSvcLabels["meepMeSvc"] = meSvcName
	                                }

					newSvcLabels["meepScenario"] = scenarioName
       	                         	newService.ObjectMeta.SetLabels(newSvcLabels)
					
       	                 	}

                        } else {
                                _, err = clientset.CoreV1().Services(service.Namespace).Update(&newService)
                                log.Info("No active scenario in Redis DB")
                        }

                        origData, err := json.Marshal(service)
                        if err != nil {
                                return err
                        }

                        newData, err := json.Marshal(newService)
                        if err != nil {
                                return err
                        }

                        patchBytes, err := strategicpatch.CreateTwoWayMergePatch(origData, newData, corev1.Service{})
                        if err != nil {
                                return err
                        }


                        _, err = clientset.CoreV1().Services(service.Namespace).Patch(service.Name, types.StrategicMergePatchType, patchBytes)
                        if err != nil {
                                return err
                        }

		}
	}
	return nil
}

func initializeStatefulSet(statefulset *v1beta1.StatefulSet, c *config, clientset *kubernetes.Clientset) error {
	log.Info("statefulset ", statefulset.Name)
	if statefulset.ObjectMeta.GetInitializers() != nil {

		pendingInitializers := statefulset.ObjectMeta.GetInitializers().Pending

		if initializerName == pendingInitializers[0].Name {
			log.Info("Initializing a new statefulset: ", statefulset.Name)

			var newStatefulset v1beta1.StatefulSet
			err := deepcopy.Copy(&newStatefulset, statefulset)
			if err != nil {
				return err
			}

			// Remove self from the list of pending Initializers while preserving ordering.
			if len(pendingInitializers) == 1 {
				newStatefulset.ObjectMeta.Initializers = nil
			} else {
				newStatefulset.ObjectMeta.Initializers.Pending = append(pendingInitializers[:0], pendingInitializers[1:]...)
			}

			//comparing to see if the pod should be updated
			//looking at the pod name if it contains the name of the scenario
			//

			jsonScenarioName := ""
			scenarioName := ""
			isPresent := false

			jsonScenarioFull, err := RedisDBJsonGetEntry("ctrl-engine:active", "")
			if err == nil {
				jsonScenarioName, err = RedisDBJsonGetEntry("ctrl-engine:active", "name")
				if err == nil {
					//removing extra character '\'
					scenarioName = jsonScenarioName[1 : len(jsonScenarioName)-1]
					isPresent = strings.Contains(statefulset.Name, scenarioName)
				} else {
					log.Error(err.Error())
				}
			} else {
				log.Error(err.Error())
			}

			if isPresent == false {
				log.Info("Required criteria not met for pod name; skipping sidecar container injection: ", scenarioName, " not in ", statefulset.Name)
				_, err = clientset.AppsV1beta1().StatefulSets(statefulset.Namespace).Update(&newStatefulset)
                                if err != nil {
                                        log.Error(err.Error())
                                }
                        } else {

                                // Modify the Statefulset's Pod template to include the sidecar container
                                // and configuration volume. Then patch the original statefulset.
                                newStatefulset.Spec.Template.Spec.Containers = append(statefulset.Spec.Template.Spec.Containers, c.Containers...)
                                newStatefulset.Spec.Template.Spec.Volumes = append(statefulset.Spec.Template.Spec.Volumes, c.Volumes...)
                                origPodLabels := statefulset.Spec.Template.ObjectMeta.GetLabels()
                                newPodLabels := make(map[string]string)

                                for k, v := range origPodLabels {
                                        newPodLabels[k] = v
                                }
                                if newPodLabels["meepApp"] == "" {
                                        log.Info("This pod does not already have a sidecar added by the virtual engine")
                                        //finding the meepApp based on the name of the pod and the name of the scenario.. meep-{scenarioName}-{uniqueName == meepApp}
                                        meepAppName := statefulset.Name[len(scenarioName)+6:]
                                        newPodLabels["meepApp"] = meepAppName
                                        newPodLabels["meepOrigin"] = "scenario"
                                        newPodLabels["meepScenario"] = scenarioName
                                        newPodLabels["processId"] = meepAppName

                                        svcName, meSvcName := getMeSvcName("statefulset", jsonScenarioFull, meepAppName)

                                        if svcName != "" {
                                                newPodLabels["meepSvc"] = svcName
                                        }

                                        if meSvcName != "" {
                                                newPodLabels["meepMeSvc"] = meSvcName
                                        }
                                } else {
                                        log.Info("This pod already has a sidecar added by the virtual engine, no need to do anything")
                                }

                                newStatefulset.Spec.Template.ObjectMeta.SetLabels(newPodLabels)
                        }

	                origData, err := json.Marshal(statefulset)
       		        if err != nil {
               		        return err
	                }

	                newData, err := json.Marshal(newStatefulset)
	                if err != nil {
	                        return err
	                }

	                patchBytes, err := strategicpatch.CreateTwoWayMergePatch(origData, newData, v1beta1.StatefulSet{})
	                if err != nil {
	                        return err
	                }

	                _, err = clientset.AppsV1beta1().StatefulSets(statefulset.Namespace).Patch(statefulset.Name, types.StrategicMergePatchType, patchBytes)
	                if err != nil {
	                        return err
			}
		}
	}
	return nil
}

func initializeDeployment(deployment *v1beta1.Deployment, c *config, clientset *kubernetes.Clientset) error {
	log.Info("deployment ", deployment.Name)
	if deployment.ObjectMeta.GetInitializers() != nil {
		pendingInitializers := deployment.ObjectMeta.GetInitializers().Pending

		if initializerName == pendingInitializers[0].Name {
			log.Info("Initializing a new deployment: ", deployment.Name)

			//o, err := runtime.NewScheme().DeepCopy(deployment)
			var newDeployment v1beta1.Deployment
			err := deepcopy.Copy(&newDeployment, deployment)
			if err != nil {
				return err
			}

			// Remove self from the list of pending Initializers while preserving ordering.
			if len(pendingInitializers) == 1 {
				newDeployment.ObjectMeta.Initializers = nil
			} else {
				newDeployment.ObjectMeta.Initializers.Pending = append(pendingInitializers[:0], pendingInitializers[1:]...)
			}

			//comparing to see if the pod should be updated
			//looking at the pod name if it contains the name of the scenario
			//

			jsonScenarioName := ""
			scenarioName := ""
			isPresent := false

			jsonScenarioFull, err := RedisDBJsonGetEntry("ctrl-engine:active", "")
			if err == nil {
				jsonScenarioName, err = RedisDBJsonGetEntry("ctrl-engine:active", "name")
				if err == nil {
					//removing extra character '\'
					scenarioName = jsonScenarioName[1 : len(jsonScenarioName)-1]
					isPresent = strings.Contains(deployment.Name, scenarioName)
				} else {
					log.Error(err.Error())
				}
			} else {
				log.Error(err.Error())
			}

			if isPresent == false {
				log.Info("Required criteria not met for pod name; skipping sidecar container injection: ", scenarioName, " not in ", deployment.Name)
				_, err = clientset.AppsV1beta1().Deployments(deployment.Namespace).Update(&newDeployment)
				if err != nil {
					log.Error(err.Error())
				}
			} else {

				// Modify the Deployment's Pod template to include the sidecar container
				// and configuration volume. Then patch the original deployment.
				newDeployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, c.Containers...)
				newDeployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, c.Volumes...)
				origPodLabels := deployment.Spec.Template.ObjectMeta.GetLabels()
				newPodLabels := make(map[string]string)

				for k, v := range origPodLabels {
					newPodLabels[k] = v
				}
				if newPodLabels["meepApp"] == "" {
					log.Info("This pod does not already have a sidecar added by the virtual engine")
					//finding the meepApp based on the name of the pod and the name of the scenario.. meep-{scenarioName}-{uniqueName == meepApp}
					//str := strings.TrimPrefix(deployment.Name, scenarioName)
					//mylen := len(str)
					meepAppName := deployment.Name[len(scenarioName)+6:]
					newPodLabels["meepApp"] = meepAppName
					newPodLabels["meepOrigin"] = "scenario"
					newPodLabels["meepScenario"] = scenarioName
					newPodLabels["processId"] = meepAppName

					svcName, meSvcName := getMeSvcName("pod", jsonScenarioFull, meepAppName)

					if svcName != "" {
						newPodLabels["meepSvc"] = svcName
					}

					if meSvcName != "" {
						newPodLabels["meepMeSvc"] = meSvcName
					}


	                               	newDeployment.Spec.Template.ObjectMeta.SetLabels(newPodLabels)
	                                var envVars []corev1.EnvVar
	                                var envVar corev1.EnvVar
	                                envVar.Name = "MEEP_POD_NAME"
	                                envVar.Value = meepAppName
	                                envVars = append(envVars, envVar)
	                                newDeployment.Spec.Template.Spec.Containers[1].Env = envVars

				} else {
					log.Info("This pod already has a sidecar added by the virtual engine, no need to do anything")
				}

				newDeployment.Spec.Template.ObjectMeta.SetLabels(newPodLabels)
			}

			origData, err := json.Marshal(deployment)
			if err != nil {
				return err
			}

			newData, err := json.Marshal(newDeployment)
			if err != nil {
				return err
			}

			patchBytes, err := strategicpatch.CreateTwoWayMergePatch(origData, newData, v1beta1.Deployment{})
			if err != nil {
				return err
			}

			_, err = clientset.AppsV1beta1().Deployments(deployment.Namespace).Patch(deployment.Name, types.StrategicMergePatchType, patchBytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func configmapToConfig(configmap *corev1.ConfigMap) (*config, error) {
	var c config
	err := yaml.Unmarshal([]byte(configmap.Data["config"]), &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
