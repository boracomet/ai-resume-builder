const APP_LANGUAGE_KEY = "appLanguage";
const APP_SETUP_COMPLETE_KEY = "appSetupComplete";
const TARGET_TRANSLATE_LANG_KEY = "targetTranslateLang";
const TARGET_TRANSLATE_LANGS_KEY = "targetTranslateLangs";
const TARGET_TRANSLATE_LANG_MANUAL_KEY = "targetTranslateLangManual";
const ACTIVE_TARGET_TRANSLATE_LANG_KEY = "activeTargetTranslateLang";

const translations = {
  tr: {
    // Toolbar
    save: "Kaydet",
    pdfDownload: "PDF İndir",
    pdfPreparing: "PDF hazırlanıyor...",
    aiAssistant: "AI Asistan",
    darkMode: "Karanlık mod",
    lightMode: "Açık mod",
    selectProfile: "Profil seç",
    newProfile: "Yeni profil",
    duplicateProfile: "Profili Kopyala",
    deleteProfile: "Profili sil",
    editingLanguage: "Düzenleme Dili",
    translateToEn: "İngilizceye Çevir",
    translateToTr: "Türkçeye Çevir",
    translating: "Çevriliyor...",
    languageEmptyHint: "İngilizce içerik henüz oluşturulmadı",
    livePreview: "Canlı Önizleme",
    cvPreview: "CV Önizleme",
    close: "Kapat",
    copy: "Kopyala",

    // Sidebar nav
    navProfile: "Profil",
    navSettings: "Ayarlar",
    navAppSettings: "Uygulama Ayarları",
    navDataBackup: "Veri Yedekleme",
    navGroupCv: "CV",
    navGroupApp: "Uygulama",
    settingsTitle: "Ayarlar",
    aiAssistantSettings: "AI Asistan",
    aiAssistantEnabled: "AI Asistanı göster",
    aiAssistantEnabledHint: "Kapatıldığında sol menüdeki AI Asistan sekmesi ve üstteki AI butonu gizlenir.",
    aiAssistantDisabledToast: "AI Asistan ayarlardan kapalı. Ayarlar'dan açabilirsiniz.",
    navPersonal: "Kişisel Bilgiler",
    navSummary: "Özet",
    navExperience: "Deneyim",
    navEducation: "Eğitim",
    navProjects: "Projeler",
    navSkills: "Beceriler",
    navPhoto: "Profil Fotoğrafı",
    navAI: "AI Asistan",
    formSections: "Form bölümleri",

    // Profile & app settings
    profileSettings: "Profil",
    profileSettingsHint: "Seçili CV profilinin adını düzenleyin.",
    appSettings: "Uygulama Ayarları",
    profileName: "Profil Adı",
    profileNamePlaceholder: "Örn: İş Başvurusu",
    generalSettingsHint: "Düzenleme dili üst araç çubuğundan yönetilir.",
    appLanguage: "Uygulama Dili",
    appLanguageHint: "Editör arayüzünün dilini seçin. CV içeriği dili ayrıca yönetilir.",
    targetTranslateLang: "Hedef Çeviri Dili",
    targetTranslateLangHint: "Araç çubuğundaki çeviri butonunun hangi dile çevireceğini belirler. İlk eklenen dil birincil hedeftir. Birden fazla dil varsa araç çubuğundan aktif hedefi seçebilirsiniz. Eklenen dillerin yanındaki × ile kaldırabilirsiniz; liste boşsa çeviri butonu gizlenir.",
    addLanguage: "Ekle",
    searchLanguages: "Dil ara veya seç...",
    removeLanguage: "Kaldır",
    languageAlreadyAdded: "Bu dil zaten eklenmiş",
    languageNotFound: "Dil bulunamadı",
    translateToLang: "Çevir: {lang}",
    cvTranslatedTo: "CV {lang} diline çevrildi",
    cvTranslatedStorageWarning: "Çeviri tamamlandı ancak CV yalnızca Türkçe/İngilizce slotlarını destekler. Sonuç kaydedilmedi.",
    langTurkish: "Türkçe",
    langEnglish: "English",

    // Personal info
    personalInfo: "Kişisel Bilgiler",
    fullName: "Ad Soyad",
    title: "Unvan",
    email: "E-posta",
    phone: "Telefon",
    linkedin: "LinkedIn",
    github: "GitHub",
    portfolio: "Portfolyo",
    location: "Konum",
    birthDate: "Doğum Tarihi",
    languages: "Diller",

    // Sections
    summary: "Özet",
    summaryPlaceholder: "Kısa profesyonel özet",
    experience: "Deneyim",
    education: "Eğitim",
    projects: "Projeler",
    skills: "Beceriler",
    profilePhoto: "Profil Fotoğrafı",
    add: "+ Ekle",
    addGroup: "+ Grup Ekle",
    remove: "Kaldır",
    delete: "Sil",

    // Photo settings
    photoPreviewAlt: "Profil fotoğrafı önizleme",
    size: "Boyut",
    borderWidth: "Çerçeve kalınlığı",
    borderColor: "Çerçeve rengi",

    // Dynamic card fields
    company: "Şirket",
    startDate: "Başlangıç",
    endDate: "Bitiş",
    duration: "Süre",
    companyDescription: "Şirket Açıklaması",
    highlights: "Maddeler (her satır bir madde)",
    program: "Program",
    institution: "Kurum",
    description: "Açıklama",
    projectName: "Proje Adı",
    url: "URL",
    category: "Kategori",
    skillsPerLine: "Beceriler (her satır bir beceri)",
    experienceN: "Deneyim {n}",
    educationN: "Eğitim {n}",
    projectN: "Proje {n}",
    skillGroupN: "Beceri Grubu {n}",
    moveUp: "Yukarı taşı",
    moveDown: "Aşağı taşı",
    dragToReorder: "Sürükleyerek sırala",

    // AI chat
    apiSetup: "API Kurulum",
    sessionCost: "Oturum Maliyeti:",
    resetSessionCost: "Sıfırla",
    resetSessionCostTitle: "Oturum maliyetini sıfırla",
    lastRequest: "Son istek:",
    send: "Gönder",
    sending: "Gönderiliyor...",
    apply: "Uygula",
    applying: "Uygulanıyor...",
    addPdf: "PDF Ekle",
    processing: "İşleniyor...",
    aiChatPlaceholder: "Örn: bana Test engineer olacağım şekilde cv yap",
    aiChatSection: "AI Sohbet",
    aiToolbar: "AI araç çubuğu",
    apiKeys: "API Anahtarları",
    apiKeyPlaceholder: "API anahtarı",
    aiSettings: "AI Ayarları",
    openaiModel: "OpenAI Model",
    fetchModels: "Modelleri Getir",
    loading: "Yükleniyor...",
    testConnection: "Bağlantıyı Test Et",
    testing: "Test ediliyor...",
    translateProvider: "Çeviri Sağlayıcısı",
    googleTranslate: "Google Translate",
    openai: "OpenAI",
    apiKeysHint: "Anahtarlar sunucuda .env ile yapılandırılabilir. Tarayıcıda girilen değerler localStorage'a kaydedilir ve sunucu anahtarına göre önceliklidir.",
    apiKeyLocal: "Tarayıcıda kayıtlı (öncelikli)",
    apiKeyEnv: ".env ile yapılandırıldı",
    apiKeyEmpty: "Yapılandırılmadı",
    thinking: "Düşünüyor",
    showText: "Metni Göster",
    hideText: "Metni Gizle",
    applyToCv: "CV'ye Uygula",
    modelsLoading: "Modeller yükleniyor...",
    modelsFetching: "Modeller getiriliyor...",
    modelsAlreadyLoading: "Modeller zaten yükleniyor...",
    modelsLoaded: "{count} model yüklendi",
    defaultModelUsed: "Varsayılan model kullanılıyor",
    defaultSuffix: "(varsayılan)",
    modelSelectorsNotFound: "Model seçiciler bulunamadı",
    openaiNotConfigured: "OpenAI yapılandırılmadı (.env veya arayüzden API anahtarı girin)",
    openaiKeyRequired: "OpenAI API anahtarı gerekli (.env veya arayüzden)",
    connectionSuccess: "Bağlantı başarılı",
    openaiConnectionSuccess: "OpenAI bağlantısı başarılı",
    onlyPdfAllowed: "Yalnızca PDF dosyaları yüklenebilir",
    pdfTooLarge: "PDF dosyası 10 MB'dan büyük olamaz",
    pdfProcessing: "PDF işleniyor...",
    pdfPagesOcr: "{pages} sayfa OCR ile okundu",
    pdfPagesExtracted: "{pages} sayfa metin çıkarıldı",
    pdfProcessed: "PDF işlendi: {pages} sayfa, {chars} karakter",
    textAddedToChat: "Metin sohbete eklendi. OpenAI anahtarı ile gönderin veya metni kopyalayın.",
    noContentToApply: "Uygulanacak içerik yok",
    aiChangesApplied: "AI değişiklikleri uygulandı",
    selectProfileFirst: "Önce bir profil seçin",
    profileDataCopied: "Profil verileri kopyalandı",
    whichProfile: "Hangi profilden almak istiyorsunuz?",
    noOtherProfiles: "Kopyalanabilecek başka profil bulunamadı.",
    dataCopiedFromProfile: "Seçilen profilden veriler kopyalandı.",
    noResponse: "Yanıt alınamadı",
    error: "Hata",
    action: "İşlem",
    ocrMethod: "OCR (Tesseract)",
    pdfTextLayer: "PDF metin katmanı",
    pdfUploaded: "[CV PDF yüklendi: {name}]",
    pdfMeta: "Sayfa: {pages} | Karakter: {chars} | Yöntem: {method}",
    pdfTextHeader: "--- PDF METNİ ---",
    pdfExtractedSummary: "PDF'den {pages} sayfa metin çıkarıldı ({chars} karakter, {method}).",
    fileLabel: "Dosya: {name}",
    applyPdfPrompt: "Yukarıdaki PDF metnini mevcut CV profilime uygula. Eksik alanları mantıklı şekilde doldur.",

    // Translate progress
    translateProgress1: "Çevirinizi yapıyorum...",
    translateProgress2: "Özet çevriliyor...",
    translateProgress3: "Deneyimler çevriliyor...",
    translateProgress4: "Eğitim çevriliyor...",
    translateProgress5: "Beceriler çevriliyor...",
    translateProgress6: "Kaydediliyor...",
    translateComplete: "Çeviri tamamlandı!",
    translateCompleteCost: "Çeviri tamamlandı! Maliyet: {cost}{tokens}",
    translateCompleteGoogle: "Çeviri tamamlandı! (Google Translate — ücretsiz API kotası)",
    translateFailed: "Çeviri başarısız: {error}",
    costLabel: "Maliyet: {cost}",
    tokenSuffix: " ({prompt} + {completion} token)",
    unknownError: "Bilinmeyen hata",

    // AI action labels
    actionCreateProfile: "Yeni profil oluşturuluyor",
    actionDuplicateProfile: "Profil kopyalanıyor",
    actionUpdateContent: "Profil içeriği güncelleniyor",
    actionTranslateProfile: "Profil çevriliyor",
    actionSwitchLanguage: "Dil değiştiriliyor",
    actionGetProfile: "Profil okunuyor",
    actionApplyToProfile: "İçerik profile uygulanıyor",
    actionRequestSelection: "Profil seçimi gösteriliyor",
    actionCopyFromProfile: "Profilden veri kopyalanıyor",
    profileCreated: "Profil oluşturuldu: {name}",
    profileUpdated: "Profil güncellendi.",

    // Toasts
    profileSaved: "Profil kaydedildi",
    newProfileCreated: "Yeni profil oluşturuldu",
    profileDuplicated: "Profil kopyalandı",
    profileDeleted: "Profil silindi",
    profilesLoadError: "Profiller yüklenemedi. Sayfayı yenileyin.",
    confirmDelete: "Bu profili silmek istediğinize emin misiniz?",
    enterTrContent: "Önce Türkçe içerik girin",
    enterEnContent: "Önce İngilizce içerik girin",
    cvTranslatedEn: "CV İngilizceye çevrildi",
    cvTranslatedTr: "CV Türkçeye çevrildi",
    pdfDownloaded: "PDF indirildi",
    photoUploaded: "Fotoğraf yüklendi",
    photoRemoved: "Fotoğraf kaldırıldı",
    newProfileCreatedNamed: "Yeni profil oluşturuldu: {name}",
    profileTranslatedEn: "Profil İngilizceye çevrildi",
    profileTranslatedTr: "Profil Türkçeye çevrildi",
    profileUpdatedToast: "Profil güncellendi",
    newCv: "Yeni CV",

    // Data backup
    backupSection: "Veri Yedekleme",
    backupHint: "Profillerinizi JSON olarak dışa aktarın veya yedekten geri yükleyin.",
    exportProfile: "JSON Dışa Aktar",
    exportAll: "Tümünü Yedekle",
    importProfile: "JSON İçe Aktar",
    importMergeConfirm: "Yedek dosyasındaki profiller mevcut profillere eklensin mi? (İsim çakışmalarında \" (import)\" eki eklenir)",
    importReplaceConfirm: "Tüm mevcut profiller silinip yedek dosyasındakilerle değiştirilecek. Bu işlem geri alınamaz. Devam etmek istiyor musunuz?",
    importSuccess: "{count} profil içe aktarıldı",
    importError: "İçe aktarma başarısız: {error}",
    exportProfileError: "Profil dışa aktarılamadı",
    exportAllError: "Yedek oluşturulamadı",
    invalidBackupFile: "Geçersiz yedek dosyası",
    resetSection: "Verileri Sıfırla",
    resetHint: "Tüm profilleri silip yalnızca örnek Full Stack Developer profilini bırakır.",
    resetButton: "Sıfırla",
    resetConfirm: "Tüm profiller silinecek, yalnızca örnek Full Stack Developer profili kalacak. Emin misiniz?",
    resetSuccess: "Veriler sıfırlandı, örnek profil yüklendi",
    resetError: "Sıfırlama başarısız",

    // Unsaved changes
    unsavedTitle: "Kaydedilmemiş değişiklikler",
    unsavedMessage: "Kaydedilmemiş değişiklikleriniz var. Çıkmadan önce kaydetmek ister misiniz?",
    unsavedSave: "Kaydet",
    unsavedDiscard: "Kaydetme",
    unsavedCancel: "İptal",
    autoSaved: "Otomatik kaydedildi",

    // First setup
    welcomeTitle: "AI Resume Builder'a hoş geldiniz",
    welcomeTagline: "ATS uyumlu CV, AI asistan, PDF export",
    welcomeSubtitle: "Başlamak için uygulama dilinizi ve temanızı seçin. CV içeriği dili ayrıca ayarlanabilir.",
    themeLabel: "Tema",
    themeLight: "Açık",
    themeDark: "Koyu",
    getStarted: "Başla",
  },
  en: {
    save: "Save",
    pdfDownload: "Download PDF",
    pdfPreparing: "Preparing PDF...",
    aiAssistant: "AI Assistant",
    darkMode: "Dark mode",
    lightMode: "Light mode",
    selectProfile: "Select profile",
    newProfile: "New profile",
    duplicateProfile: "Duplicate Profile",
    deleteProfile: "Delete profile",
    editingLanguage: "Editing Language",
    translateToEn: "Translate to English",
    translateToTr: "Translate to Turkish",
    translating: "Translating...",
    languageEmptyHint: "English content has not been created yet",
    livePreview: "Live Preview",
    cvPreview: "CV Preview",
    close: "Close",
    copy: "Copy",

    navProfile: "Profile",
    navSettings: "Settings",
    navAppSettings: "Application Settings",
    navDataBackup: "Data Backup",
    navGroupCv: "CV",
    navGroupApp: "App",
    settingsTitle: "Settings",
    aiAssistantSettings: "AI Assistant",
    aiAssistantEnabled: "Show AI Assistant",
    aiAssistantEnabledHint: "When off, the AI Assistant sidebar item and toolbar button are hidden.",
    aiAssistantDisabledToast: "AI Assistant is disabled in Settings. Enable it there to use it.",
    navPersonal: "Personal Information",
    navSummary: "Summary",
    navExperience: "Experience",
    navEducation: "Education",
    navProjects: "Projects",
    navSkills: "Skills",
    navPhoto: "Profile Photo",
    navAI: "AI Assistant",
    formSections: "Form sections",

    profileSettings: "Profile",
    profileSettingsHint: "Edit the name of the selected CV profile.",
    appSettings: "Application Settings",
    profileName: "Profile Name",
    profileNamePlaceholder: "e.g. Job Application",
    generalSettingsHint: "Editing language is managed from the top toolbar.",
    appLanguage: "App Language",
    appLanguageHint: "Choose the editor interface language. CV content language is managed separately.",
    targetTranslateLang: "Target Translation Language",
    targetTranslateLangHint: "Sets which language the toolbar translate button converts your CV into. The first added language is the primary target. When multiple languages are added, pick the active target from the toolbar. Remove languages with the × on each chip; when the list is empty, the translate button is hidden.",
    addLanguage: "Add",
    searchLanguages: "Search or select a language...",
    removeLanguage: "Remove",
    languageAlreadyAdded: "This language is already added",
    languageNotFound: "Language not found",
    translateToLang: "Translate: {lang}",
    cvTranslatedTo: "CV translated to {lang}",
    cvTranslatedStorageWarning: "Translation completed, but the CV only supports Turkish/English slots. Result was not saved.",
    langTurkish: "Türkçe",
    langEnglish: "English",

    personalInfo: "Personal Information",
    fullName: "Full Name",
    title: "Title",
    email: "Email",
    phone: "Phone",
    linkedin: "LinkedIn",
    github: "GitHub",
    portfolio: "Portfolio",
    location: "Location",
    birthDate: "Date of Birth",
    languages: "Languages",

    summary: "Summary",
    summaryPlaceholder: "Brief professional summary",
    experience: "Experience",
    education: "Education",
    projects: "Projects",
    skills: "Skills",
    profilePhoto: "Profile Photo",
    add: "+ Add",
    addGroup: "+ Add Group",
    remove: "Remove",
    delete: "Delete",

    photoPreviewAlt: "Profile photo preview",
    size: "Size",
    borderWidth: "Border width",
    borderColor: "Border color",

    company: "Company",
    startDate: "Start",
    endDate: "End",
    duration: "Duration",
    companyDescription: "Company Description",
    highlights: "Bullet points (one per line)",
    program: "Program",
    institution: "Institution",
    description: "Description",
    projectName: "Project Name",
    url: "URL",
    category: "Category",
    skillsPerLine: "Skills (one per line)",
    experienceN: "Experience {n}",
    educationN: "Education {n}",
    projectN: "Project {n}",
    skillGroupN: "Skill Group {n}",
    moveUp: "Move up",
    moveDown: "Move down",
    dragToReorder: "Drag to reorder",

    apiSetup: "API Setup",
    sessionCost: "Session Cost:",
    resetSessionCost: "Reset",
    resetSessionCostTitle: "Reset session cost",
    lastRequest: "Last request:",
    send: "Send",
    sending: "Sending...",
    apply: "Apply",
    applying: "Applying...",
    addPdf: "Add PDF",
    processing: "Processing...",
    aiChatPlaceholder: "e.g. build my CV for a Test engineer role",
    aiChatSection: "AI Chat",
    aiToolbar: "AI toolbar",
    apiKeys: "API Keys",
    apiKeyPlaceholder: "API key",
    aiSettings: "AI Settings",
    openaiModel: "OpenAI Model",
    fetchModels: "Fetch Models",
    loading: "Loading...",
    testConnection: "Test Connection",
    testing: "Testing...",
    translateProvider: "Translation Provider",
    googleTranslate: "Google Translate",
    openai: "OpenAI",
    apiKeysHint: "Keys can be configured on the server via .env. Values entered in the browser are saved to localStorage and take priority over server keys.",
    apiKeyLocal: "Saved in browser (priority)",
    apiKeyEnv: "Configured via .env",
    apiKeyEmpty: "Not configured",
    thinking: "Thinking",
    showText: "Show Text",
    hideText: "Hide Text",
    applyToCv: "Apply to CV",
    modelsLoading: "Loading models...",
    modelsFetching: "Fetching models...",
    modelsAlreadyLoading: "Models are already loading...",
    modelsLoaded: "{count} models loaded",
    defaultModelUsed: "Using default model",
    defaultSuffix: "(default)",
    modelSelectorsNotFound: "Model selectors not found",
    openaiNotConfigured: "OpenAI not configured (enter API key via .env or UI)",
    openaiKeyRequired: "OpenAI API key required (.env or UI)",
    connectionSuccess: "Connection successful",
    openaiConnectionSuccess: "OpenAI connection successful",
    onlyPdfAllowed: "Only PDF files can be uploaded",
    pdfTooLarge: "PDF file cannot exceed 10 MB",
    pdfProcessing: "Processing PDF...",
    pdfPagesOcr: "{pages} pages read via OCR",
    pdfPagesExtracted: "{pages} pages of text extracted",
    pdfProcessed: "PDF processed: {pages} pages, {chars} characters",
    textAddedToChat: "Text added to chat. Send with OpenAI key or copy the text.",
    noContentToApply: "No content to apply",
    aiChangesApplied: "AI changes applied",
    selectProfileFirst: "Select a profile first",
    profileDataCopied: "Profile data copied",
    whichProfile: "Which profile would you like to copy from?",
    noOtherProfiles: "No other profiles available to copy.",
    dataCopiedFromProfile: "Data copied from selected profile.",
    noResponse: "No response received",
    error: "Error",
    action: "Action",
    ocrMethod: "OCR (Tesseract)",
    pdfTextLayer: "PDF text layer",
    pdfUploaded: "[CV PDF uploaded: {name}]",
    pdfMeta: "Pages: {pages} | Characters: {chars} | Method: {method}",
    pdfTextHeader: "--- PDF TEXT ---",
    pdfExtractedSummary: "Extracted text from {pages} pages ({chars} characters, {method}).",
    fileLabel: "File: {name}",
    applyPdfPrompt: "Apply the PDF text above to my current CV profile. Fill in missing fields logically.",

    translateProgress1: "Translating your CV...",
    translateProgress2: "Translating summary...",
    translateProgress3: "Translating experience...",
    translateProgress4: "Translating education...",
    translateProgress5: "Translating skills...",
    translateProgress6: "Saving...",
    translateComplete: "Translation complete!",
    translateCompleteCost: "Translation complete! Cost: {cost}{tokens}",
    translateCompleteGoogle: "Translation complete! (Google Translate — free API quota)",
    translateFailed: "Translation failed: {error}",
    costLabel: "Cost: {cost}",
    tokenSuffix: " ({prompt} + {completion} tokens)",
    unknownError: "Unknown error",

    actionCreateProfile: "Creating new profile",
    actionDuplicateProfile: "Duplicating profile",
    actionUpdateContent: "Updating profile content",
    actionTranslateProfile: "Translating profile",
    actionSwitchLanguage: "Switching language",
    actionGetProfile: "Reading profile",
    actionApplyToProfile: "Applying content to profile",
    actionRequestSelection: "Showing profile selection",
    actionCopyFromProfile: "Copying data from profile",
    profileCreated: "Profile created: {name}",
    profileUpdated: "Profile updated.",

    profileSaved: "Profile saved",
    newProfileCreated: "New profile created",
    profileDuplicated: "Profile duplicated",
    profileDeleted: "Profile deleted",
    profilesLoadError: "Could not load profiles. Please refresh the page.",
    confirmDelete: "Are you sure you want to delete this profile?",
    enterTrContent: "Enter Turkish content first",
    enterEnContent: "Enter English content first",
    cvTranslatedEn: "CV translated to English",
    cvTranslatedTr: "CV translated to Turkish",
    pdfDownloaded: "PDF downloaded",
    photoUploaded: "Photo uploaded",
    photoRemoved: "Photo removed",
    newProfileCreatedNamed: "New profile created: {name}",
    profileTranslatedEn: "Profile translated to English",
    profileTranslatedTr: "Profile translated to Turkish",
    profileUpdatedToast: "Profile updated",
    newCv: "New CV",

    backupSection: "Data Backup",
    backupHint: "Export your profiles as JSON or restore from a backup file.",
    exportProfile: "Export JSON",
    exportAll: "Backup All",
    importProfile: "Import JSON",
    importMergeConfirm: "Add profiles from the backup file to your existing profiles? (Duplicate names get an \" (import)\" suffix)",
    importReplaceConfirm: "All existing profiles will be deleted and replaced with the backup. This cannot be undone. Continue?",
    importSuccess: "{count} profile(s) imported",
    importError: "Import failed: {error}",
    exportProfileError: "Could not export profile",
    exportAllError: "Could not create backup",
    invalidBackupFile: "Invalid backup file",
    resetSection: "Reset Data",
    resetHint: "Deletes all profiles and keeps only the example Full Stack Developer profile.",
    resetButton: "Reset",
    resetConfirm: "All profiles will be deleted and only the example Full Stack Developer profile will remain. Are you sure?",
    resetSuccess: "Data reset; example profile loaded",
    resetError: "Reset failed",

    unsavedTitle: "Unsaved changes",
    unsavedMessage: "You have unsaved changes. Save them before leaving?",
    unsavedSave: "Save",
    unsavedDiscard: "Don't save",
    unsavedCancel: "Cancel",
    autoSaved: "Auto-saved",

    welcomeTitle: "Welcome to AI Resume Builder",
    welcomeTagline: "ATS-friendly CV, AI assistant, PDF export",
    welcomeSubtitle: "Choose your app language and theme to get started. CV content language can be set separately.",
    themeLabel: "Theme",
    themeLight: "Light",
    themeDark: "Dark",
    getStarted: "Get Started",
  },
};

