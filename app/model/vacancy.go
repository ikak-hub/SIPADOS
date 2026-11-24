package model

type Vacancy struct {
	ID          int     `json:"id"`
	PeriodID    int     `json:"period_id"`
	CourseID    int     `json:"course_id"`
	LecturerID  int     `json:"lecturer_id"`
	Quota       int     `json:"quota"`
	MinIPK      float64 `json:"min_ipk"`
	Description string  `json:"description"`
}
