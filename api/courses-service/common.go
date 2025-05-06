package main

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertToObjectId(id string) (primitive.ObjectID, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid id")
	}
	return objId, nil
}
