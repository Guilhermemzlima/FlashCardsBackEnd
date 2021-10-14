package playlist_usecase

import (
	"encoding/json"
	"fmt"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/errors"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/deck"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/playlist"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/repository/playlist_repository"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase/card_usecase"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase/deck_usecase"
	"github.com/go-playground/validator"
	"github.com/imdario/mergo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IPlaylistUseCase interface {
	Create(userId string, playlist *playlist.Playlist) (result *playlist.Playlist, err error)
	FindByUserIdAndPublic(userId string) (result []*playlist.Playlist, count int64, err error)
	FindById(userId, id string) (result *playlist.Playlist, err error)
	FindByUserId(userId string) (result []*playlist.Playlist, count int64, err error)
	Delete(id, userId string) (result *playlist.Playlist, err error)
	Update(id, userId string, isPartial bool, playlist *playlist.Playlist) (*playlist.Playlist, error)
	AddDeckToPlaylist(id, userId string, deckId string) (*playlist.Playlist, error)
}
type PlaylistUseCase struct {
	validator   *validator.Validate
	repo        playlist_repository.IPlaylistRepository
	deckUseCase deck_usecase.DeckUseCase
	cardUseCase card_usecase.CardUseCase
}

func NewPlaylistUseCase(playlistRepository playlist_repository.IPlaylistRepository, cardUseCase card_usecase.CardUseCase, deckRepo deck_usecase.DeckUseCase, validator *validator.Validate) PlaylistUseCase {
	return PlaylistUseCase{
		validator:   validator,
		repo:        playlistRepository,
		deckUseCase: deckRepo,
		cardUseCase: cardUseCase,
	}
}

func (uc PlaylistUseCase) Create(userId string, playlist *playlist.Playlist) (result *playlist.Playlist, err error) {
	playlist.UserId = userId
	playlist.LastUpdate = time.Now()

	err = uc.validator.Struct(playlist)
	if err != nil {
		playlistBytes, _ := json.Marshal(playlist)
		log.Logger.Errorf("Error to validate input:\n %v;\n error: %v", string(playlistBytes), err.Error())
		return nil, &errors.InvalidPayload{Err: err}
	}

	result, err = uc.repo.Persist(playlist)
	if err != nil {
		log.Logger.Errorw("Playlist creation error", "Error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}

	return result, nil
}

func (uc PlaylistUseCase) FindById(userId, id string) (result *playlist.Playlist, err error) {
	objectID, err := uc.parseToObjectID(id)
	if err != nil {
		return nil, err
	}

	result, err = uc.repo.FindById(userId, &objectID, false)
	if err != nil {
		log.Logger.Errorw("playlist not found", "Error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}
	return result, nil
}

func (uc PlaylistUseCase) FindByUserId(userId string) (result []*playlist.Playlist, count int64, err error) {
	result, err = uc.repo.FindByUserId(userId)
	if err != nil {
		log.Logger.Errorw("playlist not found", "Error", err.Error())
		return nil, 0, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}

	count, err = uc.repo.Count(userId)
	if err != nil {
		log.Logger.Errorw("Count playlists error", "Error", err.Error())
		return nil, 0, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}

	return result, count, nil
}

func (uc PlaylistUseCase) FindByUserIdAndPublic(userId string) (result []*playlist.Playlist, count int64, err error) {
	result, err = uc.repo.FindByUserIdAndPublic(userId)
	if err != nil {
		log.Logger.Errorw("playlist not found", "Error", err.Error())
		return nil, 0, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}

	count, err = uc.repo.Count(userId)
	if err != nil {
		log.Logger.Errorw("Count playlists error", "Error", err.Error())
		return nil, 0, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}
	return result, count, nil
}

func (uc PlaylistUseCase) parseToObjectID(id string) (objID primitive.ObjectID, err error) {
	if id == "" {
		err := errors.WrapWithMessage(errors.ErrInvalidPayload, "id is required")
		log.Logger.Errorw("id is required")
		return objID, err
	}

	objID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Logger.Errorw("invalid Id", "Error", err.Error())
		return objID, errors.WrapWithMessage(errors.ErrInvalidPayload, err.Error())
	}

	return
}

func (uc PlaylistUseCase) Delete(id, userId string) (result *playlist.Playlist, err error) {
	objectID, err := uc.parseToObjectID(id)
	if err != nil {
		return nil, err
	}

	result, err = uc.repo.Delete(userId, &objectID)
	if err != nil {
		log.Logger.Errorw("Remove playlist error", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}
	return
}

func (uc PlaylistUseCase) Update(id, userId string, isPartial bool, playlist *playlist.Playlist) (*playlist.Playlist, error) {
	objectID, err := uc.parseToObjectID(id)
	if err != nil {
		return nil, err
	}

	playlist.UserId = userId
	playlist.LastUpdate = time.Now()

	savedPlaylist, err := uc.repo.FindById(userId, &objectID, true)
	if err != nil {
		log.Logger.Errorw("playlist not found", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}

	if isPartial {
		if err := mergo.Merge(playlist, savedPlaylist); err != nil {
			log.Logger.Errorf("Error to merge entities")
			return nil, &errors.BadRequest{Err: err}
		}
	}

	err = uc.validator.Struct(playlist)
	if err != nil {
		playlistBytes, _ := json.Marshal(playlist)
		log.Logger.Errorf("Error to validate input:\n %v;\n error: %v", string(playlistBytes), err.Error())
		return nil, &errors.InvalidPayload{Err: err}
	}

	result, err := uc.repo.Update(&objectID, userId, playlist)
	if err != nil {
		log.Logger.Errorw("update error", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}

	if result == nil {
		return nil, errors.WrapWithMessage(errors.ErrNotFound, fmt.Sprintf("id %s not found", id))
	}

	return result, nil
}

func (uc PlaylistUseCase) AddDeckToPlaylist(id, userId string, deckId string) (*playlist.Playlist, error) {
	var preview deck.DeckPreview
	savedDeck, err := uc.deckUseCase.FindById(userId, deckId)
	if err != nil {
		log.Logger.Errorw("playlistToSave not found", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}

	if savedDeck != nil {
		preview = deck.DeckPreview{
			Id:       savedDeck.Id.Hex(),
			ImageURL: savedDeck.ImageURL,
			Name:     savedDeck.Name,
			UserId:   savedDeck.UserId,
		}
	}

	err = uc.validator.Struct(preview)
	if err != nil {
		playlistBytes, _ := json.Marshal(preview)
		log.Logger.Errorf("Error to validate input:\n %v;\n error: %v", string(playlistBytes), err.Error())
		return nil, &errors.InvalidPayload{Err: err}
	}

	savedPlaylist, err := uc.FindById(userId, id)
	if err != nil {
		log.Logger.Errorw("playlist not found", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}

	savedPlaylist.Decks = append(savedPlaylist.Decks, preview)

	result, err := uc.Update(id, userId, true, savedPlaylist)
	if err != nil {
		log.Logger.Errorw("update error", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}

	if result == nil {
		return nil, errors.WrapWithMessage(errors.ErrNotFound, fmt.Sprintf("id %s not found", id))
	}

	return result, nil
}
