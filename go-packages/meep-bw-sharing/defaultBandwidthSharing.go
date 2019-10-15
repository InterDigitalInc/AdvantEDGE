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

package bws

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

// DebugConfiguration -
type DebugConfiguration struct {
	IsPercentage bool
	LogVerbose   bool
}

// SegmentConfiguration -
type SegmentConfiguration struct {
	MaxBwPerInactiveFlow      float64
	MinActivityThreshold      float64
	IncrementalStep           float64
	InactivityIncrementalStep float64
	TolerationThreshold       float64
	ActionUpperThreshold      float64
}

// BandwidthSharingSegment -
type BandwidthSharingSegment struct {
	Name                  string
	MaxThroughput         float64
	MaxFairShareBwPerFlow float64
	CurrentThroughput     float64
	Config                SegmentConfiguration
	Flows                 []*BandwidthSharingFlow
}

// BandwidthSharingFlow -
type BandwidthSharingFlow struct {
	Name                          string
	SrcNetworkElement             string
	DstNetworkElement             string
	MaximumThroughput             float64 //config
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
	Path                          *Path
}

// Path -
type Path struct {
	Name     string
	Segments []*BandwidthSharingSegment
}

// NetElem -
type NetElem struct {
	Name                     string
	Type                     string
	PhyLocName               string
	PoaName                  string
	ZoneName                 string
	DomainName               string
	MaxThroughput            float64
	PhyLocMaxThroughput      float64
	PoaMaxThroughput         float64
	IntraZoneMaxThroughput   float64
	InterZoneMaxThroughput   float64
	InterDomainMaxThroughput float64
}

// BwSharing -
type DefaultBwSharingAlgorithm struct {
	ParentBwSharing         *BwSharing
	BandwidthSharingFlowMap map[string]*BandwidthSharingFlow
	SegmentsMap             map[string]*BandwidthSharingSegment
	DefaultConfigs          SegmentConfiguration
	DefaultDebugConfigs     DebugConfiguration
}

// allocateBandwidthSharing - allocated structures
func (this *DefaultBwSharingAlgorithm) allocateBandwidthSharing() {
	this.BandwidthSharingFlowMap = make(map[string]*BandwidthSharingFlow)
	this.SegmentsMap = make(map[string]*BandwidthSharingSegment)
}

func (this *DefaultBwSharingAlgorithm) setParentBwSharing(bwSharing *BwSharing) {
	this.ParentBwSharing = bwSharing
}

// getBandwidthSharingFlow -
func (this *DefaultBwSharingAlgorithm) getBandwidthSharingFlow(key string) *BandwidthSharingFlow {
	return this.BandwidthSharingFlowMap[key]
}

// getBandwidthSharingSegment -
func (this *DefaultBwSharingAlgorithm) getBandwidthSharingSegment(key string) *BandwidthSharingSegment {
	return this.SegmentsMap[key]
}

// setSegmentConfiguration - set configuration attributes on a per segment basis (absolute or relative %)
func (this *DefaultBwSharingAlgorithm) setSegmentConfiguration(toSegment *SegmentConfiguration, fromSegment *SegmentConfiguration, maxThroughput float64) {

	if this.DefaultDebugConfigs.IsPercentage {
		toSegment.MaxBwPerInactiveFlow = fromSegment.MaxBwPerInactiveFlow * maxThroughput / 100
		toSegment.MinActivityThreshold = fromSegment.MinActivityThreshold * maxThroughput / 100
		toSegment.IncrementalStep = fromSegment.IncrementalStep * maxThroughput / 100
		toSegment.InactivityIncrementalStep = fromSegment.InactivityIncrementalStep * maxThroughput / 100
		toSegment.TolerationThreshold = fromSegment.TolerationThreshold * maxThroughput / 100
		toSegment.ActionUpperThreshold = fromSegment.ActionUpperThreshold * maxThroughput / 100
	} else {
		toSegment.MaxBwPerInactiveFlow = fromSegment.MaxBwPerInactiveFlow
		toSegment.MinActivityThreshold = fromSegment.MinActivityThreshold
		toSegment.IncrementalStep = fromSegment.IncrementalStep
		toSegment.InactivityIncrementalStep = fromSegment.InactivityIncrementalStep
		toSegment.TolerationThreshold = fromSegment.TolerationThreshold
		toSegment.ActionUpperThreshold = fromSegment.ActionUpperThreshold
	}
}

