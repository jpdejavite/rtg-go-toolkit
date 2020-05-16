package firestore

import (
	"context"
	"encoding/base64"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// IDBFirestore firestore db interface
type IDBFirestore interface {
	ConnectToDatabase(base64ServiceAccount string) error
	GetDocumentData(collection string, document string) (map[string]interface{}, error)
}

// NewDBFirestore returns a new db interface
func NewDBFirestore() IDBFirestore {
	return DBFirestore{}
}

// DBFirestore implements IDBFirestore interface
type DBFirestore struct {
	appDB *firebase.App
}

// ConnectToDatabase connect to firestore database using base service account
func (dbFirestore DBFirestore) ConnectToDatabase(base64ServiceAccount string) error {
	serviceAccount := make([]byte, base64.StdEncoding.DecodedLen(len(base64ServiceAccount)))
	_, err := base64.StdEncoding.Decode(serviceAccount, []byte(base64ServiceAccount))
	if err != nil {
		return err
	}

	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsJSON(serviceAccount)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	dbFirestore.appDB = app
	defer client.Close()
	return nil
}

// GetDocumentData get firestore document data
func (dbFirestore DBFirestore) GetDocumentData(collection string, document string) (map[string]interface{}, error) {
	client, err := dbFirestore.appDB.Firestore(context.Background())
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
