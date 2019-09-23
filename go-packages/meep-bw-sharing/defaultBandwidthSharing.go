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

type DebugConfiguration struct {
	EnableTier1 bool
	EnableTier2 bool
	EnableTier3 bool
        EnableTier4 bool
}

type SegmentConfiguration struct {
	MaxBwPerInactiveFlow      float64
	MinActivityThreshold      float64
	IncrementalStep           float64
	InactivityIncrementalStep float64
	TolerationThreshold       float64
	ActionUpperThreshold      float64
}

type BandwidthSharingSegment struct {
	Name                  string
	MaxThroughput         float64
	MaxFairShareBwPerFlow float64
	CurrentThroughput     float64
	Config                SegmentConfiguration
	Flows                 []*BandwidthSharingFlow
}

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

type Path struct {
	Name                          string
	Segments                      []*BandwidthSharingSegment
}
	
type NetElem struct {
        Name                string
	Type                string
	PhyLocName          string
        PoaName             string
        ZoneName            string
	DomainName          string
	MaxThroughput       float64
	PoaMaxThroughput    float64
	ZoneMaxThroughput   float64
	DomainMaxThroughput float64
}

// Scenario service mappings
var bandwidthSharingFlowMap map[string]*BandwidthSharingFlow
//var readyToCalculate bool
var segmentsMap map[string]*BandwidthSharingSegment
var defaultConfigs SegmentConfiguration
var defaultDebugConfigs DebugConfiguration 
var testInteger int 
var logVerbose bool

func allocateBandwidthSharing() {
	bandwidthSharingFlowMap = make(map[string]*BandwidthSharingFlow)
	segmentsMap = make(map[string]*BandwidthSharingSegment)

	//for testing/debugging while developping
	//testInteger = 1
	//assignDefaultConfigs()
}


func getBandwidthSharingFlow(key string) *BandwidthSharingFlow {

	return bandwidthSharingFlowMap[key]
}

func getBandwidthSharingSegment(key string) *BandwidthSharingSegment {

        return segmentsMap[key]
}

func setSegmentConfiguration(toSegment *SegmentConfiguration, fromSegment *SegmentConfiguration) {

	toSegment.MaxBwPerInactiveFlow = fromSegment.MaxBwPerInactiveFlow
	toSegment.MinActivityThreshold = fromSegment.MinActivityThreshold 
	toSegment.IncrementalStep = fromSegment.IncrementalStep
	toSegment.InactivityIncrementalStep = fromSegment.InactivityIncrementalStep
	toSegment.TolerationThreshold = fromSegment.TolerationThreshold 
	toSegment.ActionUpperThreshold = fromSegment.ActionUpperThreshold
}

func resetBandwidthSharingFlowMaxPlannedThroughput(flow *BandwidthSharingFlow) {

					
	flow.MaxPlannedThroughput = MAX_THROUGHPUT
	flow.MaxPlannedLowerBound = MAX_THROUGHPUT
	flow.MaxPlannedUpperBound = MAX_THROUGHPUT
}

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

