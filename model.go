// model.go

package todo_api

import (
	"database/sql"
	"errors"
)

type goal struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (g *goal) getGoal(db *sql.DB) error {
	return errors.New("NOT IMPLEMENTED")
}

func (g *goal) updateGoal(db *sql.DB) error {
	return errors.New("NOT IMPLEMENTED")
}

func (g *goal) addGoal(db *sql.DB) error {
	return errors.New("NOT IMPLEMENTED")
}

func (g *goal) deleteGoal(db *sql.DB) error {
	return errors.New("NOT IMPLEMENTED")
}

func getGoals(db *sql.DB, start, count int) ([]goal, error) {
	return nil, errors.New("Not implemented")
}
