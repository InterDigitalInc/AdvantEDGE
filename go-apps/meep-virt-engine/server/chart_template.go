/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package server

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/helm"
	model "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const SERVICE_PORT_MIN = 1
const SERVICE_PORT_MAX = 65535
const SERVICE_NODE_PORT_MIN = 30000
const SERVICE_NODE_PORT_MAX = 32767
const DEFAULT_DUMMY_CONTAINER_IMAGE = "nginx"

type DeploymentTemplate struct {
	Enabled                  string
	Name                     string
	ReplicaCount             string
	ApiVersion               string
	MatchLabels              []string
	TemplateLabels           []string
	ContainerName            string
	ContainerImageRepository string
	ContainerImagePullPolicy string
	ContainerEnvEnabled      string
	ContainerEnv             []string
	ContainerCommandEnabled  string
	ContainerCommand         []string
	ContainerCommandArg      []string
	GpuEnabled               string
	GpuType                  string
	GpuCount                 string
}

type ServiceTemplate struct {
	Enabled          string
	Name             string
	Namespace        string
	Labels           []string
	Selector         []string
	Type             string
	Ports            []ServicePortTemplate
	MeServiceEnabled string
	MeServiceName    string
}

type ServicePortTemplate struct {
	Port       string
	TargetPort string
	Protocol   string
	NodePort   string
}

type ExternalTemplate struct {
	Enabled           string
	Selector          []string
	IngressServiceMap []ServiceMapTemplate
	EgressServiceMap  []ServiceMapTemplate
}

type ServiceMapTemplate struct {
	Name     string
	IP       string
	Port     string
	NodePort string
	Protocol string
}

// helm values.yaml template
type ScenarioTemplate struct {
	Deployment    DeploymentTemplate
	Service       ServiceTemplate
	External      ExternalTemplate
	NamespaceName string
}

// Service map
var serviceMap map[string]string

func addTemplateLabel(deploymentTemplate *DeploymentTemplate, label string) {
	deploymentTemplate.TemplateLabels = append(deploymentTemplate.TemplateLabels, label)
}

func addMatchLabel(deploymentTemplate *DeploymentTemplate, label string) {
	deploymentTemplate.MatchLabels = append(deploymentTemplate.MatchLabels, label)
}

func addServiceLabel(serviceTemplate *ServiceTemplate, label string) {
	serviceTemplate.Labels = append(serviceTemplate.Labels, label)
}

func addSelector(serviceTemplate *ServiceTemplate, selector string) {
	serviceTemplate.Selector = append(serviceTemplate.Selector, selector)
}

func addExtSelector(externalTemplate *ExternalTemplate, selector string) {
	externalTemplate.Selector = append(externalTemplate.Selector, selector)
}

