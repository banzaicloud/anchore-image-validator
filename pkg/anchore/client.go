package anchore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func anchoreRequest(path string, bodyParams map[string]string, method string) ([]byte, error) {
	username := os.Getenv("ANCHORE_ENGINE_USERNAME")
	password := os.Getenv("ANCHORE_ENGINE_PASSWORD")
	anchoreEngineURL := os.Getenv("ANCHORE_ENGINE_URL")
	fullURL := anchoreEngineURL + path
	client := &http.Client{}

	bodyParamJson, err := json.Marshal(bodyParams)
	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(bodyParamJson))
	if err != nil {
		glog.Fatal(err)
	}
	req.SetBasicAuth(username, password)
	glog.Infof("Sending request to %s, with params %s", fullURL, bodyParams)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to complete request to Anchore: %v", err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	glog.Info("Anchore Response Body: " + string(bodyText))
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
		glog.Error(err)
		return false
	}
	var result []map[string]map[string][]SHAResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		glog.Error(err)
		return false
	}

	// Is this the easiest way to get this info?
	resultIndex := fmt.Sprintf("docker.io/%s:latest", tag)
	return result[0][digest][resultIndex][0].Status == "pass"
}

func getImage(imageRef string) (Image, error) {
	// Tag or repo??
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

func AddImage(image string) error {
	params := map[string]string{"tag": image}
	_, err := anchoreRequest("/v1/images", params, "POST")
	if err != nil {
		return err
	}
	glog.Infof("Added image to Anchore Engine: %s", image)
	return nil
}

func CheckImage(image string) bool {
	imageParts := strings.Split(image, ":")
	tag := "latest"
	if len(imageParts) > 1 {
		tag = imageParts[1]
	}
	digest, err := getImageDigest(image)
	if err != nil {
		AddImage(image)
		return false
	}
	return getStatus(digest, tag)
}
