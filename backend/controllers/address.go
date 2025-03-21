package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/agodse21/next-go-full-stack-ecommerce/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *Application) AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		uid, exists := c.Get("uid")

		if !exists {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "user id is empty"})
			c.Abort()
			return
		}

		user_id, ok := uid.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UID format"})
			return
		}

		address, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}

		var addresses models.Address

		addresses.AddressId = primitive.NewObjectID()

		err = c.BindJSON(&addresses)

		if err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		// Validate the struct
		if err := Validate.Struct(addresses); err != nil {
			// Extract validation errors
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Field()+" is invalid")
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": errors})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}

		unwind := bson.D{{
			Key: "$unwind",
			Value: bson.D{
				primitive.E{
					Key:   "path",
					Value: "$address",
				},
			},
		}}

		group := bson.D{
			{
				Key: "$group",
				Value: bson.D{
					primitive.E{
						Key:   "_id",
						Value: "$address_id",
					},
					{
						Key: "count",
						Value: bson.D{
							{
								Key:   "$sum",
								Value: 1,
							},
						},
					},
				},
			},
		}

		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{
			match_filter, unwind, group,
		})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		var addressInfo []bson.M

		err = pointCursor.All(ctx, &addressInfo)
		if err != nil {
			panic(err)
		}

		var size int32

		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}

		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}

			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}

			_, err := UserCollection.UpdateOne(ctx, filter, update)

			if err != nil {

				c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")

			}

			c.IndentedJSON(http.StatusOK, "Address added successfully")
		} else {
			c.IndentedJSON(http.StatusInternalServerError, "Not Allowed")
		}
		defer cancel()
		ctx.Done()
	}
}

func (app *Application) EditAddress() gin.HandlerFunc {

	return func(c *gin.Context) {

		uid, exists := c.Get("uid")

		if !exists {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "user id is empty"})
			c.Abort()
			return
		}

		user_id, ok := uid.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UID format"})
			return
		}

		user_primitive_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		var editaddress models.Address

		if err := c.BindJSON(&editaddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the struct
		if err := Validate.Struct(editaddress); err != nil {
			// Extract validation errors
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Field()+" is invalid")
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": errors})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: user_primitive_id}}

		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House},
			{
				Key:   "address.0.street_name",
				Value: editaddress.Street,
			}, {

				Key:   "address.0.city",
				Value: editaddress.City,
			}, {
				Key:   "address.0.pincode",
				Value: editaddress.Pincode,
			},
		}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		defer cancel()
		ctx.Done()

		c.IndentedJSON(http.StatusOK, "Address updated successfully")

	}
}

func (app *Application) EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, exists := c.Get("uid")

		if !exists {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "user id is empty"})
			c.Abort()
			return
		}

		user_id, ok := uid.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UID format"})
			return
		}

		user_primitive_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}

		var editaddress models.Address

		if err := c.BindJSON(&editaddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the struct
		if err := Validate.Struct(editaddress); err != nil {
			// Extract validation errors
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Field()+" is invalid")
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": errors})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: user_primitive_id}}

		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editaddress.House},
			{
				Key:   "address.1.street_name",
				Value: editaddress.Street,
			}, {

				Key:   "address.1.city",
				Value: editaddress.City,
			}, {
				Key:   "address.1.pincode",
				Value: editaddress.Pincode,
			},
		}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "The work address updated  successfully")
	}
}

func (app *Application) DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		addressQueryId := c.Query("id")

		if addressQueryId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "address id is empty"})
			c.Abort()
			return
		}

		uid, exists := c.Get("uid")

		if !exists {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "user id is empty"})
			c.Abort()
			return
		}

		user_id, ok := uid.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UID format"})
			return
		}

		addressId, err := primitive.ObjectIDFromHex(addressQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server error")
			return
		}

		user_primitive_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server error")
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		filter := bson.M{"_id": user_primitive_id}
		update := bson.M{"$pull": bson.M{"address": bson.M{"_id": addressId}}}

		result, err := UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Check if the address was actually deleted
		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found or already deleted"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})

	}
}
