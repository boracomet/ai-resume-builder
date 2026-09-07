package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/boracomet/ai-resume-builder/internal/db/seed"
	"github.com/boracomet/ai-resume-builder/internal/models"
)

type CVRepository struct {
	db *sql.DB
}

func NewCVRepository(db *sql.DB) *CVRepository {
	return &CVRepository{db: db}
}

func (r *CVRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM profiles`).Scan(&count)
	return count, err
}

// EnsureExampleProfile seeds the embedded example CV when the database is empty.
func (r *CVRepository) EnsureExampleProfile() error {
	count, err := r.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	profile := SeedProfile()
	return r.Create(profile)
}

func (r *CVRepository) List() ([]models.ProfileSummary, error) {
	rows, err := r.db.Query(`
		SELECT id, name, updated_at FROM profiles ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []models.ProfileSummary
	for rows.Next() {
		var p models.ProfileSummary
		var updatedAt string
		if err := rows.Scan(&p.ID, &p.Name, &updatedAt); err != nil {
			return nil, err
		}
		p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		profiles = append(profiles, p)
	}
	return profiles, rows.Err()
}

func (r *CVRepository) GetByID(id int64) (*models.CVProfile, error) {
	var name, data, createdAt, updatedAt string
	err := r.db.QueryRow(`
		SELECT name, data, created_at, updated_at FROM profiles WHERE id = ?
	`, id).Scan(&name, &data, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var profile models.CVProfile
	if err := json.Unmarshal([]byte(data), &profile); err != nil {
		return nil, fmt.Errorf("unmarshal profile: %w", err)
	}

	profile.ID = id
	profile.Name = name
	profile.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	profile.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	profile.Normalize()
	return &profile, nil
}

func (r *CVRepository) Create(profile *models.CVProfile) error {
	now := time.Now().UTC()
	profile.CreatedAt = now
	profile.UpdatedAt = now

	data, err := json.Marshal(profilePayload(profile))
	if err != nil {
		return err
	}

	result, err := r.db.Exec(`
		INSERT INTO profiles (name, data, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, profile.Name, string(data), now.Format(time.RFC3339), now.Format(time.RFC3339))
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	profile.ID = id
	return nil
}

func (r *CVRepository) Update(profile *models.CVProfile) error {
	now := time.Now().UTC()
	profile.UpdatedAt = now

	data, err := json.Marshal(profilePayload(profile))
	if err != nil {
		return err
	}

	result, err := r.db.Exec(`
		UPDATE profiles SET name = ?, data = ?, updated_at = ? WHERE id = ?
	`, profile.Name, string(data), now.Format(time.RFC3339), profile.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CVRepository) Duplicate(id int64) (*models.CVProfile, error) {
	original, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if original == nil {
		return nil, sql.ErrNoRows
	}

	data, err := json.Marshal(profilePayload(original))
	if err != nil {
		return nil, err
	}

	var copy models.CVProfile
	if err := json.Unmarshal(data, &copy); err != nil {
		return nil, fmt.Errorf("unmarshal duplicate: %w", err)
	}

	copy.Name = duplicateProfileName(original.Name)
	if err := r.Create(&copy); err != nil {
		return nil, err
	}
	return &copy, nil
}

func (r *CVRepository) CopyFrom(targetID, sourceID int64, opts models.ProfileCopyOptions) (*models.CVProfile, error) {
	if targetID == sourceID {
		return nil, fmt.Errorf("kaynak ve hedef profil aynı olamaz")
	}

	target, err := r.GetByID(targetID)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, sql.ErrNoRows
	}

	source, err := r.GetByID(sourceID)
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, fmt.Errorf("kaynak profil bulunamadı: %d", sourceID)
	}

	if opts.CopyPhoto {
		target.PhotoBase64 = source.PhotoBase64
		target.PhotoSize = source.PhotoSize
		target.PhotoBorderWidth = source.PhotoBorderWidth
		target.PhotoBorderColor = source.PhotoBorderColor
	}
	if opts.CopyPersonal {
		target.Personal = source.Personal
	}

	copyContent := opts.CopyContentTR || opts.CopyContentEN || len(opts.Sections) > 0
	if copyContent {
		if opts.CopyContentTR || len(opts.Sections) > 0 {
			target.ContentTR = copyLocalizedContent(target.ContentTR, source.ContentTR, opts, true)
		}
		if opts.CopyContentEN || len(opts.Sections) > 0 {
			target.ContentEN = copyLocalizedContent(target.ContentEN, source.ContentEN, opts, false)
		}
	}

	target.Normalize()
	if err := r.Update(target); err != nil {
		return nil, err
	}
	return r.GetByID(targetID)
}

func copyLocalizedContent(dst, src models.LocalizedContent, opts models.ProfileCopyOptions, isTR bool) models.LocalizedContent {
	copyAll := opts.CopiesAllContent()
	if !copyAll && len(opts.Sections) == 0 {
		if isTR && opts.CopyContentTR {
			copyAll = true
		}
		if !isTR && opts.CopyContentEN {
			copyAll = true
		}
	}

	shouldCopy := func(section string) bool {
		return copyAll || opts.WantsSection(section)
	}

	if shouldCopy("summary") {
		dst.Summary = src.Summary
	}
	if shouldCopy("personalTitle") {
		dst.PersonalTitle = src.PersonalTitle
	}
	if shouldCopy("personalLanguages") {
		dst.PersonalLanguages = src.PersonalLanguages
	}
	if shouldCopy("experience") || shouldCopy("experiences") {
		dst.Experiences = append([]models.Experience(nil), src.Experiences...)
	}
	if shouldCopy("education") {
		dst.Education = append([]models.Education(nil), src.Education...)
	}
	if shouldCopy("projects") {
		dst.Projects = append([]models.Project(nil), src.Projects...)
	}
	if shouldCopy("skills") || shouldCopy("skillGroups") {
		dst.SkillGroups = append([]models.SkillGroup(nil), src.SkillGroups...)
	}
	return dst
}

func duplicateProfileName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Varsayılan CV (Kopya)"
	}
	return name + " - Kopya"
}

func (r *CVRepository) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM profiles WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CVRepository) ExportAll() ([]models.CVProfile, error) {
	summaries, err := r.List()
	if err != nil {
		return nil, err
	}

	profiles := make([]models.CVProfile, 0, len(summaries))
	for _, summary := range summaries {
		profile, err := r.GetByID(summary.ID)
		if err != nil {
			return nil, err
		}
		if profile != nil {
			profiles = append(profiles, *profile)
		}
	}
	return profiles, nil
}

func (r *CVRepository) ImportProfiles(profiles []models.CVProfile, mode string) ([]models.CVProfile, error) {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode != "replace" {
		mode = "merge"
	}

	if mode == "replace" {
		if err := r.DeleteAll(); err != nil {
			return nil, err
		}
	}

	usedNames := map[string]bool{}
	if mode == "merge" {
		existing, err := r.ExportAll()
		if err != nil {
			return nil, err
		}
		for _, p := range existing {
			usedNames[p.Name] = true
		}
	}

	imported := make([]models.CVProfile, 0, len(profiles))
	for _, profile := range profiles {
		profile.ID = 0
		profile.CreatedAt = time.Time{}
		profile.UpdatedAt = time.Time{}
		profile.Normalize()

		if mode == "merge" {
			profile.Name = uniqueImportName(profile.Name, usedNames)
		}

		if err := r.Create(&profile); err != nil {
			return nil, err
		}
		imported = append(imported, profile)
	}
	return imported, nil
}

func (r *CVRepository) DeleteAll() error {
	_, err := r.db.Exec(`DELETE FROM profiles`)
	return err
}

// ResetToExample deletes all profiles and seeds the embedded example CV.
func (r *CVRepository) ResetToExample() (*models.CVProfile, error) {
	if err := r.DeleteAll(); err != nil {
		return nil, err
	}
	if err := r.EnsureExampleProfile(); err != nil {
		return nil, err
	}
	profiles, err := r.List()
	if err != nil {
		return nil, err
	}
	if len(profiles) == 0 {
		return nil, fmt.Errorf("örnek profil oluşturulamadı")
	}
	return r.GetByID(profiles[0].ID)
}

func uniqueImportName(name string, used map[string]bool) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Imported Profile"
	}
	if !used[name] {
		used[name] = true
		return name
	}

	candidate := name + " (import)"
	if !used[candidate] {
		used[candidate] = true
		return candidate
	}

	for i := 2; ; i++ {
		candidate = fmt.Sprintf("%s (import %d)", name, i)
		if !used[candidate] {
			used[candidate] = true
			return candidate
		}
	}
}

func (r *CVRepository) UpdatePhoto(id int64, photoBase64 string) (*models.CVProfile, error) {
	profile, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, sql.ErrNoRows
	}

	profile.PhotoBase64 = photoBase64
	if err := r.Update(profile); err != nil {
		return nil, err
	}
	return profile, nil
}

func profilePayload(profile *models.CVProfile) map[string]interface{} {
	profile.Normalize()
	payload := map[string]interface{}{
		"language":         profile.Language,
		"personal":         profile.Personal,
		"contentTR":        profile.ContentTR,
		"contentEN":        profile.ContentEN,
		"photoBase64":      profile.PhotoBase64,
		"photoSize":        profile.PhotoSize,
		"photoBorderWidth": profile.PhotoBorderWidth,
		"photoBorderColor": profile.PhotoBorderColor,
	}
	if len(profile.ContentLanguages) > 0 {
		payload["contentLanguages"] = profile.ContentLanguages
	}
	return payload
}

func SeedProfile() *models.CVProfile {
	return seed.ExampleProfile()
}
