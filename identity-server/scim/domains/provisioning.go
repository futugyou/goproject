package domains

import "time"

type ProvisioningConfigurationRecordTypes uint32

const (
	ProvisioningConfigurationRecordTypesSTRING  ProvisioningConfigurationRecordTypes = 0
	ProvisioningConfigurationRecordTypesCOMPLEX ProvisioningConfigurationRecordTypes = 1
)

type ProvisioningConfigurationTypes uint32

const (
	ProvisioningConfigurationTypesAPI ProvisioningConfigurationTypes = 0
)

type ProvisioningConfigurationHistoryStatus uint32

const (
	ProvisioningConfigurationHistoryStatusFINISHED  ProvisioningConfigurationHistoryStatus = 0
	ProvisioningConfigurationHistoryStatusEXCEPTION ProvisioningConfigurationHistoryStatus = 1
)

type ProvisioningConfigurationRecord struct {
	Name         string
	Type         ProvisioningConfigurationRecordTypes
	IsArray      bool
	ValuesString []string
	Values       []ProvisioningConfigurationRecord
}

type ProvisioningConfigurationHistory struct {
	RepresentationId          string
	RepresentationVersion     int
	Description               string
	WorkflowInstanceId        string
	WorkflowId                string
	ExecutionDateTime         time.Time
	Exception                 string
	Status                    ProvisioningConfigurationHistoryStatus
	ProvisioningConfiguration *ProvisioningConfiguration
}

func ProvisioningConfigurationHistoryComplete(representationId, description, workflowInstanceId, workflowId string, version int) *ProvisioningConfigurationHistory {
	return &ProvisioningConfigurationHistory{
		RepresentationId:      representationId,
		RepresentationVersion: version,
		Description:           description,
		WorkflowInstanceId:    workflowInstanceId,
		WorkflowId:            workflowId,
		ExecutionDateTime:     time.Now().UTC(),
		Status:                ProvisioningConfigurationHistoryStatusFINISHED,
	}
}

func ProvisioningConfigurationHistoryError(representationId, description, exception string, version int) *ProvisioningConfigurationHistory {
	return &ProvisioningConfigurationHistory{
		RepresentationId:      representationId,
		RepresentationVersion: version,
		Description:           description,
		Exception:             exception,
		ExecutionDateTime:     time.Now().UTC(),
		Status:                ProvisioningConfigurationHistoryStatusEXCEPTION,
	}
}

type ProvisioningConfiguration struct {
	Id             string
	Type           ProvisioningConfigurationTypes
	ResourceType   string
	UpdateDateTime time.Time
	Records        []ProvisioningConfigurationRecord
	HistoryLst     []ProvisioningConfigurationHistory
}

func (p *ProvisioningConfiguration) Update(records []ProvisioningConfigurationRecord) {
	p.Records = []ProvisioningConfigurationRecord{}
	for _, record := range records {
		p.Records = append(p.Records, record)
	}
	p.UpdateDateTime = time.Now().UTC()
}

func (p *ProvisioningConfiguration) Complete(representationId, description, workflowInstanceId, workflowId string, version int) {
	p.UpdateDateTime = time.Now().UTC()
	p.HistoryLst = append(p.HistoryLst, *ProvisioningConfigurationHistoryComplete(representationId, description, workflowInstanceId, workflowId, version))
}

func (p *ProvisioningConfiguration) Error(representationId, description, exception string, version int) {
	p.UpdateDateTime = time.Now().UTC()
	p.HistoryLst = append(p.HistoryLst, *ProvisioningConfigurationHistoryError(representationId, description, exception, version))
}

func (p *ProvisioningConfiguration) IsRepresentationProvisioned(representationId string, version int) bool {
	for i := 0; i < len(p.HistoryLst); i++ {
		if p.HistoryLst[i].RepresentationId == representationId && p.HistoryLst[i].RepresentationVersion == version && p.HistoryLst[i].Status == ProvisioningConfigurationHistoryStatusFINISHED {
			return true
		}
	}

	return false
}
