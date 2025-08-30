package sap

import (
	"context"
	"database/sql"
)

type Stmt interface {
	Close() error
	Exec(...any) (sql.Result, error)
	Query(...any) (*sql.Rows, error)
	QueryRow(...any) *sql.Row
}

type stmt struct {
	ctx   context.Context
	db    *DB
	stmts []*sql.Stmt
}

func (s *stmt) Close() error {
	return scatter(len(s.stmts), func(i int) error {
		return s.stmts[i].Close()
	})
}

func (s *stmt) Exec(args ...any) (sql.Result, error) {
	s.db.setRecordsModified(s.ctx)
	return s.stmts[0].Exec(args...)
}

func (s *stmt) Query(args ...any) (*sql.Rows, error) {
	if s.db.recordsModified(s.ctx) {
		return s.stmts[0].Query(args...)
	}
	return s.stmts[s.db.rotate(len(s.db.connections))].Query(args...)
}

func (s *stmt) QueryRow(args ...any) *sql.Row {
	if s.db.recordsModified(s.ctx) {
		return s.stmts[0].QueryRow(args...)
	}
	return s.stmts[s.db.rotate(len(s.db.connections))].QueryRow(args...)
}