func populateScenarioTemplate(scenario model.Scenario) ([]helm.Chart, error) {

	var charts []helm.Chart
	serviceMap = map[string]string{}

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

						// Create default scenario template
						var scenarioTemplate ScenarioTemplate
						deploymentTemplate := &scenarioTemplate.Deployment
						serviceTemplate := &scenarioTemplate.Service
						externalTemplate := &scenarioTemplate.External
						setScenarioDefaults(&scenarioTemplate)

						// Fill general scenario template information
						scenarioTemplate.NamespaceName = scenario.Name
						deploymentTemplate.Name = proc.Name

						// Create charts
						if proc.UserChartLocation != "" {
							log.Debug("Processing user-defined chart for element[", proc.Name, "]")

							// Add user-defined chart
							newChart := createChart(scenario.Name+"-"+proc.Name, getFullPath(proc.UserChartLocation),
								getFullPath(proc.UserChartAlternateValues))
							charts = append(charts, newChart)
							log.Debug("user chart added ", len(charts))

							// Parse User Chart Group to find new group services
							// Create charts only for group services that do not exist yet
							// Format: <service instance name>:[group service name]:<port>:<protocol>
							if proc.UserChartGroup != "" {
								userChartGroup := strings.Split(proc.UserChartGroup, ":")
								meSvcName := userChartGroup[1]
								if meSvcName != "" {
									if _, found := serviceMap[meSvcName]; !found {
										serviceMap[meSvcName] = "meepMeSvc: " + meSvcName
										serviceTemplate.MeServiceEnabled = "true"
										serviceTemplate.MeServiceName = meSvcName
										addServiceLabel(serviceTemplate, "meepMeSvc: "+meSvcName)

										serviceTemplate.Namespace = scenario.Name
										addServiceLabel(serviceTemplate, "meepScenario: "+scenario.Name)

										// NOTE: Every service within a group must expose the same port & protocol
										var portTemplate ServicePortTemplate
										portTemplate.Port = userChartGroup[2]
										portTemplate.Protocol = userChartGroup[3]
										serviceTemplate.Ports = append(serviceTemplate.Ports, portTemplate)

										// Create chart files
										chartLocation, err := createYamlScenarioFiles(scenarioTemplate)
										if err != nil {
											log.Debug("yaml creation file process: ", err)
											return nil, err
										}

										// Create virt-engine chart for new group service
										newChart := createChart(scenario.Name+"-"+proc.Name+"-svc", chartLocation, "")
										charts = append(charts, newChart)
										log.Debug("chart added for user chart group service ", len(charts))
									}
								}
							}
						} else {
							log.Debug("Processing virt-engine chart for element[", proc.Name, "]")

							// Fill deployment template information
							deploymentTemplate.Enabled = "true"
							deploymentTemplate.ContainerName = proc.Name
							deploymentTemplate.ContainerImageRepository = proc.Image
							deploymentTemplate.ContainerImagePullPolicy = "Always"
							setEnv(deploymentTemplate, proc.Environment)
							setCommand(deploymentTemplate, proc.CommandExe, proc.CommandArguments)
							addMatchLabel(deploymentTemplate, "meepAppId: "+proc.Id)
							addTemplateLabel(deploymentTemplate, "meepAppId: "+proc.Id)

							// Enable Service template if present
							if proc.ServiceConfig != nil {

								// Add app name associated to service
								svcName := proc.ServiceConfig.Name
								serviceTemplate.Enabled = "true"
								serviceTemplate.Name = svcName
								serviceTemplate.Namespace = scenario.Name
								addSelector(serviceTemplate, "meepSvc: "+svcName)
								addServiceLabel(serviceTemplate, "meepScenario: "+scenario.Name)
								addTemplateLabel(deploymentTemplate, "meepSvc: "+svcName)

								// Create and store ME Service template only with first occurrence.
								// If it already exists then add the matching pod label but don't create the service again.
								meSvcName := proc.ServiceConfig.MeSvcName
								if meSvcName != "" {
									if _, found := serviceMap[meSvcName]; !found {
										serviceMap[meSvcName] = "meepMeSvc: " + meSvcName
										serviceTemplate.MeServiceEnabled = "true"
										serviceTemplate.MeServiceName = meSvcName
									}
									addServiceLabel(serviceTemplate, "meepMeSvc: "+meSvcName)
									addTemplateLabel(deploymentTemplate, "meepMeSvc: "+meSvcName)
								}

								for _, ports := range proc.ServiceConfig.Ports {
									var portTemplate ServicePortTemplate
									portTemplate.Port = strconv.Itoa(int(ports.Port))
									portTemplate.TargetPort = strconv.Itoa(int(ports.Port))
									portTemplate.Protocol = ports.Protocol

									// Add NodePort if service is exposed externally
									if ports.ExternalPort >= SERVICE_NODE_PORT_MIN && ports.ExternalPort <= SERVICE_NODE_PORT_MAX {
										portTemplate.NodePort = strconv.Itoa(int(ports.ExternalPort))
										serviceTemplate.Type = "NodePort"
									} else {
										serviceTemplate.Type = "ClusterIP"
									}

									serviceTemplate.Ports = append(serviceTemplate.Ports, portTemplate)
								}
							}

							// Enable GPU template if present
							if proc.GpuConfig != nil {
								deploymentTemplate.GpuEnabled = "true"
								deploymentTemplate.GpuType = proc.GpuConfig.Type_
								deploymentTemplate.GpuCount = strconv.Itoa(int(proc.GpuConfig.Count))
							}

							// Enable External template if set
							if proc.IsExternal {
								externalTemplate.Enabled = "true"
								addExtSelector(externalTemplate, "meepAppId: "+proc.Id)

								// Add ingress Service Maps, if any
								for _, serviceMap := range proc.ExternalConfig.IngressServiceMap {
									var ingressSvcMapTemplate ServiceMapTemplate
									ingressSvcMapTemplate.NodePort = strconv.Itoa(int(serviceMap.ExternalPort))
									ingressSvcMapTemplate.Port = strconv.Itoa(int(serviceMap.Port))
									ingressSvcMapTemplate.Protocol = serviceMap.Protocol
									ingressSvcMapTemplate.Name = "ingress-" + proc.Id + "-" + ingressSvcMapTemplate.NodePort

									externalTemplate.IngressServiceMap = append(externalTemplate.IngressServiceMap, ingressSvcMapTemplate)
								}

								// Add egress Service Maps, if any
								for _, serviceMap := range proc.ExternalConfig.EgressServiceMap {
									var egressSvcMapTemplate ServiceMapTemplate
									egressSvcMapTemplate.Name = serviceMap.Name
									egressSvcMapTemplate.IP = serviceMap.Ip
									egressSvcMapTemplate.Port = strconv.Itoa(int(serviceMap.Port))
									egressSvcMapTemplate.Protocol = serviceMap.Protocol

									externalTemplate.EgressServiceMap = append(externalTemplate.EgressServiceMap, egressSvcMapTemplate)
								}
							}

							// Create chart files
							chartLocation, err := createYamlScenarioFiles(scenarioTemplate)
							if err != nil {
								log.Debug("yaml creation file process: ", err)
								return nil, err
							}

							// Create virt-engine chart
							newChart := createChart(scenario.Name+"-"+proc.Name, chartLocation, "")
							charts = append(charts, newChart)
							log.Debug("chart added ", len(charts))
						}
					}
				}
			}
		}
	}

	return charts, nil
}

