package product

import (
	"context"
	"log"
	"time"

	"../user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Vote struct {
	ID          string
	Name        int
	Ingredients int
	Nutrition   int
	Serving     int
}

type VoteCheck struct {
	part     string
	barcode  string
	version  int64
	username string
	conn     *mongo.Database
}

func canVote(v VoteCheck) bool {
	v1 := v
	v1.part = v1.part + "UP"
	v2 := v
	v2.part = v2.part + "DOWN"
	return checkCanVote(v1) && checkCanVote(v2)

}

func checkCanVote(v VoteCheck) bool {
	collection := v.conn.Collection("user")
	filter := bson.D{{"_id", v.username},
		{"pointsHistory.item", v.barcode},
		{"pointsHistory.type", v.part},
		{"pointsHistory.version", v.version}}
	doc, err := collection.Find(context.TODO(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}
	f := doc.Next(context.TODO())
	if !f {
		log.Print("Can vote")
		return true
	}
	log.Print("Cant vote")
	return false
}

func VoteOnProduct(v Vote, username string, conn *mongo.Database) {
	p := GetProductInfo(v.ID, conn)
	sec := time.Now().Unix()
	vc := VoteCheck{"NAME", v.ID, p.ProductName.Stamp, username, conn}
	if canVote(vc) {
		if v.Name > 0 {
			p.ProductName.Up++
			point := user.Point{p.ID, p.ProductName.Stamp, "NAMEUP", 1, false, sec}
			user.AddPoint(point, username, conn)
		} else if v.Name < 0 {
			p.ProductName.Down--
			point := user.Point{p.ID, p.ProductName.Stamp, "NAMEDOWN", 1, false, sec}
			user.AddPoint(point, username, conn)
		}
	}
	vc = VoteCheck{"INGREDIENTS", v.ID, p.Ingredients.Stamp, username, conn}
	if canVote(vc) {

		if v.Ingredients > 0 && canVote(vc) {
			p.Ingredients.Up++
			point := user.Point{p.ID, p.Ingredients.Stamp, "INGREDIENTSUP", 1, false, sec}
			user.AddPoint(point, username, conn)
		} else if v.Ingredients < 0 {
			p.Ingredients.Down--
			point := user.Point{p.ID, p.Ingredients.Stamp, "INGREDIENTSDOWN", 1, false, sec}
			user.AddPoint(point, username, conn)
		}
	}
	vc = VoteCheck{"NUTRITION", v.ID, p.Nutrition.Stamp, username, conn}
	if canVote(vc) {

		if v.Nutrition > 0 && canVote(vc) {
			p.Nutrition.Up++
			point := user.Point{p.ID, p.Nutrition.Stamp, "NUTRITIONUP", 1, false, sec}
			user.AddPoint(point, username, conn)
		} else if v.Nutrition < 0 {
			p.Nutrition.Down--
			point := user.Point{p.ID, p.ProductName.Stamp, "NUTRITIONDOWN", 1, false, sec}
			user.AddPoint(point, username, conn)
		}
	}
	vc = VoteCheck{"SERVING", v.ID, p.Serving.Stamp, username, conn}
	if canVote(vc) {

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

}
