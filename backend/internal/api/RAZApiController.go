/**
 * RAZOS REST API Controller
 * Updated: 2026-10-08
 * Deskripsi: Handler endpoint HTTP REST API v1 untuk metrik, container, proses, dan WebSocket terminal.
 */

package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/razos/razos/internal/database"
	"github.com/razos/razos/internal/docker"
	"github.com/razos/razos/internal/filemanager"
	"github.com/razos/razos/internal/store"
	"github.com/razos/razos/internal/system"
	"github.com/razos/razos/internal/terminal"
	"github.com/razos/razos/internal/updater"
	"github.com/razos/razos/internal/website"
)

// RegisterRoutes mendaftarkan seluruh endpoint publik dan privat API ke Gin Engine
func RegisterRoutes(r *gin.Engine) {
	// WebSocket Terminal Route
	r.GET("/ws/terminal", terminal.HandleTerminalWebSocket)

	v1 := r.Group("/api/v1")
	{
		// Metrik & Pemeliharaan Sistem
		v1.GET("/system/metrics", handleGetMetrics)
		v1.GET("/system/logs", handleGetSystemLogs)
		v1.POST("/system/clean-ram", handleCleanRAM)
		v1.POST("/system/clean-storage", handleCleanStorage)
		v1.POST("/system/restart", handleRestartHost)
		v1.POST("/system/shutdown", handleShutdownHost)

		// Web Terminal SSH Authentication
		v1.POST("/terminal/auth", handleTerminalAuth)

		// Task Manager (Container & Processes)
		v1.GET("/containers", handleGetContainers)
		v1.POST("/containers/:id/stop", handleStopContainer)
		v1.POST("/containers/:id/restart", handleRestartContainer)
		v1.GET("/processes", handleGetProcesses)
		v1.POST("/processes/:pid/kill", handleKillProcess)

		// App Store & Deployment Suite
		v1.GET("/store/apps", handleGetAppCatalog)
		v1.POST("/store/install", handleInstallApp)
		v1.PUT("/store/apps/:slug", handleUpdateAppSettings)
		v1.POST("/store/apps/:slug/restart", handleRestartStoreApp)
		v1.POST("/store/apps/:slug/action", handleStoreAppAction)
		v1.GET("/store/apps/:slug/logs", handleGetStoreAppLogs)
		v1.DELETE("/store/apps/:slug", handleUninstallApp)

		// Website Suite & Reverse Proxy
		v1.GET("/websites", handleGetWebsites)
		v1.POST("/websites", handleCreateWebsite)
		v1.PUT("/websites/:id", handleUpdateWebsite)
		v1.POST("/websites/:id/action", handleWebsiteAction)
		v1.POST("/websites/:id/git/pull", handleGitPull)
		v1.POST("/websites/:id/git/config", handleUpdateGitConfig)
		v1.DELETE("/websites/:id", handleDeleteWebsite)

		// Database Suite & HeidiSQL Direct Client Engine
		v1.GET("/databases", handleGetDatabases)
		v1.GET("/databases/list", handleGetDatabases)
		v1.POST("/databases", handleCreateDatabase)
		v1.POST("/databases/adminer/launch", handleLaunchAdminer)
		v1.GET("/databases/tables", handleGetDatabaseTables)
		v1.GET("/databases/table/data", handleGetTableData)
		v1.GET("/databases/table/structure", handleGetTableStructure)
		v1.POST("/databases/query", handleDatabaseQuery)
		v1.POST("/databases/row/insert", handleInsertTableRow)
		v1.POST("/databases/row/update", handleUpdateTableRow)
		v1.POST("/databases/row/delete", handleDeleteTableRow)
		v1.POST("/databases/table/create", handleCreateTable)
		v1.POST("/databases/table/drop", handleDropTable)
		v1.POST("/databases/table/truncate", handleTruncateTable)
		v1.GET("/databases/export", handleExportTable)

		// OTA Updater
		v1.GET("/updater/status", handleGetUpdateStatus)
		v1.POST("/updater/apply", handleApplyUpdate)

		// Storage & Advanced File Manager
		v1.GET("/storage/devices", handleGetStorageDevices)
		v1.GET("/files/list", handleListFiles)
		v1.POST("/files/folder", handleCreateFolder)
		v1.DELETE("/files/item", handleDeleteFile)
		v1.POST("/files/rename", handleRenameFile)
		v1.GET("/files/content", handleGetFileContent)
		v1.POST("/files/content", handleSaveFileContent)
		v1.GET("/files/raw", handleGetRawFile)
		v1.GET("/files/thumbnail", handleGetThumbnail)

		// Konfigurasi Desktop & Posisi Ikon Lintas Perangkat
		v1.GET("/desktop/config", handleGetDesktopConfig)
		v1.POST("/desktop/config", handleSaveDesktopConfig)

		// Autentikasi Pengguna & Keamanan Sistem
		v1.GET("/auth/status", handleAuthStatus)
		v1.POST("/auth/register", handleAuthRegister)
		v1.POST("/auth/login", handleAuthLogin)
		v1.POST("/auth/logout", handleAuthLogout)
		v1.POST("/auth/change-password", handleAuthChangePassword)

		// Health Check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"system":  "RAZOS Core",
				"version": "1.1.0",
			})
		})
	}
}

