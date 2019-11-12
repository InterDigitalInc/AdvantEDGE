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

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const MAX_THROUGHPUT = 9999999999
const THROUGHPUT_UNIT = 1000000 //convert from Mbps to bps
const DEFAULT_THROUGHPUT_LINK = 1000.0

const metricsDb = 0
const moduleMetrics string = "metrics"

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
}

// SegAlgoPath -
type SegAlgoPath struct {
	Name     string
	Segments []*SegAlgoSegment
}

// SegAlgoNetElem -
type SegAlgoNetElem struct {
	Name              string
	Type              string
	PhyLocName        string
	PoaName           string
	ZoneName          string
	DomainName        string
	ConfiguredNetChar NetChar
}

// SegmentAlgorithm -
type SegmentAlgorithm struct {
	FlowMap    map[string]*SegAlgoFlow
	SegmentMap map[string]*SegAlgoSegment
	Config     SegAlgoConfig
	rc         *redis.Connector
}

// NewSegmentAlgorithm - Create, Initialize and connect
func NewSegmentAlgorithm(redisAddr string) (*SegmentAlgorithm, error) {
	// Create new instance & set default config
	var err error
	var algo SegmentAlgorithm
	algo.FlowMap = make(map[string]*SegAlgoFlow)
	algo.SegmentMap = make(map[string]*SegAlgoSegment)
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
	_ = algo.rc.DBFlush(moduleMetrics)
	log.Info("Connected to Metrics redis DB")

	return &algo, nil
}

