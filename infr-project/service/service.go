package service

import (
	"time"

	"github.com/futugyou/domaincore/domain"
)

type Service struct {
	domain.Aggregate
	Name                 string                  `json:"name"`
	Description          string                  `json:"description"`
	TechStack            []string                `json:"tech_stack,omitempty"`
	Documents            []ServiceResource       `json:"documents,omitempty"`
	Platforms            []ServicePlatform       `json:"platforms,omitempty"`
	Deployments          []DeploymentEnvironment `json:"deployments,omitempty"`
	DependsOn            []string                `json:"depends_on,omitempty"`
	ExternalDependencies []ServiceDependency     `json:"external_dependencies,omitempty"`
	MockResponses        []APIMockEntry          `json:"mock_responses,omitempty"`
	ArchDiagram          *ArchDiagram            `json:"arch_diagram,omitempty"`
	ActivityStats        *ActivityStats          `json:"activity_stats,omitempty"`
	Versions             []ServiceVersion        `json:"versions,omitempty"`
	DefaultVersion       string                  `json:"default_version"`
}

func (r Service) AggregateName() string {
	return "services"
}

type ArchDiagram struct {
	ResourceId      string `bson:"resource_id"`
	ResourceVersion int    `bson:"resource_version"`
}

type ServiceResource struct {
	ResourceId      string `bson:"resource_id"`
	ResourceVersion int    `bson:"resource_version"`
	Type            string `json:"type"` // FlowChart er doc eg.
}

type ServiceVersion struct {
	ID        string          `json:"id"` // service_id+'@'+version eg. XXXXX@1.0.2
	ServiceID string          `json:"service_id"`
	Version   string          `json:"version"`
	CreatedAt time.Time       `json:"created_at"`
	Changelog string          `json:"changelog"`
	Snapshot  ServiceSnapshot `json:"snapshot"`
}

// When the following four properties change significantly, it can choose to create a new `ServiceVersion`.
// especially `ArchDiagram` and `Documents`
type ServiceSnapshot struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ArchDiagram *ArchDiagram      `json:"arch_diagram,omitempty"`
	Documents   []ServiceResource `json:"documents,omitempty"`
}

// Currently there is only github information.
// Considering that the services are archived, it may be necessary to find new dimensions.
type ActivityStats struct {
	Commits      int `json:"commits"`
	Pulls        int `json:"pulls"`
	Releases     int `json:"releases"`
	Contributors int `json:"contributors"`
}

type DeploymentEnvironment struct {
	Name           string            `json:"name"`
	Platform       ServicePlatform   `json:"platform"`
	Endpoints      []string          `json:"endpoints,omitempty"`
	Region         string            `json:"region,omitempty"`
	Configurations map[string]string `json:"config,omitempty"`
}

type ServicePlatform struct {
	PlatformId string   `bson:"platform_id"` // eg. github id
	ProjectId  string   `bson:"project_id"`  // eg. github responsitory id/name
	Services   []string `json:"services"`    // responsitory, aks, gke, etc
}

//	{
//		"name": "Redis",
//		"type": "Cache",
//		"provider": "RedisLabs",
//		"endpoint": "redis://cache:6379",
//		"config": {
//		  "ttl": "30s"
//		}
//	}
//
// If it to display third-party dependencies that are not declared in the `Platform`, you can use `ServiceDependency`
type ServiceDependency struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Provider    string            `json:"provider"`
	Endpoint    string            `json:"endpoint"`
	Description string            `json:"description,omitempty"`
	Config      map[string]string `json:"config,omitempty"`
}

type APIMockEntry struct {
	ID        string            `json:"id"`
	ServiceID string            `json:"service_id"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers,omitempty"`
	Path      string            `json:"path"`
	Schema    string            `json:"schema"`
	Response  string            `json:"response"`
	Mode      string            `json:"mode"` // manual | auto-llm
}
