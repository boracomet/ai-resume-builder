function t(key, params) {
  return window.I18n?.t(key, params) ?? key;
}

const API_KEY_CONFIG = [
  {
    storageKey: "googleTranslateApiKey",
    inputId: "googleApiKey",
    statusId: "googleApiKeyStatus",
    envField: "googleTranslateConfigured",
  },
  {
    storageKey: "openaiApiKey",
    inputId: "openaiApiKey",
    statusId: "openaiApiKeyStatus",
    envField: "openaiConfigured",
  },
  {
    storageKey: "googleGeminiApiKey",
    inputId: "googleGeminiApiKey",
    statusId: "googleGeminiApiKeyStatus",
    envField: "googleGeminiConfigured",
  },
  {
    storageKey: "xiaomiApiKey",
    inputId: "xiaomiApiKey",
    statusId: "xiaomiApiKeyStatus",
    envField: "xiaomiConfigured",
  },
];

const state = {
  profile: emptyProfile(),
  currentId: null,
  previewTimer: null,
  previewSeq: 0,
  previewAbort: null,
  photoPreviewDirty: false,
  envSettings: {},
  activeSection: "personal",
  dirty: false,
  saving: false,
  autoSaveTimer: null,
  unsavedResolver: null,
};

const AI_ASSISTANT_ENABLED_KEY = "aiAssistantEnabled";
const ACTIVE_SECTION_KEY = "activeEditorSection";

const SECTION_IDS = [
  "personal",
  "photo",
  "summary",
  "experience",
  "education",
  "projects",
  "skills",
  "settings",
  "ai-assistant",
];

const THEME_STORAGE_KEY = "theme";

const els = {
  profileSelect: document.getElementById("profileSelect"),
  newProfileBtn: document.getElementById("newProfileBtn"),
  duplicateProfileBtn: document.getElementById("duplicateProfileBtn"),
  deleteProfileBtn: document.getElementById("deleteProfileBtn"),
  saveBtn: document.getElementById("saveBtn"),
  pdfBtn: document.getElementById("pdfBtn"),
  themeToggle: document.getElementById("themeToggle"),
  profileName: document.getElementById("profileName"),
  photoInput: document.getElementById("photoInput"),
  photoPreview: document.getElementById("photoPreview"),
  photoPreviewImg: document.getElementById("photoPreviewImg"),
  removePhotoBtn: document.getElementById("removePhotoBtn"),
  photoSize: document.getElementById("photoSize"),
  photoSizeValue: document.getElementById("photoSizeValue"),
  photoBorderWidth: document.getElementById("photoBorderWidth"),
  photoBorderWidthValue: document.getElementById("photoBorderWidthValue"),
  photoBorderColor: document.getElementById("photoBorderColor"),
  previewFrame: document.getElementById("previewFrame"),
  cvLanguage: document.getElementById("cvLanguage"),
  languageEmptyHint: document.getElementById("languageEmptyHint"),
  translateControl: document.getElementById("translateControl"),
  translateTargetSelect: document.getElementById("translateTargetSelect"),
  translateBtn: document.getElementById("translateBtn"),
  translateToTrBtn: document.getElementById("translateToTrBtn"),
  experiencesList: document.getElementById("experiencesList"),
  educationList: document.getElementById("educationList"),
  projectsList: document.getElementById("projectsList"),
  skillGroupsList: document.getElementById("skillGroupsList"),
  toast: document.getElementById("toast"),
  exportProfileBtn: document.getElementById("exportProfileBtn"),
  exportAllBtn: document.getElementById("exportAllBtn"),
  importProfileInput: document.getElementById("importProfileInput"),
  resetProfilesBtn: document.getElementById("resetProfilesBtn"),
  editorSidebar: document.querySelector(".editor-sidebar"),
  formSections: document.querySelectorAll(".form-section[data-section-id]"),
  unsavedModal: document.getElementById("unsavedModal"),
  unsavedSaveBtn: document.getElementById("unsavedSaveBtn"),
  unsavedDiscardBtn: document.getElementById("unsavedDiscardBtn"),
  unsavedCancelBtn: document.getElementById("unsavedCancelBtn"),
};

function isAiAssistantEnabled() {
  const stored = localStorage.getItem(AI_ASSISTANT_ENABLED_KEY);
  if (stored === null) return true;
  return stored === "true";
}

function setAiAssistantEnabled(enabled) {
  localStorage.setItem(AI_ASSISTANT_ENABLED_KEY, enabled ? "true" : "false");
  applyAiAssistantVisibility();
  if (!enabled && state.activeSection === "ai-assistant") {
    showSection("settings");
  }
}

function applyAiAssistantVisibility() {
  const enabled = isAiAssistantEnabled();
  document.querySelectorAll("[data-ai-nav]").forEach((el) => {
    el.classList.toggle("hidden", !enabled);
  });
  const toggleBtn = document.getElementById("aiChatToggle");
  if (toggleBtn) {
    toggleBtn.classList.toggle("hidden", !enabled);
  }
  const checkbox = document.getElementById("aiAssistantEnabled");
  if (checkbox) {
    checkbox.checked = enabled;
  }
}

function showSection(sectionId) {
  if (!SECTION_IDS.includes(sectionId)) return;
  if (sectionId === "ai-assistant" && !isAiAssistantEnabled()) {
    sectionId = "settings";
  }

  state.activeSection = sectionId;
  try {
    localStorage.setItem(ACTIVE_SECTION_KEY, sectionId);
  } catch (_) {}

  els.formSections.forEach((section) => {
    section.classList.toggle("is-active", section.dataset.sectionId === sectionId);
  });

  els.editorSidebar?.querySelectorAll(".editor-nav-item").forEach((item) => {
    const isActive = item.dataset.section === sectionId;
    item.classList.toggle("is-active", isActive);
    item.setAttribute("aria-current", isActive ? "page" : "false");
  });

  const content = document.querySelector(".editor-content");
  if (content) {
    content.classList.toggle("editor-content--ai-chat", sectionId === "ai-assistant");
    content.scrollTop = 0;
  }
}

function getSavedSection() {
  try {
    const saved = localStorage.getItem(ACTIVE_SECTION_KEY);
    if (saved && SECTION_IDS.includes(saved)) {
      if (saved === "ai-assistant" && !isAiAssistantEnabled()) {
        return "settings";
      }
      return saved;
    }
  } catch (_) {}
  return "personal";
}

function initSectionNav() {
  els.editorSidebar?.querySelectorAll(".editor-nav-item").forEach((item) => {
    item.addEventListener("click", () => {
      showSection(item.dataset.section);
    });
  });
  applyAiAssistantVisibility();
  showSection(getSavedSection());
}

function initAiAssistantSetting() {
  const checkbox = document.getElementById("aiAssistantEnabled");
  if (!checkbox) return;
  checkbox.checked = isAiAssistantEnabled();
  checkbox.addEventListener("change", () => {
    setAiAssistantEnabled(checkbox.checked);
  });
  applyAiAssistantVisibility();
}