let currentLang = "tr";
const languageChangeListeners = [];
const targetTranslateLangListeners = [];

function getDefaultTargetTranslateLang(appLang) {
  const lang = appLang === "en" ? "en" : "tr";
  return lang === "tr" ? "en" : "tr";
}

function migrateTargetTranslateLangs() {
  if (localStorage.getItem(TARGET_TRANSLATE_LANGS_KEY)) return;

  const legacy = localStorage.getItem(TARGET_TRANSLATE_LANG_KEY);
  if (legacy === "tr" || legacy === "en") {
    localStorage.setItem(TARGET_TRANSLATE_LANGS_KEY, JSON.stringify([legacy]));
    return;
  }

  localStorage.setItem(
    TARGET_TRANSLATE_LANGS_KEY,
    JSON.stringify([getDefaultTargetTranslateLang(getAppLanguage())])
  );
}

function normalizeLangCode(code) {
  return String(code || "").trim().toLowerCase();
}

function getLanguageMeta(code) {
  const normalized = normalizeLangCode(code);
  return window.WorldLanguages?.byCode?.[normalized] || null;
}

function getLanguageDisplayName(code) {
  const meta = getLanguageMeta(code);
  if (!meta) return code;
  return currentLang === "en" ? meta.english : meta.native;
}

