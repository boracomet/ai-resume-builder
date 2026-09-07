# ATA CV Builder

Go + Gin tabanlı, ATS uyumlu CV hazırlama uygulaması. Canlı önizleme, çoklu profil desteği (SQLite), Türkçe/İngilizce dil seçimi, Google Translate entegrasyonu ve headless Chrome ile PDF export sunar.

## Özellikler

- Split-panel editör: sol form, sağ canlı önizleme (300ms debounce)
- Çoklu CV profili (SQLite, `data/cv.db`)
- CV dili seçimi: Türkçe veya English (bölüm başlıkları ve etiketler otomatik lokalize edilir)
- Google Cloud Translation API ile tek tıkla İngilizceye çeviri
- ATS uyumlu tek sütun HTML şablonu
- PDF export (`chromedp` + headless Chrome)
- Profil fotoğrafı yükleme
- İlk çalıştırmada Bora Ata Türkoğlu CV seed verisi

## Gereksinimler

- Go 1.22+
- Google Chrome veya Chromium (PDF export için)
- Google Translate (opsiyonel): Cloud Translation API anahtarı

## Kurulum ve Çalıştırma

```bash
cd bora-cv
cp .env.example .env
# .env dosyasını düzenleyin
go mod tidy
go run ./cmd/server
```

Tarayıcıda açın: [http://localhost:8080](http://localhost:8080)

## Docker ile Çalıştırma

```bash
cp .env.example .env
# .env dosyasını düzenleyin (API anahtarları vb.)
docker compose up --build
```

Uygulama `http://localhost:8080` adresinde çalışır. Veritabanı `cv-data` Docker volume'ünde saklanır.

Docker imajı PDF export için Chromium içerir (`CHROME_PATH=/usr/bin/chromium-browser`).

## Ortam Değişkenleri

`.env.example` dosyasını `.env` olarak kopyalayıp düzenleyin. Sunucu başlangıcında `godotenv` ile otomatik yüklenir.

| Değişken | Varsayılan | Açıklama |
|----------|------------|----------|
| `PORT` | `8080` | HTTP portu |
| `CV_DB_PATH` | `data/cv.db` | SQLite dosya yolu |
| `GOOGLE_TRANSLATE_API_KEY` | — | Google Cloud Translation API anahtarı |
| `OPENAI_API_KEY` | — | OpenAI API anahtarı (gelecek özellikler) |
| `GOOGLE_GEMINI_API_KEY` | — | Google Gemini API anahtarı (gelecek özellikler) |
| `XIAOMI_API_KEY` | — | Xiaomi API anahtarı (gelecek özellikler) |
| `CHROME_PATH` | — | Chrome/Chromium yolu (Docker'da otomatik) |

### Google Translate API Anahtarı Nasıl Alınır?

1. [Google Cloud Console](https://console.cloud.google.com/) üzerinde bir proje oluşturun
2. **Cloud Translation API**'yi etkinleştirin
3. **APIs & Services → Credentials** bölümünden bir **API Key** oluşturun
4. Anahtarı `.env` dosyasına ekleyin:

```bash
GOOGLE_TRANSLATE_API_KEY=AIza...
```

Alternatif olarak editörde **API Anahtarları** bölümünden anahtarı girebilirsiniz. Tarayıcıda girilen değerler `localStorage`'a kaydedilir ve sunucu anahtarına göre önceliklidir.

## Dil Desteği

Her CV profili için `language` alanı (`tr` veya `en`) saklanır. Araç çubuğundaki **CV Dili** seçicisi ile bölüm başlıkları anında güncellenir:

| Türkçe | English |
|--------|---------|
| ÖZET | SUMMARY |
| DENEYİM | EXPERIENCE |
| EĞİTİM | EDUCATION |
| PROJELER | PROJECTS |
| BECERİLER | SKILLS |

**İngilizceye Çevir** butonu, Google Translate API kullanarak özet, deneyim, eğitim, proje metinlerini ve beceri grup adlarını çevirir; teknik beceri maddeleri olduğu gibi bırakılır. Çeviri sonrası profil dili `en` olarak ayarlanır ve veritabanına kaydedilir.

## API

| Method | Path | Açıklama |
|--------|------|----------|
| GET | `/` | Editör arayüzü |
| GET | `/api/profiles` | Profil listesi |
| POST | `/api/profiles` | Yeni profil |
| GET | `/api/profiles/:id` | Profil detayı |
| PUT | `/api/profiles/:id` | Profil güncelle |
| DELETE | `/api/profiles/:id` | Profil sil |
| GET | `/api/settings` | Yapılandırılmış API anahtar durumu (değerler gönderilmez) |
| POST | `/api/preview` | JSON → HTML önizleme |
| POST | `/api/pdf` | JSON → PDF indir |
| POST | `/api/profiles/:id/photo` | Fotoğraf yükle |
| POST | `/api/profiles/:id/translate` | Profili Google Translate ile çevir |

### Ayarlar Endpoint

```http
GET /api/settings
```

Yanıt örneği:

```json
{
  "googleTranslateConfigured": true,
  "openaiConfigured": false,
  "googleGeminiConfigured": false,
  "xiaomiConfigured": false
}
```

### Çeviri Endpoint

```http
POST /api/profiles/:id/translate
Content-Type: application/json

{
  "targetLang": "en",
  "apiKey": "opsiyonel-ui-override"
}
```

`apiKey` gönderilmezse sunucu `GOOGLE_TRANSLATE_API_KEY` ortam değişkenini kullanır.

## PDF Notları

- PDF üretimi sistemde kurulu Chrome/Chromium gerektirir
- İlk PDF isteği browser başlatma nedeniyle 1-2 saniye sürebilir
- Çok sayfa desteği CSS `@page` ve `page-break` kuralları ile sağlanır

## Build

```bash
go build ./...
go build -o bin/cv-server ./cmd/server
./bin/cv-server
```
