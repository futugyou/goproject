package gofile

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AccountService service

func (s *AccountService) GetAccountID(ctx context.Context) (*GetAccountIDResponse, error) {
	req, err := http.NewRequest("GET", "https://api.gofile.io/accounts/getid", nil)
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

	response := &GetAccountIDResponse{}
	if err = json.Unmarshal(all, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *AccountService) GetAccountByID(ctx context.Context, accountID string) (*GetAccountResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.gofile.io/accounts/%s", accountID), nil)
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

	response := &GetAccountResponse{}
	if err = json.Unmarshal(all, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *AccountService) ResetToken(ctx context.Context, accountID string) (*ResetTokenResponse, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.gofile.io/accounts/%s/resettoken", accountID), nil)
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

	response := &ResetTokenResponse{}
	if err = json.Unmarshal(all, response); err != nil {
		return nil, err
	}

	return response, nil
}

type GetAccountIDResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type GetAccountResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type ResetTokenResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
