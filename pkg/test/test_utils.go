package utils

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/playlist"
	"time"
)

const (
	userId             = "Totoro"
	PersistMethodName  = "Persist"
	UpdateMethodName   = "Update"
	DeleteMethodName   = "Delete"
	FindByIdMethodName = "FindById"
	FindByUserIdMethodName = "FindByUserId"
	CountMethodName    = "Count"
)

func BuildPlaylist() *playlist.Playlist {
	return &playlist.Playlist{
		ImageURL:         "https://picsum.photos/200/300",
		Name:             "Playlist de Golang",
		Description:      "Feita para estudar Golang",
		IsPrivate:        false,
		StudySuggestions: nil,
		Decks:            nil,
		UserId:           userId,
		LastUpdate:       time.Time{},
	}
}