function getTargetTranslateLangs() {
  migrateTargetTranslateLangs();
  try {
    const parsed = JSON.parse(localStorage.getItem(TARGET_TRANSLATE_LANGS_KEY) || "[]");
    if (!Array.isArray(parsed)) return [];
    const unique = [];
    parsed.forEach((code) => {
      const normalized = normalizeLangCode(code);
      if (normalized.length === 2 && !unique.includes(normalized)) {
        unique.push(normalized);
      }
    });
    return unique;
  } catch {
    return [];
  }
}

function syncActiveTargetTranslateLang(langs) {
  if (!langs.length) {
    localStorage.removeItem(ACTIVE_TARGET_TRANSLATE_LANG_KEY);
    return;
  }

  const stored = normalizeLangCode(localStorage.getItem(ACTIVE_TARGET_TRANSLATE_LANG_KEY));
  if (!stored || !langs.includes(stored)) {
    localStorage.setItem(ACTIVE_TARGET_TRANSLATE_LANG_KEY, langs[0]);
  }
}

function notifyTargetTranslateLangChange() {
  const active = getActiveTargetTranslateLang();
  targetTranslateLangListeners.forEach((fn) => fn(active));
  window.CVEditor?.refreshTranslateUI?.();
}

function setTargetTranslateLangs(langs, manual = true) {
  const unique = [];
  (langs || []).forEach((code) => {
    const normalized = normalizeLangCode(code);
    if (normalized.length === 2 && !unique.includes(normalized)) {
      unique.push(normalized);
    }
  });
  localStorage.setItem(TARGET_TRANSLATE_LANGS_KEY, JSON.stringify(unique));
  if (unique.length > 0) {
    localStorage.setItem(TARGET_TRANSLATE_LANG_KEY, unique[0]);
  } else {
    localStorage.removeItem(TARGET_TRANSLATE_LANG_KEY);
  }
  if (manual) {
    localStorage.setItem(TARGET_TRANSLATE_LANG_MANUAL_KEY, "true");
  } else {
    localStorage.removeItem(TARGET_TRANSLATE_LANG_MANUAL_KEY);
  }
  syncActiveTargetTranslateLang(unique);
  renderTargetTranslateLangChips();
  notifyTargetTranslateLangChange();
}

