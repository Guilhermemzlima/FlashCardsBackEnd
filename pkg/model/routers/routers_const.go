package routers

const (
	BasePath           = "/flashcards"
	ApiPath            = "/api/v1"
	PlaylistPath       = ApiPath + "/playlists"
	PlaylistPathId     = PlaylistPath + "/{id}"
	PlaylistPathAll    = PlaylistPath + "/all"
	PlaylistReviewPath = PlaylistPathId + "/review"

	PlaylistPathAdd   = PlaylistPathId + "/deck/add"
	DeckPath          = ApiPath + "/decks"
	DeckPathId        = DeckPath + "/{id}"
	DeckReviewPath    = DeckPathId + "/review"
	DeckPathAll       = DeckPath + "/all"
	CardPath          = ApiPath + "/cards"
	CardPathId        = CardPath + "/{id}"
	CardDeckPathId    = CardPath + "/decks" + "/{id}"
	ReviewPath        = ApiPath + "/review"
	ReviewPathId      = ReviewPath + "/{id}"
	ReviewPathIdWrong = ReviewPath + "/{id}/wrong"
	ReviewPathIdRight = ReviewPath + "/{id}/right"
)
