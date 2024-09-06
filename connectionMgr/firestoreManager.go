package connectionMgr

import (
	"blog_app/log"
	"blog_app/mycontext"
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/option"
)

type firestoreService struct {
	Client *firestore.Client
}

type FirestoreDB interface {
	ReadOne(ctx mycontext.Context, collectionName string, docID string, data interface{}) error
	ReadAll(ctx mycontext.Context, collectionName string, data interface{}) error
	CreateOne(ctx mycontext.Context, collectionName string, docID string, data interface{}) (*firestore.WriteResult, error)
	Update(ctx mycontext.Context, collectionName string, docID string, data interface{}) (*firestore.WriteResult, error)
	Delete(ctx mycontext.Context, collectionName string, docID string) (*firestore.WriteResult, error)
	Disconnect(ctx mycontext.Context) error
}

func NewFirestoreClient(projectID string, credentialsFile string) (FirestoreDB, error) {
	ctx := mycontext.NewContext()
	client, err := firestore.NewClient(context.Background(), projectID, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.GenericError(ctx, errors.New("can't connect to Firestore: "+err.Error()), nil)
		return &firestoreService{Client: client}, err
	}
	return &firestoreService{Client: client}, nil
}

func (client *firestoreService) ReadOne(ctx mycontext.Context, collectionName string, docID string, data interface{}) error {
	doc, err := client.Client.Collection(collectionName).Doc(docID).Get(ctx.Context)
	if err != nil {
		return err
	}
	return doc.DataTo(data)
}

func (client *firestoreService) ReadAll(ctx mycontext.Context, collectionName string, data interface{}) error {
	iter := client.Client.Collection(collectionName).Documents(ctx.Context)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err != nil {
			return err
		}
		err = doc.DataTo(data)
		if err != nil {
			return err
		}
	}
}

func (client *firestoreService) CreateOne(ctx mycontext.Context, collectionName string, docID string, data interface{}) (*firestore.WriteResult, error) {
	return client.Client.Collection(collectionName).Doc(docID).Set(ctx.Context, data)
}

func (client *firestoreService) Update(ctx mycontext.Context, collectionName string, docID string, data interface{}) (*firestore.WriteResult, error) {
	return client.Client.Collection(collectionName).Doc(docID).Set(ctx.Context, data, firestore.MergeAll)
}

func (client *firestoreService) Delete(ctx mycontext.Context, collectionName string, docID string) (*firestore.WriteResult, error) {
	return client.Client.Collection(collectionName).Doc(docID).Delete(ctx.Context)
}

func (client *firestoreService) Disconnect(ctx mycontext.Context) error {
	return client.Client.Close()
}