func createChart(name string, chartLocation string, valuesFile string) helm.Chart {
	var chart helm.Chart
	chart.ChartName = name
	chart.ReleaseName = "meep-" + name
	chart.Location = chartLocation
	chart.ValuesFile = valuesFile
	return chart
}

func getFullPath(path string) string {
	fullPath := path

	// Get home directory
	usr, err := user.Current()
	if err != nil {
		return fullPath
	}
	homeDir := usr.HomeDir

	// Replace ~ with home directory
	if path == "~" {
		fullPath = homeDir
	} else if strings.HasPrefix(path, "~/") {
		fullPath = filepath.Join(homeDir, path[2:])
	}
	return fullPath
}

func setScenarioDefaults(scenarioTemplate *ScenarioTemplate) {
	setDeploymentDefaults(&scenarioTemplate.Deployment)
	setServiceDefaults(&scenarioTemplate.Service)
	setExternalDefaults(&scenarioTemplate.External)
}

func setDeploymentDefaults(deploymentTemplate *DeploymentTemplate) {
	deploymentTemplate.Enabled = "false"
	deploymentTemplate.ReplicaCount = "1"
	deploymentTemplate.ApiVersion = "v1"
	deploymentTemplate.ContainerEnvEnabled = "false"
	deploymentTemplate.ContainerCommandEnabled = "false"
	deploymentTemplate.GpuEnabled = "false"
}

func setServiceDefaults(serviceTemplate *ServiceTemplate) {
	serviceTemplate.Enabled = "false"
	serviceTemplate.MeServiceEnabled = "false"
}

func setExternalDefaults(externalTemplate *ExternalTemplate) {
	externalTemplate.Enabled = "false"
}

func setEnv(deployment *DeploymentTemplate, envString string) {
	if envString != "" {
		deployment.ContainerEnvEnabled = "true"
		allVar := strings.Split(envString, ",")

		for _, oneVar := range allVar {
			nameValue := strings.Split(oneVar, "=")
			deployment.ContainerEnv = append(deployment.ContainerEnv,
				strings.TrimSpace(nameValue[0])+": "+strings.TrimSpace(nameValue[1]))
		}
	}
}

func setCommand(deployment *DeploymentTemplate, command string, commandArgs string) {
	if command != "" {
		log.Debug("command ", command)
		deployment.ContainerCommandEnabled = "true"

		// Retrieve command list
		allCmd := strings.Split(command, ",")
		for _, cmd := range allCmd {
			deployment.ContainerCommand = append(deployment.ContainerCommand, strings.TrimSpace(cmd))
		}

		// Retrieve arguments list
		allArgs := strings.Split(commandArgs, ",")
		for _, arg := range allArgs {
			deployment.ContainerCommandArg = append(deployment.ContainerCommandArg, strings.TrimSpace(arg))
		}
	}
}

func CreateYamlScenarioFile(scenario model.Scenario) error {

	//var charts []helm.Chart
	charts, err := populateScenarioTemplate(scenario)

	if err != nil {
		log.Debug("populate template : ", err)
		return err
	}

	err = deployCharts(charts)
	if err != nil {
		log.Error("charts error : ", err)
		return err
	}

	return nil
}

func createYamlScenarioFiles(scenarioTemplate ScenarioTemplate) (string, error) {

	homePath := os.Getenv("HOME")

	templateFilePath := homePath + "/.meep/template/values-template.yaml"
	templateDefaultDir := homePath + "/.meep/template/defaultDir"

	t, err := template.ParseFiles(templateFilePath)
	if err != nil {
		log.Error(err)
		return "", err
	}

	outputDirPath := homePath + "/.meep/active/" + scenarioTemplate.NamespaceName + "/" + scenarioTemplate.Deployment.Name
	log.Debug("Creation of the output path ", outputDirPath)

	_ = CopyDir(templateDefaultDir, outputDirPath)

	outputFilePath := outputDirPath + "/values.yaml"

	//creation of output file
	f, err := os.Create(outputFilePath)
	if err != nil {
		log.Debug("create file: ", err)
		return "", err
	}

	//filling the template output file
	err = t.Execute(f, scenarioTemplate)
	if err != nil {
		log.Debug("execute: ", err)
		return "", err
	}

	f.Close()
	return outputDirPath, nil
}

func deployCharts(charts []helm.Chart) error {

	err := helm.InstallCharts(charts)
	if err != nil {
		return err
	}
	return nil
}
