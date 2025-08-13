/**
 * Formats a byte value into a human-readable string with appropriate units
 * @param bytes - The number of bytes to format
 * @param decimals - Number of decimal places to show (default: 2)
 * @returns Formatted string with appropriate unit (B, KB, MB, GB, TB, PB)
 */
export function formatBytes(bytes: number, decimals: number = 2): string {
  if (bytes === 0) return '0 B';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

/**
 * Formats a byte value into a human-readable string using binary units (1024-based)
 * @param bytes - The number of bytes to format
 * @param decimals - Number of decimal places to show (default: 2)
 * @returns Formatted string with appropriate binary unit (B, KiB, MiB, GiB, TiB, PiB)
 */
export function formatBinaryBytes(bytes: number, decimals: number = 2): string {
  if (bytes === 0) return '0 B';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['B', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

/**
 * Formats a byte value into a human-readable string using decimal units (1000-based)
 * @param bytes - The number of bytes to format
 * @param decimals - Number of decimal places to show (default: 2)
 * @returns Formatted string with appropriate decimal unit (B, kB, MB, GB, TB, PB)
 */
export function formatDecimalBytes(bytes: number, decimals: number = 2): string {
  if (bytes === 0) return '0 B';

  const k = 1000;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['B', 'kB', 'MB', 'GB', 'TB', 'PB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

/**
 * Formats a Unix timestamp into a human-readable "time ago" string
 * @param timestamp - Unix timestamp (seconds since epoch)
 * @returns Formatted string like "3 days ago", "2 hours ago", etc.
 */
export function formatTimeAgo(timestamp: number): string {
  if (!timestamp) return 'Unknown';

  const now = Math.floor(Date.now() / 1000); // Current time in seconds
  const diff = now - timestamp;

  // Handle future dates
  if (diff < 0) return 'in the future';

  // Define time intervals in seconds
  const intervals = [
    { label: 'year', seconds: 31536000 },
    { label: 'month', seconds: 2592000 },
    { label: 'week', seconds: 604800 },
    { label: 'day', seconds: 86400 },
    { label: 'hour', seconds: 3600 },
    { label: 'minute', seconds: 60 },
    { label: 'second', seconds: 1 }
  ];

  for (const interval of intervals) {
    const count = Math.floor(diff / interval.seconds);
    if (count >= 1) {
      return `${count} ${interval.label}${count !== 1 ? 's' : ''} ago`;
    }
  }

  return 'just now';
}

/**
 * Formats a container name by removing the preceding slash if present
 * @param name - Container name that may start with "/"
 * @returns Container name without the preceding slash
 */
export function formatContainerName(name: string): string {
  if (!name) return '';

  // Remove leading slash if present
  return name.startsWith('/') ? name.substring(1) : name;
}
