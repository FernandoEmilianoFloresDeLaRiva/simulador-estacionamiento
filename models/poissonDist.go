package models

import "gonum.org/v1/gonum/stat/distuv"

type Poisson struct {
}

func NewPoissont() *Poisson {
	return &Poisson{}
}

func (pd *Poisson) CreateDist(lambda float64) float64 {
	poisson := distuv.Poisson{Lambda: lambda, Src: nil}
	return poisson.Rand()
}
