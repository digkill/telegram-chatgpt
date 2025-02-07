package driver

import (
	"github.com/digkill/telegram-chatgpt/internal/config"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type MysqlDriver struct {
	sqlDb *sqlx.DB
}

func (driver *MysqlDriver) GetSqlDb() *sqlx.DB {
	return driver.sqlDb
}

func (driver *MysqlDriver) GetDriverMigration(migrationTable string, databaseName string) (database.Driver, error) {
	return mysql.WithInstance(driver.sqlDb.DB, &mysql.Config{
		MigrationsTable: migrationTable,
	})
}

func (driver *MysqlDriver) GetDataBaseNameMigration() string {
	return "mysql"
}

func NewMysqlDriverDriver(config *config.DatabaseConfig) *MysqlDriver {
	connection, err := sqlx.Connect(config.Type, config.Username+":"+config.Password+"@tcp("+
		config.Host+":"+strconv.Itoa(config.Port)+")/"+config.Name+"?multiStatements=true")

	//	config.Host+":"+config.Port+")/"+config.Name+"?multiStatements=true")

	if err != nil {
		log.Fatalf("Can not open database: %v There is an error: %v", config.Name, err)
		os.Exit(1)
	}

	connection.SetMaxIdleConns(config.MaxIdleConns)
	connection.SetMaxOpenConns(config.MaxOpenConns)
	return &MysqlDriver{sqlDb: connection}
}
