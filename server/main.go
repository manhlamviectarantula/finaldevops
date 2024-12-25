package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Price float64            `json:"price" bson:"price"`
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://mongo-db:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	productCollection := client.Database("devops-1").Collection("products")
	defer client.Disconnect(context.Background())
	fmt.Println("Connected to MongoDB!")

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/api/v1/get-products", func(c *gin.Context) {
		cursor, err := productCollection.Find(c.Request.Context(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(c.Request.Context())

		var products []Product
		if err := cursor.All(c.Request.Context(), &products); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	r.POST("/api/v1/add-products", func(c *gin.Context) {
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product.ID = primitive.NewObjectID()
		_, err := productCollection.InsertOne(c.Request.Context(), product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, product)
	})

	r.PUT("/api/v1/update-product/:id", func(c *gin.Context) {
		id := c.Param("id")
		// fmt.Println("Update route hit with ID:", id)
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var updateData Product
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		update := bson.M{
			"$set": bson.M{
				"name":  updateData.Name,
				"price": updateData.Price,
			},
		}
		result, err := productCollection.UpdateOne(c.Request.Context(), bson.M{"_id": objectID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
	})

	r.DELETE("/api/v1/delete-product/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		result, err := productCollection.DeleteOne(c.Request.Context(), bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found111"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	})

	r.Run(":8082")
}
