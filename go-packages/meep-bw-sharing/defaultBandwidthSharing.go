/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
package bws

import (
	"fmt"
	"strings"
	"strconv"
	"time"

	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"

)

// DebugConfiguration -
type DebugConfiguration struct {
	EnableTier1  bool
	EnableTier2  bool
	EnableTier3  bool
        EnableTier4  bool
	IsPercentage bool

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
	Name                          string
	Segments                      []*BandwidthSharingSegment
}

// NetElem -
type NetElem struct {
        Name                string
	Type                string
	PhyLocName          string
        PoaName             string
        ZoneName            string
	DomainName          string
	MaxThroughput       float64
	PhyLocMaxThroughput float64
	PoaMaxThroughput    float64
	ZoneMaxThroughput   float64
	DomainMaxThroughput float64
}

// flows and segments mappings
var bandwidthSharingFlowMap map[string]*BandwidthSharingFlow
var segmentsMap map[string]*BandwidthSharingSegment

var defaultConfigs SegmentConfiguration
var defaultDebugConfigs DebugConfiguration 
var logVerbose bool

// allocateBandwidthSharing - allocated structures
func allocateBandwidthSharing() {
	bandwidthSharingFlowMap = make(map[string]*BandwidthSharingFlow)
	segmentsMap = make(map[string]*BandwidthSharingSegment)
}

// getBandwidthSharingFlow -
func getBandwidthSharingFlow(key string) *BandwidthSharingFlow {

	return bandwidthSharingFlowMap[key]
}

// getBandwidthSharingSegment -
func getBandwidthSharingSegment(key string) *BandwidthSharingSegment {

        return segmentsMap[key]
}

