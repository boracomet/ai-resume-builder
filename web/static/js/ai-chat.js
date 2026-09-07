const t = (key, params) => window.I18n?.t(key, params) ?? key;

const AI_MODEL_STORAGE_KEY = "openaiSelectedModel";
const TRANSLATE_PROVIDER_STORAGE_KEY = "translateProvider";
const AI_SESSION_COST_KEY = "aiSessionCostUSD";
const AI_LAST_COST_KEY = "aiLastRequestCostUSD";
const TRANSLATE_PROGRESS_INTERVAL_MS = 450;
const PDF_PREVIEW_LIMIT = 500;

function getTranslateProgressStages() {
  return window.I18n?.getTranslateProgressStages?.() ?? [
    "Çevirinizi yapıyorum...",
    "Özet çevriliyor...",
    "Deneyimler çevriliyor...",
    "Eğitim çevriliyor...",
    "Beceriler çevriliyor...",
    "Kaydediliyor...",
  ];
}

const AIChat = {
  messages: [],
  pendingSuggestion: null,
  loadingModels: false,
  uploadingPdf: false,
  thinkingEl: null,
  translateProgressEl: null,
  translateProgressContent: null,
  translateProgressStageIndex: 0,
  translateProgressDone: false,
  translateProgressTimer: null,
  els: {},

  init() {
    this.els = {
      toggleBtn: document.getElementById("aiChatToggle"),
      messages: document.getElementById("aiChatMessages"),
      input: document.getElementById("aiChatInput"),
      sendBtn: document.getElementById("aiChatSend"),
      pdfBtn: document.getElementById("aiChatPdfBtn"),
      pdfInput: document.getElementById("aiChatPdfInput"),
      pdfStatus: document.getElementById("aiChatPdfStatus"),
      applyBtn: document.getElementById("aiApplyBtn"),
      settingsModelSelect: document.getElementById("aiSettingsModelSelect"),
      settingsFetchModelsBtn: document.getElementById("aiSettingsFetchModelsBtn"),
      settingsModelsStatus: document.getElementById("aiSettingsModelsStatus"),
      settingsTestBtn: document.getElementById("aiSettingsTestBtn"),
      settingsTestStatus: document.getElementById("aiSettingsTestStatus"),
      translateProvider: document.getElementById("translateProvider"),
      sessionCostValue: document.getElementById("aiSessionCostValue"),
      lastRequestCostValue: document.getElementById("aiLastRequestCostValue"),
      sessionCostReset: document.getElementById("aiSessionCostReset"),
      apiSetupOpenBtn: document.getElementById("apiSetupOpenBtn"),
      apiSetupModal: document.getElementById("apiSetupModal"),
    };

    this.loadTranslateProvider();
    this.populateModelPlaceholder();
    this.updateCostDisplay();
    this.bindEvents();
    this.bindApiSetupModal();
  },

  formatCostUSD(amount) {
    const value = Number(amount) || 0;
    if (value >= 0.01) return `$${value.toFixed(2)}`;
    if (value >= 0.0001) return `$${value.toFixed(4)}`;
    return `$${value.toFixed(6)}`;
  },

  getSessionCost() {
    return parseFloat(sessionStorage.getItem(AI_SESSION_COST_KEY) || "0") || 0;
  },

  getLastRequestCost() {
    const raw = sessionStorage.getItem(AI_LAST_COST_KEY);
    if (raw === null || raw === "") return null;
    return parseFloat(raw) || 0;
  },

  trackCost(response) {
    const cost = Number(response?.costUSD) || 0;
    if (cost <= 0) return;

    const session = this.getSessionCost() + cost;
    sessionStorage.setItem(AI_SESSION_COST_KEY, String(session));
    sessionStorage.setItem(AI_LAST_COST_KEY, String(cost));
    this.updateCostDisplay();
  },

  resetSessionCost() {
    sessionStorage.removeItem(AI_SESSION_COST_KEY);
    sessionStorage.removeItem(AI_LAST_COST_KEY);
    this.updateCostDisplay();
  },

  updateCostDisplay() {
    if (this.els.sessionCostValue) {
      this.els.sessionCostValue.textContent = this.formatCostUSD(this.getSessionCost());
    }
    if (this.els.lastRequestCostValue) {
      const last = this.getLastRequestCost();
      this.els.lastRequestCostValue.textContent = last === null ? "—" : this.formatCostUSD(last);
    }
  },

  getApiKey() {
    if (window.CVEditor?.getOpenAIApiKey) {
      return window.CVEditor.getOpenAIApiKey();
    }
    return localStorage.getItem("openaiApiKey") || "";
  },

  getSelectedModel() {
    const fromSettings = this.els.settingsModelSelect?.value;
    const saved = localStorage.getItem(AI_MODEL_STORAGE_KEY);
    return fromSettings || saved || "gpt-4o-mini";
  },

  saveSelectedModel(model) {
    if (!model) return;
    localStorage.setItem(AI_MODEL_STORAGE_KEY, model);
    if (this.els.settingsModelSelect) this.els.settingsModelSelect.value = model;
  },

  loadTranslateProvider() {
    const saved = localStorage.getItem(TRANSLATE_PROVIDER_STORAGE_KEY) || "google";
    if (this.els.translateProvider) {
      this.els.translateProvider.value = saved;
    }
  },

  saveTranslateProvider() {
    const value = this.els.translateProvider?.value || "google";
    localStorage.setItem(TRANSLATE_PROVIDER_STORAGE_KEY, value);
  },

  isOpenAIAvailable() {
    const apiKey = this.getApiKey();
    if (apiKey) return true;
    if (window.CVEditor?.isOpenAIConfigured) {
      return window.CVEditor.isOpenAIConfigured();
    }
    return false;
  },

  populateModelPlaceholder() {
    const saved = localStorage.getItem(AI_MODEL_STORAGE_KEY) || "gpt-4o-mini";
    const select = this.els.settingsModelSelect;
    if (!select) return;
    select.innerHTML = "";
    const opt = document.createElement("option");
    opt.value = saved;
    opt.textContent = saved;
    select.appendChild(opt);
    select.value = saved;
  },

  getFetchModelButtons() {
    return [this.els.settingsFetchModelsBtn].filter(Boolean);
  },

  setFetchButtonsLoading(loading) {
    this.getFetchModelButtons().forEach((btn) => {
      btn.disabled = loading;
      btn.textContent = loading ? t("loading") : t("fetchModels");
    });
  },

  getModelStatusElements(preferredEl = null) {
    const elements = [preferredEl, this.els.settingsModelsStatus].filter(Boolean);
    return [...new Set(elements)];
  },

  reportModelsStatus(preferredEl, text, type, { toast = false } = {}) {
    this.getModelStatusElements(preferredEl).forEach((el) => this.setStatus(el, text, type));
    if (toast && window.CVEditor?.showToast) {
      window.CVEditor.showToast(text);
    }
  },

  onAppLanguageChange() {
    if (!this.loadingModels) {
      this.setFetchButtonsLoading(false);
    }
    if (!this.els.sendBtn?.disabled) {
      this.els.sendBtn.textContent = t("send");
    }
    if (!this.els.applyBtn?.disabled) {
      this.els.applyBtn.textContent = t("apply");
    }
    if (!this.uploadingPdf && this.els.pdfBtn) {
      this.els.pdfBtn.textContent = t("addPdf");
    }
    if (!this.els.settingsTestBtn?.disabled) {
      this.els.settingsTestBtn.textContent = t("testConnection");
    }
  },

  bindEvents() {
    window.I18n?.onLanguageChange(() => this.onAppLanguageChange());
    this.els.toggleBtn?.addEventListener("click", () => this.navigateToSection());
    this.els.sendBtn?.addEventListener("click", () => this.sendMessage());
    this.els.input?.addEventListener("keydown", (e) => {
      if (e.key === "Enter" && !e.shiftKey) {
        e.preventDefault();
        this.sendMessage();
      }
    });
    this.els.settingsTestBtn?.addEventListener("click", () => this.testConnection(this.els.settingsTestStatus));
    this.els.settingsFetchModelsBtn?.addEventListener("click", () => this.loadModels({
      statusEl: this.els.settingsModelsStatus,
      manual: true,
    }));
    this.els.applyBtn?.addEventListener("click", () => this.applySuggestion());
    this.els.pdfBtn?.addEventListener("click", () => this.els.pdfInput?.click());
    this.els.pdfInput?.addEventListener("change", (e) => {
      const file = e.target.files?.[0];
      e.target.value = "";
      if (file) this.uploadPDF(file);
    });
    this.els.settingsModelSelect?.addEventListener("change", () => {
      this.saveSelectedModel(this.els.settingsModelSelect.value);
    });
    this.els.translateProvider?.addEventListener("change", () => this.saveTranslateProvider());
    this.els.sessionCostReset?.addEventListener("click", () => this.resetSessionCost());
  },

  bindApiSetupModal() {
    const modal = this.els.apiSetupModal;
    if (!modal) return;

    const openModal = () => {
      modal.classList.remove("hidden");
      modal.setAttribute("aria-hidden", "false");
      document.body.classList.add("api-setup-modal-open");
      modal.querySelector("#openaiApiKey, #googleApiKey, input[type='password']")?.focus();
    };

    const closeModal = () => {
      modal.classList.add("hidden");
      modal.setAttribute("aria-hidden", "true");
      document.body.classList.remove("api-setup-modal-open");
      this.els.apiSetupOpenBtn?.focus();
    };

    this.els.apiSetupOpenBtn?.addEventListener("click", openModal);
    modal.querySelectorAll("[data-api-setup-close]").forEach((el) => {
      el.addEventListener("click", closeModal);
    });

    document.addEventListener("keydown", (e) => {
      if (e.key === "Escape" && !modal.classList.contains("hidden")) {
        e.preventDefault();
        closeModal();
      }
    });
  },

  navigateToSection() {
    if (window.CVEditor?.showSection) {
      window.CVEditor.showSection("ai-assistant");
    }
    requestAnimationFrame(() => {
      this.els.input?.focus();
    });
  },

  async loadModels(options = {}) {
    const { statusEl = null, manual = false, silent = false } = options;
    const apiKey = this.getApiKey();
    const select = this.els.settingsModelSelect;
    if (!select) {
      const message = t("modelSelectorsNotFound");
      console.error("[AIChat] loadModels:", message);
      if (manual) {
        this.reportModelsStatus(statusEl, message, "error", { toast: true });
      }
      return;
    }

    const canAutoLoad = this.isOpenAIAvailable();
    if (!manual && !canAutoLoad) {
      return;
    }

    if (manual && !canAutoLoad) {
      const message = t("openaiNotConfigured");
      console.error("[AIChat] loadModels:", message);
      this.reportModelsStatus(statusEl, message, "error", { toast: true });
      return;
    }

    if (this.loadingModels) {
      if (manual) {
        this.reportModelsStatus(statusEl, t("modelsAlreadyLoading"), "pending", { toast: false });
      }
      return;
    }
    this.loadingModels = true;
    this.setFetchButtonsLoading(true);

    select.innerHTML = `<option value="">${t("modelsLoading")}</option>`;
    select.disabled = true;

    if (!silent) {
      this.reportModelsStatus(statusEl, t("modelsFetching"), "pending", { toast: false });
    }

    try {
      const data = await API.listAIModels(apiKey);
      const models = data.models || [];
      const saved = localStorage.getItem(AI_MODEL_STORAGE_KEY);
      const defaultModel = saved && models.includes(saved) ? saved : (models[0] || "gpt-4o-mini");

      select.innerHTML = "";
      if (models.length === 0) {
        const opt = document.createElement("option");
        opt.value = defaultModel;
        opt.textContent = defaultModel;
        select.appendChild(opt);
      } else {
        models.forEach((model) => {
          const opt = document.createElement("option");
          opt.value = model;
          opt.textContent = model;
          select.appendChild(opt);
        });
      }
      select.value = defaultModel;
      select.disabled = false;
      this.saveSelectedModel(defaultModel);

      if (!silent) {
        const count = models.length;
        this.reportModelsStatus(
          statusEl,
          count > 0 ? t("modelsLoaded", { count }) : t("defaultModelUsed"),
          "success",
          { toast: manual && count > 0 },
        );
      }
    } catch (error) {
      console.error("[AIChat] loadModels failed:", error);
      const fallback = localStorage.getItem(AI_MODEL_STORAGE_KEY) || "gpt-4o-mini";
      select.innerHTML = "";
      const opt = document.createElement("option");
      opt.value = fallback;
      opt.textContent = `${fallback} ${t("defaultSuffix")}`;
      select.appendChild(opt);
      select.disabled = false;
      select.value = fallback;

      if (!silent || manual) {
        this.reportModelsStatus(statusEl, error.message, "error", { toast: manual });
      }
    } finally {
      this.loadingModels = false;
      this.setFetchButtonsLoading(false);
    }
  },

  async testConnection(statusEl) {
    if (!this.isOpenAIAvailable()) {
      const message = t("openaiKeyRequired");
      console.error("[AIChat] testConnection:", message);
      this.setStatus(statusEl, message, "error");
      if (window.CVEditor?.showToast) {
        window.CVEditor.showToast(message);
      }
      return;
    }

    const apiKey = this.getApiKey();

    const btn = this.els.settingsTestBtn;
    if (btn) {
      btn.disabled = true;
      btn.textContent = t("testing");
    }
    this.setStatus(statusEl, t("testing"), "pending");

    try {
      const result = await API.testAIConnection({
        apiKey,
        model: this.getSelectedModel(),
      });
      this.trackCost(result);
      this.setStatus(statusEl, result.message || t("connectionSuccess"), "success");
      if (window.CVEditor?.showToast) {
        window.CVEditor.showToast(t("openaiConnectionSuccess"));
      }
    } catch (error) {
      this.setStatus(statusEl, error.message, "error");
    } finally {
      if (btn) {
        btn.disabled = false;
        btn.textContent = t("testConnection");
      }
    }
  },

  setStatus(el, text, type) {
    if (!el) return;
    el.textContent = text;
    const inline = el.dataset.inlineStatus === "true"
      || el.classList.contains("ai-test-status--inline");
    el.className = inline ? "ai-test-status ai-test-status--inline" : "ai-test-status";
    if (type) el.classList.add(`ai-test-status--${type}`);
  },

  setPdfStatus(text, type = "") {
    const el = this.els.pdfStatus;
    if (!el) return;
    el.textContent = text || "";
    el.className = "ai-chat-pdf-status";
    if (type) el.classList.add(`ai-chat-pdf-status--${type}`);
  },

  formatPdfMethod(method) {
    return method === "tesseract" ? t("ocrMethod") : t("pdfTextLayer");
  },

  buildPdfContextMessage(fileName, result) {
    const methodLabel = this.formatPdfMethod(result.method);
    return [
      t("pdfUploaded", { name: fileName }),
      t("pdfMeta", { pages: result.pages, chars: result.chars, method: methodLabel }),
      "",
      t("pdfTextHeader"),
      result.text || "",
    ].join("\n");
  },

  appendPDFResult(fileName, result) {
    const methodLabel = this.formatPdfMethod(result.method);
    const summary = t("pdfExtractedSummary", {
      pages: result.pages,
      chars: result.chars,
      method: methodLabel,
    });
    const contextMessage = this.buildPdfContextMessage(fileName, result);
    this.messages.push({ role: "user", content: contextMessage });

    const wrapper = document.createElement("div");
    wrapper.className = "ai-chat-message ai-chat-message--user ai-chat-message--pdf";

    const title = document.createElement("strong");
    title.textContent = summary;
    wrapper.appendChild(title);

    const meta = document.createElement("div");
    meta.className = "ai-chat-pdf-meta";
    meta.textContent = t("fileLabel", { name: fileName });
    wrapper.appendChild(meta);

    const preview = document.createElement("pre");
    preview.className = "ai-chat-pdf-preview";
    const fullText = result.text || "";
    const truncated = fullText.length > PDF_PREVIEW_LIMIT;
    preview.textContent = truncated ? `${fullText.slice(0, PDF_PREVIEW_LIMIT)}…` : fullText;
    wrapper.appendChild(preview);

    const actions = document.createElement("div");
    actions.className = "ai-chat-pdf-actions";

    if (truncated) {
      const showBtn = document.createElement("button");
      showBtn.type = "button";
      showBtn.className = "btn btn-secondary btn-sm";
      showBtn.textContent = t("showText");
      showBtn.addEventListener("click", () => {
        const expanded = preview.classList.toggle("ai-chat-pdf-preview--expanded");
        showBtn.textContent = expanded ? t("hideText") : t("showText");
        preview.textContent = expanded ? fullText : `${fullText.slice(0, PDF_PREVIEW_LIMIT)}…`;
      });
      actions.appendChild(showBtn);
    }

    const applyBtn = document.createElement("button");
    applyBtn.type = "button";
    applyBtn.className = "btn btn-primary btn-sm";
    applyBtn.textContent = t("applyToCv");
    applyBtn.addEventListener("click", () => this.requestApplyPdfText());
    actions.appendChild(applyBtn);

    wrapper.appendChild(actions);
    this.els.messages?.appendChild(wrapper);
    this.scrollToBottom();
  },

  requestApplyPdfText() {
    const prompt = t("applyPdfPrompt");
    if (this.els.input) {
      this.els.input.value = prompt;
      this.els.input.focus();
    }
    if (this.isOpenAIAvailable()) {
      this.sendMessage();
      return;
    }
    if (window.CVEditor?.showToast) {
      window.CVEditor.showToast(t("textAddedToChat"));
    }
  },

  async uploadPDF(file) {
    if (!file) return;

    const name = (file.name || "").toLowerCase();
    if (!name.endsWith(".pdf") && file.type !== "application/pdf") {
      const message = t("onlyPdfAllowed");
      this.setPdfStatus(message, "error");
      if (window.CVEditor?.showToast) window.CVEditor.showToast(message);
      return;
    }

    const maxSize = 10 * 1024 * 1024;
    if (file.size > maxSize) {
      const message = t("pdfTooLarge");
      this.setPdfStatus(message, "error");
      if (window.CVEditor?.showToast) window.CVEditor.showToast(message);
      return;
    }

    if (this.uploadingPdf) return;
    this.uploadingPdf = true;
    if (this.els.pdfBtn) {
      this.els.pdfBtn.disabled = true;
      this.els.pdfBtn.textContent = t("processing");
    }
    this.setPdfStatus(t("pdfProcessing"), "pending");

    try {
      const result = await API.ocrPDF(file);
      const statusText = result.method === "tesseract"
        ? t("pdfPagesOcr", { pages: result.pages })
        : t("pdfPagesExtracted", { pages: result.pages });
      this.setPdfStatus(statusText, "success");
      this.appendPDFResult(file.name || "cv.pdf", result);
      if (window.CVEditor?.showToast) {
        window.CVEditor.showToast(t("pdfProcessed", { pages: result.pages, chars: result.chars }));
      }
    } catch (error) {
      this.setPdfStatus(error.message, "error");
      if (window.CVEditor?.showToast) {
        window.CVEditor.showToast(error.message);
      }
    } finally {
      this.uploadingPdf = false;
      if (this.els.pdfBtn) {
        this.els.pdfBtn.disabled = false;
        this.els.pdfBtn.textContent = t("addPdf");
      }
    }
  },

  appendMessage(role, text) {
    this.messages.push({ role, content: text });
    const bubble = document.createElement("div");
    bubble.className = `ai-chat-message ai-chat-message--${role}`;
    bubble.textContent = text;
    this.els.messages?.appendChild(bubble);
    this.scrollToBottom();
    return bubble;
  },

  addSystemMessage(text) {
    const bubble = document.createElement("div");
    bubble.className = "ai-chat-message ai-chat-message--system";
    bubble.textContent = text;
    this.els.messages?.appendChild(bubble);
    this.scrollToBottom();
    return bubble;
  },

  addAssistantMessage(text) {
    return this.appendMessage("assistant", text);
  },

  clearTranslateProgressTimer() {
    if (this.translateProgressTimer) {
      clearInterval(this.translateProgressTimer);
      this.translateProgressTimer = null;
    }
  },

  renderTranslateProgressLines() {
    if (!this.translateProgressContent) return;

    const lines = getTranslateProgressStages().slice(0, this.translateProgressStageIndex + 1);
    this.translateProgressContent.innerHTML = "";

    lines.forEach((line, index) => {
      const row = document.createElement("div");
      row.className = "ai-chat-progress-line";
      if (index < this.translateProgressStageIndex) {
        row.classList.add("ai-chat-progress-line--done");
      } else if (!this.translateProgressDone) {
        row.classList.add("ai-chat-progress-line--active");
      }
      row.textContent = index === 0 ? line : `→ ${line}`;
      this.translateProgressContent.appendChild(row);
    });

    this.scrollToBottom();
  },

  startTranslateProgress() {
    this.navigateToSection();
    this.clearTranslateProgressTimer();
    this.translateProgressDone = false;
    this.translateProgressStageIndex = 0;

    if (this.translateProgressEl) {
      this.translateProgressEl.remove();
    }

    const bubble = document.createElement("div");
    bubble.className = "ai-chat-message ai-chat-message--assistant ai-chat-progress";
    bubble.setAttribute("aria-live", "polite");
    bubble.setAttribute("aria-busy", "true");

    const content = document.createElement("div");
    content.className = "ai-chat-progress-content";
    bubble.appendChild(content);

    this.els.messages?.appendChild(bubble);
    this.translateProgressEl = bubble;
    this.translateProgressContent = content;
    this.renderTranslateProgressLines();

    this.translateProgressTimer = setInterval(() => {
      if (this.translateProgressDone) return;
      if (this.translateProgressStageIndex < getTranslateProgressStages().length - 1) {
        this.translateProgressStageIndex += 1;
        this.renderTranslateProgressLines();
      }
    }, TRANSLATE_PROGRESS_INTERVAL_MS);

    const advance = () => {
      if (this.translateProgressDone) return;
      if (this.translateProgressStageIndex < getTranslateProgressStages().length - 1) {
        this.translateProgressStageIndex += 1;
        this.renderTranslateProgressLines();
      }
    };

    return advance;
  },

  completeTranslateProgress(success, costUSD = 0, error = null, options = {}) {
    this.clearTranslateProgressTimer();
    this.translateProgressDone = true;

    const bubble = this.translateProgressEl;
    const content = this.translateProgressContent;
    if (!bubble || !content) return;

    bubble.setAttribute("aria-busy", "false");
    bubble.classList.remove("ai-chat-progress");
    content.innerHTML = "";

    let finalText;
    if (success) {
      const provider = options.provider || "google";
      const usage = options.usage;
      if (provider === "openai" && (costUSD > 0 || usage)) {
        const costLine = this.formatCostUSD(costUSD);
        const tokenLine = usage
          ? t("tokenSuffix", { prompt: usage.promptTokens || 0, completion: usage.completionTokens || 0 })
          : "";
        finalText = t("translateCompleteCost", { cost: costLine, tokens: tokenLine });
        this.trackCost({ costUSD });
      } else if (provider === "google") {
        finalText = t("translateCompleteGoogle");
        sessionStorage.setItem(AI_LAST_COST_KEY, "0");
        this.updateCostDisplay();
      } else {
        finalText = t("translateComplete");
        if (costUSD > 0) {
          finalText += ` ${t("costLabel", { cost: this.formatCostUSD(costUSD) })}`;
          this.trackCost({ costUSD });
        }
      }
      bubble.classList.add("ai-chat-progress--success");
    } else {
      finalText = t("translateFailed", { error: error || t("unknownError") });
      bubble.classList.add("ai-chat-progress--error");
    }

    const finalLine = document.createElement("div");
    finalLine.className = "ai-chat-progress-final";
    finalLine.textContent = finalText;
    content.appendChild(finalLine);

    this.messages.push({ role: "assistant", content: finalText });
    this.scrollToBottom();

    this.translateProgressEl = null;
    this.translateProgressContent = null;
    this.translateProgressStageIndex = 0;
  },

  appendProfilePicker(message, copyOptions = {}) {
    const wrapper = document.createElement("div");
    wrapper.className = "ai-chat-message ai-chat-message--assistant ai-profile-picker";

    const text = document.createElement("p");
    text.className = "ai-profile-picker-message";
    text.textContent = message || t("whichProfile");
    wrapper.appendChild(text);

    const actions = document.createElement("div");
    actions.className = "ai-profile-picker-actions";
    wrapper.appendChild(actions);

    this.els.messages?.appendChild(wrapper);
    this.scrollToBottom();

    API.listProfiles()
      .then((profiles) => {
        const currentId = window.CVEditor?.getCurrentProfileId?.() || null;
        const candidates = (profiles || []).filter((profile) => profile.id !== currentId);
        if (candidates.length === 0) {
          const empty = document.createElement("p");
          empty.className = "ai-profile-picker-empty";
          empty.textContent = t("noOtherProfiles");
          actions.appendChild(empty);
          return;
        }

        candidates.forEach((profile) => {
          const button = document.createElement("button");
          button.type = "button";
          button.className = "btn btn-secondary btn-sm ai-profile-picker-btn";
          button.textContent = profile.name;
          button.addEventListener("click", () => this.handleProfilePickerSelect(profile.id, copyOptions, wrapper));
          actions.appendChild(button);
        });
      })
      .catch((error) => {
        const errorEl = document.createElement("p");
        errorEl.className = "ai-profile-picker-empty";
        errorEl.textContent = error.message;
        actions.appendChild(errorEl);
      });
  },

  async handleProfilePickerSelect(sourceId, copyOptions, pickerEl) {
    const targetId = window.CVEditor?.getCurrentProfileId?.();
    if (!targetId) {
      if (window.CVEditor?.showToast) {
        window.CVEditor.showToast(t("selectProfileFirst"));
      }
      return;
    }

    pickerEl.querySelectorAll(".ai-profile-picker-btn").forEach((btn) => {
      btn.disabled = true;
    });

    try {
      const profile = await API.copyFromProfile(targetId, sourceId, copyOptions || {});
      if (window.CVEditor?.onProfileUpdated) {
        await window.CVEditor.onProfileUpdated(profile);
      }
      const successText = t("dataCopiedFromProfile");
      this.appendMessage("assistant", successText);
      if (window.CVEditor?.showToast) {
        window.CVEditor.showToast(t("profileDataCopied"));
      }
    } catch (error) {
      this.appendMessage("assistant", `${t("error")}: ${error.message}`);
      if (window.CVEditor?.showToast) {
        window.CVEditor.showToast(error.message);
      }
      pickerEl.querySelectorAll(".ai-profile-picker-btn").forEach((btn) => {
        btn.disabled = false;
      });
    }
  },

  appendActionStatus(action) {
    const label = this.formatActionLabel(action);
    if (!label) return;
    const bubble = document.createElement("div");
    bubble.className = "ai-chat-message ai-chat-message--action";
    bubble.textContent = label;
    this.els.messages?.appendChild(bubble);
    this.scrollToBottom();
  },

  formatActionLabel(action) {
    if (!action?.tool) return "";
    const names = {
      create_profile: t("actionCreateProfile"),
      duplicate_profile: t("actionDuplicateProfile"),
      update_profile_content: t("actionUpdateContent"),
      translate_profile: t("actionTranslateProfile"),
      switch_language: t("actionSwitchLanguage"),
      get_current_profile: t("actionGetProfile"),
      apply_to_profile: t("actionApplyToProfile"),
      request_profile_selection: t("actionRequestSelection"),
      copy_from_profile: t("actionCopyFromProfile"),
    };
    const prefix = action.status === "error" ? t("error") : t("action");
    const verb = names[action.tool] || action.tool;
    const detail = action.message ? `: ${action.message}` : "";
    return `${prefix}: ${verb}${detail}`;
  },

  scrollToBottom() {
    if (this.els.messages) {
      this.els.messages.scrollTop = this.els.messages.scrollHeight;
    }
  },

  showThinking() {
    this.removeThinking();
    const bubble = document.createElement("div");
    bubble.className = "ai-chat-message ai-chat-message--thinking ai-thinking";
    bubble.setAttribute("aria-live", "polite");
    bubble.innerHTML = `<span class="ai-thinking-label">${t("thinking")}<span class="typing-dots" aria-hidden="true"><span></span><span></span><span></span></span></span>`;
    this.els.messages?.appendChild(bubble);
    this.thinkingEl = bubble;
    this.scrollToBottom();
  },

  removeThinking() {
    if (this.thinkingEl) {
      this.thinkingEl.remove();
      this.thinkingEl = null;
    }
  },

  clearSuggestion() {
    this.pendingSuggestion = null;
    this.els.applyBtn?.classList.add("hidden");
  },

  showSuggestion(suggestion) {
    this.pendingSuggestion = suggestion;
    if (suggestion.action === "create_profile" || suggestion.action === "update_profile") {
      this.els.applyBtn?.classList.remove("hidden");
    } else {
      this.els.applyBtn?.classList.add("hidden");
    }
  },

  async sendMessage() {
    const text = this.els.input?.value.trim();
    if (!text) return;

    if (!this.isOpenAIAvailable()) {
      const message = t("openaiKeyRequired");
      console.error("[AIChat] sendMessage:", message);
      if (window.CVEditor?.showToast) {
        window.CVEditor.showToast(message);
      }
      return;
    }

    const apiKey = this.getApiKey();

    this.appendMessage("user", text);
    this.els.input.value = "";
    this.clearSuggestion();

    const chatMessages = this.messages.map((m) => ({ role: m.role, content: m.content }));
    const profileId = window.CVEditor?.getCurrentProfileId?.() || null;
    const language = window.CVEditor?.getEditingLanguage?.() || "tr";
    const translateProvider = this.els.translateProvider?.value || localStorage.getItem(TRANSLATE_PROVIDER_STORAGE_KEY) || "google";

    this.els.sendBtn.disabled = true;
    this.els.sendBtn.textContent = t("sending");
    this.els.input.disabled = true;
    this.showThinking();

    try {
      const response = await API.aiChat({
        apiKey,
        model: this.getSelectedModel(),
        messages: chatMessages,
        profileId,
        language,
        translateProvider,
      });

      this.trackCost(response);
      this.removeThinking();

      (response.actions || []).forEach((action) => this.appendActionStatus(action));

      const reply = response.message || t("noResponse");
      this.appendMessage("assistant", reply);

      if (response.showProfilePicker) {
        this.appendProfilePicker(response.profilePickerMessage, response.copyOptions || {});
      }

      if (response.profileUpdated && window.CVEditor?.onAIChatResult) {
        await window.CVEditor.onAIChatResult(response);
      }
    } catch (error) {
      this.removeThinking();
      this.appendMessage("assistant", `${t("error")}: ${error.message}`);
      if (window.CVEditor?.showToast) {
        window.CVEditor.showToast(error.message);
      }
    } finally {
      this.els.sendBtn.disabled = false;
      this.els.sendBtn.textContent = t("send");
      this.els.input.disabled = false;
      this.els.input?.focus();
    }
  },

  async applySuggestion() {
    if (!this.pendingSuggestion) return;

    const { action, profileName, language, content } = this.pendingSuggestion;
    if (!content) {
      if (window.CVEditor?.showToast) {
        window.CVEditor.showToast(t("noContentToApply"));
      }
      return;
    }

    this.els.applyBtn.disabled = true;
    this.els.applyBtn.textContent = t("applying");
    this.showThinking();
    try {
      const payload = {
        action,
        profileName,
        language,
        content,
      };
      if (action === "update_profile") {
        payload.profileId = window.CVEditor?.getCurrentProfileId?.() || null;
      }

      const profile = await API.applyAIChanges(payload);
      this.removeThinking();
      if (window.CVEditor?.onProfileUpdated) {
        await window.CVEditor.onProfileUpdated(profile);
      }
      this.appendMessage("assistant", action === "create_profile"
        ? t("profileCreated", { name: profile.name })
        : t("profileUpdated"));
      this.clearSuggestion();
      if (window.CVEditor?.showToast) {
        window.CVEditor.showToast(t("aiChangesApplied"));
      }
    } catch (error) {
      this.removeThinking();
      if (window.CVEditor?.showToast) {
        window.CVEditor.showToast(error.message);
      }
    } finally {
      this.els.applyBtn.disabled = false;
      this.els.applyBtn.textContent = t("apply");
    }
  },
};

function initAIChat() {
  AIChat.init();
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", initAIChat);
} else {
  initAIChat();
}

window.AIChat = AIChat;
