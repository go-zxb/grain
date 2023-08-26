// Copyright © 2023 Grain. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repo

import (
	"context"
	"fmt"
	utils "github.com/go-grain/go-utils"
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/grain/internal/repo/data"
	service "github.com/go-grain/grain/internal/service/system"
	model "github.com/go-grain/grain/model/system"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRepo struct {
	data.MongoDB
	rdb redis.IRedis
}

func NewMongoDBRepo(rdb redis.IRedis, mongoUrl, dbName, collectionName string) (service.ISysLogRepo, error) {
	mongoDB := MongoDBRepo{rdb: rdb}
	return &mongoDB, mongoDB.NewMongoDBRepo(mongoUrl, dbName, collectionName)
}

func (r *MongoDBRepo) CreateSysLog(operationLog *model.SysLog) error {
	_, err := r.Collection.InsertOne(context.TODO(), operationLog)
	if err != nil {
		log.Printf("Failed to create operationLog: %v", err)
		return err
	}
	return nil
}

func (r *MongoDBRepo) GetSysLogList(req *model.SysLogReq) ([]*model.SysLog, error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 || req.PageSize >= 100 {
		req.PageSize = 20
	}

	filter := bson.M{}

	if req.Name != "" {
		filter["name"] = req.Name
	}

	if req.Role != "" {
		filter["role"] = req.Role
	}

	if req.Username != "" {
		filter["username"] = req.Username
	}

	if req.QueryTime != "" {
		t := strings.Split(req.QueryTime, ",")
		if len(t) == 2 {
			filter["created_at"] = bson.M{
				"$gte": utils.GetStringToDate(t[0], utils.YMD), // 开始时间
				"$lte": utils.GetStringToDate(t[1], utils.YMD), // 结束时间
			}
		}
	}

	fmt.Printf("%#v\n", req)
	fmt.Printf("%#v\n", filter)
	options := options.Find()
	options.SetSort(bson.M{"_id": -1})
	options.SetSkip(int64((req.Page - 1) * req.PageSize))
	options.SetLimit(int64(req.PageSize))

	documents, err := r.Collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	req.Total = documents
	cursor, err := r.Collection.Find(context.TODO(), filter, options)
	if err != nil {
		log.Printf("Failed to get operationLog list: %v", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var operationLogs []*model.SysLog
	for cursor.Next(context.TODO()) {
		var operationLog model.SysLog
		err := cursor.Decode(&operationLog)
		if err != nil {
			log.Printf("Failed to decode operationLog: %v", err)
			continue
		}
		operationLogs = append(operationLogs, &operationLog)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return operationLogs, nil
}

func (r *MongoDBRepo) DeleteSysLogById(id primitive.ObjectID, uid string) error {
	filter := bson.M{"_id": id, "uid": uid}

	result, err := r.Collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Failed to delete operationLog: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("operationLog not found")
	}

	return nil
}

func (r *MongoDBRepo) DeleteSysLogByIds(ids []primitive.ObjectID, uid string) error {
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
		"uid": uid,
	}

	result, err := r.Collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Printf("Failed to delete operationLogs: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("operationLogs not found")
	}

	return nil
}
