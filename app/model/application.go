package model

import "time"

// Status lamaran. Alur: submitted -> verified -> interview -> accepted/rejected.
const (
	ApplicationSubmitted = "submitted"
	ApplicationVerified  = "verified"
	ApplicationInterview = "interview"
	ApplicationAccepted  = "accepted"
	ApplicationRejected  = "rejected"
)

type Application struct {
	ID         int       `json:"id"`
	VacancyID  int       `json:"vacancy_id"`
	StudentID  int       `json:"student_id"`
	Motivation string    `json:"motivation"`
	Status     string    `json:"status"`
	AppliedAt  time.Time `json:"applied_at"`
}
