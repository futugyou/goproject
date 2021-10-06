package main

type adsfactory struct {
}

type adsclothing struct {
	clothing
}

type adspants struct {
	pants
}

func (a *adsfactory) createClothing() iclothing {
	return &adsclothing{
		clothing: clothing{
			color: "red",
			size:  5,
		},
	}
}

func (a *adsfactory) createPants() ipants {
	return &adspants{
		pants: pants{
			material: "unkown",
		},
	}
}
