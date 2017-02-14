package cqrs

import "sort"

type DomainExpr struct {
	//Name of the domain
	Name string

	//Aggregates are all the root aggregates in this domain
	Aggregates []*AggregateExpr

	//Events are all the defined events
	Events map[string]*EventExpr

	Projections []*ProjectionExpr
}

func (d *DomainExpr) Context() string {
	return "domain"
}

// AllEvents returns all the generated events by the aggregates and handled events by the projections
func (d *DomainExpr) AllEvents() []*EventExpr {
	em := map[string]*EventExpr{}
	emKeys := []string{}
	for _, a := range d.Aggregates {
		for _, c := range a.Commands {
			for _, eventName := range c.Events {
				if _, ok := em[eventName]; !ok {
					em[eventName] = d.Event(eventName)
					emKeys = append(emKeys, eventName)
				}
			}
		}
	}

	for _, p := range d.Projections {
		for _, eventName := range p.HandlesEvent {
			if _, ok := em[eventName]; !ok {
				em[eventName] = d.Event(eventName)
				emKeys = append(emKeys, eventName)
			}
		}
	}

	//return list of events in a sorted sequence
	sort.Strings(emKeys)
	es := make([]*EventExpr, len(em))
	for i, k := range emKeys {
		es[i] = em[k]
	}
	return es
}

// Event will return the event registered under the name
func (d *DomainExpr) Event(name string) *EventExpr {
	e, ok := d.Events[name]
	if !ok {
		return nil
	}
	return e
}

func (d *DomainExpr) AllReadRepositories() (res []*RepositoryExpr) {
	for _, p := range d.Projections {
		res = append(res, p.ReadRepositories...)
	}
	return res
}