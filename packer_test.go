package gopack

import (
	"fmt"
	"reflect"
	"testing"
)

type result struct {
	packed   []*Bin
	unpacked []*Item
}

type testData struct {
	bins        []*Bin
	items       []*Item
	expectation result
}

func TestPack(t *testing.T) {
	testCase := []testData{
		{
			// first test
			bins: []*Bin{
				NewBin("Small Bin", 100, 100),
			},
			items: []*Item{
				NewItem("Item 1", 2, 2),
			},
			expectation: result{
				packed: []*Bin{
					{
						"Small Bin", 100, 100,
						[]*Item{
							{"Item 1", 2, 2, RotationType_WH, Pivot{0, 0}},
						},
					},
				},
				unpacked: []*Item{},
			},
		},
		//second test
		{
			// first test
			bins: []*Bin{
				NewBin("Small Bin", 100, 100),
			},
			items: []*Item{
				NewItem("Item 1", 2, 2),
				NewItem("Item 2", 10, 5),
				NewItem("Item 3", 20, 10),
				NewItem("Item 4", 5, 5),
			},
			expectation: result{
				packed: []*Bin{
					{
						"Small Bin", 100, 100,
						[]*Item{
							{"Item 3", 20, 10, RotationType_WH, Pivot{0, 0}},
							{"Item 2", 10, 5, RotationType_WH, Pivot{20, 0}},
							{"Item 4", 5, 5, RotationType_WH, Pivot{30, 0}},
							{"Item 1", 2, 2, RotationType_WH, Pivot{35, 0}},
						},
					},
				},
				unpacked: []*Item{},
			},
		},
	}
	for _, tc := range testCase {
		testPack(t, tc)
	}
}

func testPack(t *testing.T, td testData) {
	packer := NewPacker()
	packer.AddBin(td.bins...)
	packer.AddItem(td.items...)

	if err := packer.Pack(); err != nil {
		t.Fatalf("Got error: %v", err)
	}

	if !reflect.DeepEqual(packer.Bins, td.expectation.packed) {
		t.Errorf("\nGot:\n%+v\nwant:\n%+v", formatBins(packer.Bins, packer.UnfitItems), formatBins(td.expectation.packed, td.expectation.unpacked))
	}
}

func formatBins(bins []*Bin, unpacked []*Item) string {
	var s string
	for _, b := range bins {
		s += fmt.Sprintln(b)
		s += fmt.Sprintln(" packed items:")
		for _, i := range b.Items {
			s += fmt.Sprintln("  ", i)
		}

		s += fmt.Sprintln(" unpacked items:")
		for _, i := range unpacked {
			s += fmt.Sprintln("  ", i)
		}
	}
	return s
}
