// Copyright © 2018 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package anchore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func anchoreRequest(path string, bodyParams map[string]string, method string) ([]byte, error) {
	username := os.Getenv("ANCHORE_ENGINE_USERNAME")
	password := os.Getenv("ANCHORE_ENGINE_PASSWORD")
	anchoreEngineURL := os.Getenv("ANCHORE_ENGINE_URL")
	fullURL := anchoreEngineURL + path
	client := &http.Client{}

	bodyParamJson, err := json.Marshal(bodyParams)
	if err != nil {
		logrus.Fatal(err)
	}
	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(bodyParamJson))
	if err != nil {
		logrus.Fatal(err)
	}
	req.SetBasicAuth(username, password)
	logrus.WithFields(logrus.Fields{
		"url":        fullURL,
		"bodyParams": bodyParams,
	}).Info("Sending request")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to complete request to Anchore: %v", err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	logrus.WithFields(logrus.Fields{
		"response": string(bodyText),
	}).Info("Anchore Response Body")
	if err != nil {
		return nil, fmt.Errorf("failed to complete request to Anchore: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response from Anchore: %d", resp.StatusCode)
	}
	return bodyText, nil
}

func getStatus(digest string, tag string) bool {
	path := fmt.Sprintf("/v1/images/%s/check?history=false&detail=false&tag=%s", digest, tag)
	body, err := anchoreRequest(path, nil, "GET")
	if err != nil {
		logrus.Error(err)
		return false
	}
	var result []map[string]map[string][]SHAResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		logrus.Error(err)
		return false
	}

	resultIndex := fmt.Sprintf("docker.io/%s:latest", tag)
	return result[0][digest][resultIndex][0].Status == "pass"
}

func getImage(imageRef string) (Image, error) {
	params := map[string]string{"tag": imageRef}
	body, err := anchoreRequest("/v1/images?history=false", params, "GET")
	if err != nil {
		return Image{}, err
	}
	var images []Image
	err = json.Unmarshal(body, &images)
	if err != nil {
		return Image{}, fmt.Errorf("failed to unmarshal JSON from response: %v", err)
	}

	return images[0], nil
}
func getImageDigest(imageRef string) (string, error) {
	image, err := getImage(imageRef)
	if err != nil {
		return "", fmt.Errorf("failed to get image digest: %v", err)
	}
	return image.ImageDigest, nil
}

// AddImage add Image to Anchore
func AddImage(image string) error {
	params := map[string]string{"tag": image}
	_, err := anchoreRequest("/v1/images", params, "POST")
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"Image": image,
	}).Info("Added image to Anchore Engine")
	return nil
}

//CheckImage checking Image with Anchore
func CheckImage(image string) (string, string, string, bool) {
	imageParts := strings.Split(image, ":")
	tag := "latest"
	if len(imageParts) > 1 {
		tag = imageParts[1]
	}
	digest, err := getImageDigest(image)
	if err != nil {
		AddImage(image)
		return "", "", "", false
	}
	return imageParts[0], tag, digest, getStatus(digest, tag)
}
