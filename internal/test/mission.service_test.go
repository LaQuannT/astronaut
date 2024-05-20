package test

import (
	"context"
	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/LaQuannT/astronaut-api/internal/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddMission(t *testing.T) {
	err := clearTables(dbConn)
	if err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}
	ctx := context.TODO()

	t.Run("returns an error for invalid mission input", func(t *testing.T) {
		m := &model.Mission{}

		m, err := service.AddMission(ctx, missionRepo, m)
		if err == nil {
			t.Errorf("Expected error for invalid mission data")
		}
		assert.Nil(t, m)
	})

	t.Run("adds new mission", func(t *testing.T) {
		m := &model.Mission{
			Name:          "flight 102",
			Alias:         "",
			DateOfMission: "2022-01-01",
			Successful:    true,
		}
		m, err = service.AddMission(ctx, missionRepo, m)
		if err != nil {
			t.Fatalf("Unexpected error adding mission: %v", err)
		}
		assert.NotEqual(t, 0, m.ID)
	})

	t.Run("returns an error if mission already exists", func(t *testing.T) {
		m := &model.Mission{
			Name:          "flight 102",
			Alias:         "",
			DateOfMission: "2022-01-01",
			Successful:    false,
		}
		m, err = service.AddMission(ctx, missionRepo, m)
		if err == nil {
			t.Errorf("Expected error for duplicate mission")
		}
	})

}

func TestGetMission(t *testing.T) {
	err := clearTables(dbConn)
	if err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}
	ctx := context.TODO()

	m := &model.Mission{
		Name:          "flight 103",
		DateOfMission: "2022-01-01",
		Successful:    false,
	}
	m, err = service.AddMission(ctx, missionRepo, m)
	if err != nil {
		t.Fatalf("Unexpected error adding mission: %v", err)
	}

	t.Run("returns an error for an unknown mission", func(t *testing.T) {
		id := 34
		mission, err := service.GetMission(ctx, missionRepo, id)
		if err == nil {
			t.Errorf("Expected error for unknown mission")
		}
		assert.Nil(t, mission)
	})

	t.Run("returns a mission", func(t *testing.T) {
		id := 1
		mission, err := service.GetMission(ctx, missionRepo, id)
		if err != nil {
			t.Fatalf("Unexpected error getting mission: %v", err)
		}
		assert.NotNil(t, mission)
		assert.Equal(t, id, mission.ID)
		assert.Equal(t, m.Name, mission.Name)
	})

}

func TestGetMissions(t *testing.T) {
	err := clearTables(dbConn)
	if err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}
	ctx := context.TODO()

	t.Run("returns nil if no missions found", func(t *testing.T) {
		missions, err := service.GetMissions(ctx, missionRepo)
		if err != nil {
			t.Fatalf("Unexpected error getting missions: %v", err)
		}
		assert.Nil(t, missions)
	})

	t.Run("returns a list of missions", func(t *testing.T) {
		x := &model.Mission{
			Name:          "flight x",
			DateOfMission: "2022-01-01",
			Successful:    true,
		}
		y := &model.Mission{
			Name:          "flight y",
			DateOfMission: "2022-01-01",
			Successful:    false,
		}

		flights := []*model.Mission{x, y}
		for _, f := range flights {
			_, err := service.AddMission(ctx, missionRepo, f)
			if err != nil {
				t.Fatalf("Unexpected error adding mission: %v", err)
			}
		}
	})

	missions, err := service.GetMissions(ctx, missionRepo)
	if err != nil {
		t.Fatalf("Unexpected error getting missions: %v", err)
	}
	assert.Len(t, missions, 2)
}

