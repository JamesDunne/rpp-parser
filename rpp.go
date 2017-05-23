package rpp

import "math"

type ReaEQBand struct {
	Frequency float64
	Gain      float64
	Bandwidth float64
}

func (band ReaEQBand) Q() float64 {
	// SQRT(POWER(2,bw)) / (POWER(2,bw)-1)
	bw_2 := math.Pow(2.0, band.Bandwidth)
	return math.Sqrt(bw_2) / (bw_2 - 1.0)
}

func VolumeToDB(p float64) float64 {
	return math.Log10(p) * 20.0
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