function emptyLocalizedContent() {
  return {
    summary: "",
    personalTitle: "",
    personalLanguages: "",
    experiences: [],
    education: [],
    projects: [],
    skillGroups: [],
  };
}

function emptyProfile() {
  return {
    id: 0,
    name: t("newCv"),
    language: "tr",
    personal: {
      name: "",
      email: "",
      phone: "",
      linkedin: "",
      github: "",
      portfolio: "",
      location: "",
      birthDate: "",
    },
    contentTR: emptyLocalizedContent(),
    contentEN: emptyLocalizedContent(),
    contentLanguages: {},
    photoBase64: "",
    photoSize: 96,
    photoBorderWidth: 2,
    photoBorderColor: "#2563eb",
  };
}

const PHOTO_DEFAULTS = {
  size: 96,
  borderWidth: 2,
  borderColor: "#2563eb",
};

function normalizeProfileLang(lang) {
  const normalized = String(lang || "").trim().toLowerCase();
  if (normalized === "en") return "en";
  if (normalized.length === 2 && /^[a-z]{2}$/.test(normalized)) {
    return normalized;
  }
  return "tr";
}

function normalizeLanguage(profile) {
  profile.language = normalizeProfileLang(profile.language);
}

function getEditingLanguageOptions() {
  const base = ["tr", "en"];
  const extras = (window.I18n?.getTargetTranslateLangs?.() ?? []).filter((code) => !base.includes(code));
  const options = [...base, ...extras];
  const current = normalizeProfileLang(state.profile?.language);
  if (!options.includes(current) && current !== "tr" && current !== "en") {
    options.push(current);
  }
  return options;
}

function getEditingLanguageLabel(code) {
  if (code === "tr") return t("langTurkish");
  if (code === "en") return t("langEnglish");
  return window.I18n?.getLanguageDisplayName?.(code) || code;
}

function renderEditingLanguageOptions() {
  if (!els.cvLanguage) return;

  const options = getEditingLanguageOptions();
  const current = normalizeProfileLang(state.profile?.language || els.cvLanguage.value);
  const selected = options.includes(current) ? current : "tr";

  els.cvLanguage.innerHTML = options
    .map((code) => `<option value="${code}">${getEditingLanguageLabel(code)}</option>`)
    .join("");
  els.cvLanguage.value = selected;

  if (state.profile && !options.includes(normalizeProfileLang(state.profile.language))) {
    state.profile.language = selected;
  }
}

function getContentForLang(profile, lang) {
  const normalized = normalizeProfileLang(lang);
  if (normalized === "en") {
    if (!profile.contentEN) profile.contentEN = emptyLocalizedContent();
    return profile.contentEN;
  }
  if (normalized === "tr") {
    if (!profile.contentTR) profile.contentTR = emptyLocalizedContent();
    return profile.contentTR;
  }
  if (!profile.contentLanguages) profile.contentLanguages = {};
  if (!profile.contentLanguages[normalized]) {
    profile.contentLanguages[normalized] = emptyLocalizedContent();
  }
  return profile.contentLanguages[normalized];
}

function setContentForLang(profile, lang, content) {
  const normalized = normalizeProfileLang(lang);
  if (normalized === "en") {
    profile.contentEN = content;
    return;
  }
  if (normalized === "tr") {
    profile.contentTR = content;
    return;
  }
  if (!profile.contentLanguages) profile.contentLanguages = {};
  profile.contentLanguages[normalized] = content;
}

function getActiveContent(profile = state.profile) {
  return getContentForLang(profile, profile.language);
}

function isContentEmpty(content) {
  if (!content) return true;
  if ((content.summary || "").trim()) return false;
  if ((content.personalTitle || "").trim()) return false;
  if ((content.personalLanguages || "").trim()) return false;
  return (
    (content.experiences || []).length === 0 &&
    (content.education || []).length === 0 &&
    (content.projects || []).length === 0 &&
    (content.skillGroups || []).length === 0
  );
}

function ensureBilingualStructure(profile) {
  if (!profile.contentTR) {
    profile.contentTR = emptyLocalizedContent();
  }
  if (!profile.contentEN) {
    profile.contentEN = emptyLocalizedContent();
  }

  const hasLegacy =
    (profile.summary && profile.summary.trim()) ||
    (profile.experiences && profile.experiences.length > 0) ||
    (profile.education && profile.education.length > 0) ||
    (profile.projects && profile.projects.length > 0) ||
    (profile.skillGroups && profile.skillGroups.length > 0) ||
    (profile.personal?.title && profile.personal.title.trim()) ||
    (profile.personal?.languages && profile.personal.languages.trim());

  if (hasLegacy && isContentEmpty(profile.contentTR) && isContentEmpty(profile.contentEN)) {
    profile.contentTR = {
      summary: profile.summary || "",
      personalTitle: profile.personal?.title || "",
      personalLanguages: profile.personal?.languages || "",
      experiences: profile.experiences || [],
      education: profile.education || [],
      projects: profile.projects || [],
      skillGroups: profile.skillGroups || [],
    };
    delete profile.summary;
    delete profile.experiences;
    delete profile.education;
    delete profile.projects;
    delete profile.skillGroups;
    if (profile.personal) {
      delete profile.personal.title;
      delete profile.personal.languages;
    }
  }

  ["contentTR", "contentEN"].forEach((key) => {
    const content = profile[key];
    content.experiences = content.experiences || [];
    content.education = content.education || [];
    content.projects = content.projects || [];
    content.skillGroups = content.skillGroups || [];
    content.education.forEach((edu) => {
      if (!edu.institution && edu.school) {
        edu.institution = edu.school;
        delete edu.school;
      }
    });
  });

  if (profile.contentLanguages) {
    Object.keys(profile.contentLanguages).forEach((lang) => {
      const content = getContentForLang(profile, lang);
      content.experiences = content.experiences || [];
      content.education = content.education || [];
      content.projects = content.projects || [];
      content.skillGroups = content.skillGroups || [];
      content.education.forEach((edu) => {
        if (!edu.institution && edu.school) {
          edu.institution = edu.school;
          delete edu.school;
        }
      });
    });
  }
}

function normalizeProfile(profile) {
  ensureBilingualStructure(profile);
  normalizePhotoSettings(profile);
  normalizeLanguage(profile);
}

function getApiKeyInput(id) {
  return document.getElementById(id);
}

function getApiKeyStatusEl(id) {
  return document.getElementById(id);
}

function loadApiKeysFromStorage() {
  API_KEY_CONFIG.forEach((config) => {
    const input = getApiKeyInput(config.inputId);
    if (!input) return;
    const saved = localStorage.getItem(config.storageKey);
    if (saved) {
      input.value = saved;
    }
  });
}

function saveApiKey(config) {
  const input = getApiKeyInput(config.inputId);
  if (!input) return;
  const value = input.value.trim();
  if (value) {
    localStorage.setItem(config.storageKey, value);
  } else {
    localStorage.removeItem(config.storageKey);
  }
  updateApiKeyStatus(config);
}

