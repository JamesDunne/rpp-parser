package rpp

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	//"math"
	//"encoding/binary"
	//"unsafe"
	//"encoding/binary"
	"unicode"
)

func TestParseRPP(t *testing.T) {
	f, err := os.Open("/Users/jimdunne/Desktop/band/20170517/R1320.RPP")
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
							fmt.Printf("%s%s\n", indent, fx.VST.Path)
							fmt.Printf("%s%2X\n", indent, fx.VST.Data)

							if fx.VST.Path == "reaeq.vst.dylib" {
								for i := 0; i < len(fx.VST.Data); i++ {
									if i&3 == 0 {
										fmt.Print(" ")
									}
									if !unicode.IsPrint(rune(fx.VST.Data[i])) {
										fmt.Print(".")
										continue
									}
									fmt.Printf("%c", rune(fx.VST.Data[i]))
									//ui := binary.LittleEndian.Uint32(data[i:i + 4])
									//fmt.Printf("%d, ", ui)
									//flip := binary.BigEndian.Uint32(fx.VST.Data[0][i:i+4])
									//f0 := *(*float32)(unsafe.Pointer(&flip))
									//fmt.Printf("  [%d] = %v\n", i, f0)
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
