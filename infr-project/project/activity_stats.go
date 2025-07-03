package project

type ActivityStats struct {
	Commits      int `json:"commits"`
	Pulls        int `json:"pulls"`
	Releases     int `json:"releases"`
	Contributors int `json:"contributors"`
}
