package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	_ "github.com/go-sql-driver/mysql"
)

func GetDB(dbConfig config.Database, log logger.Logger) (*sql.DB, error) {
	log.Info("Conectando a base de datos MySQL",
		"host", dbConfig.Host,
		"port", dbConfig.Port,
		"database", dbConfig.Name,
		"driver", dbConfig.Driver)

	var dsn string

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	if dbConfig.SSL != "" {
		log.Debug("SSL habilitado para conexión a base de datos", "tls", dbConfig.SSL)
		dsn += "&tls=" + dbConfig.SSL
	}

	db, err := sql.Open(dbConfig.Driver, dsn)
	if err != nil {
		log.Error("Error abriendo conexión a base de datos",
			"error", err,
			"host", dbConfig.Host,
			"database", dbConfig.Name)
		return nil, fmt.Errorf("error to connect to database: %w", err)
	}

	log.Debug("Configurando pool de conexiones",
		"max_open_conns", dbConfig.MaxOpenConns,
		"max_idle_conns", dbConfig.MaxIdleConns,
		"conn_max_lifetime", dbConfig.ConnMaxLifetime,
		"conn_max_idle_time", dbConfig.ConnMaxIdleTime)

	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime))
	db.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTime))

	log.Info("Verificando conectividad con base de datos (ping)...")

	// Crear contexto con timeout para el ping (5 segundos)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Error("Error en ping a base de datos",
			"error", err,
			"host", dbConfig.Host,
			"database", dbConfig.Name)
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Success("Conexión a base de datos establecida exitosamente",
		"host", dbConfig.Host,
		"database", dbConfig.Name)

	return db, nil
}
