package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"rogue/arena"
	"rogue/creature"
	"time"

	tcell "github.com/gdamore/tcell/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jroimartin/gocui"
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
	room := level.GetRoom(0)
	x, y := room.GetCenter()
	player = creature.MakeCreature("human", x, y)
	room = level.GetRoom(rand.Intn(level.GetTotalRooms()-1) + 1)
	for i, v := range npcs {
		if v == nil {
			x, y = room.GetRandomPoint()
			npcs[i] = creature.MakeCreature("goblin", x, y)
		}
	}

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default text style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	// Clear screen
	s.Clear()

	// Draw the Initial Screen
	drawScreen(s)

	quit := func() {
		s.Fini()
		os.Exit(0)
	}
	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			drawScreen(s)
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			}
		}
	}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 || r == '\n' {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func drawScreen(s tcell.Screen) {

	// Inventory View : (w-21), 9 - (w-21), (h-9)
	// Message View : 1, (h-7) - (w-2), (h-2)

	w, h := s.Size()
	w--
	h--
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	s.Clear()

	for i := 1; i < (w - 22); i++ {
		s.SetContent(i, 0, tcell.RuneHLine, nil, boxStyle)
		s.SetContent(i, h-7, tcell.RuneHLine, nil, boxStyle)
		s.SetContent(i, h, tcell.RuneHLine, nil, boxStyle)
	}
	for i := (w - 21); i < w; i++ {
		s.SetContent(i, 0, tcell.RuneHLine, nil, boxStyle)
		s.SetContent(i, 7, tcell.RuneHLine, nil, boxStyle)
		s.SetContent(i, h-7, tcell.RuneHLine, nil, boxStyle)
		s.SetContent(i, h, tcell.RuneHLine, nil, boxStyle)
	}
	for i := 1; i < (h - 7); i++ {
		s.SetContent(0, i, tcell.RuneVLine, nil, boxStyle)
		s.SetContent(w-22, i, tcell.RuneVLine, nil, boxStyle)
		s.SetContent(w, i, tcell.RuneVLine, nil, boxStyle)
	}
	for i := (h - 6); i < h; i++ {
		s.SetContent(0, i, tcell.RuneVLine, nil, boxStyle)
		s.SetContent(w, i, tcell.RuneVLine, nil, boxStyle)
	}
	s.SetContent(0, 0, tcell.RuneULCorner, nil, boxStyle)
	s.SetContent(0, h-7, tcell.RuneLTee, nil, boxStyle)
	s.SetContent(0, h, tcell.RuneLLCorner, nil, boxStyle)

	s.SetContent(w-22, 0, tcell.RuneTTee, nil, boxStyle)
	s.SetContent(w-22, 7, tcell.RuneLTee, nil, boxStyle)
	s.SetContent(w-22, h-7, tcell.RuneBTee, nil, boxStyle)
	s.SetContent(w-22, h, tcell.RuneHLine, nil, boxStyle)

	s.SetContent(w, 0, tcell.RuneURCorner, nil, boxStyle)
	s.SetContent(w, 7, tcell.RuneRTee, nil, boxStyle)
	s.SetContent(w, h-7, tcell.RuneRTee, nil, boxStyle)
	s.SetContent(w, h, tcell.RuneLRCorner, nil, boxStyle)

	// Player View : (w-21), 1 - (w-1), 7  ???

	buffer := player.StatView()
	drawText(s, w-21, 1, w-1, 7, boxStyle, buffer)

	// Map View: 1, 1 - (w-23), (h-8)  ???
	px, py, _ := player.GetPosition()
	//xoffset := -px + ((w - 24) / 2)
	buffer = level.MapView(px, py, (w - 23), (h - 8))
	drawText(s, 1, 1, (w - 23), (h - 8), boxStyle, buffer)
	s.SetContent((w-23)/2+1, (h-8)/2+1, player.GetGlyph(), nil, boxStyle)
	/* for i := range npcs {
		cx, cy, _ := npcs[i].GetPosition()

	} */
}

func upArrow(g *gocui.Gui, v *gocui.View) error {
	view, _ := g.SetCurrentView("messageView")
	if player.TryMove(level, arena.Up) {
		fmt.Fprintln(view, "Player moves Up")
	}
	return nil
}

func downArrow(g *gocui.Gui, v *gocui.View) error {
	view, _ := g.SetCurrentView("messageView")
	if player.TryMove(level, arena.Down) {
		fmt.Fprintln(view, "Player moves Down.")
	}
	return nil
}

func leftArrow(g *gocui.Gui, v *gocui.View) error {
	view, _ := g.SetCurrentView("messageView")
	if player.TryMove(level, arena.Left) {
		fmt.Fprintln(view, "Player moves Left.")
	}
	return nil
}

func rightArrow(g *gocui.Gui, v *gocui.View) error {
	view, _ := g.SetCurrentView("messageView")
	if player.TryMove(level, arena.Right) {
		fmt.Fprintln(view, "Player moves Right.")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
