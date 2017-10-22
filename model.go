package rails

type Model struct {
	ID string
}

func (m *Model) Param() string {
	return m.ID
}
