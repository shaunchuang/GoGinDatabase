package handlers

import (
	"net/http"
	"strconv"
	"time"

	"golang-gin-app/internal/models"
	"golang-gin-app/internal/service"

	"github.com/gin-gonic/gin"
)

// AvailableSlotsFormHandler 處理顯示可預約時段表單的請求
func AvailableSlotsFormHandler(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 獲取所有醫師
		doctors, err := svc.GetDoctorUsers(c.Request.Context())
		if err != nil {
			c.HTML(http.StatusInternalServerError, "available_slots.html", gin.H{
				"title": "可預約時段管理",
				"error": "獲取醫師列表失敗: " + err.Error(),
			})
			return
		}

		// 獲取所有治療師
		therapists, err := svc.GetTherapistUsers(c.Request.Context())
		if err != nil {
			c.HTML(http.StatusInternalServerError, "available_slots.html", gin.H{
				"title": "可預約時段管理",
				"error": "獲取治療師列表失敗: " + err.Error(),
			})
			return
		}

		c.HTML(http.StatusOK, "available_slots.html", gin.H{
			"title":      "可預約時段管理",
			"doctors":    doctors,
			"therapists": therapists,
		})
	}
}

// GenerateAvailableSlotsHandler 處理生成可預約時段的請求
func GenerateAvailableSlotsHandler(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		doctorIDStr := c.PostForm("doctorID")
		daysStr := c.PostForm("days")
		slotsPerDayStr := c.PostForm("slotsPerDay")
		startHourStr := c.PostForm("startHour")
		slotDurationStr := c.PostForm("slotDuration")

		// 轉換參數
		doctorID, err := strconv.ParseInt(doctorIDStr, 10, 64)
		if err != nil || doctorID <= 0 {
			c.HTML(http.StatusBadRequest, "available_slots.html", gin.H{
				"title": "可預約時段管理",
				"error": "請提供有效的醫師/治療師ID",
			})
			return
		}

		days, err := strconv.Atoi(daysStr)
		if err != nil || days <= 0 || days > 365 {
			c.HTML(http.StatusBadRequest, "available_slots.html", gin.H{
				"title": "可預約時段管理",
				"error": "天數必須在1到365之間",
			})
			return
		}

		slotsPerDay, err := strconv.Atoi(slotsPerDayStr)
		if err != nil || slotsPerDay <= 0 || slotsPerDay > 24 {
			c.HTML(http.StatusBadRequest, "available_slots.html", gin.H{
				"title": "可預約時段管理",
				"error": "每天的時段數量必須在1到24之間",
			})
			return
		}

		startHour, err := strconv.Atoi(startHourStr)
		if err != nil || startHour < 0 || startHour > 23 {
			c.HTML(http.StatusBadRequest, "available_slots.html", gin.H{
				"title": "可預約時段管理",
				"error": "開始時間必須在0到23小時之間",
			})
			return
		}

		slotDuration, err := strconv.Atoi(slotDurationStr)
		if err != nil || slotDuration <= 0 || slotDuration > 240 {
			c.HTML(http.StatusBadRequest, "available_slots.html", gin.H{
				"title": "可預約時段管理",
				"error": "時段持續時間必須在1到240分鐘之間",
			})
			return
		}

		// 生成時段
		slots, err := svc.GenerateAvailableSlots(c.Request.Context(), doctorID, days, slotsPerDay, startHour, slotDuration)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "available_slots.html", gin.H{
				"title": "可預約時段管理",
				"error": "生成時段失敗: " + err.Error(),
			})
			return
		}

		// 獲取選擇的醫師/治療師用戶資訊
		var providerName string
		var allDoctors []*models.User
		var allTherapists []*models.User

		// 重新取得醫生與治療師列表
		allDoctors, _ = svc.GetDoctorUsers(c.Request.Context())
		allTherapists, _ = svc.GetTherapistUsers(c.Request.Context())
		// 尋找醫生/治療師姓名
		for _, doctor := range allDoctors {
			if doctor.ID == doctorID {
				if doctor.Username != nil {
					providerName = *doctor.Username
				} else {
					providerName = "未設定姓名"
				}
				break
			}
		}
		if providerName == "" {
			for _, therapist := range allTherapists {
				if therapist.ID == doctorID {
					if therapist.Username != nil {
						providerName = *therapist.Username
					} else {
						providerName = "未設定姓名"
					}
					break
				}
			}
		}

		// 返回結果
		message := strconv.Itoa(len(slots)) + " 個時段已經成功為 " + providerName + " 生成"

		c.HTML(http.StatusOK, "available_slots.html", gin.H{
			"title":      "可預約時段管理",
			"message":    message,
			"doctors":    allDoctors,
			"therapists": allTherapists,
			"slots":      slots,    // 返回生成的時段供顯示
			"selectedID": doctorID, // 保存選擇的醫師/治療師ID
		})
	}
}

