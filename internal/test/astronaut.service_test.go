package test

import (
	"context"
	"errors"
	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/LaQuannT/astronaut-api/internal/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddAstronaut(t *testing.T) {
	err := clearTables(dbConn)
	if err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}
	ctx := context.TODO()

	t.Run("Adds a new Astronaut", func(t *testing.T) {
		a := &model.Astronaut{
			FirstName:  "john",
			LastName:   "smith",
			Gender:     "M",
			BirthDate:  "1999-01-01",
			BirthPlace: "new york,ny",
		}
		a, err := service.AddAstronaut(ctx, a, astroRepo)
		if err != nil {
			t.Fatalf("Unexpected error adding Astronaut: %v", err)
		}
		assert.NotEqual(t, 0, a.ID)
		assert.Equal(t, 1, a.ID)
		assert.Equal(t, "john", a.FirstName)
		assert.Equal(t, "smith", a.LastName)
		assert.Equal(t, "M", a.Gender)
		assert.Equal(t, "1999-01-01", a.BirthDate)
		assert.Equal(t, "new york,ny", a.BirthPlace)
	})

	t.Run("throws an error for invalid astronaut data", func(t *testing.T) {
		a := &model.Astronaut{}
		astronaut, err := service.AddAstronaut(ctx, a, astroRepo)
		if err == nil {
			t.Fatal("Expected error for invalid astronaut data")
		}
		assert.Nil(t, astronaut)
	})
}

func TestGetAstronaut(t *testing.T) {
	err := clearTables(dbConn)
	if err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}
	ctx := context.TODO()

	a := &model.Astronaut{
		FirstName:  "jane",
		LastName:   "doe",
		Gender:     "F",
		BirthDate:  "1999-01-01",
		BirthPlace: "london,uk",
	}
	a, err = service.AddAstronaut(ctx, a, astroRepo)
	if err != nil {
		t.Fatalf("Unexpected error adding Astronaut: %v", err)
	}

	t.Run("returns nil if not found", func(t *testing.T) {
		id := 3

		astronaut, err := service.GetAstronaut(ctx, astroRepo, id)
		if err != nil {
			var apiErr *model.APIError
			switch {
			case errors.As(err, &apiErr):
				if apiErr.Message != "not found" {
					t.Errorf("Unexpected error message: %v", apiErr.Message)
				}
			default:
				t.Errorf("Unexpected error: %v", err)
			}
		}
		assert.Nil(t, astronaut)
	})

	t.Run("returns an astronaut", func(t *testing.T) {
		id := 1
		astronaut, err := service.GetAstronaut(ctx, astroRepo, id)
		if err != nil {
			t.Errorf("Unexpected error getting Astronaut: %v", err)
		}

		assert.Equal(t, astronaut.ID, a.ID)
	})
}

func TestGetAstronauts(t *testing.T) {
	err := clearTables(dbConn)
	if err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}
	ctx := context.TODO()

	t.Run("returns nil if no astronauts found", func(t *testing.T) {
		astronauts, err := service.GetAstronauts(ctx, astroRepo)
		if err != nil {
			t.Errorf("Unexpected error getting Astronauts: %v", err)
		}
		assert.Nil(t, astronauts)
	})

	t.Run("returns an list of astronauts", func(t *testing.T) {
		john := &model.Astronaut{
			FirstName:  "john",
			LastName:   "doe",
			Gender:     "M",
			BirthDate:  "1999-01-01",
			BirthPlace: "london,uk",
		}
		jane := &model.Astronaut{
			FirstName:  "jane",
			LastName:   "doe",
			Gender:     "F",
			BirthDate:  "1999-01-01",
			BirthPlace: "london,uk",
		}

		astronauts := []*model.Astronaut{john, jane}
		for _, a := range astronauts {
			_, err := service.AddAstronaut(ctx, a, astroRepo)
			if err != nil {
				t.Fatalf("Unexpected error adding Astronaut: %v", err)
			}
		}

		astronauts, err := service.GetAstronauts(ctx, astroRepo)
		if err != nil {
			t.Errorf("Unexpected error getting Astronauts: %v", err)
		}
		assert.Equal(t, 2, len(astronauts))
	})
}

