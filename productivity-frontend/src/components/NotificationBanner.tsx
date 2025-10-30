import { useState } from "react";
import type { UnviewedTrendsCount } from "../types/trend";
import { markTrendsAsViewed } from "../api/trend-analysis";

const BellIcon = () => (
    <svg 
    xmlns="http://www.w3.org/2000/svg" 
    width="24" 
    height="24" 
    viewBox="0 0 24 24" 
    fill="none" 
    stroke="currentColor" 
    strokeWidth="2" 
    strokeLinecap="round" 
    strokeLinejoin="round"
    className="notification-bell-icon"
  >
    <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
    <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
  </svg>
);

interface NotificationBannerProps {
    onDismiss: () => void;
    trendsCount: UnviewedTrendsCount | null;
    onViewTrends: () => void;
}

const NotificationBanner = ({ onDismiss, trendsCount, onViewTrends }: NotificationBannerProps) => {
    const [isLoading, setIsLoading] = useState(false);

    const handleViewTrends = async () => {
        console.log('View Trends button clicked');

        const token = localStorage.getItem('token');
        if (!token || !trendsCount) {
            console.error('No token or trends count available');
            return;
        }

        setIsLoading(true);

        try {
            const promises: Promise<void>[] = [];

            if (trendsCount.daily_count > 0) {
                promises.push(markTrendsAsViewed(token, 'daily'));
            }

            if (trendsCount.weekly_count > 0) {
                promises.push(markTrendsAsViewed(token, 'weekly'));
            }

            await Promise.all(promises);

            onViewTrends();

            onDismiss();
        } catch (error) {
            console.error('Failed to mark trends as viewed:', error);

            alert('Failed to update trends. Please try again.');
        } finally {
            setIsLoading(false);
        }
    }


    return (
        <div className="notification-banner mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
            <div className="notification-content">
                <BellIcon />
                <div className="notification-text">
                    <h4>New trends available!</h4>
                    <p>
                        {trendsCount?.daily_count && trendsCount?.weekly_count 
                            ? 'Your daily and weekly productivity trends are ready to view.'
                            : trendsCount?.daily_count 
                            ? 'Your daily productivity trends are ready to view.'
                            : 'Your weekly productivity trends are ready to view.'}
                    </p>
                </div>
            </div>
            <div className="notification-actions">
                <button 
                    className="btn-view-trends" 
                    onClick={handleViewTrends}
                    disabled={isLoading}
                >
                    {isLoading ? 'Loading...' : 'View Trends'}
                </button>
                <button className="btn-dismiss" onClick={onDismiss}>
                    Dismiss
                </button>
            </div>
        </div>
    );
};

export default NotificationBanner;