package main

import (
	"transport/lib/database"
)

// External (MTA) structs
// ==========================================================================
type MTAVehicleMonitoringResponse struct {
	Siri MTASiri
}

type MTASiri struct {
	ServiceDelivery MTAServiceDelivery
}

type MTAServiceDelivery struct {
	ResponseTimestamp         database.Timestamp
	VehicleMonitoringDelivery []MTAVehicleMonitoringDelivery
	SituationExchangeDelivery []MTASituationExchangeDelivery
}

// Vehicle Monitoring Structs

type MTAVehicleMonitoringDelivery struct {
	VehicleActivity   []MTAVehicleActivity
	ResponseTimestamp database.Timestamp
	ValidUntil        database.Timestamp
}

type MTAVehicleActivity struct {
	MonitoredVehicleJourney MTAMonitoredVehicleJourney
	RecordedAtTime          database.Timestamp
}

type MTAMonitoredVehicleJourney struct {
	LineRef                  string
	DirectionRef             int `json:"DirectionRef,string"`
	FramedVehicleJourneyRef  MTAFramedVehicleJourneyRef
	JourneyPatternRef        string
	PublishedLineName        []string
	OperatorRef              string
	OriginRef                string
	DestinationRef           string
	DestinationName          []string
	OriginAimedDepartureTime database.Timestamp
	SituationRef             []MTASituationRef
	Monitored                bool
	VehicleLocation          MTAVehicleLocation
	Bearing                  float64
	ProgressRate             string
	ProgressStatus           []string
	Occupancy                string
	VehicleRef               string
	BlockRef                 string
	MonitoredCall            MTAMonitoredCall
}

type MTAFramedVehicleJourneyRef struct {
	DataFrameRef           string
	DatedVehicleJourneyRef string
}

type MTASituationRef struct {
	SituationSimpleRef string
}

type MTAVehicleLocation struct {
	Longitude float64
	Latitude  float64
}

type MTAMonitoredCall struct {
	ExpectedArrivalTime   database.Timestamp
	ArrivalProximityText  string
	ExpectedDepartureTime database.Timestamp
	DistanceFromStop      int
	NumberOfStopsAway     int
	StopPointRef          string
	VisitNumber           int
	StopPointName         []string
}

// Situation Structs
type MTASituationExchangeDelivery struct {
	Situations MTASituation
}

type MTASituation struct {
	PtSituationElement []MTAPTSituationElement
}

type MTAPTSituationElement struct {
	PublicationWindow MTAPublicationWindow
	Severity          string
	Summary           []string
	Description       []string
	Affects           MTAAffected
	Consequences      MTAConsequence
	CreationTime      database.Timestamp
	SituationNumber   string
}

type MTAConsequence struct {
	Consequence []MTACondition
}

type MTACondition struct {
	Condition []string
}

type MTAPublicationWindow struct {
	StartTime database.Timestamp
	EndTime   database.Timestamp
}

type MTAAffected struct {
	VehicleJourneys MTAVehicleJourneys
}

type MTAVehicleJourneys struct {
	AffectedVehicleJourney []MTAAffectedVehicleJourney
}

type MTAAffectedVehicleJourney struct {
	LineRef      string
	DirectionRef int `json:"DirectionRef,string"`
}
