package models

import (
	"time"
)

// AvailableSlot 表示醫師/治療師的可預約時段
type AvailableSlot struct {
	ID            int64     `json:"id"`
	Doctor        int64     `json:"doctor"` // 醫師/治療師的用戶ID
	IsBooked      bool      `json:"is_booked"`
	SlotBeginTime time.Time `json:"slot_begin_time"`
	SlotDate      time.Time `json:"slot_date"`
	SlotEndTime   time.Time `json:"slot_end_time"`
}

// SlotGenerationRequest 表示生成時段的請求參數
type SlotGenerationRequest struct {
	DoctorID     int64 `json:"doctor_id"`     // 醫師/治療師ID
	Days         int   `json:"days"`          // 要生成的天數
	SlotsPerDay  int   `json:"slots_per_day"` // 每天要生成的時段數量
	StartHour    int   `json:"start_hour"`    // 開始時間（小時）
	EndHour      int   `json:"end_hour"`      // 結束時間（小時）
	SlotDuration int   `json:"slot_duration"` // 每個時段的持續時間（分鐘）
}
