// 日期格式化函數
function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.getFullYear() + '-' + 
           ('0' + (date.getMonth() + 1)).slice(-2) + '-' + 
           ('0' + date.getDate()).slice(-2);
}

// 顯示 SQL 語句
function showInsertSQL() {
    let sql = "";
    
    if (patientsData && patientsData.length > 0) {
        patientsData.forEach(function(patient) {
            // 主要病患資料 SQL
            sql += "-- 插入病患 " + patient.Name + "\n";
            sql += "INSERT INTO patient (name, gender, idno, age, birth, address, city, district, phone, mail, disease_id, emergency_contact, emergency_phone, emergency_relation, OTHERHISTORYDISEASE, OTHERMEDICALHISTORY)\n";
            sql += "VALUES ('" + patient.Name + "', '" + patient.Gender + "', '" + patient.IDNo + "', " + patient.Age + ", '" + formatDate(patient.Birth) + "',\n";
            sql += "'" + patient.Address + "', '" + patient.City + "', '" + patient.District + "', '" + patient.Phone + "', '" + patient.Mail + "',\n";
            sql += patient.DiseaseID + ", '" + patient.EmergencyContact + "', '" + patient.EmergencyPhone + "', '" + patient.EmergencyRelation + "',\n";
            sql += "'" + (patient.OtherHistoryDisease || "") + "', '" + (patient.OtherMedicalHistory || "") + "');\n\n";
            
            // 使用 LAST_INSERT_ID() 獲取最後插入的 ID
            sql += "-- 使用以下變數來存儲最後插入的ID\n";
            sql += "SET @last_patient_id = LAST_INSERT_ID();\n\n";
            
            // 插入病史資料
            if (patient.HistoryDiseases && patient.HistoryDiseases.length > 0) {
                sql += "-- 插入病史資料\n";
                patient.HistoryDiseases.forEach(function(disease) {
                    sql += "INSERT INTO patient_history_disease (patient_id, history_disease) VALUES (@last_patient_id, '" + disease + "');\n";
                });
                sql += '\n';
            }
            
            // 插入醫療史資料
            if (patient.MedicalHistories && patient.MedicalHistories.length > 0) {
                sql += "-- 插入醫療史資料\n";
                patient.MedicalHistories.forEach(function(history) {
                    sql += "INSERT INTO patient_medical_history (patient_id, medical_history) VALUES (@last_patient_id, '" + history + "');\n";
                });
                sql += '\n';
            }
            
            sql += "-- -------------------------------\n\n";
        });
    } else {
        sql = "// 沒有生成任何病患數據";
    }
    
    document.getElementById('sql-code').textContent = sql;
    document.getElementById('sql-modal').style.display = 'block';
}

// 複製到剪貼簿
function copyToClipboard() {
    const text = document.getElementById('sql-code').textContent;
    navigator.clipboard.writeText(text).then(function() {
        alert('已複製到剪貼簿');
    }, function() {
        alert('複製失敗');
    });
}

// 匯出 CSV
function exportToCSV() {
    let csvContent = "ID,姓名,性別,身分證字號,年齡,生日,城市,區域,地址,電話,信箱,疾病ID,緊急聯絡人,緊急聯絡人電話,緊急聯絡人關係,其他疾病史,其他醫療史,疾病史,醫療史\n";
    
    if (patientsData && patientsData.length > 0) {
        patientsData.forEach(function(patient) {
            const gender = patient.Gender === 'M' ? '男' : '女';
            const birthDate = formatDate(patient.Birth);
            const historyDiseases = patient.HistoryDiseases ? patient.HistoryDiseases.join('|') : '';
            const medicalHistories = patient.MedicalHistories ? patient.MedicalHistories.join('|') : '';
            
            csvContent += patient.ID + "," + patient.Name + "," + gender + "," + patient.IDNo + "," + patient.Age + "," + birthDate + "," + 
                        patient.City + "," + patient.District + "," + patient.Address + "," + patient.Phone + "," + patient.Mail + "," + 
                        patient.DiseaseID + "," + patient.EmergencyContact + "," + patient.EmergencyPhone + "," + patient.EmergencyRelation + "," + 
                        (patient.OtherHistoryDisease || '') + "," + (patient.OtherMedicalHistory || '') + "," + historyDiseases + "," + medicalHistories + "\n";
        });
        
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement("a");
        const url = URL.createObjectURL(blob);
        link.setAttribute("href", url);
        link.setAttribute("download", "fake_patients_" + new Date().toISOString().replace(/[:.]/g, '_') + ".csv");
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }
}

// DOM 載入完成後加入事件監聽
window.addEventListener('DOMContentLoaded', function() {
    // 如果有按鈕則添加事件監聽
    const btnExportCSV = document.querySelector('.btn-export');
    if (btnExportCSV) {
        btnExportCSV.addEventListener('click', exportToCSV);
    }
    
    const btnShowSQL = document.querySelector('.btn-info');
    if (btnShowSQL) {
        btnShowSQL.addEventListener('click', showInsertSQL);
    }
    
    const btnCopy = document.querySelector('.copy-btn');
    if (btnCopy) {
        btnCopy.addEventListener('click', copyToClipboard);
    }
    
    const closeModal = document.querySelector('#sql-modal span');
    if (closeModal) {
        closeModal.addEventListener('click', function() {
            document.getElementById('sql-modal').style.display = 'none';
        });
    }
});