package database

import (
	"github.com/digkill/telegram-chatgpt/internal/components/driver"
	"github.com/digkill/telegram-chatgpt/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"os"
)

type SqlDriver interface {
	GetSqlDb() *sqlx.DB
	GetDriverMigration(migrationTable string, databaseName string) (database.Driver, error)
	GetDataBaseNameMigration() string
}

type MigrationComponent struct {
	driver          SqlDriver
	migrationConfig *config.MigrationConfig
}

func (migration *MigrationComponent) initDatabaseInstance() (*migrate.Migrate, error) {
	driverMigration, err := migration.driver.GetDriverMigration(
		migration.migrationConfig.MigrationTable,
		migration.migrationConfig.DatabaseName,
	)

	if err != nil {
		log.Fatalf("There is an error when get a migration driver: %v", err)
		os.Exit(1)
	}
	return migrate.NewWithDatabaseInstance(
		"file://"+migration.migrationConfig.Path,
		migration.driver.GetDataBaseNameMigration(), driverMigration,
	)
}

func (migration *MigrationComponent) Up() error {
	initMigration, err := migration.initDatabaseInstance()
	if err != nil {
		log.Fatalf("There is an error when configure migration structure: %v", err)
		os.Exit(1)
	}
	return initMigration.Up()
}

func (migration *MigrationComponent) Down() error {
	initMigration, err := migration.initDatabaseInstance()
	if err != nil {
		log.Fatalf("There is an error when configure migration structure: %v", err)
		os.Exit(1)
	}
	return initMigration.Down()
}

type DbComponent struct {
	sqlDb     *sqlx.DB
	sqlDriver *SqlDriver
}

func (component *DbComponent) GetSqlDb() *sqlx.DB {
	return component.sqlDb
}

func (component *DbComponent) Migrate(config *config.MigrationConfig) {
	migrationComponent := &MigrationComponent{
		driver:          *component.sqlDriver,
		migrationConfig: config,
	}
	err := migrationComponent.Up()

	if err != nil {
		log.Errorf("There is an error when make UP migrations: %v", err)
	}
}

func NewDb(config *config.DatabaseConfig) *DbComponent {
	var sqlDriver SqlDriver

	//	if config.Type == "mysql" {
	sqlDriver = driver.NewMysqlDriverDriver(config)
	//	} else {
	//		sqlDriver = driver.NewSqLite3Driver(config)
	//	}

	return &DbComponent{
		sqlDb:     sqlDriver.GetSqlDb(),
		sqlDriver: &sqlDriver,
	}
}
