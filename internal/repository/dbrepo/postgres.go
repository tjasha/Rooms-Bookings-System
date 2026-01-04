package dbrepo

// functions for usage in handlers - should be add to postgresDBRepo interface
func (m *postgresDBRepo) AllUsers() bool {
	return true
}
