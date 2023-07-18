package pawntown_chess_results_module

type Round struct {
	Round    int
	Pairings []Pairing
}

type Pairing struct {
	White  *PairingPlayerDetails
	Black  *PairingPlayerDetails
	Result *GameResult
}

type PairingPlayerDetails struct {
	Player *Player
	Points int // Points before the round started
}
