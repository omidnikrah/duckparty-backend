package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/omidnikrah/duckparty-backend/internal/model"
	"gorm.io/gorm"
)

const (
	leaderboardJobInterval = 4 * time.Hour
	leaderboardJobTimeout  = 5 * time.Minute
)

func NewCron(ctx context.Context, db *gorm.DB, logger *slog.Logger) (gocron.Scheduler, error) {
	if db == nil {
		return nil, fmt.Errorf("db is required")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if logger == nil {
		logger = slog.Default()
	}

	leaderboardLogger := logger.With("scope", "cron", "job", "duck-leaderboard")

	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
		gocron.WithLimitConcurrentJobs(1, gocron.LimitModeWait),
	)
	if err != nil {
		return nil, fmt.Errorf("create scheduler: %w", err)
	}

	task := gocron.NewTask(func(jobCtx context.Context) {
		if ctx.Err() != nil {
			return
		}

		execCtx, cancel := context.WithTimeout(jobCtx, leaderboardJobTimeout)
		defer cancel()

		updated, err := updateDuckLeaderboard(execCtx, db)
		if err != nil {
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				leaderboardLogger.Error("failed to reconcile leaderboard", "error", err)
			}
			return
		}

		if updated > 0 {
			leaderboardLogger.Info("leaderboard synchronized", "rows", updated)
			return
		}

		leaderboardLogger.Debug("leaderboard already up to date")
	})

	if _, err := scheduler.NewJob(
		gocron.DurationJob(leaderboardJobInterval),
		task,
		gocron.WithName("duck-leaderboard"),
		gocron.WithSingletonMode(gocron.LimitModeWait),
	); err != nil {
		return nil, fmt.Errorf("schedule leaderboard job: %w", err)
	}

	scheduler.Start()

	leaderboardLogger.Info("scheduler started", "interval", leaderboardJobInterval.String())

	return scheduler, nil
}

func updateDuckLeaderboard(ctx context.Context, db *gorm.DB) (int64, error) {
	type duckRank struct {
		ID            uint
		LikesCount    int64
		DislikesCount int64
		Rank          uint
	}

	var ducks []duckRank
	if err := db.WithContext(ctx).
		Model(&model.Duck{}).
		Select("id", "likes_count", "dislikes_count", "rank").
		Order("likes_count DESC").
		Order("dislikes_count ASC").
		Order("id ASC").
		Find(&ducks).Error; err != nil {
		return 0, fmt.Errorf("fetch ducks for leaderboard: %w", err)
	}

	if len(ducks) == 0 {
		return 0, nil
	}

	var updated int64

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)

		for index, duck := range ducks {
			expectedRank := uint(index + 1)
			if duck.Rank == expectedRank {
				continue
			}

			if err := tx.Model(&model.Duck{}).
				Where("id = ?", duck.ID).
				Update("rank", expectedRank).Error; err != nil {
				return err
			}

			updated++
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("update duck ranks: %w", err)
	}

	return updated, nil
}
