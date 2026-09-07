package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boradev/bora-cv/internal/models"
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
	return map[string]interface{}{
		"language":         profile.Language,
		"personal":         profile.Personal,
		"contentTR":        profile.ContentTR,
		"contentEN":        profile.ContentEN,
		"photoBase64":      profile.PhotoBase64,
		"photoSize":        profile.PhotoSize,
		"photoBorderWidth": profile.PhotoBorderWidth,
		"photoBorderColor": profile.PhotoBorderColor,
	}
}

func SeedProfile() *models.CVProfile {
	contentTR := models.LocalizedContent{
		PersonalTitle:     "Full Stack Developer",
		PersonalLanguages: "Türkçe (Ana Dil), İngilizce (B2)",
		Summary:           "Ölçeklenebilir web uygulamaları geliştiren bir Full Stack Developer'ım. Temiz kod, sürdürülebilir yazılım mimarileri ve kullanıcı deneyimi odaklı ürünler geliştirmeye önem veriyorum. Gerçek problemleri teknolojiyle çözmeyi, yenilikçi ve etkili yazılım çözümleri üretmeyi hedefliyorum.",
		Experiences: []models.Experience{
			{
				Title:     "IT Web Support Takım Lideri",
				Company:   "Roofstacks - Yazılım Teknolojileri Şirketi",
				StartDate: "Eylül 2025",
				EndDate:   "Şubat 2026",
				Duration:  "6 Ay",
				Location:  "Uzaktan",
				Description: "Roof Stacks, AR/VR, Web3, Metaverse, e-ticaret, oyun ve NFT alanlarında faaliyet gösteren, Austin merkezli uluslararası bir teknoloji şirketidir.",
				Highlights: []string{
					"Web Support biriminin liderliğini üstlenerek üç kişilik ekibi yönettim.",
					"Web altyapısı, CMS, frontend, deployment ve teknik destek süreçlerinde görev aldım.",
					"Dış kaynaklı web hizmetlerinin şirket içine taşınmasını yönettım.",
					".NET tabanlı iç uygulamalar ve intranet çözümleri geliştirdim.",
					"Azure DevOps ve GCP Kubernetes operasyonlarına katkı sağladım.",
				},
			},
			{
				Title:     "Full Stack Developer",
				Company:   "Onlipr - Kreatif Pazarlama Ajansı",
				StartDate: "Mart 2025",
				EndDate:   "Eylül 2025",
				Duration:  "6 Ay",
				Location:  "Uzaktan, Sözleşmeli",
				Description: "Onlipr, Türkiye ve Amerika pazarlarında faaliyet gösteren bir dijital ajans ve teknoloji şirketidir.",
				Highlights: []string{
					"Web projelerinde sunucu yönetimi, frontend geliştirme, CMS geliştirme ve yayına alma süreçlerinde görev aldım.",
					"Projelerin uçtan uca geliştirilmesine katkı sağladım.",
				},
			},
			{
				Title:     "Web Developer",
				Company:   "Segnet Dijital Yazılım",
				StartDate: "Şubat 2024",
				EndDate:   "Mart 2025",
				Duration:  "13 Ay",
				Location:  "Uzaktan",
				Description: "Segnet Digital, İstanbul merkezli bir yazılım şirketidir ve işletmelere ihtiyaçlarına özel kurumsal yazılım çözümleri sunmaktadır.",
				Highlights: []string{
					"Şirket bünyesinde kurumsal web projelerinde görev aldım.",
					"GitHub üzerinden repository ve source control süreçlerini yönettim.",
					"Projelerin Google Cloud Run üzerinde container tabanlı olarak yayına alınması, deployment ve production süreçlerinde aktif rol aldım.",
				},
			},
			{
				Title:     "Web Developer",
				Company:   "Cicoop Baby",
				StartDate: "Ocak 2023",
				EndDate:   "Ocak 2024",
				Duration:  "12 Ay",
				Location:  "Uzaktan, Sözleşmeli",
				Description: "Cicoop Baby'nin kurumsal web altyapısı ve şirket içi stok yönetim sistemlerinin geliştirilmesinde görev aldım.",
				Highlights: []string{
					"PHP, MySQL ve web scraping teknolojilerini kullanarak ürün takibi, stok kontrolü, veri toplama ve süreç otomasyonu için özel web uygulamaları geliştirdim.",
					"Backend geliştirme, veritabanı yönetimi ve iş süreçlerinin dijitalleştirilmesine yönelik çözümler ürettim.",
				},
			},
			{
				Title:     "Junior Web Developer",
				Company:   "Altın Dünyası Yayın Grubu",
				StartDate: "Eylül 2020",
				EndDate:   "Ekim 2021",
				Duration:  "13 Ay",
				Location:  "Ofisten",
				Description: "Uluslararası mücevher sektörüne yönelik devlet destekli Türkiye Mücevher B2B Portalı – turkishjewellery.org projesinde Junior Web Developer olarak görev aldım.",
				Highlights: []string{
					"Frontend geliştirme, responsive web arayüzleri, görsel optimizasyon, performans iyileştirme ve kullanıcı deneyimi çalışmalarında yer aldım.",
					"Kurumsal web platformunun geliştirilmesi ve yayına hazırlanması süreçlerine katkı sağladım.",
				},
			},
		},
		Education: []models.Education{
			{
				Degree:      "Web Tasarımı ve Kodlama",
				School:      "Anadolu Üniversitesi Açıköğretim Fakültesi",
				StartDate:   "2022",
				EndDate:     "Devam Ediyor",
				Description: "Web Tasarımı ve Kodlama önlisans programında eğitimime devam ediyorum. Program kapsamında HTML, CSS, JavaScript, kullanıcı deneyimi (UX), web tabanlı uygulama geliştirme ve dijital yayıncılık alanlarında eğitim alıyorum. Ayrıca müfredat ve uygulamalı çalışmalar kapsamında .NET, Golang, Node.js, React, TypeScript, REST API, frontend ve backend development konularında bilgi ve deneyim kazanıyorum.",
			},
			{
				Degree:      "Grafik Tasarım (Çift Anadal)",
				School:      "Nişantaşı Üniversitesi Meslek Yüksekokulu",
				StartDate:   "2018",
				EndDate:     "2020",
				Description: "Adobe Photoshop, Illustrator ve InDesign kullanarak grafik tasarım, dijital medya, arayüz tasarımı ve görsel iletişim alanlarında eğitim aldım. Web tasarımı, UI/UX ve dijital ürün tasarımı konusunda temel yetkinlikler kazandım.",
			},
			{
				Degree:      "Fotoğrafçılık ve Kameramanlık (Ana Dal)",
				School:      "Nişantaşı Üniversitesi Meslek Yüksekokulu",
				StartDate:   "2017",
				EndDate:     "2019",
				Description: "Fotoğrafçılık, videografi, görsel hikâye anlatımı ve kompozisyon alanlarında eğitim aldım. Adobe Photoshop, Premiere Pro, Lightroom, Illustrator ve InDesign kullanarak görsel düzenleme, video prodüksiyonu ve dijital içerik üretimi alanlarında deneyim kazandım.",
			},
		},
		Projects: []models.Project{},
		SkillGroups: []models.SkillGroup{
			{
				Category: "Frontend",
				Skills: []string{
					"React, Next.js",
					"JavaScript / TypeScript",
					"Vue, Nuxt, Vite, Astro, Handlebars",
					"Responsive ve performans odaklı UI geliştirme",
					"Figma to Code, Pixel-perfect arayüz geliştirme",
					"SSR, SSG ve modern frontend rendering yaklaşımları",
					"Component-based frontend architecture",
					"REST API ve GraphQL API entegrasyonu",
					"Cross-browser uyumluluk",
					"Core Web Vitals ve frontend performans optimizasyonu",
				},
			},
			{
				Category: "Backend",
				Skills: []string{
					"Node.js ekosistemi",
					"TypeScript tabanlı backend geliştirme",
					"Go (Golang) ile backend servis geliştirme",
					".NET / ASP.NET ile kurumsal ve şirket içi web uygulaması geliştirme deneyimi",
					"Python (Web Scraping)",
					"RESTful API ve GraphQL API tasarımı ve geliştirme",
					"Mikroservis ve servis tabanlı mimariler",
					"CRUD tabanlı uygulama geliştirme",
				},
			},
			{
				Category: "Database",
				Skills: []string{
					"PostgreSQL",
					"MongoDB",
					"MySQL",
					"SQLite",
					"Veri modelleme ve temel veritabanı yönetimi",
				},
			},
			{
				Category: "Messaging & Event-Driven",
				Skills: []string{
					"RabbitMQ",
					"Apache Kafka",
					"Message Queue mimarileri",
					"Event-Driven Architecture",
					"Olay tabanlı sistem tasarımı",
					"Asenkron servis iletişimi",
				},
			},
			{
				Category: "DevOps & Infrastructure",
				Skills: []string{
					"Docker",
					"Docker Compose",
					"Kubernetes",
					"Container / Pod yönetimi",
					"CI/CD Pipeline yönetimi",
					"Build, Test ve Production Deployment süreçleri",
					"Google Cloud Platform (GCP)",
					"Google Cloud Run",
					"Azure DevOps",
					"AWS EC2, S3",
					"Nginx, Reverse Proxy",
					"OpenLiteSpeed",
					"Linux sistem ve sunucu yönetimi",
					"Self-hosted PaaS Coolify, Dokploy",
					"Monorepo tabanlı deployment yapıları",
					"VPN arkasında internal service ve yönetim paneli deployment süreçleri",
				},
			},
			{
				Category: "Cloud & Architecture",
				Skills: []string{
					"Cloud-native uygulama mimarileri",
					"PaaS / Serverless deployment",
					"Tekli ve çoklu Container / Pod mimarileri",
					"Intranet ve Internal Web Application geliştirme",
					"On-premise / Private Infrastructure deneyimi",
					"Mimari planlama ve servis bağımlılıklarının tasarımı",
					"Production Release Management",
					"Deployment koordinasyonu",
					"Performans, sürdürülebilirlik ve maliyet optimizasyonu",
					"Mermaid ile sistem ve deployment diyagramları",
				},
			},
			{
				Category: "AI/ML & LLM",
				Skills: []string{
					"Local LLM deployment ve model yönetimi",
					"AI Agent geliştirme",
					"Agent orchestration",
					"RAG mimarileri",
					"Knowledge-based systems",
					"Prompt Engineering",
					"AI destekli otomasyon",
					"Otonom görev yürüten AI sistemleri",
					"OpenClaw, Hermes ve benzeri Agent Framework'leri",
				},
			},
			{
				Category: "Mobile",
				Skills: []string{
					"SwiftUI ile Apple ekosisteminde uygulama geliştirme",
					"React Native ile mobil uygulama geliştirme",
					"Native ve cross-platform mobil uygulama yaklaşımları",
				},
			},
			{
				Category: "Tools & Version Control",
				Skills: []string{
					"Git",
					"GitHub",
					"GitLab",
					"Azure DevOps Repos",
					"Source Control ve Repository Management",
					"CI/CD Pipelines",
					"Agile",
					"Scrum",
					"Kanban",
				},
			},
		},
	}

	return &models.CVProfile{
		Name:     "Varsayılan CV",
		Language: "tr",
		Personal: models.PersonalInfo{
			Name:      "Bora Ata Türkoğlu",
			BirthDate: "Haziran 1997",
			Email:     "ataturkoglubora@gmail.com",
			Phone:     "+90 541 947 44 93",
			Location:  "İstanbul",
			Portfolio: "https://boraturkoglu.com",
			LinkedIn:  "https://linkedin.com/in/boracomet",
			GitHub:    "https://github.com/boracomet",
		},
		ContentTR:        contentTR,
		ContentEN:        models.EmptyLocalizedContent(),
		PhotoSize:        models.DefaultPhotoSize,
		PhotoBorderWidth: models.DefaultPhotoBorderWidth,
		PhotoBorderColor: models.DefaultPhotoBorderColor,
	}
}
