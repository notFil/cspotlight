export function timeAgo(dateInput: string | Date): string {
  const date = new Date(dateInput);
  const now = new Date();
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  const units = [
    { label: 'year', seconds: 31536000 },
    { label: 'month', seconds: 2592000 },
    { label: 'week', seconds: 604800 },
    { label: 'day', seconds: 86400 },
    { label: 'hour', seconds: 3600 },
    { label: 'minute', seconds: 60 },
    { label: 'second', seconds: 1 }
  ];

  for (const unit of units) {
    const interval = Math.floor(seconds / unit.seconds);
    if (interval >= 1) {
      return `Last ${interval} ${unit.label}${interval > 1 ? 's' : ''} ago`;
    }
  }
  return 'Never';
}

export function formatTimestamp(timestamp: string): string {
  const date = new Date(timestamp);
  const day = date.getUTCDate().toString().padStart(2, '0');
  const monthNames = ["Jan", "Feb", "Mar", "Apr", "May", "Jun",
    "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
  const month = monthNames[date.getUTCMonth()];
  const year = date.getUTCFullYear().toString().slice(2);
  const hours = date.getUTCHours().toString().padStart(2, '0');
  const minutes = date.getUTCMinutes().toString().padStart(2, '0');

  return `${day} ${month} ${year} ${hours}:${minutes} UTC`;
}