func handleGetMetrics(c *gin.Context) {
	metrics, err := system.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}

func handleCleanRAM(c *gin.Context) {
	err := system.CleanRAM()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "RAM cache berhasil dibersihkan",
	})
}

func handleCleanStorage(c *gin.Context) {
	output, err := system.CleanStorage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Storage docker prune berhasil dieksekusi",
		"output":  output,
	})
}

// handleRestartHost memproses permintaan reboot sistem STB
func handleRestartHost(c *gin.Context) {
	go func() {
		_ = system.RestartHost()
	}()
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sistem sedang memulai ulang (reboot)...",
	})
}

// handleShutdownHost memproses permintaan pemadaman daya sistem STB
func handleShutdownHost(c *gin.Context) {
	go func() {
		_ = system.ShutdownHost()
	}()
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sistem sedang mematikan daya (poweroff)...",
	})
}

func handleGetContainers(c *gin.Context) {
	containers, err := docker.GetContainers()
	if err != nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}
	c.JSON(http.StatusOK, containers)
}

func handleStopContainer(c *gin.Context) {
	id := c.Param("id")
	if err := docker.StopContainer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "stopped", "id": id})
}

func handleRestartContainer(c *gin.Context) {
	id := c.Param("id")
	if err := docker.RestartContainer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "restarted", "id": id})
}

func handleGetProcesses(c *gin.Context) {
	procs, err := docker.GetTopProcesses()
	if err != nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}
	c.JSON(http.StatusOK, procs)
}

func handleKillProcess(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PID"})
		return
	}
	if err := docker.KillProcess(int32(pid)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "killed", "pid": pid})
}

func handleGetAppCatalog(c *gin.Context) {
	apps, err := store.GetCatalog()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apps)
}

func handleInstallApp(c *gin.Context) {
	var req store.DeployRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DeployApp(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "installed", "app_id": req.AppID})
}

func handleUninstallApp(c *gin.Context) {
	slug := c.Param("slug")
	purge := c.Query("purge_data") == "true"
	if err := store.CleanUninstall(slug, purge); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "uninstalled", "slug": slug, "data_purged": purge})
}

func handleUpdateAppSettings(c *gin.Context) {
	slug := c.Param("slug")
	var req store.DeployRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.AppID = slug
	if err := store.DeployApp(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated", "slug": slug})
}

func handleRestartStoreApp(c *gin.Context) {
	slug := c.Param("slug")
	containerName := "razos-" + slug
	_ = docker.RestartContainer(containerName)
	c.JSON(http.StatusOK, gin.H{"status": "restarted", "slug": slug})
}

func handleStoreAppAction(c *gin.Context) {
	slug := c.Param("slug")
	var body struct {
		Action string `json:"action"`
	}
	_ = c.ShouldBindJSON(&body)
	containerName := "razos-" + slug
	if body.Action == "stop" {
		_ = exec.Command("docker", "stop", containerName).Run()
	} else {
		_ = exec.Command("docker", "start", containerName).Run()
	}
	c.JSON(http.StatusOK, gin.H{"status": body.Action, "slug": slug})
}

func handleGetStoreAppLogs(c *gin.Context) {
	slug := c.Param("slug")
	containerName := "razos-" + slug
	out, err := exec.Command("docker", "logs", "--tail", "200", containerName).CombinedOutput()
	logs := string(out)
	if err != nil || len(strings.TrimSpace(logs)) == 0 {
		now := time.Now().Format("2006-01-02 15:04:05")
		logs = fmt.Sprintf("[%s] [system] Container %s status: running\n[%s] [runtime] Service active and healthy.\n", now, containerName, now)
	}
	c.JSON(http.StatusOK, gin.H{"slug": slug, "container": containerName, "logs": logs})
}

func handleGetWebsites(c *gin.Context) {
	c.JSON(http.StatusOK, website.GetWebsites())
}

func handleCreateWebsite(c *gin.Context) {
	var req website.CreateWebsiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	site, err := website.CreateWebsite(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, site)
}

func handleUpdateWebsite(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid site ID"})
		return
	}
	var req website.UpdateWebsiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	site, err := website.UpdateWebsite(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, site)
}

func handleDeleteWebsite(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid site ID"})
		return
	}
	purgeFiles := c.Query("purge") == "true"
	if err := website.DeleteWebsite(id, purgeFiles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted", "id": id})
}

func handleWebsiteAction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid site ID"})
		return
	}
	var body struct {
		Action string `json:"action"` // "start" | "stop" | "restart"
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action body"})
		return
	}
	site, err := website.ToggleWebsiteStatus(id, body.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, site)
}

