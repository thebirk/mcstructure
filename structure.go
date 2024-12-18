package structure

import (
	"io"
	"math"

	"github.com/Tnze/go-mc/nbt"
)

type rawPalette struct {
	Name       string         `nbt:"Name"`
	Properties map[string]any `nbt:"Properties,omitempty"`
}

type rawBlock struct {
	State    int32          `nbt:"state"`
	Position []int32        `nbt:"pos,list"`
	Nbt      map[string]any `nbt:"nbt,omitempty"`
}

type rawEntity struct {
	Position      []float64      `nbt:"pos"`
	BlockPosition []int32        `nbt:"blockPos"`
	Nbt           map[string]any `nbt:"nbt,omitempty"`
}

type rawStructure struct {
	DataVersion int32        `nbt:"DataVersion"`
	Size        []int32      `nbt:"size,list"`
	Palette     []rawPalette `nbt:"palette"`
	Blocks      []rawBlock   `nbt:"blocks"`
	Entities    []rawEntity  `nbt:"entities"`
}

type builderPalette struct {
	name       string
	properties map[string]any
}

type builderBlock struct {
	state   int
	x, y, z int
	nbt     map[string]any
}

type StructureBuilder struct {
	palette      []builderPalette
	blocks       []builderBlock
	maxX, minX   int32
	maxY, minY   int32
	maxZ, minZ   int32
	paletteCache map[string]int
}

func NewStructureBuilder() *StructureBuilder {
	return &StructureBuilder{
		paletteCache: make(map[string]int),
	}
}

func (b *StructureBuilder) recordBlock(x, y, z int) {
	if int32(x) < b.minX {
		b.minX = int32(x)
	} else if int32(x) > b.maxX {
		b.maxX = int32(x)
	}
	if int32(y) < b.minY {
		b.minY = int32(y)
	} else if int32(y) > b.maxY {
		b.maxY = int32(y)
	}
	if int32(z) < b.minZ {
		b.minZ = int32(z)
	} else if int32(z) > b.maxZ {
		b.maxZ = int32(z)
	}
}

func (b *StructureBuilder) PlaceEntity(nbt map[string]any, x, y, z float64, blockX, blockY, blockZ int32) *StructureBuilder {
	panic("not implemented")
}

func (b *StructureBuilder) PlaceBlockEntity(id string, x, y, z int, state map[string]any, nbt map[string]any) *StructureBuilder {
	return b.placeBlock(id, x, y, z, state, nbt)
}

func (b *StructureBuilder) PlaceBlock(id string, x, y, z int, state map[string]any) *StructureBuilder {
	return b.placeBlock(id, x, y, z, state, nil)
}

func (b *StructureBuilder) placeBlock(id string, x, y, z int, state map[string]any, nbt map[string]any) *StructureBuilder {
	var pIndex int
	if state == nil {
		var ok bool
		pIndex, ok = b.paletteCache[id]
		if !ok {
			b.palette = append(b.palette, builderPalette{
				name:       id,
				properties: state,
			})

			pIndex = len(b.palette) - 1
			b.paletteCache[id] = pIndex
		}
	} else {
		b.palette = append(b.palette, builderPalette{
			name:       id,
			properties: state,
		})

		pIndex = len(b.palette) - 1
	}

	b.blocks = append(b.blocks, builderBlock{
		state: pIndex,
		x:     x,
		y:     y,
		z:     z,
		nbt:   nbt,
	})
	b.recordBlock(x, y, z)

	return b
}

func (b *StructureBuilder) Write(w io.Writer) error {
	var m rawStructure
	m.DataVersion = 4189
	m.Size = []int32{
		int32(math.Abs(float64(b.maxX)-float64(b.minX))) + 1,
		int32(math.Abs(float64(b.maxY)-float64(b.minY))) + 1,
		int32(math.Abs(float64(b.maxZ)-float64(b.minZ))) + 1,
	}

	for _, p := range b.palette {
		m.Palette = append(m.Palette, rawPalette{
			Name:       p.name,
			Properties: p.properties,
		})
	}
	for _, b := range b.blocks {
		m.Blocks = append(m.Blocks, rawBlock{
			State:    int32(b.state),
			Position: []int32{int32(b.x), int32(b.y), int32(b.z)},
			Nbt:      nil,
		})
	}

	return nbt.NewEncoder(w).Encode(m, "")
}
