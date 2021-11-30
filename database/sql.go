package database

import (
       "context"
       "fmt"
       "net/url"
       "database/sql"
)

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
     
     q := "SELECT * FROM accession_numbers WHERE organization_id = ? AND accession_number = ?"
       
}

func (db *SQLAccessionNumberDatabase) GetByObjectId(ctx context.Context, org string, id string) (*AccessionNumber, error) {

     org_id, err := db.getOrganizationId(ctx, org)

     if err != nil {
     	return nil, fmt.Errorf("Failed to resolve organization ID, %w", err)
     }

     q := "SELECT * FROM accession_numbers WHERE organization_id = ? AND object_id = ?"
}

func (db *SQLAccessionNumberDatabase) AddAccessionNumber(ctx context.Context, acc *AccessionNumber) error {

     org_id, err := db.getOrganizationId(ctx, org)

     if err != nil {

     	org_id, err = db.createOrganization(ctx, acc.OrganizationURI)

	if err != nil {
	       return fmt.Errorf("Failed to create organization, %w", err)
	}
     }	

}

func (db *SQLAccessionNumberDatabase) RemoveAccessionNumber(ctx context.Context, acc *AccessionNumber) error {

     org_id, err := db.getOrganizationId(ctx, acc.OrganizationURI)

     if err != nil {
     	return fmt.Errorf("Failed to resolve organization ID, %w", err)
     }

     q := "DELETE FROM accession_numbers WHERE organization_id = ? AND accession_number = ?"
}

func (db *SQLAccessionNumberDatabase) getOrganizationId(ctx context.Context, org string) (int64, error){

}

func (db *SQLAccessionNumberDatabase) createOrganizationId(ctx context.Context, org string) (int64, error){

}