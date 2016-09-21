package services

import (
	peopleRepo "github.com/goblimey/films/repositories/people"
	"github.com/goblimey/films/retrofit/template"
)

type Services interface {
	GetPeopleRepository() peopleRepo.Repository

	Template(operation string) template.Template

	SetPeopleRepository(dao peopleRepo.Repository)

	SetTemplates(templateMap *map[string]template.Template)
}
