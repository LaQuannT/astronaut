package model

import "context"

// array (seperated by) -- missions (,), gradute major, undergrad, almamater (;)
type (
	AstronautData struct {
		Name               string   `json:"name" csv:"Name"`
		Year               int      `json:"year" csv:"Year"`
		Group              int      `json:"group" csv:"Group"`
		Status             string   `json:"status" csv:"Status"`
		BirthDate          string   `json:"birthDate" csv:"Birth Date"`
		BirthPlace         string   `json:"birthPlace" csv:"Birth Place"`
		Gender             string   `json:"gender" csv:"Gender"`
		AlmaMater          []string `json:"almaMater" csv:"Alma Mater"`
		UndergraduateMajor []string `json:"undergraduateMajor" csv:"Undergraduate Major"`
		GraduateMajor      []string `json:"graduateMajor" csv:"Graduate Major"`
		MilitaryRank       string   `json:"militaryRank" csv:"Military Rank"`
		MilitaryBranch     string   `json:"militaryBranch" csv:"Military Branch"`
		SpaceFlights       int      `json:"spaceFlights" csv:"Space Flights"`
		SpaceFlightHours   int      `json:"spaceFlightHours" csv:"Space Flight (hr)"`
		SpaceWalks         int      `json:"spaceWalks" csv:"Space Walks"`
		SpaceWalkHours     int      `json:"spaceWalkHours" csv:"Space Walk (hr)"`
		Missions           []string `json:"missions" csv:"Missions"`
		DeathDate          string   `json:"deathDate" csv:"Death Date"`
		DeathMission       string   `json:"deathMission" csv:"Death Mission"`
	}

	Astronaut struct {
		ID         int    `json:"id"`
		FirstName  string `json:"firstName"`
		LastName   string `json:"lastName"`
		Gender     string `json:"gender"`
		BirthDate  string `json:"birthDate"`
		BirthPlace string `json:"birthPlace"`
	}

	// add method for searching by name
	AstronautRepository interface {
		CreateAstronaut(ctx context.Context, a Astronaut) error
		FindAstronautByID(ctx context.Context, id int) (*Astronaut, error)
		UpdateAstronaut(ctx context.Context, a *Astronaut) error
		DeleteAstronaut(ctx context.Context, id int) error
	}
)
