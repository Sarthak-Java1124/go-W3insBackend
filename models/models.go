package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName *string            `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  *string            `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     *string            `json:"email,omitempty" bson:"email,omitempty"`
	Password  *string            `json:"password,omitempty" bson:"password,omitempty"`
}

type NFT struct {
	ID          primitive.ObjectID `form:"_id,omitempty" bson:"_id,omitempty"`
	Headline    *string            `form:"headline" bson:"headline"`
	Description *string            `form:"description" bson:"description"`
	Hashtag     *string            `form:"hashtag" bson:"hashtag"`
	Date        *string            `form:"date" bson:"date"`
	Metadata    *string            `form:"metadata" bson:"metadata"`
	Address     *string            `form:"address" bson:"address"`
	IpfsUrl     *string            `form:"ipfsurl" bson:"ipfsurl"`
}
