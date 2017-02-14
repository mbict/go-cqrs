package cqrs

type ProjectionExpr struct {
	// Name of the projection
	Name string

	//ReadReposiotories who are based on this projection
	ReadRepositories []*RepositoryExpr

	//HandlesEvents are the names of the events this projection accepts
	HandlesEvent []string

	///DomainExpr this projection is mainly generated from
	Domain *DomainExpr
}

func (p *ProjectionExpr) Context() string {
	return "projection"
}

//AllEvents will return a list of all events handled by this projection
func (p *ProjectionExpr) AllEvents() []*EventExpr {
	list := []*EventExpr{}
	for _, eventName := range p.HandlesEvent {
		list = append(list, p.Domain.Event(eventName))
	}
	return list
}