// resetBandwidthSharingFlowMaxPlannedThroughput -
func (this *DefaultBwSharingAlgorithm) resetBandwidthSharingFlowMaxPlannedThroughput(flow *BandwidthSharingFlow) {

	flow.MaxPlannedThroughput = MAX_THROUGHPUT
	flow.MaxPlannedLowerBound = MAX_THROUGHPUT
	flow.MaxPlannedUpperBound = MAX_THROUGHPUT
}

// updateDefaultConfigAttributes
func (this *DefaultBwSharingAlgorithm) updateDefaultConfigAttributes(fieldName string, fieldValue string) {
	switch fieldName {
	case "maxBwPerInactiveFlow":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			this.DefaultConfigs.MaxBwPerInactiveFlow = value
		}
	case "minActivityThreshold":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			this.DefaultConfigs.MinActivityThreshold = value
		}
	case "incrementalStep":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			this.DefaultConfigs.IncrementalStep = value
		}
	case "inactivityIncrementalStep":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			this.DefaultConfigs.InactivityIncrementalStep = value
		}
	case "tolerationThreshold":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			this.DefaultConfigs.TolerationThreshold = value
		}
	case "actionUpperThreshold":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			this.DefaultConfigs.ActionUpperThreshold = value
		}
	case "isPercentage":
		if "yes" == fieldValue {
			this.DefaultDebugConfigs.IsPercentage = true
		} else {
			this.DefaultDebugConfigs.IsPercentage = false
		}
	case "logVerbose":
		if "yes" == fieldValue {
			this.DefaultDebugConfigs.LogVerbose = true
		}
	default:
	}
}

// getMetricsThroughputEntryHandler -
func (this *DefaultBwSharingAlgorithm) getMetricsThroughputEntryHandler(key string, fields map[string]string, userData interface{}) error {

	subKey := strings.Split(key, ":")

	for trafficFrom, throughput := range fields {

		bwInfo := this.getBandwidthSharingFlow(trafficFrom + ":" + subKey[1])
		if bwInfo != nil {
			value, _ := strconv.ParseFloat(throughput, 64)
			bwInfo.CurrentThroughput = value
		}
	}
	return nil
}

// createMetricsThroughputEntries -
func (this *DefaultBwSharingAlgorithm) createMetricsThroughputEntries(srcElem string, dstElem string) {

	var creationTime = make(map[string]interface{})
	creationTime["creationTime"] = time.Now()

	//entries are created with no values, sidecar will only fill them, otherwise, won't be cleared
	_ = this.ParentBwSharing.rcCtrlEng.SetEntry(moduleMetrics+":"+dstElem+":"+srcElem, creationTime)
	_ = this.ParentBwSharing.rcCtrlEng.SetEntry(moduleMetrics+":"+dstElem+":throughput", creationTime)
}

// deleteAllMetricsThroughputEntries -
func (this *DefaultBwSharingAlgorithm) deleteAllMetricsThroughputEntries() {

	for _, flow := range this.BandwidthSharingFlowMap {
		//entries are created with no values, sidecar will only fill them, otherwise, won't be cleared
		_ = this.ParentBwSharing.rcCtrlEng.DelEntry(moduleMetrics + ":" + flow.DstNetworkElement + ":" + flow.SrcNetworkElement)
		_ = this.ParentBwSharing.rcCtrlEng.DelEntry(moduleMetrics + ":" + flow.DstNetworkElement + ":throughput")
	}
}

