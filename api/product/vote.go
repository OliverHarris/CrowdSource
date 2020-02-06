package product

import (
	"context"

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

func canVote(v []UserVote, user string) bool {
	for _, a := range v {
		if a.User == user {
			return false
		}
	}
	return true
}

func VoteOnProduct(v Vote, username string, conn *mongo.Database) {
	p := GetProductInfo(v.ID, username, conn)
	level := user.GetLevel(username, conn)
	if canVote(p.ProductName.Users, username) {
		if v.Name > 0 {
			switch level {
			case 0:
				p.ProductName.Votes.UpLow++
			default:
				p.ProductName.Votes.UpHigh++
			}
			p.ProductName.Users = append(p.ProductName.Users, UserVote{username, true})
		} else if v.Name < 0 {
			switch level {
			case 0:
				p.ProductName.Votes.DownLow++
			default:
				p.ProductName.Votes.DownHigh++
			}
			p.ProductName.Users = append(p.ProductName.Users, UserVote{username, false})
		}
	}
	if canVote(p.Ingredients.Users, username) {

		if v.Ingredients > 0 {
			switch level {
			case 0:
				p.Ingredients.Votes.UpLow++
			default:
				p.Ingredients.Votes.UpHigh++
			}
			p.Ingredients.Users = append(p.Ingredients.Users, UserVote{username, true})
		} else if v.Ingredients < 0 {
			switch level {
			case 0:
				p.Ingredients.Votes.DownLow++
			default:
				p.Ingredients.Votes.DownHigh++
			}
			p.Ingredients.Users = append(p.Ingredients.Users, UserVote{username, false})
		}
	}
	if canVote(p.Nutrition.Users, username) {

		if v.Nutrition > 0 {
			switch level {
			case 0:
				p.Nutrition.Votes.UpLow++
			default:
				p.Nutrition.Votes.UpHigh++
			}
			p.Nutrition.Users = append(p.Nutrition.Users, UserVote{username, true})
		} else if v.Nutrition < 0 {
			switch level {
			case 0:
				p.Nutrition.Votes.DownLow++
			default:
				p.Nutrition.Votes.DownHigh++
			}
		}
		p.Nutrition.Users = append(p.Nutrition.Users, UserVote{username, false})
	}
	collection := conn.Collection("products")
	filter := bson.M{"_id": p.ID}
	collection.FindOneAndReplace(context.TODO(), filter, p)

}
