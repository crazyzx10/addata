package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//go:embed frontend/*
var frontendFS embed.FS

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbName"`
	} `json:"database"`
	MiniPrograms []struct {
		Name          string `json:"name"`
		AppID         string `json:"appId"`
		Secret        string `json:"secret"`
		FirstPullComplete bool `json:"firstPullComplete"`
	} `json:"miniPrograms"`
}

var (
	configPath string
	logChan    chan string
	mu         sync.Mutex
)

func init() {
	if runtime.GOOS == "windows" {
		configPath = filepath.Join(os.Getenv("APPDATA"), "WechatAdConfig", "config.json")
	} else {
		home, _ := os.UserHomeDir()
		configPath = filepath.Join(home, ".wechatadconfig", "config.json")
	}
	logChan = make(chan string, 1000)
}

func loadConfig() (Config, error) {
	var cfg Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func saveConfig(cfg Config) error {
	dir := filepath.Dir(configPath)
	os.MkdirAll(dir, 0755)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func getDBConnection(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)
	return sql.Open("mysql", dsn)
}

func testDatabaseConnection(w http.ResponseWriter, r *http.Request) {
	var cfg Config
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := getDBConnection(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "数据库连接成功"})
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	cfg, _ := loadConfig()
	json.NewEncoder(w).Encode(cfg)
}

