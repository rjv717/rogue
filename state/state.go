package state

import (
	"fmt"
	"math/rand"
	"rogue/arena"
	"rogue/creature"
	kb "rogue/keyboard"
	"time"
)

type State string

const (
	Prerun       State = "pre-run condition"
	Player       State = "player turn"
	NonPlayer    State = "npc turn"
	gobbletygook State = "junk entry"
)

func DoState() {
	s := Prerun
	var player *creature.Creature
	var level *arena.ArenaType

	rand.Seed(time.Now().UTC().UnixNano())
	npcs := make([]*creature.Creature, rand.Intn(6)+1)

	for {
		switch s {
		case Prerun:
			level = arena.MakeArena(0, 0, 1)
			player = creature.MakeCreature("human")
			for i, v := range npcs {
				if v == nil {
					npcs[i] = creature.MakeCreature("goblin")
				}
			}
			s = Player
		case Player:
			ascii, _, _ := kb.GetChar() // ascii, keycode, err

			switch ascii {
			case 97: // a - left
				player.TryMove(level, arena.Left)
			case 100: // d - right
				player.TryMove(level, arena.Right)
			case 112:
				panic("panic and exit.")
			case 113:
				fmt.Println("Goodbye.")
				return
			case 115: // s - down
				player.TryMove(level, arena.Down)
			case 119: // w - up
				player.TryMove(level, arena.Up)
			default:

			}
			fmt.Println(player.GetName() + " has moved.")
			player.DebugPrint()
			s = NonPlayer
		case NonPlayer:
			for i, v := range npcs {
				if v != nil {
					fmt.Println(npcs[i].GetName() + " moves.")
					npcs[i].DebugPrint()
				}
			}
			s = Player
		default:
			s = Player
		}

	}
}