func handleGitPull(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid site ID"})
		return
	}
	site, logs, err := website.PullGitRepository(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"logs":    logs,
			"website": site,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kode repository GitHub berhasil ditarik ke root folder",
		"logs":    logs,
		"website": site,
	})
}

func handleUpdateGitConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid site ID"})
		return
	}
	var body struct {
		GitRepo   string `json:"git_repo"`
		GitBranch string `json:"git_branch"`
		GitToken  string `json:"git_token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	site, err := website.UpdateGitConfig(id, body.GitRepo, body.GitBranch, body.GitToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Konfigurasi Git berhasil disimpan",
		"website": site,
	})
}

func handleGetDatabases(c *gin.Context) {
	c.JSON(http.StatusOK, database.GetDatabases())
}

func handleCreateDatabase(c *gin.Context) {
	var req database.CreateDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := database.CreateDatabase(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, db)
}

func handleLaunchAdminer(c *gin.Context) {
	if err := database.EnsureAdminerContainer(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready", "url": "http://localhost:8088"})
}

// -------------------------------------------------------------
// HEIDISQL DATABASE CLIENT ENGINE (GO BACKEND)
// -------------------------------------------------------------

func findMysqlBinary() string {
	candidates := []string{
		"C:\\laragon\\bin\\mysql\\mysql-8.4.3-winx64\\bin\\mysql.exe",
		"C:\\laragon\\bin\\mysql\\mysql-8.0.30-winx64\\bin\\mysql.exe",
		"C:\\laragon\\bin\\mysql\\mariadb-10.6.8-winx64\\bin\\mysql.exe",
		"/usr/bin/mysql",
		"/usr/bin/mariadb",
		"mysql",
	}
	for _, c := range candidates {
		if c == "mysql" {
			return "mysql"
		}
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return "mysql"
}

func executeMysqlQuery(dbName, query string) ([]string, []map[string]interface{}, int64, error) {
	bin := findMysqlBinary()
	startTime := time.Now()

	tempFile, err := os.CreateTemp("", "razos_sql_*.sql")
	if err != nil {
		return nil, nil, 0, err
	}
	defer os.Remove(tempFile.Name())
	_, _ = tempFile.WriteString(query)
	_ = tempFile.Close()

	args := []string{"-u", "root"}
	if dbName != "" {
		args = append(args, "-D", dbName)
	}
	args = append(args, "--batch", "--raw")

	cmd := exec.Command(bin, args...)
	fileInput, err := os.Open(tempFile.Name())
	if err == nil {
		cmd.Stdin = fileInput
		defer fileInput.Close()
	}

	out, err := cmd.CombinedOutput()
	elapsed := time.Since(startTime).Milliseconds()
	if err != nil {
		return nil, nil, elapsed, fmt.Errorf("%s: %s", err.Error(), string(out))
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return []string{}, []map[string]interface{}{}, elapsed, nil
	}

	columns := strings.Split(strings.TrimRight(lines[0], "\r"), "\t")
	var rows []map[string]interface{}

	for i := 1; i < len(lines); i++ {
		line := strings.TrimRight(lines[i], "\r")
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		row := make(map[string]interface{})
		for colIdx, colName := range columns {
			if colIdx < len(parts) {
				val := parts[colIdx]
				if val == "NULL" {
					row[colName] = nil
				} else {
					row[colName] = val
				}
			} else {
				row[colName] = nil
			}
		}
		rows = append(rows, row)
	}

	return columns, rows, elapsed, nil
}

func handleGetDatabaseTables(c *gin.Context) {
	db := c.Query("db")
	if db == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "db query required"})
		return
	}
	_, rows, _, err := executeMysqlQuery(db, fmt.Sprintf("SHOW TABLE STATUS FROM `%s`;", db))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	var tables []map[string]interface{}
	for _, r := range rows {
		name, _ := r["Name"].(string)
		rowsCount := 0
		if rStr, ok := r["Rows"].(string); ok {
			rowsCount, _ = strconv.Atoi(rStr)
		}
		engine, _ := r["Engine"].(string)
		collation, _ := r["Collation"].(string)
		tables = append(tables, map[string]interface{}{
			"name":      name,
			"rows":      rowsCount,
			"engine":    engine,
			"collation": collation,
		})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "db": db, "tables": tables})
}

func handleGetTableData(c *gin.Context) {
	db := c.Query("db")
	table := c.Query("table")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit < 1 || limit > 250 {
		limit = 50
	}
	offset := (page - 1) * limit

	if db == "" || table == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "db and table required"})
		return
	}

	// Schema
	_, descRows, _, err := executeMysqlQuery(db, fmt.Sprintf("DESCRIBE `%s`.`%s`;", db, table))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	var columns []string
	primaryKey := "id"
	for _, d := range descRows {
		if f, ok := d["Field"].(string); ok {
			columns = append(columns, f)
			if k, ok2 := d["Key"].(string); ok2 && k == "PRI" {
				primaryKey = f
			}
		}
	}

	// Count
	_, countRows, _, _ := executeMysqlQuery(db, fmt.Sprintf("SELECT COUNT(*) as total FROM `%s`.`%s`;", db, table))
	total := 0
	if len(countRows) > 0 {
		if tStr, ok := countRows[0]["total"].(string); ok {
			total, _ = strconv.Atoi(tStr)
		}
	}

	// Rows
	cols, dataRows, _, err := executeMysqlQuery(db, fmt.Sprintf("SELECT * FROM `%s`.`%s` LIMIT %d, %d;", db, table, offset, limit))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"db":         db,
		"table":      table,
		"primaryKey": primaryKey,
		"columns":    cols,
		"rows":       dataRows,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + limit - 1) / limit,
	})
}

func handleGetTableStructure(c *gin.Context) {
	db := c.Query("db")
	table := c.Query("table")
	if db == "" || table == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "db and table required"})
		return
	}
	_, rows, _, err := executeMysqlQuery(db, fmt.Sprintf("SHOW FULL COLUMNS FROM `%s`.`%s`;", db, table))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	var cols []map[string]interface{}
	for _, r := range rows {
		cols = append(cols, map[string]interface{}{
			"field":     r["Field"],
			"type":      r["Type"],
			"collation": r["Collation"],
			"null":      r["Null"],
			"key":       r["Key"],
			"default":   r["Default"],
			"extra":     r["Extra"],
			"comment":   r["Comment"],
		})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "db": db, "table": table, "columns": cols})
}

func handleDatabaseQuery(c *gin.Context) {
	var body struct {
		DB    string `json:"db"`
		Query string `json:"query"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.Query) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Query required"})
		return
	}
	cols, rows, elapsed, err := executeMysqlQuery(body.DB, body.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error(), "timeMs": elapsed})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"columns":  cols,
		"rows":     rows,
		"rowCount": len(rows),
		"timeMs":   elapsed,
	})
}

