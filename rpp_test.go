package rpp

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"testing"
)

func TestParseRPP(t *testing.T) {
	f, err := os.Open("/Users/jimdunne/Desktop/band/20170517/R1322.RPP")
	if err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(f)
	project, err := ParseRPP(reader)
	if err != nil {
		t.Fatal(err)
	}

	(func(project *Project, indent string) {
		fmt.Printf("%s%s\n", indent, project.Name)
		for _, track := range project.Tracks {
			fmt.Printf("%s%s\n", indent, track.Name)
			if track.Name == "Print" {
				continue
			}

			(func(track *Track, indent string) {
				if track.FXChain != nil {
					for _, fx := range track.FXChain.FX {
						if fx.VST != nil {
							data := fx.VST.Data
							fmt.Printf("%s%s\n", indent, fx.VST.Path)
							fmt.Printf("%s%2X\n", indent, data)

							if fx.VST.Path == "reaeq.vst.dylib" {
								z := 0
								_ = binary.LittleEndian.Uint32(data[z : z+4])
								z += 4
								bands := binary.LittleEndian.Uint32(data[z : z+4])
								_ = bands
								z += 4
								_ = binary.LittleEndian.Uint32(data[z : z+4])
								z += 4
								_ = binary.LittleEndian.Uint32(data[z : z+4])
								z += 4

								//fmt.Printf("%s%2X\n", indent, data[z:])
								for band := uint32(0); band < bands; band++ {
									freq := math.Float64frombits(binary.LittleEndian.Uint64(data[z : z+8]))
									z += 8
									pct := math.Float64frombits(binary.LittleEndian.Uint64(data[z : z+8]))
									gain := math.Log10(pct) * 20
									z += 8
									q := math.Float64frombits(binary.LittleEndian.Uint64(data[z : z+8]))
									z += 8
									fmt.Printf("freq=%6.1f, gain=%5.2f dB, q=%5.3f\n", freq, gain, q)
									z += 9
								}
								fmt.Printf("\n")
							}
						} else if fx.JS != nil {
							fmt.Printf("%s%s\n", indent, fx.JS.Path)
						}
					}
				}
			})(track, indent+"  ")
		}
	})(project, "")
}
