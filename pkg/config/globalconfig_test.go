package config_test

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
	mock_firestore "github.com/jpdejavite/rtg-go-toolkit/mock/firestore"
	"github.com/jpdejavite/rtg-go-toolkit/pkg/config"
)

func TestGetGlobalKeys(t *testing.T) {
	gc := config.NewGlobalConfigs(nil)
	got := gc.GetGlobalKeys()
	expect := []string{config.GatewayPublicKey, config.TokenExpirationInMinutes, config.RefreshConfigTimeoutInSeconds}
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
	refreshConfigTimeoutInSeconds := 300
	os.Setenv(config.GatewayPublicKey, "")

	dbMock.EXPECT().
		GetDocumentData("configs", "global").
		Return(map[string]interface{}{
			config.GatewayPublicKey:              gatewayPublicKey,
			config.TokenExpirationInMinutes:      tokenExpirationInMinutes,
			config.RefreshConfigTimeoutInSeconds: refreshConfigTimeoutInSeconds,
		}, nil)

	got := gc.LoadGlobalConfig()

	if got != nil {
		t.Errorf("Error not expected %v, nil expected", got)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsStr(config.GatewayPublicKey), gatewayPublicKey); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsInt(config.TokenExpirationInMinutes), tokenExpirationInMinutes); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsInt(config.RefreshConfigTimeoutInSeconds), refreshConfigTimeoutInSeconds); diff != nil {
		t.Error(diff)
	}
}

func TestLoadGlobalConfigOverrideEnvVar(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	gc := config.NewGlobalConfigs(dbMock)

	gatewayPublicKey := "GatewayPublicKey"
	tokenExpirationInMinutes := 2
	refreshConfigTimeoutInSeconds := 300

	os.Setenv(config.GatewayPublicKey, "huahuahu")

	dbMock.EXPECT().
		GetDocumentData("configs", "global").
		Return(map[string]interface{}{
			config.GatewayPublicKey:              gatewayPublicKey,
			config.TokenExpirationInMinutes:      tokenExpirationInMinutes,
			config.RefreshConfigTimeoutInSeconds: refreshConfigTimeoutInSeconds,
		}, nil)

	got := gc.LoadGlobalConfig()

	if got != nil {
		t.Errorf("Error not expected %v, nil expected", got)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsStr(config.GatewayPublicKey), "huahuahu"); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsInt(config.TokenExpirationInMinutes), tokenExpirationInMinutes); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsInt(config.RefreshConfigTimeoutInSeconds), refreshConfigTimeoutInSeconds); diff != nil {
		t.Error(diff)
	}
}

func TestLoadGlobalConfigRefreshData(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	gc := config.NewGlobalConfigs(dbMock)

	gatewayPublicKey := "GatewayPublicKey"
	tokenExpirationInMinutes := 2
	refreshConfigTimeoutInSeconds := 300
	os.Setenv(config.GatewayPublicKey, "")

	dbMock.EXPECT().
		GetDocumentData("configs", "global").
		Return(map[string]interface{}{
			config.GatewayPublicKey:              gatewayPublicKey,
			config.TokenExpirationInMinutes:      tokenExpirationInMinutes,
			config.RefreshConfigTimeoutInSeconds: 1,
		}, nil)

	dbMock.EXPECT().
		GetDocumentData("configs", "global").
		Return(map[string]interface{}{
			config.GatewayPublicKey:              gatewayPublicKey,
			config.TokenExpirationInMinutes:      tokenExpirationInMinutes,
			config.RefreshConfigTimeoutInSeconds: refreshConfigTimeoutInSeconds,
		}, nil)

	got := gc.LoadGlobalConfig()

	time.Sleep(2 * time.Second)

	if got != nil {
		t.Errorf("Error not expected %v, nil expected", got)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsStr(config.GatewayPublicKey), gatewayPublicKey); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsInt(config.TokenExpirationInMinutes), tokenExpirationInMinutes); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsInt(config.RefreshConfigTimeoutInSeconds), refreshConfigTimeoutInSeconds); diff != nil {
		t.Error(diff)
	}
}

func TestLoadGlobalConfigRefreshDataError(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	gc := config.NewGlobalConfigs(dbMock)

	gatewayPublicKey := "GatewayPublicKey"
	tokenExpirationInMinutes := 2
	os.Setenv(config.GatewayPublicKey, "")

	dbMock.EXPECT().
		GetDocumentData("configs", "global").
		Return(map[string]interface{}{
			config.GatewayPublicKey:              gatewayPublicKey,
			config.TokenExpirationInMinutes:      tokenExpirationInMinutes,
			config.RefreshConfigTimeoutInSeconds: 1,
		}, nil)

	dbMock.EXPECT().
		GetDocumentData("configs", "global").
		Return(map[string]interface{}{
			config.GatewayPublicKey:         gatewayPublicKey,
			config.TokenExpirationInMinutes: tokenExpirationInMinutes,
		}, nil)

	got := gc.LoadGlobalConfig()

	time.Sleep(2 * time.Second)

	if got != nil {
		t.Errorf("Error not expected %v, nil expected", got)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsStr(config.GatewayPublicKey), gatewayPublicKey); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsInt(config.TokenExpirationInMinutes), tokenExpirationInMinutes); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(gc.GetGlobalConfigAsInt(config.RefreshConfigTimeoutInSeconds), 1); diff != nil {
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
