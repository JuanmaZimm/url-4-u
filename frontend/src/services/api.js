const API_BASE = '/api/v1';

export async function shortenUrl(originalUrl, customAlias = '') {
  const response = await fetch(`${API_BASE}/urls`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      original_url: originalUrl,
      custom_alias: customAlias.trim() || undefined,
    }),
  });

  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.error || 'Failed to shorten URL');
  }

  return data;
}

export async function getRecentUrls(limit = 10, offset = 0) {
  const response = await fetch(`${API_BASE}/urls?limit=${limit}&offset=${offset}`);
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.error || 'Failed to fetch recent URLs');
  }
  return data.data || [];
}

export async function getUrlAnalytics(shortCode) {
  const response = await fetch(`${API_BASE}/urls/${shortCode}/analytics`);
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.error || 'Failed to fetch analytics');
  }
  return data;
}

export function getQrCodeUrl(shortCode, size = 300) {
  return `${API_BASE}/urls/${shortCode}/qr?size=${size}`;
}
