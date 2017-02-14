package spec

import (
	. "github.com/mbict/go-cqrs/gogen/dsl"
	. "github.com/mbict/gogen"
	. "github.com/mbict/gogen/dsl"
)

var CustomerModel = Usertype("Customer", func() {

})

var _ = Domain("Ordering", func() {

	Aggregate("Customer", func() {
		Command("CreateCustomer", func() {
			Params(func() {
				Attribute("Name", String)
				Attribute("Number", String)
			})

			Event("CustomerCreated", func() {
				Attribute("Name", String)
				Attribute("Number", String)
			})
		})

		Command("UpdateCustomer", func() {
			Params(func() {
				Attribute("Name", String)
				Attribute("Number", String)
			})

			Event("CustomerNameChanged", func() {
				Attribute("Name", String)
			})
			Event("CustomerNumberChanged", func() {
				Attribute("Number", String)
			})
		})

		//no params needed (except for the ID who is auto added)
		Command("RemoveCustomer", func() {
			Event("CustomerRemoved")
		})
	})

	Projection("customernumbers", func() {
		Repository("numbers", CustomerModel, func(){
			Filter(func() {
				Attribute("Id", String)
				Attribute("Number", String)
			})
		})

		HandlesEvents("CustomerCreated", "CustomerNumberChanged", "CustomerRemoved")
	})

	Projection("customers", func() {
		Filter(func() {
			Attribute("Id", String)
			Attribute("Name", String)
			Attribute("Number", String)
		})

		HandlesEvents("CustomerCreated", "CustomerNameChanged", "CustomerNumberChanged", "CustomerRemoved")
	})

})