// ViewAvailableSlotsHandler 顯示特定醫師/治療師的可預約時段
func ViewAvailableSlotsHandler(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		doctorIDStr := c.Query("doctorID")

		if doctorIDStr == "" {
			c.HTML(http.StatusBadRequest, "available_slots.html", gin.H{
				"title": "可預約時段管理",
				"error": "請提供醫師/治療師ID",
			})
			return
		}

		doctorID, err := strconv.ParseInt(doctorIDStr, 10, 64)
		if err != nil {
			c.HTML(http.StatusBadRequest, "available_slots.html", gin.H{
				"title": "可預約時段管理",
				"error": "無效的醫師/治療師ID",
			})
			return
		}

		// 獲取該醫師/治療師的所有時段
		slots, err := svc.GetAvailableSlotsByDoctor(c.Request.Context(), doctorID)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "available_slots.html", gin.H{
				"title": "可預約時段管理",
				"error": "獲取時段失敗: " + err.Error(),
			})
			return
		}

		// 獲取所有醫師和治療師
		doctors, _ := svc.GetDoctorUsers(c.Request.Context())
		therapists, _ := svc.GetTherapistUsers(c.Request.Context())
		// 找到選定的醫師/治療師名稱
		var providerName string
		for _, doctor := range doctors {
			if doctor.ID == doctorID {
				if doctor.Username != nil {
					providerName = *doctor.Username
				} else {
					providerName = "未設定姓名"
				}
				break
			}
		}
		if providerName == "" {
			for _, therapist := range therapists {
				if therapist.ID == doctorID {
					if therapist.Username != nil {
						providerName = *therapist.Username
					} else {
						providerName = "未設定姓名"
					}
					break
				}
			}
		}

		c.HTML(http.StatusOK, "available_slots.html", gin.H{
			"title":        "可預約時段管理",
			"doctors":      doctors,
			"therapists":   therapists,
			"slots":        slots,
			"selectedID":   doctorID,
			"providerName": providerName,
		})
	}
}

// EditAvailableSlotFormHandler 顯示編輯可預約時段的表單
func EditAvailableSlotFormHandler(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		slotIDStr := c.Param("id")
		slotID, err := strconv.ParseInt(slotIDStr, 10, 64)
		if err != nil {
			c.HTML(http.StatusBadRequest, "available_slots.html", gin.H{
				"title": "編輯可預約時段",
				"error": "無效的時段ID",
			})
			return
		}

		// 獲取時段信息
		slot, err := svc.GetAvailableSlotByID(c.Request.Context(), slotID)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "available_slots.html", gin.H{
				"title": "編輯可預約時段",
				"error": "獲取時段信息失敗: " + err.Error(),
			})
			return
		}

		// 獲取所有醫師和治療師
		doctors, _ := svc.GetDoctorUsers(c.Request.Context())
		therapists, _ := svc.GetTherapistUsers(c.Request.Context())

		c.HTML(http.StatusOK, "edit_slot.html", gin.H{
			"title":      "編輯可預約時段",
			"slot":       slot,
			"doctors":    doctors,
			"therapists": therapists,
		})
	}
}

