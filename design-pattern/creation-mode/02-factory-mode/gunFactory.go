package main

import "fmt"

func getGun(gunType string) (igun, error) {
	if gunType == "ak47" {
		return newGunak47(), nil
	}
	if gunType == "m16" {
		return newGunm16(), nil
	}
	return nil, fmt.Errorf("wrong gun type")
}
