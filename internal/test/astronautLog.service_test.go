package test

import (
	"context"
	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/LaQuannT/astronaut-api/internal/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddAstronautLog(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Errorf("Error clearing tables: %v", err)
	}

	ctx := context.TODO()

	a := &model.Astronaut{
		FirstName:  "john",
		LastName:   "doe",
		Gender:     "M",
		BirthDate:  "2022-01-01",
		BirthPlace: "manchester,uk",
	}
	a, err := service.AddAstronaut(ctx, a, astroRepo)
	if err != nil {
		t.Fatalf("Unexpected error adding Astronaut: %v", err)
	}

	t.Run("returns error and nil for invalid astronaut log input", func(t *testing.T) {
		al := &model.AstronautLog{}

		al, err := service.AddAstronautLog(ctx, astroLogRepo, al)
		if err == nil {
			t.Errorf("Expected error adding invalid astronaut log input")
		}
		assert.Nil(t, al)
	})

	t.Run("returns error adding invalid astronaut log for unknown astronaut id", func(t *testing.T) {
		al := &model.AstronautLog{
			AstronautID: 30,
			Status:      model.Active,
		}
		al, err := service.AddAstronautLog(ctx, astroLogRepo, al)
		if err == nil {
			t.Errorf("Expected error adding astronaut log")
		}
		assert.Nil(t, al)
	})

	t.Run("adds a new astronaut log", func(t *testing.T) {
		al := &model.AstronautLog{
			AstronautID: a.ID,
			Status:      model.Retired,
		}

		al, err := service.AddAstronautLog(ctx, astroLogRepo, al)
		if err != nil {
			t.Errorf("Unexpected error adding astronaut log: %v", err)
		}
		assert.Equal(t, a.ID, al.AstronautID)
	})
}

func TestGetAstronautLog(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Errorf("Error clearing tables: %v", err)
	}
	ctx := context.TODO()

	a := &model.Astronaut{
		FirstName:  "john",
		LastName:   "doe",
		Gender:     "M",
		BirthDate:  "2022-01-01",
		BirthPlace: "manchester,uk",
	}

	a, err := service.AddAstronaut(ctx, a, astroRepo)
	if err != nil {
		t.Fatalf("Unexpected error adding Astronaut: %v", err)
	}

	log := &model.AstronautLog{
		AstronautID: a.ID,
		Status:      model.Retired,
	}
	log, err = service.AddAstronautLog(ctx, astroLogRepo, log)
	if err != nil {
		t.Fatalf("Unexpected error adding AstronautLog: %v", err)
	}

	t.Run("returns error and nil for unknown astronaut ID", func(t *testing.T) {
		al, err := service.GetAstronautLog(ctx, astroLogRepo, 123)
		if err == nil {
			t.Errorf("Expected error getting astronaut log with unknown astronaut ID")
		}
		assert.Nil(t, al)
	})

	t.Run("returns an astronaut log", func(t *testing.T) {
		al, err := service.GetAstronautLog(ctx, astroLogRepo, a.ID)
		if err != nil {
			t.Errorf("Unexpected error getting astronaut log: %v", err)
		}
		assert.Equal(t, a.ID, al.AstronautID)
		assert.Equal(t, model.Retired, al.Status)
		assert.Equal(t, "", al.DeathDate)
	})
}

func TestGetAstronautLogs(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Errorf("Error clearing tables: %v", err)
	}
	ctx := context.TODO()

	jane := &model.Astronaut{
		FirstName:  "jane",
		LastName:   "doe",
		Gender:     "F",
		BirthDate:  "2022-01-01",
		BirthPlace: "new york, ny",
	}

	john := &model.Astronaut{
		FirstName:  "john",
		LastName:   "doe",
		Gender:     "M",
		BirthDate:  "2022-01-01",
		BirthPlace: "manchester,uk",
	}

	as := []*model.Astronaut{jane, john}
	t.Run("returns nil if no astronaut logs are found", func(t *testing.T) {
		als, err := service.GetAstronautLogs(ctx, astroLogRepo)
		if err != nil {
			t.Errorf("Unexpected error getting astronaut logs: %v", err)
		}
		assert.Nil(t, als)
	})

	t.Run("returns a list on astronaut logs", func(t *testing.T) {

		for _, a := range as {
			a, err := service.AddAstronaut(ctx, a, astroRepo)
			log := &model.AstronautLog{
				AstronautID: a.ID,
				Status:      model.Management,
			}

			if err != nil {
				t.Errorf("Unexpected error adding Astronaut: %v", err)
			}
			log, err = service.AddAstronautLog(ctx, astroLogRepo, log)
		}

		als, err := service.GetAstronautLogs(ctx, astroLogRepo)
		if err != nil {
			t.Errorf("Unexpected error getting astronaut logs: %v", err)
		}
		assert.Len(t, als, 2)
	})
}
