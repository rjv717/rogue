package main

import (
	"fmt"
	"log"
	"math/rand"
	"rogue/arena"
	"rogue/creature"
	"time"

	gocui "github.com/awesome-gocui/gocui"
	_ "github.com/go-sql-driver/mysql"
)

type NamesTable struct {
	UID  int
	Name string
	Race string
}

var player *creature.Creature
var level *arena.ArenaType
var npcs []*creature.Creature

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	npcs = make([]*creature.Creature, rand.Intn(6)+1)

	level = arena.MakeArena(0, 0, 1)
	level.Build("rooms")
	player = creature.MakeCreature("human")
	for i, v := range npcs {
		if v == nil {
			npcs[i] = creature.MakeCreature("goblin")
		}
	}

	g, err := gocui.NewGui(gocui.OutputTrue, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, upArrow); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, downArrow); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, leftArrow); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, rightArrow); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func layout(g *gocui.Gui) error {
	var v = make([]*gocui.View, 4)
	var err error
	maxX, maxY := g.Size()

	if v[0], err = g.SetView("mapView", 0, 0, maxX-(maxX/5), maxY-8, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v[0].Editable = true
		v[0].Wrap = false
		v[0].Autoscroll = false
		level.View(v[0])
		for i := range npcs {
			npcs[i].View(v[0])
		}
		player.View(v[0])
	}

	if v[1], err = g.SetView("characterView", maxX-(maxX/5)+1, 0, maxX-1, 7, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		player.StatView(v[1])
	}
	if v[2], err = g.SetView("inventoryView", maxX-(maxX/5)+1, 8, maxX-1, maxY-8, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	if v[3], err = g.SetView("messageView", 0, maxY-7, maxX-1, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v[3].Editable = true
		v[3].Wrap = false
		v[3].Autoscroll = true
	}
	return nil
}

func upArrow(g *gocui.Gui, v *gocui.View) error {
	view, _ := g.SetCurrentView("messageView")
	if player.TryMove(level, arena.Up) {
		/* v, err := g.SetCurrentView("messageView")
		if err != nil {
			fmt.Fprintln(v, "Player moves Up")
		} */
	}
	fmt.Fprintln(view, "Player moves Up")
	//g.UpdateAsync(layout)
	return nil
}

func downArrow(g *gocui.Gui, v *gocui.View) error {
	player.TryMove(level, arena.Down)
	return nil
}

func leftArrow(g *gocui.Gui, v *gocui.View) error {
	player.TryMove(level, arena.Left)
	return nil
}

func rightArrow(g *gocui.Gui, v *gocui.View) error {
	player.TryMove(level, arena.Right)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