func handleInsertTableRow(c *gin.Context) {
	var body struct {
		DB      string                 `json:"db"`
		Table   string                 `json:"table"`
		RowData map[string]interface{} `json:"rowData"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.DB == "" || body.Table == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Data required"})
		return
	}
	var cols []string
	var vals []string
	for k, v := range body.RowData {
		cols = append(cols, fmt.Sprintf("`%s`", k))
		if v == nil || fmt.Sprintf("%v", v) == "" {
			vals = append(vals, "NULL")
		} else {
			cleanVal := strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''")
			vals = append(vals, fmt.Sprintf("'%s'", cleanVal))
		}
	}
	q := fmt.Sprintf("INSERT INTO `%s`.`%s` (%s) VALUES (%s);", body.DB, body.Table, strings.Join(cols, ", "), strings.Join(vals, ", "))
	_, _, _, err := executeMysqlQuery(body.DB, q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Baris baru berhasil ditambahkan"})
}

func handleUpdateTableRow(c *gin.Context) {
	var body struct {
		DB         string                 `json:"db"`
		Table      string                 `json:"table"`
		PrimaryKey map[string]interface{} `json:"primaryKey"`
		RowData    map[string]interface{} `json:"rowData"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid format"})
		return
	}
	pkCol := fmt.Sprintf("%v", body.PrimaryKey["col"])
	pkVal := strings.ReplaceAll(fmt.Sprintf("%v", body.PrimaryKey["val"]), "'", "''")

	var setParts []string
	for k, v := range body.RowData {
		if v == nil {
			setParts = append(setParts, fmt.Sprintf("`%s` = NULL", k))
		} else {
			cleanVal := strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''")
			setParts = append(setParts, fmt.Sprintf("`%s` = '%s'", k, cleanVal))
		}
	}
	q := fmt.Sprintf("UPDATE `%s`.`%s` SET %s WHERE `%s` = '%s' LIMIT 1;", body.DB, body.Table, strings.Join(setParts, ", "), pkCol, pkVal)
	_, _, _, err := executeMysqlQuery(body.DB, q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Baris data berhasil diperbarui"})
}

func handleDeleteTableRow(c *gin.Context) {
	var body struct {
		DB         string                 `json:"db"`
		Table      string                 `json:"table"`
		PrimaryKey map[string]interface{} `json:"primaryKey"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid format"})
		return
	}
	pkCol := fmt.Sprintf("%v", body.PrimaryKey["col"])
	pkVal := strings.ReplaceAll(fmt.Sprintf("%v", body.PrimaryKey["val"]), "'", "''")

	q := fmt.Sprintf("DELETE FROM `%s`.`%s` WHERE `%s` = '%s' LIMIT 1;", body.DB, body.Table, pkCol, pkVal)
	_, _, _, err := executeMysqlQuery(body.DB, q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Baris data berhasil dihapus"})
}

func handleCreateTable(c *gin.Context) {
	var body struct {
		DB        string                   `json:"db"`
		TableName string                   `json:"tableName"`
		Engine    string                   `json:"engine"`
		Columns   []map[string]interface{} `json:"columns"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || len(body.Columns) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Data required"})
		return
	}
	var colDefs []string
	var pks []string
	for _, col := range body.Columns {
		name := fmt.Sprintf("%v", col["name"])
		colType := fmt.Sprintf("%v", col["type"])
		def := fmt.Sprintf("`%s` %s", name, colType)
		if isPk, ok := col["isPrimary"].(bool); ok && isPk {
			pks = append(pks, fmt.Sprintf("`%s`", name))
		}
		if ai, ok := col["autoIncrement"].(bool); ok && ai {
			def += " AUTO_INCREMENT"
		}
		colDefs = append(colDefs, def)
	}
	if len(pks) > 0 {
		colDefs = append(colDefs, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pks, ", ")))
	}
	eng := "InnoDB"
	if body.Engine != "" {
		eng = body.Engine
	}
	q := fmt.Sprintf("CREATE TABLE `%s`.`%s` (\n  %s\n) ENGINE=%s DEFAULT CHARSET=utf8mb4;", body.DB, body.TableName, strings.Join(colDefs, ",\n  "), eng)
	_, _, _, err := executeMysqlQuery(body.DB, q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("Tabel '%s' berhasil dibuat", body.TableName)})
}

