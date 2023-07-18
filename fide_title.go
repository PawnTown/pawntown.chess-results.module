package pawntown_chess_results_module

import (
	"errors"
	"fmt"
)

type FideTitle string

const (
	FideTitleGM   FideTitle = "GM"
	FideTitleIM   FideTitle = "IM"
	FideTitleFM   FideTitle = "FM"
	FideTitleCM   FideTitle = "CM"
	FideTitleWGM  FideTitle = "WGM"
	FideTitleWIM  FideTitle = "WIM"
	FideTitleWFM  FideTitle = "WFM"
	FideTitleWCM  FideTitle = "WCM"
	FideTitleNone FideTitle = ""
)

func parseFideTitle(title string) (FideTitle, error) {
	if !IsValidFideTitle(title) {
		return "", errors.New(fmt.Sprintf("'%s' is not a valid fide title", title))
	}

	return FideTitle(title), nil
}

func IsValidFideTitle(title string) bool {
	switch title {
	case string(FideTitleGM),
		string(FideTitleIM),
		string(FideTitleFM),
		string(FideTitleCM),
		string(FideTitleWGM),
		string(FideTitleWIM),
		string(FideTitleWFM),
		string(FideTitleWCM),
		string(FideTitleNone):
		return true
	default:
		return false
	}
}
