package arena

import (
	rooms "rogue/arena/room"
)

func (a *ArenaType) Build(s string) {

	for range a.rooms {
	Retry_room:
		aRoom := rooms.MakeRoom(a.xSize-1, a.ySize-1)
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
