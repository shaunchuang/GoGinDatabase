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

// 台灣城市和區域對應表
var taiwanCityDistricts = map[string][]string{
	"台北市": {"中正區", "大同區", "中山區", "松山區", "大安區", "萬華區", "信義區", "士林區", "北投區", "內湖區", "南港區", "文山區"},
	"新北市": {"板橋區", "三重區", "中和區", "永和區", "新莊區", "新店區", "樹林區", "鶯歌區", "三峽區", "淡水區", "汐止區", "瑞芳區"},
	"桃園市": {"桃園區", "中壢區", "平鎮區", "八德區", "楊梅區", "蘆竹區", "龜山區", "龍潭區", "大溪區", "大園區"},
	"台中市": {"中區", "東區", "南區", "西區", "北區", "北屯區", "西屯區", "南屯區", "太平區", "大里區", "霧峰區", "烏日區"},
	"台南市": {"中西區", "東區", "南區", "北區", "安平區", "安南區", "永康區", "仁德區", "歸仁區", "新化區"},
	"高雄市": {"鹽埕區", "鼓山區", "左營區", "楠梓區", "三民區", "新興區", "前金區", "苓雅區", "前鎮區", "旗津區", "小港區"},
	"基隆市": {"仁愛區", "信義區", "中正區", "中山區", "安樂區", "暖暖區", "七堵區"},
	"新竹市": {"東區", "北區", "香山區"},
	"嘉義市": {"東區", "西區"},
	"新竹縣": {"竹北市", "竹東鎮", "新埔鎮", "關西鎮", "湖口鄉", "新豐鄉"},
	"苗栗縣": {"苗栗市", "頭份市", "竹南鎮", "後龍鎮", "通霄鎮", "苑裡鎮"},
	"彰化縣": {"彰化市", "員林市", "和美鎮", "鹿港鎮", "溪湖鎮", "田中鎮"},
	"南投縣": {"南投市", "草屯鎮", "埔里鎮", "竹山鎮", "集集鎮", "名間鄉"},
	"雲林縣": {"斗六市", "斗南鎮", "虎尾鎮", "西螺鎮", "土庫鎮", "北港鎮"},
	"嘉義縣": {"太保市", "朴子市", "布袋鎮", "大林鎮", "民雄鄉", "溪口鄉"},
	"屏東縣": {"屏東市", "潮州鎮", "東港鎮", "恆春鎮", "萬丹鄉", "長治鄉"},
	"宜蘭縣": {"宜蘭市", "羅東鎮", "蘇澳鎮", "頭城鎮", "礁溪鄉", "壯圍鄉"},
	"花蓮縣": {"花蓮市", "鳳林鎮", "玉里鎮", "新城鄉", "吉安鄉", "壽豐鄉"},
	"台東縣": {"台東市", "成功鎮", "關山鎮", "卑南鄉", "鹿野鄉", "池上鄉"},
	"澎湖縣": {"馬公市", "湖西鄉", "白沙鄉", "西嶼鄉", "望安鄉", "七美鄉"},
	"金門縣": {"金城鎮", "金湖鎮", "金沙鎮", "金寧鄉", "烈嶼鄉", "烏坵鄉"},
	"連江縣": {"南竿鄉", "北竿鄉", "莒光鄉", "東引鄉"},
}

// 疾病史選項（從資料庫 history_disease 表獲取）
var historyDiseases = []string{
	"中樞神經損傷", "心血管疾病", "呼吸方面疾病", "肝臟疾病",
	"糖尿病", "腎臟病", "癌症", "免疫相關疾病",
}

// 醫療史選項
var medicalHistories = []string{
	"骨折", "開刀", "住院治療", "重大傷病", "藥物過敏",
	"食物過敏", "輸血", "化療", "放射治療", "洗腎治療",
	"器官移植", "心導管", "心律不整裝置", "人工關節", "牙科手術",
}

// 緊急聯絡人關係選項
var emergencyRelations = []string{
	"父親", "母親", "配偶", "兒子", "女兒",
	"兄弟", "姊妹", "祖父", "祖母", "叔叔",
	"阿姨", "表親", "朋友", "同事", "監護人",
}

