package controller

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ara-thesis/monarch-project-be/helper"
	"github.com/ara-thesis/monarch-project-be/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentController struct {
	Tbname      string
	Tbname_cart string
}

func (pc *PaymentController) GetPurchaseConfirm(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	// permission check
	if userData.UserRole != "ADMIN" {
		return resp.Forbidden(c, "Access Forbidden")
	}
	row, rowErr := strconv.Atoi(c.Query("row", "10"))
	if rowErr != nil {
		row = 10
	}
	if row > 100 {
		row = 100
	}
	page, pageErr := strconv.Atoi(c.Query("page", "1"))
	if pageErr != nil {
		page = 1
	}

	qyStr := fmt.Sprintf(`SELECT * FROM %s LIMIT $1 OFFSET $2`, pc.Tbname)
	resQy, resErr := db.Query(qyStr, row, (page-1)*row)

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")

}

func (pc *PaymentController) GetPurchaseConfirmById(c *fiber.Ctx) error {

	userData := c.Locals("user").(*helper.ClaimsData)

	// permission check
	if userData.UserRole != "ADMIN" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	qyStr := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, pc.Tbname)
	resQy, resErr := db.Query(qyStr, c.Query("id"))

	if resErr != nil {
		return resp.ServerError(c, resErr.Error())
	}

	return resp.Success(c, resQy, "Success Fetching Data")
}

func (pc *PaymentController) AcceptPurchaseConfirmation(c *fiber.Ctx) error {
	return nil
}

func (pc *PaymentController) DenyPurchaseConfirmation(c *fiber.Ctx) error {
	return nil
}

func (pc *PaymentController) PayCart(c *fiber.Ctx) error {

	model := new(model.PaymentModel)
	userData := c.Locals("user").(*helper.ClaimsData)
	uuid := uuid.New()

	// permission check
	if userData.UserRole != "Tourist" {
		return resp.Forbidden(c, "Access Forbidden")
	}

	// check file availability
	qyStr := fmt.Sprintf("SELECT * FROM %s WHERE created_by = $1", pc.Tbname_cart)
	checkData, checkErr := db.Query(qyStr, userData.UserId)
	if checkErr != nil {
		return resp.ServerError(c, checkErr.Error())
	}
	if checkData[0] == nil {
		return resp.NotFound(c, "Data Not Found")
	}

	// file process
	fileForm, fileErr := c.FormFile("image")
	if fileErr == nil {
		fileName := fmt.Sprintf("%s-%s", checkData[0].(map[string]interface{})["id"], fileForm.Filename)
		c.SaveFile(fileForm, fmt.Sprintf("public/payment/%s", fileName))
		model.Image = fmt.Sprintf("/api/public/payment/%s", fileName)
	}

	// insert data process
	cmdStr1 := fmt.Sprintf(`INSERT INTO %s(id, total_price, image, created_at, created_by, updated_at, updated_by)
	VALUES($1, $2, $3, $4, $5, $6, $7)`, pc.Tbname)
	cmdErr1 := db.Command(cmdStr1, uuid, model.Total_price, model.Image, time.Now(), userData.UserId, time.Now(), userData.UserId)
	if cmdErr1 != nil {
		resp.ServerError(c, "Error Adding Data: "+cmdErr1.Error())
	}

	// update data process
	cmdStr2 := fmt.Sprintf("UPDATE %s SET wait=$1, updated_by=$2, updated_at=$3 WHERE created_by = $4", pc.Tbname_cart)
	cmdErr2 := db.Command(cmdStr2, true, time.Now(), userData.UserId, userData.UserId)
	if cmdErr2 != nil {
		resp.ServerError(c, "Error Adding Data: "+cmdErr2.Error())
	}
	return resp.Created(c, "Success Adding Data")
}
