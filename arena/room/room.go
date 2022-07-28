package room

import (
	"math/rand"
)

type setter interface {
	SetCell(string, int, int)
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

func (r *RoomType) IsCollision(other *RoomType) bool {

	if r.x1-1 <= other.x2+1 && r.x2+1 >= other.x1-1 && r.y1-1 <= other.y2+1 && r.y2+1 >= other.y1-1 {
		return true
	} else {
		return false
	}
}

func (r *RoomType) GetCenter() (int, int) {
	return (r.x2 + r.x1) / 2, (r.y2 - r.y1) / 2
}

func (r *RoomType) RandomPoint() (int, int) {

	x := r.x1 + rand.Intn(r.x2-r.x1) + 1
	y := r.y1 + rand.Intn(r.y2-r.y1) + 1

	return x, y
}

func (r *RoomType) SetRoom(a setter) {
	for i := r.x1; i <= r.x2; i++ {
		for j := r.y1; j <= r.y2; j++ {
			if j == r.y1 || j == r.y2 || i == r.x1 || i == r.x2 {
				a.SetCell("wall", i, j)
			} else {
				a.SetCell("floor", i, j)
			}
		}
	}
}
