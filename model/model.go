package model

type xModel struct {
	db
}
type db interface {
	SelectPeople() ([]*Person, error)
}
type Person struct {
	Id          int64
	First, Last string
}

func New(db db) *xModel {
	return &xModel{
		db: db,
	}
}

func (m *xModel) People() ([]*Person, error) {
	return m.SelectPeople()
	var x []Person
}