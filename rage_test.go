package rage_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/faiuwle/go-parser-game/rage"
)

func TestFormatItems(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name  string
		items []string
		want  string
	}{
		{
			name:  "Three items",
			items: []string{"one", "two", "three"},
			want:  "one, two, and three",
		},
		{
			name:  "Two items",
			items: []string{"one", "two"},
			want:  "one and two",
		},
		{
			name:  "One item",
			items: []string{"one"},
			want:  "one",
		},
		{
			name:  "No items",
			items: []string{},
			want:  "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			list := rage.FormatItems(tc.items)

			if list != tc.want {
				t.Errorf("Wanted %q, got %q", tc.want, list)
			}
		})
	}
}

func TestListExits(t *testing.T) {
	t.Parallel()
	room := rage.Entity{
		Exits: map[string]rage.Exit{
			"north": {},
			"east":  {},
		},
	}

	exitString := rage.ListExits(room)
	want := "You can go east and north from here."

	if exitString != want {
		t.Errorf("Wanted %q, got %q", want, exitString)
	}
}

func TestListContents(t *testing.T) {
	t.Parallel()
	room := rage.Entity{
		Contents: []string{"key", "phone", "chocolate"},
	}

	got := room.ListContents()
	want := "You see here key, phone, and chocolate."

	if got != want {
		t.Errorf("Wanted %q, got %q", want, got)
	}
}

