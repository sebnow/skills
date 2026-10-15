package server

import (
	"context"
	"testing"
)

func TestProjectTasksGroupsByState(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	p := s.mustProject(t, "web")
	s.mustTask(t, p, "paused")
	s.mustTask(t, p, "running")
	s.mustTask(t, p, "finished")
	s.mustTask(t, p, "failed")
	dismissed := s.mustTask(t, p, "stopped")
	s.mustDismiss(t, dismissed)

	groups, err := s.projectTasks(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range []int{1, 1, 1, 2} {
		if got := len(groups[i]); got != want {
			t.Errorf("group %q has %d tasks, want %d", taskGroups[i].title, got, want)
		}
	}
}

func TestDeleteShownOnlyWithoutTasks(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	empty := s.mustProject(t, "empty")
	used := s.mustProject(t, "used")
	s.mustTask(t, used, "finished")

	for _, tc := range []struct {
		project string
		inUse   bool
	}{{empty, false}, {used, true}} {
		got, err := s.projectInUse(ctx, tc.project)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.inUse {
			t.Errorf("projectInUse(%s) = %v, want %v", tc.project, got, tc.inUse)
		}
	}
}
