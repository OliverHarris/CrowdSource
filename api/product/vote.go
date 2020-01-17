package product

import (
	"time"

	"../user"
	"go.mongodb.org/mongo-driver/mongo"
)

type Vote struct {
	ID          string
	Name        int
	Ingredients int
	Nutrition   int
	Serving     int
}

func VoteOnProduct(v Vote, username string, conn *mongo.Database) {
	p := GetProductInfo(v.ID, conn)
	sec := time.Now().Unix()
	if v.Name > 0 {
		p.ProductName.Up++
		point := user.Point{p.ID, p.ProductName.Stamp, "NAMEUP", 1, false, sec}
		user.AddPoint(point, username, conn)
	} else if v.Name < 0 {
		p.ProductName.Down--
		point := user.Point{p.ID, p.ProductName.Stamp, "NAMEDOWN", 1, false, sec}
		user.AddPoint(point, username, conn)
	}

	if v.Ingredients > 0 {
		p.Ingredients.Up++
		point := user.Point{p.ID, p.Ingredients.Stamp, "INGREDIENTSUP", 1, false, sec}
		user.AddPoint(point, username, conn)
	} else if v.Ingredients < 0 {
		p.Ingredients.Down--
		point := user.Point{p.ID, p.Ingredients.Stamp, "INGREDIENTSDOWN", 1, false, sec}
		user.AddPoint(point, username, conn)
	}

	if v.Nutrition > 0 {
		p.Nutrition.Up++
		point := user.Point{p.ID, p.Nutrition.Stamp, "NUTRITIONUP", 1, false, sec}
		user.AddPoint(point, username, conn)
	} else if v.Nutrition < 0 {
		p.Nutrition.Down--
		point := user.Point{p.ID, p.ProductName.Stamp, "NUTRITIONDOWN", 1, false, sec}
		user.AddPoint(point, username, conn)
	}

	if v.Serving > 0 {
		p.Serving.Up++
		point := user.Point{p.ID, p.Serving.Stamp, "SERVINGUP", 1, false, sec}
		user.AddPoint(point, username, conn)
	} else if v.Serving < 0 {
		p.Serving.Down--
		point := user.Point{p.ID, p.Serving.Stamp, "SERVINGDOWN", 1, false, sec}
		user.AddPoint(point, username, conn)
	}

}
