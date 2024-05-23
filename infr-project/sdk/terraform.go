package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	tfcAPIBaseURL = "https://app.terraform.io/api/v2"
	organization  = "futugyousuzu"
	workspace     = "infr-project"
)

type TerraformClient struct {
	token string
}

func NewTerraformClient(token string) *TerraformClient {
	return &TerraformClient{
		token: token,
	}
}

func (s *TerraformClient) CheckWorkspace(name string) (string, error) {
	url := fmt.Sprintf("%s/organizations/%s/workspaces", tfcAPIBaseURL, organization)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to list workspaces: %s", body)
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	for _, data := range result["data"].([]interface{}) {
		workspace := data.(map[string]interface{})
		if workspace["attributes"].(map[string]interface{})["name"].(string) == name {
			return workspace["id"].(string), nil
		}
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"attributes": map[string]interface{}{
				"name": name,
			},
			"type": "workspaces",
		},
	}
	jsonData, _ := json.Marshal(data)
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create workspace: %s", body)
	}

	json.Unmarshal(body, &result)
	return result["data"].(map[string]interface{})["id"].(string), nil
}
