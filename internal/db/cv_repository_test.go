package db

import (
	"path/filepath"
	"testing"

	"github.com/boracomet/ai-resume-builder/internal/models"
)

func TestCVRepositoryCopyFrom(t *testing.T) {
	conn, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })

	repo := NewCVRepository(conn)

	source := models.CVProfile{
		Name:     "Kaynak",
		Language: "tr",
		Personal: models.PersonalInfo{Name: "Kaynak Kişi", Email: "kaynak@example.com"},
		ContentTR: models.LocalizedContent{
			Education: []models.Education{{
				Degree:      "Grafik Tasarım",
				Institution: "Nişantaşı Üniversitesi",
				StartDate:   "2018",
				EndDate:     "2020",
				Description: "Açıklama",
			}},
		},
		PhotoBase64: "data:image/png;base64,abc",
	}
	source.Normalize()
	if err := repo.Create(&source); err != nil {
		t.Fatal(err)
	}

	target := models.CVProfile{Name: "Hedef", Language: "tr"}
	target.Normalize()
	if err := repo.Create(&target); err != nil {
		t.Fatal(err)
	}

	updated, err := repo.CopyFrom(target.ID, source.ID, models.ProfileCopyOptions{
		CopyPhoto:    true,
		CopyPersonal: true,
	})
	if err != nil {
		t.Fatalf("CopyFrom() error = %v", err)
	}

	if updated.Personal.Name != "Kaynak Kişi" {
		t.Fatalf("personal not copied: %+v", updated.Personal)
	}
	if updated.PhotoBase64 != source.PhotoBase64 {
		t.Fatal("photo not copied")
	}
}

func TestEnsureExampleProfile(t *testing.T) {
	conn, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })

	repo := NewCVRepository(conn)

	if err := repo.EnsureExampleProfile(); err != nil {
		t.Fatalf("EnsureExampleProfile() error = %v", err)
	}

	count, err := repo.Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 profile after seed, got %d", count)
	}

	profile, err := repo.GetByID(1)
	if err != nil {
		t.Fatal(err)
	}
	if profile == nil || profile.Name != "Full Stack Developer" {
		t.Fatalf("unexpected seeded profile: %+v", profile)
	}

	if err := repo.EnsureExampleProfile(); err != nil {
		t.Fatalf("EnsureExampleProfile() second call error = %v", err)
	}
	count, err = repo.Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected seed to be idempotent, got %d profiles", count)
	}
}

func TestSeedProfileExample(t *testing.T) {
	profile := SeedProfile()
	if profile == nil {
		t.Fatal("SeedProfile() returned nil")
	}
	if profile.Name != "Full Stack Developer" {
		t.Fatalf("unexpected seed name: %q", profile.Name)
	}
	if profile.Personal.Phone != "" {
		t.Fatalf("seed phone should be empty, got %q", profile.Personal.Phone)
	}
	if len(profile.ContentEN.Experiences) == 0 {
		t.Fatal("seed contentEN should have experiences")
	}
}

func TestCVRepositoryImportProfilesMerge(t *testing.T) {
	conn, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })

	repo := NewCVRepository(conn)

	existing := models.CVProfile{Name: "Full Stack Developer", Language: "tr"}
	existing.Normalize()
	if err := repo.Create(&existing); err != nil {
		t.Fatal(err)
	}

	imported := models.CVProfile{Name: "Full Stack Developer", Language: "en"}
	imported.Normalize()
	result, err := repo.ImportProfiles([]models.CVProfile{imported}, "merge")
	if err != nil {
		t.Fatalf("ImportProfiles() error = %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 imported profile, got %d", len(result))
	}
	if result[0].Name != "Full Stack Developer (import)" {
		t.Fatalf("unexpected import name: %q", result[0].Name)
	}

	count, err := repo.Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("expected 2 profiles after merge import, got %d", count)
	}
}

func TestCVRepositoryResetToExample(t *testing.T) {
	conn, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })

	repo := NewCVRepository(conn)

	first := models.CVProfile{Name: "Profil 1", Language: "tr"}
	first.Normalize()
	if err := repo.Create(&first); err != nil {
		t.Fatal(err)
	}

	second := models.CVProfile{Name: "Profil 2", Language: "en"}
	second.Normalize()
	if err := repo.Create(&second); err != nil {
		t.Fatal(err)
	}

	profile, err := repo.ResetToExample()
	if err != nil {
		t.Fatalf("ResetToExample() error = %v", err)
	}
	if profile == nil || profile.Name != "Full Stack Developer" {
		t.Fatalf("unexpected reset profile: %+v", profile)
	}

	count, err := repo.Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 profile after reset, got %d", count)
	}
}