function getApiKeyValue(config) {
  const input = getApiKeyInput(config.inputId);
  const fromInput = input?.value.trim() || "";
  if (fromInput) return fromInput;
  const fromStorage = (localStorage.getItem(config.storageKey) || "").trim();
  return fromStorage;
}

function isOpenAIConfigured() {
  return !!getOpenAIApiKey() || !!state.envSettings.openaiConfigured;
}

function getOpenAIApiKey() {
  const config = API_KEY_CONFIG[1];
  return getApiKeyValue(config);
}

function getTranslateProvider() {
  const select = document.getElementById("translateProvider");
  const fromSelect = select?.value;
  if (fromSelect) return fromSelect;
  return localStorage.getItem("translateProvider") || "google";
}

function getGoogleApiKey() {
  const config = API_KEY_CONFIG[0];
  return getApiKeyValue(config);
}

function updateApiKeyStatus(config) {
  const statusEl = getApiKeyStatusEl(config.statusId);
  const input = getApiKeyInput(config.inputId);
  if (!statusEl || !input) return;

  const hasLocal = Boolean(input.value.trim() || localStorage.getItem(config.storageKey));
  const hasEnv = Boolean(state.envSettings[config.envField]);

  statusEl.className = "api-key-status";
  if (hasLocal) {
    statusEl.textContent = t("apiKeyLocal");
    statusEl.classList.add("api-key-status--local");
  } else if (hasEnv) {
    statusEl.textContent = t("apiKeyEnv");
    statusEl.classList.add("api-key-status--env");
  } else {
    statusEl.textContent = t("apiKeyEmpty");
    statusEl.classList.add("api-key-status--empty");
  }
}

function updateAllApiKeyStatuses() {
  API_KEY_CONFIG.forEach(updateApiKeyStatus);
}

async function loadEnvSettings() {
  try {
    state.envSettings = await API.getSettings();
  } catch (error) {
    console.warn("Ayarlar yüklenemedi:", error);
    state.envSettings = {};
  }
  updateAllApiKeyStatuses();
}

function normalizePhotoSettings(profile) {
  const legacy =
    !profile.photoSize &&
    !profile.photoBorderWidth &&
    !profile.photoBorderColor;

  if (!profile.photoSize || profile.photoSize <= 0) {
    profile.photoSize = PHOTO_DEFAULTS.size;
  }
  if (!profile.photoBorderColor) {
    profile.photoBorderColor = PHOTO_DEFAULTS.borderColor;
  }
  if (legacy || profile.photoBorderWidth < 0) {
    profile.photoBorderWidth = PHOTO_DEFAULTS.borderWidth;
  }
}

function applyPhotoPreviewStyles() {
  const borderWidth = Number(els.photoBorderWidth?.value) || 0;
  const borderColor = els.photoBorderColor?.value || PHOTO_DEFAULTS.borderColor;

  els.photoPreviewImg.style.border = `${borderWidth}px solid ${borderColor}`;
}

function syncPhotoSettingLabels() {
  if (els.photoSizeValue) {
    els.photoSizeValue.textContent = `${els.photoSize.value} px`;
  }
  if (els.photoBorderWidthValue) {
    els.photoBorderWidthValue.textContent = `${els.photoBorderWidth.value} px`;
  }
}

function getStoredTheme() {
  const theme = localStorage.getItem(THEME_STORAGE_KEY);
  return theme === "dark" || theme === "light" ? theme : "light";
}

function applyTheme(theme) {
  const resolved = theme === "dark" ? "dark" : "light";
  document.documentElement.setAttribute("data-theme", resolved);
  if (els.themeToggle) {
    const isDark = resolved === "dark";
    const label = isDark ? t("lightMode") : t("darkMode");
    els.themeToggle.setAttribute("aria-label", label);
    els.themeToggle.title = label;
  }
}

function initTheme() {
  applyTheme(getStoredTheme());
  els.themeToggle?.addEventListener("click", () => {
    const current = document.documentElement.getAttribute("data-theme") === "dark" ? "dark" : "light";
    const next = current === "dark" ? "light" : "dark";
    localStorage.setItem(THEME_STORAGE_KEY, next);
    applyTheme(next);
  });
}

function showToast(message) {
  els.toast.textContent = message;
  els.toast.classList.remove("hidden");
  clearTimeout(showToast.timer);
  showToast.timer = setTimeout(() => els.toast.classList.add("hidden"), 2500);
}

function setNested(obj, path, value) {
  const parts = path.split(".");
  let current = obj;
  for (let i = 0; i < parts.length - 1; i++) {
    current = current[parts[i]];
  }
  current[parts[parts.length - 1]] = value;
}

function getNested(obj, path) {
  return path.split(".").reduce((acc, key) => acc?.[key], obj);
}

function renderProfileSelect(profiles) {
  els.profileSelect.innerHTML = "";
  profiles.forEach((profile) => {
    const option = document.createElement("option");
    option.value = profile.id;
    option.textContent = profile.name;
    els.profileSelect.appendChild(option);
  });
  if (state.currentId) {
    els.profileSelect.value = String(state.currentId);
  }
}

function renderPhoto() {
  if (state.profile.photoBase64) {
    els.photoPreviewImg.src = state.profile.photoBase64;
    els.photoPreview.classList.remove("hidden");
    applyPhotoPreviewStyles();
  } else {
    els.photoPreview.classList.add("hidden");
    els.photoPreviewImg.removeAttribute("src");
  }
}

function renderPhotoSettings() {
  normalizePhotoSettings(state.profile);
  els.photoSize.value = String(state.profile.photoSize);
  els.photoBorderWidth.value = String(state.profile.photoBorderWidth);
  els.photoBorderColor.value = state.profile.photoBorderColor;
  syncPhotoSettingLabels();
  applyPhotoPreviewStyles();
}

function refreshTranslateUI() {
  const lang = normalizeProfileLang(state.profile?.language);
  const targets = window.I18n?.getTargetTranslateLangs?.() ?? [];
  const targetLang = window.I18n?.getTargetTranslateLang?.() ?? null;
  const showSelect = targets.length > 1;
  const showTranslate = Boolean(targetLang) && lang !== targetLang;
  const showControl = targets.length > 0 && (showSelect || showTranslate);
  const targetName = window.I18n?.getLanguageDisplayName?.(targetLang) || targetLang;
  const label = t("translateToLang", { lang: targetName });

  if (els.translateTargetSelect) {
    els.translateTargetSelect.classList.toggle("hidden", !showSelect);
    if (showSelect) {
      els.translateTargetSelect.setAttribute("aria-label", t("targetTranslateLang"));
      els.translateTargetSelect.innerHTML = targets
        .map((code) => {
          const name = window.I18n?.getLanguageDisplayName?.(code) || code;
          return `<option value="${code}">${name}</option>`;
        })
        .join("");
      if (targetLang) {
        els.translateTargetSelect.value = targetLang;
      }
    }
  }

  if (els.translateControl) {
    els.translateControl.classList.toggle("hidden", !showControl);
  }
  if (els.translateBtn) {
    els.translateBtn.classList.toggle("hidden", !showTranslate);
    if (showTranslate && !els.translateBtn.disabled) {
      els.translateBtn.textContent = label;
    }
  }
  if (els.translateToTrBtn) {
    els.translateToTrBtn.classList.add("hidden");
  }
}

