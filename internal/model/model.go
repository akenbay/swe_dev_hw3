package model

type StudentResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
	GroupName string `json:"group_name"`
}

type ScheduleResponse struct {
	ID        int    `json:"id"`
	Faculty   string `json:"faculty"`
	Group     string `json:"group"`
	Subject   string `json:"subject"`
	ClassTime string `json:"class_time"`
}
