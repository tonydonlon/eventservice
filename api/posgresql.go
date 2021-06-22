package api

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// PostgreSQLWriter implements EventWriter by writing to a PostGreSQL database
type PostgreSQLWriter struct {
	conn *sqlx.DB
}

// Init is a DB connection initialization
func (p *PostgreSQLWriter) Init() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PATH"),
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return err
	}

	// TODO pool connections
	db.SetMaxOpenConns(5)
	p.conn = db

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}

// Write inserts session events to a PostgreSQL database
func (p *PostgreSQLWriter) Write(msg SessionEvent, sessionID string) error {
	log.WithFields(logrus.Fields{
		"name":  "PostgreSQLWriter",
		"event": msg,
	}).Info("writing event")

	tx := p.conn.MustBegin()
	// TODO use NamedQuery
	switch msg.Type {
	case StartSession:
		tx.MustExec("INSERT INTO sessions (session_id, session_start) VALUES ($1, to_timestamp($2 / 1000.0))", msg.SessionID, msg.Timestamp)
	case EndSession:
		tx.MustExec("UPDATE sessions SET session_end=to_timestamp($1 / 1000.0) WHERE session_id=$2", msg.Timestamp, msg.SessionID)
	default:
		sql := `
		INSERT INTO events (event_timestamp, session_id, event_name_id)
		SELECT to_timestamp($1 / 1000.0), $2, event_name_id 
			FROM event_names en WHERE en.event_name=$3
		`
		tx.MustExec(sql, msg.Timestamp, sessionID, msg.Name)
	}
	// TODO note the queries convert typical JS unix time (on ms scale) from `Date.now()`. Should do some Date validation before sending to DB
	return tx.Commit()
}

// PostgreSQLReader implements SessionReader
type PostgreSQLReader struct {
	conn *sqlx.DB
}

// Init is a DB connection initialization
func (p *PostgreSQLReader) Init() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PATH"),
	)
	db, err := sqlx.Connect("postgres", dsn)
	// TODO pool connections
	db.SetMaxOpenConns(5)
	p.conn = db
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}

// SessionEvents retrieves a Session and all of it's set of events
func (p *PostgreSQLReader) SessionEvents(sessionID string) (*SessionEventsResponse, error) {
	sessionEvents := &SessionEventsResponse{
		Type:     TypeSession,
		Start:    0,
		End:      0,
		Children: []Event{},
	}

	sql := `
	SELECT
	    s.session_start,
	    s.session_end,
	    en.event_name,
	    e.event_timestamp
    FROM sessions s
        JOIN events e ON e.session_id = s.session_id
        JOIN event_names en ON en.event_name_id = e.event_name_id
    WHERE s.session_id = $1
    ORDER by e.event_timestamp`

	rows, err := p.conn.Queryx(sql, sessionID)
	if err != nil {
		return nil, err
	}

	type queryResult struct {
		Session_start   time.Time
		Session_end     time.Time
		Event_name      string
		Event_timestamp time.Time
	}
	rowCount := 0

	for rows.Next() {
		var q queryResult
		err = rows.StructScan(&q)
		if err != nil {
			log.Error(err)
		}

		sessionEvents.Start = q.Session_start.UnixNano() / 1000000
		sessionEvents.End = q.Session_end.UnixNano() / 1000000
		sessionEvents.Children = append(sessionEvents.Children, Event{
			Type:      TypeEvent,
			Timestamp: q.Event_timestamp.UnixNano() / 1000000,
			Name:      q.Event_name,
		})
		rowCount++
	}

	// handle no results
	if rowCount == 0 {
		return nil, nil
	}

	return sessionEvents, nil
}
