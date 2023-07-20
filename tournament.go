package pawntown_chess_results_module

import (
	"errors"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strconv"
)

type Tournament struct {
	api          *Api
	id           string
	data         TournamentData
	secInputData url.Values
}

type TournamentData struct {
	RawMeta map[string]string

	Name               string
	Organizer          string
	Federation         string
	Director           string
	ChiefArbiter       string
	DeputyChiefArbiter string
	TimeControl        string
	Location           string
	NumberOfRounds     string
	TournamentMode     string
	EloCalculation     string
	AverageElo         string
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

func (t *Tournament) UpdateSecurityData() error {
	rawData, err := t.api.get(fmt.Sprintf("/%s.aspx", t.id))
	if err != nil {
		return err
	}

	secValues := url.Values{}
	secValues.Set("cb_alleDetails", "Turnierdetails anzeigen")
	secValues.Set("txt_name", "")

	reSecData := regexp.MustCompile(`<input type="hidden".*?name="([^"]*?)".*?value="([^"]*?)"`)
	inputsData := reSecData.FindAllStringSubmatch(rawData, -1)
	for _, match := range inputsData {
		secValues.Set(match[1], match[2])
	}
	t.secInputData = secValues
	return nil
}

func (t *Tournament) Refresh() error {
	t.data = TournamentData{}

	if err := t.UpdateSecurityData(); err != nil {
		return err
	}

	rawData, err := t.api.post(fmt.Sprintf("/%s.aspx?turdet=YES", t.id), t.secInputData)
	if err != nil {
		return err
	}

	// Tournament Information
	reName := regexp.MustCompile(`<h2>(.*?)</h2>`)
	reMetaReader := regexp.MustCompile(`<tr[^>]*?><td[^>]*>([^>]*?)</td><td[^>]*>(.*?)</td></tr>`)

	t.data.Name = reName.FindStringSubmatch(rawData)[1]
	t.data.RawMeta = map[string]string{}

	metaEntries := reMetaReader.FindAllStringSubmatch(rawData, -1)
	for _, match := range metaEntries {
		key := match[1]
		value := html.UnescapeString(match[2])
		if key != "" {
			t.data.RawMeta[key] = value

			switch key {
			case "Organizer(s)":
				t.data.Organizer = value
			case "Federation":
				t.data.Federation = value
			case "Tournament director":
				t.data.Director = value
			case "Chief Arbiter":
				t.data.ChiefArbiter = value
			case "Deputy Chief Arbiter":
				t.data.DeputyChiefArbiter = value
			case "Time control (Standard)":
				t.data.TimeControl = value
			case "Location":
				t.data.Location = value
			case "Number of rounds":
				t.data.NumberOfRounds = value
			case "Tournament type":
				t.data.TournamentMode = value
			case "Rating calculation":
				t.data.EloCalculation = value
			case "Rating-Ø":
				t.data.AverageElo = value
			case "Date":
				t.data.FromTo = value
			}
		}
	}

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
