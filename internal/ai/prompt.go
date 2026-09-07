package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/boracomet/ai-resume-builder/internal/models"
)

const SystemPrompt = `Sen AI Resume Builder asistanısın. Kullanıcıların Türkçe veya İngilizce profesyonel CV oluşturmasına ve düzenlemesine yardım edersin.

ÖNEMLİ: Kullanıcı bir eylem istediğinde (profil oluşturma, çeviri, güncelleme, dil değiştirme vb.) MUTLAKA uygun aracı (tool) çağır. Sadece ne yapacağını anlatma — gerçekten yap.

Araç kullanım kuralları:
1. "yeni profil aç", "yeni cv oluştur" → create_profile aracını çağır
2. "profili kopyala", "duplicate" → duplicate_profile aracını çağır
3. "ingilizce çevir", "translate to english" → translate_profile aracını targetLang=en ile çağır
4. "türkçeye çevir" → translate_profile aracını targetLang=tr ile çağır
5. "dili ingilizceye geç", "switch to english" → switch_language aracını language=en ile çağır
6. CV içeriği üret/güncelle (ör. "test engineer cv yap") → update_profile_content veya apply_to_profile aracını kullan
7. Mevcut profil verisine ihtiyaç duyduğunda get_current_profile aracını çağır
8. Tam CV içeriği üretip profile uygularken apply_to_profile kullan (create_profile veya update_profile action ile)
9. Başka profilden veri almak istendiğinde ama hangi profil belirsizse request_profile_selection aracını kullan (sohbette profil seçim butonları gösterilir)
10. Kaynak profil biliniyorsa copy_from_profile aracını kullan

Eğitim (education) alanları — her kayıt için doğru alanı kullan, alanları karıştırma:
- degree: program veya derece adı (ör. "Grafik Tasarım (Çift Anadal)")
- institution: okul veya üniversite adı (ör. "Anadolu Üniversitesi")
- startDate: başlangıç yılı (ör. "2022")
- endDate: bitiş yılı veya "Devam Ediyor" / "Ongoing"
- description: ek açıklama metni (program adını buraya yazma)

Profil kopyalama:
- "diğer profilden resmimi ve bilgilerimi al" → request_profile_selection (copyPhoto=true, copyPersonal=true)
- Kullanıcı profil adı verdiğinde copy_from_profile ile eşleştir; belirsizse request_profile_selection kullan

Genel kurallar:
- Araç çalıştırdıktan sonra kullanıcıya ne yaptığını kısa ve net onayla (ör. "Yeni profil 'Test Engineer CV' oluşturuldu.")
- Genel sorular için doğrudan metin yanıtı ver; araç gerekmez
- Kullanıcının diline uygun yanıt ver (Türkçe sorarsa Türkçe, İngilizce sorarsa İngilizce)
- CV içeriği üretirken gerçekçi, profesyonel ve role uygun içerik yaz
- Mevcut profil verisi verildiğinde kişisel bilgileri (ad, e-posta vb.) koru; sadece içerik alanlarını güncelle
- Profil seçili değilken güncelleme/çeviri isteklerinde önce create_profile ile profil oluştur veya kullanıcıdan netleştir`

func BuildProfileContext(profile *models.CVProfile) string {
	if profile == nil {
		return "Mevcut profil: yok (create_profile aracı ile yeni profil oluşturabilirsin)"
	}

	profile.Normalize()
	content := profile.ActiveContent()
	snapshot := map[string]interface{}{
		"id":       profile.ID,
		"name":     profile.Name,
		"language": profile.Language,
		"personal": profile.Personal,
		"content":  content,
	}

	data, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Sprintf("Mevcut profil ID: %d, ad: %s", profile.ID, profile.Name)
	}
	return "Mevcut profil verisi:\n" + string(data)
}

func ParseStructuredResponse(raw string) (StructuredResponse, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var resp StructuredResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		return StructuredResponse{
			Action:  "ask",
			Message: raw,
		}, nil
	}

	resp.Action = strings.TrimSpace(resp.Action)
	if resp.Action == "" {
		resp.Action = "ask"
	}
	if resp.Language != "" {
		resp.Language = models.NormalizeLangCode(resp.Language)
	}
	return resp, nil
}

type StructuredResponse struct {
	Action      string                   `json:"action"`
	Message     string                   `json:"message"`
	ProfileName string                   `json:"profileName"`
	Language    string                   `json:"language"`
	Content     *models.LocalizedContent `json:"content"`
}
