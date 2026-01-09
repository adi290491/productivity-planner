/**
 * Format utility functions
 * Updated to handle dates and times correctly across timezones
 */

/**
 * Parse time string like "3h 20m" into total minutes
 */
export const parseTimeToMinutes = (timeStr: string): number => {
    const match = timeStr.match(/(?:(\d+)h)?(?:(\d+)m)?/);
    if (!match) return 0;
    const hours = parseInt(match[1] || '0');
    const mins = parseInt(match[2] || '0');
    return hours * 60 + mins;
  };
  
/**
 * Format minutes into human readable label like "3h 20m"
 */
export const formatMinutesToLabel = (mins: number): string => {
    const h = Math.floor(mins / 60);
    const m = mins % 60;
    return `${h}h ${m}m`;
  };
  
/**
 * Calculate duration between two ISO timestamp strings
 * Returns formatted string like "2h 15m 30s"
 */
export const formatDuration = (start: string, end: string): string => {
    if (!start || !end) return "-";
    
    try {
      const s = new Date(start);
      const e = new Date(end);
      
      // Validate dates
      if (isNaN(s.getTime()) || isNaN(e.getTime())) {
        console.error('Invalid date in formatDuration:', { start, end });
        return "-";
      }
      
      const diff = Math.floor((e.getTime() - s.getTime()) / 1000);
      const h = Math.floor(diff / 3600);
      const m = Math.floor((diff % 3600) / 60);
      const sRem = diff % 60;
      return `${h}h ${m}m ${sRem}s`;
    } catch (error) {
      console.error('Error formatting duration:', error);
      return "-";
    }
  };
  
/**
 * Format total seconds into "XH YM ZS" display
 */
export const formatTime = (totalSeconds: number): string => {
    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const seconds = Math.floor(totalSeconds % 60);

    return `${hours}H ${minutes}M ${seconds}S`;
};

/**
 * Format date string (YYYY-MM-DD) for display
 * Works correctly in all timezones
 */
export const formatDate = (dateString: string): string => {
  if (!dateString || typeof dateString !== 'string') {
    console.error('Invalid date string:', dateString);
    return 'Invalid Date';
  }

  try {
    // Parse as local date (not UTC)
    const [year, month, day] = dateString.split('-').map(Number);
    
    if (!year || !month || !day) {
      throw new Error('Invalid date format');
    }
    
    // Create date in local timezone (month is 0-indexed)
    const date = new Date(year, month - 1, day);
    
    if (isNaN(date.getTime())) {
      throw new Error('Invalid date');
    }
    
    return date.toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric' 
    });
  } catch (error) {
    console.error('Error formatting date:', dateString, error);
    return 'Invalid Date';
  }
};

/**
 * Format ISO timestamp for display
 * Shows local time in readable format
 */
export const formatTimestamp = (isoString: string): string => {
  if (!isoString) return 'Invalid Date';
  
  try {
    const date = new Date(isoString);
    
    if (isNaN(date.getTime())) {
      throw new Error('Invalid timestamp');
    }
    
    return date.toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: 'numeric',
      minute: '2-digit',
      hour12: true
    });
  } catch (error) {
    console.error('Error formatting timestamp:', isoString, error);
    return 'Invalid Date';
  }
};

/**
 * Get current date in YYYY-MM-DD format (local timezone)
 */
export const getCurrentDate = (): string => {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, '0');
  const day = String(now.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

/**
 * Validate if a string is a valid date in YYYY-MM-DD format
 */
export const isValidDate = (dateString: string): boolean => {
  if (!dateString || typeof dateString !== 'string') {
    return false;
  }
  
  const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
  if (!dateRegex.test(dateString)) {
    return false;
  }
  
  try {
    const [year, month, day] = dateString.split('-').map(Number);
    const date = new Date(year, month - 1, day);
    return !isNaN(date.getTime());
  } catch {
    return false;
  }
};

/**
 * Convert seconds to formatted time string "Xh Ym Zs"
 */
export const formatSeconds = (seconds: number): string => {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;
  
  if (h > 0) {
    return `${h}h ${m}m ${s}s`;
  } else if (m > 0) {
    return `${m}m ${s}s`;
  } else {
    return `${s}s`;
  }
};