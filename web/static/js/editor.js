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
  envSettings: {},
  activeSection: "general",
};

const SECTION_IDS = [
  "general",
  "personal",
  "summary",
  "experience",
  "education",
  "projects",
  "skills",
  "photo",
  "api-keys",
];

const THEME_STORAGE_KEY = "theme";

const els = {
  profileSelect: document.getElementById("profileSelect"),
  newProfileBtn: document.getElementById("newProfileBtn"),
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
  translateBtn: document.getElementById("translateBtn"),
  translateToTrBtn: document.getElementById("translateToTrBtn"),
  experiencesList: document.getElementById("experiencesList"),
  educationList: document.getElementById("educationList"),
  projectsList: document.getElementById("projectsList"),
  skillGroupsList: document.getElementById("skillGroupsList"),
  toast: document.getElementById("toast"),
  editorSidebar: document.querySelector(".editor-sidebar"),
  formSections: document.querySelectorAll(".form-section[data-section-id]"),
};

function showSection(sectionId) {
  if (!SECTION_IDS.includes(sectionId)) return;

  state.activeSection = sectionId;

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
    content.scrollTop = 0;
  }
}

function initSectionNav() {
  els.editorSidebar?.querySelectorAll(".editor-nav-item").forEach((item) => {
    item.addEventListener("click", () => {
      showSection(item.dataset.section);
    });
  });
  showSection(state.activeSection);
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
    name: "Yeni CV",
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

function normalizeLanguage(profile) {
  profile.language = profile.language === "en" ? "en" : "tr";
}

function activeContentKey(language) {
  return language === "en" ? "contentEN" : "contentTR";
}

function getActiveContent(profile = state.profile) {
  const key = activeContentKey(profile.language);
  if (!profile[key]) {
    profile[key] = emptyLocalizedContent();
  }
  return profile[key];
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
  });
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
  return localStorage.getItem(config.storageKey) || "";
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
    statusEl.textContent = "Tarayıcıda kayıtlı (öncelikli)";
    statusEl.classList.add("api-key-status--local");
  } else if (hasEnv) {
    statusEl.textContent = ".env ile yapılandırıldı";
    statusEl.classList.add("api-key-status--env");
  } else {
    statusEl.textContent = "Yapılandırılmadı";
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
    const label = isDark ? "Açık mod" : "Karanlık mod";
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

function updateLanguageUI() {
  const lang = state.profile.language === "en" ? "en" : "tr";
  if (els.cvLanguage) {
    els.cvLanguage.value = lang;
  }

  const enEmpty = isContentEmpty(state.profile.contentEN);
  if (els.languageEmptyHint) {
    els.languageEmptyHint.classList.toggle("hidden", !(lang === "en" && enEmpty));
  }

  if (els.translateBtn) {
    els.translateBtn.classList.toggle("hidden", lang === "en");
    els.translateBtn.textContent = "İngilizceye Çevir";
  }
  if (els.translateToTrBtn) {
    els.translateToTrBtn.classList.toggle("hidden", lang === "tr");
    els.translateToTrBtn.textContent = "Türkçeye Çevir";
  }
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
    schedulePreview();
  }
}

function reorderEducation(fromIndex, toIndex) {
  const content = getActiveContent();
  if (reorderItem(content.education, fromIndex, toIndex)) {
    renderEducation();
    schedulePreview();
  }
}

function reorderProjects(fromIndex, toIndex) {
  const content = getActiveContent();
  if (reorderItem(content.projects, fromIndex, toIndex)) {
    renderProjects();
    schedulePreview();
  }
}

function reorderSkillGroups(fromIndex, toIndex) {
  const content = getActiveContent();
  if (reorderItem(content.skillGroups, fromIndex, toIndex)) {
    renderSkillGroups();
    schedulePreview();
  }
}

function createReorderButton(label, disabled, onClick) {
  const btn = document.createElement("button");
  btn.type = "button";
  btn.className = "btn btn-secondary btn-sm btn-icon-reorder";
  btn.setAttribute("aria-label", label);
  btn.title = label;
  btn.disabled = disabled;
  btn.innerHTML = label === "Yukarı taşı"
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
  dragHandle.setAttribute("aria-label", "Sürükleyerek sırala");
  dragHandle.title = "Sürükleyerek sırala";
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
      createReorderButton("Yukarı taşı", reorder.index === 0, () => reorder.onMove(-1)),
      createReorderButton("Aşağı taşı", reorder.index === reorder.total - 1, () => reorder.onMove(1))
    );
  }

  const removeBtn = document.createElement("button");
  removeBtn.type = "button";
  removeBtn.className = "btn btn-danger btn-sm";
  removeBtn.textContent = "Sil";
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
    schedulePreview();
  });
  label.appendChild(input);
  return label;
}

