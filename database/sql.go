package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"net/url"
)

//go:embed sql.schema
var schema string

type SQLAccessionNumberDatabase struct {
	AccessionNumberDatabase
	conn *sql.DB
}

func NewSQLAccessionNumberDatabase(ctx context.Context, uri string) (AccessionNumberDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	engine := u.Host
	dsn := q.Get("dsn")

	conn, err := sql.Open(engine, dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection, %w", err)
	}

	err = conn.Ping()

	if err != nil {
		return nil, fmt.Errorf("Failed to ping database server, %w", err)
	}

	db := &SQLAccessionNumberDatabase{
		conn: conn,
	}

	return db, nil
}

func (db *SQLAccessionNumberDatabase) GetByAccessionNumber(ctx context.Context, org string, num string) (*AccessionNumber, error) {

	org_id, err := db.getOrganizationId(ctx, org)

	if err != nil {
		return nil, fmt.Errorf("Failed to resolve organization ID, %w", err)
	}

	q := "SELECT accession_number, object_id FROM accession_numbers WHERE organization_id = ? AND accession_number = ?"

	row := db.conn.QueryRowContext(ctx, q, org_id, num)

	var accession_number string
	var object_id string

	err = row.Scan(&accession_number, &object_id)

	if err != nil {
		return nil, fmt.Errorf("Failed to scan database row, %w", err)
	}

	a := &AccessionNumber{
		AccessionNumber: accession_number,
		ObjectId:        object_id,
		OrganizationURI: org,
	}

	return a, nil
}

func (db *SQLAccessionNumberDatabase) GetByObjectId(ctx context.Context, org string, id string) (*AccessionNumber, error) {

	org_id, err := db.getOrganizationId(ctx, org)

	if err != nil {
		return nil, fmt.Errorf("Failed to resolve organization ID, %w", err)
	}

	q := "SELECT * FROM accession_numbers WHERE organization_id = ? AND object_id = ?"

	row := db.conn.QueryRowContext(ctx, q, org_id, id)

	var accession_number string
	var object_id string

	err = row.Scan(&accession_number, &object_id)

	if err != nil {
		return nil, fmt.Errorf("Failed to scan database row, %w", err)
	}

	a := &AccessionNumber{
		AccessionNumber: accession_number,
		ObjectId:        object_id,
		OrganizationURI: org,
	}

	return a, nil

}

func (db *SQLAccessionNumberDatabase) AddAccessionNumber(ctx context.Context, acc *AccessionNumber) error {

	org_id, err := db.getOrganizationId(ctx, acc.OrganizationURI)

	if err != nil {

		org_id, err = db.createOrganization(ctx, acc.OrganizationURI)

		if err != nil {
			return fmt.Errorf("Failed to create organization, %w", err)
		}
	}

	q := "INSERT INTO accession_numbers (accession_number, object_id, organization_id) VALUES(?, ?, ?)"

	_, err = db.conn.ExecContext(ctx, q, org_id, acc.AccessionNumber)

	if err != nil {
		return fmt.Errorf("Failed to create accession number, %w", err)
	}

	return nil
}

func (db *SQLAccessionNumberDatabase) RemoveAccessionNumber(ctx context.Context, acc *AccessionNumber) error {

	org_id, err := db.getOrganizationId(ctx, acc.OrganizationURI)

	if err != nil {
		return fmt.Errorf("Failed to resolve organization ID, %w", err)
	}

	q := "DELETE FROM accession_numbers WHERE organization_id = ? AND accession_number = ?"

	_, err = db.conn.ExecContext(ctx, q, org_id, acc.AccessionNumber)

	if err != nil {
		return fmt.Errorf("Failed to delete row, %w", err)
	}

	return nil
}

func (db *SQLAccessionNumberDatabase) getOrganizationId(ctx context.Context, org string) (int64, error) {

	q := "SELECT organization_id FROM organizations WHERE organization_uri = ?"
	row := db.conn.QueryRowContext(ctx, q, org)

	var organization_id int64
	err := row.Scan(&organization_id)

	if err != nil {
		return 0, fmt.Errorf("Failed to scan database row, %w", err)
	}

	return organization_id, nil
}

func (db *SQLAccessionNumberDatabase) createOrganization(ctx context.Context, org string) (int64, error) {

	org_id := NewOrganizationId()

	q := "INSERT INTO organization (organization_uri, organization_id) VALUES(?, ?)"

	_, err := db.conn.ExecContext(ctx, q, org, org_id)

	if err != nil {
		return 0, fmt.Errorf("Failed to create organization, %w", err)
	}

	return org_id, nil
}