function updateLanguageUI() {
  renderEditingLanguageOptions();
  const lang = normalizeProfileLang(state.profile.language);
  if (els.cvLanguage) {
    els.cvLanguage.value = lang;
  }

  const activeEmpty = isContentEmpty(getContentForLang(state.profile, lang));
  if (els.languageEmptyHint) {
    els.languageEmptyHint.classList.toggle("hidden", !(lang !== "tr" && activeEmpty));
  }

  refreshTranslateUI();
}

function renderStaticFields() {
  const content = getActiveContent();
  els.profileName.value = state.profile.name || "";
  document.querySelectorAll("[data-field]").forEach((input) => {
    const value = getNested(state.profile, input.dataset.field);
    input.value = value ?? "";
  });
  document.querySelectorAll("[data-localized]").forEach((input) => {
    const value = content[input.dataset.localized];
    input.value = value ?? "";
  });
  updateLanguageUI();
  renderPhotoSettings();
  renderPhoto();
}

function moveItem(array, index, direction) {
  const newIndex = index + direction;
  if (newIndex < 0 || newIndex >= array.length) return false;
  const [item] = array.splice(index, 1);
  array.splice(newIndex, 0, item);
  return true;
}

function reorderItem(array, fromIndex, toIndex) {
  if (
    fromIndex === toIndex ||
    fromIndex < 0 ||
    toIndex < 0 ||
    fromIndex >= array.length ||
    toIndex >= array.length
  ) {
    return false;
  }
  const [item] = array.splice(fromIndex, 1);
  array.splice(toIndex, 0, item);
  return true;
}

function reorderExperiences(fromIndex, toIndex) {
  const content = getActiveContent();
  if (reorderItem(content.experiences, fromIndex, toIndex)) {
    renderExperiences();
    onEditorChange();
  }
}

function reorderEducation(fromIndex, toIndex) {
  const content = getActiveContent();
  if (reorderItem(content.education, fromIndex, toIndex)) {
    renderEducation();
    onEditorChange();
  }
}

function reorderProjects(fromIndex, toIndex) {
  const content = getActiveContent();
  if (reorderItem(content.projects, fromIndex, toIndex)) {
    renderProjects();
    onEditorChange();
  }
}

function reorderSkillGroups(fromIndex, toIndex) {
  const content = getActiveContent();
  if (reorderItem(content.skillGroups, fromIndex, toIndex)) {
    renderSkillGroups();
    onEditorChange();
  }
}

function createReorderButton(direction, disabled, onClick) {
  const label = direction === "up" ? t("moveUp") : t("moveDown");
  const btn = document.createElement("button");
  btn.type = "button";
  btn.className = "btn btn-secondary btn-sm btn-icon-reorder";
  btn.setAttribute("aria-label", label);
  btn.title = label;
  btn.disabled = disabled;
  btn.innerHTML = direction === "up"
    ? '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" aria-hidden="true"><path d="M12 5l-7 7h14l-7-7z"/></svg>'
    : '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" aria-hidden="true"><path d="M12 19l7-7H5l7 7z"/></svg>';
  btn.addEventListener("click", onClick);
  return btn;
}

function attachCardDragReorder(card, index, onReorder) {
  let dragHandleActive = false;

  const dragHandle = document.createElement("button");
  dragHandle.type = "button";
  dragHandle.className = "card-drag-handle";
  dragHandle.setAttribute("aria-label", t("dragToReorder"));
  dragHandle.title = t("dragToReorder");
  dragHandle.innerHTML = '<svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><circle cx="9" cy="6" r="1.5"/><circle cx="15" cy="6" r="1.5"/><circle cx="9" cy="12" r="1.5"/><circle cx="15" cy="12" r="1.5"/><circle cx="9" cy="18" r="1.5"/><circle cx="15" cy="18" r="1.5"/></svg>';

  dragHandle.addEventListener("mousedown", () => {
    dragHandleActive = true;
  });
  dragHandle.addEventListener("mouseup", () => {
    dragHandleActive = false;
  });
  dragHandle.addEventListener("mouseleave", () => {
    dragHandleActive = false;
  });

  card.draggable = true;
  card.dataset.index = String(index);

  card.addEventListener("dragstart", (event) => {
    if (!dragHandleActive) {
      event.preventDefault();
      return;
    }
    event.dataTransfer.effectAllowed = "move";
    event.dataTransfer.setData("text/plain", String(index));
    card.classList.add("card--dragging");
  });

  card.addEventListener("dragend", () => {
    dragHandleActive = false;
    card.classList.remove("card--dragging");
    card.closest(".dynamic-list")?.querySelectorAll(".card--drag-over").forEach((el) => {
      el.classList.remove("card--drag-over");
    });
  });

  card.addEventListener("dragover", (event) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = "move";
  });

  card.addEventListener("dragenter", (event) => {
    event.preventDefault();
    if (!card.classList.contains("card--dragging")) {
      card.classList.add("card--drag-over");
    }
  });

  card.addEventListener("dragleave", (event) => {
    if (!card.contains(event.relatedTarget)) {
      card.classList.remove("card--drag-over");
    }
  });

  card.addEventListener("drop", (event) => {
    event.preventDefault();
    card.classList.remove("card--drag-over");
    const fromIndex = Number(event.dataTransfer.getData("text/plain"));
    if (!Number.isNaN(fromIndex) && fromIndex !== index) {
      onReorder(fromIndex, index);
    }
  });

  return dragHandle;
}

function createCard(title, onRemove, contentBuilder, reorder = null) {
  const card = document.createElement("div");
  card.className = "card";

  const header = document.createElement("div");
  header.className = "card-header";

  const headerLeft = document.createElement("div");
  headerLeft.className = "card-header__left";

  if (reorder) {
    headerLeft.appendChild(attachCardDragReorder(card, reorder.index, reorder.onReorder));
  }

  const heading = document.createElement("h3");
  heading.textContent = title;
  headerLeft.appendChild(heading);

  const actions = document.createElement("div");
  actions.className = "card-header__actions";

  if (reorder) {
    actions.append(
      createReorderButton("up", reorder.index === 0, () => reorder.onMove(-1)),
      createReorderButton("down", reorder.index === reorder.total - 1, () => reorder.onMove(1))
    );
  }

  const removeBtn = document.createElement("button");
  removeBtn.type = "button";
  removeBtn.className = "btn btn-danger btn-sm";
  removeBtn.textContent = t("delete");
  removeBtn.addEventListener("click", onRemove);
  actions.appendChild(removeBtn);

  header.append(headerLeft, actions);
  card.appendChild(header);
  contentBuilder(card);
  return card;
}

