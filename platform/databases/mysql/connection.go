package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/EstebanGitPro/motogo-backend/config"
	_ "github.com/go-sql-driver/mysql"
)

func GetDB(dbConfig config.Database) (*sql.DB, error) {
	var dsn string

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	if dbConfig.SSL != "" {
		dsn += "&tls=" + dbConfig.SSL
	}

	db, err := sql.Open(dbConfig.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("error to connect to database: %w", err)
	}

	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime))
	db.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTime))
	

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}