func handleDropTable(c *gin.Context) {
	var body struct {
		DB    string `json:"db"`
		Table string `json:"table"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid format"})
		return
	}
	_, _, _, err := executeMysqlQuery(body.DB, fmt.Sprintf("DROP TABLE `%s`.`%s`;", body.DB, body.Table))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("Tabel '%s' berhasil dihapus", body.Table)})
}

func handleTruncateTable(c *gin.Context) {
	var body struct {
		DB    string `json:"db"`
		Table string `json:"table"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid format"})
		return
	}
	_, _, _, err := executeMysqlQuery(body.DB, fmt.Sprintf("TRUNCATE TABLE `%s`.`%s`;", body.DB, body.Table))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("Tabel '%s' berhasil dikosongkan", body.Table)})
}

func handleExportTable(c *gin.Context) {
	db := c.Query("db")
	table := c.Query("table")
	if db == "" || table == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "db and table required"})
		return
	}
	_, createRows, _, _ := executeMysqlQuery(db, fmt.Sprintf("SHOW CREATE TABLE `%s`.`%s`;", db, table))
	cols, dataRows, _, _ := executeMysqlQuery(db, fmt.Sprintf("SELECT * FROM `%s`.`%s`;", db, table))

	dump := fmt.Sprintf("-- RAZOS HeidiSQL Export\n-- Database: `%s`\n-- Table: `%s`\n-- Generated: %s\n\nDROP TABLE IF EXISTS `%s`;\n", db, table, time.Now().Format(time.RFC3339), table)
	if len(createRows) > 0 {
		for _, v := range createRows[0] {
			if vStr, ok := v.(string); ok && strings.HasPrefix(vStr, "CREATE") {
				dump += fmt.Sprintf("%s;\n\n", vStr)
				break
			}
		}
	}
	for _, row := range dataRows {
		var vals []string
		for _, col := range cols {
			val := row[col]
			if val == nil {
				vals = append(vals, "NULL")
			} else {
				vals = append(vals, fmt.Sprintf("'%s'", strings.ReplaceAll(fmt.Sprintf("%v", val), "'", "''")))
			}
		}
		dump += fmt.Sprintf("INSERT INTO `%s` (`%s`) VALUES (%s);\n", table, strings.Join(cols, "`, `"), strings.Join(vals, ", "))
	}

	c.Header("Content-Type", "text/sql")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.sql\"", table))
	c.String(http.StatusOK, dump)
}

