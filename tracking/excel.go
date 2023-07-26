package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

const (
	startFromRow       = 3
	trackingType       = "createTrackingOrder"
	version            = "1.0"
	gruziruINN         = "7729632882"
	gruziruKPP         = "771401001"
	loadingRouteType   = "loading"
	unloadingRouteType = "unloading"
	radius             = 3000
	leadingTime        = -2
	trailingTime       = 24
)

type ExcelData struct {
	SenderContractorInn              string `validate:"required,numeric,len=10"`
	SenderContractorKpp              string `validate:"required,numeric,len=9"`
	ExternalNumber                   string `validate:"required"`
	VehicleNo                        string `validate:"required"`
	TrailerNo                        string
	DriverFio                        string `validate:"required"`
	DriverPhone                      string `validate:"required,numeric,startswith=79,len=11"`
	LoadingAddress                   string `validate:"required"`
	LoadingPlannedTimeOfArrival      time.Time
	LoadingPlannedEndTimeOfArrival   time.Time
	UnloadingAddress                 string `validate:"required"`
	UnloadingPlannedTimeOfArrival    time.Time
	UnloadingPlannedEndTimeOfArrival time.Time
}

func (d ExcelData) ToOrder(from string) TrackingOrder {
	loadingExtCode := md5.Sum([]byte(d.LoadingAddress))
	unloadingExtCode := md5.Sum([]byte(d.UnloadingAddress))
	return TrackingOrder{
		Body: Body{
			Route: []Route{
				{
					Type: loadingRouteType,
					Location: Location{
						Address:      d.LoadingAddress,
						ExternalCode: hex.EncodeToString(loadingExtCode[:]),
					},
					TargetRadius:            radius,
					ParkingRadius:           radius,
					PlannedTimeOfArrival:    d.LoadingPlannedTimeOfArrival,
					PlannedEndTimeOfArrival: d.LoadingPlannedEndTimeOfArrival,
				},
				{
					Type: unloadingRouteType,
					Location: Location{
						Address:      d.UnloadingAddress,
						ExternalCode: hex.EncodeToString(unloadingExtCode[:]),
					},
					TargetRadius:            radius,
					ParkingRadius:           radius,
					PlannedTimeOfArrival:    d.UnloadingPlannedTimeOfArrival,
					PlannedEndTimeOfArrival: d.UnloadingPlannedEndTimeOfArrival,
				},
			},
			Driver: Driver{
				Fio:   d.DriverFio,
				Phone: d.DriverPhone,
			},
			ClientInn:       d.SenderContractorInn,
			TrailerNo:       d.TrailerNo,
			VehicleNo:       d.VehicleNo,
			ExternalId:      d.ExternalNumber,
			TrackingEndDt:   d.LoadingPlannedEndTimeOfArrival.Add(trailingTime * time.Hour),
			ExternalNumber:  d.ExternalNumber,
			TrackingStartDt: d.LoadingPlannedTimeOfArrival.Add(leadingTime * time.Hour),
		},
		Header: Header{
			Type:                  trackingType,
			Uuid:                  uuid.New().String(),
			IsTest:                false,
			Version:               version,
			CreatedAt:             time.Now(),
			ManualSend:            false,
			SenderSystemCode:      from,
			SenderContractorInn:   d.SenderContractorInn,
			SenderContractorKpp:   d.SenderContractorKpp,
			TransportOperatorInn:  gruziruINN,
			ReceiverContractorInn: gruziruINN,
			ReceiverContractorKpp: gruziruKPP,
		},
	}
}

type ExcelReader struct {
	in *excelize.File
}

func NewExcelReader(filename string) (xls *ExcelReader, err error) {
	xls = &ExcelReader{}
	xls.in, err = excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("open excel file: %w", err)
	}

	defer func() {
		err = xls.in.Close()
		if err != nil {
			xls = nil
			err = fmt.Errorf("close excel file: %w", err)
		}
	}()

	return
}

func (r ExcelReader) GetData() (ords []ExcelData, err error) {
	if r.in == nil {
		return nil, fmt.Errorf("excel file is nil")
	}

	ords = make([]ExcelData, 0, 4)
	for _, sheet := range r.in.GetSheetMap() {
		ord, err := r.processSheet(sheet)
		if err != nil {
			return nil, fmt.Errorf("process sheet %s: %w", sheet, err)
		}
		ords = append(ords, ord...)
	}
	return ords, nil
}

func (r ExcelReader) processSheet(sheet string) ([]ExcelData, error) {
	rows, err := r.in.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}

	ords := make([]ExcelData, 0, len(rows)-startFromRow)
	for rowNum := startFromRow - 1; rowNum < len(rows); rowNum++ {
		row := rows[rowNum]
		to, err := r.readRow(row)
		if err != nil {
			return nil, fmt.Errorf("read row %d: %w", rowNum, err)
		}
		ords = append(ords, to)
	}
	return ords, nil
}

func (r ExcelReader) readRow(row []string) (data ExcelData, err error) {
	msk := time.FixedZone("UTC+3", 3*60*60)

	data.SenderContractorInn = row[1]
	data.SenderContractorKpp = row[2]
	data.ExternalNumber = row[3]
	data.VehicleNo = row[4]
	data.TrailerNo = row[5]
	data.DriverFio = row[6]
	data.DriverPhone = row[7]
	data.LoadingAddress = row[8]
	data.LoadingPlannedTimeOfArrival, err = time.ParseInLocation("02/01/06 15:04", row[9], msk)
	data.LoadingPlannedEndTimeOfArrival, err = time.ParseInLocation("02/01/06 15:04", row[10], msk)
	data.UnloadingAddress = row[11]
	data.UnloadingPlannedTimeOfArrival, err = time.ParseInLocation("02/01/06 15:04", row[12], msk)
	data.UnloadingPlannedEndTimeOfArrival, err = time.ParseInLocation("02/01/06 15:04", row[13], msk)

	val := validator.New()
	if err = val.Struct(data); err != nil {
		return ExcelData{}, fmt.Errorf("validate: %w", err)
	}

	return
}
