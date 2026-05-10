package lcp

import (
	"context"
	"database/sql"
	"time"

	domain "github.com/Mehrbod2002/lcp/internal/domain/lcp"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresPublicationRepository struct {
	db *sql.DB
}

type postgresLicenseRepository struct {
	db *sql.DB
}

func OpenPostgres(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func NewPostgresPublicationRepository(db *sql.DB) PublicationRepository {
	return &postgresPublicationRepository{db: db}
}

func NewPostgresLicenseRepository(db *sql.DB) LicenseRepository {
	return &postgresLicenseRepository{db: db}
}

func (r *postgresPublicationRepository) Save(ctx context.Context, pub *domain.Publication) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO publications (id, title, file_path, encrypted_path, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			file_path = EXCLUDED.file_path,
			encrypted_path = EXCLUDED.encrypted_path
	`, pub.ID, pub.Title, pub.FilePath, pub.EncryptedPath, pub.CreatedAt)
	return err
}

func (r *postgresPublicationRepository) FindAll(ctx context.Context) ([]*domain.Publication, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, file_path, encrypted_path, created_at
		FROM publications
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pubs []*domain.Publication
	for rows.Next() {
		pub, err := scanPublication(rows)
		if err != nil {
			return nil, err
		}
		pubs = append(pubs, pub)
	}
	return pubs, rows.Err()
}

func (r *postgresPublicationRepository) FindByID(ctx context.Context, id string) (*domain.Publication, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, title, file_path, encrypted_path, created_at
		FROM publications
		WHERE id = $1
	`, id)
	pub, err := scanPublication(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return pub, err
}

func (r *postgresLicenseRepository) Save(ctx context.Context, license *domain.License) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO licenses (
			id, publication_id, user_id, passphrase, hint, publication_url,
			right_print, right_copy, start_date, end_date, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			passphrase = EXCLUDED.passphrase,
			hint = EXCLUDED.hint,
			publication_url = EXCLUDED.publication_url,
			right_print = EXCLUDED.right_print,
			right_copy = EXCLUDED.right_copy,
			start_date = EXCLUDED.start_date,
			end_date = EXCLUDED.end_date
	`, license.ID, license.PublicationID, license.UserID, license.Passphrase, license.Hint,
		license.PublicationURL, license.RightPrint, license.RightCopy, license.StartDate,
		license.EndDate, license.CreatedAt)
	return err
}

func (r *postgresLicenseRepository) FindByID(ctx context.Context, id string) (*domain.License, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, publication_id, user_id, passphrase, hint, publication_url,
			right_print, right_copy, start_date, end_date, created_at
		FROM licenses
		WHERE id = $1
	`, id)
	license, err := scanLicense(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return license, err
}

func (r *postgresLicenseRepository) FindByPublication(ctx context.Context, publicationID *string) ([]*domain.License, error) {
	query := `
		SELECT id, publication_id, user_id, passphrase, hint, publication_url,
			right_print, right_copy, start_date, end_date, created_at
		FROM licenses
	`
	args := []interface{}{}
	if publicationID != nil {
		query += " WHERE publication_id = $1"
		args = append(args, *publicationID)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var licenses []*domain.License
	for rows.Next() {
		license, err := scanLicense(rows)
		if err != nil {
			return nil, err
		}
		licenses = append(licenses, license)
	}
	return licenses, rows.Err()
}

func (r *postgresLicenseRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM licenses WHERE id = $1", id)
	return err
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanPublication(row rowScanner) (*domain.Publication, error) {
	pub := &domain.Publication{}
	err := row.Scan(&pub.ID, &pub.Title, &pub.FilePath, &pub.EncryptedPath, &pub.CreatedAt)
	return pub, err
}

func scanLicense(row rowScanner) (*domain.License, error) {
	license := &domain.License{}
	err := row.Scan(&license.ID, &license.PublicationID, &license.UserID, &license.Passphrase,
		&license.Hint, &license.PublicationURL, &license.RightPrint, &license.RightCopy,
		&license.StartDate, &license.EndDate, &license.CreatedAt)
	return license, err
}
