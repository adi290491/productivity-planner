export const parseTimeToMinutes = (timeStr: string): number => {
    const match = timeStr.match(/(?:(\d+)h)?(?:(\d+)m)?/);
    if (!match) return 0;
    const hours = parseInt(match[1] || '0');
    const mins = parseInt(match[2] || '0');
    return hours * 60 + mins;
  };
  
  export const formatMinutesToLabel = (mins: number): string => {
    const h = Math.floor(mins / 60);
    const m = mins % 60;
    return `${h}h ${m}m`;
  };
  
  export const formatDuration = (start: string, end: string): string => {

    if (!start || !end) return "-";
    
    const s = new Date(start);
    const e = new Date(end);
    const diff = Math.floor((e.getTime() - s.getTime()) / 1000);
    const h = Math.floor(diff / 3600);
    const m = Math.floor((diff % 3600) / 60);
    const sRem = diff % 60;
    return `${h}h ${m}m ${sRem}s`;
  };
  
  export const formatTime = (totalSeconds: number) => {
    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = Math.floor(totalSeconds % 60);

    return `${hours}H ${minutes}M ${seconds}S`
};