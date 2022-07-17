// Package mapper contains functions to map row types from pg.*
// to game.*
package mapper

import (
	pg "github.com/lordvidex/wordle-wf/internal/db/pg/gen"
	"github.com/lordvidex/wordle-wf/internal/game"
	"github.com/lordvidex/wordle-wf/internal/words"
)

// GetPlayersInGame helps to map row types from PlayerInGame to Player
// func GetPlayersInGame(players []*pg.) []*game.Player {
// 	playerSessions := make([]*game.Player, len(players))
// 	for i, gp := range players {
// 		var player game.Player
// 		player.ID = gp.ID
// 		player.Name = gp.Name
// 		player.User = &game.Player{
// 			ID: gp.UserID,
// 		}
// 		if gp.UserName.Valid {
// 			player.User.Name = gp.UserName.String
// 		}
// 		if gp.Password.Valid {
// 			player.User.Password = gp.Password.String
// 		}
// 		if gp.Email.Valid {
// 			player.User.Email = gp.Email.String
// 		}
// 		playerSessions[i] = &player
// 	}
// 	return playerSessions
// }

// FindByIdRow helps to map row types from pg.FindByIdRow to game.Game
func FindByIdRow(row *pg.FindByIdRow) *game.Game {
	gm := &game.Game{}
	gm.ID = row.ID
	gm.InviteID = row.InviteID
	gm.Word = words.Word{Word: row.Word}
	gm.StartTime = row.StartTime
	gm.PlayerCount = int(row.PlayerCount)

	if row.EndTime.Valid {
		gm.EndTime = &row.EndTime.Time
	}
	// check if each field in row is valid, and add
	// their values to gm.Settings
	if row.WordLength.Valid {
		gm.Settings.WordLength = int(row.WordLength.Int16)
	}
	if row.Trials.Valid {
		gm.Settings.Trials = int(row.Trials.Int16)
	}
	if row.MaxPlayerCount.Valid {
		gm.Settings.MaxPlayerCount = int(row.MaxPlayerCount.Int16)
	}
	if row.HasAnalytics.Valid {
		gm.Settings.Analytics = row.HasAnalytics.Bool
	}
	if row.ShouldRecordTime.Valid {
		gm.Settings.RecordTime = row.ShouldRecordTime.Bool
	}
	if row.CanViewOpponentsSessions.Valid {
		gm.Settings.ViewOpponentsSessions = row.CanViewOpponentsSessions.Bool
	}
	return gm
}

func GetPlayersResultInGame(players []*pg.Player) []*game.Session {
	playerSessions := make([]*game.Session, len(players))
	for i, gp := range players {
		var playerSession *game.Session
		player := &game.Player{
			ID:        gp.ID,
			Name:      gp.Name,
			Points:    gp.Points.Int64,
			Email:     gp.Email,
			Password:  gp.Password,
			IsDeleted: gp.IsDeleted,
		}
		playerSession = game.NewSession(player)
		playerSessions[i] = playerSession
	}
	return playerSessions
}
