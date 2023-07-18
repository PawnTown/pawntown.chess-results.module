package pawntown_chess_results_module

import (
	"errors"
	"fmt"
)

type GameResult string

const (
	GameResultWhiteWon FideTitle = "1 - 0"
	GameResultBlackWon FideTitle = "0 - 1"
	GameResultDraw     FideTitle = "½ - ½"
)

func parseGameResult(gameResult string) (GameResult, error) {
	if !IsValidGameResult(gameResult) {
		return "", errors.New(fmt.Sprintf("'%s' is not a valid game result", gameResult))
	}

	return GameResult(gameResult), nil
}

func IsValidGameResult(gameResult string) bool {
	switch gameResult {
	case string(GameResultWhiteWon),
		string(GameResultBlackWon),
		string(GameResultDraw):
		return true
	default:
		return false
	}
}
