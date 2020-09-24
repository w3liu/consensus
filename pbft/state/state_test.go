package state

import (
	"github.com/w3liu/consensus/pbft/config"
	"testing"
)

func TestCheckMajor23(t *testing.T) {
	cfg, err := config.Init("../config/config.toml")
	if err != nil {
		t.Fatal(err)
	}
	s := NewState(cfg)

	t.Log(s.checkMajor23(3))
}
