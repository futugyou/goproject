package main

type gunak47 struct {
	gun
}

func newGunak47() igun {
	return &gunak47{
		gun: gun{
			name:  "ak47",
			power: 4,
		},
	}
}
