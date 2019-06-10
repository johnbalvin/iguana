package clients

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

//ClientFirestore returns client to work with firestore
func ClientFirestore() *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, ProjectID)
	if err != nil {
		log.Fatal("clientefilestore -> client -> err:", err)
	}
	return client
}
