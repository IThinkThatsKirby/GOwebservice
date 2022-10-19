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
var db = []DataBase{}

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
	// spend points return { "payer": <string>, "points": <integer> }
	app.Post("/spend")
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
type DataBase struct {
	Payer  string
	Points int
	//sortable timestamp when transaction happened.
	UnixTime int64
}

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

// Handle Database CRUD
func dbCRUD(payer string, points int, tUnix int64) {
	crudDATA := DataBase{
		payer,
		points,
		tUnix,
	}
	db = append(db, crudDATA)
	//sort db by points old to new
	// WARNING if UNIX TIME IS THE SAME IT WILL MESS WITH SORTING !!
	// todo: fix UNIX TIME so its never identical or improve sorting
	sortDB()
}

// BLAZINGLY FAST JSON MARSHAL in function
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
func updatePoints() (payerPoints map[string]int) {
	// current payers points {payer: points}
	payerPoints = map[string]int{"kirby": 0}
	for i, points := range db {
		payerPoints[points.Payer] += db[i].Points
	}
	return payerPoints
}

// sorts the database of transactions by oldest to newest :D
func sortDB() []DataBase {
	sort.Slice(db, func(i, j int) bool {
		return db[i].UnixTime < db[j].UnixTime
	})
	return db
}
