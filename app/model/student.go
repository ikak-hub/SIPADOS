package model

type Student struct {
	ID             int     `json:"id"`
	UserID         int     `json:"user_id"`
	NIM            string  `json:"nim"`
	StudyProgramID int     `json:"study_program_id"`
	Semester       int     `json:"semester"`
	IPK            float64 `json:"ipk"`
}
