package seed

import "final-by-me/internal/models"

func EPLTeams() []models.Team {
	return []models.Team{
		{Code: "ARS", Name: "Arsenal"},
		{Code: "AVL", Name: "Aston Villa"},
		{Code: "BOU", Name: "Bournemouth"},
		{Code: "BRE", Name: "Brentford"},
		{Code: "BHA", Name: "Brighton"},
		{Code: "CHE", Name: "Chelsea"},
		{Code: "CRY", Name: "Crystal Palace"},
		{Code: "EVE", Name: "Everton"},
		{Code: "FUL", Name: "Fulham"},
		{Code: "IPS", Name: "Ipswich Town"},
		{Code: "LEI", Name: "Leicester City"},
		{Code: "LIV", Name: "Liverpool"},
		{Code: "MCI", Name: "Manchester City"},
		{Code: "MUN", Name: "Manchester United"},
		{Code: "NEW", Name: "Newcastle United"},
		{Code: "NFO", Name: "Nottingham Forest"},
		{Code: "SOU", Name: "Southampton"},
		{Code: "TOT", Name: "Tottenham Hotspur"},
		{Code: "WHU", Name: "West Ham United"},
		{Code: "WOL", Name: "Wolves"},
	}
}
