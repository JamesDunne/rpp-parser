package rpp

import "math"

type ReaEQBand struct {
	Frequency float64
	Gain float64
	Q float64
}

func (band ReaEQBand) Bandwidth() float64 {
	q_sqr := band.Q * band.Q
	bw := math.Log2((2.0*q_sqr+1.0)/(2.0*q_sqr) + math.Sqrt(math.Pow((2.0*q_sqr+1.0)/q_sqr, 2.0)/4.0-1.0))
	return bw
}

type ReaEQ struct {
	Bands []ReaEQBand
}

type VST struct {
	Name string
	Path string
	Data []byte

	ReaEQ *ReaEQ
}

type JS struct {
	Name       string
	Path       string
	Parameters []*float64
}

type FX struct {
	Bypass bool
	VST    *VST
	JS     *JS
}

type FXChain struct {
	FX []*FX
}

type Track struct {
	Name        string
	Volume      float64
	Pan         float64
	InvertPhase bool
	FXChain     *FXChain
}

type Project struct {
	Name   string
	Tracks []*Track
}
