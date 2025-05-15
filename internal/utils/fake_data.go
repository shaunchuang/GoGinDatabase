package utils

import (
	"database/sql"
	"fmt"
	"golang-gin-app/internal/models"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// 常見中文姓氏（已擴充、去重）
var chineseSurnames = []string{
	"陳", "林", "黃", "張", "李", "王", "吳", "劉",
	"蔡", "楊", "許", "鄭", "謝", "洪", "郭", "曾",
	"邱", "廖", "江", "賴", "何", "呂", "余", "葉",
	"趙", "鍾", "潘", "田", "游", "莊", "方", "石",
	"章", "蔣", "唐", "韓", "蕭", "鄧", "梁", "宋",
	"范", "彭", "曹", "魏", "夏", "賀", "姚", "蘇",
	"杜", "龔",
}

// 常見中文名字用字
var chineseNameRunes = []rune(
	"偉芳娜敏俊艷悅蓉婷峰翔志誠豪欣怡淑芬秀英君建國宏文雄強美蘭珍梅慧琳玲玉環瑩" +
		"雯琪凱安宸瑋語嫣詩涵雅庭睿哲梓子宜萱彥廷啟航詠晴知淇奕辰晉銘遠瑞昕曉彤弘嘉祺瑤軒靜凡筱宇霖念慈萍思源雨薇芷若依蔓惜霏煌洛旭筠羿恆孟心昌逸飛毅",
)

// 用於生成醫師和治療師的唯一編號
var (
	doctorNumberMutex  sync.Mutex
	therapyNumberMutex sync.Mutex
	lastDoctorNumber   = 0
	lastTherapyNumber  = 0
	isInitialized      = false
)

// 角色ID常數，對應資料庫中的角色ID
const (
	USER_ROLE_ID    = 1
	ADMIN_ROLE_ID   = 2
	DOCTOR_ROLE_ID  = 3
	DTX_PSY_ROLE_ID = 4
	DTX_ST_ROLE_ID  = 5
	DTX_OT_ROLE_ID  = 6
	DTX_PI_ROLE_ID  = 7
)

func init() {
	// 先用時間當種子
	rand.Seed(time.Now().UnixNano())
	// 註冊一個自訂的 gofakeit 欄位生成功能 "{chinese_name}"
	gofakeit.AddFuncLookup("chinese_name", gofakeit.Info{
		Display:     "Chinese Name",
		Category:    "person",
		Description: "隨機產生中文姓名（姓＋1~2字名）",
		Example:     "陳大文",
		Output:      "string",
		Params:      []gofakeit.Param{},
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {
			// 隨機挑一個姓氏
			surname := chineseSurnames[rand.Intn(len(chineseSurnames))]
			// 決定名字長度 1 或 2 個字
			nameLen := rand.Intn(2) + 1
			name := ""
			for i := 0; i < nameLen; i++ {
				name += string(chineseNameRunes[rand.Intn(len(chineseNameRunes))])
			}
			return surname + name, nil
		},
	})

	// 註冊一個自訂的 gofakeit 欄位生成功能 "{taiwan_phone}"
	gofakeit.AddFuncLookup("taiwan_phone", gofakeit.Info{
		Display:     "Taiwan Phone",
		Category:    "person",
		Description: "隨機產生台灣手機號碼格式",
		Example:     "0912345678",
		Output:      "string",
		Params:      []gofakeit.Param{},
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {
			// 台灣手機前兩碼通常是 09
			prefix := "09"
			// 隨機生成 8 位數字
			digits := ""
			for i := 0; i < 8; i++ {
				digits += fmt.Sprintf("%d", rand.Intn(10))
			}
			return prefix + digits, nil
		},
	})
}

// InitializeCounters 從數據庫初始化計數器，防止重複帳號
func InitializeCounters(db *sql.DB) error {
	if isInitialized {
		return nil
	}

	// 獲取當前醫師編號的最大值
	doctorQuery := `SELECT MAX(CAST(SUBSTRING(account, 7) AS UNSIGNED)) FROM user WHERE account LIKE 'doctor%'`
	var maxDoctorNum sql.NullInt64
	err := db.QueryRow(doctorQuery).Scan(&maxDoctorNum)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get max doctor number: %v", err)
	}

	// 獲取當前治療師編號的最大值
	therapyQuery := `SELECT MAX(CAST(SUBSTRING(account, 8) AS UNSIGNED)) FROM user WHERE account LIKE 'therapy%'`
	var maxTherapyNum sql.NullInt64
	err = db.QueryRow(therapyQuery).Scan(&maxTherapyNum)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get max therapy number: %v", err)
	}

	// 設置計數器初始值
	doctorNumberMutex.Lock()
	if maxDoctorNum.Valid {
		lastDoctorNumber = int(maxDoctorNum.Int64)
	}
	doctorNumberMutex.Unlock()

	therapyNumberMutex.Lock()
	if maxTherapyNum.Valid {
		lastTherapyNumber = int(maxTherapyNum.Int64)
	}
	therapyNumberMutex.Unlock()

	isInitialized = true
	return nil
}

