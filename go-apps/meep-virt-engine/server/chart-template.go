/*
 * Copyright (c) 2019  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/helm"
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
)

const serviceNodePortMin = 30000
const serviceNodePortMax = 32767
const trueStr = "true"
const falseStr = "false"

// DeploymentTemplate - Deployment Template
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
	PlacementId              string
	CpuEnabled               string
	CpuMin                   string
	CpuMax                   string
	MemoryEnabled            string
	MemoryMin                string
	MemoryMax                string
}

// ServiceTemplate - Service Template
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

// ServicePortTemplate - Service Port Template
type ServicePortTemplate struct {
	Port       string
	TargetPort string
	Protocol   string
	NodePort   string
}

// ExternalTemplate -  External Template
type ExternalTemplate struct {
	Enabled           string
	Selector          []string
	IngressServiceMap []IngressServiceTemplate
	EgressServiceMap  []EgressServiceTemplate
}

// IngressServiceTemplate - Ingress Service Template
type IngressServiceTemplate struct {
	Name     string
	Port     string
	NodePort string
	Protocol string
}

// EgressServiceTemplate - Egress Service Template
type EgressServiceTemplate struct {
	Name      string
	MeSvcName string
	IP        string
	Port      string
	Protocol  string
}

// ScenarioTemplate -helm values.yaml template
type ScenarioTemplate struct {
	Name       string
	Deployment DeploymentTemplate
	Service    ServiceTemplate
	External   ExternalTemplate
	Namespace  string
}

// SandboxTemplate -helm values.yaml template
type SandboxTemplate struct {
	SandboxName    string
	Namespace      string
	HostUrl        string
	UserSwagger    string
	UserSwaggerDir string
	HttpsOnly      bool
	AuthEnabled    bool
	IsMepService   bool
	MepName        string
}

// Deploy - Generate charts & deploy single process or entire scenario
func Deploy(sandboxName string, procName string, model *mod.Model) error {

	// Create scenario charts
	charts, err := generateScenarioCharts(sandboxName, procName, model)
	if err != nil {
		log.Debug("Error creating scenario charts: ", err)
		return err
	}
	log.Debug("Created ", len(charts), " scenario charts")

	// Deploy all charts
	err = deployCharts(charts, sandboxName)
	if err != nil {
		log.Error("Error deploying charts: ", err)
		return err
	}

	return nil
}

func getMepService(proc *dataModel.Process) string {
	// !!! Temporary patch for MEP Service configuration !!!
	// Use well-known Edge App Environment variable to obtain MEP Service Name
	if proc != nil && proc.Environment != "" {
		allVar := strings.Split(proc.Environment, ",")
		for _, oneVar := range allVar {
			nameValue := strings.Split(oneVar, "=")
			if nameValue[0] == "MEEP_MEP_SERVICE" {
				return nameValue[1]
			}
		}
	}
	return ""
}

func generateScenarioCharts(sandboxName string, procName string, model *mod.Model) (charts []helm.Chart, err error) {
	serviceMap := map[string]string{}

	procNames := model.GetNodeNames("CLOUD-APP")
	procNames = append(procNames, model.GetNodeNames("EDGE-APP")...)
	procNames = append(procNames, model.GetNodeNames("UE-APP")...)
	for _, name := range procNames {
		// Check if single process is being added
		if procName != "" && name != procName {
			continue
		}

		// Retrieve node process information from model
		node := model.GetNode(name)
		if node == nil {
			err = errors.New("Error finding process: " + name)
			return nil, err
		}
		proc, ok := node.(*dataModel.Process)
		if !ok {
			err = errors.New("Error casting process: " + name)
			return nil, err
		}
		ctx := model.GetNodeContext(name)
		if ctx == nil {
			err = errors.New("Error finding context for process: " + name)
			return nil, err
		}

		scenarioName := model.GetScenarioName()

		// Create default scenario template
		var scenarioTemplate ScenarioTemplate
		deploymentTemplate := &scenarioTemplate.Deployment
		serviceTemplate := &scenarioTemplate.Service
		externalTemplate := &scenarioTemplate.External
		setScenarioDefaults(&scenarioTemplate)

		// Fill general scenario template information
		scenarioTemplate.Name = scenarioName
		scenarioTemplate.Namespace = sandboxName
		deploymentTemplate.Name = proc.Name

		// Create charts
		if proc.UserChartLocation != "" {
			log.Debug("Processing user-defined chart for element[", proc.Name, "]")

			// Add user-defined chart
			nc := newChart(proc.Name, sandboxName, scenarioName,
				getFullPath(proc.UserChartLocation), getFullPath(proc.UserChartAlternateValues))
			charts = append(charts, nc)
			log.Debug("user chart added ", len(charts))

			// Parse User Chart Group to find new group services
			// Create charts only for group services that do not exist yet
			// Format: <service instance name>:[group service name]:<port>:<protocol>
			if proc.UserChartGroup != "" {
				userChartGroup := strings.Split(proc.UserChartGroup, ":")
				meSvcName := userChartGroup[1]
				if meSvcName != "" {
					// NOTE: Every service within a group must expose the same port & protocol
					var portTemplate ServicePortTemplate
					portTemplate.Port = userChartGroup[2]
					portTemplate.Protocol = userChartGroup[3]
					serviceTemplate.Ports = append(serviceTemplate.Ports, portTemplate)

					c, err := createMeSvcChart(sandboxName, scenarioName, meSvcName, serviceTemplate.Ports)
					if err != nil {
						log.Debug("Failed to create ME Svc chart: ", err)
						return nil, err
					}
					if c != nil {
						charts = append(charts, *c)
						log.Debug("chart added for group service: ", meSvcName, " len:", len(charts))
					}
				}
			}
		} else if mepService := getMepService(proc); mepService != "" {
			log.Debug("Processing MEP Service chart for element[", proc.Name, "]")

			// Create Sandbox template
			var sandboxTemplate SandboxTemplate
			sandboxTemplate.SandboxName = sandboxName
			sandboxTemplate.Namespace = sandboxName
			sandboxTemplate.HostUrl = ve.hostUrl
			sandboxTemplate.HttpsOnly = ve.httpsOnly
			sandboxTemplate.AuthEnabled = ve.authEnabled
			sandboxTemplate.IsMepService = true

			// Get MEP Name
			mepName := ctx.Parents[mod.PhyLoc]
			sandboxTemplate.MepName = mepName

			// Create chart
			chartName := proc.Name
			chartLocation, _, err := createChart(chartName, sandboxName, scenarioName, mepService, sandboxTemplate)
			if err != nil {
				log.Debug("yaml creation file process: ", err)
				return nil, err
			}

			// validate if there is user value override
			userValueFile := "/user-values/" + mepName + "/" + mepService + ".yaml"
			if _, err := os.Stat(userValueFile); err != nil {
				// path/to/file does not exists
				// Note: according to https://helm.sh/docs/chart_template_guide/values_files/
				//       the order of precedence is: (lowest) default values.yaml
				//                                            then user value file
				//                                            then individual --set params (highest)
				//       Therefore, --set flags may interfere with user overrides
				userValueFile = ""
			}

			// Add chart to list
			c := newChart(chartName, sandboxName, scenarioName, chartLocation, userValueFile)
			charts = append(charts, c)
			log.Debug("chart added ", len(charts))

		} else {
			log.Debug("Processing virt-engine chart for element[", proc.Name, "]")

			// Fill deployment template information
			deploymentTemplate.Enabled = trueStr
			deploymentTemplate.ContainerName = proc.Name
			deploymentTemplate.ContainerImageRepository = proc.Image
			deploymentTemplate.ContainerImagePullPolicy = "Always"
			setEnv(deploymentTemplate, proc.Environment)
			setCommand(deploymentTemplate, proc.CommandExe, proc.CommandArguments)
			addMatchLabel(deploymentTemplate, "meepAppId: "+proc.Id)
			addTemplateLabel(deploymentTemplate, "meepAppId: "+proc.Id)
			deploymentTemplate.PlacementId = proc.PlacementId

			// Enable Service template if present
			if proc.ServiceConfig != nil {

				// Add app name associated to service
				svcName := proc.ServiceConfig.Name
				serviceTemplate.Enabled = trueStr
				serviceTemplate.Name = svcName
				serviceTemplate.Namespace = scenarioName
				addSelector(serviceTemplate, "meepSvc: "+svcName)
				addServiceLabel(serviceTemplate, "meepScenario: "+scenarioName)
				addTemplateLabel(deploymentTemplate, "meepSvc: "+svcName)

				// Add ports
				for _, ports := range proc.ServiceConfig.Ports {
					var portTemplate ServicePortTemplate
					portTemplate.Port = strconv.Itoa(int(ports.Port))
					portTemplate.TargetPort = strconv.Itoa(int(ports.Port))
					portTemplate.Protocol = ports.Protocol

					// Add NodePort if service is exposed externally
					if ports.ExternalPort >= serviceNodePortMin && ports.ExternalPort <= serviceNodePortMax {
						portTemplate.NodePort = strconv.Itoa(int(ports.ExternalPort))
						serviceTemplate.Type = "NodePort"
					} else {
						serviceTemplate.Type = "ClusterIP"
					}

					serviceTemplate.Ports = append(serviceTemplate.Ports, portTemplate)
				}

				// Create ME Service chart on first occurrence
				meSvcName := proc.ServiceConfig.MeSvcName
				if meSvcName != "" {
					c, err := createMeSvcChart(sandboxName, scenarioName, meSvcName, serviceTemplate.Ports)
					if err != nil {
						log.Debug("Failed to create ME Svc chart: ", err)
						return nil, err
					}
					if c != nil {
						charts = append(charts, *c)
						log.Debug("chart added for group service: ", meSvcName, " len:", len(charts))
					}

					// Add ME Svc service & pod labels
					addServiceLabel(serviceTemplate, "meepMeSvc: "+meSvcName)
					addTemplateLabel(deploymentTemplate, "meepMeSvc: "+meSvcName)
				}
			}

			// Enable GPU template if present
			if proc.GpuConfig != nil {
				deploymentTemplate.GpuEnabled = trueStr
				deploymentTemplate.GpuType = proc.GpuConfig.Type_
				deploymentTemplate.GpuCount = strconv.Itoa(int(proc.GpuConfig.Count))
			}

			// Enable CPU template if present
			if proc.CpuConfig != nil {
				deploymentTemplate.CpuEnabled = trueStr
				if proc.CpuConfig.Min != 0 {
					deploymentTemplate.CpuMin = strconv.FormatFloat(float64(proc.CpuConfig.Min), 'f', -1, 32)
				}
				if proc.CpuConfig.Max != 0 {
					deploymentTemplate.CpuMax = strconv.FormatFloat(float64(proc.CpuConfig.Max), 'f', -1, 32)
				}
			}

			// Enable Memory template if present
			if proc.MemoryConfig != nil {
				deploymentTemplate.MemoryEnabled = trueStr
				if proc.MemoryConfig.Min != 0 {
					deploymentTemplate.MemoryMin = strconv.Itoa(int(proc.MemoryConfig.Min)) + "Mi"
				}
				if proc.MemoryConfig.Max != 0 {
					deploymentTemplate.MemoryMax = strconv.Itoa(int(proc.MemoryConfig.Max)) + "Mi"
				}
			}

			// Enable External template if set
			if proc.IsExternal {
				externalTemplate.Enabled = trueStr
				addExtSelector(externalTemplate, "meepAppId: "+proc.Id)

				// Add ingress Service Maps, if any
				for _, svcMap := range proc.ExternalConfig.IngressServiceMap {
					var ingressSvcTemplate IngressServiceTemplate
					ingressSvcTemplate.NodePort = strconv.Itoa(int(svcMap.ExternalPort))
					ingressSvcTemplate.Port = strconv.Itoa(int(svcMap.Port))
					ingressSvcTemplate.Protocol = svcMap.Protocol
					ingressSvcTemplate.Name = "ingress-" + proc.Id + "-" + ingressSvcTemplate.NodePort

					externalTemplate.IngressServiceMap = append(externalTemplate.IngressServiceMap, ingressSvcTemplate)
				}

				// Add egress Service Maps, if any
				for _, svcMap := range proc.ExternalConfig.EgressServiceMap {
					var egressSvcTemplate EgressServiceTemplate
					egressSvcTemplate.Name = svcMap.Name
					egressSvcTemplate.IP = svcMap.Ip
					egressSvcTemplate.Port = strconv.Itoa(int(svcMap.Port))
					egressSvcTemplate.Protocol = svcMap.Protocol

					// Create and store ME Service template only with first occurrence.
					// If it already exists then add the matching pod label but don't create the service again.
					meSvcName := svcMap.MeSvcName
					if meSvcName != "" {
						if _, found := serviceMap[meSvcName]; !found {
							serviceMap[meSvcName] = "meepMeSvc: " + meSvcName
							egressSvcTemplate.MeSvcName = meSvcName
						}
					}

					externalTemplate.EgressServiceMap = append(externalTemplate.EgressServiceMap, egressSvcTemplate)
				}
			}

			// Create virt-engine chart
			chartName := proc.Name
			chartLocation, _, err := createChart(chartName, sandboxName, scenarioName, "", scenarioTemplate)
			if err != nil {
				log.Debug("yaml creation file process: ", err)
				return nil, err
			}
			c := newChart(chartName, sandboxName, scenarioName, chartLocation, "")
			charts = append(charts, c)
			log.Debug("chart added ", len(charts))
		}
	}

	return charts, nil
}

// Create ME Svc chart
func createMeSvcChart(sandboxName string, scenarioName string, meSvcName string, ports []ServicePortTemplate) (*helm.Chart, error) {

	// Create default scenario template
	var scenarioTemplate ScenarioTemplate
	serviceTemplate := &scenarioTemplate.Service
	setScenarioDefaults(&scenarioTemplate)

	// Fill general scenario template information
	scenarioTemplate.Namespace = scenarioName

	// Fill ME Svc template information
	serviceTemplate.MeServiceEnabled = trueStr
	serviceTemplate.MeServiceName = meSvcName
	serviceTemplate.Namespace = scenarioName
	serviceTemplate.Ports = ports
	addServiceLabel(serviceTemplate, "meepMeSvc: "+meSvcName)
	addServiceLabel(serviceTemplate, "meepScenario: "+scenarioName)

	// Create virt-engine chart for new group service
	chartName := "me-svc-" + meSvcName
	chartLocation, isNew, err := createChart(chartName, sandboxName, scenarioName, "", scenarioTemplate)
	if err != nil {
		log.Debug("yaml creation file process: ", err)
		return nil, err
	}
	if !isNew {
		log.Debug("Ignoring existing chart")
		return nil, nil
	}
	c := newChart(chartName, sandboxName, scenarioName, chartLocation, "")
	return &c, nil
}

func deployCharts(charts []helm.Chart, sandboxName string) error {
	err := helm.InstallCharts(charts, sandboxName)
	if err != nil {
		return err
	}
	return nil
}

func createChart(chartName, sandboxName, scenarioName, serviceName string, templateData interface{}) (outChart string, isNew bool, err error) {
	isNew = true

	// Determine source templates & destination chart location
	var templateChart string
	if scenarioName == "" && serviceName == "" {
		// Sandbox chart
		templateChart = "/templates/sandbox/" + chartName
		outChart = "/charts/" + sandboxName + "/sandbox/" + chartName
	} else if scenarioName != "" && serviceName == "" {
		// Scenario Chart
		templateChart = "/templates/scenario/meep-virt-chart-templates"
		outChart = "/charts/" + sandboxName + "/scenario/" + scenarioName + "/" + chartName
	} else if scenarioName != "" && serviceName != "" {
		// Service Chart
		templateChart = "/templates/sandbox/" + serviceName
		outChart = "/charts/" + sandboxName + "/scenario/" + scenarioName + "/" + chartName
	} else {
		return "", isNew, errors.New("Unsupported chart type")
	}
	templateValues := templateChart + "/values-template.yaml"
	outValues := outChart + "/values.yaml"

	// Create template object from template values file
	t, err := template.ParseFiles(templateValues)
	if err != nil {
		log.Error(err)
		return "", isNew, err
	}

	// Remove old chart if it already exists
	if _, err := os.Stat(outChart); err == nil {
		log.Debug("Removing old chart from path: ", outChart)
		os.RemoveAll(outChart)
		isNew = false
	}

	// Create new chart folder
	log.Debug("Creation of the output chart path: ", outChart)
	_ = CopyDir(templateChart, outChart)

	// Create new chart values file
	f, err := os.Create(outValues)
	if err != nil {
		log.Debug("create file: ", err)
		return "", isNew, err
	}

	// Fill new chart values file using template data
	err = t.Execute(f, templateData)
	if err != nil {
		log.Debug("execute: ", err)
		return "", isNew, err
	}

	f.Close()
	return outChart, isNew, nil
}

func newChart(chartName string, sandboxName string, scenarioName string, chartLocation string, valuesFile string) helm.Chart {
	var chart helm.Chart

	// Create release name by adding scenario prefix
	if scenarioName == "" {
		chart.ReleaseName = chartName
	} else {
		chart.ReleaseName = "meep-" + scenarioName + "-" + chartName
	}

	chart.Name = chartName
	chart.Namespace = sandboxName
	chart.Location = chartLocation
	chart.ValuesFile = valuesFile
	return chart
}

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

func getFullPath(path string) string {
	fullPath := path
	if path != "" && !strings.HasPrefix(path, "/") {
		fullPath = filepath.Join("/data/user-charts/", path)
	}
	return fullPath
}

func setScenarioDefaults(scenarioTemplate *ScenarioTemplate) {
	setDeploymentDefaults(&scenarioTemplate.Deployment)
	setServiceDefaults(&scenarioTemplate.Service)
	setExternalDefaults(&scenarioTemplate.External)
}

func setDeploymentDefaults(deploymentTemplate *DeploymentTemplate) {
	deploymentTemplate.Enabled = falseStr
	deploymentTemplate.ReplicaCount = "1"
	deploymentTemplate.ApiVersion = "v1"
	deploymentTemplate.ContainerEnvEnabled = falseStr
	deploymentTemplate.ContainerCommandEnabled = falseStr
	deploymentTemplate.GpuEnabled = falseStr
	deploymentTemplate.CpuEnabled = falseStr
	deploymentTemplate.MemoryEnabled = falseStr
	deploymentTemplate.CpuMin = ""
	deploymentTemplate.CpuMax = ""
	deploymentTemplate.MemoryMin = ""
	deploymentTemplate.MemoryMax = ""
}

func setServiceDefaults(serviceTemplate *ServiceTemplate) {
	serviceTemplate.Enabled = falseStr
	serviceTemplate.MeServiceEnabled = falseStr
}

func setExternalDefaults(externalTemplate *ExternalTemplate) {
	externalTemplate.Enabled = falseStr
}

func setEnv(deployment *DeploymentTemplate, envString string) {
	if envString != "" {
		deployment.ContainerEnvEnabled = trueStr
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
		deployment.ContainerCommandEnabled = trueStr

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

func generateSandboxCharts(sandboxName string) (charts []helm.Chart, err error) {

	// Create Sandbox template
	var sandboxTemplate SandboxTemplate
	sandboxTemplate.SandboxName = sandboxName
	sandboxTemplate.Namespace = sandboxName
	sandboxTemplate.HostUrl = ve.hostUrl
	sandboxTemplate.UserSwagger = ve.userSwagger
	sandboxTemplate.UserSwaggerDir = ve.userSwaggerDir
	sandboxTemplate.HttpsOnly = ve.httpsOnly
	sandboxTemplate.AuthEnabled = ve.authEnabled
	sandboxTemplate.IsMepService = false

	// Create sandbox charts
	for pod := range ve.sboxPods {
		var chartLocation string
		chartLocation, _, err = createChart(pod, sandboxName, "", "", sandboxTemplate)
		if err != nil {
			return
		}
		// validate if there is user value override
		userValueFile := "/user-values/" + pod + ".yaml"
		if _, err := os.Stat(userValueFile); err != nil {
			// path/to/file does not exists
			// Note: according to https://helm.sh/docs/chart_template_guide/values_files/
			//       the order of precedence is: (lowest) default values.yaml
			//                                            then user value file
			//                                            then individual --set params (highest)
			//       Therefore, --set flags may interfere with user overrides
			userValueFile = ""
		}

		chart := newChart(pod, sandboxName, "", chartLocation, userValueFile)
		charts = append(charts, chart)
	}

	return charts, nil
}

func deploySandbox(name string) error {

	// Create sandbox charts
	charts, err := generateSandboxCharts(name)
	if err != nil {
		log.Debug("Error creating sandbox charts: ", err)
		return err
	}
	log.Debug("Created ", len(charts), " sandbox charts")

	// Deploy all charts
	err = deployCharts(charts, name)
	if err != nil {
		log.Error("Error deploying charts: ", err)
		return err
	}

	return nil
}