func TestSearchAstronautByName(t *testing.T) {
	ctx := context.TODO()
	t.Run("returns nil if no astronauts found", func(t *testing.T) {
		target := "smith"
		astronauts, err := service.SearchAstronautByName(ctx, astroRepo, target)
		if err != nil {
			t.Errorf("Unexpected error getting Astronauts: %v", err)
		}
		assert.Nil(t, astronauts)
	})

	t.Run("returns a list of astronauts with last name matching search target", func(t *testing.T) {
		target := "doe"
		astronauts, err := service.SearchAstronautByName(ctx, astroRepo, target)
		if err != nil {
			t.Errorf("Unexpected error getting Astronauts: %v", err)
		}
		assert.Equal(t, 2, len(astronauts))
	})

	t.Run("returns an astronaut with first and last name matching search target", func(t *testing.T) {
		target := "jane doe"
		astronauts, err := service.SearchAstronautByName(ctx, astroRepo, target)
		if err != nil {
			t.Errorf("Unexpected error getting Astronauts: %v", err)
		}
		assert.Equal(t, 1, len(astronauts))
		assert.Equal(t, "jane", astronauts[0].FirstName)
		assert.Equal(t, "doe", astronauts[0].LastName)
	})
}

func TestUpdateAstronaut(t *testing.T) {
	err := clearTables(dbConn)
	if err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}
	ctx := context.TODO()

	a := &model.Astronaut{
		FirstName:  "james",
		LastName:   "smith",
		Gender:     "M",
		BirthDate:  "1999-01-01",
		BirthPlace: "london,uk",
	}

	a, err = service.AddAstronaut(ctx, a, astroRepo)
	if err != nil {
		t.Fatalf("Unexpected error adding Astronaut: %v", err)
	}

	t.Run("returns an error for invalid update astronaut input", func(t *testing.T) {
		astronaut := &model.Astronaut{
			ID:        1,
			FirstName: "",
			LastName:  a.LastName,
			Gender:    "",
			BirthDate: a.BirthDate,
		}

		if err := service.UpdateAstronaut(ctx, astronaut, astroRepo); err == nil {
			t.Errorf("Expected error for invalid update astronaut")
		}
	})

	t.Run("returns an updated astronaut", func(t *testing.T) {
		astronaut := &model.Astronaut{
			ID:         1,
			FirstName:  "cindy",
			LastName:   a.LastName,
			Gender:     "F",
			BirthDate:  a.BirthDate,
			BirthPlace: "salford,uk",
		}

		if err := service.UpdateAstronaut(ctx, astronaut, astroRepo); err != nil {
			t.Errorf("Unexpected error updating Astronaut: %v", err)
		}
		a, err = service.GetAstronaut(ctx, astroRepo, astronaut.ID)
		if err != nil {
			t.Errorf("Unexpected error getting updated Astronaut: %v", err)
		}
		assert.Equal(t, astronaut.ID, a.ID)
		assert.Equal(t, astronaut.FirstName, a.FirstName)
		assert.Equal(t, astronaut.Gender, a.Gender)
		assert.Equal(t, astronaut.BirthPlace, a.BirthPlace)
	})
}

func TestDeleteAstronaut(t *testing.T) {
	ctx := context.TODO()

	t.Run("returns an error for unknown astronaut ID", func(t *testing.T) {
		id := 7
		if err := service.DeleteAstronaut(ctx, astroRepo, id); err == nil {
			t.Errorf("Expected error for unknown astronaut ID")
		}
	})

	t.Run("deletes an astronaut", func(t *testing.T) {
		id := 1
		if err := service.DeleteAstronaut(ctx, astroRepo, id); err != nil {
			t.Errorf("Unexpected error deleting Astronaut: %v", err)
		}
	})
}
