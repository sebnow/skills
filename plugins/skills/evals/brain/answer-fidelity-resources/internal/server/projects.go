package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	errUnknownProject = errors.New("unknown project")
	errProjectInUse   = errors.New("project in use")
)

// Project is a set of defaults that tasks started in it take.
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Pool        string `json:"pool"`
	Branch      string `json:"branch"`
	Environment string `json:"environment"`
}

// projectInUse reports whether any task belongs to the project.
func (s *Store) projectInUse(ctx context.Context, id string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT count(*) FROM tasks WHERE project = ?`, id).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("count tasks of project %q: %w", id, err)
	}
	return n > 0, nil
}

// deleteProject removes a project. It refuses while a task belongs to
// it, when it returns errProjectInUse, or it does not exist, when it
// returns errUnknownProject.
func (s *Store) deleteProject(ctx context.Context, id string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("delete project %q: %w", id, err)
	}
	defer tx.Rollback()
	var tasks int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM tasks WHERE project = ?`, id).Scan(&tasks); err != nil {
		return fmt.Errorf("delete project %q: %w", id, err)
	}
	if tasks > 0 {
		return fmt.Errorf("%w: %d tasks belong to project %q", errProjectInUse, tasks, id)
	}
	res, err := tx.ExecContext(ctx, `DELETE FROM projects WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete project %q: %w", id, err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("%w: %q", errUnknownProject, id)
	}
	return tx.Commit()
}

// readProject returns one project, or errUnknownProject.
func (s *Store) readProject(ctx context.Context, id string) (Project, error) {
	var p Project
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, pool, branch, environment FROM projects WHERE id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.Pool, &p.Branch, &p.Environment)
	if errors.Is(err, sql.ErrNoRows) {
		return Project{}, fmt.Errorf("%w: %q", errUnknownProject, id)
	}
	if err != nil {
		return Project{}, fmt.Errorf("read project %q: %w", id, err)
	}
	return p, nil
}
