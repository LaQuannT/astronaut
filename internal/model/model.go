package model

import (
	"context"
	"time"
)

type status string

const (
	Active     status = "active"
	Retired    status = "retired"
	Management status = "management"
	Deceased   status = "deceased"
)

// array (seperated by) -- missions (,), gradute major, undergrad, almamater (;) or refactor CSV
//
//	to fit my structs?
type AstronautData struct {
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

type Astronaut struct {
	ID         int    `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Gender     string `json:"gender"`
	BirthDate  string `json:"birthDate"`
	BirthPlace string `json:"birthPlace"`
}

func (a *Astronaut) Valid() (map[string]string, bool) {
	m := make(map[string]string)
	if a.FirstName == "" {
		m["FirstName"] = "first name must not be empty"
	}
	if a.LastName == "" {
		m["LastName"] = "last name must not be empty"
	}
	if a.Gender != "F" && a.Gender != "M" {
		m["Gender"] = "gender must be either 'M' or 'F'"
	}
	if a.BirthDate == "" {
		m["BirthDate"] = "birth date must not be empty"
	} else if a.BirthDate != "" {
		_, err := time.Parse(time.DateOnly, a.BirthDate)
		if err != nil {
			m["BirthDate"] = "birth date must be a valid date yyyy-mm-dd"
		}

	}
	if a.BirthPlace == "" {
		m["BirthPlace"] = "birth place must not be empty"
	}

	if len(m) > 0 {
		return m, false
	}
	return m, true
}

type Mission struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Alias         string `json:"alias"`
	DateOfMission string `json:"dateOfMission"`
	Successful    bool   `json:"successful"`
}

func (m *Mission) Valid() (map[string]string, bool) {
	problems := make(map[string]string)
	if m.Name == "" {
		problems["Name"] = "name must not be empty"
	}
	if m.DateOfMission == "" {
		problems["DateOfMission"] = "dateOfMission must not be empty"
	} else if m.DateOfMission != "" {
		_, err := time.Parse(time.DateOnly, m.DateOfMission)
		if err != nil {
			problems["DateOfMission"] = "dateOfMission must be a valid date yyyy-mm-dd"
		}
	}
	if len(problems) > 0 {
		return problems, false
	}
	return problems, true
}

type AstronautLog struct {
	AstronautID      int
	SpaceFlights     int
	SpaceFlightHours int
	SpaceWalks       int
	SpaceWalkHours   int
	Status           status
	DeathMissionID   int
	DeathDate        string
}

func (a *AstronautLog) Valid() (map[string]string, bool) {
	m := make(map[string]string)

	if a.AstronautID == 0 {
		m["astronaut_id"] = "astronaut_id must not be empty"
	}
	if a.Status != Active || a.Status != Retired || a.Status != Management || a.Status != Deceased {
		m["status"] = "status must be one of active, retired, management or deceased"
	}

	if a.DeathDate != "" {
		_, err := time.Parse(time.DateOnly, a.DeathDate)
		if err != nil {
			m["death_date"] = "death_date must be a valid date yyyy-mm-dd"
		}
	}

	if len(m) > 0 {
		return m, false
	}
	return m, true
}

type (
	User struct {
		ID        int    `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password,omitempty"`
		APIKey    string `json:"apiKey,omitempty"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}

	MilitaryLog struct {
		AstronautID int
		Branch      string
		Rank        string
		Retired     bool
	}

	Major struct {
		ID     int
		Course string
	}

	AlmaMater struct {
		ID     int
		School string
	}

	AcademicLog struct {
		AstronautID     int
		AlmaMaters      []*AlmaMater
		UnderGradMajors []*Major
		GradMajors      []*Major
	}

	AstronautRepository interface {
		CreateAstronaut(ctx context.Context, a *Astronaut) error
		FindAstronautByID(ctx context.Context, id int) (*Astronaut, error)
		UpdateAstronaut(ctx context.Context, a *Astronaut) error
		DeleteAstronaut(ctx context.Context, id int) error
		FindAstronauts(ctx context.Context) ([]*Astronaut, error)
		FindAstronautByName(ctx context.Context, name string) ([]*Astronaut, error)
	}

	UserRepository interface {
		CreateUser(ctx context.Context, u *User) error
		FindUserByID(ctx context.Context, id int) (*User, error)
		FindUserByEmail(ctx context.Context, email string) (*User, error)
		FindAllUsers(ctx context.Context) ([]*User, error)
		UpdateUser(ctx context.Context, u *User) error
		DeleteUser(ctx context.Context, id int) error
		RestUserPassword(ctx context.Context, hash string, id int) error
		GenerateNewUserAPIKey(ctx context.Context, id int) (string, error)
		GiveAdminPrivileges(ctx context.Context, id int) error
		RevokeAdminPrivileges(ctx context.Context, id int) error
	}

	MissionRepository interface {
		CreateMission(ctx context.Context, m *Mission) error
		FindMissionByID(ctx context.Context, id int) (*Mission, error)
		FindMissionByNameOrAlias(ctx context.Context, target string) ([]*Mission, error)
		FindAllMissions(ctx context.Context) ([]*Mission, error)
		UpdateMission(ctx context.Context, m *Mission) error
		CreateAstronautMission(ctx context.Context, astronautID, missionID int) error
		FindMissionsByAstronaut(ctx context.Context, astronautID int) ([]*Mission, error)
		DeleteAstronautMission(ctx context.Context, astronautID, missionID int) error
		DeleteMission(ctx context.Context, missionID int) error
	}

	MilitaryLogRepository interface {
		CreateMilitaryLog(ctx context.Context, m *MilitaryLog) error
		FindMilitaryLog(ctx context.Context, astronautID int) (*MilitaryLog, error)
		FindAllMilitaryLogs(ctx context.Context) ([]*MilitaryLog, error)
		UpdateMilitaryLog(ctx context.Context, m *MilitaryLog) error
		DeleteMilitaryLog(ctx context.Context, astronautID int) error
	}

	AcademicLogRepository interface {
		CreateMajor(ctx context.Context, m *Major) error
		CreateAlmaMater(ctx context.Context, a *AlmaMater) error
		AddUnderGradMajor(ctx context.Context, astronautID, majorID int) error
		AddGradMajor(ctx context.Context, astronautID, majorID int) error
		AddAstronautAlmaMater(ctx context.Context, astronautID, almaMaterID int) error
		UpdateMajor(ctx context.Context, m *Major) error
		UpdateAlmaMater(ctx context.Context, a *AlmaMater) error
		FindMajorByID(ctx context.Context, id int) (*Major, error)
		FindAlmaMaterByID(ctx context.Context, id int) (*AlmaMater, error)
		FindAstronautUnderGradMajors(ctx context.Context, astronautID int) ([]*Major, error)
		FindAstronautGradMajors(ctx context.Context, astronautID int) ([]*Major, error)
		FindAstronautAlmaMaters(ctx context.Context, astronautID int) ([]*AlmaMater, error)
		DeleteMajor(ctx context.Context, id int) error
		DeleteAstronautUnderGradMajor(ctx context.Context, astronautID, majorID int) error
		DeleteAstronautGradMajor(ctx context.Context, astronautID, majorID int) error
		DeleteAlmaMater(ctx context.Context, id int) error
		DeleteAstronautAlmaMater(ctx context.Context, astronautID, majorID int) error
		GetAcademicLog(ctx context.Context, astronautID int) (*AcademicLog, error)
	}

	AstronautLogRepository interface {
		CreateAstronautLog(ctx context.Context, a *AstronautLog) error
		FindAstronautLogById(ctx context.Context, astronautID int) (*AstronautLog, error)
		FindAstronautLogs(ctx context.Context) ([]*AstronautLog, error)
		UpdateAstronautLog(ctx context.Context, a *AstronautLog) error
		DeleteAstronautLog(ctx context.Context, astronautID int) error
	}
)
