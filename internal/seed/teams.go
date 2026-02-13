package seed

import "final-by-me/internal/models"

// Add more teams anytime.
func TeamsAll() []models.Team {
	return []models.Team{
		// EPL
		// EPL (20 teams)
		{Code: "ARS", Name: "Arsenal", League: "EPL"},
		{Code: "AVL", Name: "Aston Villa", League: "EPL"},
		{Code: "BOU", Name: "Bournemouth", League: "EPL"},
		{Code: "BRE", Name: "Brentford", League: "EPL"},
		{Code: "BHA", Name: "Brighton", League: "EPL"},
		{Code: "BUR", Name: "Burnley", League: "EPL"},
		{Code: "CHE", Name: "Chelsea", League: "EPL"},
		{Code: "CRY", Name: "Crystal Palace", League: "EPL"},
		{Code: "EVE", Name: "Everton", League: "EPL"},
		{Code: "FUL", Name: "Fulham", League: "EPL"},
		{Code: "LIV", Name: "Liverpool", League: "EPL"},
		{Code: "LUT", Name: "Luton Town", League: "EPL"},
		{Code: "MCI", Name: "Manchester City", League: "EPL"},
		{Code: "MUN", Name: "Manchester United", League: "EPL"},
		{Code: "NEW", Name: "Newcastle", League: "EPL"},
		{Code: "NFO", Name: "Nottingham Forest", League: "EPL"},
		{Code: "SHU", Name: "Sheffield United", League: "EPL"},
		{Code: "TOT", Name: "Tottenham", League: "EPL"},
		{Code: "WHU", Name: "West Ham", League: "EPL"},
		{Code: "WOL", Name: "Wolves", League: "EPL"},

		// LaLiga
		// LaLiga (20 teams)
		{Code: "RMA", Name: "Real Madrid", League: "LaLiga"},
		{Code: "FCB", Name: "Barcelona", League: "LaLiga"},
		{Code: "ATM", Name: "Atletico Madrid", League: "LaLiga"},
		{Code: "SEV", Name: "Sevilla", League: "LaLiga"},
		{Code: "VAL", Name: "Valencia", League: "LaLiga"},
		{Code: "RSO", Name: "Real Sociedad", League: "LaLiga"},
		{Code: "VIL", Name: "Villarreal", League: "LaLiga"},
		{Code: "ATH", Name: "Athletic Bilbao", League: "LaLiga"},
		{Code: "BET", Name: "Real Betis", League: "LaLiga"},
		{Code: "OSA", Name: "Osasuna", League: "LaLiga"},
		{Code: "GIR", Name: "Girona", League: "LaLiga"},
		{Code: "GET", Name: "Getafe", League: "LaLiga"},
		{Code: "ALM", Name: "Almeria", League: "LaLiga"},
		{Code: "CAD", Name: "Cadiz", League: "LaLiga"},
		{Code: "CEL", Name: "Celta Vigo", League: "LaLiga"},
		{Code: "MLL", Name: "Mallorca", League: "LaLiga"},
		{Code: "RAY", Name: "Rayo Vallecano", League: "LaLiga"},
		{Code: "GRA", Name: "Granada", League: "LaLiga"},
		{Code: "ALV", Name: "Alaves", League: "LaLiga"},
		{Code: "LAS", Name: "Las Palmas", League: "LaLiga"},

		// SerieA
		// SerieA (20 teams)
		{Code: "INT", Name: "Inter", League: "SerieA"},
		{Code: "ACM", Name: "AC Milan", League: "SerieA"},
		{Code: "JUV", Name: "Juventus", League: "SerieA"},
		{Code: "NAP", Name: "Napoli", League: "SerieA"},
		{Code: "ROM", Name: "Roma", League: "SerieA"},
		{Code: "LAZ", Name: "Lazio", League: "SerieA"},
		{Code: "ATA", Name: "Atalanta", League: "SerieA"},
		{Code: "FIO", Name: "Fiorentina", League: "SerieA"},
		{Code: "BOL", Name: "Bologna", League: "SerieA"},
		{Code: "TOR", Name: "Torino", League: "SerieA"},
		{Code: "UDI", Name: "Udinese", League: "SerieA"},
		{Code: "SAS", Name: "Sassuolo", League: "SerieA"},
		{Code: "MON", Name: "Monza", League: "SerieA"},
		{Code: "GEN", Name: "Genoa", League: "SerieA"},
		{Code: "EMP", Name: "Empoli", League: "SerieA"},
		{Code: "LEC", Name: "Lecce", League: "SerieA"},
		{Code: "CAG", Name: "Cagliari", League: "SerieA"},
		{Code: "VER", Name: "Verona", League: "SerieA"},
		{Code: "SAL", Name: "Salernitana", League: "SerieA"},
		{Code: "FRO", Name: "Frosinone", League: "SerieA"},

		// KPL (Kazakhstan Premier League)
		// KPL
		{Code: "AST", Name: "Astana", League: "KPL"},
		{Code: "KAI", Name: "Kairat", League: "KPL"},
		{Code: "ORD", Name: "Ordabasy", League: "KPL"},
		{Code: "AKT", Name: "Aktobe", League: "KPL"},
		{Code: "TOB", Name: "Tobol", League: "KPL"},
		{Code: "ATY", Name: "Atyrau", League: "KPL"},
		{Code: "KYZ", Name: "Kyzylzhar", League: "KPL"},
		{Code: "SHA", Name: "Shakhter Karagandy", League: "KPL"},
		{Code: "ZHE", Name: "Zhetysu", League: "KPL"},
		{Code: "KAS", Name: "Kaspiy", League: "KPL"},
		{Code: "MAK", Name: "Maktaaral", League: "KPL"},
		{Code: "TUR", Name: "Turan", League: "KPL"},
	}
}

func Leagues() []string {
	return []string{"EPL", "LaLiga", "SerieA", "KPL"}
}
