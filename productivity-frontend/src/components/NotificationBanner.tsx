import React from "react";

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

const NotificationBanner = ({ onDismiss }: { onDismiss: () => void}) => {
    const handleViewTrends = () => {
        console.log('View Trends button clicked');
    }

    return (
        <div className="notification-banner mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
            <div className="notification-content">
                <BellIcon />
                <div className="notification-text">
                    <h4>New trends available!</h4>
                    <p>Your daily and weekly productivity trends are ready to view.</p>
                </div>
            </div>
            <div className="notification-actions">
                <button className="btn-view-trends" onClick={handleViewTrends}>
                    View Trends
                </button>
                <button className="btn-dismiss" onClick={onDismiss}>
                    Dismiss
                </button>
            </div>
        </div>
    );
};

export default NotificationBanner;