package cqrs

type AggregateExpr struct {
	// Name of the aggregate
	Name string

	//Commands this aggregate can process
	Commands []*CommandExpr

	//DomainExpr is the domain this aggregate belongs to
	Domain *DomainExpr
}

func (a *AggregateExpr) Context() string {
	return "aggregate"
}

//AllEvents will return a list of all events generated from all the commands in this aggregate
func (a *AggregateExpr) AllEvents() []*EventExpr {
	em := map[string]*EventExpr{}
	emKeys := []string{}
	for _, command := range a.Commands {
		for _, eventName := range command.Events {
			if _, ok := em[eventName]; !ok {
				em[eventName] = a.Domain.Event(eventName)
				emKeys = append(emKeys, eventName)
			}
		}
	}

	//return list of events in a sorted sequence
	//sort.Strings(emKeys)
	es := make([]*EventExpr, len(em))
	for i, k := range emKeys {
		es[i] = em[k]
	}
	return es
}