function addTargetTranslateLang(code) {
  const normalized = normalizeLangCode(code);
  if (!getLanguageMeta(normalized)) return false;
  const langs = getTargetTranslateLangs();
  if (langs.includes(normalized)) return false;
  setTargetTranslateLangs([...langs, normalized], true);
  return true;
}

function removeTargetTranslateLang(code) {
  const normalized = normalizeLangCode(code);
  const langs = getTargetTranslateLangs().filter((item) => item !== normalized);
  setTargetTranslateLangs(langs, true);
}

function resolveLanguageFromInput(value) {
  const trimmed = String(value || "").trim();
  if (!trimmed) return null;

  const codeInParen = trimmed.match(/\(([a-z]{2})\)\s*$/i);
  if (codeInParen && getLanguageMeta(codeInParen[1])) {
    return normalizeLangCode(codeInParen[1]);
  }

  const byCode = normalizeLangCode(trimmed);
  if (getLanguageMeta(byCode)) return byCode;

  const lower = trimmed.toLowerCase();
  const list = window.WorldLanguages?.list || [];
  const exact = list.find(
    (lang) =>
      lang.native.toLowerCase() === lower ||
      lang.english.toLowerCase() === lower ||
      `${lang.native} (${lang.english})`.toLowerCase() === lower
  );
  if (exact) return exact.code;

  const partial = list.find(
    (lang) =>
      lang.native.toLowerCase().includes(lower) ||
      lang.english.toLowerCase().includes(lower) ||
      lang.code.includes(lower)
  );
  return partial?.code || null;
}

