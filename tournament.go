package pawntown_chess_results_module

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Tournament struct {
	api  *Api
	id   string
	data TournamentData
}

type TournamentData struct {
	Name               string
	Organizer          string
	Federation         string
	Director           string
	ChiefArbiter       string
	DeputyChiefArbiter string
	TimeControl        string
	Location           string
	NumberOfRounds     int
	TournamentMode     string
	EloCalculation     string
	AverageElo         int
	FromTo             string
	Rounds             []Round
	Players            []Player
}

func parsePlayerRow(rowData string) (Player, error) {
	player := Player{}

	playerCellRegex := regexp.MustCompile(`\<td[^>]*\>(.*)\<\/td\>`)
	playerCells := playerCellRegex.FindAllStringSubmatch(rowData, -1)
	for index, match := range playerCells {
		switch index {
		case 0:
			num, err := strconv.Atoi(match[1])
			if err != nil {
				return Player{}, err
			}
			player.Num = num
		case 1:
			title, err := parseFideTitle(match[1])
			if err != nil {
				return Player{}, err
			}
			player.FideTitle = title
		case 2:
			name := match[1]
			if len(name) < 1 {
				return Player{}, errors.New("no name is present")
			}
			player.Name = name
		case 3:
			fideId, err := strconv.Atoi(match[1])
			if err != nil {
				return Player{}, err
			}
			player.FideID = fideId
		case 4:
			player.Country = match[1]
		case 5:
			player.Club = match[1]
		}
	}

	return player, nil
}

func newTournament(api *Api, id string) (*Tournament, error) {
	tournament := &Tournament{
		api: api,
		id:  id,
	}

	if err := tournament.Refresh(); err != nil {
		return nil, err
	}

	return tournament, nil
}

func (t *Tournament) Refresh() error {
	t.data = TournamentData{}

	rawData, err := t.api.get(fmt.Sprintf("https://chess-results.com/%s.aspx?turdet=YES", t.id))
	if err != nil {
		return err
	}

	// Tournament Information
	reName := regexp.MustCompile(`<h2>(.*?)</h2>`)
	reOrganizer := regexp.MustCompile(`<tr><td[^>]*>Veranstalter</td><td[^>]*>(.*?)</td></tr>`)
	reFederation := regexp.MustCompile(`<tr><td[^>]*>F&ouml;deration</td><td[^>]*>(.*?)</td></tr>`)
	reDirector := regexp.MustCompile(`<tr><td[^>]*>Turnierdirektor</td><td[^>]*>(.*?)</td></tr>`)
	reChiefArbiter := regexp.MustCompile(`<tr[^>]*><td[^>]*>Hauptschiedsrichter</td><td[^>]*>(.*?)</td></tr>`)
	reDeputyChiefArbiter := regexp.MustCompile(`<tr[^>]*><td[^>]*>Deputy Hauptschiedsrichter</td><td[^>]*>(.*?)</td></tr>`)
	reTimeControl := regexp.MustCompile(`<tr><td[^>]*>Bedenkzeit \(Standard\)</td><td[^>]*>(.*?)</td></tr>`)
	reLocation := regexp.MustCompile(`<tr><td[^>]*>Ort</td><td[^>]*>(.*?)</td></tr>`)
	reNumberOfRounds := regexp.MustCompile(`<tr><td[^>]*>Rundenanzahl</td><td[^>]*>(.*?)</td></tr>`)
	reTournamentMode := regexp.MustCompile(`<tr><td[^>]*>Turniermodus</td><td[^>]*>(.*?)</td></tr>`)
	reEloCalculation := regexp.MustCompile(`<tr><td[^>]*>Elorechnung</td><td[^>]*>(.*?)</td></tr>`)
	reAverageElo := regexp.MustCompile(`<tr><td[^>]*>Eloschnitt</td><td[^>]*>(.*?)</td></tr>`)
	reFromTo := regexp.MustCompile(`<tr><td[^>]*>Von</td><td[^>]*>(.*?)</td></tr>`)

	// Match and assign the values
	t.data.Name = reName.FindStringSubmatch(rawData)[1]
	t.data.Organizer = reOrganizer.FindStringSubmatch(rawData)[1]
	t.data.Federation = reFederation.FindStringSubmatch(rawData)[1]
	t.data.Director = reDirector.FindStringSubmatch(rawData)[1]
	t.data.ChiefArbiter = reChiefArbiter.FindStringSubmatch(rawData)[1]
	t.data.DeputyChiefArbiter = reDeputyChiefArbiter.FindStringSubmatch(rawData)[1]
	t.data.TimeControl = reTimeControl.FindStringSubmatch(rawData)[1]
	t.data.Location = reLocation.FindStringSubmatch(rawData)[1]
	t.data.NumberOfRounds, _ = strconv.Atoi(reNumberOfRounds.FindStringSubmatch(rawData)[1])
	t.data.TournamentMode = reTournamentMode.FindStringSubmatch(rawData)[1]
	t.data.EloCalculation = reEloCalculation.FindStringSubmatch(rawData)[1]
	t.data.AverageElo, _ = strconv.Atoi(reAverageElo.FindStringSubmatch(rawData)[1])
	t.data.FromTo = reFromTo.FindStringSubmatch(rawData)[1]

	// Player Information
	rePlayerRows := regexp.MustCompile(`\<tr[^>]*\>(.*)\<\/tr\>`)
	var players []Player
	playerRows := rePlayerRows.FindAllStringSubmatch(rawData, -1)
	for _, match := range playerRows {
		player, err := parsePlayerRow(match[1])
		if err != nil {
			players = append(players, player)
		}
	}
	t.data.Players = players

	return nil
}
