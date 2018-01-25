package rpp

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestParseRPP(t *testing.T) {
	f, err := os.Open("R1425.RPP")
	if err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(f)
	project, err := ParseRPP(reader)
	if err != nil {
		t.Fatal(err)
	}

	for _, track := range project.Tracks {
		fmt.Printf("%s\n", track.Name)
		if track.Name == "Print" {
			continue
		}
		fmt.Printf("  Volume = %f\n", track.Volume)
		fmt.Printf("  Pan = %f\n", track.Pan)

		if track.FXChain != nil {
			fmt.Printf("  FX\n")
			for _, fx := range track.FXChain.FX {
				if fx.VST != nil {
					data := fx.VST.Data
					fmt.Printf("    %s\n", fx.VST.Path)
					if fx.VST.ReaEQ != nil {
						for _, band := range fx.VST.ReaEQ.Bands {
							//fmt.Printf("      [%d] freq=%7.1f Hz, gain=%6.2f dB, bw=%5.3f, q=%6.3f\n", i, band.Frequency, band.Gain, band.Bandwidth, band.Q())
							fmt.Printf("      {\"freq\":%f, \"gain\":%f, \"q\":%f},\n", band.Frequency, band.Gain, band.Q())
						}
					} else {
						fmt.Printf("      %2X\n", data)
					}
				} else if fx.JS != nil {
					fmt.Printf("    %s\n", fx.JS.Path)
				}
			}
		}
	}
}