function getActiveTargetTranslateLang() {
  const langs = getTargetTranslateLangs();
  if (!langs.length) return null;

  const stored = normalizeLangCode(localStorage.getItem(ACTIVE_TARGET_TRANSLATE_LANG_KEY));
  if (stored && langs.includes(stored)) return stored;

  localStorage.setItem(ACTIVE_TARGET_TRANSLATE_LANG_KEY, langs[0]);
  return langs[0];
}

function setActiveTargetTranslateLang(code) {
  const normalized = normalizeLangCode(code);
  const langs = getTargetTranslateLangs();
  if (!langs.includes(normalized)) return;
  localStorage.setItem(ACTIVE_TARGET_TRANSLATE_LANG_KEY, normalized);
  notifyTargetTranslateLangChange();
}

function getTargetTranslateLang() {
  return getActiveTargetTranslateLang();
}

function setTargetTranslateLang(lang, manual = true) {
  const normalized = normalizeLangCode(lang);
  const langs = getTargetTranslateLangs().filter((item) => item !== normalized);
  setTargetTranslateLangs([normalized, ...langs], manual);
}

function isTargetLangStorageSupported(code) {
  const normalized = normalizeLangCode(code);
  return normalized.length === 2 && Boolean(getLanguageMeta(normalized));
}

