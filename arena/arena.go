package arena

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
)

type Direction int64

const (
	Left Direction = iota
	Right
	Up
	Down
)

const standardSize = 80

type ArenaType struct {
	x, y, z, xSize, ySize int
	field                 [][]*Cell
	rooms                 []*RoomType
}

func MakeArena(x, y, z int) *ArenaType {
	aArena := new(ArenaType)

	aArena.x = x
	aArena.y = y
	aArena.z = z
	aArena.xSize = standardSize
	aArena.ySize = standardSize

	aArena.field = make([][]*Cell, standardSize)
	for i := range aArena.field {
		aArena.field[i] = make([]*Cell, standardSize)
	}
	numRooms := rand.Intn(standardSize/10) + (standardSize / 10)
	aArena.rooms = make([]*RoomType, numRooms)

	for i := 0; i < numRooms; i++ {

	Make_room:

		aRoom := MakeRoom(aArena.xSize-1, aArena.ySize-1)
		for j := 0; j < numRooms; j++ {
			if aArena.rooms[j] != nil {
				if aArena.rooms[j].IsCollision(aRoom) {
					i--
					if i < 0 {
						i = 0
					}
					goto Make_room
				}
			} else {
				goto Next
			}
		}

	Next:

		aArena.rooms[i] = aRoom
	}

	for _, v := range aArena.rooms {
		v.SetRoom(aArena)
	}
	for _, v := range aArena.rooms {
		randomRoom := rand.Intn(numRooms)
		room := aArena.rooms[randomRoom]
		v.makeCorridor(room, aArena)
	}

	return aArena
}

/* func (a *ArenaType) Build(s string) {

	for range a.rooms {
	Retry_room:
		aRoom := MakeRoom(a.xSize-1, a.ySize-1)
		for i, v := range a.rooms {
			if v == nil {
				a.rooms[i] = aRoom
				break
			} else if aRoom.IsCollision(a.rooms[i]) {
				goto Retry_room
			}
		}
	}

	for i := range a.rooms {
		a.rooms[i].SetRoom(a)
	}

}
*/

func (a *ArenaType) IsLocationPassable(x, y int) bool {
	if a.field[x][y] != nil {
		return a.field[x][y].IsPassable()
	}
	return false
}

func (a *ArenaType) GetTotalRooms() int {
	return len(a.rooms)
}

func (a *ArenaType) GetRoom(index int) *RoomType {
	if index >= len(a.rooms) {
		return a.rooms[0]
	} else {
		return a.rooms[index]
	}
}

func (a *ArenaType) getcell(x, y int) (*Cell, error) {
	if x >= 0 && y >= 0 && x < len(a.field[0]) && y < len(a.field) {
		return a.field[x][y], nil
	}
	return nil, errors.New("out of bounds")
}

func (a *ArenaType) SetCell(typeOfCell string, x, y int) {
	var ct CellType

	switch typeOfCell {
	case "floor":
		ct = Floor
	case "upstairs":
		ct = UpStairs
	case "downstairs":
		ct = DownStairs
	default:
		ct = Wall
	}

	if a.field[x][y] == nil {
		a.field[x][y] = MakeCell(ct, x, y, a.z)
	}
}

func (a *ArenaType) MapView(x, y, w, h int) string {

	var buffer string
	xOffset := w / 2
	yOffset := h / 2

	for j := y - yOffset + 1; j <= y-yOffset+h; j++ {
		for i := x - xOffset + 1; i <= x-xOffset+w; i++ {
			if i >= 0 && i < len(a.field) && j >= 0 && j < len(a.field[0]) {
				if a.field[i][j] != nil {
					buffer = fmt.Sprintf("%s%c", buffer, a.field[i][j].Display())
				} else {
					buffer = fmt.Sprintf("%s%c", buffer, ' ')
				}
			} else {
				buffer = fmt.Sprintf("%s%c", buffer, 'X')
			}
		}
	}

	return buffer
}

type RoomType struct {
	x1, y1, x2, y2 int
}

func MakeRoom(xRange, yRange int) *RoomType {
	aRoom := new(RoomType)

	xSize := rand.Intn(6) + 4
	ySize := rand.Intn(6) + 4

	aRoom.x1 = rand.Intn(xRange - xSize - 1)
	aRoom.y1 = rand.Intn(yRange - ySize - 1)
	aRoom.x2 = aRoom.x1 + xSize
	aRoom.y2 = aRoom.y1 + ySize

	return aRoom
}

