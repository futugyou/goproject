package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/futugyou/extensions"

	"github.com/futugyou/platformservice/application/service"
	"github.com/futugyou/platformservice/options"
)

type VaultService struct {
	client *http.Client
	apiURL struct {
		GetVaultsByIDs string
		ShowVaultRaw   string
		CreateVault    string
	}
}

func NewVaultService(opts *options.Options) *VaultService {
	return &VaultService{
		client: &http.Client{Timeout: 10 * time.Second},
		apiURL: struct {
			GetVaultsByIDs string
			ShowVaultRaw   string
			CreateVault    string
		}{
			GetVaultsByIDs: opts.GetVaultsByIDs,
			ShowVaultRaw:   opts.ShowVaultRaw,
			CreateVault:    opts.CreateVault,
		},
	}
}

func (s *VaultService) doRequest(ctx context.Context, method, url string, body any, out any) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	token, ok := extensions.JWTFrom(ctx)
	if !ok || token == "" {
		return errors.New("missing jwt token in context")
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("vault service returned %d: %s", resp.StatusCode, string(data))
	}

	if out != nil {
		if strOut, ok := out.(*string); ok {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("read response body: %w", err)
			}
			*strOut = string(data)
			return nil
		}

		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

func (s *VaultService) GetVaultsByIDs(ctx context.Context, ids []string) ([]service.VaultView, error) {
	var result []service.VaultView
	err := s.doRequest(ctx, "POST", s.apiURL.GetVaultsByIDs, map[string][]string{"ids": ids}, &result)
	return result, err
}

func (s *VaultService) ShowVaultRawValue(ctx context.Context, vaultId string) (string, error) {
	url := fmt.Sprintf(s.apiURL.ShowVaultRaw, vaultId)
	var result string
	err := s.doRequest(ctx, "POST", url, nil, &result)
	return result, err
}

func (s *VaultService) CreateVault(ctx context.Context, aux service.CreateVaultRequest) (*service.VaultView, error) {
	var view service.VaultView
	err := s.doRequest(ctx, "POST", s.apiURL.CreateVault, aux, &view)
	if err != nil {
		return nil, err
	}
	return &view, nil
}
