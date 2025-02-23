package gofile

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type ServerService service

func (s *ServerService) GetServer(ctx context.Context) (*ServerResponse, error) {
	req, err := http.NewRequest("GET", "https://api.gofile.io/servers", nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := s.client.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &ServerResponse{}
	if err = json.Unmarshal(all, response); err != nil {
		return nil, err
	}

	return response, nil
}

type ServerResponse struct {
	Status string     `json:"status"`
	Data   ServerData `json:"data"`
}

type ServerData struct {
	Servers        []Server `json:"servers"`
	ServersAllZone []Server `json:"serversAllZone"`
}

type Server struct {
	Name string `json:"name"`
	Zone string `json:"zone"`
}
