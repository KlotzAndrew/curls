package ai_test

import (
	"curls/ai"
	"curls/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

func TestNextMove(t *testing.T) {
	tests := []struct {
		file string
		move models.Move
	}{
		// {"sample_move.json", models.Up},
		// {"tiny_move.json", models.Up},
		// {"small_spin.json", models.Up},
		// {"consider_food.json", models.Right},
		// {"search_food.json", models.Up},
		{"search_food_nearest.json", models.Up},
	}

	for _, tt := range tests {
		game := getFile(t, tt.file)
		move := ai.NextMove(game)

		if move.Move != tt.move {
			t.Errorf("got %s wanted %s for input %s", move, tt.move, tt.file)
		}
	}
}

func getFile(t *testing.T, name string) models.GameRequest {
	data, err := ioutil.ReadFile("fixtures/" + name)
	if err != nil {
		t.Fatal(err)
	}

	game := models.GameRequest{}
	err = json.Unmarshal(data, &game)
	if err != nil {
		log.Fatal(err)
	}

	return game
}
