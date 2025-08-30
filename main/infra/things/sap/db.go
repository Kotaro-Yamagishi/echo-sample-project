package sap

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"strings"
	"sync/atomic"
	"time"
)

type DB struct {
	connections []*sql.DB
	counter     uint64
}

func Open(driverName, dataSourceNames string) (*DB, error) {
	connections := strings.Split(dataSourceNames, ";")
	db := &DB{connections: make([]*sql.DB, len(connections))}

	err := scatter(len(db.connections), func(i int) (err error) {
		db.connections[i], err = sql.Open(driverName, connections[i])
		return err
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Close() error {
	return scatter(len(db.connections), func(i int) error {
		return db.connections[i].Close()
	})
}

func (db *DB) Driver() driver.Driver {
	return db.Primary().Driver()
}

func (db *DB) Begin() (*sql.Tx, error) {
	panic("not supported")
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	db.setRecordsModified(ctx)
	return db.Primary().BeginTx(ctx, opts)
}

func (db *DB) Exec(query string, args ...any) (sql.Result, error) {
	panic("not supported")
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	db.setRecordsModified(ctx)
	return db.Primary().ExecContext(ctx, query, args...)
}

func (db *DB) Ping() error {
	return scatter(len(db.connections), func(i int) error {
		return db.connections[i].Ping()
	})
}

func (db *DB) PingContext(ctx context.Context) error {
	return scatter(len(db.connections), func(i int) error {
		return db.connections[i].PingContext(ctx)
	})
}

func (db *DB) Prepare(query string) (Stmt, error) {
	panic("not supported")
}

func (db *DB) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	if db.recordsModified(ctx) {
		return db.Primary().PrepareContext(ctx, query)
	}

	stmts := make([]*sql.Stmt, len(db.connections))

	err := scatter(len(db.connections), func(i int) (err error) {
		stmts[i], err = db.connections[i].PrepareContext(ctx, query)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &stmt{
		ctx:   ctx,
		db:    db,
		stmts: stmts,
	}, nil
}

func (db *DB) Query(query string, args ...any) (*sql.Rows, error) {
	panic("not supported")
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if db.recordsModified(ctx) {
		return db.Primary().QueryContext(ctx, query, args...)
	}
	return db.ReadReplica().QueryContext(ctx, query, args...)
}

func (db *DB) QueryRow(query string, args ...any) *sql.Row {
	panic("not supported")
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if db.recordsModified(ctx) {
		return db.Primary().QueryRowContext(ctx, query, args...)
	}
	return db.ReadReplica().QueryRowContext(ctx, query, args...)
}

func (db *DB) SetMaxIdleConns(n int) {
	for i := range db.connections {
		db.connections[i].SetMaxIdleConns(n)
	}
}

func (db *DB) SetMaxOpenConns(n int) {
	for i := range db.connections {
		db.connections[i].SetMaxOpenConns(n)
	}
}

func (db *DB) SetConnMaxIdleTime(d time.Duration) {
	for i := range db.connections {
		db.connections[i].SetConnMaxIdleTime(d)
	}
}

func (db *DB) SetConnMaxLifetime(d time.Duration) {
	for i := range db.connections {
		db.connections[i].SetConnMaxLifetime(d)
	}
}

func (db *DB) ReadReplica() *sql.DB {
	return db.connections[db.rotate(len(db.connections))]
}

func (db *DB) Primary() *sql.DB {
	return db.connections[0]
}

func (db *DB) rotate(n int) int {
	if n <= 1 {
		return 0
	}
	return int(1 + (atomic.AddUint64(&db.counter, 1) % uint64(n-1)))
}

func (db *DB) setRecordsModified(ctx context.Context) {
	c := FromContext(ctx)
	if c != nil {
		c.SetRecordsModified()
	}
}

func (db *DB) recordsModified(ctx context.Context) bool {
	c := FromContext(ctx)
	if c != nil {
		return c.RecordsModified()
	}
	return false
}
