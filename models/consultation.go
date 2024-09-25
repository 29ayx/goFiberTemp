package models

type AppointmentForm struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	PatientEmail string `json:"patient_email"`
	Time      string `json:"time"`
	Date 	  string `json:"date"`
	Status	string `json:"status"`
	PregMonth string `json:"preg_month"` //Pregnancy month
	Specialist string `json:"specialist"`
	PrefferedDoctorGender string `json:"preferred_doctor_gender"`
	PrefferedCommnuicationMode string `json:"preferred_communication_mode"`
	Message string `json:"message"`

}

type Consultation struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	DoctorID    uint   `json:"doctor_id"`
	MotherID    uint   `json:"mother_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Status      string `json:"status"`
}

type Response struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ConsultationID uint `json:"consultation_id"`
	DoctorID    uint   `json:"doctor_id"`
	MotherID    uint   `json:"mother_id"`
	Content     string `json:"content"`
}