// populateBandwidthSharingFlow  - creation of a flow
func (this *DefaultBwSharingAlgorithm) populateBandwidthSharingFlow(flowName string, srcElement *NetElem, destElement *NetElem, maxBw float64) {
	bwSharingInfo := this.BandwidthSharingFlowMap[flowName]
	//maxBw is the min of the 2 ends if a max is not forced
	if maxBw == 0 {
		if srcElement.MaxThroughput < destElement.MaxThroughput {
			maxBw = srcElement.MaxThroughput
		} else {
			maxBw = destElement.MaxThroughput
		}
	}

	if bwSharingInfo == nil {
		bwSharingInfo = new(BandwidthSharingFlow)
		bwSharingInfo.Name = flowName
		bwSharingInfo.SrcNetworkElement = srcElement.Name
		bwSharingInfo.DstNetworkElement = destElement.Name
		bwSharingInfo.MaximumThroughput = maxBw
		this.BandwidthSharingFlowMap[flowName] = bwSharingInfo
	} else {
		if bwSharingInfo.Name == flowName && bwSharingInfo.SrcNetworkElement == srcElement.Name && bwSharingInfo.DstNetworkElement == destElement.Name {
			bwSharingInfo.MaximumThroughput = maxBw
		} else {
			log.Error("bwSharingElement already exists but not the same info, something is wrong!")
		}
	}
	//recreate the path it is part of (the segments might differ from previous ones)
	bwSharingInfo.Path = this.createPath(flowName, srcElement, destElement)
}

// createPath -
func (this *DefaultBwSharingAlgorithm) createPath(flowName string, srcElement *NetElem, destElement *NetElem) *Path {

	//Tier 1 -- check if they are in the same poa
	//Tier 2 -- check if they are in the same zone, but different poa
	//Tier 3 -- check if they are in the same domain, but different zone
	//Tier 4 -- check if they are in different domains

	direction := ""
	segmentName := ""
	var segment *BandwidthSharingSegment

	path := new(Path)
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
		segment = this.createSegment(segmentName, flowName, srcElement.PhyLocMaxThroughput)
		path.Segments = append(path.Segments, segment)

		if srcElement.Type == "UE" {
			//segment 2A
			segmentName = srcElement.PoaName + "-" + direction
			segment = this.createSegment(segmentName, flowName, srcElement.PoaMaxThroughput)
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
			segment = this.createSegment(segmentName, flowName, destElement.PoaMaxThroughput)
			path.Segments = append(path.Segments, segment)
		}

		//segment 3A or 3B
		segmentName = destElement.PhyLocName + "-" + direction
		segment = this.createSegment(segmentName, flowName, destElement.PhyLocMaxThroughput)
		path.Segments = append(path.Segments, segment)
	}
	//	}

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
			segment = this.createSegment(segmentName, flowName, srcElement.PhyLocMaxThroughput)
			path.Segments = append(path.Segments, segment)

			//segment 2A
			segmentName = srcElement.ZoneName + "-" + srcElement.PhyLocName + "-" + direction
		} else {
			//segment 2B
			segmentName = srcElement.ZoneName + "-" + srcElement.PoaName + "-" + direction
		}
		if srcElement.ZoneName != "" {
			segment = this.createSegment(segmentName, flowName, srcElement.IntraZoneMaxThroughput)
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
			segment = this.createSegment(segmentName, flowName, destElement.PhyLocMaxThroughput)
			path.Segments = append(path.Segments, segment)

			//segment 3A
			segmentName = destElement.ZoneName + "-" + destElement.PhyLocName + "-" + direction
		} else {
			//segment 3B
			segmentName = destElement.ZoneName + "-" + destElement.PoaName + "-" + direction
		}
		if destElement.ZoneName != "" {
			segment = this.createSegment(segmentName, flowName, destElement.IntraZoneMaxThroughput)
			path.Segments = append(path.Segments, segment)
		}
		//	}
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
			segment = this.createSegment(segmentName, flowName, srcElement.InterZoneMaxThroughput)
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
			segment = this.createSegment(segmentName, flowName, destElement.InterZoneMaxThroughput)
			path.Segments = append(path.Segments, segment)
		}
	}
	//        }

	//Tier 4
	if srcElement.DomainName != destElement.DomainName {
		//segments to interDomain backbone
		//1 possibility
		//InterZoneBackbone->InterDomainBackbone
		//1A. InterZone backbone -> Domain backbone
		direction = "uplink"
		//segment 1A
		segmentName = srcElement.DomainName + "-" + direction
		segment = this.createSegment(segmentName, flowName, srcElement.InterDomainMaxThroughput)
		path.Segments = append(path.Segments, segment)

		//segments from interDomain backbone
		//1 possibility
		//InterDomainBackbone->InterZoneBackbone
		//2A. Domain backbone -> InterZone backbone
		direction = "downlink"
		//segment 2A
		segmentName = destElement.DomainName + "-" + direction
		segment = this.createSegment(segmentName, flowName, destElement.InterDomainMaxThroughput)
		path.Segments = append(path.Segments, segment)

		//when going through interdomain, only to/from a cloud app is supported
		if destElement.Type == "CLOUD" {
			segmentName = destElement.PhyLocName + "-" + direction
			segment = this.createSegment(segmentName, flowName, destElement.PhyLocMaxThroughput)
			path.Segments = append(path.Segments, segment)
		} else {
			if srcElement.Type == "CLOUD" {
				direction = "uplink"
				segmentName = srcElement.PhyLocName + "-" + direction
				segment = this.createSegment(segmentName, flowName, srcElement.PhyLocMaxThroughput)
				path.Segments = append(path.Segments, segment)
			} else {
				return nil
			}
		}
	}

	return path
}

