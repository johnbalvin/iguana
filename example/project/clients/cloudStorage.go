package clients

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
)


var bucketName = ProjectID + ".appspot.com"

var url = "https://storage.googleapis.com/" + bucketName + "/"

//ClientStorage returns the google cloud storage
func ClientStorage() *storage.BucketHandle {
	ctx := context.Background()
	cliente, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("clientestorage -> err:", err)
	}
	return cliente.Bucket(bucketName)
}

//BucketURL returns the first full path of any url file
func BucketURL() string {
	return url
}
