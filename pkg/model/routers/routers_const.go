package routers

const (
	BasePath       = "/flashcards"
	ApiPath        = "/api/v1"
	PlaylistPath   = ApiPath + "/playlists"
	PlaylistPathId = PlaylistPath + "/{id}"
	PlaylistPathAll = PlaylistPath + "/all"
	DeckPath       = ApiPath + "/decks"
	DeckPathId     = DeckPath + "/{id}"
	CardPath       = ApiPath + "/cards"
	CardPathId     = CardPath + "/{id}"
)
