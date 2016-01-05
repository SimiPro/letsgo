package common
import (
	"golang.org/x/net/context"
	"google.golang.org/cloud/datastore"
	"log"
)

func GetDatastore(ctx context.Context) datastore.Client {
	client, err := datastore.NewClient(ctx, PROEJCT_ID)
	if err != nil {
		log.Fatalf("Failed to create datastore client: %v", client)
	}
	return client
}