package gopack

import (
	"fmt"
	"math"
	"sort"
)

type Bin struct {
	Name   string
	Width  int64
	Height int64

	Items []*Item
}

type BinSlice []*Bin

func (bs BinSlice) Len() int {
	return len(bs)
}

func (bs BinSlice) Less(i, j int) bool {
	return bs[i].GetVolume() < bs[j].GetVolume()
}

func (bs BinSlice) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

func NewBin(name string, w, h int64) *Bin {
	return &Bin{
		Name:   name,
		Width:  w,
		Height: h,
		Items:  make([]*Item, 0),
	}
}

func (b *Bin) GetName() string {
	return b.Name
}

func (b *Bin) GetWidth() int64 {
	return b.Width
}

func (b *Bin) GetHeight() int64 {
	return b.Height
}

func (b *Bin) GetVolume() int64 {
	return b.Width * b.Height
}

func (b *Bin) PutItem(item *Item, p Pivot) (fit bool) {
	item.Position = p
	for i := 0; i < 2; i++ {
		item.RotationType = RotationType(i)
		d := item.GetDimension()
		if b.GetWidth() < p[0]+d[0] || b.GetHeight() < p[1]+d[1] {
			continue
		}
		fit = true

		for _, ib := range b.Items {
			if ib.Intersect(item) {
				fit = false
				break
			}
		}

		if fit {
			b.Items = append(b.Items, item)
		}

		return
	}

	return
}

func (b *Bin) String() string {
	return fmt.Sprintf("%s(%vx%v)", b.GetName(), b.GetWidth(), b.GetHeight())
}

type RotationType int

const (
	RotationType_WH RotationType = iota
	RotationType_HW
)

var RotationTypeStrings = [...]string{
	"RotationType_WH (w,h)",
	"RotationType_HW (h,w)",
}

func (rt RotationType) String() string {
	return RotationTypeStrings[rt]
}

type Axis int

const (
	WidthAxis Axis = iota
	HeightAxis
)

type Pivot [2]int64
type Dimension [2]int64

func (pv Pivot) String() string {
	return fmt.Sprintf("%v,%v", pv[0], pv[1])
}

var startPosition = Pivot{0, 0}

type Item struct {
	Name   string
	Width  int64
	Height int64

	RotationType RotationType
	Position     Pivot
}

type ItemSlice []*Item

func (is ItemSlice) Len() int {
	return len(is)
}

func (is ItemSlice) Less(i, j int) bool {
	return is[i].GetVolume() > is[j].GetVolume()
}

func (is ItemSlice) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

func NewItem(name string, w, h int64) *Item {
	return &Item{
		Name:   name,
		Width:  w,
		Height: h,
	}
}

func (i *Item) GetName() string {
	return i.Name
}

func (i *Item) GetWidth() int64 {
	return i.Width
}

func (i *Item) GetHeight() int64 {
	return i.Height
}

func (i *Item) GetVolume() int64 {
	return i.Width * i.Height
}

func (i *Item) GetDimension() (d Dimension) {
	switch i.RotationType {
	case RotationType_WH:
		d = Dimension{i.GetWidth(), i.GetHeight()}
	case RotationType_HW:
		d = Dimension{i.GetHeight(), i.GetWidth()}
	}
	return
}

func (i *Item) Intersect(i2 *Item) bool {
	return rectIntersect(i, i2, WidthAxis, HeightAxis)
}

func rectIntersect(i1, i2 *Item, x, y Axis) bool {
	d1 := i1.GetDimension()
	d2 := i2.GetDimension()

	cx1 := i1.Position[x] + d1[x]/2
	cy1 := i1.Position[y] + d1[y]/2
	cx2 := i2.Position[x] + d2[x]/2
	cy2 := i2.Position[y] + d2[y]/2

	ix := math.Max(float64(cx1), float64(cx2)) - math.Min(float64(cx1), float64(cx2))
	iy := math.Max(float64(cy1), float64(cy2)) - math.Min(float64(cy1), float64(cy2))

	return ix < float64((d1[x]+d2[x])/2) && iy < float64((d1[y]+d2[y])/2)
}

func (i *Item) String() string {
	return fmt.Sprintf("%s(%vx%v) pos(%s) rt(%s)", i.GetName(), i.GetWidth(), i.GetHeight(), i.Position, i.RotationType)
}

type Packer struct {
	Bins       []*Bin
	Items      []*Item
	UnfitItems []*Item
}

func NewPacker() *Packer {
	return &Packer{
		Bins:       make([]*Bin, 0),
		Items:      make([]*Item, 0),
		UnfitItems: make([]*Item, 0),
	}
}

func (p *Packer) AddBin(bins ...*Bin) {
	p.Bins = append(p.Bins, bins...)
}

func (p *Packer) AddItem(items ...*Item) {
	p.Items = append(p.Items, items...)
}

func (p *Packer) Pack() error {
	sort.Sort(BinSlice(p.Bins))
	sort.Sort(ItemSlice(p.Items))

	for len(p.Items) > 0 {
		bin := p.FindFittedBin(p.Items[0])
		if bin == nil {
			p.unfitItem()
			continue
		}

		p.Items = p.packToBin(bin, p.Items)
	}

	return nil
}

func (p *Packer) unfitItem() {
	if len(p.Items) == 0 {
		return
	}
	p.UnfitItems = append(p.UnfitItems, p.Items[0])
	p.Items = p.Items[1:]
}

func (p *Packer) packToBin(b *Bin, items []*Item) (unpacked []*Item) {
	if !b.PutItem(items[0], startPosition) {

		if b2 := p.getBiggerBinThan(b); b2 != nil {
			return p.packToBin(b2, items)
		}

		return p.Items
	}

	for _, i := range items[1:] {
		var fitted bool
	lookup:
		for pt := 0; pt < 2; pt++ {
			for _, ib := range b.Items {
				var pv Pivot
				switch Axis(pt) {
				case WidthAxis:
					pv = Pivot{ib.Position[0] + ib.GetWidth(), ib.Position[1]}
				case HeightAxis:
					pv = Pivot{ib.Position[0], ib.Position[1] + ib.GetHeight()}
				}

				if b.PutItem(i, pv) {
					fitted = true
					break lookup
				}
			}
		}

		if !fitted {
			for b2 := p.getBiggerBinThan(b); b2 != nil; b2 = p.getBiggerBinThan(b) {
				left := p.packToBin(b2, append(b2.Items, i))
				if len(left) == 0 {
					b = b2
					fitted = true
					break
				}
			}

			if !fitted {
				unpacked = append(unpacked, i)
			}
		}
	}

	return
}

func (p *Packer) getBiggerBinThan(b *Bin) *Bin {
	v := b.GetVolume()
	for _, b2 := range p.Bins {
		if b2.GetVolume() > v {
			return b2
		}
	}
	return nil
}

func (p *Packer) FindFittedBin(i *Item) *Bin {
	for _, b := range p.Bins {
		if !b.PutItem(i, startPosition) {
			continue
		}

		if len(b.Items) == 1 && b.Items[0] == i {
			b.Items = []*Item{}
		}

		return b
	}
	return nil
}

func Display_packed(bins []*Bin) {
	for _, bin := range bins {
		fmt.Println(bin)
		fmt.Println(" packed items:")
		for _, i := range bin.Items {
			fmt.Println("  ", i)
		}
	}
}