function bindInput(labelText, value, onChange, options = {}) {
  const label = document.createElement("label");
  label.textContent = labelText;
  const input = document.createElement(options.type === "textarea" ? "textarea" : "input");
  if (options.type === "textarea") {
    input.rows = options.rows || 3;
  } else if (options.type) {
    input.type = options.type;
  }
  input.value = value || "";
  input.addEventListener("input", (e) => {
    onChange(e.target.value);
    onEditorChange();
  });
  label.appendChild(input);
  return label;
}

function renderExperiences() {
  const content = getActiveContent();
  els.experiencesList.innerHTML = "";
  const total = content.experiences.length;
  content.experiences.forEach((exp, index) => {
    const card = createCard(t("experienceN", { n: index + 1 }), () => {
      content.experiences.splice(index, 1);
      renderExperiences();
      onEditorChange();
    }, (container) => {
      const fields = document.createElement("div");
      fields.className = "card-grid";
      fields.append(
        bindInput(t("title"), exp.title, (v) => (exp.title = v)),
        bindInput(t("company"), exp.company, (v) => (exp.company = v)),
        bindInput(t("startDate"), exp.startDate, (v) => (exp.startDate = v)),
        bindInput(t("endDate"), exp.endDate, (v) => (exp.endDate = v)),
        bindInput(t("duration"), exp.duration, (v) => (exp.duration = v)),
        bindInput(t("location"), exp.location, (v) => (exp.location = v)),
        bindInput(t("companyDescription"), exp.description, (v) => (exp.description = v), { type: "textarea", rows: 2 }),
        bindInput(t("highlights"), (exp.highlights || []).join("\n"), (v) => {
          // Keep raw lines while typing so preview stays in sync; trim empties on save via normalize.
          exp.highlights = v.split("\n");
        }, { type: "textarea", rows: 4 })
      );
      container.appendChild(fields);
    }, {
      index,
      total,
      onMove: (direction) => {
        if (moveItem(content.experiences, index, direction)) {
          renderExperiences();
          schedulePreview();
        }
      },
      onReorder: reorderExperiences,
    });
    els.experiencesList.appendChild(card);
  });
}

function renderEducation() {
  const content = getActiveContent();
  els.educationList.innerHTML = "";
  const total = content.education.length;
  content.education.forEach((edu, index) => {
    const card = createCard(t("educationN", { n: index + 1 }), () => {
      content.education.splice(index, 1);
      renderEducation();
      onEditorChange();
    }, (container) => {
      const fields = document.createElement("div");
      fields.className = "card-grid";
      fields.append(
        bindInput(t("program"), edu.degree, (v) => (edu.degree = v)),
        bindInput(t("institution"), edu.institution, (v) => (edu.institution = v)),
        bindInput(t("startDate"), edu.startDate, (v) => (edu.startDate = v)),
        bindInput(t("endDate"), edu.endDate, (v) => (edu.endDate = v)),
        bindInput(t("description"), edu.description, (v) => (edu.description = v), { type: "textarea", rows: 3 })
      );
      container.appendChild(fields);
    }, {
      index,
      total,
      onMove: (direction) => {
        if (moveItem(content.education, index, direction)) {
          renderEducation();
          schedulePreview();
        }
      },
      onReorder: reorderEducation,
    });
    els.educationList.appendChild(card);
  });
}

function renderProjects() {
  const content = getActiveContent();
  els.projectsList.innerHTML = "";
  const total = content.projects.length;
  content.projects.forEach((project, index) => {
    const card = createCard(t("projectN", { n: index + 1 }), () => {
      content.projects.splice(index, 1);
      renderProjects();
      onEditorChange();
    }, (container) => {
      const fields = document.createElement("div");
      fields.className = "card-grid";
      fields.append(
        bindInput(t("projectName"), project.name, (v) => (project.name = v)),
        bindInput(t("url"), project.url, (v) => (project.url = v)),
        bindInput(t("description"), project.description, (v) => (project.description = v), { type: "textarea", rows: 3 })
      );
      container.appendChild(fields);
    }, {
      index,
      total,
      onMove: (direction) => {
        if (moveItem(content.projects, index, direction)) {
          renderProjects();
          schedulePreview();
        }
      },
      onReorder: reorderProjects,
    });
    els.projectsList.appendChild(card);
  });
}

function renderSkillGroups() {
  const content = getActiveContent();
  els.skillGroupsList.innerHTML = "";
  const total = content.skillGroups.length;
  content.skillGroups.forEach((group, index) => {
    const card = createCard(t("skillGroupN", { n: index + 1 }), () => {
      content.skillGroups.splice(index, 1);
      renderSkillGroups();
      onEditorChange();
    }, (container) => {
      const fields = document.createElement("div");
      fields.className = "card-grid";
      fields.append(
        bindInput(t("category"), group.category, (v) => (group.category = v)),
        bindInput(t("skillsPerLine"), (group.skills || []).join("\n"), (v) => {
          group.skills = v.split("\n").map((line) => line.trim()).filter(Boolean);
        }, { type: "textarea", rows: 5 })
      );
      container.appendChild(fields);
    }, {
      index,
      total,
      onMove: (direction) => {
        if (moveItem(content.skillGroups, index, direction)) {
          renderSkillGroups();
          schedulePreview();
        }
      },
      onReorder: reorderSkillGroups,
    });
    els.skillGroupsList.appendChild(card);
  });
}

function renderAllLists() {
  renderExperiences();
  renderEducation();
  renderProjects();
  renderSkillGroups();
}

function renderForm() {
  renderStaticFields();
  renderAllLists();
}

// Dynamic card inputs close over the content objects they were built from, so a
// profile swap must rebuild the form; otherwise those fields keep writing to an
// orphaned object graph and edits silently disappear until reload.
function setProfile(profile, options = {}) {
  state.profile = profile;
  normalizeProfile(state.profile);
  state.photoPreviewDirty = false;
  if (options.render !== false) {
    renderForm();
  }
}

// Applied instead of swapping state.profile after a write, so in-flight edits and
// the live input bindings survive.
function applyServerMeta(saved) {
  if (!saved) return;
  if (saved.id) state.profile.id = saved.id;
  if (saved.createdAt) state.profile.createdAt = saved.createdAt;
  if (saved.updatedAt) state.profile.updatedAt = saved.updatedAt;
}