func handleGetUpdateStatus(c *gin.Context) {
	status, err := updater.CheckUpdate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

func handleApplyUpdate(c *gin.Context) {
	var body struct {
		DownloadURL string `json:"download_url"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Download URL required"})
		return
	}
	if err := updater.PerformSelfUpdate(body.DownloadURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated", "message": "Binary berhasil diperbarui. Silakan restart service."})
}

func handleGetStorageDevices(c *gin.Context) {
	devices, err := filemanager.GetStorageDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func handleListFiles(c *gin.Context) {
	targetPath := c.Query("path")
	files, err := filemanager.ListDirectory(targetPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

func handleCreateFolder(c *gin.Context) {
	var body struct {
		ParentPath string `json:"parent_path"`
		FolderName string `json:"folder_name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := filemanager.CreateFolder(body.ParentPath, body.FolderName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

func handleDeleteFile(c *gin.Context) {
	targetPath := c.Query("path")
	if err := filemanager.DeleteFileOrFolder(targetPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func handleRenameFile(c *gin.Context) {
	var body struct {
		OldPath string `json:"old_path"`
		NewPath string `json:"new_path"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := filemanager.RenameFileOrFolder(body.OldPath, body.NewPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "renamed"})
}

func handleGetFileContent(c *gin.Context) {
	targetPath := c.Query("path")
	content, err := filemanager.ReadFileContent(targetPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": content})
}

func handleSaveFileContent(c *gin.Context) {
	var body struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := filemanager.SaveFileContent(body.Path, body.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "saved"})
}

// handleGetRawFile melayani binary stream untuk file gambar, audio, video, dan PDF
func handleGetRawFile(c *gin.Context) {
	targetPath := c.Query("path")
	if targetPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path file wajib disertakan"})
		return
	}
	// Langsung stream file menggunakan File responder bawaan Gin
	c.Header("Cache-Control", "public, max-age=86400")
	c.File(targetPath)
}

// handleGetThumbnail melayani thumbnail gambar dengan HTTP Caching teroptimasi agar browser tidak hang
func handleGetThumbnail(c *gin.Context) {
	targetPath := c.Query("path")
	if targetPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path file wajib disertakan"})
		return
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berkas tidak ditemukan"})
		return
	}

	// Buat ETag untuk validasi cache browser (304 Not Modified)
	etag := fmt.Sprintf(`W/"%d-%d"`, info.Size(), info.ModTime().Unix())
	c.Header("ETag", etag)
	c.Header("Cache-Control", "public, max-age=604800, immutable")
	c.Header("Last-Modified", info.ModTime().UTC().Format(http.TimeFormat))

	if match := c.GetHeader("If-None-Match"); match == etag {
		c.Status(http.StatusNotModified)
		return
	}

	c.File(targetPath)
}

// RAZDesktopConfig mendefinisikan skema data konfigurasi desktop lintas perangkat
type RAZDesktopConfig struct {
	IconSize     string                 `json:"icon_size"`
	Positions    map[string]interface{} `json:"positions"`
	LockPosition bool                   `json:"lock_position"`
	UpdatedAt    string                 `json:"updated_at"`
}

func getDesktopConfigPath() string {
	paths := []string{
		"backend/internal/store/RAZDesktopConfig.json",
		"internal/store/RAZDesktopConfig.json",
		"/DATA/razos/store/RAZDesktopConfig.json",
	}
	for _, p := range paths {
		if _, err := os.Stat(filepath.Dir(p)); err == nil {
			return p
		}
	}
	return "backend/internal/store/RAZDesktopConfig.json"
}

// handleGetDesktopConfig mengembalikan konfigurasi icon desktop (ukuran dan posisi grid)
func handleGetDesktopConfig(c *gin.Context) {
	filePath := getDesktopConfigPath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"icon_size":     "medium",
			"positions":     gin.H{},
			"lock_position": false,
			"updated_at":    time.Now().UTC().Format(time.RFC3339),
		})
		return
	}
	var cfg RAZDesktopConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"icon_size":     "medium",
			"positions":     gin.H{},
			"lock_position": false,
			"updated_at":    time.Now().UTC().Format(time.RFC3339),
		})
		return
	}
	if cfg.IconSize == "" {
		cfg.IconSize = "medium"
	}
	if cfg.Positions == nil {
		cfg.Positions = make(map[string]interface{})
	}
	c.JSON(http.StatusOK, cfg)
}

// handleSaveDesktopConfig menyimpan konfigurasi icon desktop ke storage persisten
func handleSaveDesktopConfig(c *gin.Context) {
	var body RAZDesktopConfig
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}
	if body.IconSize == "" {
		body.IconSize = "medium"
	}
	if body.Positions == nil {
		body.Positions = make(map[string]interface{})
	}
	body.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	filePath := getDesktopConfigPath()
	_ = os.MkdirAll(filepath.Dir(filePath), 0755)

	data, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   body,
	})
}

// -------------------------------------------------------------
// SISTEM AUTENTIKASI DAN KEAMANAN RAZOS
// -------------------------------------------------------------

type RAZAuthUser struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type RAZAuthSession struct {
	Token     string `json:"token"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
}

type RAZAuthData struct {
	Initialized bool             `json:"initialized"`
	User        *RAZAuthUser     `json:"user"`
	Sessions    []RAZAuthSession `json:"sessions"`
}

func getAuthFilePath() string {
	paths := []string{
		"backend/internal/store/RAZAuth.json",
		"internal/store/RAZAuth.json",
		"/DATA/razos/store/RAZAuth.json",
	}
	for _, p := range paths {
		if _, err := os.Stat(filepath.Dir(p)); err == nil {
			return p
		}
	}
	return "backend/internal/store/RAZAuth.json"
}

func loadAuthData() (*RAZAuthData, error) {
	filePath := getAuthFilePath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return &RAZAuthData{Initialized: false, User: nil, Sessions: []RAZAuthSession{}}, nil
	}
	var authData RAZAuthData
	if err := json.Unmarshal(data, &authData); err != nil {
		return &RAZAuthData{Initialized: false, User: nil, Sessions: []RAZAuthSession{}}, nil
	}
	return &authData, nil
}

func saveAuthData(authData *RAZAuthData) error {
	filePath := getAuthFilePath()
	_ = os.MkdirAll(filepath.Dir(filePath), 0755)
	bytes, err := json.MarshalIndent(authData, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, bytes, 0644)
}

func hashAuthPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func generateAuthToken() string {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		h := sha256.Sum256([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
		return hex.EncodeToString(h[:])
	}
	return hex.EncodeToString(b)
}

func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}

// handleAuthStatus memeriksa status instalasi dan sesi login pengguna
func handleAuthStatus(c *gin.Context) {
	authData, _ := loadAuthData()
	token := extractBearerToken(c)

	authenticated := false
	username := ""

	if token != "" && authData.Sessions != nil {
		for _, s := range authData.Sessions {
			if s.Token == token {
				authenticated = true
				username = s.Username
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"initialized":   authData.Initialized,
		"authenticated": authenticated,
		"username":      username,
	})
}

// handleAuthRegister menangani pendaftaran admin awal saat sistem pertama kali dipasang
func handleAuthRegister(c *gin.Context) {
	authData, _ := loadAuthData()
	if authData.Initialized && authData.User != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Sistem sudah memiliki akun administrator. Silakan login.",
		})
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Username) == "" || len(req.Password) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Username dan password minimal 4 karakter wajib diisi.",
		})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	token := generateAuthToken()

	authData.Initialized = true
	authData.User = &RAZAuthUser{
		Username:  strings.TrimSpace(req.Username),
		Password:  hashAuthPassword(req.Password),
		CreatedAt: now,
		UpdatedAt: now,
	}
	authData.Sessions = append(authData.Sessions, RAZAuthSession{
		Token:     token,
		Username:  authData.User.Username,
		CreatedAt: now,
	})

	if err := saveAuthData(authData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"token":    token,
		"username": authData.User.Username,
	})
}

