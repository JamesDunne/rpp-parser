package rpp

type VST struct {
	Name string
	Path string
	Data []byte
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
