package main

type gunm16 struct {
	gun
}

func newGunm16() igun {
	return &gunm16{
		gun: gun{
			name:  "m16",
			power: 5,
		},
	}
}