function renderExperiences() {
  const content = getActiveContent();
  els.experiencesList.innerHTML = "";
  const total = content.experiences.length;
  content.experiences.forEach((exp, index) => {
    const card = createCard(`Deneyim ${index + 1}`, () => {
      content.experiences.splice(index, 1);
      renderExperiences();
      schedulePreview();
    }, (container) => {
      const fields = document.createElement("div");
      fields.className = "card-grid";
      fields.append(
        bindInput("Unvan", exp.title, (v) => (exp.title = v)),
        bindInput("Şirket", exp.company, (v) => (exp.company = v)),
        bindInput("Başlangıç", exp.startDate, (v) => (exp.startDate = v)),
        bindInput("Bitiş", exp.endDate, (v) => (exp.endDate = v)),
        bindInput("Süre", exp.duration, (v) => (exp.duration = v)),
        bindInput("Konum", exp.location, (v) => (exp.location = v)),
        bindInput("Şirket Açıklaması", exp.description, (v) => (exp.description = v), { type: "textarea", rows: 2 }),
        bindInput("Maddeler (her satır bir madde)", (exp.highlights || []).join("\n"), (v) => {
          exp.highlights = v.split("\n").map((line) => line.trim()).filter(Boolean);
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
    const card = createCard(`Eğitim ${index + 1}`, () => {
      content.education.splice(index, 1);
      renderEducation();
      schedulePreview();
    }, (container) => {
      const fields = document.createElement("div");
      fields.className = "card-grid";
      fields.append(
        bindInput("Program", edu.degree, (v) => (edu.degree = v)),
        bindInput("Kurum", edu.school, (v) => (edu.school = v)),
        bindInput("Başlangıç", edu.startDate, (v) => (edu.startDate = v)),
        bindInput("Bitiş", edu.endDate, (v) => (edu.endDate = v)),
        bindInput("Açıklama", edu.description, (v) => (edu.description = v), { type: "textarea", rows: 3 })
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
    const card = createCard(`Proje ${index + 1}`, () => {
      content.projects.splice(index, 1);
      renderProjects();
      schedulePreview();
    }, (container) => {
      const fields = document.createElement("div");
      fields.className = "card-grid";
      fields.append(
        bindInput("Proje Adı", project.name, (v) => (project.name = v)),
        bindInput("URL", project.url, (v) => (project.url = v)),
        bindInput("Açıklama", project.description, (v) => (project.description = v), { type: "textarea", rows: 3 })
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
    const card = createCard(`Beceri Grubu ${index + 1}`, () => {
      content.skillGroups.splice(index, 1);
      renderSkillGroups();
      schedulePreview();
    }, (container) => {
      const fields = document.createElement("div");
      fields.className = "card-grid";
      fields.append(
        bindInput("Kategori", group.category, (v) => (group.category = v)),
        bindInput("Beceriler (her satır bir beceri)", (group.skills || []).join("\n"), (v) => {
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

function collectProfileFromForm(options = {}) {
  const editingLang = options.language ?? (els.cvLanguage?.value === "en" ? "en" : "tr");
  const content = state.profile[activeContentKey(editingLang)] || emptyLocalizedContent();
  state.profile[activeContentKey(editingLang)] = content;

  state.profile.name = els.profileName.value.trim() || "Yeni CV";
  if (!options.skipLanguageUpdate && els.cvLanguage) {
    state.profile.language = els.cvLanguage.value === "en" ? "en" : "tr";
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
  const currentLang = state.profile.language === "en" ? "en" : "tr";
  const normalizedNext = nextLang === "en" ? "en" : "tr";
  if (currentLang === normalizedNext) {
    return;
  }

  collectProfileFromForm({ language: currentLang, skipLanguageUpdate: true });
  state.profile.language = normalizedNext;
  if (els.cvLanguage) {
    els.cvLanguage.value = normalizedNext;
  }
  renderForm();
  schedulePreview();
}

async function updatePreview() {
  const profile = collectProfileFromForm();
  try {
    const html = await API.preview(profile);
    els.previewFrame.srcdoc = html;
  } catch (error) {
    console.error(error);
  }
}

function schedulePreview() {
  clearTimeout(state.previewTimer);
  state.previewTimer = setTimeout(updatePreview, 300);
}

async function loadProfiles(selectId) {
  const profiles = await API.listProfiles();
  renderProfileSelect(profiles);

  if (profiles.length === 0) {
    const created = await API.createProfile(emptyProfile());
    state.currentId = created.id;
    state.profile = created;
    normalizeProfile(state.profile);
  } else {
    const targetId = selectId || state.currentId || profiles[0].id;
    const profile = await API.getProfile(targetId);
    state.currentId = profile.id;
    state.profile = profile;
    normalizeProfile(state.profile);
  }

  renderForm();
  await updatePreview();
}

async function saveProfile() {
  const profile = collectProfileFromForm();
  try {
    const saved = await API.updateProfile(state.currentId, profile);
    state.profile = saved;
    showToast("Profil kaydedildi");
    const profiles = await API.listProfiles();
    renderProfileSelect(profiles);
    els.profileSelect.value = String(state.currentId);
  } catch (error) {
    showToast(error.message);
  }
}

async function createNewProfile() {
  const profile = emptyProfile();
  profile.name = `CV ${new Date().toLocaleDateString("tr-TR")}`;
  const created = await API.createProfile(profile);
  state.currentId = created.id;
  state.profile = created;
  const profiles = await API.listProfiles();
  renderProfileSelect(profiles);
  els.profileSelect.value = String(created.id);
  renderForm();
  await updatePreview();
  showToast("Yeni profil oluşturuldu");
}

async function deleteCurrentProfile() {
  if (!state.currentId) return;
  if (!confirm("Bu profili silmek istediğinize emin misiniz?")) return;

  try {
    await API.deleteProfile(state.currentId);
    state.currentId = null;
    await loadProfiles();
    showToast("Profil silindi");
  } catch (error) {
    showToast(error.message);
  }
}

async function translateProfile(targetLang) {
  if (!state.currentId) return;

  collectProfileFromForm();
  const sourceLang = targetLang === "en" ? "tr" : "en";
  const sourceContent = state.profile[activeContentKey(sourceLang)];
  if (isContentEmpty(sourceContent)) {
    showToast(sourceLang === "tr" ? "Önce Türkçe içerik girin" : "Önce İngilizce içerik girin");
    return;
  }

  const btn = targetLang === "en" ? els.translateBtn : els.translateToTrBtn;
  const defaultLabel = targetLang === "en" ? "İngilizceye Çevir" : "Türkçeye Çevir";

  try {
    if (btn) {
      btn.disabled = true;
      btn.textContent = "Çevriliyor...";
    }
    saveApiKey(API_KEY_CONFIG[0]);

    await API.updateProfile(state.currentId, state.profile);

    const translated = await API.translateProfile(state.currentId, {
      targetLang,
      apiKey: getGoogleApiKey(),
    });

    state.profile = translated;
    normalizeProfile(state.profile);
    renderForm();
    await updatePreview();
    showToast(targetLang === "en" ? "CV İngilizceye çevrildi" : "CV Türkçeye çevrildi");
  } catch (error) {
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
    els.pdfBtn.textContent = "PDF hazırlanıyor...";
    await API.downloadPDF(profile);
    showToast("PDF indirildi");
  } catch (error) {
    showToast(error.message);
  } finally {
    els.pdfBtn.disabled = false;
    els.pdfBtn.textContent = "PDF İndir";
  }
}

async function handlePhotoUpload(file) {
  if (!file || !state.currentId) return;
  try {
    const updated = await API.uploadPhoto(state.currentId, file);
    state.profile = updated;
    renderPhoto();
    schedulePreview();
    showToast("Fotoğraf yüklendi");
  } catch (error) {
    showToast(error.message);
  }
}

function initEventListeners() {
  els.profileSelect.addEventListener("change", async () => {
    state.currentId = Number(els.profileSelect.value);
    state.profile = await API.getProfile(state.currentId);
    normalizeProfile(state.profile);
    renderForm();
    await updatePreview();
  });

  els.cvLanguage?.addEventListener("change", () => {
    switchEditingLanguage(els.cvLanguage.value);
  });

  API_KEY_CONFIG.forEach((config) => {
    const input = getApiKeyInput(config.inputId);
    if (!input) return;
    input.addEventListener("change", () => saveApiKey(config));
    input.addEventListener("blur", () => saveApiKey(config));
    input.addEventListener("input", () => updateApiKeyStatus(config));
  });

  els.translateBtn?.addEventListener("click", () => translateProfile("en"));
  els.translateToTrBtn?.addEventListener("click", () => translateProfile("tr"));

  els.newProfileBtn.addEventListener("click", createNewProfile);
  els.deleteProfileBtn.addEventListener("click", deleteCurrentProfile);
  els.saveBtn.addEventListener("click", saveProfile);
  els.pdfBtn.addEventListener("click", downloadPDF);

  els.profileName.addEventListener("input", schedulePreview);
  document.querySelectorAll("[data-field]").forEach((input) => {
    input.addEventListener("input", schedulePreview);
  });
  document.querySelectorAll("[data-localized]").forEach((input) => {
    input.addEventListener("input", schedulePreview);
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
          school: "",
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
      schedulePreview();
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
      schedulePreview();
    });
  });

  els.photoBorderColor.addEventListener("input", () => {
    applyPhotoPreviewStyles();
    schedulePreview();
  });

  els.removePhotoBtn.addEventListener("click", async () => {
    state.profile.photoBase64 = "";
    renderPhoto();
    schedulePreview();
    if (state.currentId) {
      try {
        await API.updateProfile(state.currentId, collectProfileFromForm());
        showToast("Fotoğraf kaldırıldı");
      } catch (error) {
        showToast(error.message);
      }
    }
  });
}

document.addEventListener("DOMContentLoaded", async () => {
  initTheme();
  loadApiKeysFromStorage();
  initSectionNav();
  initEventListeners();
  try {
    await loadEnvSettings();
    await loadProfiles();
  } catch (error) {
    showToast(error.message);
  }
});