function collectProfileFromForm(options = {}) {
  const editingLang = normalizeProfileLang(options.language ?? els.cvLanguage?.value ?? state.profile.language);
  const content = getContentForLang(state.profile, editingLang);
  setContentForLang(state.profile, editingLang, content);

  state.profile.name = els.profileName.value.trim() || t("newCv");
  if (!options.skipLanguageUpdate && els.cvLanguage) {
    state.profile.language = normalizeProfileLang(els.cvLanguage.value);
  }

  document.querySelectorAll("[data-field]").forEach((input) => {
    let value = input.value;
    if (input.type === "range") {
      value = Number(value);
    }
    setNested(state.profile, input.dataset.field, value);
  });

  document.querySelectorAll("[data-localized]").forEach((input) => {
    content[input.dataset.localized] = input.value;
  });

  normalizeProfile(state.profile);
  return state.profile;
}

function switchEditingLanguage(nextLang) {
  const currentLang = normalizeProfileLang(state.profile.language);
  const normalizedNext = normalizeProfileLang(nextLang);
  if (currentLang === normalizedNext) {
    return;
  }

  collectProfileFromForm({ language: currentLang, skipLanguageUpdate: true });
  state.profile.language = normalizedNext;
  if (els.cvLanguage) {
    els.cvLanguage.value = normalizedNext;
  }
  renderForm();
  markDirty();
  schedulePreview();
}

async function updatePreview() {
  collectProfileFromForm();
  const seq = ++state.previewSeq;

  if (state.previewAbort) {
    try {
      state.previewAbort.abort();
    } catch (_) {}
  }
  state.previewAbort = new AbortController();
  const { signal } = state.previewAbort;

  // Don't resend multi-MB photoBase64 on every keystroke — server reuses stored photo.
  const payload = {
    id: state.profile.id || state.currentId || 0,
    name: state.profile.name,
    language: state.profile.language,
    personal: state.profile.personal,
    contentTR: state.profile.contentTR,
    contentEN: state.profile.contentEN,
    contentLanguages: state.profile.contentLanguages,
    photoSize: state.profile.photoSize,
    photoBorderWidth: state.profile.photoBorderWidth,
    photoBorderColor: state.profile.photoBorderColor,
    photoBase64: state.photoPreviewDirty ? (state.profile.photoBase64 || "") : "",
    useStoredPhoto: !state.photoPreviewDirty,
  };

  try {
    const html = await API.preview(payload, { signal });
    if (seq !== state.previewSeq || signal.aborted) return;
    els.previewFrame.srcdoc = html;
  } catch (error) {
    if (error?.name === "AbortError" || signal.aborted) return;
    console.error(error);
  }
}

function schedulePreview() {
  clearTimeout(state.previewTimer);
  state.previewTimer = setTimeout(() => {
    updatePreview();
  }, 400);
}

function markDirty() {
  if (!state.currentId) return;
  state.dirty = true;
  updateSaveButtonHint();
}

function markClean() {
  state.dirty = false;
  updateSaveButtonHint();
}

function isDirty() {
  return Boolean(state.dirty && state.currentId);
}

function updateSaveButtonHint() {
  if (!els.saveBtn) return;
  if (state.dirty) {
    els.saveBtn.classList.add("is-dirty");
    els.saveBtn.title = t("unsavedTitle");
  } else {
    els.saveBtn.classList.remove("is-dirty");
    els.saveBtn.title = t("save");
  }
}

function onEditorChange() {
  markDirty();
  schedulePreview();
}

function hideUnsavedModal() {
  if (!els.unsavedModal) return;
  els.unsavedModal.classList.add("hidden");
  els.unsavedModal.setAttribute("aria-hidden", "true");
  document.body.classList.remove("confirm-modal-open");
}

function showUnsavedModal() {
  return new Promise((resolve) => {
    if (!els.unsavedModal) {
      resolve("cancel");
      return;
    }
    state.unsavedResolver = resolve;
    els.unsavedModal.classList.remove("hidden");
    els.unsavedModal.setAttribute("aria-hidden", "false");
    document.body.classList.add("confirm-modal-open");
    els.unsavedSaveBtn?.focus();
  });
}

function resolveUnsavedModal(action) {
  hideUnsavedModal();
  const resolve = state.unsavedResolver;
  state.unsavedResolver = null;
  if (resolve) resolve(action);
}

async function confirmUnsavedChanges() {
  if (!isDirty()) return "discard";
  return showUnsavedModal();
}

async function guardUnsavedThen(actionFn) {
  const decision = await confirmUnsavedChanges();
  if (decision === "cancel") return false;
  if (decision === "save") {
    const ok = await saveProfile({ silent: false });
    if (!ok) return false;
  } else {
    markClean();
  }
  await actionFn();
  return true;
}

function initUnsavedModal() {
  els.unsavedSaveBtn?.addEventListener("click", () => resolveUnsavedModal("save"));
  els.unsavedDiscardBtn?.addEventListener("click", () => resolveUnsavedModal("discard"));
  els.unsavedModal?.querySelectorAll("[data-unsaved-cancel]").forEach((el) => {
    el.addEventListener("click", () => resolveUnsavedModal("cancel"));
  });

  document.addEventListener("keydown", (e) => {
    if (e.key !== "Escape") return;
    if (!els.unsavedModal || els.unsavedModal.classList.contains("hidden")) return;
    e.preventDefault();
    resolveUnsavedModal("cancel");
  });

  window.addEventListener("beforeunload", (e) => {
    if (!isDirty()) return;
    e.preventDefault();
    e.returnValue = "";
  });
}

function startAutoSave() {
  if (state.autoSaveTimer) clearInterval(state.autoSaveTimer);
  state.autoSaveTimer = setInterval(() => {
    if (!isDirty() || state.saving || !state.currentId) return;
    saveProfile({ silent: true, auto: true });
  }, 60_000);
}

async function loadProfiles(selectId) {
  let profiles = await API.listProfiles();
  renderProfileSelect(profiles);

  if (profiles.length === 0) {
    showToast(t("profilesLoadError"));
    return;
  }

  const targetId = selectId || state.currentId || profiles[0].id;
  const exists = profiles.some((profile) => profile.id === targetId);
  const profile = await API.getProfile(exists ? targetId : profiles[0].id);
  state.currentId = profile.id;
  setProfile(profile);
  markClean();
  await updatePreview();
}

async function saveProfile(options = {}) {
  if (!state.currentId || state.saving) return false;
  const profile = collectProfileFromForm();
  state.saving = true;
  try {
    const saved = await API.updateProfile(state.currentId, profile);
    applyServerMeta(saved);
    state.photoPreviewDirty = false;
    markClean();
    if (options.auto) {
      showToast(t("autoSaved"));
    } else if (!options.silent) {
      showToast(t("profileSaved"));
    }
    const profiles = await API.listProfiles();
    renderProfileSelect(profiles);
    els.profileSelect.value = String(state.currentId);
    return true;
  } catch (error) {
    showToast(error.message);
    return false;
  } finally {
    state.saving = false;
  }
}

async function createNewProfile() {
  await guardUnsavedThen(async () => {
    const profile = emptyProfile();
    profile.name = `CV ${new Date().toLocaleDateString("tr-TR")}`;
    const created = await API.createProfile(profile);
    state.currentId = created.id;
    setProfile(created);
    const profiles = await API.listProfiles();
    renderProfileSelect(profiles);
    els.profileSelect.value = String(created.id);
    markClean();
    await updatePreview();
    showToast(t("newProfileCreated"));
  });
}

