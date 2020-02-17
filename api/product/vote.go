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

func Trust(v PerVote) (int, int) {
	up := v.UpHigh*30 + v.UpLow*20
	down := v.DownHigh*30 + v.DownLow*20
	if up > 100 {
		up = 100
	}
	if down > 100 {
		down = 100
	}

	return up, down
}

func confirmed(users []UserVote, up bool, conn *mongo.Database) {
	for _, u := range users {
		if u.Up == up {
			user.PointsForUpdate(u.User, conn)
		} else {
			user.PointsForDeny(u.User, conn)
		}
	}
}

func voteComplete(v PerVote) bool {
	tot := v.UpHigh*2 + v.UpLow
	totd := v.DownHigh*2 + v.DownLow
	return tot > 5 || totd > 5
}

func upWon(v PerVote) bool {
	tot := v.UpHigh*2 + v.UpLow
	totd := v.DownHigh*2 + v.DownLow
	return tot > totd
}

func VoteOnProduct(v Vote, username string, conn *mongo.Database) {
	p := GetProductInfo(v.ID, username, conn)
	level := user.GetLevel(username, conn)
	if canVote(p.ProductName.Users, username) && !voteComplete(p.ProductName.Votes) {
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
		if voteComplete(p.ProductName.Votes) {
			confirmed(p.ProductName.Users, upWon(p.ProductName.Votes), conn)
		}
	}
	if canVote(p.Ingredients.Users, username) && !voteComplete(p.Ingredients.Votes) {

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
		if voteComplete(p.Ingredients.Votes) {
			confirmed(p.Ingredients.Users, upWon(p.Ingredients.Votes), conn)
		}
	}
	if canVote(p.Nutrition.Users, username) && !voteComplete(p.Nutrition.Votes) {

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
			p.Nutrition.Users = append(p.Nutrition.Users, UserVote{username, false})
		}

		if voteComplete(p.Nutrition.Votes) {
			confirmed(p.Nutrition.Users, upWon(p.Nutrition.Votes), conn)
		}
	}
	collection := conn.Collection("products")
	filter := bson.M{"_id": p.ID}
	collection.FindOneAndReplace(context.TODO(), filter, p)

}
