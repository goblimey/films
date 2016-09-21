package services

import (
	peopleRepo "github.com/goblimey/films/repositories/people"
	"github.com/goblimey/films/retrofit/template"
)

type ConcreteServices struct {
	peopleRepo  peopleRepo.Repository
	templateMap *map[string]template.Template
}

func (cs ConcreteServices) GetPeopleRepository() peopleRepo.Repository {
	return cs.peopleRepo
}

// Template returns an HTML template, given a CRUD operation (Index, Edit etc).
func (cs ConcreteServices) Template(operation string) template.Template {
	return (*cs.templateMap)[operation]
}

func (cs *ConcreteServices) SetPeopleRepository(repo peopleRepo.Repository) {
	cs.peopleRepo = repo
}

func (cs *ConcreteServices) SetTemplates(
	templateMap *map[string]template.Template) {

	cs.templateMap = templateMap
}
