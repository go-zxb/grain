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

package data

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

func (m *MongoDB) NewMongoDBRepo(mongoUrl, dbName, collection string) (err error) {
	clientOptions := options.Client().ApplyURI(mongoUrl)

	m.Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	err = m.Client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}

	m.Database = m.Client.Database(dbName)
	m.Collection = m.Database.Collection(collection)

	return nil
}
