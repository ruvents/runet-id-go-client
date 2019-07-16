package models

/*
	{
	  "DeviceNumber": "214",
	  "HallId": 286,
	  "Id": 739,
	  "EventId": null,
	  "PlaceId": 1112,
	  "Title": "A6",
	  "Description": null,
	  "SuccessCount": 0,
	  "ErrorCount": 0,
	  "CreationTime": "2019-07-14 08:37:10"
	}
*/
type PaperlessOstrovAssociation struct {
	DeviceNumber string `json:"DeviceNumber"`
	HallID       uint32 `json:"HallId"`
	HallName     string `json:"Title"`
}
