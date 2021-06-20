package writers

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/tonydonlon/eventservice/api"
)

// PostgreSQLWriter implements EventWriter and writes to a PostGreSQL database
type PostgreSQLWriter struct {
	conn *sqlx.DB
}

// Init is a DB connection initialization
func (p *PostgreSQLWriter) Init() error {
	// TODO pool connections
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PATH"),
	)
	db, err := sqlx.Connect("postgres", dsn)
	p.conn = db
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}

// Write insert session events to a PostgreSQL database
func (p *PostgreSQLWriter) Write(msg api.Event, sessionID string) error {
	log.WithFields(logrus.Fields{
		"name":  "PostgreSQLWriter",
		"event": msg,
	}).Info("writing event")

	tx := p.conn.MustBegin()
	// TODO use NamedQuery
	switch msg.Type {
	case api.StartSession:
		tx.MustExec("INSERT INTO sessions (session_id, session_start) VALUES ($1, to_timestamp($2 / 1000.0))", msg.SessionID, msg.Time)
	case api.EndSession:
		tx.MustExec("UPDATE sessions SET session_end=to_timestamp($1 / 1000.0) WHERE session_id=$2", msg.Time, msg.SessionID)
	default:
		sql := `
		INSERT INTO events (event_timestamp, session_id, event_name_id)
		SELECT to_timestamp($1 / 1000.0), $2, event_name_id 
			FROM event_names en WHERE en.event_name=$3
		`
		tx.MustExec(sql, msg.Time, sessionID, msg.Name)
	}
	tx.Commit()

	return nil
}
