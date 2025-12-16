package repository

import (
	"context"
	"testing"
	"time"

	"updatemanager/internal/models"
	"updatemanager/pkg/database"
)

var updateRolloutRepo *UpdateRolloutRepository

func setupUpdateRolloutTestDB(t *testing.T) {
	ctx := context.Background()

	cfg := &database.Config{
		URI:      "mongodb://admin:admin123@localhost:27017/updatemanager_test?authSource=admin",
		Database: "updatemanager_test",
		Timeout:  10 * time.Second,
	}

	db, err := database.Connect(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	testDB = db
	testCtx = ctx
	updateRolloutRepo = NewUpdateRolloutRepository(db.Collection("update_rollouts"))
}

func teardownUpdateRolloutTestDB(t *testing.T) {
	if testDB != nil {
		_ = testDB.Collection("update_rollouts").Drop(testCtx)
		_ = testDB.Disconnect(testCtx)
	}
}

func TestUpdateRolloutCreate(t *testing.T) {
	setupUpdateRolloutTestDB(t)
	defer teardownUpdateRolloutTestDB(t)

	rollout := &models.UpdateRollout{
		EndpointID:  "endpoint-123",
		ProductID:   "test-product",
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		Status:      models.RolloutStatusPending,
		InitiatedBy: "admin-user",
		Progress:    0,
	}

	err := updateRolloutRepo.Create(testCtx, rollout)
	if err != nil {
		t.Fatalf("Failed to create update rollout: %v", err)
	}

	if rollout.ID.IsZero() {
		t.Error("Update rollout ID was not set after creation")
	}

	t.Logf("Created update rollout with ID: %s", rollout.ID.Hex())
}

func TestUpdateRolloutUpdateStatus(t *testing.T) {
	setupUpdateRolloutTestDB(t)
	defer teardownUpdateRolloutTestDB(t)

	rollout := &models.UpdateRollout{
		EndpointID:  "endpoint-456",
		ProductID:   "test-product",
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		Status:      models.RolloutStatusPending,
		InitiatedBy: "admin-user",
		Progress:    0,
	}

	err := updateRolloutRepo.Create(testCtx, rollout)
	if err != nil {
		t.Fatalf("Failed to create update rollout: %v", err)
	}

	err = updateRolloutRepo.UpdateStatus(testCtx, rollout.ID, models.RolloutStatusInProgress, "")
	if err != nil {
		t.Fatalf("Failed to update rollout status: %v", err)
	}

	retrieved, err := updateRolloutRepo.GetByID(testCtx, rollout.ID)
	if err != nil {
		t.Fatalf("Failed to get update rollout: %v", err)
	}

	if retrieved.Status != models.RolloutStatusInProgress {
		t.Errorf("Status not updated: got %s, want %s", retrieved.Status, models.RolloutStatusInProgress)
	}

	if retrieved.StartedAt == nil {
		t.Error("StartedAt should be set when status is in_progress")
	}

	t.Logf("Updated rollout status: %+v", retrieved)
}

func TestUpdateRolloutUpdateProgress(t *testing.T) {
	setupUpdateRolloutTestDB(t)
	defer teardownUpdateRolloutTestDB(t)

	rollout := &models.UpdateRollout{
		EndpointID:  "endpoint-789",
		ProductID:   "test-product",
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		Status:      models.RolloutStatusInProgress,
		InitiatedBy: "admin-user",
		Progress:    0,
	}

	err := updateRolloutRepo.Create(testCtx, rollout)
	if err != nil {
		t.Fatalf("Failed to create update rollout: %v", err)
	}

	err = updateRolloutRepo.UpdateProgress(testCtx, rollout.ID, 50)
	if err != nil {
		t.Fatalf("Failed to update rollout progress: %v", err)
	}

	retrieved, err := updateRolloutRepo.GetByID(testCtx, rollout.ID)
	if err != nil {
		t.Fatalf("Failed to get update rollout: %v", err)
	}

	if retrieved.Progress != 50 {
		t.Errorf("Progress not updated: got %d, want 50", retrieved.Progress)
	}

	t.Logf("Updated rollout progress: %+v", retrieved)
}
