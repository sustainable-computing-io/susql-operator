/*
Copyright 2023, 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
)

func queryCarbonIntensity(url string, location string, filter string, conv2J float64) (float64, error) {
	queryUrl := fmt.Sprintf(url, location)

	response, err := http.Get(queryUrl)
	if err != nil {
		return 0.0, fmt.Errorf("queryCarbonIntensity: %w\nURL=%s", err, queryUrl)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0.0, fmt.Errorf("queryCarbonIntensity: %w\nURL=%s\nresponse=%s", err, queryUrl, string(responseData))
	}

	length := gjson.Get(string(responseData), "#").Int() - 1
	index := strconv.Itoa(int(length))

	newFilter := index + "." + filter
	carbonIntensityString := gjson.Get(string(responseData), newFilter).String()

	carbonIntensityFloat, err := strconv.ParseFloat(carbonIntensityString, 64)
	if err != nil {
		return 0.0, fmt.Errorf("queryCarbonIntensity: %w\nURL=%s\nresponse=%s\nfilter=%s\nresult=%s", err, queryUrl, string(responseData), newFilter, carbonIntensityString)
	}

	// return nil error since no error
	return carbonIntensityFloat * conv2J, nil
}

func querySimpleCarbonIntensity(url string, location string, filter string, conv2J float64) (float64, error) {
	queryUrl := fmt.Sprintf(url, location)

	response, err := http.Get(fmt.Sprintf(url, location))
	if err != nil {
		return 0.0, fmt.Errorf("querySimpleCarbonIntensity: %w\nURL=%s", err, queryUrl)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0.0, fmt.Errorf("queryCarbonSimpleIntensity: %w\nURL=%s\nresponse=%s", err, queryUrl, string(responseData))
	}

	carbonIntensityString := gjson.Get(string(responseData), filter).String()

	carbonIntensityFloat, err := strconv.ParseFloat(carbonIntensityString, 64)
	if err != nil {
		return 0.0, fmt.Errorf("queryCarbonSimpleIntensity: %w\nURL=%s\nresponse=%s\nfilter=%s\nresult=%s", err, queryUrl, string(responseData), filter, carbonIntensityString)
	}

	// return nil error since no error
	return carbonIntensityFloat * conv2J, nil
}
