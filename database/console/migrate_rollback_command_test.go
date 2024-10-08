package console

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ormcontract "github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/database/gorm"
	configmock "github.com/goravel/framework/mocks/config"
	consolemocks "github.com/goravel/framework/mocks/console"
	"github.com/goravel/framework/support/docker"
	"github.com/goravel/framework/support/env"
)

func TestMigrateRollbackCommand(t *testing.T) {
	if env.IsWindows() {
		t.Skip("Skipping tests of using docker")
	}

	var (
		mockConfig *configmock.Config
		query      ormcontract.Query
	)

	beforeEach := func() {
		mockConfig = &configmock.Config{}
	}

	tests := []struct {
		name  string
		setup func()
	}{
		{
			name: "mysql",
			setup: func() {
				var err error
				docker := gorm.NewMysqlDocker(docker.Mysql())
				query, err = docker.New()
				assert.Nil(t, err)
				mockConfig = docker.MockConfig
				createMysqlMigrations()

			},
		},
		{
			name: "postgres",
			setup: func() {
				var err error
				docker := gorm.NewPostgresDocker(docker.Postgres())
				query, err = docker.New()
				assert.Nil(t, err)
				mockConfig = docker.MockConfig
				createPostgresMigrations()
			},
		},
		{
			name: "sqlserver",
			setup: func() {
				var err error
				docker := gorm.NewSqlserverDocker(docker.Sqlserver())
				query, err = docker.New()
				assert.Nil(t, err)
				mockConfig = docker.MockConfig
				createSqlserverMigrations()
			},
		},
		{
			name: "sqlite",
			setup: func() {
				var err error
				docker := gorm.NewSqliteDocker(docker.Sqlite())
				query, err = docker.New()
				assert.Nil(t, err)
				mockConfig = docker.MockConfig
				createSqliteMigrations()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			beforeEach()
			test.setup()

			mockContext := &consolemocks.Context{}
			mockContext.On("Option", "step").Return("1").Once()

			migrateCommand := NewMigrateCommand(mockConfig)
			assert.Nil(t, migrateCommand.Handle(mockContext))

			var agent Agent
			err := query.Where("name", "goravel").FirstOrFail(&agent)
			assert.Nil(t, err)
			assert.True(t, agent.ID > 0)

			migrateRollbackCommand := NewMigrateRollbackCommand(mockConfig)
			assert.Nil(t, migrateRollbackCommand.Handle(mockContext))

			var agent1 Agent
			err = query.Where("name", "goravel").FirstOrFail(&agent1)
			assert.Error(t, err)

			mockContext.AssertExpectations(t)
			removeMigrations()
		})
	}
}