// createSegment -
func (this *DefaultBwSharingAlgorithm) createSegment(segmentName string, flowName string, maxThroughput float64) *BandwidthSharingSegment {
	segment := this.getBandwidthSharingSegment(segmentName)
	if segment == nil {
		segment = new(BandwidthSharingSegment)
		segment.Name = segmentName
		this.SegmentsMap[segmentName] = segment
		this.setSegmentConfiguration(&segment.Config, &this.DefaultConfigs, maxThroughput)

	}
	segment.MaxThroughput = maxThroughput
	flow := this.getBandwidthSharingFlow(flowName)
	if flow != nil {
		segment.Flows = append(segment.Flows, flow)
	}
	return segment
}

// tickerFunction - function called periodically to get metrics, calculate bandwidth values per flow per segment
//func (this DefaultBwSharingAlgorithm) tickerFunction(bw *BwSharing, algo *DefaultBwSharingAlgorithm) {
func (this *DefaultBwSharingAlgorithm) tickerFunction() {
	var start time.Time
	var elapsed time.Duration
	if this.DefaultDebugConfigs.LogVerbose {
		start = time.Now()
		log.Info("******************************************************************************************")
		elapsed = time.Since(start)
		log.WithFields(log.Fields{
			"meep.log.component": "meep-bw-sharing",
			"meep.time.location": "time to print",
			"meep.time.exec":     elapsed,
		}).Info("Measurements log")
	}

	keyName := moduleMetrics + ":*:throughput"
	err := this.ParentBwSharing.rcCtrlEng.ForEachEntry(keyName, this.getMetricsThroughputEntryHandler, nil)
	if err != nil {
		log.Error("Failed to get entries: ", err)
		return
	}

	if this.DefaultDebugConfigs.LogVerbose {
		elapsed = time.Since(start) - elapsed
		log.WithFields(log.Fields{
			"meep.log.component": "meep-bw-sharing",
			"meep.time.location": "time to update metrics",
			"meep.time.exec":     elapsed,
		}).Info("Measurements log")
	}

	this.reCalculateThroughputs()

	if this.DefaultDebugConfigs.LogVerbose {
		elapsed = time.Since(start) - elapsed
		log.WithFields(log.Fields{
			"meep.log.component": "meep-bw-sharing",
			"meep.time.location": "time to recalculate",
			"meep.time.exec":     elapsed,
		}).Info("Measurements log")
	}
}

// reCalculateThroughputs -
func (this *DefaultBwSharingAlgorithm) reCalculateThroughputs() {

	//reset every planned throughput values for every flow since they will start to populate those
	this.resetAllBandwidthSharingFlowMaxPlannedThroughput()

	//all segments determined by the scenario
	for _, segment := range this.SegmentsMap {
		this.updateMaxFairShareBwPerFlow(segment)

		unusedBw, list := this.needToReevaluate(segment)

		if list != nil {
			if this.DefaultDebugConfigs.LogVerbose {
				log.Info("Segment ", segment.Name, " reevaluation result - BW unused: ", unusedBw, "***flows to evaluate***: ", printBwFlowsNameFromList(list))
			}

			this.recalculateSegment(segment, list, unusedBw)

			if this.DefaultDebugConfigs.LogVerbose {
				printBwFlows(segment)
			}
		}
	}

	//apply to all flows
	this.updateAllBandwidthSharingFlow()
}

