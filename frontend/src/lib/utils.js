// src/lib/utils.js

/**
 * Given an ISO timestamp, returns a human‑friendly time string.
 * @param {string} isoString
 * @returns {string}
 */
export function formatMessageTime(isoString) {
  const d = new Date(isoString);
  // e.g. "3:45 PM"
  return d.toLocaleTimeString([], {
    hour: 'numeric',
    minute: '2-digit',
  });
}
