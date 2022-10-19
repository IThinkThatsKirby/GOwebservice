package main

import (
	"sort"
	"time"

	// "sort"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// instance of mem database of in-coming calls
var db = DataBase{}

func main() {
	// instance of Fiber w/BLAZINGLY fast JSON marshals
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	// please don't break it
	app.Use(recover.New())
	// --> ROUTES <--
	// send me calls
	app.Post("/", process)
	// spend points return []{ "payer": <string>, "points": -<points int> }
	app.Post("/spend", spendit)
	//points balance for all payers
	app.Get("/points", totalPoints)
	// Open a port for HTTP reqeusts
	app.Listen(":5000")
}

// struct for the JSON calls
type Transaction struct {
	Payer     string `json:"payer"`
	Points    int    `json:"points"`
	Timestamp string `json:"timestamp"`
}

// in memory DataBase structure
type Data struct {
	Payer  string
	Points int
	//sortable timestamp when transaction happened.
	UnixTime int64
}
type DataBase []*Data

// Handler JSON unmarsheling
func process(c *fiber.Ctx) error {
	data := new(Transaction)
	c.BodyParser(data)
	t, _ := time.Parse(time.RFC3339, data.Timestamp)
	//Convert time to unix for sorting transactions by oldest
	tUnix := t.UnixNano()
	dbCRUD(data.Payer, data.Points, tUnix)
	return c.SendString("We received: " + strIt(data))
}

// incoming req shape
type Spendings struct {
	Points int `json:"points"`
}

// res shape{"payer",points Spent:}
type ResShape struct {
	Payer  string
	Points int
}

// type resLog []*ResShape

// Spend it
func spendit(c *fiber.Ctx) error {
	type spent map[string]int
	res := spent{}
	spendReq := new(Spendings)
	c.BodyParser(spendReq)
	for i, data := range db {
		if data.Points >= spendReq.Points {
			db[i].Points -= spendReq.Points
			// logging for response
			res[data.Payer] = (-1 * spendReq.Points)
			// building res [{resData}]
			spendReq.Points = 0
			println(strIt(res))
			return c.JSON(res)
		} else if data.Points < spendReq.Points {
			// logging for response
			res[data.Payer] = (-1 * db[i].Points)
			// building res [{resData}]
			// subtract from points in order of oldest to newest
			spendReq.Points -= db[i].Points
			db[i].Points = 0
		}
		if spendReq.Points == 0 {
			return c.JSON(res)
		}
		continue
	}
	return c.JSON(res)
}

// Handle Database CRUD
func dbCRUD(payer string, points int, tUnix int64) {
	crudDATA := Data{
		payer,
		points,
		tUnix,
	}
	db = append(db, &crudDATA)
	//sort db by points old to new
	// WARNING if UNIX TIME IS THE SAME IT WILL MESS WITH SORTING !!
	// todo: fix UNIX TIME so its never identical or improve sorting
	sort.Stable(ByUnix{db})
}

// BLAZINGLY FAST JSON MARSHAL
func strIt(data any) (str string) {
	strByte, _ := json.Marshal(data)
	str = string(strByte)
	return str
}

// find current balance of payers points
func totalPoints(c *fiber.Ctx) error {
	return c.JSON(updatePoints())
}

// get payer points totals
func updatePoints() map[string]int {
	// current payers points {payer: points}
	sort.Stable(ByUnix{db})
	type Totals map[string]int
	res := Totals{}
	for i, value := range db {
		res[db[i].Payer] += value.Points
	}
	return res
}

// sorts the database of transactions by oldest to newest :D
type ByUnix struct{ DataBase }

// defining sort interface doings.
func (db ByUnix) Len() int           { return len(db.DataBase) }
func (db ByUnix) Swap(i, j int)      { db.DataBase[i], db.DataBase[j] = db.DataBase[j], db.DataBase[i] }
func (db ByUnix) Less(i, j int) bool { return db.DataBase[i].UnixTime < db.DataBase[j].UnixTime }
