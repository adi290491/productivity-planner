import { forwardRef, useEffect, useImperativeHandle, useRef, useState } from "react";

export interface TimerHandle {
    getElapsedSeconds: () => number;
    resetTimer: () => void;
}

interface TimerProps {
    isRunning: boolean
    initialSeconds?: number;
}

const Timer = forwardRef<TimerHandle, TimerProps>(({ isRunning, initialSeconds = 0 }, ref) => {
    const [elapsedSeconds, setElapsedSeconds] = useState(initialSeconds);
    const startTimeRef = useRef<number | null>(null);
    const pausedTimeRef = useRef<number>(initialSeconds);

    useImperativeHandle(ref, () => ({
        getElapsedSeconds: () => elapsedSeconds,
        resetTimer: () => {
            setElapsedSeconds(0);
            startTimeRef.current = null;
            pausedTimeRef.current = 0;
        }
    }));

    useEffect(() => {
        if (isRunning) {
            startTimeRef.current = Date.now() - (pausedTimeRef.current * 1000);

            const intervalId = setInterval(() => {
                if (startTimeRef.current !== null) {
                    const now = Date.now();
                    const elapsed = Math.floor((now - startTimeRef.current) / 1000);
                    setElapsedSeconds(elapsed);
                }
            }, 100);

            return () => clearInterval(intervalId);
        } else {
            pausedTimeRef.current = elapsedSeconds;
            startTimeRef.current = null;
        }
    }, [isRunning]);

    const formatTime = (totalSeconds: number): string => {
    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const seconds = totalSeconds % 60;
    return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
  };

  return (
    <div className="timer-display">
      <div className="time-text">{formatTime(elapsedSeconds)}</div>
    </div>
  );
});

Timer.displayName = 'Timer';

export default Timer;