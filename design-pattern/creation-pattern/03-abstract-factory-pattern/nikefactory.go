package main

type nikefactory struct {
}

type nikeclothing struct {
	clothing
}

type nikepants struct {
	pants
}

func (a *nikefactory) createClothing() iclothing {
	return &nikeclothing{
		clothing: clothing{
			color: "bule",
			size:  1,
		},
	}
}

func (a *nikefactory) createPants() ipants {
	return &nikepants{
		pants: pants{
			material: "wool",
		},
	}
}
