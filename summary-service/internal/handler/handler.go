package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/adi290491/productivity-planner/summary-service/internal/service"
	"github.com/adi290491/productivity-planner/summary-service/pkg/httperr"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) GetDailySummary(w http.ResponseWriter, r *http.Request) {

	userID := strings.TrimSpace(r.Header.Get("X-USER-ID"))
	if userID == "" {
		slog.Warn("Daily summary request missing user ID",
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
		httperr.HandleError(w, http.StatusUnauthorized, "user id is missing")
		return
	}

	queryDate := r.URL.Query().Get("date")
	slog.Info("Processing daily summary request",
		"user_id", userID,
		"date", queryDate,
	)

	summaryResponse, err := h.svc.GetDailySessionSummary(r.Context(), userID, queryDate)

	if err != nil {
		if strings.Contains(err.Error(), "invalid date format") {
			slog.Warn("Invalid date format in request",
				"user_id", userID,
				"date", queryDate,
				"error", err,
			)
			httperr.HandleError(w, http.StatusBadRequest, err.Error())
			return
		}

		if strings.Contains(err.Error(), "no sessions found for the given day") {
			httperr.HandleError(w, http.StatusNoContent, fmt.Sprintf("no sessions found for user: %s on date: %s", userID, queryDate))
			return
		}

		slog.Error("Failed to get daily summary",
			"user_id", userID,
			"date", queryDate,
			"error", err,
		)
		httperr.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	slog.Info("Daily summary retrieved successfully",
		"user_id", userID,
		"date", summaryResponse.Date,
		"total_time", summaryResponse.TotalTime,
	)

	httperr.RespondJSON(w, http.StatusOK, summaryResponse)

}

func (h *Handler) GetWeeklySummary(w http.ResponseWriter, r *http.Request) {

	userID := strings.TrimSpace(r.Header.Get("X-USER-ID"))
	if userID == "" {
		slog.Warn("Daily summary request missing user ID",
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
		httperr.HandleError(w, http.StatusUnauthorized, "user id is missing")
		return
	}

	startDate := r.URL.Query().Get("start_date")

	summary, err := h.svc.GetWeeklySessionSummary(r.Context(), userID, startDate)
	if err != nil {
		if strings.Contains(err.Error(), "invalid date format") {
			slog.Warn("Invalid date format in request",
				"user_id", userID,
				"start_date", startDate,
				"error", err,
			)
			httperr.HandleError(w, http.StatusBadRequest, err.Error())
			return
		}
		if strings.Contains(err.Error(), "no sessions found") {
			slog.Info("No sessions found for user",
				"user_id", userID,
				"start_date", startDate,
			)
			httperr.HandleError(w, http.StatusNoContent, fmt.Sprintf("no sessions found for user: %s", userID))
			return
		}
		slog.Error("Failed to get weekly summary",
			"user_id", userID,
			"start_date", startDate,
			"error", err,
		)
		httperr.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	slog.Info("Weekly summary retrieved successfully",
		"user_id", userID,
		"start_date", summary.StartDate,
		"end_date", summary.EndDate,
		"total_time", summary.TotalTime,
		"daily_summaries_count", len(summary.DailySummaries),
	)

	httperr.RespondJSON(w, http.StatusOK, summary)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	httperr.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"service":   "summary-service",
		"timestamp": time.Now(),
	})
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	httperr.RespondJSON(w, http.StatusOK, map[string]interface{}{"status": "ready"})
}