async function duplicateCurrentProfile() {
  if (!state.currentId) return;

  await guardUnsavedThen(async () => {
    try {
      const duplicated = await API.duplicateProfile(state.currentId);
      state.currentId = duplicated.id;
      setProfile(duplicated);
      const profiles = await API.listProfiles();
      renderProfileSelect(profiles);
      els.profileSelect.value = String(duplicated.id);
      markClean();
      await updatePreview();
      showToast(t("profileDuplicated"));
    } catch (error) {
      showToast(error.message);
    }
  });
}

async function deleteCurrentProfile() {
  if (!state.currentId) return;
  if (!confirm(t("confirmDelete"))) return;

  try {
    await API.deleteProfile(state.currentId);
    state.currentId = null;
    markClean();
    await loadProfiles();
    showToast(t("profileDeleted"));
  } catch (error) {
    showToast(error.message);
  }
}

async function translateProfile() {
  if (!state.currentId) return;

  const targetLang = window.I18n?.getTargetTranslateLang?.();
  if (!targetLang) return;

  const storageSupported = window.I18n?.isTargetLangStorageSupported?.(targetLang) ?? true;
  collectProfileFromForm();
  const sourceLang = normalizeProfileLang(state.profile.language);
  const sourceContent = getContentForLang(state.profile, sourceLang);
  if (isContentEmpty(sourceContent)) {
    const sourceName = getEditingLanguageLabel(sourceLang);
    showToast(t("enterLangContent", { lang: sourceName }) || (sourceLang === "tr" ? t("enterTrContent") : t("enterEnContent")));
    return;
  }

  const btn = els.translateBtn;
  const targetName = window.I18n?.getLanguageDisplayName?.(targetLang) || targetLang;
  const defaultLabel = t("translateToLang", { lang: targetName });

  const provider = getTranslateProvider();
  if (isAiAssistantEnabled()) {
    showSection("ai-assistant");
    window.AIChat?.startTranslateProgress?.();
  }

  try {
    if (btn) {
      btn.disabled = true;
      btn.textContent = t("translating");
    }
    saveApiKey(API_KEY_CONFIG[0]);

    await API.updateProfile(state.currentId, state.profile);

    const translated = await API.translateProfile(state.currentId, {
      targetLang,
      sourceLang,
      apiKey: provider === "openai" ? getOpenAIApiKey() : getGoogleApiKey(),
      provider,
      model: localStorage.getItem("openaiSelectedModel") || "gpt-4o-mini",
    });

    window.AIChat?.completeTranslateProgress?.(true, translated.costUSD || 0, null, {
      provider,
      usage: translated.usage,
    });

    if (storageSupported && translated.storageSupported !== false) {
      setProfile(translated);
      markClean();
      await updatePreview();
      if (targetLang === "en") {
        showToast(t("cvTranslatedEn"));
      } else if (targetLang === "tr") {
        showToast(t("cvTranslatedTr"));
      } else {
        showToast(t("cvTranslatedTo", { lang: targetName }));
      }
    } else {
      showToast(t("cvTranslatedStorageWarning"));
    }
  } catch (error) {
    window.AIChat?.completeTranslateProgress?.(false, 0, error.message);
    showToast(error.message);
  } finally {
    if (btn) {
      btn.disabled = false;
      btn.textContent = defaultLabel;
    }
  }
}

async function downloadPDF() {
  const profile = collectProfileFromForm();
  try {
    els.pdfBtn.disabled = true;
    els.pdfBtn.textContent = t("pdfPreparing");
    // Persist latest edits before export so PDF matches what you see.
    if (state.currentId && isDirty()) {
      await saveProfile({ silent: true });
    }
    await API.downloadPDF(profile);
    showToast(t("pdfDownloaded"));
  } catch (error) {
    showToast(error.message);
  } finally {
    els.pdfBtn.disabled = false;
    els.pdfBtn.textContent = t("pdfDownload");
  }
}

async function exportCurrentProfileBackup() {
  if (!state.currentId) {
    showToast(t("selectProfileFirst"));
    return;
  }

  try {
    const profile = collectProfileFromForm();
    await API.exportCurrentProfile(state.currentId, profile.name || state.profile?.name);
  } catch (error) {
    showToast(t("exportProfileError") + (error.message ? `: ${error.message}` : ""));
  }
}

async function exportAllProfilesBackup() {
  try {
    await API.exportAllProfiles();
  } catch (error) {
    showToast(t("exportAllError") + (error.message ? `: ${error.message}` : ""));
  }
}

async function resetAllProfiles() {
  await guardUnsavedThen(async () => {
    if (!window.confirm(t("resetConfirm"))) {
      return;
    }

    try {
      const result = await API.resetAllProfiles();
      const selectId = result?.profile?.id;
      await loadProfiles(selectId);
      showToast(t("resetSuccess"));
    } catch (error) {
      showToast(t("resetError") + (error.message ? `: ${error.message}` : ""));
    }
  });
}

async function handleImportBackup(file) {
  if (!file) return;

  await guardUnsavedThen(async () => {
    const merge = window.confirm(t("importMergeConfirm"));
    let mode = "merge";
    if (!merge) {
      if (!window.confirm(t("importReplaceConfirm"))) {
        return;
      }
      mode = "replace";
    }

    try {
      const result = await API.importProfiles(file, mode);
      const imported = Array.isArray(result?.profiles) ? result.profiles : [];
      const selectId = imported[0]?.id;
      await loadProfiles(selectId);
      const count = result?.count ?? imported.length;
      showToast(t("importSuccess", { count }));
    } catch (error) {
      showToast(t("importError", { error: error.message || t("unknownError") }));
    } finally {
      if (els.importProfileInput) {
        els.importProfileInput.value = "";
      }
    }
  });
}

async function handlePhotoUpload(file) {
  if (!file || !state.currentId) return;
  try {
    // Only the photo is persisted by this endpoint, so keep the rest of the
    // in-progress edits (and their dirty state) untouched.
    const updated = await API.uploadPhoto(state.currentId, file);
    state.profile.photoBase64 = updated?.photoBase64 || "";
    applyServerMeta(updated);
    state.photoPreviewDirty = false;
    renderPhoto();
    schedulePreview();
    showToast(t("photoUploaded"));
  } catch (error) {
    showToast(error.message);
  }
}

