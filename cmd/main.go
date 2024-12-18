package main

import (
	"compress/gzip"
	"os"

	structure "github.com/thebirk/mcstructure"
	"github.com/thebirk/mcstructure/blocks"
)

func main() {
	f, err := os.OpenFile("test.nbt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b := structure.NewStructureBuilder()
	/*
		b.PlaceBlock(blocks.Dirt.NamespacedName, 0, 0, 0, nil)
		b.PlaceBlock("minecraft:grass_block", 0, 1, 0, map[string]any{
			"snowy": "false",
		})
		b.PlaceBlock("minecraft:dandelion", 0, 2, 0, nil)
		b.PlaceBlock("minecraft:dirt", 3, 0, 0, nil)
		b.PlaceBlock(blocks.DiamondBlock.NamespacedName, 3, 1, 0, nil)
	*/
	b.PlaceBlock(blocks.RedGlazedTerracotta.NamespacedName, 0, -1, 0, nil)
	b.PlaceBlock(blocks.StoneBricks.NamespacedName, 1, -1, 0, nil)
	b.PlaceBlock(blocks.StoneBricks.NamespacedName, -1, -1, 0, nil)
	b.PlaceBlock(blocks.StoneBricks.NamespacedName, 0, -1, 1, nil)
	b.PlaceBlock(blocks.StoneBricks.NamespacedName, 0, -1, -1, nil)

	gw := gzip.NewWriter(f)
	defer gw.Close()
	if err := b.Write(gw); err != nil {
		panic(err)
	}
}
