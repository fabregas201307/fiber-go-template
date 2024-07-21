package queries

import (
	"github.com/fabregas201307/fiber-go-template/app/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// BondQueries struct for queries from Bond model.
type BondQueries struct {
	*sqlx.DB
}

// GetBonds method for getting all bonds.
func (q *BondQueries) GetBonds() ([]models.Bond, error) {
	// Define bonds variable.
	bonds := []models.Bond{}

	// Define query string.
	query := `SELECT * FROM bonds`

	// Send query to database.
	err := q.Select(&bonds, query)
	if err != nil {
		// Return empty object and error.
		return bonds, err
	}

	// Return query result.
	return bonds, nil
}

// GetBondsByAuthor method for getting all bonds by given author.
func (q *BondQueries) GetBondsByAuthor(author string) ([]models.Bond, error) {
	// Define bonds variable.
	bonds := []models.Bond{}

	// Define query string.
	query := `SELECT * FROM bonds WHERE author = $1`

	// Send query to database.
	err := q.Get(&bonds, query, author)
	if err != nil {
		// Return empty object and error.
		return bonds, err
	}

	// Return query result.
	return bonds, nil
}

// GetBond method for getting one bond by given ID.
func (q *BondQueries) GetBond(id uuid.UUID) (models.Bond, error) {
	// Define bond variable.
	bond := models.Bond{}

	// Define query string.
	query := `SELECT * FROM bonds WHERE id = $1`

	// Send query to database.
	err := q.Get(&bond, query, id)
	if err != nil {
		// Return empty object and error.
		return bond, err
	}

	// Return query result.
	return bond, nil
}

// CreateBond method for creating bond by given Bond object.
func (q *BondQueries) CreateBond(b *models.Bond) error {
	// Define query string.
	query := `INSERT INTO bonds VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	// Send query to database.
	_, err := q.Exec(query, b.ID, b.CreatedAt, b.UpdatedAt, b.UserID, b.Title, b.Author, b.BondStatus, b.BondAttrs)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateBond method for updating bond by given Bond object.
func (q *BondQueries) UpdateBond(id uuid.UUID, b *models.Bond) error {
	// Define query string.
	query := `UPDATE bonds SET updated_at = $2, title = $3, author = $4, bond_status = $5, bond_attrs = $6 WHERE id = $1`

	// Send query to database.
	_, err := q.Exec(query, id, b.UpdatedAt, b.Title, b.Author, b.BondStatus, b.BondAttrs)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteBond method for delete bond by given ID.
func (q *BondQueries) DeleteBond(id uuid.UUID) error {
	// Define query string.
	query := `DELETE FROM bonds WHERE id = $1`

	// Send query to database.
	_, err := q.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
