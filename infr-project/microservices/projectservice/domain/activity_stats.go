package domain

// Currently there is only github information.
// Considering that the services are archived, it may be necessary to find new dimensions.
type ActivityStats struct {
	Commits      int `json:"commits"`
	Pulls        int `json:"pulls"`
	Releases     int `json:"releases"`
	Contributors int `json:"contributors"`
}