function onTargetTranslateLangChange(fn) {
  targetTranslateLangListeners.push(fn);
}

function syncTargetTranslateLangPicker() {
  const datalist = document.getElementById("targetTranslateLangList");
  if (!datalist) return;

  const added = new Set(getTargetTranslateLangs());
  const options = (window.WorldLanguages?.list || [])
    .filter((lang) => !added.has(lang.code))
    .map((lang) => {
      const label = currentLang === "en" ? lang.english : lang.native;
      return `<option value="${label} (${lang.code})"></option>`;
    });
  datalist.innerHTML = options.join("");
}

function renderTargetTranslateLangChips() {
  const container = document.getElementById("targetTranslateLangChips");
  if (!container) return;

  const langs = getTargetTranslateLangs();
  container.innerHTML = langs
    .map((code) => {
      const label = getLanguageDisplayName(code);
      const removeLabel = t("removeLanguage");
      return `<span class="lang-chip" data-lang="${code}">
        <span class="lang-chip__label">${label}</span>
        <button type="button" class="lang-chip__remove" data-remove-lang="${code}" aria-label="${removeLabel}" title="${removeLabel}">×</button>
      </span>`;
    })
    .join("");

  container.querySelectorAll("[data-remove-lang]").forEach((btn) => {
    btn.addEventListener("click", () => {
      removeTargetTranslateLang(btn.dataset.removeLang);
    });
  });

  syncTargetTranslateLangPicker();
}