// resetAllBandwidthSharingFlowMaxPlannedThroughput -
func (this *DefaultBwSharingAlgorithm) resetAllBandwidthSharingFlowMaxPlannedThroughput() {
	for _, flow := range this.BandwidthSharingFlowMap {
		this.resetBandwidthSharingFlowMaxPlannedThroughput(flow)
	}
}

// updateAllBandwidthSharingFlow -
func (this *DefaultBwSharingAlgorithm) updateAllBandwidthSharingFlow() {

	for _, flow := range this.BandwidthSharingFlowMap {

		if flow.MaxPlannedThroughput != flow.AllocatedThroughput && flow.MaxPlannedThroughput != MAX_THROUGHPUT {
			log.Info("Update allocated bandwidth for ", flow.Name, " to ", flow.MaxPlannedThroughput)
			flow.AllocatedThroughput = flow.MaxPlannedThroughput
			flow.AllocatedThroughputLowerBound = flow.MaxPlannedLowerBound
			flow.AllocatedThroughputUpperBound = flow.MaxPlannedUpperBound
			this.ParentBwSharing.updateFilter(flow.DstNetworkElement, flow.SrcNetworkElement, flow.AllocatedThroughput)
			this.ParentBwSharing.applyFilter()
		}
	}
}

// recalculateSegment -
func (this *DefaultBwSharingAlgorithm) recalculateSegment(segment *BandwidthSharingSegment, flowsToEvaluate []*BandwidthSharingFlow, unusedBw float64) {

	nbEvaluatedFlowsLeft := len(flowsToEvaluate)

	if segment.CurrentThroughput > segment.MaxThroughput || nbEvaluatedFlowsLeft >= 1 {

		//category 1 flows
		for _, flow := range flowsToEvaluate {
			if flow.CurrentThroughput+segment.Config.IncrementalStep > segment.MaxFairShareBwPerFlow {
				flow.PlannedThroughput = segment.MaxFairShareBwPerFlow //category 2 or 3
			} else {
				if flow.CurrentThroughput <= segment.Config.MinActivityThreshold {

					flow.PlannedThroughput = segment.Config.MaxBwPerInactiveFlow
					flow.PlannedUpperBound = segment.Config.InactivityIncrementalStep
					flow.PlannedLowerBound = 0
				} else {
					flow.PlannedThroughput = flow.CurrentThroughput + segment.Config.IncrementalStep
					if flow.PlannedThroughput > flow.MaximumThroughput {
						flow.PlannedThroughput = flow.MaximumThroughput
					}
					flow.PlannedUpperBound = flow.PlannedThroughput - segment.Config.ActionUpperThreshold
					flow.PlannedLowerBound = flow.PlannedUpperBound - segment.Config.TolerationThreshold
					//lower bound cannot be less than min threshold
					if flow.PlannedLowerBound < segment.Config.MinActivityThreshold {
						flow.PlannedLowerBound = segment.Config.MinActivityThreshold
					}
				}
				nbEvaluatedFlowsLeft--
				if flow.PlannedThroughput != segment.Config.MaxBwPerInactiveFlow {
					unusedBw -= flow.PlannedThroughput
				}
			}
		}

		var extra float64

		if nbEvaluatedFlowsLeft > 0 {

			//category 2 flows
			for _, flow := range flowsToEvaluate {

				if flow.PlannedThroughput == segment.MaxFairShareBwPerFlow {
					if flow.CurrentThroughput < segment.MaxFairShareBwPerFlow {

						nbEvaluatedFlowsLeft--
						if nbEvaluatedFlowsLeft == 0 { //allocate everything of what is left
							flow.PlannedThroughput = unusedBw
							if flow.PlannedThroughput > flow.MaximumThroughput {
								flow.PlannedThroughput = flow.MaximumThroughput
							}
							flow.PlannedUpperBound = flow.PlannedThroughput
							flow.PlannedLowerBound = flow.PlannedThroughput - segment.Config.TolerationThreshold
							//lower bound cannot be less than min threshold
							if flow.PlannedLowerBound < segment.Config.MinActivityThreshold {
								flow.PlannedLowerBound = segment.Config.MinActivityThreshold
							}
						} else {
							flow.PlannedThroughput = flow.CurrentThroughput + segment.Config.IncrementalStep
							if flow.PlannedThroughput > flow.MaximumThroughput {
								flow.PlannedThroughput = flow.MaximumThroughput
							}

							flow.PlannedUpperBound = flow.PlannedThroughput - segment.Config.ActionUpperThreshold
							flow.PlannedLowerBound = flow.PlannedUpperBound - segment.Config.TolerationThreshold
							//lower bound cannot be less than min threshold
							if flow.PlannedLowerBound < segment.Config.MinActivityThreshold {
								flow.PlannedLowerBound = segment.Config.MinActivityThreshold
							}

						}
						unusedBw -= flow.PlannedThroughput
					}
				}
			}

			if nbEvaluatedFlowsLeft > 0 {
				if nbEvaluatedFlowsLeft >= 1 {
					extra = (unusedBw - float64(nbEvaluatedFlowsLeft)*segment.MaxFairShareBwPerFlow) / float64(nbEvaluatedFlowsLeft)
				} else {
					extra = 0
				}

				//category 3
				for _, flow := range flowsToEvaluate {
					if flow.PlannedThroughput == segment.MaxFairShareBwPerFlow && flow.CurrentThroughput >= segment.MaxFairShareBwPerFlow {
						flow.PlannedThroughput = segment.MaxFairShareBwPerFlow + extra
						if flow.PlannedThroughput > flow.MaximumThroughput {
							flow.PlannedThroughput = flow.MaximumThroughput
						}

						flow.PlannedUpperBound = flow.PlannedThroughput - segment.Config.ActionUpperThreshold
						flow.PlannedLowerBound = flow.PlannedUpperBound - segment.Config.TolerationThreshold
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
			if flow.CurrentThroughput > segment.Config.MinActivityThreshold {
				flow.PlannedThroughput = segment.MaxFairShareBwPerFlow
				if flow.PlannedThroughput > flow.MaximumThroughput {
					flow.PlannedThroughput = flow.MaximumThroughput
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

// deallocateBandwidthSharing -
func (this *DefaultBwSharingAlgorithm) deallocateBandwidthSharing() {
	this.BandwidthSharingFlowMap = nil
	this.SegmentsMap = nil
}

// initDefaultConfigAttributes -
func (this *DefaultBwSharingAlgorithm) initDefaultConfigAttributes() {
	this.DefaultConfigs.MaxBwPerInactiveFlow = 2.0
	this.DefaultConfigs.MinActivityThreshold = 0.3
	this.DefaultConfigs.IncrementalStep = 3.0
	this.DefaultConfigs.InactivityIncrementalStep = 1.0
	this.DefaultConfigs.ActionUpperThreshold = 1.0
	this.DefaultConfigs.TolerationThreshold = 4.0
}

// updateMaxFairShareBwPerFlow -
func (this *DefaultBwSharingAlgorithm) updateMaxFairShareBwPerFlow(segment *BandwidthSharingSegment) {

	nbActiveConnections := 0
	for _, flow := range segment.Flows {
		if flow.CurrentThroughput >= segment.Config.MinActivityThreshold {
			nbActiveConnections++
		}
	}
	if nbActiveConnections >= 1 {
		segment.MaxFairShareBwPerFlow = segment.MaxThroughput / float64(nbActiveConnections)
	} else {
		segment.MaxFairShareBwPerFlow = MAX_THROUGHPUT
	}
}

// needToReevaluate - determines which flows must be recalculated for bandwidth sharing within the segment
func (this *DefaultBwSharingAlgorithm) needToReevaluate(segment *BandwidthSharingSegment) (unusedBw float64, list []*BandwidthSharingFlow) {

	unusedBw = segment.MaxThroughput

	//how many active connections that needs to be taken into account
	for _, flow := range segment.Flows {

		if flow.CurrentThroughput < flow.AllocatedThroughputLowerBound || flow.CurrentThroughput > flow.AllocatedThroughputUpperBound || flow.CurrentThroughput >= segment.MaxFairShareBwPerFlow {
			//			resetBandwidthSharingFlowMaxPlannedThroughput(flow)
			list = append(list, flow)

		} else {
			//no need to reevalute this one, so removing its allocated bw from the available one
			unusedBw -= flow.AllocatedThroughput
		}
		if flow.CurrentThroughput < segment.Config.MinActivityThreshold {
			//we just re-add the bw for inactive connections
			unusedBw += flow.AllocatedThroughput
		}
	}
	return unusedBw, list
}

// printBwFlowsNameFromList -
func printBwFlowsNameFromList(list []*BandwidthSharingFlow) string {

	str := ""
	for _, flow := range list {
		str += flow.Name + "."
	}
	return str
}

// printBwFlows -
func printBwFlows(segment *BandwidthSharingSegment) {

	log.Info("Flows on segment ", segment.Name)

	for _, flow := range segment.Flows {
		log.Info(printBwFlow(flow))
	}
}

// printBwFlow -
func printBwFlow(flow *BandwidthSharingFlow) string {

	s0 := fmt.Sprintf("%x", &flow)
	s1 := flow.Name + "(" + s0 + ")"
	s2 := fmt.Sprintf("%f", flow.MaximumThroughput)
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

	str := s1 + ": " + "Current: " + s6 + " - Max: " + s2 + " - Allocated: " + s3a + "[" + s4a + "-" + s5a + "]" + " - MaxPlanned: " + s3m + "[" + s4m + "-" + s5m + "]" + " - Planned: " + s3p + "[" + s4p + "-" + s5p + "] "

	str += printPath(flow.Path)

	return str
}

// printPath -
func printPath(path *Path) string {
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

// parseScenario -
func (this *DefaultBwSharingAlgorithm) parseScenario(scenario ceModel.Scenario) {
	var netElemList []NetElem

	if scenario.Name != "" {
		//reinitialise structures
		this.allocateBandwidthSharing()

		// Parse Domains
		for _, domain := range scenario.Deployment.Domains {

			// Parse Zones
			for _, zone := range domain.Zones {

				// Parse Network Locations
				for _, nl := range zone.NetworkLocations {

					// Parse Physical locations
					for _, pl := range nl.PhysicalLocations {

						// Parse Processes
						for _, proc := range pl.Processes {

							element := new(NetElem)
							element.Name = proc.Name
							element.Type = pl.Type_
							// Update element information based on current location characteristics
							element.PhyLocName = pl.Name
							element.DomainName = domain.Name
							element.ZoneName = zone.Name
							if pl.Type_ == "UE" || pl.Type_ == "FOG" {
								element.PoaName = nl.Name
								element.PoaMaxThroughput = float64(nl.TerminalLinkThroughput)
							}
							element.PhyLocMaxThroughput = float64(pl.LinkThroughput)
							element.MaxThroughput = float64(proc.AppThroughput)
							element.IntraZoneMaxThroughput = float64(zone.EdgeFogThroughput)
							element.InterZoneMaxThroughput = float64(domain.InterZoneThroughput)
							element.InterDomainMaxThroughput = float64(scenario.Deployment.InterDomainThroughput)

							//to support scenarios without this info (compatibility with old scenarios)
							if element.MaxThroughput == 0 {
								element.MaxThroughput = DEFAULT_THROUGHPUT_LINK
							}
							if element.PhyLocMaxThroughput == 0 {
								element.PhyLocMaxThroughput = DEFAULT_THROUGHPUT_LINK
							}

							netElemList = append(netElemList, *element)
						}
					}
				}
			}
		}

		//bandwidth sharing creation elements section
		//all the dummy processes meant for calculation can be ignored, others must have bandwidth flows created between each of them
		for _, elemSrc := range netElemList {
			for _, elemDest := range netElemList {
				if elemSrc.Name != elemDest.Name {
					this.populateBandwidthSharingFlow(elemSrc.Name+":"+elemDest.Name, &elemSrc, &elemDest, 0)
					//create entries in DB that will be populated by the sidecar
					this.createMetricsThroughputEntries(elemSrc.Name, elemDest.Name)

				}
			}
		}

		if this.DefaultDebugConfigs.LogVerbose {
			log.Info("Segments map: ", this.SegmentsMap)
			log.Info("Flows map: ", this.BandwidthSharingFlowMap)
		}
	} else {
		//metrics created while parsing a scenario needs to be cleared when parsing nil scenario (stop scenario)
		this.deleteAllMetricsThroughputEntries()
	}

}
