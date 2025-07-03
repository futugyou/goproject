package service

import "time"

type Service struct {
	ID                   string                  `json:"id"`
	Name                 string                  `json:"name"`
	Description          string                  `json:"description"`
	TechStack            []string                `json:"tech_stack,omitempty"`
	Documents            []ServiceResource       `json:"documents,omitempty"`
	Platforms            []ServicePlatform       `json:"platforms,omitempty"`
	Deployments          []DeploymentEnvironment `json:"deployments,omitempty"`
	DependsOn            []string                `json:"depends_on,omitempty"`
	ExternalDependencies []ServiceDependency     `json:"external_dependencies,omitempty"`
	MockResponses        []APIMockEntry          `json:"mock_responses,omitempty"`
	FlowDiagram          *FlowDiagram            `json:"flow_diagram,omitempty"`
	ActivityStats        *ServiceActivityStats   `json:"activity_stats,omitempty"`
	Versions             []ServiceVersion        `json:"versions,omitempty"`
	DefaultVersion       string                  `json:"default_version"`
}

type FlowDiagram struct {
	Type            string `json:"type"`
	Description     string `json:"description,omitempty"`
	ResourceId      string `bson:"resource_id"`
	ResourceVersion int    `bson:"resource_version"`
}

type ServiceResource struct {
	Type            string `json:"type"`
	Description     string `json:"description,omitempty"`
	ResourceId      string `bson:"resource_id"`
	ResourceVersion int    `bson:"resource_version"`
}

type ServiceVersion struct {
	ID        string          `json:"id"`
	ServiceID string          `json:"service_id"`
	Version   string          `json:"version"`
	CreatedAt time.Time       `json:"created_at"`
	Changelog string          `json:"changelog"`
	Snapshot  ServiceSnapshot `json:"snapshot"`
}

type ServiceSnapshot struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	TechStack   []string          `json:"tech_stack"`
	Documents   []ServiceResource `json:"documents"`
}

type ServiceActivityStats struct {
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
	PlatformId  string            `bson:"platform_id"`
	ProjectId   string            `bson:"project_id"`
	Name        string            `json:"name"`
	Provider    string            `json:"provider"`
	Endpoint    string            `json:"endpoint"`
	Description string            `json:"description,omitempty"`
	Config      map[string]string `json:"config,omitempty"`
	Services    []string          `json:"services"` // aks, gke, etc
}

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
