package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/boracomet/ai-resume-builder/internal/db"
	"github.com/boracomet/ai-resume-builder/internal/handlers"
	"github.com/boracomet/ai-resume-builder/web"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	gin.SetMode(gin.ReleaseMode)

	dbPath := os.Getenv("CV_DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join("data", "cv.db")
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("veritabanı açılamadı: %v", err)
	}
	defer conn.Close()

	repo := db.NewCVRepository(conn)
	if err := seedIfEmpty(repo); err != nil {
		log.Fatalf("seed verisi yüklenemedi: %v", err)
	}

	cvHandler := handlers.NewCVHandler(repo)
	previewHandler := handlers.NewPreviewHandler(repo)
	pdfHandler := handlers.NewPDFHandler()
	translateHandler := handlers.NewTranslateHandler(repo)
	aiHandler := handlers.NewAIHandler(repo)
	ocrHandler := handlers.NewOCRHandler()
	settingsHandler := handlers.NewSettingsHandler()
	backupHandler := handlers.NewBackupHandler(repo)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	staticFS, err := fs.Sub(web.Static, "static")
	if err != nil {
		log.Fatalf("static dosyalar yüklenemedi: %v", err)
	}
	router.StaticFS("/static", http.FS(staticFS))

	router.GET("/", serveIndex)

	api := router.Group("/api")
	{
		api.GET("/profiles", cvHandler.ListProfiles)
		api.POST("/profiles", cvHandler.CreateProfile)
		api.GET("/profiles/:id", cvHandler.GetProfile)
		api.PUT("/profiles/:id", cvHandler.UpdateProfile)
		api.DELETE("/profiles/:id", cvHandler.DeleteProfile)
		api.POST("/profiles/:id/duplicate", cvHandler.DuplicateProfile)
		api.POST("/profiles/:id/copy-from/:sourceId", cvHandler.CopyFromProfile)
		api.POST("/profiles/:id/photo", cvHandler.UploadPhoto)
		api.GET("/profiles/:id/export", backupHandler.ExportProfile)
		api.GET("/export", backupHandler.ExportAll)
		api.POST("/import", backupHandler.Import)
		api.POST("/reset", backupHandler.Reset)
		api.POST("/profiles/:id/translate", translateHandler.TranslateProfile)
		api.GET("/ai/models", aiHandler.ListModels)
		api.POST("/ai/test", aiHandler.TestConnection)
		api.POST("/ai/chat", aiHandler.Chat)
		api.POST("/ai/apply", aiHandler.Apply)
		api.POST("/ocr/pdf", ocrHandler.ProcessPDF)
		api.GET("/settings", settingsHandler.GetStatus)
		api.POST("/preview", previewHandler.Preview)
		api.POST("/pdf", pdfHandler.Generate)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("AI Resume Builder http://localhost:%s adresinde çalışıyor", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("sunucu başlatılamadı: %v", err)
	}
}

func serveIndex(c *gin.Context) {
	data, err := web.Static.ReadFile("static/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "index.html bulunamadı")
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}

func seedIfEmpty(repo *db.CVRepository) error {
	return repo.EnsureExampleProfile()
}
