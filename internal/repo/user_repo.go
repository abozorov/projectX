package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/abozorov/projectX/internal/models"
	"github.com/abozorov/projectX/pkg/errs"
)

type UIUserRepo interface {
	GetAll(ctx context.Context) ([]models.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	Create(ctx context.Context, u models.User) error
	Update(ctx context.Context, u models.User) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewpostgresRepo(db *sql.DB) UIUserRepo {
	return &postgresRepo{
		db: db,
	}
}

func execAnalysis(res sql.Result, err error) error {
	if err != nil {
		return fmt.Errorf("transaction.Exec(userQuerry): %w", distributor(err))
	}
	if rows, err := res.RowsAffected(); err != nil || rows == 0 {
		if rows == 0 {
			return fmt.Errorf("transaction.Exec: %w", errs.ErrUserNotFound)
		}
		return fmt.Errorf("transaction.Exec: %w", distributor(err))
	}
	return nil
}

func (r *postgresRepo) GetAll(ctx context.Context) ([]models.User, error) {
	// // sleep
	// time.Sleep(time.Second * 15)

	// get users
	const query = `
		SELECT  u.id,
				u.name,
				a.login,
				u.created_at
		FROM users u
		JOIN auth a on u.id = a.user_id
		WHERE u.is_active`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", distributor(err))
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() { // 5 = 1,2,3,4,5
		var user models.User
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Login,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		users = append(users, user)
	}

	// return data
	return users, nil
}

func (r *postgresRepo) Create(ctx context.Context, usr models.User) error {
	// check for exist
	const userQuery = `
	INSERT INTO users(name)
    VALUES ($1)
    RETURNING id`

	const authQuery = `
	INSERT INTO auth(user_id, login, password)
	VALUES ($1, $2, $3)`

	// add user
	transaction, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Create r.db.Begin: %w", distributor(err))
	}
	defer transaction.Rollback()

	var id int
	if err = transaction.QueryRow(userQuery, usr.Name).Scan(&id); err != nil {
		return fmt.Errorf("Create transaction.QueryRow.Scan: %w", distributor(err))
	}
	// log.Println("id=", id)
	// log.Println("user=", usr)

	if _, err = transaction.Exec(authQuery, id, usr.Login, usr.Password); err != nil {
		return fmt.Errorf("repo.Create transaction.Exec(authQuery): %w", distributor(err))
	}
	transaction.Commit()

	return nil
}

func (r *postgresRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	// load all users
	const query = `
		SELECT  u.id,
				u.name,
				a.login,
				u.created_at
		FROM users u
		JOIN auth a on u.id = a.user_id
		WHERE u.id = $1 AND u.is_active`

	// get by id
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Login,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("rows.Scan: %w", distributor(err))
	}

	return user, nil
}

func (r *postgresRepo) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	// load user
	const query = `
		SELECT  u.id,
				u.name,
				a.login,
				a.password,
				u.created_at
		FROM users u
		JOIN auth a on u.id = a.user_id
		WHERE a.login = $1 AND u.is_active`

	// get by login
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&user.ID,
		&user.Name,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		// log.Println(distributor(err).Error())
		return nil, fmt.Errorf("rows.Scan: %w", distributor(err))
	}

	return user, nil
}

func (r *postgresRepo) Update(ctx context.Context, usr models.User) error {
	// update with id
	const userQuery = `
		UPDATE users
		SET name=$1
		WHERE id=$2 AND is_active=true`

	const authQuery = `
		UPDATE auth
		SET login=$1
		WHERE user_id=$2`

	transaction, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("r.db.Begin: %w", distributor(err))
	}
	defer transaction.Rollback()

	if err = execAnalysis(transaction.ExecContext(ctx, userQuery, usr.Name, usr.ID)); err != nil {
		return fmt.Errorf("repo.Crerate userQuery: %w", err)
	}

	if err = execAnalysis(transaction.ExecContext(ctx, authQuery, usr.Login, usr.ID)); err != nil {
		return fmt.Errorf("repo.Crerate authQuery: %w", err)
	}
	transaction.Commit()

	// write data
	return nil
}

func (r *postgresRepo) DeleteUser(ctx context.Context, id int) error {
	// update with id
	const query = `
	UPDATE users
	SET is_active=false
	WHERE id=$1 AND is_active`

	if err := execAnalysis(r.db.ExecContext(ctx, query, id)); err != nil {
		return fmt.Errorf("DeleteUsers.db.Exec %w", err)
	}

	// write data
	return nil
}
