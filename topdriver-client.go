package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	SearchType                = "3"
	WeekDays                  = "2%2C3%2C4%2C5%2C6%2C7%2C1"
	Timings                   = "-1"
	PageIndex                 = "1"
	ClassAvailabilityEndpoint = "https://topdriversignals.com/Customer/BTWScheduling/RefineSearch"
)

type TopdriverClient struct {
	client          *http.Client
	cookie          string
	InstructorId    string
	Timings         string
	PickupLocation  string
	DropoffLocation string
}

func NewTopdriverClient(cookie, instructorId, timings, pickupLocation, dropoffLocation string) *TopdriverClient {
	return &TopdriverClient{
		client:          &http.Client{},
		cookie:          cookie,
		InstructorId:    instructorId,
		Timings:         timings,
		PickupLocation:  pickupLocation,
		DropoffLocation: dropoffLocation,
	}
}

func (t *TopdriverClient) GetAvailableDays() (string, error) {
	data := url.Values{}
	data.Set("SearchType", SearchType)
	data.Set("InstructorID", t.InstructorId)
	data.Set("WeekDays", WeekDays)
	data.Set("Timings", Timings)
	data.Set("PageIndex", PageIndex)
	data.Set("PickupLocation", t.PickupLocation)
	data.Set("DropoffLocation", t.DropoffLocation)
	data.Set("SingleDate", "")
	data.Set("EC", "101")

	req, err := http.NewRequest("POST", ClassAvailabilityEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("cookie", t.cookie)

	res, err := t.client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func ParseRes(res string) (bool, error) {
	return strings.Contains(res, "No appointment is available for your search"), nil

	// if strings.Contains(res, "No appointment is available for your search") {
	// 	return []string{}, nil
	// }

	// Just for now say that if that above string is not present, lessons are available
}