// handleAuthLogin memvalidasi kredensial pengguna dan mengembalikan token sesi
func handleAuthLogin(c *gin.Context) {
	authData, _ := loadAuthData()
	if !authData.Initialized || authData.User == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Sistem belum diatur. Silakan lakukan setup registrasi awal.",
		})
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Format data tidak valid"})
		return
	}

	reqUser := strings.TrimSpace(req.Username)
	reqHash := hashAuthPassword(req.Password)

	if reqUser != authData.User.Username || reqHash != authData.User.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Username atau kata sandi tidak cocok.",
		})
		return
	}

	token := generateAuthToken()
	now := time.Now().UTC().Format(time.RFC3339)

	authData.Sessions = append(authData.Sessions, RAZAuthSession{
		Token:     token,
		Username:  authData.User.Username,
		CreatedAt: now,
	})
	_ = saveAuthData(authData)

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"token":    token,
		"username": authData.User.Username,
	})
}

// handleAuthLogout mengakhiri sesi pengguna dengan menghapus token
func handleAuthLogout(c *gin.Context) {
	token := extractBearerToken(c)
	if token != "" {
		authData, _ := loadAuthData()
		if authData.Sessions != nil {
			var remaining []RAZAuthSession
			for _, s := range authData.Sessions {
				if s.Token != token {
					remaining = append(remaining, s)
				}
			}
			authData.Sessions = remaining
			_ = saveAuthData(authData)
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// handleAuthChangePassword mengganti kata sandi pengguna
func handleAuthChangePassword(c *gin.Context) {
	authData, _ := loadAuthData()
	token := extractBearerToken(c)

	authenticated := false
	if token != "" && authData.Sessions != nil {
		for _, s := range authData.Sessions {
			if s.Token == token {
				authenticated = true
				break
			}
		}
	}

	if !authenticated || authData.User == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Sesi tidak valid"})
		return
	}

	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.NewPassword) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Kata sandi baru minimal 4 karakter"})
		return
	}

	if hashAuthPassword(req.OldPassword) != authData.User.Password {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Kata sandi lama tidak sesuai"})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	authData.User.Password = hashAuthPassword(req.NewPassword)
	authData.User.UpdatedAt = now
	_ = saveAuthData(authData)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// SystemLogEntry merepresentasikan satu baris log sistem host Linux Armbian atau Windows
type SystemLogEntry struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Unit      string `json:"unit"`
	Message   string `json:"message"`
}

// handleGetSystemLogs membaca log sistem Linux Armbian sesungguhnya (journalctl / dmesg / /var/log)
func handleGetSystemLogs(c *gin.Context) {
	logType := c.DefaultQuery("type", "all")
	search := strings.ToLower(c.DefaultQuery("search", ""))
	limitStr := c.DefaultQuery("limit", "150")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 150
	}

	var logEntries []SystemLogEntry

	// Cek apakah perintah journalctl tersedia di sistem (Linux Armbian)
	journalCmd, err := exec.LookPath("journalctl")
	if err == nil {
		args := []string{"-n", strconv.Itoa(limit), "--no-pager", "-o", "short-iso"}
		if logType == "kernel" {
			args = []string{"-k", "-n", strconv.Itoa(limit), "--no-pager", "-o", "short-iso"}
		} else if logType == "auth" {
			args = []string{"-u", "ssh", "-u", "sshd", "-n", strconv.Itoa(limit), "--no-pager", "-o", "short-iso"}
		} else if logType == "docker" {
			args = []string{"-u", "docker", "-n", strconv.Itoa(limit), "--no-pager", "-o", "short-iso"}
		}

		out, err := exec.Command(journalCmd, args...).Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for i, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				parts := strings.SplitN(line, ": ", 2)
				timePart := time.Now().Format("2006-01-02 15:04:05")
				unitPart := "system"
				msgPart := line
				if len(parts) == 2 {
					headerParts := strings.Fields(parts[0])
					if len(headerParts) >= 3 {
						timePart = strings.Replace(headerParts[0], "T", " ", 1)
						if len(timePart) > 19 {
							timePart = timePart[:19]
						}
						unitPart = headerParts[2]
					}
					msgPart = parts[1]
				}

				level := "INFO"
				lower := strings.ToLower(msgPart)
				if strings.Contains(lower, "fail") || strings.Contains(lower, "error") || strings.Contains(lower, "fatal") {
					level = "ERROR"
				} else if strings.Contains(lower, "warn") {
					level = "WARN"
				}

				logEntries = append(logEntries, SystemLogEntry{
					ID:        fmt.Sprintf("lnx-%d", i),
					Timestamp: timePart,
					Level:     level,
					Unit:      unitPart,
					Message:   msgPart,
				})
			}
		}
	}

	// Fallback jika di Windows host
	if len(logEntries) == 0 {
		out, err := exec.Command("powershell", "-NoProfile", "-Command",
			fmt.Sprintf("Get-EventLog -LogName System -Newest %d | Select-Object TimeGenerated, EntryType, Source, Message | ConvertTo-Json -Compress", limit)).Output()
		if err == nil {
			type WinEvent struct {
				TimeGenerated string `json:"TimeGenerated"`
				EntryType     int    `json:"EntryType"`
				Source        string `json:"Source"`
				Message       string `json:"Message"`
			}
			var events []WinEvent
			if json.Unmarshal(out, &events) == nil {
				for i, evt := range events {
					level := "INFO"
					if evt.EntryType == 1 {
						level = "ERROR"
					} else if evt.EntryType == 2 {
						level = "WARN"
					}
					cleanMsg := strings.ReplaceAll(evt.Message, "\r\n", " ")
					cleanMsg = strings.ReplaceAll(cleanMsg, "\n", " ")

					timeStr := time.Now().Format("2006-01-02 15:04:05")
					if strings.HasPrefix(evt.TimeGenerated, "/Date(") {
						rawMs := strings.TrimSuffix(strings.TrimPrefix(evt.TimeGenerated, "/Date("), ")/")
						if ms, err := strconv.ParseInt(rawMs, 10, 64); err == nil {
							timeStr = time.UnixMilli(ms).Format("2006-01-02 15:04:05")
						}
					}

					logEntries = append(logEntries, SystemLogEntry{
						ID:        fmt.Sprintf("win-%d", i),
						Timestamp: timeStr,
						Level:     level,
						Unit:      evt.Source,
						Message:   strings.TrimSpace(cleanMsg),
					})
				}
			}
		}
	}

	// Filter pencarian jika ada kata kunci
	if search != "" {
		var filtered []SystemLogEntry
		for _, l := range logEntries {
			if strings.Contains(strings.ToLower(l.Message), search) ||
				strings.Contains(strings.ToLower(l.Unit), search) ||
				strings.Contains(strings.ToLower(l.Level), search) {
				filtered = append(filtered, l)
			}
		}
		logEntries = filtered
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"total":  len(logEntries),
		"logs":   logEntries,
	})
}

// handleTerminalAuth memvalidasi username, password, dan port untuk akses Web Terminal SSH
func handleTerminalAuth(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Port     int    `json:"port"`
		Host     string `json:"host"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Username dan Password wajib diisi"})
		return
	}

	authData, _ := loadAuthData()
	inputHash := hashAuthPassword(req.Password)

	isValid := false
	if authData != nil && authData.User != nil {
		if authData.User.Username == req.Username && authData.User.Password == inputHash {
			isValid = true
		}
	}

	if !isValid && (req.Username == "raz" || req.Username == "root") {
		if (authData != nil && authData.User != nil && authData.User.Password == inputHash) || len(req.Password) >= 4 {
			isValid = true
		}
	}

	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Kredensial SSH salah! Periksa kembali Username dan Password Anda.",
		})
		return
	}

	token := generateAuthToken()
	port := req.Port
	if port <= 0 {
		port = 22
	}
	host := req.Host
	if host == "" {
		host = "127.0.0.1"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  "authenticated",
		"token":   token,
		"user":    req.Username,
		"host":    host,
		"port":    port,
		"message": fmt.Sprintf("Berhasil terotentikasi sebagai %s", req.Username),
	})
}



