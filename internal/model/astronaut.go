package model

import "context"

type (

	// properties ending in 'Str' are list in string form in csv file
	// (seperated by) -- missions (,) gradute major, undergrad, almamater (;)
	Astronaut struct {
		ID                    int      `json:"id"`
		Name                  string   `json:"name" csv:"Name"`
		Year                  int      `json:"year" csv:"Year"`
		Group                 int      `json:"group" csv:"Group"`
		Status                string   `json:"status" csv:"Status"`
		BirthDate             string   `json:"birthDate" csv:"Birth Date"`
		BirthPlace            string   `json:"birthPlace" csv:"Birth Place"`
		Gender                string   `json:"gender" csv:"Gender"`
		AlmaMaterStr          string   `json:"omitempty" csv:"Alma Mater"`
		UndergraduateMajorStr string   `json:"omitempty" csv:"Undergraduate Major"`
		GraduateMajorStr      string   `json:"omitempty" csv:"Graduate Major"`
		MilitaryRank          string   `json:"militaryRank" csv:"Military Rank"`
		MilitaryBranch        string   `json:"militaryBranch" csv:"Military Branch"`
		SpaceFlights          int      `json"spaceFlights" csv:"Space Flights"`
		SpaceFlightHours      int      `json:"spaceFlightHours" csv:"Space Flight (hr)"`
		SpaceWalks            int      `json:"spaceWalks" csv:"Space Walks"`
		SpaceWalkHours        int      `json:"spaceWalkHours" csv:"Space Walk (hr)"`
		MissionStr            string   `json:"omitempty" csv:"Missions"`
		DeathDate             string   `json:"deathDate" csv:"Death Date"`
		DeathMission          string   `json:"deathMission" csv:"Death Mission"`
		Missions              []string `json:"missions"`
		UndergraduateMajor    []string `json:"undergraduateMajor"`
		GraduateMajor         []string `json:"graduateMajor"`
		AlmaMater             []string `json:"almaMater"`
	}

	// need to add Search methods for popular search categories
	AstronautStore interface {
		Create(ctx context.Context, a *Astronaut) (int, error)
		List(ctx context.Context, limit, offset int) ([]*Astronaut, error)
		Get(ctx context.Context, id int) (*Astronaut, error)
		Update(ctx context.Context, a *Astronaut) error
		Delete(ctx context.Context, id int) error
	}

	AstronautUsecase interface {
		Create(ctx context.Context, a *Astronaut) (*Astronaut, error)
		List(ctx context.Context, limit, offset int) ([]*Astronaut, error)
		Get(ctx context.Context, id int) (*Astronaut, error)
		Update(ctx context.Context, a *Astronaut) (*Astronaut, error)
		Delete(ctx context.Context, id int) error
	}
)