func TestSearchMissionName(t *testing.T) {
	ctx := context.TODO()

	t.Run("returns nil for an unknown mission", func(t *testing.T) {
		target := "seek and find"

		missions, err := service.SearchMissionName(ctx, missionRepo, target)
		if err != nil {
			t.Errorf("Unexpected error for unknown mission: %v", err)
		}
		assert.Nil(t, missions)
	})

	t.Run("returns a list of missions containing target name", func(t *testing.T) {
		target := "flight"

		missions, err := service.SearchMissionName(ctx, missionRepo, target)
		if err != nil {
			t.Fatalf("Unexpected error getting missions: %v", err)
		}
		assert.Len(t, missions, 2)
	})

	t.Run("returns an mission matching exact target name", func(t *testing.T) {
		target := "flight y"
		missions, err := service.SearchMissionName(ctx, missionRepo, target)
		if err != nil {
			t.Fatalf("Unexpected error getting missions: %v", err)
		}
		assert.Len(t, missions, 1)
		assert.Equal(t, missions[0].Name, target)
	})
}

func TestUpdateMission(t *testing.T) {
	err := clearTables(dbConn)
	if err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}

	ctx := context.TODO()

	m := &model.Mission{
		Name:          "flight 340",
		DateOfMission: "2022-01-01",
		Successful:    true,
	}

	_, err = service.AddMission(ctx, missionRepo, m)
	if err != nil {
		t.Fatalf("Unexpected error adding mission: %v", err)
	}

	t.Run("returns an error for invalid mission input", func(t *testing.T) {
		mission := &model.Mission{
			ID:            1,
			Name:          "",
			Alias:         "space flight",
			DateOfMission: "",
			Successful:    false,
		}

		if err := service.UpdateMission(ctx, missionRepo, mission); err == nil {
			t.Errorf("Expected error for invalid mission data")
		}
	})

	t.Run("returns an updated mission", func(t *testing.T) {
		mission := &model.Mission{
			ID:            1,
			Name:          "SpaceForce Flight 1",
			Alias:         "force 12",
			DateOfMission: m.DateOfMission,
			Successful:    m.Successful,
		}

		if err := service.UpdateMission(ctx, missionRepo, mission); err != nil {
			t.Errorf("Unexpected error updating mission: %v", err)
		}

		m, err := service.GetMission(ctx, missionRepo, mission.ID)
		if err != nil {
			t.Errorf("Unexpected error getting mission: %v", err)
		}
		assert.Equal(t, mission.Name, m.Name)
		assert.Equal(t, mission.Alias, m.Alias)
	})
}

func TestRegisterAstronautToMission(t *testing.T) {
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

	_, err = service.AddAstronaut(ctx, a, astroRepo)
	if err != nil {
		t.Fatalf("Unexpected error adding Astronaut: %v", err)
	}

	m := &model.Mission{
		Name:          "flight 64",
		DateOfMission: "2022-01-01",
		Successful:    true,
	}
	_, err = service.AddMission(ctx, missionRepo, m)
	if err != nil {
		t.Fatalf("Unexpected error adding mission: %v", err)
	}

	t.Run("returns an error for an unknown astronaut", func(t *testing.T) {
		unknownAstronautID := 44
		if err := service.RegisterAstronautToMission(ctx, missionRepo, unknownAstronautID, m.ID); err == nil {
			t.Errorf("Expected error for unknown astronaut")
		}
	})

	t.Run("returns an error for an unknown mission", func(t *testing.T) {
		unknownMissionID := 90
		if err := service.RegisterAstronautToMission(ctx, missionRepo, a.ID, unknownMissionID); err == nil {
			t.Errorf("Expected error for unknown mission")
		}

		t.Run("registers a astronaut to a mission", func(t *testing.T) {
			if err := service.RegisterAstronautToMission(ctx, missionRepo, a.ID, m.ID); err != nil {
				t.Errorf("unexpected error registering astronaut to mission: %v", err)
			}
		})

	})
}

