package creature

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	ar "rogue/arena"
	at "rogue/creature/attribute"

	_ "github.com/go-sql-driver/mysql"
)

type Creature struct {
	name, race string
	attribs    map[string]*at.Attribute
	x, y, z    int
	glyph      rune
}

func MakeCreature(n string, x int, y int) *Creature {

	var newCreature = Creature{}

	newCreature.attribs = make(map[string]*at.Attribute)
	newCreature.z = 1
	newCreature.x = x
	newCreature.y = y
	newCreature.race = n

	db, err := sql.Open("mysql", "rodney:Akhen@t0n@tcp(127.0.0.1:3306)/roguelike")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	res, err2 := db.Query("SELECT tbl_character_attribute_templates.name, tbl_character_attribute_templates.attr_template FROM (tbl_character_attribute_templates INNER JOIN tbl_creatures ON tbl_character_attribute_templates.creature = tbl_creatures.UID) WHERE tbl_creatures.name = ?", n)

	if err2 != nil {
		log.Fatal(err2)
	}

	defer res.Close()

	for res.Next() {
		var attrName string
		var template int
		err3 := res.Scan(&attrName, &template)

		if err3 != nil {
			log.Fatal(err3)
		}

		newCreature.attribs[attrName] = at.MakeAttribute(template)
	}

	res2, err4 := db.Query("SELECT tbl_names.name FROM (tbl_names INNER JOIN tbl_creatures ON tbl_names.race = tbl_creatures.UID) WHERE tbl_creatures.name = ?", n)

	if err4 != nil {
		log.Fatal(err4)
	}

	defer res2.Close()

	names := make([]string, 0)
	for res2.Next() {
		var creatureName string

		err5 := res2.Scan(&creatureName)

		if err5 != nil {
			log.Fatal(err5)
		}

		names = append(names, creatureName)
	}

	if len(names) == 0 {
		switch n {
		case "human":
			newCreature.name = "Player"
		default:
			newCreature.name = "Creature"
		}
	} else {
		position := rand.Intn(len(names) - 1)
		newCreature.name = names[position]
	}

	res3 := db.QueryRow("SELECT tbl_displayables.glyph FROM (tbl_displayables INNER JOIN tbl_creatures ON tbl_displayables.UID = tbl_creatures.display) WHERE tbl_creatures.name = ?", n)

	glyph := make([]byte, 2)

	err6 := res3.Scan(&glyph)

	if err6 != nil {
		log.Fatal(err6)
	}

	newCreature.glyph = rune(glyph[0])

	return &newCreature
}

func (c *Creature) GetName() string {
	return c.name
}

func (c *Creature) GetPosition() (int, int, int) {
	return c.x, c.y, c.z
}

func (c *Creature) GetGlyph() rune {
	return c.glyph
}

func (c *Creature) DebugPrint() {
	fmt.Printf("%s {\n\trace:\t%s\n", c.name, c.race)
	for key, value := range c.attribs {
		fmt.Printf("%s: \t%s\n", key, value.DebugPrint())
	}
	fmt.Printf("location :\n\tx: %d\ty: %d\tz: %d\n\n", c.x, c.y, c.z)
}

func (c *Creature) TryMove(a *ar.ArenaType, d ar.Direction) bool {
	switch d {
	case ar.Up:
		if a.IsLocationPassable(c.x, c.y-1) {
			c.y--
			return true
		}
	case ar.Down:
		if a.IsLocationPassable(c.x, c.y+1) {
			c.y++
			return true
		}
	case ar.Left:
		if a.IsLocationPassable(c.x-1, c.y) {
			c.x--
			return true
		}
	case ar.Right:
		if a.IsLocationPassable(c.x+1, c.y) {
			c.x++
			return true
		}
	}

	return false
}

func (c *Creature) NamePosView() string {

	buffer := fmt.Sprintf("Name: %q\nPos:  x-%d y-%d\n", c.name, c.x, c.y)

	return buffer
}

func (c *Creature) StatView() string {

	var buffer string

	buffer = fmt.Sprintf("Strength: %d\n", c.attribs["strength"].Get())
	buffer += fmt.Sprintf("Intelligence: %d\n", c.attribs["intelligence"].Get())
	buffer += fmt.Sprintf("Dexterity: %d\n", c.attribs["dexterity"].Get())
	buffer += fmt.Sprintf("Magic:  %d/%d\n", c.attribs["magic"].Get(), c.attribs["magic"].GetMax())
	buffer += fmt.Sprintf("Health: %d/%d\n", c.attribs["health"].Get(), c.attribs["health"].GetMax())
	buffer += fmt.Sprintf("Hunger: %d/%d\n", c.attribs["hunger"].Get(), c.attribs["hunger"].GetMax())

	return buffer
}
