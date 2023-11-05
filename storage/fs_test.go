package storage

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/muzykantov/tgpt/chat"
)

func TestSaveAndLoadHistory(t *testing.T) {
	// Setup.
	ctx := context.Background()
	baseDir, err := os.MkdirTemp("", "test_histories")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(baseDir) // Clean up.

	fs := FS{BaseDir: baseDir}
	history := &chat.History{
		ID: chat.ID{
			User:  123,
			Chat:  456,
			Model: "test-model",
		},
		Log: []chat.Message{
			{User: "Hello, Assistant!", Assistant: "Hello, User!"},
		},
	}

	// Execute SaveHistory.
	err = fs.SaveHistory(ctx, history)
	if err != nil {
		t.Fatalf("SaveHistory failed: %s", err)
	}

	// Execute LoadHistory.
	loadedHistory, err := fs.LoadHistory(ctx, history.ID)
	if err != nil {
		t.Fatalf("LoadHistory failed: %s", err)
	}

	// Assert.
	if !reflect.DeepEqual(history, loadedHistory) {
		t.Errorf("Loaded history %+v does not match saved history %+v", loadedHistory, history)
	}
}

func TestSaveAndLoadStatistics(t *testing.T) {
	// Setup.
	ctx := context.Background()
	baseDir, err := os.MkdirTemp("", "test_statistics")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(baseDir) // Clean up.

	fs := FS{BaseDir: baseDir}
	statistics := &chat.Statistics{
		ID: chat.ID{
			User:  789,
			Chat:  1011,
			Model: "test-model",
		},
		LastMessage: 0.5,
		Daily:       5.0,
		Monthly:     map[time.Month]chat.Cost{time.January: 150.0},
		Total:       155.5,
	}

	// Execute SaveStatistics.
	err = fs.SaveStatistics(ctx, statistics)
	if err != nil {
		t.Fatalf("SaveStatistics failed: %s", err)
	}

	// Execute LoadStatistics.
	loadedStatistics, err := fs.LoadStatistics(ctx, statistics.ID)
	if err != nil {
		t.Fatalf("LoadStatistics failed: %s", err)
	}

	// Assert.
	if !reflect.DeepEqual(statistics, loadedStatistics) {
		t.Errorf(
			"Loaded statistics %+v does not match saved statistics %+v",
			loadedStatistics,
			statistics,
		)
	}
}
