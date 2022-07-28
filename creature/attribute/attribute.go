package attribute

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Attribute struct {
	min, max, base, current int
	chanceChange            float32
}

type template struct {
	min, max, numberOfDice, sizeOfDice, bonus int
	chanceOfChange                            float32
}

func MakeAttribute(attrNum int) *Attribute {
	var attrib Attribute
	var temp template
	var rolls []int

	rand.Seed(time.Now().UTC().UnixNano())
	db, err := sql.Open("mysql", "rodney:Akhen@t0n@tcp(127.0.0.1:3306)/roguelike")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	res := db.QueryRow("SELECT min, max, number_of_dice, size_of_die, bonus, chance_of_change FROM tbl_attribute_templates where UID = ?", attrNum)

	err2 := res.Scan(&temp.min, &temp.max, &temp.numberOfDice, &temp.sizeOfDice, &temp.bonus, &temp.chanceOfChange)
	if err2 != nil {
		log.Fatal(err2)
	}

	if temp.numberOfDice > 0 {
		rolls = make([]int, temp.numberOfDice)
		for i := range rolls {
			rolls[i] = rand.Intn(temp.sizeOfDice) + 1
		}
	}

	attrib.min = temp.min
	attrib.max = temp.max
	attrib.chanceChange = temp.chanceOfChange
	if temp.numberOfDice <= 0 {
		attrib.current = temp.max
		attrib.base = temp.max
	} else {
		for i := 0; i < temp.numberOfDice; i++ {
			if temp.numberOfDice > 0 {
				attrib.base += rolls[i]
			} else {
				attrib.base = temp.max
			}
		}
		attrib.base += temp.bonus
		attrib.current = attrib.base
	}

	return &attrib
}

func (a *Attribute) DebugPrint() string {
	return fmt.Sprintf("min: %d  max:%d  base:%d  current:%d  chance of change:%f%%", a.min, a.max, a.base, a.current, (a.chanceChange * 100))
}

func (a *Attribute) Get() int {
	return a.current
}

func (a *Attribute) GetMax() int {
	return a.max
}
