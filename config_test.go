package rage_test

import (
	"os"
	"strings"
	"testing"

	"github.com/faiuwle/go-parser-game/rage"
	"github.com/google/go-cmp/cmp"
)

func TestReadValidEntitiesFromConfig(t *testing.T) {
	configFile, err := os.Open("testdata/testrooms.json")
	if err != nil {
		t.Fatal(err)
	}
	defer configFile.Close()

	want := commonGameData
	var got rage.GameData

	got, err = rage.ReadConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestReadConfigErrorsOnInvalidJson(t *testing.T) {
	_, err := rage.ReadConfig(strings.NewReader("}"))
	if err == nil {
		t.Fatal("got no error, config is valid json")
	}
}
