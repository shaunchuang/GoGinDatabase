package handlers

import (
	"database/sql"
	"fmt"
	"golang-gin-app/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateFakePatientsFormHandler 處理 GET /fake-patients 路由，顯示生成假病患資料的表單
func GenerateFakePatientsFormHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "fake_patients.html", gin.H{
			"title": "產生假病患資料",
		})
	}
}

// GenerateFakePatientsHandler 處理 POST /fake-patients 路由，生成假病患資料並顯示
func GenerateFakePatientsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 獲取表單參數
		countStr := c.PostForm("count")
		count, err := strconv.Atoi(countStr)
		if err != nil || count <= 0 || count > 100 {
			c.HTML(http.StatusBadRequest, "fake_patients.html", gin.H{
				"title": "產生假病患資料",
				"error": "請輸入有效的數量（1-100）",
			})
			return
		}

		// 檢查是否要直接插入到資料庫
		insertToDBStr := c.PostForm("insertToDB")
		insertToDB := insertToDBStr == "true"

		// 生成假病患資料
		patients, err := utils.GenerateFakePatients(count)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "fake_patients.html", gin.H{
				"title": "產生假病患資料",
				"error": "生成假病患資料失敗: " + err.Error(),
			})
			return
		}

		// 如果需要插入到資料庫
		var successCount int
		var errorMessages []string
		var ids []int64

		if insertToDB {
			tx, err := db.Begin()
			if err != nil {
				c.HTML(http.StatusInternalServerError, "fake_patients.html", gin.H{
					"title":    "產生假病患資料",
					"error":    "資料庫交易開始失敗: " + err.Error(),
					"patients": patients,
				})
				return
			}
			defer tx.Rollback()

			for _, patient := range patients {
				// 插入主病患資料
				var patientID int64
				err := tx.QueryRow(`
					INSERT INTO patient (name, gender, idno, age, birth, address, city, district, 
						phone, mail, disease_id, emergency_contact, emergency_phone, emergency_relation,
						OTHERHISTORYDISEASE, OTHERMEDICALHISTORY)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
					RETURNING ID`,
					patient.Name, patient.Gender, patient.IDNo, patient.Age, patient.Birth,
					patient.Address, patient.City, patient.District, patient.Phone, patient.Mail,
					patient.DiseaseID, patient.EmergencyContact, patient.EmergencyPhone,
					patient.EmergencyRelation, patient.OtherHistoryDisease, patient.OtherMedicalHistory,
				).Scan(&patientID)

				if err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("插入病患 %s 失敗: %v", patient.Name, err))
					continue
				}

				// 更新 ID 以便在前端顯示
				patient.ID = patientID
				ids = append(ids, patientID) // 插入病史資料
				for _, historyDisease := range patient.HistoryDiseases {
					// 根據疾病名稱查詢對應的 ID
					var diseaseID int64
					err = tx.QueryRow(`
						SELECT ID FROM history_disease 
						WHERE disease_name = ?`,
						historyDisease,
					).Scan(&diseaseID)
					if err != nil {
						errorMessages = append(errorMessages, fmt.Sprintf("查詢疾病 %s 的ID失敗: %v", historyDisease, err))
						continue
					}

					_, err = tx.Exec(`
						INSERT INTO patient_history_disease (patient_id, history_disease, disease_id)
						VALUES (?, ?, ?)`,
						patientID, historyDisease, diseaseID,
					)
					if err != nil {
						errorMessages = append(errorMessages, fmt.Sprintf("插入病患 %s 的病史資料失敗: %v", patient.Name, err))
					}
				}

				// 插入醫療史資料
				for _, medicalHistory := range patient.MedicalHistories {
					_, err = tx.Exec(`
						INSERT INTO patient_medical_history (patient_id, medical_history)
						VALUES (?, ?)`,
						patientID, medicalHistory,
					)
					if err != nil {
						errorMessages = append(errorMessages, fmt.Sprintf("插入病患 %s 的醫療史資料失敗: %v", patient.Name, err))
					}
				}

				successCount++
			}

			// 如果全部成功，提交事務
			if len(errorMessages) == 0 {
				if err := tx.Commit(); err != nil {
					c.HTML(http.StatusInternalServerError, "fake_patients.html", gin.H{
						"title":    "產生假病患資料",
						"error":    "無法提交資料庫交易: " + err.Error(),
						"patients": patients,
					})
					return
				}
			} else {
				if err := tx.Commit(); err != nil {
					errorMessages = append(errorMessages, "無法提交資料庫交易: "+err.Error())
				}
			}
		} // 返回結果
		c.HTML(http.StatusOK, "fake_patients.html", gin.H{
			"title":        "產生假病患資料",
			"patients":     patients,
			"count":        count,
			"successCount": successCount,
			"insertToDB":   insertToDB,
			"errors":       errorMessages,
			"generated":    true,
			"timestamp":    time.Now().Format("2006-01-02 15:04:05"),
		})
	}
}
