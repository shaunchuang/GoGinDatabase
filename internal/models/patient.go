package models

import "time"

// Patient 定義患者資料模型
type Patient struct {
	ID                  int64     `json:"id"`
	Address             string    `json:"address"`
	Age                 int       `json:"age"`
	Birth               time.Time `json:"birth"`
	City                string    `json:"city"`
	DiseaseID           int64     `json:"disease_id"`
	District            string    `json:"district"`
	EmergencyContact    string    `json:"emergency_contact"`
	EmergencyPhone      string    `json:"emergency_phone"`
	EmergencyRelation   string    `json:"emergency_relation"`
	Gender              string    `json:"gender"`
	IDNo                string    `json:"idno"`
	Mail                string    `json:"mail"`
	Name                string    `json:"name"`
	OtherHistoryDisease string    `json:"other_history_disease"`
	OtherMedicalHistory string    `json:"other_medical_history"`
	Phone               string    `json:"phone"`
	HistoryDiseases     []string  `json:"history_diseases,omitempty"`
	MedicalHistories    []string  `json:"medical_histories,omitempty"`
}

// PatientHistoryDisease 定義患者病史資料模型
type PatientHistoryDisease struct {
	PatientID      int64  `json:"patient_id"`
	HistoryDisease string `json:"history_disease"`
	DiseaseID      int64  `json:"disease_id"`
}

// HistoryDisease 定義疾病史資料模型
type HistoryDisease struct {
	ID          int64  `json:"id"`
	DiseaseName string `json:"disease_name"`
}

// PatientMedicalHistory 定義患者醫療史資料模型
type PatientMedicalHistory struct {
	PatientID      int64  `json:"patient_id"`
	MedicalHistory string `json:"medical_history"`
}
