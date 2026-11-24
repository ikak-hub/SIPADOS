package model

import "time"

// Status periode rekrutmen. Alur: draft -> open -> closed.
const (
	PeriodDraft  = "draft"
	PeriodOpen   = "open"
	PeriodClosed = "closed"
)

// RecruitmentPeriod adalah entitas periode rekrutmen asdos.
type RecruitmentPeriod struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatePeriodRequest dipakai untuk POST. Tanggal berformat YYYY-MM-DD.
// Status tidak dikirim klien; periode baru selalu berstatus draft.
type CreatePeriodRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// ReplacePeriodRequest dipakai untuk PUT; seluruh field wajib ada.
type ReplacePeriodRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Status    string `json:"status"`
}

// PatchPeriodRequest dipakai untuk PATCH; field nil berarti tidak diubah.
type PatchPeriodRequest struct {
	Name      *string `json:"name"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	Status    *string `json:"status"`
}
