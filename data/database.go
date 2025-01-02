package data

import (
	"context"
	"docker-go-test/models"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pgInstance *models.Postgres
	pgOnce     sync.Once
)

// CreateDBPool creates a new database connection pool.
func CreateDBPool(ctx context.Context, connString string) (*models.Postgres, error) {
	pgOnce.Do(func() {
		dbpool, err := pgxpool.New(ctx, connString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
			os.Exit(1)
		}

		pgInstance = &models.Postgres{Db: dbpool}
	})

	return pgInstance, nil
}

// QueryGreeting queries the greeting message from the database.
func QueryGreeting(ctx context.Context, pg *models.Postgres) (string, error) {
	var greeting string
	err := pg.Db.QueryRow(ctx, "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		return "", fmt.Errorf("QueryRow failed: %v", err)
	}
	return greeting, nil
}

// Test Single Query queries from the database.
func QuerySingleTest(ctx context.Context, pg *models.Postgres) (*models.ClientProfile, error) {
	var client models.ClientProfile
	err := pg.Db.QueryRow(ctx, "SELECT * FROM Account WHERE id = 1").Scan(&client.Id, &client.FirstName, &client.LastName, &client.Token)
	if err != nil {
		return nil, fmt.Errorf("QueryRow failed: %v", err)
	}
	return &client, nil
}

// Test Multiple Query queries from the database.
func QueryMultiTest(ctx context.Context, pg *models.Postgres) ([]models.ClientProfile, error) {
	rows, err := pg.Db.Query(ctx, "SELECT * FROM Account")
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.ClientProfile])
}

func GetUser(pg *models.Postgres, ctx context.Context, acctId int) (*models.ClientProfile, error) {
	query := `
	SELECT * FROM Account WHERE id = @acctId
	`
	var client models.ClientProfile
	// Define the named arguments for the query.
	args := pgx.NamedArgs{
		"acctId": acctId,
	}

	err := pg.Db.QueryRow(ctx, query, args).Scan(&client.Id, &client.FirstName, &client.LastName, &client.Token)
	if err != nil {
		return nil, fmt.Errorf("QueryRow failed: %v", err)
	}
	return &client, nil
}

func UpdateUser(pg *models.Postgres, ctx context.Context, acctId int, accountDetails models.ClientProfile) error {
	query := `
			UPDATE Account 
			SET id = @Id, FirstName = @FirstName, LastName = @LastName, Token = @Token
			WHERE id = @Id
			`
	args := pgx.NamedArgs{
		"Id":        accountDetails.Id,
		"FirstName": accountDetails.FirstName,
		"LastName":  accountDetails.LastName,
		"Token":     accountDetails.Token,
	}
	_, err := pg.Db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to update row: %w", err)
	}

	return nil
}

func InsertUser(pg *models.Postgres, ctx context.Context, accountDetails models.ClientProfile) error {
	query := `INSERT INTO Account (id, FirstName, LastName, Token) VALUES (@Id, @FirstName, @LastName, @Token)`
	args := pgx.NamedArgs{
		"Id":        accountDetails.Id,
		"FirstName": accountDetails.FirstName,
		"LastName":  accountDetails.LastName,
		"Token":     accountDetails.Token,
	}
	_, err := pg.Db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

func DeleteUser(pg *models.Postgres, ctx context.Context, acctId int) error {
	query := `
	DELETE FROM Account WHERE id = @Id
	`
	args := pgx.NamedArgs{
		"Id": acctId,
	}
	_, err := pg.Db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error deleting row: %w", err)
	}

	return nil
}
