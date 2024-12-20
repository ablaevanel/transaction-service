package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`        
	Amount      float64            `json:"amount" bson:"amount"`           
	Description string             `json:"description" bson:"description"` 
	Category    string             `json:"category" bson:"category"`    
	Date        time.Time          `json:"date" bson:"date"`
}
