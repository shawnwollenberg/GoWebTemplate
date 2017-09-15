package db

import (
	"database/sql"
	_"github.com/denisenkom/go-mssqldb"
	"flag"
	"WebTemplate/model"
	"github.com/jmoiron/sqlx"
)

var debug = flag.Bool("debug", false, "enable debugging")
var password = flag.String("password", "xxx", "the database password")
//var port *int = flag.Int("port", 1433, "the database port")
var server = flag.String("server", "", "the database server")
var user = flag.String("user", "", "the database user")

type Config struct {
	ConnectString string
}

func InitDb(cfg Config) (*pgDb, error) {
	if dbConn, err := sqlx.Connect("postgres", cfg.ConnectString); err != nil {
		return nil, err
	} else {
		p := &pgDb{dbConn: dbConn}
		if err := p.dbConn.Ping(); err != nil {
			return nil, err
		}
		if err := p.createTablesIfNotExist(); err != nil {
			return nil, err
		}
		if err := p.prepareSqlStatements(); err != nil {
			return nil, err
		}
		return p, nil
	}
}

type pgDb struct {
	dbConn *sqlx.DB
	sqlSelectPeople *sqlx.Stmt
	sqlInsertPerson *sqlx.NamedStmt
	sqlSelectPerson *sql.Stmt
}

func (p *pgDb) createTablesIfNotExist() error {
	create_sql := `
		IF  NOT EXISTS (SELECT * FROM sys.objects
		WHERE object_id = OBJECT_ID(N'[dbo].[shawntestpeople]') AND type in (N'U'))
		BEGIN
		CREATE TABLE shawntestpeople(
		testid int NOT NULL PRIMARY KEY,
		first varchar(50) NOT NULL,
		last varchar(50) NOT NULL)
		END
    `
	if rows, err := p.dbConn.Query(create_sql); err != nil {
		return err
	} else {
		rows.Close()
	}
	return nil
}

func (p *pgDb) prepareSqlStatements() (err error) {

	if p.sqlSelectPeople, err = p.dbConn.Preparex(
		"SELECT id, first, last FROM people",
	); err != nil {
		return err
	}
	if p.sqlInsertPerson, err = p.dbConn.PrepareNamed(
		"INSERT INTO people (first, last) VALUES (:first, :last) " +
			"RETURNING id, first, last",
	); err != nil {
		return err
	}
	if p.sqlSelectPerson, err = p.dbConn.Prepare(
		"SELECT id, first, last FROM people WHERE id = $1",
	); err != nil {
		return err
	}

	return nil
}

func (p *pgDb) SelectPeople() ([]*model.Person, error) {
	people := make([]*model.Person, 0)
	if err := p.sqlSelectPeople.Select(&people); err != nil {
		return nil, err
	}
	return people, nil
}