function initEventListeners() {
  els.profileSelect.addEventListener("change", async () => {
    const nextId = Number(els.profileSelect.value);
    const previousId = state.currentId;
    if (previousId && nextId === previousId) return;

    if (isDirty()) {
      els.profileSelect.value = String(previousId);
      const decision = await confirmUnsavedChanges();
      if (decision === "cancel") return;
      if (decision === "save") {
        const ok = await saveProfile({ silent: false });
        if (!ok) return;
      } else {
        markClean();
      }
      els.profileSelect.value = String(nextId);
    }

    state.currentId = nextId;
    setProfile(await API.getProfile(state.currentId));
    markClean();
    await updatePreview();
  });

  els.cvLanguage?.addEventListener("change", () => {
    switchEditingLanguage(els.cvLanguage.value);
  });

  API_KEY_CONFIG.forEach((config) => {
    const input = getApiKeyInput(config.inputId);
    if (!input) return;
    input.addEventListener("change", () => {
      saveApiKey(config);
      if (config.storageKey === "openaiApiKey") {
        window.AIChat?.loadModels?.({ silent: true });
      }
    });
    input.addEventListener("blur", () => {
      saveApiKey(config);
      if (config.storageKey === "openaiApiKey") {
        window.AIChat?.loadModels?.({ silent: true });
      }
    });
    input.addEventListener("input", () => updateApiKeyStatus(config));
  });

  els.translateTargetSelect?.addEventListener("change", () => {
    const code = els.translateTargetSelect.value;
    window.I18n?.setActiveTargetTranslateLang?.(code);
  });

  els.translateBtn?.addEventListener("click", () => translateProfile());

  els.newProfileBtn.addEventListener("click", createNewProfile);
  els.duplicateProfileBtn.addEventListener("click", duplicateCurrentProfile);
  els.deleteProfileBtn.addEventListener("click", deleteCurrentProfile);
  els.saveBtn.addEventListener("click", () => saveProfile());
  els.pdfBtn.addEventListener("click", downloadPDF);
  els.exportProfileBtn?.addEventListener("click", exportCurrentProfileBackup);
  els.exportAllBtn?.addEventListener("click", exportAllProfilesBackup);
  els.importProfileInput?.addEventListener("change", (event) => {
    const file = event.target.files?.[0];
    if (file) {
      handleImportBackup(file);
    }
  });
  els.resetProfilesBtn?.addEventListener("click", resetAllProfiles);

  els.profileName.addEventListener("input", onEditorChange);
  document.querySelectorAll("[data-field]").forEach((input) => {
    input.addEventListener("input", onEditorChange);
  });
  document.querySelectorAll("[data-localized]").forEach((input) => {
    input.addEventListener("input", onEditorChange);
  });

  document.querySelectorAll("[data-add]").forEach((button) => {
    button.addEventListener("click", () => {
      const content = getActiveContent();
      const type = button.dataset.add;
      if (type === "experiences") {
        content.experiences.push({
          title: "",
          company: "",
          startDate: "",
          endDate: "",
          duration: "",
          location: "",
          description: "",
          highlights: [],
        });
        renderExperiences();
      } else if (type === "education") {
        content.education.push({
          degree: "",
          institution: "",
          startDate: "",
          endDate: "",
          description: "",
        });
        renderEducation();
      } else if (type === "projects") {
        content.projects.push({ name: "", url: "", description: "" });
        renderProjects();
      } else if (type === "skillGroups") {
        content.skillGroups.push({ category: "", skills: [] });
        renderSkillGroups();
      }
      onEditorChange();
    });
  });

  els.photoInput.addEventListener("change", async (event) => {
    const file = event.target.files?.[0];
    if (file) {
      await handlePhotoUpload(file);
    }
    event.target.value = "";
  });

  [els.photoSize, els.photoBorderWidth].forEach((input) => {
    input.addEventListener("input", () => {
      syncPhotoSettingLabels();
      applyPhotoPreviewStyles();
      onEditorChange();
    });
  });

  els.photoBorderColor.addEventListener("input", () => {
    applyPhotoPreviewStyles();
    onEditorChange();
  });

  els.removePhotoBtn.addEventListener("click", async () => {
    state.profile.photoBase64 = "";
    state.photoPreviewDirty = true;
    renderPhoto();
    schedulePreview();
    if (state.currentId) {
      try {
        applyServerMeta(await API.updateProfile(state.currentId, collectProfileFromForm()));
        state.photoPreviewDirty = false;
        markClean();
        showToast(t("photoRemoved"));
      } catch (error) {
        markDirty();
        showToast(error.message);
      }
    } else {
      markDirty();
    }
  });
}

async function onProfileUpdated(profile) {
  state.currentId = profile.id;
  setProfile(profile);
  const profiles = await API.listProfiles();
  renderProfileSelect(profiles);
  els.profileSelect.value = String(profile.id);
  markClean();
  await updatePreview();
}

async function onAIChatResult(response) {
  const profile = response.profile;
  if (!profile?.id) {
    if (response.profileId) {
      await loadProfiles(response.profileId);
    }
    return;
  }

  state.currentId = profile.id;
  setProfile(profile, { render: false });

  if (response.language) {
    state.profile.language = normalizeProfileLang(response.language);
    renderEditingLanguageOptions();
    if (els.cvLanguage) {
      els.cvLanguage.value = state.profile.language;
    }
  }

  const profiles = await API.listProfiles();
  renderProfileSelect(profiles);
  els.profileSelect.value = String(profile.id);
  renderForm();
  markClean();
  await updatePreview();

  const createAction = (response.actions || []).find(
    (a) => a.tool === "create_profile" && a.status === "success",
  );
  const translateAction = (response.actions || []).find(
    (a) => a.tool === "translate_profile" && a.status === "success",
  );

  if (createAction) {
    showToast(t("newProfileCreatedNamed", { name: profile.name }));
  } else if (translateAction) {
    showToast(response.language === "en" ? t("profileTranslatedEn") : t("profileTranslatedTr"));
  } else if (response.profileUpdated) {
    showToast(t("profileUpdatedToast"));
  }
}

window.CVEditor = {
  getCurrentProfileId: () => state.currentId,
  getEditingLanguage: () => normalizeProfileLang(state.profile?.language),
  getProfile: () => state.profile,
  getOpenAIApiKey,
  isOpenAIConfigured,
  isAiAssistantEnabled,
  showToast,
  onProfileUpdated,
  onAIChatResult,
  showSection,
  refreshTranslateUI,
};

function onAppLanguageChange() {
  applyTheme(getStoredTheme());
  updateLanguageUI();
  updateAllApiKeyStatuses();
  renderAllLists();
}

function onTargetTranslateLangsChanged() {
  renderEditingLanguageOptions();
  refreshTranslateUI();
}

document.addEventListener("DOMContentLoaded", async () => {
  window.I18n?.initI18n();
  window.I18n?.onLanguageChange(onAppLanguageChange);
  window.I18n?.onTargetTranslateLangChange?.(() => onTargetTranslateLangsChanged());
  initTheme();
  loadApiKeysFromStorage();
  initSectionNav();
  initAiAssistantSetting();
  initUnsavedModal();
  initEventListeners();
  startAutoSave();
  try {
    await loadEnvSettings();
    await loadProfiles();
    refreshTranslateUI();
    window.AIChat?.loadModels?.({ silent: true });
  } catch (error) {
    showToast(error.message);
  }
});
