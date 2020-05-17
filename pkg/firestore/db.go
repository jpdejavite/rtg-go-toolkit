package firestore

import (
	"context"
	"encoding/base64"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// IDBFirestore firestore db interface
type IDBFirestore interface {
	GetDocumentData(collection string, document string) (map[string]interface{}, error)
}

// NewDBFirestore returns a new db interface
func NewDBFirestore(AppDB firebase.App) IDBFirestore {
	return DBFirestore{AppDB}
}

// DBFirestore implements IDBFirestore interface
type DBFirestore struct {
	AppDB firebase.App
}

// ConnectToDatabase connect to firestore database using base service account
func ConnectToDatabase(base64ServiceAccount string) (IDBFirestore, error) {
	serviceAccount := make([]byte, base64.StdEncoding.DecodedLen(len(base64ServiceAccount)))
	_, err := base64.StdEncoding.Decode(serviceAccount, []byte(base64ServiceAccount))
	if err != nil {
		return nil, err
	}

	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsJSON(serviceAccount)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return DBFirestore{AppDB: *app}, nil
}

// GetDocumentData get firestore document data
func (dbFirestore DBFirestore) GetDocumentData(collection string, document string) (map[string]interface{}, error) {
	client, err := dbFirestore.AppDB.Firestore(context.Background())
	if err != nil {
		return nil, err
	}
	docSnap, err := client.Collection(collection).Doc(document).Get(context.Background())
	if err != nil {
		return nil, err
	}
	if docSnap != nil && (*docSnap).Data() != nil {
		return (*docSnap).Data(), nil
	}
	return nil, nil
}