// UpdateAvailableSlotHandler 處理更新可預約時段的請求
func UpdateAvailableSlotHandler(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		slotIDStr := c.Param("id")
		slotID, err := strconv.ParseInt(slotIDStr, 10, 64)
		if err != nil {
			c.HTML(http.StatusBadRequest, "available_slots.html", gin.H{
				"title": "更新可預約時段",
				"error": "無效的時段ID",
			})
			return
		}

		// 獲取現有時段
		slot, err := svc.GetAvailableSlotByID(c.Request.Context(), slotID)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "available_slots.html", gin.H{
				"title": "更新可預約時段",
				"error": "獲取時段信息失敗: " + err.Error(),
			})
			return
		}

		// 解析表單數據
		doctorIDStr := c.PostForm("doctorID")
		dateStr := c.PostForm("date")
		beginTimeStr := c.PostForm("beginTime")
		endTimeStr := c.PostForm("endTime")
		isBookedStr := c.PostForm("isBooked")

		doctorID, err := strconv.ParseInt(doctorIDStr, 10, 64)
		if err != nil || doctorID <= 0 {
			c.HTML(http.StatusBadRequest, "edit_slot.html", gin.H{
				"title": "更新可預約時段",
				"error": "無效的醫師/治療師ID",
				"slot":  slot,
			})
			return
		}

		// 解析日期和時間
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.HTML(http.StatusBadRequest, "edit_slot.html", gin.H{
				"title": "更新可預約時段",
				"error": "無效的日期格式",
				"slot":  slot,
			})
			return
		}

		// 解析開始時間和結束時間
		beginTime, err := time.Parse("15:04", beginTimeStr)
		if err != nil {
			c.HTML(http.StatusBadRequest, "edit_slot.html", gin.H{
				"title": "更新可預約時段",
				"error": "無效的開始時間格式",
				"slot":  slot,
			})
			return
		}

		endTime, err := time.Parse("15:04", endTimeStr)
		if err != nil {
			c.HTML(http.StatusBadRequest, "edit_slot.html", gin.H{
				"title": "更新可預約時段",
				"error": "無效的結束時間格式",
				"slot":  slot,
			})
			return
		}

		// 組合完整的日期時間
		slotBeginTime := time.Date(
			date.Year(), date.Month(), date.Day(),
			beginTime.Hour(), beginTime.Minute(), 0, 0, time.Local)

		slotEndTime := time.Date(
			date.Year(), date.Month(), date.Day(),
			endTime.Hour(), endTime.Minute(), 0, 0, time.Local)

		// 檢查時間順序
		if slotBeginTime.After(slotEndTime) {
			c.HTML(http.StatusBadRequest, "edit_slot.html", gin.H{
				"title": "更新可預約時段",
				"error": "開始時間不能晚於結束時間",
				"slot":  slot,
			})
			return
		}

		// 更新時段
		slot.Doctor = doctorID
		slot.SlotDate = date
		slot.SlotBeginTime = slotBeginTime
		slot.SlotEndTime = slotEndTime
		slot.IsBooked = isBookedStr == "true"

		if err := svc.UpdateAvailableSlot(c.Request.Context(), slot); err != nil {
			c.HTML(http.StatusInternalServerError, "edit_slot.html", gin.H{
				"title": "更新可預約時段",
				"error": "更新時段失敗: " + err.Error(),
				"slot":  slot,
			})
			return
		}

		// 重定向回查看頁面
		c.Redirect(http.StatusFound, "/available-slots/view?doctorID="+doctorIDStr)
	}
}

// DeleteAvailableSlotHandler 處理刪除可預約時段的請求
func DeleteAvailableSlotHandler(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		slotIDStr := c.Param("id")
		doctorIDStr := c.Query("doctorID") // 保存醫師/治療師ID用於重定向

		slotID, err := strconv.ParseInt(slotIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的時段ID"})
			return
		}

		// 刪除時段
		if err := svc.DeleteAvailableSlot(c.Request.Context(), slotID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除時段失敗: " + err.Error()})
			return
		}

		// 如果是AJAX請求，返回JSON響應
		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusOK, gin.H{"success": true})
			return
		}

		// 否則重定向回查看頁面
		c.Redirect(http.StatusFound, "/available-slots/view?doctorID="+doctorIDStr)
	}
}
