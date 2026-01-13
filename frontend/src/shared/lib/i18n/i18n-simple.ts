/**
 * Simplified i18n Configuration for Vietnamese error messages
 * (Simplified version without full i18next setup for MVP)
 */

export const translations = {
  vi: {
    errors: {
      network_error: "Lỗi kết nối mạng. Vui lòng kiểm tra kết nối internet của bạn.",
      unknown_error: "Đã xảy ra lỗi không xác định. Vui lòng thử lại sau.",
      validation_error: "Dữ liệu không hợp lệ. Vui lòng kiểm tra lại.",
      not_found: "Không tìm thấy tài nguyên yêu cầu.",
      server_error: "Lỗi máy chủ. Vui lòng thử lại sau.",
      insufficient_words: "Không đủ từ vựng cho cấu hình bạn đã chọn. Vui lòng chọn cấu hình khác.",
      same_language: "Ngôn ngữ nguồn và ngôn ngữ đích phải khác nhau.",
    },
    common: {
      loading: "Đang tải...",
      error: "Lỗi",
      success: "Thành công",
    },
  },
  en: {
    errors: {
      network_error: "Network connection error. Please check your internet connection.",
      unknown_error: "An unknown error occurred. Please try again later.",
      validation_error: "Invalid data. Please check again.",
      not_found: "Requested resource not found.",
      server_error: "Server error. Please try again later.",
      insufficient_words: "Insufficient vocabulary for the selected configuration. Please choose a different configuration.",
      same_language: "Source and target languages must be different.",
    },
    common: {
      loading: "Loading...",
      error: "Error",
      success: "Success",
    },
  },
};

export type Language = 'vi' | 'en';

let currentLanguage: Language = 'vi';

export const setLanguage = (lang: Language) => {
  currentLanguage = lang;
};

export const getLanguage = (): Language => {
  return currentLanguage;
};

export const t = (key: string, lang?: Language): string => {
  const lng = lang || currentLanguage;
  const keys = key.split('.');
  let value: unknown = translations[lng];
  
  for (const k of keys) {
    if (value && typeof value === 'object' && k in value) {
      value = (value as Record<string, unknown>)[k];
    } else {
      return key; // Return key if translation not found
    }
  }
  
  return typeof value === 'string' ? value : key;
};