func (r *RoomType) makeCorridor(r1 *RoomType, a *ArenaType) {
	var x, y, ox2, oy1 int
	var err error
	var c *Cell

	x1, y1 := r.GetCenter()
	x2, y2 := r1.GetCenter()
	ox2 = x2
	oy1 = y1

	if x2 < x1 {
		temp := x2
		x2 = x1
		x1 = temp
	}
	if y2 < y1 {
		temp := y2
		y2 = y1
		y1 = temp
	}
	for x = x1; x <= x2; x++ {
		c, err = a.getcell(x, oy1)
		if err == nil {
			if c == nil || !c.IsPassable() {
				a.field[x][oy1] = MakeCell(Floor, x, oy1, a.z)
			}
		}
		c, err = a.getcell(x, oy1-1)
		if err == nil {
			if c == nil {
				a.field[x][oy1-1] = MakeCell(Wall, x, oy1-1, a.z)
			}
		}
		c, err = a.getcell(x, oy1+1)
		if err == nil {
			if c == nil {
				a.field[x][oy1+1] = MakeCell(Wall, x, oy1+1, a.z)
			}
		}
	}
	c, err = a.getcell(ox2-1, oy1-1)
	if err == nil {
		if c == nil {
			a.field[ox2-1][oy1-1] = MakeCell(Wall, ox2-1, oy1-1, a.z)
		}
	}
	c, err = a.getcell(ox2-1, oy1+1)
	if err == nil {
		if c == nil {
			a.field[ox2-1][oy1+1] = MakeCell(Wall, ox2-1, oy1+1, a.z)
		}
	}
	c, err = a.getcell(ox2+1, oy1-1)
	if err == nil {
		if c == nil {
			a.field[ox2+1][oy1-1] = MakeCell(Wall, ox2+1, oy1-1, a.z)
		}
	}
	c, err = a.getcell(ox2+1, oy1+1)
	if err == nil {
		if c == nil {
			a.field[ox2+1][oy1+1] = MakeCell(Wall, ox2+1, oy1+1, a.z)
		}
	}
	for y = y1; y <= y2; y++ {
		c, err = a.getcell(ox2, y)
		if err == nil {
			if c == nil || !c.IsPassable() {
				a.field[ox2][y] = MakeCell(Floor, ox2, y, a.z)
			}
		}
		c, err = a.getcell(ox2-1, y)
		if err == nil {
			if c == nil {
				a.field[ox2-1][y] = MakeCell(Wall, ox2-1, y, a.z)
			}
		}
		c, err = a.getcell(ox2+1, y)
		if err == nil {
			if c == nil {
				a.field[ox2+1][y] = MakeCell(Wall, ox2+1, y, a.z)
			}
		}
	}
}

func (r *RoomType) IsCollision(other *RoomType) bool {

	if r.x1-1 <= other.x2+1 && r.x2+1 >= other.x1-1 && r.y1-1 <= other.y2+1 && r.y2+1 >= other.y1-1 {
		return true
	} else {
		return false
	}
}

func (r *RoomType) GetCenter() (int, int) {
	return (r.x2 + r.x1) / 2, (r.y2 + r.y1) / 2
}

func (r *RoomType) GetRandomPoint() (int, int) {

	x := r.x1 + rand.Intn(r.x2-r.x1-1) + 1
	y := r.y1 + rand.Intn(r.y2-r.y1-1) + 1

	return x, y
}

func (r *RoomType) SetRoom(a *ArenaType) {
	for i := r.y1; i <= r.y2; i++ {
		for j := r.x1; j <= r.x2; j++ {
			if j == r.x1 || j == r.x2 || i == r.y1 || i == r.y2 {
				a.SetCell("wall", j, i)
			} else {
				a.SetCell("floor", j, i)
			}
		}
	}
}

type CellType int64

const (
	Wall CellType = iota
	Floor
	UpStairs
	DownStairs
)

type Cell struct {
	cellType   CellType
	x, y, z    int
	glyph      rune
	isPassable bool
	isMutable  bool
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

	var glyph []uint8
	var passable bool

	res := db.QueryRow("SELECT tbl_displayables.glyph FROM (tbl_displayables INNER JOIN tbl_celltype ON tbl_displayables.UID = tbl_celltype.display) WHERE tbl_celltype.description = ?", n)

	err2 := res.Scan(&glyph)
	if err2 != nil {
		log.Fatal(err2)
	}

	res = db.QueryRow("SELECT is_passable FROM tbl_celltype WHERE description = ?", n)
	err2 = res.Scan(&passable)

	if err2 != nil {
		log.Fatal(err2)
	}

	aCell := &Cell{cellType, x, y, z, rune(glyph[0]), passable, true}

	return aCell
}

func (c *Cell) Fix() {
	c.isMutable = false
}

func (c *Cell) IsFixed() bool {
	return c.isMutable
}

func (c *Cell) IsPassable() bool {
	return c.isPassable
}

func (c *Cell) Display() rune {
	return c.glyph
}
