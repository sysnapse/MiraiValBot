package valorant

import (
	"encoding/json"
	"fmt"
	"strings"
)

type (
	ApiError struct {
		Errors []Error
	}

	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Details string `json:"details"`
	}

	Resp struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data"`
		Errors []Error     `json:"errors"`
	}

	AccountDetails struct {
		PUuid         string            `json:"puuid"`
		Region        string            `json:"region"`
		AccountLevel  int               `json:"account_level"`
		Name          string            `json:"name"`
		Tag           string            `json:"tag"`
		Card          map[string]string `json:"card"`
		LastUpdate    string            `json:"last_update"`
		LastUpdateRaw int64             `json:"last_update_raw"`
	}

	Localization struct {
		LocalizeItems
		Version string `json:"version"`
	}

	LocalizeItems map[string]LocalizeItem

	LocalizeItem struct {
		Name           string            `json:"name"`
		LocalizedNames map[string]string `json:"localizedNames"`
	}

	MatchData struct {
		MetaData MatchMetaData             `json:"metadata"`
		Players  []MatchPlayer             `json:"players"`
		Teams    map[string]MatchTeamStats `json:"teams"`
		Rounds   []MatchRound              `json:"rounds"`
	}

	MatchMetaData struct {
		Map              string `json:"map"`
		GameVersion      string `json:"game_version"`
		GameLength       int64  `json:"game_length"`
		GameStart        int64  `json:"game_start"`
		GameStartPatched string `json:"game_start_patched"`
		RoundsPlayed     int    `json:"rounds_played"`
		Mode             string `json:"mode"`
		Queue            string `json:"queue"`
		SeasonId         string `json:"season_id"`
		Platform         string `json:"platform"`
		MatchId          string `json:"matchid"`
		Region           string `json:"region"`
		Cluster          string `json:"cluster"`
	}

	MatchPlayer struct {
		PUuid              string           `json:"puuid"`
		Name               string           `json:"name"`
		Tag                string           `json:"tag"`
		Team               string           `json:"team"`
		Level              int              `json:"level"`
		Character          string           `json:"character"`
		CurrentTier        int              `json:"currenttier"`
		CurrentTierPatched int              `json:"currenttier_patched"`
		PlayerCard         string           `json:"player_card"`
		PlayerTitle        string           `json:"player_title"`
		PartyId            string           `json:"party_id"`
		SessionPlayTime    map[string]int64 `json:"session_playtime"`

		Assets struct {
			Card  map[string]string `json:"card"`
			Agent map[string]string `json:"agent"`
		} `json:"assets"`

		Behaviour struct {
			AfkRounds    int `json:"afk_rounds"`
			FriendlyFire struct {
				Incoming int `json:"incoming"`
				Outgoing int `json:"outgoing"`
			} `json:"friendly_fire"`
			RoundInSpawn int `json:"round_in_spawn"`
		} `json:"behaviour"`

		Platform struct {
			Type string `json:"type"`
			OS   struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"os"`
		} `json:"platform"`

		AbilityCasts map[string]int `json:"ability_casts"`

		Stats struct {
			Score     int `json:"score"`
			Kills     int `json:"kills"`
			Deaths    int `json:"deaths"`
			Assists   int `json:"assists"`
			BodyShots int `json:"bodyshots"`
			Headshots int `json:"headshots"`
			LegShots  int `json:"legshots"`
		} `json:"stats"`

		Economy struct {
			Spent        EconomyInfo `json:"spent"`
			LoadoutValue EconomyInfo `json:"loadout_value"`
		} `json:"economy"`

		DamageMade     int64 `json:"damage_made"`
		DamageReceived int64 `json:"damage_received"`
	}

	EconomyInfo struct {
		OverAll int64 `json:"overall"`
		Overage int64 `json:"overage"`
	}

	MatchTeamStats struct {
		HasWon     bool `json:"has_won"`
		RoundsWon  int  `json:"rounds_won"`
		RoundsLost int  `json:"rounds_lost"`
	}

	MatchRound struct {
		WinningTeam string `json:"winning_team"`
		EndType     string `json:"end_type"`
		BombPlanted bool   `json:"bomb_planted"`
		BombDefused bool   `json:"bomb_defused"`
		PlantEvents struct {
			PlantLocation          Location         `json:"plant_location"`
			PlantedBy              EventOwner       `json:"planted_by"`
			PlantSite              string           `json:"plant_site"`
			PlantTimeInRound       int64            `json:"plant_time_in_round"`
			PlayerLocationsOnPlant []PlayerLocation `json:"player_locations_on_plant"`
		} `json:"plant_events"`
		DefuseEvents struct {
			DefuseLocation         Location         `json:"defuse_location"`
			DefusedBy              EventOwner       `json:"defused_by"`
			DefuseTimeInRound      int64            `json:"defuse_time_in_round"`
			PlayerLocationOnDefuse []PlayerLocation `json:"player_locations_on_defuse"`
		}
		PlayerStats []MatchPlayerStats `json:"player_stats"`
	}

	MatchPlayerStats struct {
		AbilityCasts map[string]int `json:"ability_casts"`
		TeamPlayer
		DamageEvents []interface{} `json:"damage_events"`
		Damage       int           `json:"damage"`
		BodyShots    int           `json:"bodyshots"`
		Headshots    int           `json:"headshots"`
		LegShots     int           `json:"legshots"`
		KillsEvents  []interface{} `json:"kills_events"`
		Kills        int           `json:"kills"`
		Score        int           `json:"score"`
		Economy      struct {
			LoadoutValue int `json:"loadout_value"`
			Weapon       struct {
				Weapon Equipment `json:"weapon"`
				Armor  Equipment `json:"armor"`
			} `json:"weapon"`
			Remaining int `json:"remaining"`
			Spent     int `json:"spent"`
		}
		WasAfk        bool `json:"was_afk"`
		WasPenalized  bool `json:"was_penalized"`
		StayedInSpawn bool `json:"stayed_in_spawn"`
	}

	Equipment struct {
		Id     string            `json:"id"`
		Name   string            `json:"name"`
		Assets map[string]string `json:"assets"`
	}

	EventOwner struct {
		PUuid       string `json:"puuid"`
		DisplayName string `json:"display_name"`
		Team        string `json:"team"`
	}

	TeamPlayer struct {
		PlayerPUuid       string `json:"player_puuid"`
		PlayerDisplayName string `json:"player_display_name"`
		PlayerTeam        string `json:"player_team"`
	}

	PlayerLocation struct {
		TeamPlayer
		Location    Location `json:"location"`
		ViewRadians float64  `json:"view_radians"`
	}

	Location struct {
		X int64 `json:"x"`
		Y int64 `json:"y"`
	}

	GameStatus struct {
		Maintenances []MaintainInfo `json:"maintenances"`
		Incidents    []MaintainInfo `json:"incidents"`
	}

	MaintainInfo struct {
		CreatedAt string `json:"created_at"`
		ArchiveAt string `json:"archive_at"`
		Updates   []struct {
			CreatedAt        string        `json:"created_at"`
			ArchiveAt        string        `json:"archive_at"`
			Publish          bool          `json:"publish"`
			Id               int           `json:"id"`
			Translations     []I18NContent `json:"translations"`
			PublishLocations []string      `json:"publish_locations"`
			Author           string        `json:"author"`
		} `json:"updates"`
		Platforms         []string      `json:"platforms"`
		UpdatedAt         string        `json:"updated_at"`
		Id                int           `json:"id"`
		Titles            []I18NContent `json:"titles"`
		MaintenanceStatus string        `json:"maintenance_status"`
		IncidentSeverity  string        `json:"incident_severity"`
	}

	I18NContent struct {
		Content string `json:"content"`
		Locale  string `json:"locale"`
	}
)

func (resp *Resp) ParseData(t interface{}) error {
	return json.Unmarshal(resp.Data.([]byte), t)
}

func (err *ApiError) Error() string {
	messages := make([]string, len(err.Errors))
	for i, e := range err.Errors {
		messages[i] = e.Message
	}
	return fmt.Sprintf("API Errors: %s", strings.Join(messages, ", "))
}
