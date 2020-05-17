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

func TestLoadConfigWhenGetDocumentDataReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	c := config.NewConfigs(dbMock)

	app := "myapp"
	keys := []string{"config1", "config2"}

	expect := errors.New("access error")
	dbMock.EXPECT().
		GetDocumentData("configs", app).
		Return(nil, expect)

	got := c.LoadConfig(app, keys)

	if diff := deep.Equal(got, expect); diff != nil {
		t.Error(diff)
	}
}

func TestLoadConfigWhenGetDocumentDataReturnsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	c := config.NewConfigs(dbMock)

	app := "myapp"
	keys := []string{"config1", "config2"}

	expect := errors.New("no data in config")
	dbMock.EXPECT().
		GetDocumentData("configs", app).
		Return(nil, nil)

	got := c.LoadConfig(app, keys)

	if diff := deep.Equal(got, expect); diff != nil {
		t.Error(diff)
	}
}

func TestLoadConfigWhenConfigIsMissing(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	c := config.NewConfigs(dbMock)

	app := "myapp"
	keys := []string{"config1", "config2"}

	expect := fmt.Errorf("missing config %s", "config2")
	dbMock.EXPECT().
		GetDocumentData("configs", app).
		Return(map[string]interface{}{
			"config1": "vla",
		}, nil)

	got := c.LoadConfig(app, keys)

	if diff := deep.Equal(got, expect); diff != nil {
		t.Error(diff)
	}
}

func TestLoadConfigAllOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	c := config.NewConfigs(dbMock)

	app := "myapp"
	keys := []string{"config1", "config2", "config3", "config4", "config5", "config6", "config7", "config8", "config9"}

	config1 := "jahwidh93u"
	config2 := 12491
	config3 := int64(63123)
	config4 := float64(861792831)
	config5 := 12491
	config6 := int64(1263123)
	config7 := float64(861792831)
	config8 := "1263123"
	config9 := "861792831"

	dbMock.EXPECT().
		GetDocumentData("configs", app).
		Return(map[string]interface{}{
			"config1": config1,
			"config2": config2,
			"config3": config3,
			"config4": config4,
			"config5": config5,
			"config6": config6,
			"config7": config7,
			"config8": config8,
			"config9": config9,
		}, nil)

	got := c.LoadConfig(app, keys)

	if got != nil {
		t.Errorf("Error not expected %v, nil expected", got)
	} else if diff := deep.Equal(c.GetConfigAsStr("config1"), config1); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(c.GetConfigAsInt("config2"), config2); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(c.GetConfigAsInt("config3"), int(config3)); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(c.GetConfigAsInt("config4"), int(config4)); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(c.GetConfigAsInt64("config5"), int64(config5)); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(c.GetConfigAsInt64("config6"), config6); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(c.GetConfigAsInt64("config7"), int64(config7)); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(c.GetConfigAsInt("config8"), 0); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(c.GetConfigAsInt64("config9"), int64(0)); diff != nil {
		t.Error(diff)
	}
}

func TestLoadConfigAllOkOverrideEnvVar(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	c := config.NewConfigs(dbMock)

	app := "myapp"
	keys := []string{"config1", "new_config"}

	config1 := "jahwidh93u"
	config3 := "2jf1023"
	os.Setenv("new_config", "ejw19208o")

	dbMock.EXPECT().
		GetDocumentData("configs", app).
		Return(map[string]interface{}{
			"config1":    config1,
			"new_config": config3,
		}, nil)

	got := c.LoadConfig(app, keys)

	if got != nil {
		t.Errorf("Error not expected %v, nil expected", got)
	} else if diff := deep.Equal(c.GetConfigAsStr("config1"), config1); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(c.GetConfigAsStr("new_config"), "ejw19208o"); diff != nil {
		t.Error(diff)
	}
}

func TestLoadConfigAllOkRefreshConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	c := config.NewConfigs(dbMock)

	app := "myapp"
	keys := []string{"config1", "config2"}

	config1 := "jahwidh93u"
	config2 := 12491
	config3 := 280312

	dbMock.EXPECT().
		GetDocumentData("configs", app).
		Return(map[string]interface{}{
			"config1": config1,
			"config2": config2,
			config.RefreshConfigTimeoutInSeconds: 1,
		}, nil)

	dbMock.EXPECT().
		GetDocumentData("configs", app).
		Return(map[string]interface{}{
			"config1": config1,
			"config2": config3,
			config.RefreshConfigTimeoutInSeconds: 300,
		}, nil)

	got := c.LoadConfig(app, keys)

	time.Sleep(2 * time.Second)

	if got != nil {
		t.Errorf("Error not expected %v, nil expected", got)
	} else if diff := deep.Equal(c.GetConfigAsStr("config1"), config1); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(c.GetConfigAsInt("config2"), config3); diff != nil {
		t.Error(diff)
	}
}

func TestLoadConfigAllOkRefreshConfigError(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	c := config.NewConfigs(dbMock)

	app := "myapp"
	keys := []string{"config1", "config2"}

	config1 := "jahwidh93u"
	config2 := 12491

	dbMock.EXPECT().
		GetDocumentData("configs", app).
		Return(map[string]interface{}{
			"config1": config1,
			"config2": config2,
			config.RefreshConfigTimeoutInSeconds: 1,
		}, nil)

	dbMock.EXPECT().
		GetDocumentData("configs", app).
		Return(map[string]interface{}{
			"config1": config1,
		}, nil)

	got := c.LoadConfig(app, keys)

	time.Sleep(2 * time.Second)

	if got != nil {
		t.Errorf("Error not expected %v, nil expected", got)
	} else if diff := deep.Equal(c.GetConfigAsStr("config1"), config1); diff != nil {
		t.Error(diff)
	} else if diff := deep.Equal(c.GetConfigAsInt("config2"), config2); diff != nil {
		t.Error(diff)
	}
}

func TestGetConfigAsIntEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	c := config.NewConfigs(dbMock)

	if diff := deep.Equal(c.GetConfigAsInt(config.TokenExpirationInMinutes), 0); diff != nil {
		t.Error(diff)
	}
}

func TestGetConfigAsInt64Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	c := config.NewConfigs(dbMock)

	if diff := deep.Equal(c.GetConfigAsInt64(config.TokenExpirationInMinutes), int64(0)); diff != nil {
		t.Error(diff)
	}
}

func TestGetConfigAsStrEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_firestore.NewMockIDBFirestore(ctrl)
	c := config.NewConfigs(dbMock)

	if diff := deep.Equal(c.GetConfigAsStr(config.GatewayPublicKey), ""); diff != nil {
		t.Error(diff)
	}
}