func TestListContentsOmitPlayer(t *testing.T) {
	t.Parallel()
	g, err := rage.NewGame(getCommonGameData(), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	livingRoom := g.GetEntity("Living Room")
	got := livingRoom.ListContents()
	want := "You see here key."
	if got != want {
		t.Errorf("Wanted %q, got %q", want, got)
	}
}

func TestListEmptyContents(t *testing.T) {
	t.Parallel()
	g, err := rage.NewGame(getCommonGameData(), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	bathroom := g.GetEntity("Bathroom")
	got := bathroom.ListContents()
	want := ""
	if got != want {
		t.Errorf("Wanted %q, got %q", want, got)
	}
}

func TestTakeItemSucceedsIfItemIsPresent(t *testing.T) {
	t.Parallel()
	g, err := rage.NewGame(getCommonGameData(), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	response := g.TakeItem("key")
	want := "You take the key."

	if response != want {
		t.Errorf("Wanted %q, got %q", want, response)
	}
	if !g.Player.Contains("key") {
		t.Error("Key was not transferred to inventory.")
	}
	if g.PlayerLocation().Contains("key") {
		t.Errorf("Key is still in the room.")
	}
}

func TestTakeItemFailsIfItemIsNotPresent(t *testing.T) {
	t.Parallel()
	g, err := rage.NewGame(getCommonGameData(), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	response := g.TakeItem("phone")
	want := "I can't see that here."

	if response != want {
		t.Errorf("Wanted %q, got %q", want, response)
	}

	if g.Player.Contains("phone") {
		t.Errorf("Player took the phone.")
	}
}

func TestPlayerLocationReturnsNameOfRoomWherePlayerIs(t *testing.T) {
	t.Parallel()
	game, err := rage.NewGame(getCommonGameData(), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	got := game.PlayerLocation()
	if got.Name != "Living Room" {
		t.Errorf("expected player to be in \"Living Room\" but got %q", got.Name)
	}
}

func TestPlayerCanUseExitToMoveBetweenRooms(t *testing.T) {
	t.Parallel()
	game, err := rage.NewGame(getCommonGameData(), io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	cmd := rage.Command{Action: "south"}
	start := game.PlayerLocation()
	err = game.Do(cmd)
	finish := game.PlayerLocation()
	if err != nil {
		t.Fatal(err)
	}

	if start == finish {
		t.Error("player did not move rooms")
	}
}

func TestPlayerCannotPassExitWithoutKey(t *testing.T) {
	t.Parallel()

	writer := bytes.Buffer{}
	game, err := rage.NewGame(getCommonGameData(), &writer)
	if err != nil {
		t.Fatal(err)
	}

	cmd := rage.Command{Action: "north"}
	start := game.PlayerLocation()
	_ = game.Do(cmd)
	finish := game.PlayerLocation()

	want := "The door appears to be locked."
	got := strings.TrimSpace(writer.String())
	if got != want {
		t.Error(cmp.Diff(want, got))
	}

	if start != finish {
		t.Error("player was able to move through exit without key")
	}
}

func TestPlayerCannotPassExitWithoutKeyAndSeesDefaultFailureMessage(t *testing.T) {
	t.Parallel()
	writer := bytes.Buffer{}
	game, err := rage.NewGame(getCommonGameData(), &writer)
	if err != nil {
		t.Fatal(err)
	}

	cmd := rage.Command{Action: "east"}
	start := game.PlayerLocation()
	_ = game.Do(cmd)
	finish := game.PlayerLocation()

	want := rage.DefaultExitFailureMessage
	got := strings.TrimSpace(writer.String())
	if got != want {
		t.Error(cmp.Diff(want, got))
	}

	if start != finish {
		t.Error("player was able to move through exit without key")
	}
}

func TestGetEntityPanicsWhenGivenNonExistentEntity(t *testing.T) {
	t.Parallel()
	game, err := rage.NewGame(getCommonGameData(), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("GetEntity did not panic")
		}
	}()
	game.GetEntity("whatever")
}

func TestGetEntityReturnsEntityWhenItExists(t *testing.T) {
	t.Parallel()
	want := &rage.Entity{
		Name:     "player",
		Location: "Living Room",
		Kind:     "Character",
	}
	game, err := rage.NewGame(getCommonGameData(), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	got := game.GetEntity("player")
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestNewGameReturnsErrorWithInconsistentData(t *testing.T) {
	t.Parallel()
	testCases := map[string]rage.GameData{
		"non-existent location": {
			"object": rage.Entity{
				Name:     "object",
				Location: "non-existent",
				Kind:     "Thing",
			},
			"player": {
				Name: "player",
			},
		},
		"non-existant object in room": {
			"room": rage.Entity{
				Name: "room",
				Contents: []string{
					"non-existent",
				},
				Kind: "Room",
			},
			"player": {
				Name: "player",
			},
		},
		"empty location on thing": {
			"object": rage.Entity{
				Name:     "object",
				Location: "",
				Kind:     "Thing",
			},
			"player": {
				Name: "player",
			},
		},
		"non-empty location on room": {
			"room": rage.Entity{
				Name:     "room",
				Location: "somewhere",
				Kind:     "Room",
			},
			"somewhere": rage.Entity{
				Name:     "somewhere",
				Location: "room",
				Kind:     "Room",
			},
			"player": {
				Name: "player",
			},
		},
		"no player in game data": {},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := rage.NewGame(tc, io.Discard)
			if err == nil {
				t.Fatal("Did not fail with invalid game data")
			}
		})
	}
}

func TestNewGameCreatesGameFromConsistentData(t *testing.T) {
	t.Parallel()
	_, err := rage.NewGame(getCommonGameData(), io.Discard)
	if err != nil {
		t.Fatalf("Errored on consistent game data: %s", err)
	}
}

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

func getCommonGameData() rage.GameData {
	return rage.GameData{
		"Living Room": {
			Name:        "Living Room",
			Description: "The living room",
			Kind:        "Room",
			Exits: map[string]rage.Exit{
				"north": {
					Destination:    "Bedroom",
					Requires:       "key",
					FailureMessage: "The door appears to be locked.",
				},
				"south": {
					Destination: "Bathroom",
				},
				"east": {
					Destination: "Bedroom",
					Requires:    "key",
				},
			},
			Contents: []string{"key", "player"},
		},
		"Bedroom": {
			Name:        "Bedroom",
			Description: "The bedroom",
			Kind:        "Room",
		},
		"Bathroom": {
			Name:        "Bathroom",
			Description: "The bathroom",
			Kind:        "Room",
		},
		"key": {
			Name:     "key",
			Location: "Living Room",
			Kind:     "Thing",
		},
		"player": {
			Name:     "player",
			Location: "Living Room",
			Kind:     "Character",
		},
	}
}

func TestInvalidCommandFails(t *testing.T) {
	command := "take key and run"

	_, err := rage.Parse(command)

	if err != rage.ErrorInvalidCommand {
		t.Error(err)
	}
}

func TestTakeKeySucceeds(t *testing.T) {
	command := "take key"

	cmd, err := rage.Parse(command)

	if cmd.Action != "take" || cmd.Noun != "key" {
		t.Errorf("Command did not parse successfully")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestSetupBuildDir_CreatesTmpDirWithJSONAndGoAndReturnsPath(t *testing.T) {
	t.Parallel()
	buildPath, err := rage.SetupBuildDir("testdata/adventure.json")
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(buildPath + "/adventure.json")
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(buildPath + "/main.go")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompileProducesBinaryGivenJsonDataReturnsNoError(t *testing.T) {
	binPath := t.TempDir() + "/adventure"

	// TODO create compile method that returns executable, test that file functions properly
	err := rage.Compile("testdata/adventure.json", binPath)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(binPath)
	want := "The living room"
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	got := string(output)
	if !strings.Contains(got, want) {
		t.Errorf("want output to contain %q, but got:\n%s", want, got)
	}
}

func TestExecGoBuild_ProducesBinaryThatWeCanRunAndSeeCorrectOutput(t *testing.T) {
	t.Parallel()
	buildDir := "testdata/build"
	binPath := t.TempDir() + "/adventure"
	err := rage.ExecGoBuild(buildDir, binPath)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(binPath)
	want := "The living room"
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	got := string(output)
	if !strings.Contains(got, want) {
		t.Errorf("want output to contain %q, but got:\n%s", want, got)
	}
}
