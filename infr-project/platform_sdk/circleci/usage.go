package circleci

func (s *CircleciClient) CreateUsageExport(org_id string, start string, end string, shared_org_ids []string) (*UsageExport, error) {
	path := "/organizations/" + org_id + "/usage_export_job"
	request := struct {
		Start        string   `json:"start"`
		End          string   `json:"end"`
		SharedOrgIds []string `json:"shared_org_ids"`
	}{
		Start:        start,
		End:          end,
		SharedOrgIds: shared_org_ids,
	}
	result := &UsageExport{}
	if err := s.http.Post(path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) GetUsageExport(org_id string, usage_export_job_id string) (*UsageExport, error) {
	path := "/organizations/" + org_id + "/usage_export_job/" + usage_export_job_id

	result := &UsageExport{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

type UsageExport struct {
	UsageExportJobID string   `json:"usage_export_job_id"`
	State            string   `json:"state"`
	Start            string   `json:"start"`
	End              string   `json:"end"`
	DownloadUrls     []string `json:"download_urls"`
	Message          *string  `json:"message,omitempty"`
}
