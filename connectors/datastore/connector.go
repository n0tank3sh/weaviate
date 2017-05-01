/*                          _       _
 *__      _____  __ ___   ___  __ _| |_ ___
 *\ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
 * \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
 *  \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
 *
 * Copyright © 2016 Weaviate. All rights reserved.
 * LICENSE: https://github.com/weaviate/weaviate/blob/master/LICENSE
 * AUTHOR: Bob van Luijt (bob@weaviate.com)
 * See www.weaviate.com for details
 * Contact: @weaviate_iot / yourfriends@weaviate.com
 */

package datastore

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/iterator"

	uuid "github.com/satori/go.uuid"

	"cloud.google.com/go/datastore"
)

type Datastore struct{}

type Object struct {
	Uuid         string // uuid, also used in Object's id
	Owner        string // uuid of the owner
	RefType      string // type, as defined
	CreateTimeMs int64  // creation time in ms
	Object       string // the JSON object, id will be collected from current uuid
	Deleted      bool   // if true, it does not exsist anymore
}

// Connect to datastore
func (f Datastore) Connect() bool {
	return true
}

// Add item to DB
func (f Datastore) Add(owner string, refType string, object string) string {

	// Setx your Google Cloud Platform project ID.
	ctx := context.Background()

	projectID := "weaviate-dev-001"

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the kind for the new entity.
	kind := "weaviate"

	// Sets the name/ID for the new entity.
	uuid := fmt.Sprintf("%v", uuid.NewV4())

	// Creates a Key instance.
	taskKey := datastore.NameKey(kind, uuid, nil)

	// Creates a Task instance.
	task := Object{
		Uuid:         uuid,
		Owner:        owner,
		RefType:      refType,
		CreateTimeMs: time.Now().UnixNano() / int64(time.Millisecond),
		Object:       object,
		Deleted:      false,
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, taskKey, &task); err != nil {
		log.Fatalf("Failed to save task: %v", err)
	}

	// return the ID that is used to create.
	return uuid

}

func (f Datastore) Get(Uuid string) string {

	// Setx your Google Cloud Platform project ID.
	ctx := context.Background()
	projectID := "weaviate-dev-001"

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	query := datastore.NewQuery("weaviate").Filter("Uuid =", Uuid).Order("-CreateTimeMs").Limit(1)
	it := client.Run(ctx, query)

	count := 0

	for {
		var object Object
		_, err := it.Next(&object)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching next object: %v", err)
		}

		count++
		fmt.Printf("Uid %q, Time %d, Object %d\n", object.Uuid, object.CreateTimeMs, object.Object)
	}

	println(count)

	return Uuid
}
