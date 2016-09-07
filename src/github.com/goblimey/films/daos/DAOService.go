package people

import (
	peopleDAO "github.com/goblimey/films/daos/people"
)

// DAOService provides Data Access Objects (DAOs).
type DAOService struct {
	DAOField peopleDAO.DAO
}

// MakeDAOService is a factory function that creates and returns a DAO service.
func MakeDAOService(dao peopleDAO.DAO) DAOService {
	return DAOService{dao}
}

// SetDAO sets the DAO.
func (dfi *DAOService) SetDAO(dao peopleDAO.DAO) {
	dfi.DAOField = dao
}

// DAO returns a pointer to the DAO.
func (dfi DAOService) DAO() peopleDAO.DAO {
	return dfi.DAOField
}
