package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/boradev/bora-cv/internal/db"
	"github.com/boradev/bora-cv/internal/handlers"
	"github.com/boradev/bora-cv/web"
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
	previewHandler := handlers.NewPreviewHandler()
	pdfHandler := handlers.NewPDFHandler()
	translateHandler := handlers.NewTranslateHandler(repo)
	settingsHandler := handlers.NewSettingsHandler()

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
		api.POST("/profiles/:id/photo", cvHandler.UploadPhoto)
		api.POST("/profiles/:id/translate", translateHandler.TranslateProfile)
		api.GET("/settings", settingsHandler.GetStatus)
		api.POST("/preview", previewHandler.Preview)
		api.POST("/pdf", pdfHandler.Generate)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("ATA CV Builder http://localhost:%s adresinde çalışıyor", port)
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
	count, err := repo.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	profile := db.SeedProfile()
	return repo.Create(profile)
}
