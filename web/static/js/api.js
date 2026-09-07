const API = {
  async request(path, options = {}) {
    const response = await fetch(path, {
      headers: {
        "Content-Type": "application/json",
        ...(options.headers || {}),
      },
      ...options,
    });

    if (!response.ok) {
      let message = "İstek başarısız";
      try {
        const data = await response.json();
        message = data.error || message;
      } catch (_) {
        message = await response.text();
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

  async preview(profile) {
    const response = await this.request("/api/preview", {
      method: "POST",
      body: JSON.stringify(profile),
    });
    return response.text();
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
    anchor.download = "bora-cv.pdf";
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
};
