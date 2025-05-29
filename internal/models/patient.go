package models

import "time"

// Patient 定義患者資料模型
type Patient struct {
	ID                  int64     `json:"id"`
	UserID              int64     `json:"user_id"`
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
	HistoryDiseases     []string  `json:"history_diseases,omitempty"`  // 用於顯示，實際儲存在 historyDisease 欄位
	MedicalHistories    []string  `json:"medical_histories,omitempty"` // 用於顯示，實際儲存在 medicalHistory 欄位
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

// HistoryDiseaseEnum 疾病史枚舉類型
type HistoryDiseaseEnum string

// MedicalHistoryEnum 醫療史枚舉類型
type MedicalHistoryEnum string

// 疾病史枚舉值
const (
	HistoryHypertension    HistoryDiseaseEnum = "hypertension"
	HistoryDiabetes        HistoryDiseaseEnum = "diabetes"
	HistoryHyperlipidemia  HistoryDiseaseEnum = "hyperlipidemia"
	HistoryHeartDisease    HistoryDiseaseEnum = "heartDisease"
	HistoryStroke          HistoryDiseaseEnum = "stroke"
	HistoryCancer          HistoryDiseaseEnum = "cancer"
	HistoryCopd            HistoryDiseaseEnum = "copd"
	HistoryAsthma          HistoryDiseaseEnum = "asthma"
	HistorySleepApnea      HistoryDiseaseEnum = "sleepApnea"
	HistoryOsteoporosis    HistoryDiseaseEnum = "osteoporosis"
	HistoryArthritis       HistoryDiseaseEnum = "arthritis"
	HistoryKidneyDisease   HistoryDiseaseEnum = "kidneyDisease"
	HistoryLiverDisease    HistoryDiseaseEnum = "liverDisease"
	HistoryDepression      HistoryDiseaseEnum = "depression"
	HistoryAnxiety         HistoryDiseaseEnum = "anxiety"
	HistoryBipolarDisorder HistoryDiseaseEnum = "bipolarDisorder"
	HistorySchizophrenia   HistoryDiseaseEnum = "schizophrenia"
	HistoryDementia        HistoryDiseaseEnum = "dementia"
	HistoryParkinsons      HistoryDiseaseEnum = "parkinsonsDisease"
	HistoryEpilepsy        HistoryDiseaseEnum = "epilepsy"
	HistoryMigraine        HistoryDiseaseEnum = "migraine"
	HistoryThyroidDisease  HistoryDiseaseEnum = "thyroidDisease"
	HistoryGout            HistoryDiseaseEnum = "gout"
	HistoryOther           HistoryDiseaseEnum = "other"
)

// 醫療史枚舉值
const (
	MedicalHypertension    MedicalHistoryEnum = "hypertension"
	MedicalDiabetes        MedicalHistoryEnum = "diabetes"
	MedicalHyperlipidemia  MedicalHistoryEnum = "hyperlipidemia"
	MedicalHeartDisease    MedicalHistoryEnum = "heartDisease"
	MedicalStroke          MedicalHistoryEnum = "stroke"
	MedicalCancer          MedicalHistoryEnum = "cancer"
	MedicalCopd            MedicalHistoryEnum = "copd"
	MedicalAsthma          MedicalHistoryEnum = "asthma"
	MedicalSleepApnea      MedicalHistoryEnum = "sleepApnea"
	MedicalOsteoporosis    MedicalHistoryEnum = "osteoporosis"
	MedicalArthritis       MedicalHistoryEnum = "arthritis"
	MedicalKidneyDisease   MedicalHistoryEnum = "kidneyDisease"
	MedicalLiverDisease    MedicalHistoryEnum = "liverDisease"
	MedicalDepression      MedicalHistoryEnum = "depression"
	MedicalAnxiety         MedicalHistoryEnum = "anxiety"
	MedicalBipolarDisorder MedicalHistoryEnum = "bipolarDisorder"
	MedicalSchizophrenia   MedicalHistoryEnum = "schizophrenia"
	MedicalDementia        MedicalHistoryEnum = "dementia"
	MedicalParkinsons      MedicalHistoryEnum = "parkinsonsDisease"
	MedicalEpilepsy        MedicalHistoryEnum = "epilepsy"
	MedicalMigraine        MedicalHistoryEnum = "migraine"
	MedicalThyroidDisease  MedicalHistoryEnum = "thyroidDisease"
	MedicalGout            MedicalHistoryEnum = "gout"
	MedicalOther           MedicalHistoryEnum = "other"
)
