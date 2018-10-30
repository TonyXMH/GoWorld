package config

import (
	"encoding/json"
	"github.com/TonyXMH/GoWorld/gwlog"
	"testing"
)

func TestLoad(t *testing.T) {
	config := Get()
	gwlog.Debug("goworld config: \n%s", config)
	if config == nil {
		t.FailNow()
	}
	if config.dispatcher.Ip == "" {
		t.Errorf("dispatch ip not found")
	}
	if config.dispatcher.Port == 0 {
		t.Errorf("dispatcher port not found")
	}
	for gameid, gameConfig := range config.games {
		if gameConfig.Ip == "" {
			t.Errorf("game %d ip not found", gameid)
		}
		if gameConfig.Port == 0 {
			t.Errorf("game %d port not found", gameid)
		}
	}

	for gateid, gateConfig := range config.gates {
		if gateConfig.Ip == "" {
			t.Errorf("gate %d ip not found", gateid)
		}
		if gateConfig.Port == 0 {
			t.Errorf("gate %d port not found", gateid)
		}
	}
	gwlog.Info("read goworld config:%v", config)
}
func TestReload(t *testing.T) {
	config := Get()
	config = Reload()
	gwlog.Debug("goworld config:\n%s", config)
}

func TestGetDispatcher(t *testing.T) {
	cfg := GetDispatcher()
	cfgStr, _ := json.Marshal(cfg)
	gwlog.Info("dispatcher config:%s", string(cfgStr))
}

func TestGetGame(t *testing.T) {
	for id := 1; id <= 10; id++ {
		cfg := GetGame(id)
		if cfg == nil {
			gwlog.Error("Game %d not found", id)
		} else {
			gwlog.Info("Game %d config: %v", id, cfg)
		}
	}
}

func TestGetGate(t *testing.T) {
	for id := 1; id <= 10; id++ {
		cfg := GetGate(id)
		if cfg == nil {
			gwlog.Error("Gate %d not found", id)
		} else {
			gwlog.Info("Gate %d config: %v", id, cfg)
		}
	}
}