// ProcessScenario -
func (algo *SegmentAlgorithm) ProcessScenario(model *mod.Model) error {
	var netElemList []SegAlgoNetElem

	// Process empty scenario
	if model.GetScenarioName() == "" {
		// Remove any existing metrics
		algo.deleteMetricsEntries()
		//reset the map
		algo.FlowMap = make(map[string]*SegAlgoFlow)
	}

	// Clear segment & flow maps
	algo.SegmentMap = make(map[string]*SegAlgoSegment)
	// Process active scenario
	procNames := model.GetNodeNames("CLOUD-APP")
	procNames = append(procNames, model.GetNodeNames("EDGE-APP")...)
	procNames = append(procNames, model.GetNodeNames("UE-APP")...)

	// Create NetElem for each scenario process
	for _, name := range procNames {
		// Retrieve node & context from model
		node := model.GetNode(name)
		if node == nil {
			err := errors.New("Error finding process: " + name)
			return err
		}
		proc, ok := node.(*ceModel.Process)
		if !ok {
			err := errors.New("Error casting process: " + name)
			return err
		}
		ctx := model.GetNodeContext(name)
		if ctx == nil {
			err := errors.New("Error getting context for process: " + name)
			return err
		}
		nodeCtx, ok := ctx.(*mod.NodeContext)
		if !ok {
			err := errors.New("Error casting context for process: " + name)
			return err
		}

		// Create & populate new element
		element := new(SegAlgoNetElem)
		element.Name = proc.Name
		element.PhyLocName = nodeCtx.Parents[mod.PhyLoc]
		element.DomainName = nodeCtx.Parents[mod.Domain]

		// Type-specific values
		element.Type = model.GetNodeType(element.PhyLocName)
		if element.Type == "UE" || element.Type == "FOG" {
			element.PoaName = nodeCtx.Parents[mod.NetLoc]
		}
		if element.Type != "DC" {
			element.ZoneName = nodeCtx.Parents[mod.Zone]
		}

		// Set max App Net chars (use default if set to 0)
		element.ConfiguredNetChar.Throughput = float64(proc.AppThroughput)
		if element.ConfiguredNetChar.Throughput == 0 {
			element.ConfiguredNetChar.Throughput = DEFAULT_THROUGHPUT_LINK
		}
		element.ConfiguredNetChar.Latency = float64(proc.AppLatency)
		element.ConfiguredNetChar.Jitter = float64(proc.AppLatencyVariation)
		element.ConfiguredNetChar.PacketLoss = float64(proc.AppPacketLoss)

		// Add element to list
		netElemList = append(netElemList, *element)
	}

	// Create all flows using Network Element list
	for _, elemSrc := range netElemList {
		for _, elemDest := range netElemList {
			if elemSrc.Name != elemDest.Name {
				// Create flow
				algo.populateFlow(elemSrc.Name+":"+elemDest.Name, &elemSrc, &elemDest, 0, model)

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
	keyName := moduleMetrics + ":*:throughput"
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

			flow.AllocatedThroughput = flow.MaxPlannedThroughput
			flow.AllocatedThroughputLowerBound = flow.MaxPlannedLowerBound
			flow.AllocatedThroughputUpperBound = flow.MaxPlannedUpperBound
			flow.AppliedNetChar.Throughput = flow.AllocatedThroughput
			updateNeeded = true
		}

		if (flow.ComputedLatency != flow.AppliedNetChar.Latency) ||
			(flow.ComputedJitter != flow.AppliedNetChar.Jitter) ||
			(flow.ComputedPacketLoss != flow.AppliedNetChar.PacketLoss) {
			if algo.Config.LogVerbose {
				log.Info("Update other netchars for ", flow.Name, " to ", flow.ComputedLatency, "-", flow.ComputedJitter, "-", flow.ComputedPacketLoss, " from ", flow.AppliedNetChar.Latency, "-", flow.AppliedNetChar.Jitter, "-", flow.AppliedNetChar.PacketLoss)
			}

			flow.AppliedNetChar.Latency = flow.ComputedLatency
			flow.AppliedNetChar.Jitter = flow.ComputedJitter
			flow.AppliedNetChar.PacketLoss = flow.ComputedPacketLoss
			updateNeeded = true
		}

		if updateNeeded {
			netchar := NetChar{flow.AppliedNetChar.Latency, flow.AppliedNetChar.Jitter, flow.AppliedNetChar.PacketLoss, flow.AppliedNetChar.Throughput}
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
			"meep.log.component": moduleName,
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
	_ = algo.rc.SetEntry(moduleMetrics+":"+dstElem+":"+srcElem, creationTime)
	_ = algo.rc.SetEntry(moduleMetrics+":"+dstElem+":throughput", creationTime)
}

// deleteMetricsEntries -
func (algo *SegmentAlgorithm) deleteMetricsEntries() {
	for _, flow := range algo.FlowMap {
		// Entries are created with no values, sidecar will only fill them, otherwise, won't be cleared
		_ = algo.rc.DelEntry(moduleMetrics + ":" + flow.DstNetElem + ":" + flow.SrcNetElem)
		_ = algo.rc.DelEntry(moduleMetrics + ":" + flow.DstNetElem + ":throughput")
	}
}

// populateFlow - Create/Update flow
func (algo *SegmentAlgorithm) populateFlow(flowName string, srcElement *SegAlgoNetElem, destElement *SegAlgoNetElem, maxBw float64, model *mod.Model) {

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
		if srcElement.ConfiguredNetChar.Throughput < destElement.ConfiguredNetChar.Throughput {
			maxBw = srcElement.ConfiguredNetChar.Throughput
		} else {
			maxBw = destElement.ConfiguredNetChar.Throughput
		}
	}
	flow.ConfiguredNetChar.Throughput = maxBw
	flow.ConfiguredNetChar.Latency = 0
	flow.ConfiguredNetChar.Jitter = 0
	flow.ConfiguredNetChar.PacketLoss = 0
	// Create a new path for this flow
	flow.Path = algo.createPath(flowName, srcElement, destElement, model)
}

// createPath -
func (algo *SegmentAlgorithm) createPath(flowName string, srcElement *SegAlgoNetElem, destElement *SegAlgoNetElem, model *mod.Model) *SegAlgoPath {

	direction := ""
	segmentName := ""
	var segment *SegAlgoSegment

	path := new(SegAlgoPath)
	path.Name = flowName

	//app segment ul, dl
	direction = "uplink"
	segmentName = srcElement.Name + direction
	segment = algo.createSegment(segmentName, flowName, srcElement.Name, model)
	path.Segments = append(path.Segments, segment)
	direction = "downlink"
	segmentName = destElement.Name + direction
	segment = algo.createSegment(segmentName, flowName, destElement.Name, model)
	path.Segments = append(path.Segments, segment)

	//node segment ul, dl
	direction = "uplink"
	segmentName = srcElement.PhyLocName + direction
	segment = algo.createSegment(segmentName, flowName, srcElement.PhyLocName, model)
	path.Segments = append(path.Segments, segment)
	direction = "downlink"
	segmentName = destElement.PhyLocName + direction
	segment = algo.createSegment(segmentName, flowName, destElement.PhyLocName, model)
	path.Segments = append(path.Segments, segment)

	//if on same node, return
	if srcElement.PhyLocName == destElement.PhyLocName {
		return path
	}

	//network location ul, dl
	if srcElement.Type == "UE" {
		direction = "uplink"
		segmentName = srcElement.PoaName + direction
		segment = algo.createSegment(segmentName, flowName, srcElement.PoaName, model)
		path.Segments = append(path.Segments, segment)
	}

	if destElement.Type == "UE" {
		direction = "downlink"
		segmentName = destElement.PoaName + direction
		segment = algo.createSegment(segmentName, flowName, destElement.PoaName, model)
		path.Segments = append(path.Segments, segment)
	}

	//if on same network location (poa), return
	if srcElement.PoaName == destElement.PoaName {
		return path
	}

	//zone ul, dl
	if srcElement.Type != "DC" {
		direction = "uplink"
		segmentName = srcElement.ZoneName + direction
		segment = algo.createSegment(segmentName, flowName, srcElement.ZoneName, model)
		path.Segments = append(path.Segments, segment)

	}

	if destElement.Type != "DC" {
		direction = "downlink"
		segmentName = destElement.ZoneName + direction
		segment = algo.createSegment(segmentName, flowName, destElement.ZoneName, model)
		path.Segments = append(path.Segments, segment)

	}

	//if in same zone, return
	if srcElement.ZoneName == destElement.ZoneName {
		return path
	}

	//domain ul, dl
	if srcElement.Type != "DC" {
		direction = "uplink"
		segmentName = srcElement.DomainName + direction
		segment = algo.createSegment(segmentName, flowName, srcElement.DomainName, model)
		path.Segments = append(path.Segments, segment)

	}

	if destElement.Type != "DC" {
		direction = "downlink"
		segmentName = destElement.DomainName + direction
		segment = algo.createSegment(segmentName, flowName, destElement.DomainName, model)
		path.Segments = append(path.Segments, segment)

	}

	//if in same domain, return
	if srcElement.DomainName == destElement.DomainName {
		return path
	}

	//cloud ul, dl
	if srcElement.Type == "DC" {
		direction = "uplink"
		segmentName = model.GetScenarioName() + "-cloud-" + direction
		segment = algo.createSegment(segmentName, flowName, model.GetScenarioName(), model)
		path.Segments = append(path.Segments, segment)

	}

	if destElement.Type == "DC" {
		direction = "downlink"
		segmentName = model.GetScenarioName() + "-cloud-" + direction
		segment = algo.createSegment(segmentName, flowName, model.GetScenarioName(), model)
		path.Segments = append(path.Segments, segment)

	}

	return path
}

/*

// createPath -
func (algo *SegmentAlgorithm) createPath(flowName string, srcElement *SegAlgoNetElem, destElement *SegAlgoNetElem, model *mod.Model) *SegAlgoPath {

	//Tier 1 -- check if they are in the same poa
	//Tier 2 -- check if they are in the same zone, but different poa
	//Tier 3 -- check if they are in the same domain, but different zone
	//Tier 4 -- check if they are in different domains

	direction := ""
	segmentName := ""
	var segment *SegAlgoSegment

	path := new(SegAlgoPath)
	path.Name = flowName

	//Tier 1
	if srcElement.PoaName != "" {
		//segments from element to POA
		//2 possibilities
		//UE->POA
		//segments for srcElement(app) -> 1A. UE-Node uplink-> 2A. POA-TermLink uplink
		//FOG-POA
		//segments for srcElement(app) -> 1B. FogNode uplink
		direction = "uplink"
		//segment 1A or 1B
		segmentName = srcElement.PhyLocName + "-" + direction
		segment = algo.createSegment(segmentName, flowName, srcElement.PhyLocName, model)
		path.Segments = append(path.Segments, segment)

		if srcElement.Type == "UE" {
			//segment 2A
			segmentName = srcElement.PoaName + "-" + direction
			segment = algo.createSegment(segmentName, flowName, srcElement.PoaName, model)
			path.Segments = append(path.Segments, segment)
		}
	}

	if destElement.PoaName != "" {
		//segments from POA to element
		//2 possibilities
		//POA->FOG
		//3A. Fog-Node downlink -> dstElement(app)
		//POA-UE
		//2B. POA-TermLink downlink -> 3B. UE-Node downlink -> dstElement(app)
		direction = "downlink"
		if destElement.Type == "UE" {
			//segment 2B
			segmentName = destElement.PoaName + "-" + direction
			segment = algo.createSegment(segmentName, flowName, destElement.PoaName, model)
			path.Segments = append(path.Segments, segment)
		}

		//segment 3A or 3B
		segmentName = destElement.PhyLocName + "-" + direction
		segment = algo.createSegment(segmentName, flowName, destElement.PhyLocName, model)
		path.Segments = append(path.Segments, segment)
	}

	//Tier 2
	//if same zone, different POA, OR no POA at all (Edge-Edge)
	if (srcElement.PoaName != destElement.PoaName) || (srcElement.PoaName == "" && destElement.PoaName == "") {
		//segments to intraZone backbone
		//2 possibilities
		//EDGE->IntraZoneBackbone
		//srcElement(app) -> 1A. Edge-Node uplink -> 2A. IntraZone uplink
		//POA->IntraZoneBackbone
		//2B. IntraZone uplink
		direction = "uplink"
		if srcElement.Type == "EDGE" {
			//segment 1A
			segmentName = srcElement.PhyLocName + "-" + direction
			segment = algo.createSegment(segmentName, flowName, srcElement.PhyLocName, model)
			path.Segments = append(path.Segments, segment)

			//segment 2A
			segmentName = srcElement.ZoneName + "-" + srcElement.PhyLocName + "-" + direction
		} else {
			//segment 2B
			segmentName = srcElement.ZoneName + "-" + srcElement.PoaName + "-" + direction
		}
		if srcElement.ZoneName != "" {
			segment = algo.createSegment(segmentName, flowName, srcElement.ZoneName, model)
			path.Segments = append(path.Segments, segment)
		}

		//segments from intraZone backbone
		//2 possibilities
		//IntraZoneBackbone->EDGE
		//3A. IntraZone downlink -> 4A. Edge-Node downlink -> srcElement(app)
		//IntraZoneBackbone->POA
		//3B. IntraZone downlink
		direction = "downlink"
		if destElement.Type == "EDGE" {
			//segment 4A
			segmentName = destElement.PhyLocName + "-" + direction
			segment = algo.createSegment(segmentName, flowName, destElement.PhyLocName, model)
			path.Segments = append(path.Segments, segment)

			//segment 3A
			segmentName = destElement.ZoneName + "-" + destElement.PhyLocName + "-" + direction
		} else {
			//segment 3B
			segmentName = destElement.ZoneName + "-" + destElement.PoaName + "-" + direction
		}
		if destElement.ZoneName != "" {
			segment = algo.createSegment(segmentName, flowName, destElement.ZoneName, model)
			path.Segments = append(path.Segments, segment)
		}
	}

	//Tier 3
	if srcElement.ZoneName != destElement.ZoneName {
		//segments to interZone backbone
		//1 possibility
		//Zone->InterZoneBackbone
		//1A. Zone uplink -> InterZone backbone (if zone exist)
		direction = "uplink"
		//segment 1A
		if srcElement.ZoneName != "" {
			segmentName = srcElement.ZoneName + "-" + direction
			segment = algo.createSegment(segmentName, flowName, srcElement.DomainName, model)
			path.Segments = append(path.Segments, segment)
		}

		//segments from interZone backbone
		//1 possibility
		//InterZoneBackbone->Zone
		//2A. InterZone backbone -> Zone downlink (if zone exist)
		direction = "downlink"
		//segment 2A
		if destElement.ZoneName != "" {
			segmentName = destElement.ZoneName + "-" + direction
			segment = algo.createSegment(segmentName, flowName, destElement.DomainName, model)
			path.Segments = append(path.Segments, segment)
		}
	}

	//Tier 4
	if srcElement.DomainName != destElement.DomainName {
		//segments to interDomain backbone
		//1 possibility
		//InterZoneBackbone->InterDomainBackbone
		//1A. InterZone backbone -> Domain backbone
		direction = "uplink"
		//segment 1A
		segmentName = srcElement.DomainName + "-" + direction
		segment = algo.createSegment(segmentName, flowName, model.GetScenarioName(), model)
		path.Segments = append(path.Segments, segment)

		//segments from interDomain backbone
		//1 possibility
		//InterDomainBackbone->InterZoneBackbone
		//2A. Domain backbone -> InterZone backbone
		direction = "downlink"
		//segment 2A
		segmentName = destElement.DomainName + "-" + direction
		segment = algo.createSegment(segmentName, flowName, model.GetScenarioName(), model)
		path.Segments = append(path.Segments, segment)

		//when going through interdomain, either from/to the cloud or another domain, if not cloud, already handled in other tiers sections
		if destElement.Type == "CLOUD" {
			segmentName = destElement.PhyLocName + "-" + direction
			segment = algo.createSegment(segmentName, flowName, destElement.PhyLocName, model)
			path.Segments = append(path.Segments, segment)
		} else if srcElement.Type == "CLOUD" {
			direction = "uplink"
			segmentName = srcElement.PhyLocName + "-" + direction
			segment = algo.createSegment(segmentName, flowName, srcElement.PhyLocName, model)
			path.Segments = append(path.Segments, segment)
		}
	}

	return path
}
*/

// createSegment -
func (algo *SegmentAlgorithm) createSegment(segmentName string, flowName string, elemName string, model *mod.Model) *SegAlgoSegment {
	// Create new segment if it does not exist
	segment := algo.SegmentMap[segmentName]
	if segment == nil {
		segment = new(SegAlgoSegment)
		segment.Name = segmentName

		// Retrieve max throughput from model using model scenario element name
		nc := getNetChars(elemName, model)
		segment.ConfiguredNetChar = nc
		maxThroughput := nc.Throughput
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
		flow := algo.FlowMap[trafficFrom+":"+subKey[1]]
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
				//first time it finds a value, it applies it directly
				flow.ComputedPacketLoss = segment.ConfiguredNetChar.PacketLoss
			} else {
				flow.ComputedPacketLoss += (flow.ComputedPacketLoss * (1 - segment.ConfiguredNetChar.PacketLoss))
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
		if flow.CurrentThroughput < flow.AllocatedThroughputLowerBound || flow.CurrentThroughput > flow.AllocatedThroughputUpperBound || flow.CurrentThroughput >= segment.MaxFairShareBwPerFlow {
			//			resetFlowMaxPlannedThroughput(flow)
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
func getNetChars(elemName string, model *mod.Model) (nc NetChar) {
	// Get Node
	node := model.GetNode(elemName)
	if node == nil {
		log.Error("Error finding element: " + elemName)
		return nc
	}

	maxThroughput := 0.0
	latency := 0.0
	jitter := 0.0
	packetLoss := 0.0
	// Get max throughput based on Node Type, as well as other netcharse
	if p, ok := node.(*ceModel.Process); ok {
		maxThroughput = float64(p.AppThroughput)
		latency = float64(p.AppLatency)
		jitter = float64(p.AppLatencyVariation)
		packetLoss = float64(p.AppPacketLoss)
	} else if pl, ok := node.(*ceModel.PhysicalLocation); ok {
		maxThroughput = float64(pl.LinkThroughput)
		latency = float64(pl.LinkLatency)
		jitter = float64(pl.LinkLatencyVariation)
		packetLoss = float64(pl.LinkPacketLoss)
	} else if nl, ok := node.(*ceModel.NetworkLocation); ok {
		maxThroughput = float64(nl.TerminalLinkThroughput)
		latency = float64(nl.TerminalLinkLatency)
		jitter = float64(nl.TerminalLinkLatencyVariation)
		packetLoss = float64(nl.TerminalLinkPacketLoss)
	} else if zone, ok := node.(*ceModel.Zone); ok {
		maxThroughput = float64(zone.EdgeFogThroughput)
		latency = float64(zone.EdgeFogLatency)
		jitter = float64(zone.EdgeFogLatencyVariation)
		packetLoss = float64(zone.EdgeFogPacketLoss)
	} else if domain, ok := node.(*ceModel.Domain); ok {
		maxThroughput = float64(domain.InterZoneThroughput)
		latency = float64(domain.InterZoneLatency)
		jitter = float64(domain.InterZoneLatencyVariation)
		packetLoss = float64(domain.InterZonePacketLoss)
	} else if deployment, ok := node.(*ceModel.Deployment); ok {
		maxThroughput = float64(deployment.InterDomainThroughput)
		latency = float64(deployment.InterDomainLatency)
		jitter = float64(deployment.InterDomainLatencyVariation)
		packetLoss = float64(deployment.InterDomainPacketLoss)
	} else {
		log.Error("Error casting element: " + elemName)
	}

	// For compatiblity reasons, set to default value if 0
	if maxThroughput == 0 {
		maxThroughput = DEFAULT_THROUGHPUT_LINK
	}

	nc.Throughput = maxThroughput
	nc.Latency = latency
	nc.Jitter = jitter
	nc.PacketLoss = packetLoss

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

	str := s1 + ": " + "Current: " + s6 + " - Configured: [" + s2t + "-" + s2l + "-" + s2j + "-" + s2p + "] Allocated: " + s3a + "[" + s4a + "-" + s5a + "]" + " - MaxPlanned: " + s3m + "[" + s4m + "-" + s5m + "]" + " - Planned: " + s3p + "[" + s4p + "-" + s5p + "] Computed Net Char: [" + s7l + "-" + s7j + "-" + s7p + "] Applied Net Char: [" + s8l + "-" + s8j + "-" + s8p + "]"
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

// printElement -
func printElement(elem *SegAlgoNetElem) string {
	str := elem.Name + "-" + elem.Type + "-" + elem.PhyLocName + "-" + elem.PoaName + "-" + elem.ZoneName + "-" + elem.DomainName
	return str
}
