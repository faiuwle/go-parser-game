package rage_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/faiuwle/go-parser-game/rage"
	"github.com/google/go-cmp/cmp"
)

func TestReadValidEntitiesFromConfig(t *testing.T) {
	t.Parallel()

	path := t.TempDir() + "/rage_data.json"
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(getCommonGameData())
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	configFile, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer configFile.Close()

	want := getCommonGameData()
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
