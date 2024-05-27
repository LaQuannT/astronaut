package test

import (
	"context"
	"testing"

	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/LaQuannT/astronaut-api/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestAddMilitaryLog(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}

	ctx := context.TODO()

	a := &model.Astronaut{
		FirstName:  "ryan",
		LastName:   "thomas",
		Gender:     "M",
		BirthDate:  "1999-01-01",
		BirthPlace: "usa",
	}

	a, err := service.AddAstronaut(ctx, a, astroRepo)
	if err != nil {
		t.Fatalf("Unexpected error adding astronuat: %v", err)
	}

	t.Run("returns error when adding a military log for unknown astronaut", func(t *testing.T) {
		ml := &model.MilitaryLog{
			AstronautID: 78,
			Branch:      "air force",
			Rank:        "major",
			Retired:     false,
		}

		ml, err := service.AddMilitaryLog(ctx, militaryRepo, ml)
		if err == nil {
			t.Error("Expected an error adding Military Log for unknown astronaut")
		}
		assert.Nil(t, ml)
	})

	t.Run("returns an error for invalid military log input", func(t *testing.T) {
		ml := &model.MilitaryLog{
			AstronautID: a.ID,
		}

		ml, err := service.AddMilitaryLog(ctx, militaryRepo, ml)
		if err == nil {
			t.Error("Expected an error adding Military Log for unknown astronaut")
		}
		assert.Nil(t, ml)
	})

	t.Run("adds a new military log", func(t *testing.T) {
		ml := &model.MilitaryLog{
			AstronautID: a.ID,
			Branch:      "air force",
			Rank:        "major",
			Retired:     false,
		}

		ml, err := service.AddMilitaryLog(ctx, militaryRepo, ml)
		if err != nil {
			t.Errorf("Unexpected error adding astronaut military log: %v", err)
		}

		assert.Equal(t, a.ID, ml.AstronautID)
	})
}

func TestGetMilitaryLog(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}
	ctx := context.TODO()

	a := &model.Astronaut{
		FirstName:  "john",
		LastName:   "doe",
		Gender:     "M",
		BirthDate:  "1999-01-01",
		BirthPlace: "usa",
	}
	a, err := service.AddAstronaut(ctx, a, astroRepo)
	if err != nil {
		t.Fatalf("Unexpected error adding astronaut: %v", err)
	}
	ml := &model.MilitaryLog{
		AstronautID: a.ID,
		Branch:      "air force",
		Rank:        "major",
		Retired:     false,
	}
	ml, err = service.AddMilitaryLog(ctx, militaryRepo, ml)
	if err != nil {
		t.Errorf("Unexpected error adding military log: %v", err)
	}

	t.Run("returns error when getting military log for unknown astronaut", func(t *testing.T) {
		astronautID := 67

		log, err := service.GetMilitaryLog(ctx, militaryRepo, astronautID)
		if err == nil {
			t.Error("Expected an error getting military log for unknown astronaut")
		}
		assert.Nil(t, log)
	})

	t.Run("returns an military log", func(t *testing.T) {
		log, err := service.GetMilitaryLog(ctx, militaryRepo, a.ID)
		if err != nil {
			t.Errorf("Unexpected error getting military log: %v", err)
		}
		assert.Equal(t, ml.AstronautID, log.AstronautID)
		assert.Equal(t, ml.Branch, log.Branch)
		assert.Equal(t, ml.Rank, log.Rank)
		assert.Equal(t, ml.Retired, log.Retired)
	})
}

