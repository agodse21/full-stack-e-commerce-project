package main

import (
	"log"
	"net/http"
	"os"

	"github.com/agodse21/next-go-full-stack-ecommerce/backend/controllers"
	"github.com/agodse21/next-go-full-stack-ecommerce/backend/database"
	"github.com/agodse21/next-go-full-stack-ecommerce/backend/middlewares"
	"github.com/agodse21/next-go-full-stack-ecommerce/backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

//	bson.D = it is a array of key value pairs,where type E struct {
//	    Key   string
//	    Value interface{}
//	}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.Default()

	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middlewares.Authentication())

	router.POST("/user/addaddress", app.AddAddress())
	router.PATCH("/user/edithomeaddress", app.EditAddress())
	router.PATCH("/user/editworkaddress", app.EditWorkAddress())
	router.DELETE("/user/deleteaddress", app.DeleteAddress())

	router.GET("/getproduct/:id", controllers.SearchProductById())
	router.GET("/addtocart", app.AddToCart())
	router.DELETE("/removeitem", app.RemoveItem())
	router.GET("/getcartitems", app.GetItemsFromCart())
	router.GET("/addtowishlist", app.AddToWishlist())
	router.DELETE("/removefromwishlist", app.RemoveItemFromWishlist())
	router.GET("/getwishlistitems", app.GetItemsFromWishlist())
	router.GET("/checkout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow frontend
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // Required if using authentication
	}).Handler(router)

	// log.Fatal(router.Run(":"+port, handler))
	log.Fatal(http.ListenAndServe(":"+port, handler))

}
