package person

// Person represents a person.  It has an ID, a forename and a surname.
type Person interface { 
	// ID() gets the id of the person
	ID() uint64	
	//Forename gets the forename of the person
	Forename() string 
	// Surname gets the surname of the person
	Surname() string
	// String gets the person as a String
	String() string
	// SetID sets the id to the given value
	SetID(id uint64)
	// SetForename sets the forename of the person
	SetForename(forename string)
	// SetSurname sets the surname of the person
	SetSurname(surname string)
}