func TestGetMilitaryLogs(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}
	ctx := context.TODO()

	jane := &model.Astronaut{
		FirstName:  "jane",
		LastName:   "doe",
		Gender:     "F",
		BirthDate:  "1999-01-01",
		BirthPlace: "usa",
	}
	john := &model.Astronaut{
		FirstName:  "john",
		LastName:   "doe",
		Gender:     "M",
		BirthDate:  "1999-01-01",
		BirthPlace: "usa",
	}

	janeMl := &model.MilitaryLog{
		Branch:  "navy",
		Rank:    "sergeant",
		Retired: false,
	}

	johnMl := &model.MilitaryLog{
		Branch:  "army",
		Rank:    "staff sergeant",
		Retired: false,
	}

	t.Run("returns nil and an error if no astronaut military logs found", func(t *testing.T) {
		mls, err := service.GetMilitaryLogs(ctx, militaryRepo)
		if err != nil {
			t.Errorf("Unexpected error getting military logs: %v", err)
		}
		assert.Nil(t, mls)
	})

	t.Run("returns a list of military logs", func(t *testing.T) {
		as := []*model.Astronaut{jane, john}
		mls := []*model.MilitaryLog{janeMl, johnMl}

		for i, a := range as {
			a, err := service.AddAstronaut(ctx, a, astroRepo)
			if err != nil {
				t.Fatalf("Unexpected error adding military log: %v", err)
			}
			mls[i].AstronautID = a.ID
			_, err = service.AddMilitaryLog(ctx, militaryRepo, mls[i])
			if err != nil {
				t.Fatalf("Unexpected error adding military log: %v", err)
			}
		}

		logs, err := service.GetMilitaryLogs(ctx, militaryRepo)
		if err != nil {
			t.Errorf("Unexpected error getting military logs: %v", err)
		}
		assert.Len(t, logs, 2)
	})
}

func TestUpdateMilitaryLog(t *testing.T) {
	err := clearTables(dbConn)
	if err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}

	ctx := context.TODO()

	a := &model.Astronaut{
		FirstName:  "johnny",
		LastName:   "bravo",
		Gender:     "M",
		BirthDate:  "1999-01-01",
		BirthPlace: "usa",
	}

	a, err = service.AddAstronaut(ctx, a, astroRepo)
	if err != nil {
		t.Fatalf("Unexpected error adding astronaut: %v", err)
	}

	log := &model.MilitaryLog{
		AstronautID: a.ID,
		Branch:      "army",
		Rank:        "major",
		Retired:     true,
	}
	log, err = service.AddMilitaryLog(ctx, militaryRepo, log)
	if err != nil {
		t.Errorf("Unexpected error adding military log: %v", err)
	}

	t.Run("return error when trying to update an unknown military log", func(t *testing.T) {
		ml := &model.MilitaryLog{
			AstronautID: 33,
			Branch:      "air force",
			Rank:        "major",
			Retired:     false,
		}
		err = service.UpdateMilitaryLog(ctx, militaryRepo, ml)
		if err == nil {
			t.Error("Expected an error updating military log")
		}
	})

	t.Run("returns error when trying to update an military log with invalid input", func(t *testing.T) {
		ml := &model.MilitaryLog{
			AstronautID: a.ID,
		}
		err = service.UpdateMilitaryLog(ctx, militaryRepo, ml)
		if err == nil {
			t.Error("Expected an error updating military log")
		}
	})

	t.Run("updates an existing military log", func(t *testing.T) {
		ml := &model.MilitaryLog{
			AstronautID: a.ID,
			Branch:      "navy",
			Rank:        "sergeant",
			Retired:     false,
		}
		err = service.UpdateMilitaryLog(ctx, militaryRepo, ml)
		if err != nil {
			t.Errorf("Unexpected error updating military log: %v", err)
		}

		ml, err = service.GetMilitaryLog(ctx, militaryRepo, a.ID)
		if err != nil {
			t.Errorf("Unexpected error getting military log: %v", err)
		}
		assert.Equal(t, a.ID, ml.AstronautID)
		assert.NotEqual(t, log.Branch, ml.Branch)
		assert.NotEqual(t, log.Rank, ml.Rank)
	})
}

func TestDeleteMilitaryLog(t *testing.T) {
	ctx := context.TODO()

	t.Run("returns error when trying to delete an unknown military log", func(t *testing.T) {
		astronautID := 67
		err := service.DeleteMilitaryLog(ctx, militaryRepo, astronautID)
		if err == nil {
			t.Error("Expected an error deleting military log")
		}
	})

	t.Run("deletes an existing military log", func(t *testing.T) {
		astronautID := 1
		err := service.DeleteMilitaryLog(ctx, militaryRepo, astronautID)
		if err != nil {
			t.Errorf("Unexpected error deleting military log: %v", err)
		}
	})
}