func saveConfigHandler(w http.ResponseWriter, r *http.Request) {
	var newCfg Config
	if err := json.NewDecoder(r.Body).Decode(&newCfg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oldCfg, _ := loadConfig()
	if newCfg.Database.Host == "" {
		newCfg.Database = oldCfg.Database
	}
	if len(newCfg.MiniPrograms) == 0 {
		newCfg.MiniPrograms = oldCfg.MiniPrograms
	}

	if err := saveConfig(newCfg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func getToken(appID, secret string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appID, secret)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Errcode     int    `json:"errcode"`
		Errmsg      string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Errcode != 0 {
		return "", fmt.Errorf("%s", result.Errmsg)
	}
	return result.AccessToken, nil
}

func initDatabase(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS wechat_ad_summary (
			id INT AUTO_INCREMENT PRIMARY KEY,
			mini_program_name VARCHAR(255) NOT NULL,
			date DATE NOT NULL,
			ad_slot_id VARCHAR(100),
			ad_slot_name VARCHAR(255),
			ad_slot_type VARCHAR(100),
			total_request_count INT,
			total_served_count INT,
			total_click_count INT,
			total_revenue DECIMAL(10,4),
			total_ecpm DECIMAL(10,4),
			total_impression_rate DECIMAL(10,4),
			total_click_rate DECIMAL(10,4),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_date (date),
			INDEX idx_mini_program (mini_program_name),
			UNIQUE KEY unique_summary (mini_program_name, date, ad_slot_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS wechat_ad_unit_list (
			id INT AUTO_INCREMENT PRIMARY KEY,
			mini_program_name VARCHAR(255) NOT NULL,
			ad_slot_id VARCHAR(100) NOT NULL,
			ad_slot_name VARCHAR(255),
			ad_slot_type VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_mini_program (mini_program_name),
			INDEX idx_ad_slot_id (ad_slot_id),
			UNIQUE KEY unique_unit (mini_program_name, ad_slot_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func syncAdUnitList(db *sql.DB, miniProgramName, accessToken string) error {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/ad/get?access_token=%s", accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		List []struct {
			AdSlotId   string `json:"ad_slot_id"`
			AdSlotName string `json:"ad_slot_name"`
			AdSlotType string `json:"ad_slot_type"`
		} `json:"list"`
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Errcode != 0 {
		return fmt.Errorf("%s", result.Errmsg)
	}

	for _, unit := range result.List {
		_, err := db.Exec(`
			INSERT INTO wechat_ad_unit_list 
			(mini_program_name, ad_slot_id, ad_slot_name, ad_slot_type) 
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE 
			ad_slot_name = VALUES(ad_slot_name), 
			ad_slot_type = VALUES(ad_slot_type)
		`, miniProgramName, unit.AdSlotId, unit.AdSlotName, unit.AdSlotType)
		if err != nil {
			return err
		}
	}
	return nil
}

func syncSummaryData(db *sql.DB, miniProgramName, accessToken, startDate, endDate string, firstPull bool) error {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/ad/getdata?access_token=%s", accessToken)
	payload := map[string]string{
		"start_date": startDate,
		"end_date":   endDate,
	}
	jsonPayload, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		List []struct {
			Date           string  `json:"date"`
			AdSlotId       string  `json:"ad_slot_id"`
			AdSlotName     string  `json:"ad_slot_name"`
			AdSlotType     string  `json:"ad_slot_type"`
			ReqCount       int     `json:"req_count"`
			ServedCount    int     `json:"served_count"`
			ClickCount     int     `json:"click_count"`
			Revenue        float64 `json:"revenue"`
			Ecpm           float64 `json:"ecpm"`
			ImpressionRate float64 `json:"impression_rate"`
			ClickRate      float64 `json:"click_rate"`
		} `json:"list"`
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Errcode != 0 {
		return fmt.Errorf("%s", result.Errmsg)
	}

	for _, item := range result.List {
		_, err := db.Exec(`
			INSERT INTO wechat_ad_summary 
			(mini_program_name, date, ad_slot_id, ad_slot_name, ad_slot_type, 
			 total_request_count, total_served_count, total_click_count, 
			 total_revenue, total_ecpm, total_impression_rate, total_click_rate) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE 
			ad_slot_name = VALUES(ad_slot_name),
			ad_slot_type = VALUES(ad_slot_type),
			total_request_count = VALUES(total_request_count),
			total_served_count = VALUES(total_served_count),
			total_click_count = VALUES(total_click_count),
			total_revenue = VALUES(total_revenue),
			total_ecpm = VALUES(total_ecpm),
			total_impression_rate = VALUES(total_impression_rate),
			total_click_rate = VALUES(total_click_rate)
		`, miniProgramName, item.Date, item.AdSlotId, item.AdSlotName, item.AdSlotType,
			item.ReqCount, item.ServedCount, item.ClickCount,
			item.Revenue, item.Ecpm, item.ImpressionRate, item.ClickRate)
		if err != nil {
			return err
		}
	}
	return nil
}

func sendLog(msg string) {
	select {
	case logChan <- msg:
	default:
	}
}

func executeFetch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	go func() {
		defer func() {
			logChan <- "DONE"
		}()

		cfg, err := loadConfig()
		if err != nil {
			sendLog(fmt.Sprintf("加载配置失败: %v", err))
			return
		}

		db, err := getDBConnection(cfg)
		if err != nil {
			sendLog(fmt.Sprintf("连接数据库失败: %v", err))
			return
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			sendLog(fmt.Sprintf("数据库连接测试失败: %v", err))
			return
		}
		sendLog("数据库连接成功")

		if err := initDatabase(db); err != nil {
			sendLog(fmt.Sprintf("初始化数据库失败: %v", err))
			return
		}
		sendLog("数据库初始化完成")

		today := time.Now()
		endDate := today.Format("2006-01-02")
		startDate := today.AddDate(0, 0, -7).Format("2006-01-02")

		for i, mp := range cfg.MiniPrograms {
			sendLog(fmt.Sprintf("正在处理小程序: %s (%d/%d)", mp.Name, i+1, len(cfg.MiniPrograms)))

			token, err := getToken(mp.AppID, mp.Secret)
			if err != nil {
				sendLog(fmt.Sprintf("获取Token失败(%s): %v", mp.Name, err))
				continue
			}
			sendLog(fmt.Sprintf("获取Token成功(%s)", mp.Name))

			if err := syncAdUnitList(db, mp.Name, token); err != nil {
				sendLog(fmt.Sprintf("同步广告位列表失败(%s): %v", mp.Name, err))
			} else {
				sendLog(fmt.Sprintf("同步广告位列表成功(%s)", mp.Name))
			}

			actualStartDate := startDate
			if !mp.FirstPullComplete {
				actualStartDate = today.AddDate(0, 0, -30).Format("2006-01-02")
				sendLog(fmt.Sprintf("首次拉取(%s), 同步最近30天数据", mp.Name))
			}

			if err := syncSummaryData(db, mp.Name, token, actualStartDate, endDate, !mp.FirstPullComplete); err != nil {
				sendLog(fmt.Sprintf("同步数据失败(%s): %v", mp.Name, err))
			} else {
				sendLog(fmt.Sprintf("同步数据成功(%s)", mp.Name))
				if !mp.FirstPullComplete {
					mu.Lock()
					for j := range cfg.MiniPrograms {
						if cfg.MiniPrograms[j].AppID == mp.AppID {
							cfg.MiniPrograms[j].FirstPullComplete = true
							break
						}
					}
					saveConfig(cfg)
					mu.Unlock()
				}
			}

			time.Sleep(500 * time.Millisecond)
		}

		sendLog("所有任务完成")
	}()

	for {
		select {
		case msg := <-logChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
			if msg == "DONE" {
				return
			}
		case <-r.Context().Done():
			return
		}
	}
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

func main() {
	http.Handle("/", http.FileServer(http.FS(frontendFS)))
	http.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getConfig(w, r)
		} else if r.Method == http.MethodPost {
			saveConfigHandler(w, r)
		}
	})
	http.HandleFunc("/api/test-connection", testDatabaseConnection)
	http.HandleFunc("/api/execute", executeFetch)

	port := 28384
	addr := fmt.Sprintf(":%d", port)
	url := fmt.Sprintf("http://localhost:%d/frontend/", port)

	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser(url)
	}()

	log.Printf("服务器启动: %s", url)
	log.Fatal(http.ListenAndServe(addr, nil))
}

