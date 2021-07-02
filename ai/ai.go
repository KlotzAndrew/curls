package ai

import "curls/models"

func NextMove(game models.GameRequest) models.MoveResponse {
	return models.MoveResponse{Move: models.Up, Shout: ""}
}
