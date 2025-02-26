package commands

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	`github.com/Permify/permify/internal/keys`
	"github.com/Permify/permify/internal/repositories/mocks"
	"github.com/Permify/permify/pkg/database"
	"github.com/Permify/permify/pkg/dsl/compiler"
	"github.com/Permify/permify/pkg/dsl/schema"
	"github.com/Permify/permify/pkg/logger"
	base "github.com/Permify/permify/pkg/pb/base/v1"
	`github.com/Permify/permify/pkg/token`
	"github.com/Permify/permify/pkg/tuple"
)

var _ = Describe("check-command", func() {
	var checkCommand *CheckCommand
	l := logger.New("debug")

	// DRIVE SAMPLE

	driveSchema := `
entity user {}

entity organization {
	relation admin @user
}

entity folder {
	relation org @organization
	relation creator @user
	relation collaborator @user

	action read = collaborator
	action update = collaborator
	action delete = creator or org.admin
}

entity doc {
	relation org @organization
	relation parent @folder
	relation owner @user
	
	action read = (owner or parent.collaborator) or org.admin
	action update = owner and org.admin
	action delete = owner or org.admin
}
`

	Context("Drive Sample: Check", func() {
		It("Drive Sample: Case 1", func() {
			var err error

			// SCHEMA

			schemaReader := new(mocks.SchemaReader)

			var sch *base.IndexedSchema
			sch, err = compiler.NewSchema(driveSchema)
			Expect(err).ShouldNot(HaveOccurred())

			var en *base.EntityDefinition
			en, err = schema.GetEntityByName(sch, "doc")
			Expect(err).ShouldNot(HaveOccurred())

			schemaReader.On("ReadSchemaDefinition", "doc", "noop").Return(en, "noop", nil).Times(1)

			// RELATIONSHIPS

			relationshipReader := new(mocks.RelationshipReader)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "doc",
					Ids:  []string{"1"},
				},
				Relation: "owner",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "doc",
						Id:   "1",
					},
					Relation: "owner",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "2",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "doc",
					Ids:  []string{"1"},
				},
				Relation: "parent",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "doc",
						Id:   "1",
					},
					Relation: "parent",
					Subject: &base.Subject{
						Type:     "folder",
						Id:       "1",
						Relation: tuple.ELLIPSIS,
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "folder",
					Ids:  []string{"1"},
				},
				Relation: "collaborator",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "folder",
						Id:   "1",
					},
					Relation: "collaborator",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "1",
						Relation: "",
					},
				},
				{
					Entity: &base.Entity{
						Type: "folder",
						Id:   "1",
					},
					Relation: "collaborator",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "3",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "doc",
					Ids:  []string{"1"},
				},
				Relation: "org",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "doc",
						Id:   "1",
					},
					Relation: "org",
					Subject: &base.Subject{
						Type:     "organization",
						Id:       "1",
						Relation: tuple.ELLIPSIS,
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "organization",
					Ids:  []string{"1"},
				},
				Relation: "admin",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "organization",
						Id:   "1",
					},
					Relation: "admin",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "1",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			checkCommand = NewCheckCommand(keys.NewNoopCheckCommandKeys(), schemaReader, relationshipReader, l)

			req := &base.PermissionCheckRequest{
				Entity:        &base.Entity{Type: "doc", Id: "1"},
				Subject:       &base.Subject{Type: tuple.USER, Id: "1"},
				Permission:    "read",
				SnapToken:     token.NewNoopToken().Encode().String(),
				SchemaVersion: "noop",
			}

			var response *base.PermissionCheckResponse
			response, err = checkCommand.Execute(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(base.PermissionCheckResponse_RESULT_ALLOWED).Should(Equal(response.GetCan()))
		})

		It("Drive Sample: Case 2", func() {
			var err error

			// SCHEMA

			schemaReader := new(mocks.SchemaReader)

			var sch *base.IndexedSchema
			sch, err = compiler.NewSchema(driveSchema)
			Expect(err).ShouldNot(HaveOccurred())

			var en *base.EntityDefinition
			en, err = schema.GetEntityByName(sch, "doc")
			Expect(err).ShouldNot(HaveOccurred())

			schemaReader.On("ReadSchemaDefinition", "doc", "noop").Return(en, "noop", nil).Times(1)

			// RELATIONSHIPS

			relationshipReader := new(mocks.RelationshipReader)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "doc",
					Ids:  []string{"1"},
				},
				Relation: "owner",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "doc",
						Id:   "1",
					},
					Relation: "owner",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "2",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "doc",
					Ids:  []string{"1"},
				},
				Relation: "org",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "doc",
						Id:   "1",
					},
					Relation: "org",
					Subject: &base.Subject{
						Type:     "organization",
						Id:       "1",
						Relation: tuple.ELLIPSIS,
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "organization",
					Ids:  []string{"1"},
				},
				Relation: "admin",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "organization",
						Id:   "1",
					},
					Relation: "admin",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "1",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			checkCommand = NewCheckCommand(keys.NewNoopCheckCommandKeys(), schemaReader, relationshipReader, l)

			req := &base.PermissionCheckRequest{
				Entity:        &base.Entity{Type: "doc", Id: "1"},
				Subject:       &base.Subject{Type: tuple.USER, Id: "1"},
				Permission:    "update",
				SnapToken:     token.NewNoopToken().Encode().String(),
				SchemaVersion: "noop",
			}

			var response *base.PermissionCheckResponse
			response, err = checkCommand.Execute(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(base.PermissionCheckResponse_RESULT_DENIED).Should(Equal(response.GetCan()))
		})

		It("Drive Sample: Case 3", func() {
			var err error

			// SCHEMA

			schemaReader := new(mocks.SchemaReader)

			var sch *base.IndexedSchema
			sch, err = compiler.NewSchema(driveSchema)
			Expect(err).ShouldNot(HaveOccurred())

			var en *base.EntityDefinition
			en, err = schema.GetEntityByName(sch, "doc")
			Expect(err).ShouldNot(HaveOccurred())

			schemaReader.On("ReadSchemaDefinition", "doc", "noop").Return(en, "noop", nil).Times(1)

			// RELATIONSHIPS

			relationshipReader := new(mocks.RelationshipReader)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "doc",
					Ids:  []string{"1"},
				},
				Relation: "owner",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "doc",
						Id:   "1",
					},
					Relation: "owner",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "2",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "doc",
					Ids:  []string{"1"},
				},
				Relation: "parent",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "doc",
						Id:   "1",
					},
					Relation: "parent",
					Subject: &base.Subject{
						Type:     "folder",
						Id:       "1",
						Relation: tuple.ELLIPSIS,
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "folder",
					Ids:  []string{"1"},
				},
				Relation: "collaborator",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "folder",
						Id:   "1",
					},
					Relation: "collaborator",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "7",
						Relation: "",
					},
				},
				{
					Entity: &base.Entity{
						Type: "folder",
						Id:   "1",
					},
					Relation: "collaborator",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "3",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "doc",
					Ids:  []string{"1"},
				},
				Relation: "org",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "doc",
						Id:   "1",
					},
					Relation: "org",
					Subject: &base.Subject{
						Type:     "organization",
						Id:       "1",
						Relation: tuple.ELLIPSIS,
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "organization",
					Ids:  []string{"1"},
				},
				Relation: "admin",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "organization",
						Id:   "1",
					},
					Relation: "admin",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "7",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			checkCommand = NewCheckCommand(keys.NewNoopCheckCommandKeys(), schemaReader, relationshipReader, l)

			req := &base.PermissionCheckRequest{
				Entity:        &base.Entity{Type: "doc", Id: "1"},
				Subject:       &base.Subject{Type: tuple.USER, Id: "1"},
				Permission:    "read",
				SnapToken:     token.NewNoopToken().Encode().String(),
				SchemaVersion: "noop",
			}

			var response *base.PermissionCheckResponse
			response, err = checkCommand.Execute(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(base.PermissionCheckResponse_RESULT_DENIED).Should(Equal(response.GetCan()))
		})
	})

	// GITHUB SAMPLE

	githubSchema := `
	entity user {}
	
	entity organization {
		relation admin @user
		relation member @user
	
		action create_repository = admin or member
		action delete = admin
	}
	
	entity repository {
		relation parent @organization
		relation owner @user
	
		action push   = owner
	  action read   = owner and (parent.admin or parent.member)
	  action delete = parent.member and (parent.admin or owner)
	}
	`

	Context("Github Sample: Check", func() {
		It("Github Sample: Case 1", func() {

			var err error

			// SCHEMA

			schemaReader := new(mocks.SchemaReader)

			var sch *base.IndexedSchema
			sch, err = compiler.NewSchema(githubSchema)
			Expect(err).ShouldNot(HaveOccurred())

			var en *base.EntityDefinition
			en, err = schema.GetEntityByName(sch, "repository")
			Expect(err).ShouldNot(HaveOccurred())

			schemaReader.On("ReadSchemaDefinition", "repository", "noop").Return(en, "noop", nil).Times(1)

			// RELATIONSHIPS

			relationshipReader := new(mocks.RelationshipReader)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "repository",
					Ids:  []string{"1"},
				},
				Relation: "owner",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "repository",
						Id:   "1",
					},
					Relation: "owner",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "2",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			checkCommand = NewCheckCommand(keys.NewNoopCheckCommandKeys(), schemaReader, relationshipReader, l)

			req := &base.PermissionCheckRequest{
				Entity:        &base.Entity{Type: "repository", Id: "1"},
				Subject:       &base.Subject{Type: tuple.USER, Id: "1"},
				Permission:    "push",
				SnapToken:     token.NewNoopToken().Encode().String(),
				SchemaVersion: "noop",
			}

			var response *base.PermissionCheckResponse
			response, err = checkCommand.Execute(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(base.PermissionCheckResponse_RESULT_DENIED).Should(Equal(response.GetCan()))
		})

		It("Github Sample: Case 2", func() {

			var err error

			// SCHEMA

			schemaReader := new(mocks.SchemaReader)

			var sch *base.IndexedSchema
			sch, err = compiler.NewSchema(githubSchema)
			Expect(err).ShouldNot(HaveOccurred())

			var en *base.EntityDefinition
			en, err = schema.GetEntityByName(sch, "repository")
			Expect(err).ShouldNot(HaveOccurred())

			schemaReader.On("ReadSchemaDefinition", "repository", "noop").Return(en, "noop", nil).Times(1)

			// RELATIONSHIPS

			relationshipReader := new(mocks.RelationshipReader)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "repository",
					Ids:  []string{"1"},
				},
				Relation: "owner",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "repository",
						Id:   "1",
					},
					Relation: "owner",
					Subject: &base.Subject{
						Type:     "organization",
						Id:       "2",
						Relation: "admin",
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "organization",
					Ids:  []string{"2"},
				},
				Relation: "admin",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "organization",
						Id:   "2",
					},
					Relation: "admin",
					Subject: &base.Subject{
						Type:     "organization",
						Id:       "3",
						Relation: "member",
					},
				},
				{
					Entity: &base.Entity{
						Type: "organization",
						Id:   "2",
					},
					Relation: "admin",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "3",
						Relation: "",
					},
				},
				{
					Entity: &base.Entity{
						Type: "organization",
						Id:   "2",
					},
					Relation: "admin",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "8",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "organization",
					Ids:  []string{"3"},
				},
				Relation: "member",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "organization",
						Id:   "3",
					},
					Relation: "member",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "1",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			checkCommand = NewCheckCommand(keys.NewNoopCheckCommandKeys(), schemaReader, relationshipReader, l)

			req := &base.PermissionCheckRequest{
				Entity:        &base.Entity{Type: "repository", Id: "1"},
				Subject:       &base.Subject{Type: tuple.USER, Id: "1"},
				Permission:    "push",
				SnapToken:     token.NewNoopToken().Encode().String(),
				SchemaVersion: "noop",
			}

			var response *base.PermissionCheckResponse
			response, err = checkCommand.Execute(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(base.PermissionCheckResponse_RESULT_ALLOWED).Should(Equal(response.GetCan()))
		})

		It("Github Sample: Case 3", func() {

			var err error

			// SCHEMA

			schemaReader := new(mocks.SchemaReader)

			var sch *base.IndexedSchema
			sch, err = compiler.NewSchema(githubSchema)
			Expect(err).ShouldNot(HaveOccurred())

			var en *base.EntityDefinition
			en, err = schema.GetEntityByName(sch, "repository")
			Expect(err).ShouldNot(HaveOccurred())

			schemaReader.On("ReadSchemaDefinition", "repository", "noop").Return(en, "noop", nil).Times(1)

			// RELATIONSHIPS

			relationshipReader := new(mocks.RelationshipReader)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "repository",
					Ids:  []string{"1"},
				},
				Relation: "parent",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "repository",
						Id:   "1",
					},
					Relation: "parent",
					Subject: &base.Subject{
						Type:     "organization",
						Id:       "8",
						Relation: tuple.ELLIPSIS,
					},
				},
			}...), nil).Times(2)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "organization",
					Ids:  []string{"8"},
				},
				Relation: "member",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "organization",
						Id:   "8",
					},
					Relation: "member",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "1",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "organization",
					Ids:  []string{"8"},
				},
				Relation: "admin",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "organization",
						Id:   "8",
					},
					Relation: "admin",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "2",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			relationshipReader.On("QueryRelationships", &base.TupleFilter{
				Entity: &base.EntityFilter{
					Type: "repository",
					Ids:  []string{"1"},
				},
				Relation: "owner",
			}, token.NewNoopToken().Encode().String()).Return(database.NewTupleCollection([]*base.Tuple{
				{
					Entity: &base.Entity{
						Type: "repository",
						Id:   "1",
					},
					Relation: "owner",
					Subject: &base.Subject{
						Type:     tuple.USER,
						Id:       "7",
						Relation: "",
					},
				},
			}...), nil).Times(1)

			checkCommand = NewCheckCommand(keys.NewNoopCheckCommandKeys(), schemaReader, relationshipReader, l)

			req := &base.PermissionCheckRequest{
				Entity:        &base.Entity{Type: "repository", Id: "1"},
				Subject:       &base.Subject{Type: tuple.USER, Id: "1"},
				Permission:    "delete",
				SnapToken:     token.NewNoopToken().Encode().String(),
				SchemaVersion: "noop",
			}

			var response *base.PermissionCheckResponse
			response, err = checkCommand.Execute(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(base.PermissionCheckResponse_RESULT_DENIED).Should(Equal(response.GetCan()))
		})
	})
})
