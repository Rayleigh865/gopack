package main

import (
	"github.com/Rayleigh865/gopack"
	"log"
)

func main() {
	p := gopack.NewPacker()

	// Add bin
	p.AddBin(gopack.NewBin("Small Bin", 100, 100))

	// Add Items
	p.AddItem(gopack.NewItem("Item 1", 2, 2))
	p.AddItem(gopack.NewItem("Item 2", 10, 5))
	p.AddItem(gopack.NewItem("Item 3", 20, 10))
	p.AddItem(gopack.NewItem("Item 4", 5, 5))

	// Solve
	if err := p.Pack(); err != nil {
		log.Fatal(err)
	}

	//show results
	gopack.Display_packed(p.Bins)
}