// GenerateFakePatients 生成指定數量的假病患資料
func GenerateFakePatients(count int) ([]*models.Patient, error) {
	patients := make([]*models.Patient, 0, count)

	for i := 0; i < count; i++ {
		// 隨機生成身分證字號
		idno := generateTaiwanID()

		// 基本個人信息
		gender := ""
		if rand.Intn(2) == 0 {
			gender = "M"
		} else {
			gender = "F"
		}

		// 隨機生成出生日期（18-90歲）
		minAge := 18
		maxAge := 90
		age := rand.Intn(maxAge-minAge) + minAge
		birthYear := time.Now().Year() - age
		birthMonth := rand.Intn(12) + 1
		birthDay := rand.Intn(28) + 1 // 簡化處理，避免月份天數問題
		birth := time.Date(birthYear, time.Month(birthMonth), birthDay, 0, 0, 0, 0, time.Local)

		// 隨機選擇城市和區域
		city := getRandomKey(taiwanCityDistricts)
		districts := taiwanCityDistricts[city]
		district := districts[rand.Intn(len(districts))]

		// 隨機生成地址
		address := fmt.Sprintf("%s%d號",
			[]string{"中山路", "中正路", "復興路", "和平路", "民生路", "建國路", "光明街", "仁愛路", "忠孝路", "信義路"}[rand.Intn(10)],
			rand.Intn(200)+1,
		)

		// 隨機選擇病史和醫療史
		otherHistoryDisease := ""
		otherMedicalHistory := ""

		// 少數情況下有自定義病史和醫療史
		if rand.Intn(10) < 3 {
			otherHistoryDisease = []string{"", "家族有糖尿病史", "曾有嚴重過敏", "青少年哮喘", "其他慢性疾病"}[rand.Intn(5)]
		}
		if rand.Intn(10) < 3 {
			otherMedicalHistory = []string{"", "曾動過小手術", "曾做過重大手術", "有長期用藥", "最近有服用特殊藥物"}[rand.Intn(5)]
		}

		// 隨機選擇緊急聯絡人關係
		relation := emergencyRelations[rand.Intn(len(emergencyRelations))]

		patient := &models.Patient{
			ID:                  int64(i + 1),
			UserID:              int64(rand.Intn(100) + 1), // 隨機分配一個用戶ID
			Name:                generateChineseName(),
			Gender:              gender,
			IDNo:                idno,
			Age:                 age,
			Birth:               birth,
			Address:             address,
			City:                city,
			District:            district,
			Phone:               generateTaiwanPhone(),
			Mail:                gofakeit.Email(),
			DiseaseID:           int64(rand.Intn(10) + 1),
			EmergencyContact:    generateChineseName(),
			EmergencyPhone:      generateTaiwanPhone(),
			EmergencyRelation:   relation,
			OtherHistoryDisease: otherHistoryDisease,
			OtherMedicalHistory: otherMedicalHistory,
		}

		// 隨機選擇0-3個病史
		historyDiseaseCount := rand.Intn(4)
		if historyDiseaseCount > 0 {
			patient.HistoryDiseases = make([]string, 0, historyDiseaseCount)
			// 創建一個副本以便進行隨機選擇而不重複
			availableHistoryDiseases := make([]string, len(historyDiseases))
			copy(availableHistoryDiseases, historyDiseases)

			for j := 0; j < historyDiseaseCount; j++ {
				if len(availableHistoryDiseases) == 0 {
					break
				}
				// 隨機選擇一個病史項目
				idx := rand.Intn(len(availableHistoryDiseases))
				patient.HistoryDiseases = append(patient.HistoryDiseases, availableHistoryDiseases[idx])
				// 從可用列表中移除已選項目以避免重複
				availableHistoryDiseases = append(availableHistoryDiseases[:idx], availableHistoryDiseases[idx+1:]...)
			}
		}

		// 隨機選擇0-3個醫療史
		medicalHistoryCount := rand.Intn(4)
		if medicalHistoryCount > 0 {
			patient.MedicalHistories = make([]string, 0, medicalHistoryCount)
			// 創建一個副本以便進行隨機選擇而不重複
			availableMedicalHistories := make([]string, len(medicalHistories))
			copy(availableMedicalHistories, medicalHistories)

			for j := 0; j < medicalHistoryCount; j++ {
				if len(availableMedicalHistories) == 0 {
					break
				}
				// 隨機選擇一個醫療史項目
				idx := rand.Intn(len(availableMedicalHistories))
				patient.MedicalHistories = append(patient.MedicalHistories, availableMedicalHistories[idx])
				// 從可用列表中移除已選項目以避免重複
				availableMedicalHistories = append(availableMedicalHistories[:idx], availableMedicalHistories[idx+1:]...)
			}
		}

		patients = append(patients, patient)
	}

	return patients, nil
}

// generateTaiwanID 生成台灣身分證字號
func generateTaiwanID() string {
	// 第一個字母代表地區
	letters := "ABCDEFGHJKLMNPQRSTUVXYWZIO"
	firstLetter := string(letters[rand.Intn(len(letters))])

	// 第二個數字代表性別（1男性，2女性）
	genderNum := rand.Intn(2) + 1

	// 隨機生成其余7個數字
	restNums := ""
	for i := 0; i < 7; i++ {
		restNums += fmt.Sprintf("%d", rand.Intn(10))
	}

	// 最後一個檢查碼先假設為0（實際上應該有正確的檢查碼計算）
	return fmt.Sprintf("%s%d%s%d", firstLetter, genderNum, restNums, rand.Intn(10))
}

// getRandomKey 從map中隨機選擇一個鍵
func getRandomKey(m map[string][]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys[rand.Intn(len(keys))]
}

// generateChineseName 產生隨機中文姓名
func generateChineseName() string {
	// 隨機挑一個姓氏
	surname := chineseSurnames[rand.Intn(len(chineseSurnames))]
	// 決定名字長度 1 或 2 個字
	nameLen := rand.Intn(2) + 1
	name := ""
	for i := 0; i < nameLen; i++ {
		name += string(chineseNameRunes[rand.Intn(len(chineseNameRunes))])
	}
	return surname + name
}

// generateTaiwanPhone 產生隨機台灣手機號碼
func generateTaiwanPhone() string {
	// 台灣手機前兩碼通常是 09
	prefix := "09"
	// 隨機生成 8 位數字
	digits := ""
	for i := 0; i < 8; i++ {
		digits += fmt.Sprintf("%d", rand.Intn(10))
	}
	return prefix + digits
}
