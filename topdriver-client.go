package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	LoginUserEndpoint         = "https://topdriversignals.com/Customer/StudentLogin/LoginUser"
)

type TopdriverClient struct {
	client          *http.Client
	username        string
	password        string
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

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if result["SessionExpired"] == true {
		authToken, err := t.GetNewAuthToken()
		if err != nil {
			return "", err
		}
		t.cookie = authToken
		log.Println("fetched new auth token, updating for next try.")
		return "", nil
	}

	stringRes := fmt.Sprintf("%v", result["AvaialbleCalendarResult"])
	// log.Println("result: ", result)
	// log.Println("result: ", result["AvailableCalendarResult"])
	return stringRes, nil
}

func (t *TopdriverClient) GetNewAuthToken() (string, error) {
	log.Println("Fetching new auth token")
	// data := url.Values{}
	// data.Set("u", t.username)
	// data.Set("p", t.password)
	// data.Set("chk", WeekDays)

	// req, err := http.NewRequest("POST", LoginUserEndpoint, strings.NewReader(data.Encode()))
	// if err != nil {
	// 	return "", err
	// }

	// req.Header.Set("accept", "*/*")
	// req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	// // req.Header.Set("sec-fetch-dest", "empty")
	// // req.Header.Set("sec-fetch-mode", "cors")
	// // req.Header.Set("sec-fetch-site", "same-origin")

	// res, err := t.client.Do(req)
	// if err != nil {
	// 	return "", err
	// }

	// log.Println(res.Cookies())
	// // log.Println(res.Header["set-cookie"])
	url := "https://topdriversignals.com/Customer/StudentLogin/LoginUser"
	method := "POST"

	payload := strings.NewReader("u=Rho264796&p=191261226&chk=false")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	return res.Cookies()[0].Value, nil
}

func ParseRes(res string) (bool, error) {
	return strings.Contains(res, `<input type="hidden" id="hdnAvailableDates" value=`), nil
	// return false, nil
	// if strings.Contains(res, "No appointment is available for your search") {
	// 	return []string{}, nil
	// }

	// Just for now say that if that above string is not present, lessons are available
}
