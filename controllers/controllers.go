package controllers

// The `import "github.com/gin-gonic/gin"` statement in the Go code is importing the `gin` package from
// the `github.com/gin-gonic` repository. This package is commonly used for building web applications
// and APIs in Go using the Gin framework.
import (
	"context"
	"e-CommerceBackend/database"
	"e-CommerceBackend/models"
	"e-CommerceBackend/utils"

	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func SignUp(c *gin.Context) {
	var body models.User

	// Parse JSON
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if user already exists
	userCollection := database.MongoClient.Database("topmateapp").Collection("users")
	isPresent, err := userCollection.CountDocuments(context.TODO(), bson.M{"email": body.Email})
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	if isPresent > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
		return
	}

	// Hash password
	hashedPassword := utils.HashPassword(*body.Password)
	body.Password = &hashedPassword

	// Insert user into database
	_, err = userCollection.InsertOne(context.TODO(), body)
	if err != nil {
		fmt.Println("Error inserting user into database:", err)
		c.JSON(500, gin.H{"message": "Failed to insert into database"})
		return
	}

	// Success response
	c.JSON(200, gin.H{"message": "Signup successful"})
}

func Login(c *gin.Context) {
	var RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := c.BindJSON(&RequestBody)

	if err != nil {
		c.JSON(500, gin.H{"message": "Unable to bind the json format"})
		return
	}
	isUser := database.MongoClient.Database("topmateapp").Collection("users")
	count, err := isUser.CountDocuments(context.TODO(), bson.M{"email": RequestBody.Email})
	if count == 0 {
		c.JSON(500, gin.H{"message": "No User found"})
		return
	}
	var presentUser models.User
	isUser.FindOne(context.TODO(), bson.M{"email": RequestBody.Email}).Decode(&presentUser)
	comparePassword := utils.VerifyPassword(*presentUser.Password, RequestBody.Password)

	if !comparePassword {
		c.JSON(500, gin.H{"message": "Wrong Password"})
		return
	}
	token, err := utils.GenerateJwt(*presentUser.FirstName)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error in Generating the jwt token"})
		return
	}

	c.SetCookie("token", token, 3600*24*10, "/", "localhost", false, true)

	c.JSON(200, gin.H{"message": "Login Succesfully"})
}

func FormSubmit(c *gin.Context) {
	var body models.NFT

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	file, err := c.FormFile("image")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Image Cannot be Uploaded"})

	}

	ipfsUrl, err := utils.UploadToPinata(file)

	metadata := map[string]interface{}{
		"headline":    body.Headline,
		"description": body.Description,
		"image":       ipfsUrl,
		"date":        body.Date,
	}
	metadataUrl, err := utils.UploadJSONToPinata(metadata)
	if err != nil {
		fmt.Println("The error in getting the metadata is : ", metadataUrl)

	}

	fmt.Println("The metadata url is : ", metadataUrl)
	if err != nil {
		fmt.Println("The error in getting the ipfS url is ", err)
	}
	fmt.Println("The ipfs url is : ", ipfsUrl)

	savePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Image cannot be saved"})
	}

	nftcollection := database.MongoClient.Database("topmateapp").Collection("nfts")
	body.IpfsUrl = &ipfsUrl
	response, err := nftcollection.InsertOne(context.TODO(), body)
	fmt.Println("The response is : ", response)
	if err != nil {
		c.JSON(500, "Error in Saving Data in DB")
	}
	c.JSON(200, gin.H{
		"message":     "Form Submitted Successfully",
		"success":     true,
		"ipfsURL":     metadataUrl,
		"imageURL":    ipfsUrl,
		"title":       body.Headline,
		"headline":    body.Hashtag,
		"description": body.Description,
		"date":        body.Date,
	})
}

func GetFormData(c *gin.Context) {
	address := c.Param("address")
	filter := bson.M{"address": address}

	cursor, err := database.MongoClient.Database("topmateapp").Collection("nfts").Find(context.TODO(), filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error finding documents"})
		return
	}
	defer cursor.Close(context.TODO())

	var results []models.NFT
	for cursor.Next(context.TODO()) {
		var nft models.NFT
		if err := cursor.Decode(&nft); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding document"})
			return
		}
		fmt.Println("The nft is : ", nft)
		results = append(results, nft)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cursor error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"cardData": results,
	})
}
