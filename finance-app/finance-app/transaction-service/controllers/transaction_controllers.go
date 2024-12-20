package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"transaction-service/db"
	"transaction-service/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

func getUserIDFromToken(c *gin.Context) (primitive.ObjectID, error) {
	tokenString := c.GetHeader("Authorization")[7:] 
	if tokenString == "" {
		log.Println("Authorization header is missing or invalid")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return primitive.ObjectID{}, errors.New("missing token")
	}

	claims := jwt.MapClaims{}
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return primitive.ObjectID{}, err
	}

	userIDHex := claims["id"].(string)
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return userID, nil
}

func GetTransactionByID(c *gin.Context) {
    userID, err := getUserIDFromToken(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    transactionID := c.Param("id")
    objectID, err := primitive.ObjectIDFromHex(transactionID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
        return
    }

    var transaction models.Transaction
    err = db.TransactionsCollection.FindOne(context.Background(), bson.M{
        "_id":     objectID,
        "user_id": userID,
    }).Decode(&transaction)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

func GetStatistics(c *gin.Context) {
    userID, err := getUserIDFromToken(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    now := time.Now()
    startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
    endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond) 

    incomeFilter := bson.M{
        "user_id":  userID,
        "category": "income",
        "date": bson.M{
            "$gte": startOfMonth,
            "$lte": endOfMonth,
        },
    }

    expenseFilter := bson.M{
        "user_id":  userID,
        "category": "expense",
        "date": bson.M{
            "$gte": startOfMonth,
            "$lte": endOfMonth,
        },
    }

    incomeCursor, err := db.TransactionsCollection.Aggregate(context.Background(), []bson.M{
        {"$match": incomeFilter},
        {"$group": bson.M{"_id": nil, "total_income": bson.M{"$sum": "$amount"}}},
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calculating income"})
        return
    }

    expenseCursor, err := db.TransactionsCollection.Aggregate(context.Background(), []bson.M{
        {"$match": expenseFilter},
        {"$group": bson.M{"_id": nil, "total_expense": bson.M{"$sum": "$amount"}}},
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calculating expense"})
        return
    }

    var incomeResult, expenseResult []bson.M
    if err := incomeCursor.All(context.Background(), &incomeResult); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing income data"})
        return
    }
    if err := expenseCursor.All(context.Background(), &expenseResult); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing expense data"})
        return
    }

    var expenseTotal float64 = 0
    var incomeTotal float64 = 0
    if len(incomeResult) > 0 {
        if val, ok := incomeResult[0]["total_income"].(float64); ok {
            incomeTotal = val
        }
    }
    if len(expenseResult) > 0 {
        if val, ok := expenseResult[0]["total_expense"].(float64); ok {
            expenseTotal = val
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "income":  incomeTotal,
        "expense": expenseTotal,
    })
}

func GetTransactions(c *gin.Context) {
    userID, err := getUserIDFromToken(c)
    if err != nil {
        log.Println("Error getting userID from token:", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }
    
    log.Println("Fetching transactions for userID:", userID)

    now := time.Now()
    startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
    endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond) 

    filter := bson.M{
        "user_id": userID,
        "date": bson.M{
            "$gte": startOfMonth,
            "$lte": endOfMonth,
        },
    }

    cursor, err := db.TransactionsCollection.Find(context.Background(), filter)
    if err != nil {
        log.Println("Error finding transactions:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
        return
    }
    defer cursor.Close(context.Background())

    var transactions []models.Transaction
    if err := cursor.All(context.Background(), &transactions); err != nil {
        log.Println("Error reading transactions:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading transactions"})
        return
    }

    log.Printf("Found %d transactions", len(transactions))

    c.JSON(http.StatusOK,  gin.H{"transactions": transactions})
}


func AddTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

    transaction.ID = primitive.NewObjectID()
	transaction.UserID = primitive.ObjectID(userID)
	transaction.Date = time.Now()

	_, err = db.TransactionsCollection.InsertOne(context.Background(), transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction added successfully"})
}

func DeleteTransaction(c *gin.Context) {
	transactionID := c.Param("id")
	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	_, err = db.TransactionsCollection.DeleteOne(context.Background(), bson.M{
		"_id":     objectID,
		"user_id": userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

func EditTransaction(c *gin.Context) {
	transactionID := c.Param("id")
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	log.Printf("Updating transaction: %v", objectID)

	update := bson.M{
		"$set": bson.M{
			"amount":      transaction.Amount,
			"description": transaction.Description,
			"category":    transaction.Category,
			"date":        time.Now(),
		},
	}

	_, err = db.TransactionsCollection.UpdateOne(context.Background(), bson.M{
		"_id":     objectID,
		"user_id": userID,
	}, update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
}
