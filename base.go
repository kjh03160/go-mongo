package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type BaseModel struct {
	CreatedDate time.Time `json:"created_date" bson:"created_date"`
	UpdatedDate time.Time `json:"updated_date" bson:"updated_date"`
}

func (model *BaseModel) Create() {
	model.CreatedDate = time.Now()
	model.UpdatedDate = time.Now()
}

func (model *BaseModel) Save() {
	model.UpdatedDate = time.Now()
}

func StructToBson(s interface{}) (*bson.D, error) {
	var update bson.D
	b, err := bson.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("struct bson marshal err:%w, struct:%+v", err, s)
	}
	err = bson.Unmarshal(b, &update)
	if err != nil {
		return nil, fmt.Errorf("struct bson unmarshal err:%w, struct:%+v", err, s)
	}
	return &update, nil
}
