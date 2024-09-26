package db

import (
	"context"
	"errors"

	"github.com/kodernubie/gocommon/conf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConn struct {
	db *mongo.Database
}

func (o *MongoConn) Save(obj interface{}) error {

	var err error

	collectionName := GetTableName(obj)

	if collectionName == "" {
		return errors.New("no collection defined")
	}

	id := GetIDValue(obj)

	if id == "" {

		insertHandler, ok := obj.(BeforeInsertHandler)

		if ok {
			insertHandler.BeforeCreate()
		}

		_, err = o.db.Collection(collectionName).InsertOne(context.Background(), obj)
	} else {
		_, err = o.db.Collection(collectionName).ReplaceOne(context.Background(), bson.M{"_id": id}, obj)
	}

	return err
}

func (o *MongoConn) Update(table interface{}, filter interface{}, update interface{}) error {

	collectionName, ok := table.(string)

	if !ok {
		collectionName = GetTableName(table)
	}

	if collectionName == "" {
		return errors.New("no collection defined")
	}

	o.db.Collection(collectionName).UpdateMany(context.Background(), filter, update)
	return nil
}

func (o *MongoConn) Delete(obj interface{}) error {

	var err error

	collectionName := GetTableName(obj)

	if collectionName == "" {
		return errors.New("no collection defined")
	}

	id := GetIDValue(obj)

	if id != "" {

		_, err = o.db.Collection(collectionName).DeleteOne(context.Background(), bson.M{"_id": id})
	}

	return err
}

func (o *MongoConn) DeleteMany(table interface{}, filter interface{}) error {

	var err error

	collectionName, ok := table.(string)

	if !ok {
		collectionName = GetTableName(table)
	}

	if collectionName == "" {
		return errors.New("no collection defined")
	}

	_, err = o.db.Collection(collectionName).DeleteMany(context.Background(), filter)

	return err
}

func (o *MongoConn) FindById(out interface{}, id string) error {

	if out == nil {
		return errors.New("out can't be nil")
	}

	collection := GetTableName(out)

	if collection == "" {
		return errors.New("invalid collection")
	}

	ret := o.db.Collection(collection).FindOne(context.Background(), bson.M{"_id": id})

	if ret != nil {
		ret.Decode(out)
		return nil
	}

	return errors.New("not found")
}

func (o *MongoConn) FindOne(out interface{}, filter interface{}, opt ...FindOption) error {

	if out == nil {
		return errors.New("out can't be nil")
	}

	collection := GetTableName(out)

	if collection == "" {
		return errors.New("invalid collection")
	}

	var targetOpt *options.FindOneOptions

	ret := o.db.Collection(collection).FindOne(context.Background(), filter, targetOpt)

	if ret == nil {
		return errors.New("not found")
	}

	return ret.Decode(out)
}

func (o *MongoConn) Find(out interface{}, filter interface{}, findOptions ...FindOption) error {

	collection := GetTableName(out)

	if collection == "" {
		return errors.New("invalid collection")
	}

	var opts *options.FindOptions

	if len(findOptions) > 0 {

		opts = options.Find()

		if findOptions[0].Skip > 0 {
			opts = opts.SetSkip(int64(findOptions[0].Skip))
		}

		if findOptions[0].Limit > 0 {
			opts = opts.SetLimit(int64(findOptions[0].Limit))
		}

		if len(findOptions[0].Order) > 0 {

			sortList := bson.D{}

			for _, sort := range findOptions[0].Order {

				if sort.Order == "asc" {
					sortList = append(sortList, bson.E{sort.Field, 1})
				} else {
					sortList = append(sortList, bson.E{sort.Field, -1})
				}

			}

			opts = opts.SetSort(sortList)
		}
	}

	cursor, err := o.db.Collection(collection).Find(context.Background(), filter, opts)

	if err != nil {

		return err
	}

	err = cursor.All(context.Background(), out)
	return err
}

func init() {

	RegConnCreator("mongo", func(configName string) (Connection, error) {

		conn := &MongoConn{}

		var err error
		clientOptions := options.Client().ApplyURI(conf.Str("DB_" + configName + "_URL"))
		client, err := mongo.Connect(context.Background(), clientOptions)

		if err == nil {
			conn.db = client.Database(conf.Str("DB_"+configName+"_DBNAME", "default"))
		}

		return conn, err
	})
}
