package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExcelReader_GetData(t *testing.T) {
	msk := time.FixedZone("UTC+3", 3*60*60)

	type fields struct {
		filename string
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []ExcelData
		wantErr  bool
	}{
		{
			name:   "template",
			fields: fields{filename: "template.xlsx"},
			wantData: []ExcelData{
				{
					SenderContractorInn:              "9701037090",
					SenderContractorKpp:              "770101001",
					ExternalNumber:                   "KO231207_112",
					VehicleNo:                        "Н001СТ99",
					TrailerNo:                        "АН703447",
					DriverFio:                        "Абрамов Николай Владимирович",
					DriverPhone:                      "79153231010",
					LoadingAddress:                   "г. Краснодар, ул. Шувалова д. 5",
					LoadingPlannedTimeOfArrival:      time.Date(2023, 1, 22, 15, 0, 0, 0, msk),
					LoadingPlannedEndTimeOfArrival:   time.Date(2023, 1, 22, 21, 0, 0, 0, msk),
					UnloadingAddress:                 "г. Новороссийск ул. Портовая 8",
					UnloadingPlannedTimeOfArrival:    time.Date(2023, 1, 24, 06, 0, 0, 0, msk),
					UnloadingPlannedEndTimeOfArrival: time.Date(2023, 1, 24, 12, 0, 0, 0, msk),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewExcelReader(tt.fields.filename)
			require.NoError(t, err)

			gotData, err := r.GetData()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// fmt.Println(gotData)
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("GetData() gotData = %v, want %v", gotData, tt.wantData)
			}
			ord, err := json.Marshal(gotData[0].ToOrder("test"))
			require.NoError(t, err)
			fmt.Println(string(ord))
		})
	}
}
