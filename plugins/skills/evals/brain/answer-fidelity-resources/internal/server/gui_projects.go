package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// taskGroups is the order the project overview shows task states in.
var taskGroups = []struct {
	title  string
	states []string
}{
	{"Needs attention", []string{"paused"}},
	{"Running", []string{"queued", "running"}},
	{"Finished", []string{"finished"}},
	{"Ended", []string{"stopped", "failed"}},
}

type taskRow struct {
	ID        string
	State     string
	Command   string
	Dismissed bool
}

// projectTasks lists every task in the project, grouped by taskGroups.
// Dismissed tasks stay in their group.
func (s *Store) projectTasks(ctx context.Context, project string) ([][]taskRow, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, state, command, dismissed IS NOT NULL FROM tasks WHERE project = ? ORDER BY created DESC`, project)
	if err != nil {
		return nil, fmt.Errorf("list tasks of project %q: %w", project, err)
	}
	defer rows.Close()
	groups := make([][]taskRow, len(taskGroups))
	for rows.Next() {
		var r taskRow
		if err := rows.Scan(&r.ID, &r.State, &r.Command, &r.Dismissed); err != nil {
			return nil, fmt.Errorf("list tasks of project %q: %w", project, err)
		}
		for i, g := range taskGroups {
			for _, st := range g.states {
				if st == r.State {
					groups[i] = append(groups[i], r)
				}
			}
		}
	}
	return groups, rows.Err()
}

// postDeleteProject handles the Delete button. The button is shown only
// while no task belongs to the project; the 409 covers a stale page.
func (h *handler) postDeleteProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("project")
	err := h.store.deleteProject(r.Context(), id)
	switch {
	case errors.Is(err, errProjectInUse):
		http.Error(w, "Tasks still belong to this project.", http.StatusConflict)
	case errors.Is(err, errUnknownProject):
		http.NotFound(w, r)
	case err != nil:
		h.fail(w, err)
	default:
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
	}
}
