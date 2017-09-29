package models

// Team represent team configuration
type Team struct {
	ID       int
	Name     string
	Instance string
}

// NewTeam returns initialized Team object
func NewTeam() *Team {
	return &Team{}
}

// Authentication returns Team when succeed authenticate
func (t *Team) Authentication(email, password string) (*Team, error) {
	d, err := NewDatastore()
	if err != nil {
		return nil, err
	}

	row, err := d.findTeamByEmailAndPassword(email, password)
	if err != nil {
		return nil, err
	}
	err = row.Scan(&t.ID, &t.Name, &t.Instance)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Get returns Team that assigned fields
func (t *Team) Get(id int) (*Team, error) {
	d, err := NewDatastore()
	if err != nil {
		return nil, err
	}

	row, err := d.findTeamByID(id)
	if err != nil {
		return nil, err
	}
	err = row.Scan(&t.ID, &t.Name, &t.Instance)
	if err != nil {
		return nil, err
	}
	return t, err
}

// Register register new team to the datastore
func (t *Team) Register(id int, name, email, password, instance string) error {
	d, err := NewDatastore()
	if err != nil {
		return err
	}
	return d.saveTeams(id, name, email, password, instance)
}
