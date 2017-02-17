package dsl

//
//import (
//	. "gopkg.in/check.v1"
//	. "stalling/stalling/design/cqrs"
//	"stalling/stalling/dslengine"
//	"stalling/stalling/design/common"
//)
//
//type DslSuite struct{}
//
//var _ = Suite(&DslSuite{})
//
//func (s *DslSuite) SetUpTest(c *C) {
//	*Root = *NewCqrsDefinition()
//	dslengine.Reset()
//
//	createDsl()
//
//	err := dslengine.Run()
//	if err != nil {
//		c.Fatal("failed to create example dsl, error : `%s`", err.Error())
//	}
//}
//
//func createDsl() {
//
//	RootAggregate("item", func() {
//		HandlesCommands(func() {
//			CommandExpr("createItem", func() {
//				Attribute("name", String, func() {
//
//				})
//				Attribute("price", Float)
//				Attribute("status", func() {
//					Enum("active", "pending", "inactive")
//				})
//
//				Required("name", "price")
//			})
//			CommandExpr("updateItem", func() {
//				Attribute("name", String)
//				Attribute("price", Float)
//			})
//			CommandExpr("deleteItem")
//		})
//
//		GenerateEvents(func() {
//			Event("itemCreated", func() {
//				Attribute("name", String)
//				Attribute("price", Float)
//			})
//			Event("itemNameChanged", func() {
//				Attribute("name", String)
//			})
//			Event("itemPriceChanged", func() {
//				Attribute("price", Float)
//			})
//			Event("itemDeleted")
//		})
//	})
//
//	ProjectionExpr("item", func() {
//		HandlesEvents(func() {
//			Event("itemCreated")
//			Event("itemNameChanged")
//			Event("itemDeleted")
//		})
//	})
//}
//
//func (s *DslSuite) TestRootHasEvents(c *C) {
//	c.Assert(Root.Events, HasLen, 4)
//	c.Assert(Root.Events["itemcreated"], NotNil)
//	c.Assert(Root.Events["itemnamechanged"], NotNil)
//	c.Assert(Root.Events["itempricechanged"], NotNil)
//	c.Assert(Root.Events["itemdeleted"], NotNil)
//}
//
//func (s *DslSuite) TestRootHasProjections(c *C) {
//	c.Assert(Root.Projections, HasLen, 1)
//	c.Assert(Root.Projections["item"], NotNil)
//}
//
//func (s *DslSuite) TestRootHasAggregates(c *C) {
//	c.Assert(Root.Aggregates, HasLen, 1)
//	c.Assert(Root.Aggregates["item"], NotNil)
//}
//
//func (s *DslSuite) TestRootHasCommands(c *C) {
//	c.Assert(Root.Commands, HasLen, 3)
//	c.Assert(Root.Commands["createitem"], NotNil)
//	c.Assert(Root.Commands["updateitem"], NotNil)
//	c.Assert(Root.Commands["deleteitem"], NotNil)
//}
//
//func (s *DslSuite) TestAggregateItemHasEvents(c *C) {
//	c.Assert(Root.Aggregates["item"].GeneratesEvents, HasLen, 4)
//	c.Assert(Root.Aggregates["item"].GeneratesEvents["itemcreated"], NotNil)
//	c.Assert(Root.Aggregates["item"].GeneratesEvents["itemnamechanged"], NotNil)
//	c.Assert(Root.Aggregates["item"].GeneratesEvents["itempricechanged"], NotNil)
//	c.Assert(Root.Aggregates["item"].GeneratesEvents["itemdeleted"], NotNil)
//}
//
//func (s *DslSuite) TestAggregateItemHasCommands(c *C) {
//	c.Assert(Root.Aggregates["item"].HandlesCommands, HasLen, 3)
//	c.Assert(Root.Aggregates["item"].HandlesCommands["createitem"], NotNil)
//	c.Assert(Root.Aggregates["item"].HandlesCommands["updateitem"], NotNil)
//	c.Assert(Root.Aggregates["item"].HandlesCommands["deleteitem"], NotNil)
//}
//
//func (s *DslSuite) TestProjectionItemHandlesEvents(c *C) {
//	c.Assert(Root.Projections["item"].HandlesEvents, HasLen, 3)
//	c.Assert(Root.Projections["item"].HandlesEvents["itemcreated"], NotNil)
//	c.Assert(Root.Projections["item"].HandlesEvents["itemnamechanged"], NotNil)
//	c.Assert(Root.Projections["item"].HandlesEvents["itemdeleted"], NotNil)
//}
//
//func (s *DslSuite) TestCommandHasAttributes(c *C) {
//	c.Assert(Root.Commands["createitem"].Attribute().Type.(common.Object), HasLen, 3)
//	c.Assert(Root.Commands["updateitem"].Attribute().Type.(common.Object), HasLen, 2)
//	c.Assert(Root.Commands["deleteitem"].Attribute().Type.(common.Object), HasLen, 0)
//}
//
//func (s *DslSuite) TestEventHasAttributes(c *C) {
//	c.Assert(Root.Events["itemcreated"].Attribute().Type.(common.Object), HasLen, 2)
//	c.Assert(Root.Events["itemnamechanged"].Attribute().Type.(common.Object), HasLen, 1)
//	c.Assert(Root.Events["itempricechanged"].Attribute().Type.(common.Object), HasLen, 1)
//	c.Assert(Root.Events["itemdeleted"].Attribute().Type.(common.Object), HasLen, 0)
//}
//
//func (s *DslSuite) TestCommandCreateItemDefinition(c *C) {
//	commandDefinition := Root.Commands["createitem"]
//	params, ok := common.ToObject(commandDefinition.Attribute().Type)
//
//	c.Assert(ok, Equals, true)
//
//	c.Assert(commandDefinition.Name, Equals, "createItem")
//	c.Assert(params["name"], NotNil)
//	c.Assert(params["name"].Type, Equals, String)
//	c.Assert(params["price"], NotNil)
//	c.Assert(params["price"].Type, Equals, Float)
//	c.Assert(params["status"], NotNil)
//	//c.Assert(params["name"].Type, Equals, String)
//}
//
//func (s *DslSuite) TestCommandUpdateItemDefinition(c *C) {
//	commandDefinition := Root.Commands["updateitem"]
//	params, ok := common.ToObject(commandDefinition.Attribute().Type)
//
//	c.Assert(ok, Equals, true)
//
//	c.Assert(commandDefinition.Name, Equals, "updateItem")
//	c.Assert(params["name"], NotNil)
//	c.Assert(params["name"].Type, Equals, String)
//	c.Assert(params["price"], NotNil)
//	c.Assert(params["price"].Type, Equals, Float)
//}
