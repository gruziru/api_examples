package main

import (
	"time"
)

type TrackingOrder struct {
	Body         Body        `json:"body"`
	Header       Header      `json:"header"`
	ReceiverData interface{} `json:"receiverData"`
}

type Header struct {
	Type                  string    `json:"type"`
	Uuid                  string    `json:"uuid"`
	IsTest                bool      `json:"isTest"`
	Version               string    `json:"version"`
	CreatedAt             time.Time `json:"createdAt"`
	ManualSend            bool      `json:"manualSend"`
	SenderSystemCode      string    `json:"senderSystemCode"`
	SenderContractorInn   string    `json:"senderContractorInn"`
	SenderContractorKpp   string    `json:"senderContractorKpp"`
	TransportOperatorInn  string    `json:"transportOperatorInn"`
	ReceiverContractorInn string    `json:"receiverContractorInn"`
	ReceiverContractorKpp string    `json:"receiverContractorKpp"`
}

type Body struct {
	Route           []Route   `json:"route"`
	Driver          Driver    `json:"driver"`
	ClientInn       string    `json:"clientInn"`
	TrailerNo       string    `json:"trailerNo"`
	VehicleNo       string    `json:"vehicleNo"`
	ExternalId      string    `json:"externalId"`
	TrackingEndDt   time.Time `json:"trackingEndDt"`
	ExternalNumber  string    `json:"externalNumber"`
	TrackingStartDt time.Time `json:"trackingStartDt"`
}

type Route struct {
	Type                    string    `json:"type"`
	Location                Location  `json:"location"`
	TargetRadius            int       `json:"targetRadius"`
	ParkingRadius           int       `json:"parkingRadius"`
	PlannedTimeOfArrival    time.Time `json:"plannedTimeOfArrival"`
	PlannedEndTimeOfArrival time.Time `json:"plannedEndTimeOfArrival"`
}

type Location struct {
	Address      string `json:"address"`
	ExternalCode string `json:"externalCode"`
}

type Driver struct {
	Fio   string `json:"fio"`
	Phone string `json:"phone"`
}
