package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	documentKey struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	namespace struct {
		Db   string `bson:"db"`
		Coll string `bson:"coll"`
	}

	changeID struct {
		Data string `bson:"_data"`
	}

	StreamCompany struct {
		ID            changeID            `bson:"_sid"`
		Ns            namespace           `bson:"ns"`
		OperationType string              `bson:"operationType"`
		ClusterTime   primitive.Timestamp `bson:"clusterTime"`
		DocumentKey   documentKey         `bson:"documentKey"`
		FullDocument  Company             `bson:"fullDocument"`
	}

	StreamContact struct {
		ID            changeID            `bson:"_sid"`
		Ns            namespace           `bson:"ns"`
		OperationType string              `bson:"operationType"`
		ClusterTime   primitive.Timestamp `bson:"clusterTime"`
		DocumentKey   documentKey         `bson:"documentKey"`
		FullDocument  Contact             `bson:"fullDocument"`
	}
)