func TestGetMissionsByAstronaut(t *testing.T) {
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

	f64 := &model.Mission{
		Name:          "flight 64",
		DateOfMission: "2022-01-01",
		Successful:    true,
	}

	f65 := &model.Mission{
		Name:          "flight 65",
		DateOfMission: "2022-01-01",
		Successful:    false,
	}

	ms := []*model.Mission{f64, f65}
	for _, m := range ms {
		m, err = service.AddMission(ctx, missionRepo, m)
		if err != nil {
			t.Fatalf("Unexpected error adding mission: %v", err)
		}
		err = service.RegisterAstronautToMission(ctx, missionRepo, a.ID, m.ID)
		if err != nil {
			t.Fatalf("Unexpected error registering astronaut to mission: %v", err)
		}
	}

	t.Run("returns nil for unknown astronaut", func(t *testing.T) {
		unknownAstronautID := 89
		missions, err := service.GetMissionsByAstronaut(ctx, missionRepo, unknownAstronautID)
		if err != nil {
			t.Error("Expected error getting missions")
		}
		assert.Nil(t, missions)
	})

	t.Run("returns a list of missions an astronaut is registered on", func(t *testing.T) {

		missions, err := service.GetMissionsByAstronaut(ctx, missionRepo, a.ID)
		if err != nil {
			t.Errorf("Unexpected error getting missions: %v", err)
		}

		assert.Len(t, missions, 2)
	})
}

func TestRemoveAstronautFromMission(t *testing.T) {
	err := clearTables(dbConn)
	if err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}

	ctx := context.TODO()

	a := &model.Astronaut{
		FirstName:  "sarah",
		LastName:   "tucker",
		Gender:     "F",
		BirthDate:  "1999-01-01",
		BirthPlace: "london,uk",
	}

	a, err = service.AddAstronaut(ctx, a, astroRepo)
	if err != nil {
		t.Fatalf("Unexpected error adding Astronaut: %v", err)
	}

	m := &model.Mission{
		Name:          "space one",
		DateOfMission: "2022-01-01",
		Successful:    true,
	}

	m, err = service.AddMission(ctx, missionRepo, m)
	if err != nil {
		t.Fatalf("Unexpected error adding mission: %v", err)
	}

	err = service.RegisterAstronautToMission(ctx, missionRepo, a.ID, m.ID)
	if err != nil {
		t.Fatalf("Unexpected error registering astronaut to mission: %v", err)
	}

	t.Run("returns an error for unknown astronaut", func(t *testing.T) {
		astronautID := 48
		missionID := 1

		err := service.RemoveAstronautFromMission(ctx, missionRepo, astronautID, missionID)
		if err == nil {
			t.Errorf("Expected error for unknown astronaut")
		}
	})

	t.Run("returns an error for an unknown mission", func(t *testing.T) {
		missionID := 90
		astronautID := 1
		err := service.RemoveAstronautFromMission(ctx, missionRepo, astronautID, missionID)
		if err == nil {
			t.Errorf("Expected error for unknown mission")
		}
	})

	t.Run("removes an astronaut from a mission", func(t *testing.T) {
		missionID := 1
		astronautID := 1

		err := service.RemoveAstronautFromMission(ctx, missionRepo, astronautID, missionID)
		if err != nil {
			t.Errorf("Unexpected error removing mission: %v", err)
		}

		missons, err := service.GetMissionsByAstronaut(ctx, missionRepo, a.ID)
		if err != nil {
			t.Errorf("Unexpected error getting missions: %v", err)
		}
		assert.Len(t, missons, 0)

	})
}

func TestDeleteMission(t *testing.T) {
	err := clearTables(dbConn)
	if err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}
	ctx := context.TODO()

	m := &model.Mission{
		Name:          "space 3",
		DateOfMission: "2022-01-01",
		Successful:    true,
	}

	m, err = service.AddMission(ctx, missionRepo, m)
	if err != nil {
		t.Fatalf("Unexpected error adding mission: %v", err)
	}

	t.Run("returns nil for unknown mission", func(t *testing.T) {
		missionID := 90
		err := service.DeleteMission(ctx, missionRepo, missionID)
		if err == nil {
			t.Errorf("Expected error deleting mission")
		}
	})

	t.Run("deletes a mission", func(t *testing.T) {
		err := service.DeleteMission(ctx, missionRepo, m.ID)
		if err != nil {
			t.Errorf("Unexpected error deleting mission: %v", err)
		}
		mission, err := service.GetMission(ctx, missionRepo, m.ID)
		if err == nil {
			t.Error("Expected error getting an unknown mission")
		}
		assert.Nil(t, mission)
	})
}
