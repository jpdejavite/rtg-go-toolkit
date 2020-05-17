package config_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
	mock_firestore "github.com/jpdejavite/rtg-go-toolkit/mock/firestore"
	"github.com/jpdejavite/rtg-go-toolkit/pkg/config"
)

func TestGetGlobalKeys(t *testing.T) {
	gc := config.NewGlobalConfigs(nil)
	got := gc.GetGlobalKeys()
	expect := []string{config.GatewayPublicKey, config.TokenExpirationInMinutes}
	if diff := deep.Equal(got, expect); diff != nil {
		t.Error(diff)
	}
}

func TestLoadGlobalConfigWhenGetDocumentDataReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	gc := config.NewGlobalConfigs(dbMock)

	expect := errors.New("access error")
	dbMock.EXPECT().
		GetDocumentData("configs", "global").
		Return(nil, expect)

	got := gc.LoadGlobalConfig()

	if diff := deep.Equal(got, expect); diff != nil {
		t.Error(diff)
	}
}

func TestLoadGlobalConfigWhenGetDocumentDataReturnsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	gc := config.NewGlobalConfigs(dbMock)

	expect := errors.New("no data in global config")
	dbMock.EXPECT().
		GetDocumentData("configs", "global").
		Return(nil, nil)

	got := gc.LoadGlobalConfig()

	if diff := deep.Equal(got, expect); diff != nil {
		t.Error(diff)
	}
}

func TestLoadGlobalConfigWhenGlobalConfigIsMissing(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	gc := config.NewGlobalConfigs(dbMock)

	expect := fmt.Errorf("missing global config %s", config.TokenExpirationInMinutes)
	dbMock.EXPECT().
		GetDocumentData("configs", "global").
		Return(map[string]interface{}{
			config.GatewayPublicKey: "vla",
		}, nil)

	got := gc.LoadGlobalConfig()

	if diff := deep.Equal(got, expect); diff != nil {
		t.Error(diff)
	}
}

func TestLoadGlobalConfigAllOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	gc := config.NewGlobalConfigs(dbMock)

	gatewayPublicKey := "GatewayPublicKey"
	tokenExpirationInMinutes := 2

	dbMock.EXPECT().
		GetDocumentData("configs", "global").
		Return(map[string]interface{}{
			config.GatewayPublicKey:         gatewayPublicKey,
			config.TokenExpirationInMinutes: tokenExpirationInMinutes,
		}, nil)

	got := gc.LoadGlobalConfig()

	if got != nil {
		t.Errorf("Error not expetcted %v, nil expected", got)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsStr(config.GatewayPublicKey), gatewayPublicKey); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsInt(config.TokenExpirationInMinutes), tokenExpirationInMinutes); diff != nil {
		t.Error(diff)
	}
}

func TestGetGlobalConfigAsIntEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	gc := config.NewGlobalConfigs(dbMock)

	if diff := deep.Equal(gc.GetGlobalConfigAsInt(config.TokenExpirationInMinutes), 0); diff != nil {
		t.Error(diff)
	}
}

func TestGetGlobalConfigAsStrEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	gc := config.NewGlobalConfigs(dbMock)

	if diff := deep.Equal(gc.GetGlobalConfigAsStr(config.GatewayPublicKey), ""); diff != nil {
		t.Error(diff)
	}
}
