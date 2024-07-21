package controllers

import (
	"time"

	"github.com/fabregas201307/fiber-go-template/pkg/utils"
	"github.com/fabregas201307/fiber-go-template/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Getbonds func gets all exists bonds.
// @Description Get all exists bonds.
// @Summary get all exists bonds
// @Tags bonds
// @Accept json
// @Produce json
// @Success 200 {array} models.bond
// @Router /v1/bonds [get]
func Getbonds(c *fiber.Ctx) error {
	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get all bonds.
	bonds, err := db.Getbonds()
	if err != nil {
		// Return, if bonds not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "bonds were not found",
			"count": 0,
			"bonds": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": len(bonds),
		"bonds": bonds,
	})
}

// Getbond func gets bond by given ID or 404 error.
// @Description Get bond by given ID.
// @Summary get bond by given ID
// @Tags bond
// @Accept json
// @Produce json
// @Param id path string true "bond ID"
// @Success 200 {object} models.bond
// @Router /v1/bond/{id} [get]
func Getbond(c *fiber.Ctx) error {
	// Catch bond ID from URL.
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get bond by ID.
	bond, err := db.Getbond(id)
	if err != nil {
		// Return, if bond not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "bond with the given ID is not found",
			"bond":  nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"bond":  bond,
	})
}

// Createbond func for creates a new bond.
// @Description Create a new bond.
// @Summary create a new bond
// @Tags bond
// @Accept json
// @Produce json
// @Param title body string true "Title"
// @Param author body string true "Author"
// @Param user_id body string true "User ID"
// @Param bond_attrs body models.bondAttrs true "bond attributes"
// @Success 200 {object} models.bond
// @Security ApiKeyAuth
// @Router /v1/bond [post]
func Createbond(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current bond.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Set credential `bond:create` from JWT data of current bond.
	credential := claims.Credentials[repository.bondCreateCredential]

	// Only user with `bond:create` credential can create a new bond.
	if !credential {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, check credentials of your token",
		})
	}

	// Create new bond struct
	bond := &models.bond{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(bond); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a bond model.
	validate := utils.NewValidator()

	// Set initialized default data for bond:
	bond.ID = uuid.New()
	bond.CreatedAt = time.Now()
	bond.UserID = claims.UserID
	bond.bondStatus = 1 // 0 == draft, 1 == active

	// Validate bond fields.
	if err := validate.Struct(bond); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create bond by given model.
	if err := db.Createbond(bond); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"bond":  bond,
	})
}

// Updatebond func for updates bond by given ID.
// @Description Update bond.
// @Summary update bond
// @Tags bond
// @Accept json
// @Produce json
// @Param id body string true "bond ID"
// @Param title body string true "Title"
// @Param author body string true "Author"
// @Param user_id body string true "User ID"
// @Param bond_status body integer true "bond status"
// @Param bond_attrs body models.bondAttrs true "bond attributes"
// @Success 202 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/bond [put]
func Updatebond(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current bond.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Set credential `bond:update` from JWT data of current bond.
	credential := claims.Credentials[repository.bondUpdateCredential]

	// Only bond creator with `bond:update` credential can update his bond.
	if !credential {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, check credentials of your token",
		})
	}

	// Create new bond struct
	bond := &models.bond{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(bond); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if bond with given ID is exists.
	foundedbond, err := db.Getbond(bond.ID)
	if err != nil {
		// Return status 404 and bond not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "bond with this ID not found",
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his bond.
	if foundedbond.UserID == userID {
		// Set initialized default data for bond:
		bond.UpdatedAt = time.Now()

		// Create a new validator for a bond model.
		validate := utils.NewValidator()

		// Validate bond fields.
		if err := validate.Struct(bond); err != nil {
			// Return, if some fields are not valid.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   utils.ValidatorErrors(err),
			})
		}

		// Update bond by given ID.
		if err := db.Updatebond(foundedbond.ID, bond); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Return status 201.
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"error": false,
			"msg":   nil,
		})
	} else {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, only the creator can delete his bond",
		})
	}
}

// Deletebond func for deletes bond by given ID.
// @Description Delete bond by given ID.
// @Summary delete bond by given ID
// @Tags bond
// @Accept json
// @Produce json
// @Param id body string true "bond ID"
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/bond [delete]
func Deletebond(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current bond.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Set credential `bond:delete` from JWT data of current bond.
	credential := claims.Credentials[repository.bondDeleteCredential]

	// Only bond creator with `bond:delete` credential can delete his bond.
	if !credential {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, check credentials of your token",
		})
	}

	// Create new bond struct
	bond := &models.bond{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(bond); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a bond model.
	validate := utils.NewValidator()

	// Validate bond fields.
	if err := validate.StructPartial(bond, "id"); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if bond with given ID is exists.
	foundedbond, err := db.Getbond(bond.ID)
	if err != nil {
		// Return status 404 and bond not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "bond with this ID not found",
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his bond.
	if foundedbond.UserID == userID {
		// Delete bond by given ID.
		if err := db.Deletebond(foundedbond.ID); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, only the creator can delete his bond",
		})
	}
}
