package controller

import (
	"fmt"
	"time"

	"github.com/ara-thesis/monarch-project-be/helper"
	"github.com/ara-thesis/monarch-project-be/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReviewHandler struct{}

//////////////////////
// fetch all comment
//////////////////////
func (pinf *ReviewHandler) GetComment(c *fiber.Ctx) error {

	placeId := c.Query("placeid")
	if placeId == "" {
		return resp.BadRequest(c, "Need more query (placeid is needed)")
	}

	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE place_id = $1", tbname["placeinfo"])
	resQy, resErr := db.Query(qyStr, placeId)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

////////////////////
// add new comment
////////////////////
func (pinf *ReviewHandler) AddComment(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)
	model := new(model.CommentModel)
	uuid := uuid.New()

	// permission check
	if userData.UserRole != roleId["T"] {
		return resp.Forbidden(c, "Access Forbidden")
	}

	if reqErr := c.BodyParser(model); reqErr != nil {
		return resp.ServerError(c, reqErr.Error())
	}

	// db process
	cmdMainStr := fmt.Sprintf(`
	INSERT INTO %s(id, place_id, comment, created_at, created_by, updated_at, updated_by) VALUES($1, $2, $3, $4, $5, $6)`, tbname["review"])
	resMainErr := db.Command(
		cmdMainStr, uuid, model.Place_Id, model.Comment, time.Now(), userData.UserId, time.Now(), userData.UserId,
	)
	if resMainErr != nil {
		return resp.ServerError(c, "Error Adding Data: "+resMainErr.Error())
	}

	return resp.Created(c, "Success Adding Data")
}

///////////////////
// delete comment
///////////////////
func (pinf *ReviewHandler) DeleteCommentAdmin(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	// permission check
	if userData.UserRole != "PLACE MANAGER" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// check data availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tbname["review"])
	checkData, checkErr := db.Query(qyStr, c.Params("id"))
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// check for permitted user
	if userData.Id != checkData[0].(map[string]interface{})["created_by"] {
		return resp.Forbidden(c, "Delete Forbidden")
	}

	// delete process
	cmdStr := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tbname["news"])
	resErr := db.Command(cmdStr, c.Params("id"))
	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, nil, "Success Delete Data")

}
