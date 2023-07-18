package pawntown_chess_results_module

import (
	"testing"
)

func TestChannelPathMatch(t *testing.T) {
	t.Run("With Params", func(t *testing.T) {
		tournamentId := "tnr790507"

		api := New(nil)
		_, err := api.Tournament(tournamentId)

		if err != nil {
			t.Errorf("api.Tournament(%s) = ([...], err), want ([...], nil)", tournamentId)
		}
	})
}
