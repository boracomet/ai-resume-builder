package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPreviewRequestBinding(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"id":15,"name":"test","photoBase64":"data:image/png;base64,abc","useStoredPhoto":true,"personal":{"name":"Bora"}}`

	var viaStdlib previewRequest
	if err := json.Unmarshal([]byte(body), &viaStdlib); err != nil {
		t.Fatalf("stdlib: %v", err)
	}
	if viaStdlib.ID != 15 {
		t.Fatalf("stdlib id: got %d", viaStdlib.ID)
	}
	if viaStdlib.PhotoBase64 != "data:image/png;base64,abc" {
		t.Fatalf("stdlib photo: got %q", viaStdlib.PhotoBase64)
	}
	if !viaStdlib.UseStoredPhoto {
		t.Fatal("stdlib: UseStoredPhoto should be true")
	}

	r := gin.New()
	var viaGin previewRequest
	r.POST("/preview", func(c *gin.Context) {
		if err := c.ShouldBindJSON(&viaGin); err != nil {
			t.Errorf("gin bind: %v", err)
			c.Status(400)
			return
		}
		c.Status(200)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/preview", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("gin status: %d", w.Code)
	}
	if viaGin.ID != 15 || viaGin.PhotoBase64 != "data:image/png;base64,abc" || !viaGin.UseStoredPhoto {
		t.Fatalf("gin: ID=%d Photo=%q UseStored=%v", viaGin.ID, viaGin.PhotoBase64, viaGin.UseStoredPhoto)
	}
}

func TestPreviewRequestUseStoredWithoutInlinePhoto(t *testing.T) {
	body := `{"id":15,"name":"test","photoBase64":"","useStoredPhoto":true}`
	var req previewRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatal(err)
	}
	if req.ID != 15 || !req.UseStoredPhoto || req.PhotoBase64 != "" {
		t.Fatalf("got ID=%d UseStored=%v Photo=%q", req.ID, req.UseStoredPhoto, req.PhotoBase64)
	}
}