func populateBandwidthSharingFlow(flowName string, srcElement *NetElem, destElement *NetElem, maxBw float64) {
	bwSharingInfo := bandwidthSharingFlowMap[flowName]
	//maxBw is the min of the 2 ends
	if srcElement.MaxThroughput < destElement.MaxThroughput {
		maxBw = srcElement.MaxThroughput
	} else {
		maxBw = destElement.MaxThroughput
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

func createPath(flowName string, srcElement *NetElem, destElement *NetElem) *Path {

	//Tier 1 -- check if they are in the same poa
	//Tier 2 -- check if they are in the same zone
	//Tier 3 -- check if they are in the same domain
	//Tier 4 -- check if they are in different domains

	//communication between UE-apps and EDGE-apps only
	if !((srcElement.Type == "UE" && destElement.Type != "UE") || (srcElement.Type != "UE" && destElement.Type == "UE")) {
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
					//segments for srcElement(app) -> 1A. UE-Node egress-> 2A. POA-TermLink ingress -> 3A. Fog-Node ingress -> dstElement(app)
					//segments for srcElement(app) -> 1B. FogNode egress -> 2B. POA-TermLink egress -> 3B. UE-Node ingress -> dstElement(app)

					//segment 1A or 1B
				        segmentName := srcElement.PhyLocName + "-egress"
					segment := createSegment(segmentName, flowName, srcElement.MaxThroughput, &defaultConfigs)
					path.Segments = append(path.Segments, segment)

					//segment 2A or 2B
					if srcElement.Type == "UE" {
						segmentName = srcElement.PoaName + "-ingress"
					} else {
						segmentName = srcElement.PoaName + "-egress"
					}
                                        segment = createSegment(segmentName, flowName, srcElement.PoaMaxThroughput, &defaultConfigs)
                                        path.Segments = append(path.Segments, segment)

					//segment 3A or 3B
                                        segmentName = destElement.PhyLocName + "-ingress"
                                        segment = createSegment(segmentName, flowName, destElement.MaxThroughput, &defaultConfigs)
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

func createSegment(segmentName string, flowName string, maxThroughput float64, config *SegmentConfiguration) (*BandwidthSharingSegment) {
        segment := getBandwidthSharingSegment(segmentName)
        if segment == nil {
 	        segment = new(BandwidthSharingSegment)
                segment.Name = segmentName
                segmentsMap[segmentName] = segment
		setSegmentConfiguration(&segment.Config, config)

        }
        segment.MaxThroughput = maxThroughput
        flow := getBandwidthSharingFlow(flowName)
        if flow != nil {
                segment.Flows = append(segment.Flows, flow)
        }
	return segment
}

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

func reCalculateThroughputs(updateFilter func(string, string, float64), applyFilter func()) {
	if testInteger != 0 {
		switch testInteger {
		case 1:
		//	assignForTest()
			res := assignForTestSegment()
			if res {
				testInteger+=10
			}
		case 2:
			assignForTest2()
			testInteger++
		case 3:
			assignForTest3()
			testInteger++
		case 4:
			assignForTest4()
			testInteger++
		case 5:
			assignForTest5()
			testInteger++
		case 6:
			assignForTest6()
			testInteger++

		default:
		}
	}

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

func resetAllBandwidthSharingFlowMaxPlannedThroughput() {
        for _, flow := range bandwidthSharingFlowMap {
                resetBandwidthSharingFlowMaxPlannedThroughput(flow)
        }
}

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
			if nbEvaluatedFlowsLeft >= 1 {
				extra = (unusedBw - float64(nbEvaluatedFlowsLeft)*segment.MaxFairShareBwPerFlow) / float64(nbEvaluatedFlowsLeft)
			} else {
				extra = 0
			}

			//category 2 flows
			for _, flow := range flowsToEvaluate {

				if flow.PlannedThroughput == segment.MaxFairShareBwPerFlow {
//					if flow.CurrentThroughput < (segment.MaxFairShareBwPerFlow + extra) {
					if flow.CurrentThroughput < segment.MaxFairShareBwPerFlow {
					
						nbEvaluatedFlowsLeft--
						if nbEvaluatedFlowsLeft == 0 { //allocate everything of what is left
							flow.PlannedThroughput = unusedBw
							flow.PlannedUpperBound = unusedBw
							flow.PlannedLowerBound = unusedBw - segment.Config.TolerationThreshold
							//lower bound cannot be less than min threshold
							if flow.PlannedLowerBound < segment.Config.MinActivityThreshold {
								flow.PlannedLowerBound = segment.Config.MinActivityThreshold
							}
						} else {
							flow.PlannedThroughput = flow.CurrentThroughput + segment.Config.IncrementalStep
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
						flow.PlannedUpperBound = flow.PlannedThroughput - segment.Config.ActionUpperThreshold
						flow.PlannedLowerBound = flow.PlannedUpperBound - segment.Config.TolerationThreshold
						unusedBw -= flow.PlannedThroughput
					}
				}
			}
		}
	}
	//we allocate all the bw to active users and very low values to inactive ones if there is any residual
	//using a minimum value that is close but not exactly 0, since we use float operations and approximation may not lead to a perfect 0{
	if unusedBw >= 1 {

		for _, flow := range flowsToEvaluate {
			if flow.CurrentThroughput > segment.Config.MinActivityThreshold {
				flow.PlannedThroughput = segment.MaxFairShareBwPerFlow
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

func cleanUp() {
	//nothing allocated that should be cleared
}

func assignForTestSegment() bool {

        var segment BandwidthSharingSegment
        segment.Name = "S1"
        segment.MaxThroughput = 200.0
	nbFlow := 0
	flow := bandwidthSharingFlowMap["ue1-iperf:zone1-fog1-iperf"]
	if flow != nil {
	        segment.Flows = append(segment.Flows, flow)
		nbFlow++
	}
        flow = bandwidthSharingFlowMap["ue2-iperf:zone1-fog1-iperf"]
        if flow != nil {
        	segment.Flows = append(segment.Flows, flow)
		nbFlow++
	}
        flow = bandwidthSharingFlowMap["ue3-iperf:zone1-fog1-iperf"]
        if flow != nil {
	        segment.Flows = append(segment.Flows, flow)
		nbFlow++
	}
        flow = bandwidthSharingFlowMap["ue4-iperf:zone1-fog1-iperf"]
        if flow != nil {
	        segment.Flows = append(segment.Flows, flow)
		nbFlow++
	}
	if nbFlow == 4 {
	        assignDefaultSegmentConfigs(&segment)
	        //allSegments = append(allSegments, &segment)
		return true
	}
	return false
}

func assignForTest() {

	var b1 BandwidthSharingFlow
	b1.Name = "B1"
	b1.MaximumThroughput = 100.0
	b1.CurrentThroughput = 0.1
	bandwidthSharingFlowMap[b1.Name] = &b1

	var b2 BandwidthSharingFlow
	b2.Name = "B2"
	b2.MaximumThroughput = 100.0
	b2.CurrentThroughput = 0.1
	bandwidthSharingFlowMap[b2.Name] = &b2

	var b3 BandwidthSharingFlow
	b3.Name = "B3"
	b3.MaximumThroughput = 100.0
	b3.CurrentThroughput = 0.1
	bandwidthSharingFlowMap[b3.Name] = &b3

	var b4 BandwidthSharingFlow
	b4.Name = "B4"
	b4.MaximumThroughput = 100.0
	b4.CurrentThroughput = 0.1
	bandwidthSharingFlowMap[b4.Name] = &b4

	var segment BandwidthSharingSegment
	segment.Name = "S1"
	segment.MaxThroughput = 100.0
	segment.Flows = append(segment.Flows, &b1)
	segment.Flows = append(segment.Flows, &b2)
	segment.Flows = append(segment.Flows, &b3)
	segment.Flows = append(segment.Flows, &b4)
	assignDefaultSegmentConfigs(&segment)
	//allSegments = append(allSegments, &segment)

}

func assignDefaultSegmentConfigs(segment *BandwidthSharingSegment) {
        segment.Config = defaultConfigs
}

func initDefaultConfigAttributes() {
        defaultConfigs.MaxBwPerInactiveFlow = 2.0
        defaultConfigs.MinActivityThreshold = 0.3
        defaultConfigs.IncrementalStep = 3.0
        defaultConfigs.InactivityIncrementalStep = 1.0
        defaultConfigs.ActionUpperThreshold = 1.0
        defaultConfigs.TolerationThreshold = 4.0
}

func updateForTest(flowName string, currentThroughput float64) {

	bandwidthSharingFlowMap[flowName].CurrentThroughput = currentThroughput

}

func assignForTest2() {

	updateForTest("B1", 2.0)
}

/*
func assignForTest2() {

	updateForTest("B1", 10.0)
	updateForTest("B2", 20.0)
	updateForTest("B3", 26.0)
	updateForTest("B4", 50.0)
}
*/
func assignForTest3() {

	updateForTest("B1", 80.0)
	updateForTest("B2", 2.0)
	updateForTest("B3", 2.0)
}

func assignForTest4() {

	updateForTest("B1", 80.0)
	updateForTest("B2", 5.0)
	updateForTest("B3", 5.0)
}

func assignForTest5() {

	updateForTest("B1", 80.0)
	updateForTest("B2", 8.0)
	updateForTest("B3", 8.0)
}

func assignForTest6() {

	updateForTest("B1", 78.0)
	updateForTest("B2", 10.0)
	updateForTest("B3", 7.2)
}

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

func printBwFlowsNameFromList(list []*BandwidthSharingFlow) string {

	str := ""
        for _, flow := range list {
                str += flow.Name + "."
        }
	return str
}

func printBwFlows(segment *BandwidthSharingSegment) {

	log.Info("Flows on segment ", segment.Name)

	for _, flow := range segment.Flows {
		log.Info(printBwFlow(flow))
	}
}

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

func parseScenario(scenario ceModel.Scenario) {
        log.Debug("parseScenario")

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
						element.MaxThroughput = 100.0
                                		netElemList = append(netElemList, *element)
					}
				}
                        }
                }
        }

        //bandwidth sharing creation elements section
        //all the dummy processes meant for calculation can be ignored, others must havd bandwidth flows created between each of them
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

