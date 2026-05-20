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
	Database      DatabaseConfig      `json:"database"`
	Settings      SettingsConfig      `json:"settings"`
	MiniPrograms  []MiniProgramConfig `json:"miniPrograms"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type SettingsConfig struct {
	StartDate string `json:"startDate"`
}

type MiniProgramConfig struct {
	Name     string `json:"name"`
	AppID    string `json:"appid"`
	AppSecret string `json:"appsecret"`
}

var (
	configPath string
	mu         sync.Mutex
)

func init() {
	if runtime.GOOS == "windows" {
		configPath = filepath.Join(os.Getenv("APPDATA"), "WechatAdConfig", "config.json")
	} else {
		home, _ := os.UserHomeDir()
		configPath = filepath.Join(home, ".wechatadconfig", "config.json")
	}
}

func loadConfig() (Config, error) {
	var cfg Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 返回默认配置
			return Config{
				Database: DatabaseConfig{
					Host: "localhost",
					Port: 3306,
					User: "root",
					Database: "ad_data",
				},
				Settings: SettingsConfig{
					StartDate: "2025-07-01",
				},
				MiniPrograms: []MiniProgramConfig{},
			}, nil
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
	port := 3306
	if cfg.Database.Port > 0 {
		port = cfg.Database.Port
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		port,
		cfg.Database.Database,
	)
	return sql.Open("mysql", dsn)
}

func testDatabaseConnection(w http.ResponseWriter, r *http.Request) {
	var cfg DatabaseConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, _ := loadConfig()
	config.Database = cfg

	db, err := getDBConnection(config)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "连接成功",
	})
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

	// 合并配置
	if newCfg.Database.Host == "" {
		newCfg.Database = oldCfg.Database
	} else {
		if newCfg.Database.Port == 0 {
			newCfg.Database.Port = oldCfg.Database.Port
			if newCfg.Database.Port == 0 {
				newCfg.Database.Port = 3306
			}
		}
	}

	if len(newCfg.MiniPrograms) == 0 {
		newCfg.MiniPrograms = oldCfg.MiniPrograms
	}

	if newCfg.Settings.StartDate == "" {
		newCfg.Settings = oldCfg.Settings
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
		`CREATE TABLE IF NOT EXISTS publisher_adunit_general (
			id INT AUTO_INCREMENT PRIMARY KEY,
			mini_program_name VARCHAR(255) NOT NULL,
			stat_date VARCHAR(8) NOT NULL,
			adslot_id VARCHAR(100),
			adslot_name VARCHAR(255),
			adslot_type VARCHAR(100),
			req_cnt INT,
			show_cnt INT,
			click_cnt INT,
			revenue DECIMAL(10,4),
			ecpm DECIMAL(10,4),
			show_rate DECIMAL(10,4),
			click_rate DECIMAL(10,4),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_stat_date (stat_date),
			INDEX idx_mini_program (mini_program_name),
			UNIQUE KEY unique_summary (mini_program_name, stat_date, adslot_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS publisher_adpos_general (
			id INT AUTO_INCREMENT PRIMARY KEY,
			mini_program_name VARCHAR(255) NOT NULL,
			adslot_id VARCHAR(100),
			adpos_id VARCHAR(100),
			adpos_name VARCHAR(255),
			req_cnt INT,
			show_cnt INT,
			click_cnt INT,
			revenue DECIMAL(10,4),
			stat_date VARCHAR(8) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_stat_date (stat_date),
			UNIQUE KEY unique_adpos (mini_program_name, stat_date, adslot_id, adpos_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS publisher_settlement (
			id INT AUTO_INCREMENT PRIMARY KEY,
			mini_program_name VARCHAR(255) NOT NULL,
			settle_date VARCHAR(8) NOT NULL,
			settle_amount DECIMAL(12,4),
			tax_amount DECIMAL(12,4),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_settle_date (settle_date),
			UNIQUE KEY unique_settle (mini_program_name, settle_date)
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

func syncAdUnitList(db *sql.DB, miniProgramName, accessToken string, logChan chan<- string) error {
	logChan <- fmt.Sprintf("✅ [%s] 正在获取广告位列表...", miniProgramName)
	
	apiURL := fmt.Sprintf("https://api.weixin.qq.com/publisher/stat?action=get_adunit_list&access_token=%s", accessToken)
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("获取广告位列表失败: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Errcode   int `json:"errcode"`
		Errmsg    string `json:"errmsg"`
		AdunitList []struct {
			AdslotId   string `json:"adslot_id"`
			AdslotName string `json:"adslot_name"`
			AdslotType string `json:"adslot_type"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析广告位列表失败: %v", err)
	}

	if result.Errcode != 0 {
		return fmt.Errorf("微信API错误: %s", result.Errmsg)
	}

	logChan <- fmt.Sprintf("✅ [%s] 广告位列表同步完成，共 %d 个广告位", miniProgramName, len(result.AdunitList))
	return nil
}

func syncSummaryData(db *sql.DB, miniProgramName, accessToken, startDate, endDate string, logChan chan<- string) error {
	logChan <- fmt.Sprintf("✅ [%s] 正在获取 %s 至 %s 的数据...", miniProgramName, startDate, endDate)
	
	startDateFormatted := strings.ReplaceAll(startDate, "-", "")
	endDateFormatted := strings.ReplaceAll(endDate, "-", "")
	
	apiURL := fmt.Sprintf("https://api.weixin.qq.com/publisher/stat?action=get_summary&access_token=%s&begin_date=%s&end_date=%s",
		accessToken, startDateFormatted, endDateFormatted)
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("获取汇总数据失败: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Errcode int `json:"errcode"`
		Errmsg  string `json:"errmsg"`
		Data    []struct {
			StatDate    string  `json:"stat_date"`
			AdslotId    string  `json:"adslot_id"`
			AdslotName  string  `json:"adslot_name"`
			AdslotType  string  `json:"adslot_type"`
			ReqCnt      int     `json:"req_cnt"`
			ShowCnt     int     `json:"show_cnt"`
			ClickCnt    int     `json:"click_cnt"`
			Revenue     float64 `json:"revenue"`
			Ecpm        float64 `json:"ecpm"`
			ShowRate    float64 `json:"show_rate"`
			ClickRate   float64 `json:"click_rate"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析汇总数据失败: %v", err)
	}

	if result.Errcode != 0 {
		return fmt.Errorf("微信API错误: %s", result.Errmsg)
	}

	for _, item := range result.Data {
		_, err := db.Exec(`INSERT INTO publisher_adunit_general 
			(mini_program_name, stat_date, adslot_id, adslot_name, adslot_type,
			req_cnt, show_cnt, click_cnt, revenue, ecpm, show_rate, click_rate)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE 
			req_cnt=VALUES(req_cnt), show_cnt=VALUES(show_cnt), click_cnt=VALUES(click_cnt),
			revenue=VALUES(revenue), ecpm=VALUES(ecpm), show_rate=VALUES(show_rate), click_rate=VALUES(click_rate)`,
			miniProgramName, item.StatDate, item.AdslotId, item.AdslotName, item.AdslotType,
			item.ReqCnt, item.ShowCnt, item.ClickCnt, item.Revenue, item.Ecpm, item.ShowRate, item.ClickRate)
		if err != nil {
			logChan <- fmt.Sprintf("❌ [%s] 保存数据 %s 失败: %v", miniProgramName, item.StatDate, err)
		}
	}

	logChan <- fmt.Sprintf("✅ [%s] 数据同步完成，共 %d 条记录", miniProgramName, len(result.Data))
	return nil
}

func executeFetch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	ctx := r.Context()

	// 创建日志通道
	logChan := make(chan string, 100)
	doneChan := make(chan struct{})

	go func() {
		defer close(logChan)
		defer close(doneChan)

		logChan <- "📊 开始执行数据拉取任务..."

		select {
		case <-ctx.Done():
			logChan <- "⚠️ 任务已中断"
			return
		default:
		}

		cfg, err := loadConfig()
		if err != nil {
			logChan <- fmt.Sprintf("❌ 加载配置失败: %v", err)
			return
		}

		if cfg.Database.Host == "" {
			logChan <- "❌ 数据库配置为空，请先配置数据库"
			return
		}

		logChan <- "🔌 正在连接数据库..."
		db, err := getDBConnection(cfg)
		if err != nil {
			logChan <- fmt.Sprintf("❌ 连接数据库失败: %v", err)
			return
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			logChan <- fmt.Sprintf("❌ 数据库连接测试失败: %v", err)
			return
		}
		logChan <- "✅ 数据库连接成功"

		logChan <- "📝 正在初始化数据库表..."
		if err := initDatabase(db); err != nil {
			logChan <- fmt.Sprintf("❌ 初始化数据库失败: %v", err)
			return
		}
		logChan <- "✅ 数据库初始化完成"

		if len(cfg.MiniPrograms) == 0 {
			logChan <- "⚠️ 没有配置小程序，任务结束"
			return
		}

		for i, mp := range cfg.MiniPrograms {
			select {
			case <-ctx.Done():
				logChan <- "⚠️ 任务已中断"
				return
			default:
			}

			logChan <- fmt.Sprintf("📱 正在处理小程序 %s (%d/%d)...", mp.Name, i+1, len(cfg.MiniPrograms))

			logChan <- fmt.Sprintf("🔑 [%s] 正在获取 Access Token...", mp.Name)
			token, err := getToken(mp.AppID, mp.AppSecret)
			if err != nil {
				logChan <- fmt.Sprintf("❌ [%s] 获取Token失败: %v", mp.Name, err)
				continue
			}
			logChan <- fmt.Sprintf("✅ [%s] 获取Token成功", mp.Name)

			if err := syncAdUnitList(db, mp.Name, token, logChan); err != nil {
				logChan <- fmt.Sprintf("❌ [%s] 同步广告位列表失败: %v", mp.Name, err)
			}

			endDate := time.Now().Format("2006-01-02")
			startDate := cfg.Settings.StartDate
			if startDate == "" {
				startDate = "2025-07-01"
			}

			if err := syncSummaryData(db, mp.Name, token, startDate, endDate, logChan); err != nil {
				logChan <- fmt.Sprintf("❌ [%s] 同步数据失败: %v", mp.Name, err)
			}

			time.Sleep(500 * time.Millisecond)
		}

		logChan <- "🎉 所有任务执行完成！"
	}()

	// 实时发送日志
	for {
		select {
		case msg, ok := <-logChan:
			if !ok {
				return
			}
			fmt.Fprintf(w, "%s\n", msg)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-ctx.Done():
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
	// 配置路由
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

	// 自动打开浏览器
	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser(url)
	}()

	log.Printf("🚀 服务器已启动: %s", url)
	log.Printf("💡 按 Ctrl+C 停止服务器")
	log.Fatal(http.ListenAndServe(addr, nil))
}