// setSegmentConfiguration - set configuration attributes on a per segment basis (absolute or relative %)
func setSegmentConfiguration(toSegment *SegmentConfiguration, fromSegment *SegmentConfiguration, isPercentage bool, maxThroughput float64) {

	if isPercentage {
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
func resetBandwidthSharingFlowMaxPlannedThroughput(flow *BandwidthSharingFlow) {

					
	flow.MaxPlannedThroughput = MAX_THROUGHPUT
	flow.MaxPlannedLowerBound = MAX_THROUGHPUT
	flow.MaxPlannedUpperBound = MAX_THROUGHPUT
}

// updateDefaultConfigAttributes
func updateDefaultConfigAttributes(fieldName string, fieldValue string) {
        switch(fieldName) {
        case "maxBwPerInactiveFlow":
		value, err := strconv.ParseFloat(fieldValue, 64)
		if err == nil {
			defaultConfigs.MaxBwPerInactiveFlow = value
		}
	case "minActivityThreshold":
                value, err := strconv.ParseFloat(fieldValue, 64)
                if err == nil {
                        defaultConfigs.MinActivityThreshold = value
                }
	case "incrementalStep":
                value, err := strconv.ParseFloat(fieldValue, 64)
                if err == nil {
                        defaultConfigs.IncrementalStep = value
                }
	case "inactivityIncrementalStep":
                value, err := strconv.ParseFloat(fieldValue, 64)
                if err == nil {
                        defaultConfigs.InactivityIncrementalStep = value
                }
	case "tolerationThreshold":
                value, err := strconv.ParseFloat(fieldValue, 64)
                if err == nil {
                        defaultConfigs.TolerationThreshold = value
                }
	case "actionUpperThreshold":
                value, err := strconv.ParseFloat(fieldValue, 64)
                if err == nil {
                        defaultConfigs.ActionUpperThreshold = value
                }
	case "isPercentage":
                if "yes" == fieldValue {
                        defaultDebugConfigs.IsPercentage = true
                } else {
                        defaultDebugConfigs.IsPercentage = false
                }
        case "enableTier1":
                if "yes" == fieldValue {
                        defaultDebugConfigs.EnableTier1 = true
                } else {
			defaultDebugConfigs.EnableTier1 = false
		}
        case "enableTier2":
                if "yes" == fieldValue {
                        defaultDebugConfigs.EnableTier2 = true
                } else {
                        defaultDebugConfigs.EnableTier2 = false
                }
        case "enableTier3":
                if "yes" == fieldValue {
                        defaultDebugConfigs.EnableTier3 = true
                } else {
                        defaultDebugConfigs.EnableTier3 = false
                }
	default:
	}
}

// getMetricsThroughputEntryHandler -
func getMetricsThroughputEntryHandler(key string, fields map[string]string, userData interface{}) error {


        subKey :=  strings.Split(key, ":")

        for trafficFrom, throughput := range fields {

                bwInfo := getBandwidthSharingFlow(trafficFrom + ":" + subKey[1])
		if bwInfo != nil {
                	value, _ := strconv.ParseFloat(throughput, 64)
	                bwInfo.CurrentThroughput = value
		}
        }
        return nil
}

// populateBandwidthSharingFlow  - creation of a flow
func populateBandwidthSharingFlow(flowName string, srcElement *NetElem, destElement *NetElem, maxBw float64) {
	bwSharingInfo := bandwidthSharingFlowMap[flowName]
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
                bandwidthSharingFlowMap[flowName] = bwSharingInfo
        } else {
		if bwSharingInfo.Name == flowName && bwSharingInfo.SrcNetworkElement == srcElement.Name && bwSharingInfo.DstNetworkElement == destElement.Name {
			bwSharingInfo.MaximumThroughput = maxBw
		} else {
			log.Error("bwSharingElement already exists but not the same info, something is wrong!")
		}
	}
	//recreate the path it is part of (the segments might differ from previous ones)
	bwSharingInfo.Path = createPath(flowName, srcElement, destElement)
}

// createPath -
func createPath(flowName string, srcElement *NetElem, destElement *NetElem) *Path {

	//Tier 1 -- check if they are in the same poa
	//Tier 2 -- check if they are in the same zone
	//Tier 3 -- check if they are in the same domain
	//Tier 4 -- check if they are in different domains

	direction := ""
	if srcElement.Type == "UE" && destElement.Type != "UE" {
		direction = "uplink"
	} else {
		if srcElement.Type != "UE" && destElement.Type == "UE" {
       			direction = "downlink"
        	}
	}

        //communication between UE-apps and EDGE-apps only
	if direction == "" {
		return nil
	}

	path := new(Path) 
	path.Name = flowName
	if srcElement.DomainName == destElement.DomainName {
		if srcElement.ZoneName == destElement.ZoneName {
			if srcElement.PoaName == destElement.PoaName && srcElement.PoaName != "" {
				//Tier 1
				if defaultDebugConfigs.EnableTier1 {
					//2 possibilities
					//segments for srcElement(app) -> 1A. UE-Node uplink-> 2A. POA-TermLink uplink -> 3A. Fog-Node uplink -> dstElement(app)
					//segments for srcElement(app) -> 1B. FogNode downlink -> 2B. POA-TermLink downlink -> 3B. UE-Node downlink -> dstElement(app)

					//segment 1A or 1B
				        segmentName := srcElement.PhyLocName + "-" + direction
					segment := createSegment(segmentName, flowName, srcElement.PhyLocMaxThroughput, &defaultConfigs, defaultDebugConfigs.IsPercentage)
					path.Segments = append(path.Segments, segment)

					//segment 2A or 2B
					segmentName = srcElement.PoaName + "-" + direction
/*					if srcElement.Type == "UE" {
						segmentName = srcElement.PoaName + "-ingress"
					} else {
						segmentName = srcElement.PoaName + "-egress"
					}
*/
                                        segment = createSegment(segmentName, flowName, srcElement.PoaMaxThroughput, &defaultConfigs, defaultDebugConfigs.IsPercentage)
                                        path.Segments = append(path.Segments, segment)

					//segment 3A or 3B
                                        segmentName = destElement.PhyLocName + "-" + direction
                                        segment = createSegment(segmentName, flowName, destElement.PhyLocMaxThroughput, &defaultConfigs, defaultDebugConfigs.IsPercentage)
                                        path.Segments = append(path.Segments, segment)
				} 
			} else {
				//Tier 2
				if defaultDebugConfigs.EnableTier2 {

                                }
			}
		} else {
			//Tier 3
			if defaultDebugConfigs.EnableTier3 {

                        }
		}
	} else {
		//Tier 4
		if defaultDebugConfigs.EnableTier4 {

                }
	}
	return path
}

// createSegment -
func createSegment(segmentName string, flowName string, maxThroughput float64, config *SegmentConfiguration, isPercentage bool) (*BandwidthSharingSegment) {
        segment := getBandwidthSharingSegment(segmentName)
        if segment == nil {
 	        segment = new(BandwidthSharingSegment)
                segment.Name = segmentName
                segmentsMap[segmentName] = segment
		setSegmentConfiguration(&segment.Config, config, isPercentage, maxThroughput)

        }
        segment.MaxThroughput = maxThroughput
        flow := getBandwidthSharingFlow(flowName)
        if flow != nil {
                segment.Flows = append(segment.Flows, flow)
        }
	return segment
}

// tickerFunction - function called periodically to get metrics, calculate bandwidth values per flow per segment  
func tickerFunction(rcCtrlEng *redis.Connector, globalLogVerbose bool, updateFilter func(string, string, float64), applyFilter func()) {

	var start time.Time
	var elapsed time.Duration
	logVerbose = globalLogVerbose

	if logVerbose {
		start = time.Now()
       		log.Info("******************************************************************************************")
	        elapsed = time.Since(start)
	        log.WithFields(log.Fields{
	        	"meep.log.component": "meep-bw-sharing",
		        "meep.time.location": "time to print",
		        "meep.time.exec":     elapsed,
	        	}).Info("Measurements log")
	}

        keyName := moduleMetrics+":*:throughput"
        err := rcCtrlEng.ForEachEntry(keyName, getMetricsThroughputEntryHandler, nil)
        if err != nil {
                log.Error("Failed to get entries: ", err)
                return
        }

	if logVerbose {
	        elapsed = time.Since(start) - elapsed
	        log.WithFields(log.Fields{
	        	"meep.log.component": "meep-bw-sharing",
		        "meep.time.location": "time to update metrics",
		        "meep.time.exec":     elapsed,
		        }).Info("Measurements log")
	}

        reCalculateThroughputs(updateFilter, applyFilter)

	if logVerbose {
	        elapsed = time.Since(start) - elapsed
	        log.WithFields(log.Fields{
	        	"meep.log.component": "meep-bw-sharing",
		        "meep.time.location": "time to recalculate",
		        "meep.time.exec":     elapsed,
		        }).Info("Measurements log")
	}
}

// reCalculateThroughputs -
func reCalculateThroughputs(updateFilter func(string, string, float64), applyFilter func()) {

	//reset every planned throughput values for every flow since they will start to populate those
	resetAllBandwidthSharingFlowMaxPlannedThroughput()

	//all segments determined by the scenario
	for _, segment := range segmentsMap {
		updateMaxFairShareBwPerFlow(segment)

		unusedBw, list := needToReevaluate(segment)
	
		if list != nil {
	                if logVerbose {
       	                	log.Info("Segment ", segment.Name, " reevaluation result - BW unused: ", unusedBw, "***flows to evaluate***: ", printBwFlowsNameFromList(list))
       		        }

			recalculateSegment(segment, list, unusedBw)

			if logVerbose {
				printBwFlows(segment)
			}
		}
	}

	//apply to all flows
	updateAllBandwidthSharingFlow(updateFilter, applyFilter)
}

// resetAllBandwidthSharingFlowMaxPlannedThroughput -
func resetAllBandwidthSharingFlowMaxPlannedThroughput() {
        for _, flow := range bandwidthSharingFlowMap {
                resetBandwidthSharingFlowMaxPlannedThroughput(flow)
        }
}

// updateAllBandwidthSharingFlow -
func updateAllBandwidthSharingFlow(updateFilter func(string, string, float64), applyFilter func()) {

	for _, flow := range bandwidthSharingFlowMap {

		//flow := bandwidthSharingFlowMap[key]
		if flow.MaxPlannedThroughput != flow.AllocatedThroughput && flow.MaxPlannedThroughput != MAX_THROUGHPUT {
			log.Info("Update allocated bandwidth for ", flow.Name, " to ", flow.MaxPlannedThroughput)
			flow.AllocatedThroughput = flow.MaxPlannedThroughput
			flow.AllocatedThroughputLowerBound = flow.MaxPlannedLowerBound
			flow.AllocatedThroughputUpperBound = flow.MaxPlannedUpperBound
			updateFilter(flow.DstNetworkElement, flow.SrcNetworkElement, flow.AllocatedThroughput)
			applyFilter()
		}
	}
}

// recalculateSegment -
func recalculateSegment(segment *BandwidthSharingSegment, flowsToEvaluate []*BandwidthSharingFlow, unusedBw float64) {

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

// cleanUp -
func cleanUp() {
	//nothing allocated that should be cleared
}

// initDefaultConfigAttributes -
func initDefaultConfigAttributes() {
        defaultConfigs.MaxBwPerInactiveFlow = 2.0
        defaultConfigs.MinActivityThreshold = 0.3
        defaultConfigs.IncrementalStep = 3.0
        defaultConfigs.InactivityIncrementalStep = 1.0
        defaultConfigs.ActionUpperThreshold = 1.0
        defaultConfigs.TolerationThreshold = 4.0
}

// updateMaxFairShareBwPerFlow -
func updateMaxFairShareBwPerFlow(segment *BandwidthSharingSegment) {

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
func needToReevaluate(segment *BandwidthSharingSegment) (unusedBw float64, list []*BandwidthSharingFlow) {

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

	str += printPath(flow)

	return str
}

// printPath -
func printPath(flow *BandwidthSharingFlow) string {
	str := ""
	first := true
	if flow.Path != nil {
		str = "Path: "
		for _, segment := range flow.Path.Segments {
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
func parseScenario(scenario ceModel.Scenario) {

	var netElemList []NetElem
	//reinitialise structures
	allocateBandwidthSharing()
	
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
                                populateBandwidthSharingFlow(elemSrc.Name+":"+elemDest.Name, &elemSrc, &elemDest, 0)
                        }
		}
        }

	if logVerbose {
		log.Info("Segments map: ", segmentsMap)
	        log.Info("Flows map: ", bandwidthSharingFlowMap)
	}
}

