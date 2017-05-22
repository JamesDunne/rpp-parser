package rpp

import (
	"bufio"
	"fmt"
	"os"
	"testing"
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
			(func(track *Track, indent string) {
				if track.FXChain != nil {
					for _, fx := range track.FXChain.FX {
						if fx.VST != nil {
							fmt.Printf("%s%s\n", indent, fx.VST.Path)
						} else if fx.JS != nil {
							fmt.Printf("%s%s\n", indent, fx.JS.Path)
						}
					}
				}
			})(track, indent + "  ")
		}
	})(project, "")
}
