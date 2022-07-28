package cell

import (
	"database/sql"
	"log"
)

type CellType int64

const (
	Wall CellType = iota
	Floor
	UpStairs
	DownStairs
)

type Cell struct {
	cellType  CellType
	x, y, z   int
	glyph     uint8
	isMutable bool
}

func MakeCell(cellType CellType, x, y, z int) *Cell {

	db, err := sql.Open("mysql", "rodney:Akhen@t0n@tcp(127.0.0.1:3306)/roguelike")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var n string

	switch cellType {
	case Wall:
		n = "dungeon wall"
	case Floor:
		n = "dungeon floor"
	case UpStairs:
		n = "dungeon up staircase"
	case DownStairs:
		n = "dungeon down staircase"
	}

	var glyph uint8

	res := db.QueryRow("SELECT tbl_displayable.glyph FROM (tbl_displayable INNER JOIN tbl_celltype ON tbl_diplayable.UID = tbl_celltype.display) WHERE tbl_celltype.description = ?", n)

	err2 := res.Scan(&glyph)
	if err2 != nil {
		log.Fatal(err2)
	}

	aCell := Cell{cellType, x, y, z, glyph, true}

	return &aCell
}

func (c *Cell) Fix() {
	c.isMutable = false
}

func (c *Cell) IsFixed() bool {
	return c.isMutable
}

func (c *Cell) IsPassable() bool {
	var passable int

	db, err := sql.Open("mysql", "rodney:Akhen@t0n@tcp(127.0.0.1:3306)/roguelike")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	res := db.QueryRow("SELECT is_passable FROM tbl_celltype WHERE UID = ?", c.cellType)
	err2 := res.Scan(&passable)

	if err2 != nil {
		log.Fatal(err2)
	}

	if passable == 1 {
		return true
	}
	return true
}

func (c *Cell) Display() uint8 {
	return c.glyph
}
