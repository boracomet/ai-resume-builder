const API = {
  downloadBlob(blob, filename) {
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = filename;
    anchor.click();
    URL.revokeObjectURL(url);
  },

  getFilenameFromResponse(response, fallback) {
    const disposition = response.headers.get("content-disposition") || "";
    const utfMatch = disposition.match(/filename\*=UTF-8''([^;]+)/i);
    if (utfMatch) {
      try {
        return decodeURIComponent(utfMatch[1].trim());
      } catch (_) {}
    }
    const quotedMatch = disposition.match(/filename="([^"]+)"/i);
    if (quotedMatch) return quotedMatch[1];
    const plainMatch = disposition.match(/filename=([^;]+)/i);
    if (plainMatch) return plainMatch[1].trim().replace(/^"|"$/g, "");
    return fallback;
  },

  sanitizeFilename(name) {
    return (name || "profile")
      .trim()
      .replace(/[^\w\s.-]/g, "")
      .replace(/\s+/g, "-")
      .replace(/-+/g, "-")
      .replace(/^-|-$/g, "") || "profile";
  },

  backupDateSuffix() {
    return new Date().toISOString().slice(0, 10);
  },
  async request(path, options = {}) {
    const { headers: optionHeaders, ...rest } = options;
    let response;
    try {
      response = await fetch(path, {
        ...rest,
        headers: {
          "Content-Type": "application/json",
          ...(optionHeaders || {}),
        },
      });
    } catch (error) {
      if (error?.name === "AbortError") throw error;
      throw error;
    }

    if (!response.ok) {
      let message = "İstek başarısız";
      try {
        const data = await response.json();
        message = data.error || message;
      } catch (_) {
        const text = (await response.text()).trim();
        if (text) message = text;
      }
      if (response.status === 404) {
        message = "AI API uç noktası bulunamadı. Sunucuyu yeniden başlatıp tekrar deneyin.";
      } else if (response.status === 502 && message === "İstek başarısız") {
        message = "OpenAI servisine ulaşılamadı";
      }
      throw new Error(message);
    }

    const contentType = response.headers.get("content-type") || "";
    if (contentType.includes("application/json")) {
      return response.json();
    }
    return response;
  },

  listProfiles() {
    return this.request("/api/profiles");
  },

  getProfile(id) {
    return this.request(`/api/profiles/${id}`);
  },

  createProfile(profile) {
    return this.request("/api/profiles", {
      method: "POST",
      body: JSON.stringify(profile),
    });
  },

  updateProfile(id, profile) {
    return this.request(`/api/profiles/${id}`, {
      method: "PUT",
      body: JSON.stringify(profile),
    });
  },

  deleteProfile(id) {
    return this.request(`/api/profiles/${id}`, { method: "DELETE" });
  },

  duplicateProfile(id) {
    return this.request(`/api/profiles/${id}/duplicate`, { method: "POST" });
  },

  copyFromProfile(targetId, sourceId, options = {}) {
    return this.request(`/api/profiles/${targetId}/copy-from/${sourceId}`, {
      method: "POST",
      body: JSON.stringify(options),
    });
  },

  async preview(profile, options = {}) {
    const response = await this.request("/api/preview", {
      method: "POST",
      body: JSON.stringify(profile),
      signal: options.signal,
    });
    if (typeof response?.text === "function") {
      return response.text();
    }
    return String(response ?? "");
  },

  async downloadPDF(profile) {
    const response = await fetch("/api/pdf", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(profile),
    });

    if (!response.ok) {
      let message = "PDF oluşturulamadı";
      try {
        const data = await response.json();
        message = data.error || message;
      } catch (_) {}
      throw new Error(message);
    }

    const blob = await response.blob();
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement("a");
    anchor.href = url;
    const stamp = new Date().toISOString().slice(0, 19).replace(/[:T]/g, "-");
    const baseName = (profile?.personal?.name || "resume")
      .toString()
      .trim()
      .toLowerCase()
      .replace(/\s+/g, "-")
      .replace(/[^a-z0-9-_]/gi, "") || "resume";
    anchor.download = `${baseName}-${stamp}.pdf`;
    anchor.click();
    URL.revokeObjectURL(url);
  },

  async uploadPhoto(id, file) {
    const formData = new FormData();
    formData.append("photo", file);

    const response = await fetch(`/api/profiles/${id}/photo`, {
      method: "POST",
      body: formData,
    });

    if (!response.ok) {
      const data = await response.json();
      throw new Error(data.error || "Fotoğraf yüklenemedi");
    }

    return response.json();
  },

  translateProfile(id, payload) {
    return this.request(`/api/profiles/${id}/translate`, {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },

  getSettings() {
    return this.request("/api/settings");
  },

  listAIModels(apiKey) {
    const trimmed = (apiKey || "").trim();
    const params = trimmed ? `?apiKey=${encodeURIComponent(trimmed)}` : "";
    return this.request(`/api/ai/models${params}`);
  },

  testAIConnection(payload) {
    return this.request("/api/ai/test", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },

  aiChat(payload) {
    return this.request("/api/ai/chat", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },

  applyAIChanges(payload) {
    return this.request("/api/ai/apply", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },

  async ocrPDF(file) {
    const formData = new FormData();
    formData.append("file", file);

    const response = await fetch("/api/ocr/pdf", {
      method: "POST",
      body: formData,
    });

    if (!response.ok) {
      let message = "PDF işlenemedi";
      try {
        const data = await response.json();
        message = data.error || message;
      } catch (_) {
        const text = (await response.text()).trim();
        if (text) message = text;
      }
      throw new Error(message);
    }

    return response.json();
  },

  async exportCurrentProfile(id, profileName) {
    const response = await fetch(`/api/profiles/${id}/export`);

    if (!response.ok) {
      let message = "Profil dışa aktarılamadı";
      try {
        const data = await response.json();
        message = data.error || message;
      } catch (_) {}
      throw new Error(message);
    }

    const blob = await response.blob();
    const fallback = `${this.sanitizeFilename(profileName)}.json`;
    const filename = this.getFilenameFromResponse(response, fallback);
    this.downloadBlob(blob, filename);
  },

  async exportAllProfiles() {
    const response = await fetch("/api/export");

    if (!response.ok) {
      let message = "Yedek oluşturulamadı";
      try {
        const data = await response.json();
        message = data.error || message;
      } catch (_) {}
      throw new Error(message);
    }

    const blob = await response.blob();
    const fallback = `ai-resume-builder-backup-${this.backupDateSuffix()}.json`;
    const filename = this.getFilenameFromResponse(response, fallback);
    this.downloadBlob(blob, filename);
  },

  async importProfiles(file, mode = "merge") {
    const text = await file.text();
    let payload;
    try {
      payload = JSON.parse(text);
    } catch (_) {
      throw new Error("Geçersiz yedek dosyası");
    }

    const profiles = Array.isArray(payload?.profiles)
      ? payload.profiles
      : Array.isArray(payload)
        ? payload
        : null;

    if (!profiles || profiles.length === 0) {
      throw new Error("Geçersiz yedek dosyası");
    }

    return this.request("/api/import", {
      method: "POST",
      body: JSON.stringify({ mode, profiles }),
    });
  },

  resetAllProfiles() {
    return this.request("/api/reset", { method: "POST" });
  },
};
