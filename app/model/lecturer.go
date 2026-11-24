package model

type Lecturer struct {
	ID             int    `json:"id"`
	UserID         int    `json:"user_id"`
	NIDN           string `json:"nidn"`
	StudyProgramID int    `json:"study_program_id"`
}
