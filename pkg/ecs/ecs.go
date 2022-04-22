package ecs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Retrieve ECS TaskID Metadata in V4 format
func GetTaskV4() (string, error) {
	metadataURL := os.Getenv("ECS_CONTAINER_METADATA_URI_V4")
	if metadataURL == "" {
		return "", errors.New("missing metadata uri v4")
	}

	response, err := http.Get(fmt.Sprintf("%s/task", metadataURL))
	if err != nil {
		return "", err
	}

	if response.Body == nil {
		return "", errors.New("failed to get task info: response is missing")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	response.Body.Close()

	taskInfo := &taskInfo{}
	err = json.Unmarshal(body, taskInfo)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal into task metadata v4: %w", err)
	}

	taskID, err := parseTaskID(taskInfo.TaskARN)
	if err != nil {
		return "", err
	}

	return taskID, nil
}

type taskInfo struct {
	TaskARN string
}

func parseTaskID(taskARN string) (string, error) {
	index := strings.LastIndex(taskARN, "/")
	if index == -1 {
		return "", errors.New("invalid Task ARN, can't find '/'")
	}
	adjustedIndex := index + 1
	if adjustedIndex >= len(taskARN) {
		return "", errors.New("invalid Task ARN, it can't end in '/'")
	}
	return taskARN[adjustedIndex:], nil
}
