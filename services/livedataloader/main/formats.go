package main

// Internal Format Structs
// ==========================================================================
type VehicleJourney struct {
	LineRef                  string
	DirectionRef             int
	TripID                   string
	PublishedLineName        string
	OperatorRef              string
	OriginRef                string
	DestinationRef           string
	OriginAimedDepartureTime Timestamp
	SituationRef             []string
	Longitude                float64
	Latitude                 float64
	ProgressRate             string
	Occupancy                string
	VehicleRef               string
	ExpectedArrivalTime      Timestamp
	ExpectedDepartureTime    Timestamp
	DistanceFromStop         int
	NumberOfStopsAway        int
	StopPointRef             string
	Timestamp                Timestamp
}

// External (MTA) structs
// ==========================================================================
type MTAVehicleMonitoringResponse struct {
	Siri MTASiri
}

type MTASiri struct {
	ServiceDelivery MTAServiceDelivery
}

type MTAServiceDelivery struct {
	ResponseTimestamp         Timestamp
	VehicleMonitoringDelivery []MTAVehicleMonitoringDelivery
	SituationExchangeDelivery []MTASituationExchangeDelivery
}

// Vehicle Monitoring Structs

type MTAVehicleMonitoringDelivery struct {
	VehicleActivity   []MTAVehicleActivity
	ResponseTimestamp Timestamp
	ValidUntil        Timestamp
}

type MTAVehicleActivity struct {
	MonitoredVehicleJourney MTAMonitoredVehicleJourney
	RecordedAtTime          Timestamp
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
	OriginAimedDepartureTime Timestamp
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
	ExpectedArrivalTime   Timestamp
	ArrivalProximityText  string
	ExpectedDepartureTime Timestamp
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
	CreationTime      Timestamp
	SituationNumber   string
}

type MTAConsequence struct {
	Consequence []MTACondition
}

type MTACondition struct {
	Condition []string
}

type MTAPublicationWindow struct {
	StartTime Timestamp
	EndTime   Timestamp
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