// getNextDoctorNumber 獲取下一個醫師編號
func getNextDoctorNumber() int {
	doctorNumberMutex.Lock()
	defer doctorNumberMutex.Unlock()
	lastDoctorNumber++
	return lastDoctorNumber
}

// getNextTherapyNumber 獲取下一個治療師編號
func getNextTherapyNumber() int {
	therapyNumberMutex.Lock()
	defer therapyNumberMutex.Unlock()
	lastTherapyNumber++
	return lastTherapyNumber
}

// GenerateFakeUser creates a fake user with realistic data
func GenerateFakeUser(userType string) *models.User {
	now := time.Now()
	pastDate := gofakeit.DateRange(now.AddDate(-5, 0, 0), now)
	steamID := gofakeit.UUID()
	username, _ := gofakeit.Generate("{chinese_name}")

	var account, email string

	// 根據使用者類型產生不同的 account 和 email
	if strings.ToLower(userType) == "doctor" {
		number := getNextDoctorNumber()
		account = fmt.Sprintf("doctor%d", number)
		email = fmt.Sprintf("%s@example.com", account)
	} else { // therapy
		number := getNextTherapyNumber()
		account = fmt.Sprintf("therapy%d", number)
		email = fmt.Sprintf("%s@example.com", account)
	}

	// 使用我們自定義的 taiwan_phone 生成器
	phoneNumber, _ := gofakeit.Generate("{taiwan_phone}")

	return &models.User{
		Account:       account,
		CreateTime:    now,
		Email:         email,
		LastLoginDate: &pastDate,
		Password:      "$2a$12$qDECDR.WBiP2Xueb5ftW3.LKslm6.Gs7oeTH1T3SmnUxpucvUm8sW",
		Status:        "APPROVED",
		SteamID:       &steamID,
		TelCell:       phoneNumber,
		// 產生像「張偉」「李欣怡」這樣的中文姓名
		Username: username,
		// 初始化空角色陣列
		Roles: make([]*models.Role, 0),
	}
}

// GenerateFakeUsers creates a slice of fake users
func GenerateFakeUsers(count int, userType string) []*models.User {
	users := make([]*models.User, count)
	for i := 0; i < count; i++ {
		users[i] = GenerateFakeUser(userType)
	}
	return users
}

// GetDefaultRoleIDsForUserType 根據使用者類型獲取預設角色 ID 列表
// 此函數可在沒有前端角色選擇時提供預設值
func GetDefaultRoleIDsForUserType(userType string) []int64 {
	switch strings.ToLower(userType) {
	case "doctor":
		return []int64{DOCTOR_ROLE_ID} // 醫師只需要醫師角色
	case "therapy":
		// 對於治療師，隨機選擇一種治療師角色，不含 USER 角色
		therapistRoles := []int64{DTX_PSY_ROLE_ID, DTX_ST_ROLE_ID, DTX_OT_ROLE_ID, DTX_PI_ROLE_ID}
		randomIndex := rand.Intn(len(therapistRoles))
		return []int64{therapistRoles[randomIndex]}
	default:
		return []int64{} // 不自動分配任何角色
	}
}
