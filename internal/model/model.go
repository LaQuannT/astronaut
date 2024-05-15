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

	// add method for searching by name
	AstronautRepository interface {
		CreateAstronaut(ctx context.Context, a Astronaut) error
		FindAstronautByID(ctx context.Context, id int) (*Astronaut, error)
		UpdateAstronaut(ctx context.Context, a *Astronaut) error
		DeleteAstronaut(ctx context.Context, id int) error
	}

	UserRepository interface {
		CreateUser(ctx context.Context, u User) error
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
)