function syncTargetTranslateLangSelect() {
  renderTargetTranslateLangChips();
}

function detectBrowserLanguage() {
  const nav = (navigator.language || navigator.userLanguage || "tr").toLowerCase();
  return nav.startsWith("tr") ? "tr" : "en";
}

function getAppLanguage() {
  const stored = localStorage.getItem(APP_LANGUAGE_KEY);
  if (stored === "tr" || stored === "en") return stored;
  return detectBrowserLanguage();
}

function t(key, params = {}) {
  let text = translations[currentLang]?.[key] ?? translations.tr[key] ?? key;
  Object.entries(params).forEach(([k, v]) => {
    text = text.replace(new RegExp(`\\{${k}\\}`, "g"), String(v));
  });
  return text;
}

function applyTranslations() {
  document.documentElement.lang = currentLang;

  document.querySelectorAll("[data-i18n]").forEach((el) => {
    const key = el.dataset.i18n;
    if (key) el.textContent = t(key);
  });

  document.querySelectorAll("[data-i18n-placeholder]").forEach((el) => {
    el.placeholder = t(el.dataset.i18nPlaceholder);
  });

  document.querySelectorAll("[data-i18n-title]").forEach((el) => {
    el.title = t(el.dataset.i18nTitle);
  });

  document.querySelectorAll("[data-i18n-aria-label]").forEach((el) => {
    el.setAttribute("aria-label", t(el.dataset.i18nAriaLabel));
  });

  document.querySelectorAll("[data-i18n-alt]").forEach((el) => {
    el.alt = t(el.dataset.i18nAlt);
  });

  const appLangSelect = document.getElementById("appLanguage");
  if (appLangSelect && appLangSelect.value !== currentLang) {
    appLangSelect.value = currentLang;
  }

  syncTargetTranslateLangSelect();

  const setupLangSelect = document.getElementById("setupLanguage");
  if (setupLangSelect && setupLangSelect.value !== currentLang) {
    setupLangSelect.value = currentLang;
  }

  syncSetupThemeRadios();
}

function getStoredTheme() {
  const theme = localStorage.getItem("theme");
  return theme === "dark" || theme === "light" ? theme : "light";
}

function applySetupTheme(theme) {
  const resolved = theme === "dark" ? "dark" : "light";
  document.documentElement.setAttribute("data-theme", resolved);
}

