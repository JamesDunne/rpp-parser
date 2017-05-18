package rpp

import (
	"testing"
	"os"
	"bufio"
	"fmt"
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

	dump(project, "")
}

func dump(node *RPP, indent string) {
	fmt.Printf("%s%s\n", indent, node.Name)
	for _, child := range node.Children {
		dump(child, indent + "  ")
	}
}