package model

type Course struct {
	ID             int    `json:"id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	StudyProgramID int    `json:"study_program_id"`
}