function syncSetupThemeRadios() {
  const theme = getStoredTheme();
  document.querySelectorAll('input[name="setupTheme"]').forEach((input) => {
    input.checked = input.value === theme;
  });
}

function setAppLanguage(lang) {
  const normalized = lang === "en" ? "en" : "tr";
  currentLang = normalized;
  localStorage.setItem(APP_LANGUAGE_KEY, normalized);
  if (localStorage.getItem(TARGET_TRANSLATE_LANG_MANUAL_KEY) !== "true") {
    setTargetTranslateLangs([getDefaultTargetTranslateLang(normalized)], false);
  }
  applyTranslations();
  languageChangeListeners.forEach((fn) => fn(normalized));
}

function onLanguageChange(fn) {
  languageChangeListeners.push(fn);
}

function isSetupComplete() {
  return localStorage.getItem(APP_SETUP_COMPLETE_KEY) === "true";
}

function shouldShowSetup() {
  if (isSetupComplete()) return false;
  if (localStorage.getItem(APP_LANGUAGE_KEY)) return false;
  // Migration: existing users who already have preferences
  if (localStorage.getItem("theme") || localStorage.getItem("openaiApiKey")) {
    localStorage.setItem(APP_SETUP_COMPLETE_KEY, "true");
    localStorage.setItem(APP_LANGUAGE_KEY, getAppLanguage());
    return false;
  }
  return true;
}

function hideSetupModal() {
  const modal = document.getElementById("appSetupModal");
  if (!modal) return;
  modal.classList.add("hidden");
  modal.setAttribute("aria-hidden", "true");
  document.body.classList.remove("app-setup-modal-open");
}

function showSetupModal() {
  const modal = document.getElementById("appSetupModal");
  if (!modal) return;
  modal.classList.remove("hidden");
  modal.setAttribute("aria-hidden", "false");
  document.body.classList.add("app-setup-modal-open");
}

function completeSetup(lang, theme) {
  setAppLanguage(lang || getAppLanguage());
  const resolvedTheme = theme === "dark" || theme === "light" ? theme : getStoredTheme();
  localStorage.setItem("theme", resolvedTheme);
  applySetupTheme(resolvedTheme);
  localStorage.setItem(APP_SETUP_COMPLETE_KEY, "true");
  hideSetupModal();
}

function initSetupModal() {
  const modal = document.getElementById("appSetupModal");
  if (!modal) return;

  const langSelect = document.getElementById("setupLanguage");
  const startBtn = document.getElementById("setupStartBtn");

  syncSetupThemeRadios();

  langSelect?.addEventListener("change", () => {
    setAppLanguage(langSelect.value);
  });

  document.querySelectorAll('input[name="setupTheme"]').forEach((input) => {
    input.addEventListener("change", () => {
      if (input.checked) {
        applySetupTheme(input.value);
      }
    });
  });

  startBtn?.addEventListener("click", () => {
    const themeInput = document.querySelector('input[name="setupTheme"]:checked');
    completeSetup(langSelect?.value || currentLang, themeInput?.value || getStoredTheme());
  });

  if (shouldShowSetup()) {
    showSetupModal();
  }
}

function initTargetTranslateLangSetting() {
  const searchInput = document.getElementById("targetTranslateLangSearch");
  const addBtn = document.getElementById("targetTranslateLangAdd");
  if (!searchInput || !addBtn) return;

  renderTargetTranslateLangChips();

  const handleAdd = () => {
    const code = resolveLanguageFromInput(searchInput.value);
    if (!code) {
      window.showToast?.(t("languageNotFound"));
      return;
    }
    if (getTargetTranslateLangs().includes(code)) {
      window.showToast?.(t("languageAlreadyAdded"));
      return;
    }
    addTargetTranslateLang(code);
    searchInput.value = "";
    searchInput.focus();
  };

  addBtn.addEventListener("click", handleAdd);
  searchInput.addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
      event.preventDefault();
      handleAdd();
    }
  });
}

function initAppLanguageSetting() {
  const select = document.getElementById("appLanguage");
  if (!select) return;
  select.value = currentLang;
  select.addEventListener("change", () => {
    setAppLanguage(select.value);
  });
}

function initI18n() {
  currentLang = getAppLanguage();
  applyTranslations();
  initAppLanguageSetting();
  initTargetTranslateLangSetting();
  initSetupModal();
}

function getTranslateProgressStages() {
  return [
    t("translateProgress1"),
    t("translateProgress2"),
    t("translateProgress3"),
    t("translateProgress4"),
    t("translateProgress5"),
    t("translateProgress6"),
  ];
}

window.I18n = {
  t,
  getAppLanguage,
  setAppLanguage,
  getTargetTranslateLang,
  getActiveTargetTranslateLang,
  getTargetTranslateLangs,
  setTargetTranslateLang,
  setActiveTargetTranslateLang,
  setTargetTranslateLangs,
  addTargetTranslateLang,
  removeTargetTranslateLang,
  getLanguageDisplayName,
  isTargetLangStorageSupported,
  onTargetTranslateLangChange,
  applyTranslations,
  onLanguageChange,
  initI18n,
  isSetupComplete,
  completeSetup,
  shouldShowSetup,
  getTranslateProgressStages,
};
