package data

import (
	"database/sql"
	"time"

	"github.com/AustinMusiku/Greenlight/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // Use the - directive
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`           // Add the omitempty directive
	Runtime   Runtime   `json:"runtime,omitempty,string"` // Add the omitempty directive // Add the string directive
	Genres    []string  `json:"genres,omitempty"`         // Add the omitempty directive
	Version   int32     `json:"version"`
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type MovieModel struct {
	DB *sql.DB
}

// Placeholder method for inserting a new record in the movies table.
func (m MovieModel) Insert(movie *Movie) error {
	return nil
}

// Placeholder method for fetching a specific record from the movies table.
func (m MovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

// Placeholder method for updating a specific record in the movies table.
func (m MovieModel) Update(movie *Movie) error {
	return nil
}

// Placeholder method for deleting a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
	return nil
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	// Title checks
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	// Year checks
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	// Runtime checks
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	// Genres checks
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
