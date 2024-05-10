package handler

import (
	connectteam "ConnectTeam/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type createGameInput struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
}

func (h *Handler) createGame(c *gin.Context) {
	var input createGameInput
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if access != string(connectteam.UserAccess) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.services.GetUserActivePlan(id)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, "user has no plan")
		return
	}

	game, err := h.services.CreateGame(id, input.StartDate, input.Name)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.CreateParticipant(game.CreatorId, game.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":              game.Id,
		"name":            game.Name,
		"start_date":      game.StartDate,
		"status":          game.Status,
		"invitation_code": game.InvitationCode,
	})
}

//type getUserGamesInput struct {
//	Limit  int `json:"limit" binding:"required"`
//	Offset int `json:"offset" binding:"required"`
//}

type getUserGamesResponse struct {
	Data []connectteam.Game `json:"data"`
}

func (h *Handler) getCreatedGames(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if access != string(connectteam.UserAccess) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//_, err = h.services.GetUserPlan(id)
	//if err != nil {
	//	newErrorResponse(c, http.StatusForbidden, "user has no plan")
	//	return
	//}

	games, err := h.services.Game.GetCreatedGames(page, id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getUserGamesResponse{
		Data: games,
	})
}

type getResultsResponse struct {
	Data []connectteam.UserResult `json:"data"`
}

func (h *Handler) getResults(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if access != string(connectteam.UserAccess) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	gameId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	results, err := h.services.GetResults(gameId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getResultsResponse{
		Data: results,
	})
}

func (h *Handler) inviteMemberToGame(c *gin.Context) {
	creatorId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input inviteUserRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	gameId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.InviteUserToGame(gameId, input.UserId, creatorId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)

}

func (h *Handler) deleteGameFromGameList(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if access != string(connectteam.UserAccess) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	gameId, err := uuid.Parse(c.Param("id"))

	game, err := h.services.Game.GetGame(gameId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if id == game.CreatorId {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	err = h.services.Game.DeleteGameFromGameList(gameId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

type changeGameStartDateInput struct {
	StartDate string `json:"start_date" binding:"required"`
}

func (h *Handler) changeGameStartDate(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	gameId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	game, err := h.services.Game.GetGame(gameId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if game.CreatorId != id {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
	}

	var input changeGameStartDateInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.ChangeStartDate(gameId, input.StartDate)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

type changeGameNameInput struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) changeGameName(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	gameId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	game, err := h.services.Game.GetGame(gameId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if game.CreatorId != id {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
	}

	var input changeGameNameInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.ChangeGameName(gameId, input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getGame(c *gin.Context) {

	gameId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	game, err := h.services.Game.GetGame(gameId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":              game.Id,
		"name":            game.Name,
		"creator_id":      game.CreatorId,
		"start_date":      game.StartDate,
		"status":          game.Status,
		"invitation_code": game.InvitationCode,
	})
}

func (h *Handler) addUserAsParticipant(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	invitationCode := c.Param("code")

	var game connectteam.Game
	game, err = h.services.Game.GetGameWithInvitationCode(invitationCode)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if game.Id == uuid.Nil {
		newErrorResponse(c, http.StatusNotFound, "incorrect invitation code")
		return
	}

	err = h.services.Game.CreateParticipant(id, game.Id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) validateGameInvitationCode(c *gin.Context) {
	code := c.Param("code")
	game, err := h.services.GetGameWithInvitationCode(code)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if game.Id == uuid.Nil || len(code) == 0 {
		newErrorResponse(c, http.StatusNotFound, "incorrect invitation code")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":              game.Id,
		"name":            game.Name,
		"creator_id":      game.CreatorId,
		"start_date":      game.StartDate,
		"status":          game.Status,
		"invitation_code": game.InvitationCode,
	})
}

func (h *Handler) getTagsResults(c *gin.Context) {
	gameId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	userId, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	tags, err := h.services.GetTagsUsers(userId, gameId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getTagsResponse{
		Data: tags,
	})

}

type getGameMembersResponse struct {
	Members []connectteam.UserPublic `json:"members"`
}

func (h *Handler) cancelGame(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	gameId, err := uuid.Parse(c.Param("id"))

	err = h.services.CancelGame(gameId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getGameMembers(c *gin.Context) {
	gameId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	members, err := h.services.Game.GetGameParticipants(gameId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getGameMembersResponse{
		Members: members,
	})
}

//func (h *Handler) addTopicToGame(c *gin.Context) {
//	id, err := getUserId(c)
//	if err != nil {
//		newErrorResponse(c, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	gameId, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		newErrorResponse(c, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	game, err := h.services.GetGame(gameId)
//	if err != nil {
//		newErrorResponse(c, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	if game.CreatorId != id {
//		newErrorResponse(c, http.StatusForbidden, err.Error())
//		return
//	}
//
//}

func (h *Handler) getGames(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//_, err = h.services.GetUserPlan(id)
	//if err != nil {
	//	newErrorResponse(c, http.StatusForbidden, "user has no plan")
	//	return
	//}
	games, err := h.services.Game.GetGames(page, id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getUserGamesResponse{
		Data: games,
	})

}
