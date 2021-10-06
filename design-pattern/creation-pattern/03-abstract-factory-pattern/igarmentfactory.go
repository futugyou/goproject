package main

type igarmentfactory interface {
	createClothing() iclothing
	createPants() ipants
}

func newigarmentfactory(brand string) (igarmentfactory, error) {
	if brand == "ads" {
		return &adsfactory{}, nil
	}
	if brand == "nike" {
		return &nikefactory{}, nil
	}
	return nil, nil
}
