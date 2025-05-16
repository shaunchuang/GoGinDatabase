package service

import (
	"context"
	"fmt"
	"golang-gin-app/internal/models"
	"time"
)

// GenerateAvailableSlots 根據指定條件生成可預約時段
func (s *Service) GenerateAvailableSlots(ctx context.Context, doctorID int64, days int, slotsPerDay int, startHour int, slotDuration int) ([]*models.AvailableSlot, error) {
	// 基本參數驗證
	if doctorID <= 0 {
		return nil, fmt.Errorf("醫師/治療師ID必須大於0")
	}
	if days <= 0 || days > 365 {
		return nil, fmt.Errorf("天數必須在1到365之間")
	}
	if slotsPerDay <= 0 || slotsPerDay > 24 {
		return nil, fmt.Errorf("每天時段數必須在1到24之間")
	}
	if startHour < 0 || startHour > 23 {
		return nil, fmt.Errorf("開始時間必須在0到23之間")
	}
	if slotDuration <= 0 || slotDuration > 240 {
		return nil, fmt.Errorf("每個時段的持續時間必須在1到240分鐘之間")
	}

	// 生成預約時段
	slots := make([]*models.AvailableSlot, 0, days*slotsPerDay)
	today := time.Now().Truncate(24 * time.Hour) // 今天日期，去掉時間部分

	// 逐天生成
	for day := 0; day < days; day++ {
		slotDate := today.AddDate(0, 0, day)

		// 計算該天的時段
		for slotIndex := 0; slotIndex < slotsPerDay; slotIndex++ {
			// 計算開始時間
			beginHour := startHour + (slotIndex*slotDuration)/60
			beginMinute := (slotIndex * slotDuration) % 60

			if beginHour >= 24 {
				// 如果時段超過了當天，就停止為這一天生成時段
				break
			}

			// 計算結束時間
			endHour := beginHour
			endMinute := beginMinute + slotDuration

			// 處理分鐘進位
			if endMinute >= 60 {
				endHour += endMinute / 60
				endMinute = endMinute % 60
			}

			// 如果結束時間超過24小時，則調整為23:59
			if endHour >= 24 {
				endHour = 23
				endMinute = 59
			}

			// 創建時間對象
			beginTime := time.Date(
				slotDate.Year(), slotDate.Month(), slotDate.Day(),
				beginHour, beginMinute, 0, 0, time.Local)
			endTime := time.Date(
				slotDate.Year(), slotDate.Month(), slotDate.Day(),
				endHour, endMinute, 0, 0, time.Local)

			// 創建時段對象
			slot := &models.AvailableSlot{
				Doctor:        doctorID,
				IsBooked:      false,
				SlotBeginTime: beginTime,
				SlotDate:      slotDate,
				SlotEndTime:   endTime,
			}
			slots = append(slots, slot)
		}
	}

	// 批量保存到數據庫
	if err := s.repo.BatchCreateAvailableSlots(ctx, slots); err != nil {
		return nil, fmt.Errorf("保存預約時段失敗: %v", err)
	}

	return slots, nil
}

// GetAvailableSlotsByDoctor 獲取指定醫師的可預約時段
func (s *Service) GetAvailableSlotsByDoctor(ctx context.Context, doctorID int64) ([]*models.AvailableSlot, error) {
	return s.repo.GetAvailableSlotsByDoctor(ctx, doctorID)
}

// GetDoctorUsers 獲取具有醫師角色的用戶
func (s *Service) GetDoctorUsers(ctx context.Context) ([]*models.User, error) {
	// 醫師角色ID為3，根據utils.go中的常數
	return s.repo.GetUserByRoleID(ctx, 3) // DOCTOR_ROLE_ID = 3
}

// GetTherapistUsers 獲取具有治療師角色的用戶
func (s *Service) GetTherapistUsers(ctx context.Context) ([]*models.User, error) {
	// 合併所有治療師角色的用戶
	// 治療師角色ID為4,5,6,7，根據utils.go中的常數
	therapistRoles := []int64{4, 5, 6, 7} // DTX_PSY_ROLE_ID = 4, DTX_ST_ROLE_ID = 5, DTX_OT_ROLE_ID = 6, DTX_PI_ROLE_ID = 7

	allTherapists := make([]*models.User, 0)
	uniqueIDs := make(map[int64]bool)

	for _, roleID := range therapistRoles {
		users, err := s.repo.GetUserByRoleID(ctx, roleID)
		if err != nil {
			return nil, fmt.Errorf("獲取角色ID %d 的用戶失敗: %v", roleID, err)
		}

		// 避免重複用戶
		for _, user := range users {
			if _, exists := uniqueIDs[user.ID]; !exists {
				uniqueIDs[user.ID] = true
				allTherapists = append(allTherapists, user)
			}
		}
	}

	return allTherapists, nil
}

// UpdateAvailableSlot 更新可預約時段
func (s *Service) UpdateAvailableSlot(ctx context.Context, slot *models.AvailableSlot) error {
	// 驗證必要字段
	if slot.ID <= 0 {
		return fmt.Errorf("無效的時段ID")
	}
	if slot.Doctor <= 0 {
		return fmt.Errorf("無效的醫師/治療師ID")
	}

	// 檢查時間格式
	if slot.SlotBeginTime.After(slot.SlotEndTime) {
		return fmt.Errorf("開始時間不能晚於結束時間")
	}

	return s.repo.UpdateAvailableSlot(ctx, slot)
}

// DeleteAvailableSlot 刪除可預約時段
func (s *Service) DeleteAvailableSlot(ctx context.Context, slotID int64) error {
	if slotID <= 0 {
		return fmt.Errorf("無效的時段ID")
	}

	// 先檢查時段是否存在
	slot, err := s.repo.GetAvailableSlotByID(ctx, slotID)
	if err != nil {
		return err
	}

	// 如果時段已被預約，不允許刪除
	if slot.IsBooked {
		return fmt.Errorf("該時段已被預約，無法刪除")
	}

	return s.repo.DeleteAvailableSlot(ctx, slotID)
}

// GetAvailableSlotByID 通過ID獲取時段
func (s *Service) GetAvailableSlotByID(ctx context.Context, slotID int64) (*models.AvailableSlot, error) {
	return s.repo.GetAvailableSlotByID(ctx, slotID)
}
