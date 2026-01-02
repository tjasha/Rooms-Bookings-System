package forms

type errors map[string][]string

// adds error message for a given field
// field = which fields form has
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// check if the field has an error and return first one
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	return es[0]
}
