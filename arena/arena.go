package arena

import (
	"fmt"
	"math/rand"
	cell "rogue/arena/cell"
	room "rogue/arena/room"

	"github.com/awesome-gocui/gocui"
)

type Direction int64

const (
	Left Direction = iota
	Right
	Up
	Down
)

const standardSize = 160

type ArenaType struct {
	x, y, z, xSize, ySize int
	field                 [][]*cell.Cell
	rooms                 []*room.RoomType
}

func MakeArena(x, y, z int) *ArenaType {
	aArena := new(ArenaType)

	aArena.x = x
	aArena.y = y
	aArena.z = z
	aArena.xSize = standardSize
	aArena.ySize = standardSize

	aArena.field = make([][]*cell.Cell, standardSize)
	for i := range aArena.field {
		aArena.field[i] = make([]*cell.Cell, standardSize)
	}
	aArena.rooms = make([]*room.RoomType, rand.Intn(24)+4)

	return aArena
}

func (a *ArenaType) IsLocationPassable(x, y int) bool {
	if a.field[x][y] != nil {
		return a.field[x][y].IsPassable()
	}
	return false
}

func (a *ArenaType) SetCell(typeOfCell string, x, y int) {
	var ct cell.CellType

	switch typeOfCell {
	case "floor":
		ct = cell.Floor
	case "upstairs":
		ct = cell.UpStairs
	case "downstairs":
		ct = cell.DownStairs
	default:
		ct = cell.Wall
	}

	if a.field[x][y] != nil {
		a.field[x][y] = cell.MakeCell(ct, x, y, a.z)
	}
}

func (a *ArenaType) View(v *gocui.View) {

	v.Clear()
	for i := range a.field {
		for j := range a.field[i] {
			if a.field[i][j] == nil {
				fmt.Fprintf(v, ".")
			} else {
				fmt.Fprintf(v, "%c", a.field[i][j].Display())
			}
		}
		fmt.Fprintf(v, "\n")
	}
}
