/*
 * Copyright (c) 2019  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use algo file except in compliance with the License.
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

package netchar

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const MAX_THROUGHPUT = 9999999999
const THROUGHPUT_UNIT = 1000000 //convert from Mbps to bps
const DEFAULT_THROUGHPUT_LINK = 1000.0

const metricsDb = 0
const metricsKey string = "metrics:"

// SegAlgoConfig - Segment Algorithm Config
type SegAlgoConfig struct {
	// Segment config
	MaxBwPerInactiveFlow      float64
	MaxBwPerInactiveFlowFloor float64
	MinActivityThreshold      float64
	IncrementalStep           float64
	InactivityIncrementalStep float64
	TolerationThreshold       float64
	ActionUpperThreshold      float64

	// Debug Config
	IsPercentage bool
	LogVerbose   bool
}

// SegAlgoSegment -
type SegAlgoSegment struct {
	Name                      string
	ConfiguredNetChar         NetChar
	MaxFairShareBwPerFlow     float64
	CurrentThroughput         float64
	MaxBwPerInactiveFlow      float64
	MinActivityThreshold      float64
	IncrementalStep           float64
	InactivityIncrementalStep float64
	TolerationThreshold       float64
	ActionUpperThreshold      float64
	Flows                     []*SegAlgoFlow
}

// SegAlgoFlow -
type SegAlgoFlow struct {
	Name                          string
	SrcNetElem                    string
	DstNetElem                    string
	ConfiguredNetChar             NetChar
	AppliedNetChar                NetChar
	ComputedLatency               float64
	ComputedJitter                float64
	ComputedPacketLoss            float64
	AllocatedThroughput           float64 //allocated
	AllocatedThroughputLowerBound float64 //allocated
	AllocatedThroughputUpperBound float64 //allocated
	MaxPlannedThroughput          float64
	MaxPlannedLowerBound          float64
	MaxPlannedUpperBound          float64
	PlannedThroughput             float64
	PlannedLowerBound             float64
	PlannedUpperBound             float64
	CurrentThroughput             float64 //measured
	CurrentThroughputEgress       float64 //measured
	Path                          *SegAlgoPath
	UpdateRequired                bool
}

// SegAlgoPath -
type SegAlgoPath struct {
	Name         string
	Segments     []*SegAlgoSegment
	Disconnected bool
}

// SegAlgoNetElem -
type SegAlgoNetElem struct {
	Name              string
	Type              string
	PhyLocName        string
	PoaName           string
	ZoneName          string
	DomainName        string
	ConfiguredNetChar ElemNetChar
}

// SegmentAlgorithm -
type SegmentAlgorithm struct {
	Name              string
	Namespace         string
	BaseKey           string
	FlowMap           map[string]*SegAlgoFlow
	SegmentMap        map[string]*SegAlgoSegment
	ConnectivityModel string
	Config            SegAlgoConfig
	rc                *redis.Connector
}

// NewSegmentAlgorithm - Create, Initialize and connect
func NewSegmentAlgorithm(name string, namespace string, redisAddr string) (*SegmentAlgorithm, error) {
	// Create new instance & set default config
	var err error
	var algo SegmentAlgorithm
	algo.Name = name
	algo.Namespace = namespace
	algo.BaseKey = dkm.GetKeyRoot(namespace) + metricsKey
	algo.FlowMap = make(map[string]*SegAlgoFlow)
	algo.SegmentMap = make(map[string]*SegAlgoSegment)
	algo.ConnectivityModel = mod.ConnectivityModelOpen
	algo.Config.MaxBwPerInactiveFlow = 20.0
	algo.Config.MaxBwPerInactiveFlowFloor = 6.0
	algo.Config.MinActivityThreshold = 0.3
	algo.Config.IncrementalStep = 3.0
	algo.Config.InactivityIncrementalStep = 1.0
	algo.Config.ActionUpperThreshold = 1.0
	algo.Config.TolerationThreshold = 4.0
	algo.Config.IsPercentage = true

	// Create connection to Metrics Redis DB & flush entries
	algo.rc, err = redis.NewConnector(redisAddr, metricsDb)
	if err != nil {
		log.Error("Failed connection to Metrics redis DB. Error: ", err)
		return nil, err
	}
	_ = algo.rc.DBFlush(algo.BaseKey)
	log.Info("Connected to Metrics redis DB")

	return &algo, nil
}

// ProcessScenario -
func (algo *SegmentAlgorithm) ProcessScenario(model *mod.Model, pduSessions map[string]map[string]*dataModel.PduSessionInfo) error {
	var netElemList []SegAlgoNetElem

	// Process empty scenario
	if model.GetScenarioName() == "" {
		// Remove any existing metrics
		algo.deleteMetricsEntries()
		//reset the map
		algo.FlowMap = make(map[string]*SegAlgoFlow)
	}

	// Get scenario connectivity model
	algo.ConnectivityModel = model.GetConnectivityModel()

	// Clear segment & flow maps
	algo.SegmentMap = make(map[string]*SegAlgoSegment)
	// Process active scenario
	procNames := model.GetNodeNames("CLOUD-APP", "EDGE-APP", "UE-APP")

	// Create NetElem for each scenario process
	for _, name := range procNames {
		// Retrieve node & context from model
		node := model.GetNode(name)
		if node == nil {
			err := errors.New("Error finding process: " + name)
			return err
		}
		proc, ok := node.(*dataModel.Process)
		if !ok {
			err := errors.New("Error casting process: " + name)
			return err
		}
		ctx := model.GetNodeContext(name)
		if ctx == nil {
			err := errors.New("Error getting context for process: " + name)
			return err
		}

		// Create & populate new element
		element := new(SegAlgoNetElem)
		element.Name = proc.Name
		element.PhyLocName = ctx.Parents[mod.PhyLoc]
		element.DomainName = ctx.Parents[mod.Domain]

		// Type-specific values
		element.Type = model.GetNodeType(element.PhyLocName)
		if element.Type == "UE" || element.Type == "FOG" {
			element.PoaName = ctx.Parents[mod.NetLoc]
		}
		if element.Type != "DC" {
			element.ZoneName = ctx.Parents[mod.Zone]
		}

		deployment := model.GetNodeParent(element.DomainName).(*dataModel.Deployment)

		// Set max App Net chars (use default if set to 0)
		element.ConfiguredNetChar.Latency = float64(proc.NetChar.Latency)
		element.ConfiguredNetChar.Jitter = float64(proc.NetChar.LatencyVariation)
		element.ConfiguredNetChar.Distribution = deployment.NetChar.LatencyDistribution //set global value
		element.ConfiguredNetChar.ThroughputDl = float64(proc.NetChar.ThroughputUl)
		element.ConfiguredNetChar.ThroughputUl = float64(proc.NetChar.ThroughputUl)
		element.ConfiguredNetChar.PacketLoss = float64(proc.NetChar.PacketLoss)

		if element.ConfiguredNetChar.ThroughputUl == 0 {
			element.ConfiguredNetChar.ThroughputUl = DEFAULT_THROUGHPUT_LINK
		}
		if element.ConfiguredNetChar.ThroughputDl == 0 {
			element.ConfiguredNetChar.ThroughputDl = DEFAULT_THROUGHPUT_LINK
		}

		// Add element to list
		netElemList = append(netElemList, *element)
	}

	// Create all flows using Network Element list
	for _, elemSrc := range netElemList {
		for _, elemDest := range netElemList {
			if elemSrc.Name != elemDest.Name {
				// Create flow
				algo.populateFlow(elemSrc.Name+":"+elemDest.Name, &elemSrc, &elemDest, 0, model, pduSessions)

				// Create DB entry to begin collecting metrics for this flow
				algo.createMetricsEntry(elemSrc.Name, elemDest.Name)
			}
		}
	}

	// Log segments & flows in Verbose mode
	if algo.Config.LogVerbose {
		log.Info("Segments map: ", algo.SegmentMap)
		log.Info("Flows map: ", algo.FlowMap)
	}

	return nil
}

// CalculateNetChar - Run algorithm to recalculate network characteristics using latest scenario & metrics
func (algo *SegmentAlgorithm) CalculateNetChar() []FlowNetChar {
	var updatedNetCharList []FlowNetChar
	currentTime := time.Now()
	algo.logTimeLapse(&currentTime, "time to print")

	// Update flow with latest metrics
	keyName := algo.BaseKey + "*:throughput"
	err := algo.rc.ForEachEntry(keyName, algo.getMetricsThroughputEntryHandler, nil)
	if err != nil {
		log.Error("Failed to get entries: ", err)
		return updatedNetCharList
	}
	algo.logTimeLapse(&currentTime, "time to update metrics")

	// Recalculate segment BW allocation for each flow
	algo.reCalculateNetChar()
	algo.logTimeLapse(&currentTime, "time to recalculate throughput")

	// Prepare list of updated flows
	for _, flow := range algo.FlowMap {
		updateNeeded := false
		if flow.MaxPlannedThroughput != flow.AllocatedThroughput && flow.MaxPlannedThroughput != MAX_THROUGHPUT {
			if algo.Config.LogVerbose {
				log.Info("Update allocated bandwidth for ", flow.Name, " to ", flow.MaxPlannedThroughput, " from ", flow.AllocatedThroughput)
			}
			if flow.MaxPlannedThroughput >= 0 {
				flow.AllocatedThroughput = flow.MaxPlannedThroughput
				flow.AllocatedThroughputLowerBound = flow.MaxPlannedLowerBound
				flow.AllocatedThroughputUpperBound = flow.MaxPlannedUpperBound
				flow.AppliedNetChar.Throughput = flow.AllocatedThroughput
				updateNeeded = true
				if flow.MaxPlannedThroughput == 0 {
					log.Error("Impossible 0 result: ", printFlow(flow))
				}
			} else {
				log.Error("Impossible negative result: ", printFlow(flow))
			}
		}

		if (flow.ComputedLatency != flow.AppliedNetChar.Latency) ||
			(flow.ComputedJitter != flow.AppliedNetChar.Jitter) ||
			(flow.ComputedPacketLoss != flow.AppliedNetChar.PacketLoss) ||
			(flow.ConfiguredNetChar.Distribution != flow.AppliedNetChar.Distribution) {
			if algo.Config.LogVerbose {
				log.Info("Update other netchars for ", flow.Name, " to ", flow.ComputedLatency, "-", flow.ComputedJitter, "-", flow.ComputedPacketLoss, " from ", flow.AppliedNetChar.Latency, "-", flow.AppliedNetChar.Jitter, "-", flow.AppliedNetChar.PacketLoss, "-", flow.AppliedNetChar.Distribution)
			}

			flow.AppliedNetChar.Latency = flow.ComputedLatency
			flow.AppliedNetChar.Jitter = flow.ComputedJitter
			flow.AppliedNetChar.PacketLoss = flow.ComputedPacketLoss
			flow.AppliedNetChar.Distribution = flow.ConfiguredNetChar.Distribution
			updateNeeded = true
		}

		if updateNeeded {
			netchar := NetChar{flow.AppliedNetChar.Latency, flow.AppliedNetChar.Jitter, flow.AppliedNetChar.PacketLoss, flow.AppliedNetChar.Throughput, flow.ConfiguredNetChar.Distribution}
			flowNetChar := FlowNetChar{flow.SrcNetElem, flow.DstNetElem, netchar}
			updatedNetCharList = append(updatedNetCharList, flowNetChar)
		}
	}
	return updatedNetCharList
}

// SetConfigAttribute
func (algo *SegmentAlgorithm) SetConfigAttribute(fieldName string, fieldValue string) {
	switch fieldName {
	case "maxBwPerInactiveFlow":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			algo.Config.MaxBwPerInactiveFlow = value
		}
	case "maxBwPerInactiveFlowFloor":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			algo.Config.MaxBwPerInactiveFlowFloor = value
		}
	case "minActivityThreshold":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			algo.Config.MinActivityThreshold = value
		}
	case "incrementalStep":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			algo.Config.IncrementalStep = value
		}
	case "inactivityIncrementalStep":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			algo.Config.InactivityIncrementalStep = value
		}
	case "tolerationThreshold":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			algo.Config.TolerationThreshold = value
		}
	case "actionUpperThreshold":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			algo.Config.ActionUpperThreshold = value
		}
	case "isPercentage":
		if "yes" == fieldValue {
			algo.Config.IsPercentage = true
		} else {
			algo.Config.IsPercentage = false
		}
	case "logVerbose":
		if "yes" == fieldValue {
			algo.Config.LogVerbose = true
		}
	default:
	}
}

// logTimeLapse -
func (algo *SegmentAlgorithm) logTimeLapse(currentTime *time.Time, message string) {
	if algo.Config.LogVerbose {
		elapsed := time.Since(*currentTime)
		log.WithFields(log.Fields{
			"meep.log.component": algo.Name,
			"meep.time.location": message,
			"meep.time.exec":     elapsed,
		}).Info("Measurements log")
		*currentTime = time.Now()
	}
}

// createMetricsEntry -
func (algo *SegmentAlgorithm) createMetricsEntry(srcElem string, dstElem string) {
	var creationTime = make(map[string]interface{})
	creationTime["creationTime"] = time.Now()

	// Entries are created with no values, sidecar will only fill them, otherwise, won't be cleared
	_ = algo.rc.SetEntry(algo.BaseKey+dstElem+":throughput", creationTime)
}

// deleteMetricsEntries -
func (algo *SegmentAlgorithm) deleteMetricsEntries() {
	for _, flow := range algo.FlowMap {
		// Entries are created with no values, sidecar will only fill them, otherwise, won't be cleared
		_ = algo.rc.DelEntry(algo.BaseKey + flow.DstNetElem + ":throughput")
	}
}

// populateFlow - Create/Update flow
func (algo *SegmentAlgorithm) populateFlow(flowName string, srcElement *SegAlgoNetElem, destElement *SegAlgoNetElem, maxBw float64,
	model *mod.Model, pduSessions map[string]map[string]*dataModel.PduSessionInfo) {

	// Use existing flow if present or Create new flow
	flow := algo.FlowMap[flowName]
	if flow == nil {
		flow = new(SegAlgoFlow)
		flow.Name = flowName
		flow.SrcNetElem = srcElement.Name
		flow.DstNetElem = destElement.Name
		algo.FlowMap[flowName] = flow
	} else if flow.Name != flowName || flow.SrcNetElem != srcElement.Name && flow.DstNetElem != destElement.Name {
		log.Error("Flow already exists but not the same info, something is wrong!")
	}

	// Set maxBw to the minimum of the 2 ends if a max is not forced
	if maxBw == 0 {
		if srcElement.ConfiguredNetChar.ThroughputUl < destElement.ConfiguredNetChar.ThroughputDl {
			maxBw = srcElement.ConfiguredNetChar.ThroughputUl
		} else {
			maxBw = destElement.ConfiguredNetChar.ThroughputDl
		}
	}

	flow.ConfiguredNetChar.Throughput = maxBw
	//using distribution to pass it down, since it is global, they all have the same data at this point, so use any elements distribution
	flow.ConfiguredNetChar.Distribution = srcElement.ConfiguredNetChar.Distribution
	flow.ConfiguredNetChar.Latency = 0
	flow.ConfiguredNetChar.Jitter = 0
	flow.ConfiguredNetChar.PacketLoss = 0
	// Create a new path for this flow
	oldPath := flow.Path
	flow.Path = algo.createPath(flowName, srcElement, destElement, model, pduSessions)
	flow.UpdateRequired = algo.comparePath(oldPath, flow.Path)
}

func (algo *SegmentAlgorithm) comparePath(oldPath *SegAlgoPath, newPath *SegAlgoPath) bool {

	if oldPath == nil {
		return true
	}

	if len(oldPath.Segments) != len(newPath.Segments) {
		return true
	}

	for index, newSegment := range newPath.Segments {
		if newSegment.Name != oldPath.Segments[index].Name {
			return true
		}
	}
	return false
}

// createPath -
func (algo *SegmentAlgorithm) createPath(flowName string, srcElement *SegAlgoNetElem, destElement *SegAlgoNetElem,
	model *mod.Model, pduSessions map[string]map[string]*dataModel.PduSessionInfo) *SegAlgoPath {

	direction := ""
	var segment *SegAlgoSegment

	path := new(SegAlgoPath)
	path.Name = flowName
	path.Disconnected = false

	//app segment ul, dl
	direction = "uplink"
	segment = algo.createSegment(srcElement.Name, direction, flowName, model)
	path.Segments = append(path.Segments, segment)
	direction = "downlink"
	segment = algo.createSegment(destElement.Name, direction, flowName, model)
	path.Segments = append(path.Segments, segment)

	//node segment ul, dl
	direction = "uplink"
	segment = algo.createSegment(srcElement.PhyLocName, direction, flowName, model)
	path.Segments = append(path.Segments, segment)
	direction = "downlink"
	segment = algo.createSegment(destElement.PhyLocName, direction, flowName, model)
	path.Segments = append(path.Segments, segment)

	//if on same node, return
	if srcElement.PhyLocName == destElement.PhyLocName {
		return path
	}

	// Check if src or dest Physical location is disconnected
	// NOTE: Does not apply to apps on same physical node
	var srcPhyLoc *dataModel.PhysicalLocation
	srcPhyLocNode := model.GetNode(srcElement.PhyLocName)
	if srcPhyLocNode != nil {
		var ok bool
		if srcPhyLoc, ok = srcPhyLocNode.(*dataModel.PhysicalLocation); ok {
			path.Disconnected = path.Disconnected || !srcPhyLoc.Connected
		}
	}
	var destPhyLoc *dataModel.PhysicalLocation
	destPhyLocNode := model.GetNode(destElement.PhyLocName)
	if destPhyLocNode != nil {
		var ok bool
		if destPhyLoc, ok = destPhyLocNode.(*dataModel.PhysicalLocation); ok {
			path.Disconnected = path.Disconnected || !destPhyLoc.Connected
		}
	}

	// If in PDU Connectivity mode, check if src or dest UE has PDU connectivity to DN
	// NOTE: For LADN, additionally verify that UE and edge app are in the same zone
	if !path.Disconnected && algo.ConnectivityModel == mod.ConnectivityModelPdu {
		if mod.IsUe(srcPhyLoc.Type_) {
			pduMap, ok := pduSessions[srcPhyLoc.Name]
			if !ok || mod.IsUe(destPhyLoc.Type_) || destPhyLoc.DataNetwork == nil {
				path.Disconnected = true
			} else if destPhyLoc.DataNetwork.Ladn && srcElement.ZoneName != destElement.ZoneName {
				// LADN & not in same zone
				path.Disconnected = true
			} else {
				// Find matching DNN
				var found bool
				for _, pdu := range pduMap {
					if pdu.Dnn == destPhyLoc.DataNetwork.Dnn {
						found = true
						break
					}
				}
				if !found {
					path.Disconnected = true
				}
			}
		}

		if mod.IsUe(destPhyLoc.Type_) {
			pduMap, ok := pduSessions[destPhyLoc.Name]
			if !ok || mod.IsUe(srcPhyLoc.Type_) || srcPhyLoc.DataNetwork == nil {
				path.Disconnected = true
			} else if srcPhyLoc.DataNetwork.Ladn && srcElement.ZoneName != destElement.ZoneName {
				// LADN & not in same zone
				path.Disconnected = true
			} else {
				// Find matching DNN
				var found bool
				for _, pdu := range pduMap {
					if pdu.Dnn == srcPhyLoc.DataNetwork.Dnn {
						found = true
						break
					}
				}
				if !found {
					path.Disconnected = true
				}
			}
		}
	}

	//network location ul, dl
	if srcElement.Type == "UE" {
		direction = "uplink"
		segment = algo.createSegment(srcElement.PoaName, direction, flowName, model)
		path.Segments = append(path.Segments, segment)
	}

	if destElement.Type == "UE" {
		direction = "downlink"
		segment = algo.createSegment(destElement.PoaName, direction, flowName, model)
		path.Segments = append(path.Segments, segment)
	}

	//if on same network location (poa), return
	if srcElement.PoaName == destElement.PoaName {
		return path
	}

	//zone ul, dl
	if srcElement.Type != "DC" {
		direction = "uplink"
		segment = algo.createSegment(srcElement.ZoneName, direction, flowName, model)
		path.Segments = append(path.Segments, segment)

	}

	if destElement.Type != "DC" {
		direction = "downlink"
		segment = algo.createSegment(destElement.ZoneName, direction, flowName, model)
		path.Segments = append(path.Segments, segment)

	}

	//if in same zone, return
	if srcElement.ZoneName == destElement.ZoneName {
		return path
	}

	//domain ul, dl
	if srcElement.Type != "DC" {
		direction = "uplink"
		segment = algo.createSegment(srcElement.DomainName, direction, flowName, model)
		path.Segments = append(path.Segments, segment)

	}

	if destElement.Type != "DC" {
		direction = "downlink"
		segment = algo.createSegment(destElement.DomainName, direction, flowName, model)
		path.Segments = append(path.Segments, segment)

	}

	//if in same domain, return
	if srcElement.DomainName == destElement.DomainName {
		return path
	}

	//interdomain

	//computing twice while in the interdomain
	direction = "uplink"
	segment = algo.createSegment(model.GetScenarioName(), direction, flowName, model)
	path.Segments = append(path.Segments, segment)

	direction = "downlink"
	segment = algo.createSegment(model.GetScenarioName(), direction, flowName, model)
	path.Segments = append(path.Segments, segment)

	return path
}

// createSegment -
func (algo *SegmentAlgorithm) createSegment(elemName string, direction string, flowName string, model *mod.Model) *SegAlgoSegment {
	// Create new segment if it does not exist
	segmentName := elemName + direction
	segment := algo.SegmentMap[segmentName]
	if segment == nil {
		segment = new(SegAlgoSegment)
		segment.Name = segmentName

		// Retrieve max throughput from model using model scenario element name
		nc := getNetChars(elemName, model)
		ncThroughput := 0.0
		if direction == "uplink" {
			ncThroughput = float64(nc.ThroughputUl)
		} else {
			ncThroughput = float64(nc.ThroughputDl)
		}

		segment.ConfiguredNetChar.Latency = float64(nc.Latency)
		segment.ConfiguredNetChar.Jitter = float64(nc.LatencyVariation)
		segment.ConfiguredNetChar.PacketLoss = float64(nc.PacketLoss)
		segment.ConfiguredNetChar.Throughput = float64(ncThroughput)

		maxThroughput := ncThroughput
		// Initialize segment-specific BW attributes from Algo config
		if algo.Config.IsPercentage {
			segment.MaxBwPerInactiveFlow = algo.Config.MaxBwPerInactiveFlow * maxThroughput / 100
			if segment.MaxBwPerInactiveFlow < algo.Config.MaxBwPerInactiveFlowFloor {
				segment.MaxBwPerInactiveFlow = algo.Config.MaxBwPerInactiveFlowFloor
			}
			segment.MinActivityThreshold = algo.Config.MinActivityThreshold * maxThroughput / 100
			segment.IncrementalStep = algo.Config.IncrementalStep * maxThroughput / 100
			segment.InactivityIncrementalStep = algo.Config.InactivityIncrementalStep * maxThroughput / 100
			segment.TolerationThreshold = algo.Config.TolerationThreshold * maxThroughput / 100
			segment.ActionUpperThreshold = algo.Config.ActionUpperThreshold * maxThroughput / 100
		} else {
			segment.MaxBwPerInactiveFlow = algo.Config.MaxBwPerInactiveFlow
			segment.MinActivityThreshold = algo.Config.MinActivityThreshold
			segment.IncrementalStep = algo.Config.IncrementalStep
			segment.InactivityIncrementalStep = algo.Config.InactivityIncrementalStep
			segment.TolerationThreshold = algo.Config.TolerationThreshold
			segment.ActionUpperThreshold = algo.Config.ActionUpperThreshold
		}

		// Add segment to map
		algo.SegmentMap[segmentName] = segment
	}

	// Add flow to segment flow map
	flow := algo.FlowMap[flowName]
	if flow != nil {
		segment.Flows = append(segment.Flows, flow)
	} else {
		log.Error("Missing flow: ", flowName)
	}

	return segment
}

// getMetricsThroughputEntryHandler -
func (algo *SegmentAlgorithm) getMetricsThroughputEntryHandler(key string, fields map[string]string, userData interface{}) error {
	subKey := strings.Split(key, ":")
	for trafficFrom, throughput := range fields {
		flow := algo.FlowMap[trafficFrom+":"+subKey[len(subKey)-2]]
		if flow != nil {
			value, _ := strconv.ParseFloat(throughput, 64)
			flow.CurrentThroughput = value
		}
	}
	return nil
}

// reCalculateNetChar -
func (algo *SegmentAlgorithm) reCalculateNetChar() {
	//reset every planned throughput values for every flow since they will start to populate those
	for _, flow := range algo.FlowMap {
		resetComputedNetChar(flow)

		// For flows passing through a disconnected Physical location, set Packet loss to 100%
		if flow.Path != nil && flow.Path.Disconnected {
			flow.ComputedPacketLoss = 100
		}
	}

	//all segments determined by the scenario
	for _, segment := range algo.SegmentMap {

		//throughput specific
		updateMaxFairShareBwPerFlow(segment)
		unusedBw, list := needToReevaluate(segment)

		if list != nil {
			if algo.Config.LogVerbose {
				log.Info("Segment ", segment.Name, " reevaluation result - BW unused: ", unusedBw, "***Flows to evaluate***: ", printFlowNamesFromList(list))
			}

			recalculateSegmentBw(segment, list, unusedBw)
		}

		//latency, jitter, packet-loss computation for each flow in each segment
		for _, flow := range segment.Flows {
			flow.ComputedLatency += segment.ConfiguredNetChar.Latency
			flow.ComputedJitter += segment.ConfiguredNetChar.Jitter
			if flow.ComputedPacketLoss == 0 {
				flow.ComputedPacketLoss = segment.ConfiguredNetChar.PacketLoss
			} else if segment.ConfiguredNetChar.PacketLoss != 0 {
				flow.ComputedPacketLoss += (segment.ConfiguredNetChar.PacketLoss * ((100 - flow.ComputedPacketLoss) / 100))
			}
		}
		if algo.Config.LogVerbose {
			printFlows(segment)
		}
	}
}

// resetComputedNetChar -
func resetComputedNetChar(flow *SegAlgoFlow) {
	flow.MaxPlannedThroughput = MAX_THROUGHPUT
	flow.MaxPlannedLowerBound = MAX_THROUGHPUT
	flow.MaxPlannedUpperBound = MAX_THROUGHPUT
	flow.ComputedLatency = 0
	flow.ComputedJitter = 0
	flow.ComputedPacketLoss = 0
}

// recalculateSegmentBw -
func recalculateSegmentBw(segment *SegAlgoSegment, flowsToEvaluate []*SegAlgoFlow, unusedBw float64) {
	nbEvaluatedflowsLeft := len(flowsToEvaluate)
	if segment.CurrentThroughput > segment.ConfiguredNetChar.Throughput || nbEvaluatedflowsLeft >= 1 {

		//category 1 Flows
		for _, flow := range flowsToEvaluate {
			if flow.CurrentThroughput+segment.IncrementalStep > segment.MaxFairShareBwPerFlow {
				flow.PlannedThroughput = segment.MaxFairShareBwPerFlow //category 2 or 3
			} else {
				if flow.CurrentThroughput <= segment.MinActivityThreshold {
					flow.PlannedThroughput = segment.MaxBwPerInactiveFlow
					flow.PlannedUpperBound = segment.InactivityIncrementalStep
					flow.PlannedLowerBound = 0
				} else {
					flow.PlannedThroughput = flow.CurrentThroughput + segment.IncrementalStep
					if flow.PlannedThroughput > flow.ConfiguredNetChar.Throughput {
						flow.PlannedThroughput = flow.ConfiguredNetChar.Throughput
					}
					flow.PlannedUpperBound = flow.PlannedThroughput - segment.ActionUpperThreshold
					flow.PlannedLowerBound = flow.PlannedUpperBound - segment.TolerationThreshold
					//lower bound cannot be less than min threshold
					if flow.PlannedLowerBound < segment.MinActivityThreshold {
						flow.PlannedLowerBound = segment.MinActivityThreshold
					}
				}
				nbEvaluatedflowsLeft--
				if flow.PlannedThroughput != segment.MaxBwPerInactiveFlow {
					unusedBw -= flow.PlannedThroughput
				}
			}
		}

		var extra float64

		if nbEvaluatedflowsLeft > 0 {

			//category 2 Flows
			for _, flow := range flowsToEvaluate {
				if flow.PlannedThroughput == segment.MaxFairShareBwPerFlow {
					if flow.CurrentThroughput < segment.MaxFairShareBwPerFlow {
						nbEvaluatedflowsLeft--
						if nbEvaluatedflowsLeft == 0 { //allocate everything of what is left
							flow.PlannedThroughput = unusedBw
							if flow.PlannedThroughput > flow.ConfiguredNetChar.Throughput {
								flow.PlannedThroughput = flow.ConfiguredNetChar.Throughput
							}
							flow.PlannedUpperBound = flow.PlannedThroughput
							flow.PlannedLowerBound = flow.PlannedThroughput - segment.TolerationThreshold
							//lower bound cannot be less than min threshold
							if flow.PlannedLowerBound < segment.MinActivityThreshold {
								flow.PlannedLowerBound = segment.MinActivityThreshold
							}
						} else {
							flow.PlannedThroughput = flow.CurrentThroughput + segment.IncrementalStep
							if flow.PlannedThroughput > flow.ConfiguredNetChar.Throughput {
								flow.PlannedThroughput = flow.ConfiguredNetChar.Throughput
							}
							flow.PlannedUpperBound = flow.PlannedThroughput - segment.ActionUpperThreshold
							flow.PlannedLowerBound = flow.PlannedUpperBound - segment.TolerationThreshold
							//lower bound cannot be less than min threshold
							if flow.PlannedLowerBound < segment.MinActivityThreshold {
								flow.PlannedLowerBound = segment.MinActivityThreshold
							}
						}
						unusedBw -= flow.PlannedThroughput
					}
				}
			}

			if nbEvaluatedflowsLeft > 0 {
				if nbEvaluatedflowsLeft >= 1 {
					extra = (unusedBw - float64(nbEvaluatedflowsLeft)*segment.MaxFairShareBwPerFlow) / float64(nbEvaluatedflowsLeft)
				} else {
					extra = 0
				}

				//category 3
				for _, flow := range flowsToEvaluate {
					if flow.PlannedThroughput == segment.MaxFairShareBwPerFlow && flow.CurrentThroughput >= segment.MaxFairShareBwPerFlow {
						flow.PlannedThroughput = segment.MaxFairShareBwPerFlow + extra
						if flow.PlannedThroughput > flow.ConfiguredNetChar.Throughput {
							flow.PlannedThroughput = flow.ConfiguredNetChar.Throughput
						}
						flow.PlannedUpperBound = flow.PlannedThroughput - segment.ActionUpperThreshold
						flow.PlannedLowerBound = flow.PlannedUpperBound - segment.TolerationThreshold
						unusedBw -= flow.PlannedThroughput
					}
				}
			}
		}
	}
	//we allocate all the bw to active users and very low values to inactive ones if there is any residual
	//using a minimum value that is close but not exactly 0, since we use float operations and approximation may not lead to a perfect
	if unusedBw >= 1 {
		for _, flow := range flowsToEvaluate {
			if flow.CurrentThroughput > segment.MinActivityThreshold {
				flow.PlannedThroughput = segment.MaxFairShareBwPerFlow
				if flow.PlannedThroughput > flow.ConfiguredNetChar.Throughput {
					flow.PlannedThroughput = flow.ConfiguredNetChar.Throughput
				}
				flow.PlannedLowerBound = 0
				flow.PlannedUpperBound = 0
			}
		}
	}

	//update or not the throughput
	for _, flow := range flowsToEvaluate {
		if flow.PlannedThroughput < flow.MaxPlannedThroughput {
			if flow.PlannedThroughput <= 0 {
				log.Error("Max : ", flow.PlannedThroughput, "---", flow.MaxPlannedThroughput)
			}
			flow.MaxPlannedThroughput = flow.PlannedThroughput
			flow.MaxPlannedLowerBound = flow.PlannedLowerBound
			flow.MaxPlannedUpperBound = flow.PlannedUpperBound
		}
	}

}

// needToReevaluate - determines which Flows must be recalculated for bandwidth sharing within the segment
func needToReevaluate(segment *SegAlgoSegment) (unusedBw float64, list []*SegAlgoFlow) {
	unusedBw = segment.ConfiguredNetChar.Throughput

	//how many active connections that needs to be taken into account
	for _, flow := range segment.Flows {
		if flow.CurrentThroughput < flow.AllocatedThroughputLowerBound ||
			flow.CurrentThroughput > flow.AllocatedThroughputUpperBound ||
			flow.CurrentThroughput >= segment.MaxFairShareBwPerFlow ||
			flow.UpdateRequired {
			list = append(list, flow)
		} else {
			//no need to reevalute algo one, so removing its allocated bw from the available one
			if flow.CurrentThroughput >= segment.MinActivityThreshold {
				unusedBw -= flow.AllocatedThroughput
			}
		}
	}
	return unusedBw, list
}

// updateMaxFairShareBwPerFlow -
func updateMaxFairShareBwPerFlow(segment *SegAlgoSegment) {
	nbActiveConnections := 0
	for _, flow := range segment.Flows {
		if flow.CurrentThroughput >= segment.MinActivityThreshold {
			nbActiveConnections++
		}
	}
	if nbActiveConnections >= 1 {
		segment.MaxFairShareBwPerFlow = segment.ConfiguredNetChar.Throughput / float64(nbActiveConnections)
	} else {
		segment.MaxFairShareBwPerFlow = MAX_THROUGHPUT
	}
}

// getNetChars - Retrieve all network characteristics from model for provided element name
func getNetChars(elemName string, model *mod.Model) *dataModel.NetworkCharacteristics {

	nc := new(dataModel.NetworkCharacteristics)

	// Get Node
	node := model.GetNode(elemName)
	if node == nil {
		log.Error("Error finding element: " + elemName)
		return nc
	}

	// Get max throughput based on Node Type, as well as other netcharse
	if p, ok := node.(*dataModel.Process); ok {
		nc = p.NetChar
	} else if pl, ok := node.(*dataModel.PhysicalLocation); ok {
		nc = pl.NetChar
	} else if nl, ok := node.(*dataModel.NetworkLocation); ok {
		nc = nl.NetChar
	} else if zone, ok := node.(*dataModel.Zone); ok {
		nc = zone.NetChar
	} else if domain, ok := node.(*dataModel.Domain); ok {
		nc = domain.NetChar
	} else if deployment, ok := node.(*dataModel.Deployment); ok {
		nc = deployment.NetChar
	} else {
		log.Error("Error casting element: " + elemName)
	}

	// For compatiblity reasons, set to default value if 0
	if nc.ThroughputUl == 0 {
		nc.ThroughputUl = DEFAULT_THROUGHPUT_LINK
	}
	if nc.ThroughputDl == 0 {
		nc.ThroughputDl = DEFAULT_THROUGHPUT_LINK
	}

	return nc
}

// printFlowNamesFromList -
func printFlowNamesFromList(list []*SegAlgoFlow) string {
	str := ""
	for _, flow := range list {
		str += flow.Name + "."
	}
	return str
}

// printFlows -
func printFlows(segment *SegAlgoSegment) {
	log.Info("Flows on segment ", segment.Name)
	for _, flow := range segment.Flows {
		log.Info(printFlow(flow))
	}
}

// printFlow -
func printFlow(flow *SegAlgoFlow) string {
	s0 := fmt.Sprintf("%x", &flow)
	s1 := flow.Name + "(" + s0 + ")"
	s2t := fmt.Sprintf("%f", flow.ConfiguredNetChar.Throughput)
	s2l := fmt.Sprintf("%f", flow.ConfiguredNetChar.Latency)
	s2j := fmt.Sprintf("%f", flow.ConfiguredNetChar.Jitter)
	s2d := flow.ConfiguredNetChar.Distribution
	s2p := fmt.Sprintf("%f", flow.ConfiguredNetChar.PacketLoss)
	s3a := fmt.Sprintf("%f", flow.AllocatedThroughput)
	s4a := fmt.Sprintf("%f", flow.AllocatedThroughputLowerBound)
	s5a := fmt.Sprintf("%f", flow.AllocatedThroughputUpperBound)
	s3m := fmt.Sprintf("%f", flow.MaxPlannedThroughput)
	s4m := fmt.Sprintf("%f", flow.MaxPlannedLowerBound)
	s5m := fmt.Sprintf("%f", flow.MaxPlannedUpperBound)
	s3p := fmt.Sprintf("%f", flow.PlannedThroughput)
	s4p := fmt.Sprintf("%f", flow.PlannedLowerBound)
	s5p := fmt.Sprintf("%f", flow.PlannedUpperBound)
	s6 := fmt.Sprintf("%f", flow.CurrentThroughput)
	s7l := fmt.Sprintf("%f", flow.ComputedLatency)
	s7j := fmt.Sprintf("%f", flow.ComputedJitter)
	s7p := fmt.Sprintf("%f", flow.ComputedPacketLoss)
	s8l := fmt.Sprintf("%f", flow.AppliedNetChar.Latency)
	s8j := fmt.Sprintf("%f", flow.AppliedNetChar.Jitter)
	s8p := fmt.Sprintf("%f", flow.AppliedNetChar.PacketLoss)
	s8d := flow.AppliedNetChar.Distribution

	str := s1 + ": " + "Current: " + s6 + " - Configured: [" + s2t + "-" + s2l + "-" + s2j + "-" + s2p + "-" + s2d + "] Allocated: " + s3a + "[" + s4a + "-" + s5a + "]" + " - MaxPlanned: " + s3m + "[" + s4m + "-" + s5m + "]" + " - Planned: " + s3p + "[" + s4p + "-" + s5p + "] Computed Net Char: [" + s7l + "-" + s7j + "-" + s7p + "] Applied Net Char: [" + s8l + "-" + s8j + "-" + s8p + "-" + s8d + "]"
	str += printPath(flow.Path)
	return str
}

// printPath -
func printPath(path *SegAlgoPath) string {
	str := ""
	first := true
	if path != nil {
		str = "Path: "
		for _, segment := range path.Segments {
			if first {
				str += segment.Name
				first = false
			} else {
				str += "..." + segment.Name
			}
		}
	}
	return str
}

// // printElement -
// func printElement(elem *SegAlgoNetElem) string {
// 	str := elem.Name + "-" + elem.Type + "-" + elem.PhyLocName + "-" + elem.PoaName + "-" + elem.ZoneName + "-" + elem.DomainName
// 	return str
// }
