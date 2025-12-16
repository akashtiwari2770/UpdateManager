import React, { useState, useContext } from 'react';
import { Link } from 'react-router-dom';
import { NotificationCenter } from '@/components/notifications';
import { useNotifications } from '@/hooks/useNotifications';
import { Notification } from '@/types';
import { ToastContext } from '@/components/notifications/ToastContainer';

export const Header: React.FC = () => {
  const [notificationCenterOpen, setNotificationCenterOpen] = useState(false);
  const toastContext = useContext(ToastContext);
  const showToast = toastContext?.showToast;

  const handleNewNotification = (notification: Notification) => {
    if (showToast) {
      showToast(notification.message, 'info');
    }
  };

  useNotifications({
    pollInterval: 30000, // Poll every 30 seconds
    onNewNotification: handleNewNotification,
    enabled: true,
  });

  return (
    <header className="bg-white border-b border-gray-200 shadow-sm relative">
      <div className="px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-6">
          <Link to="/" className="flex items-center gap-2">
            <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">UM</span>
            </div>
            <span className="text-xl font-semibold text-gray-900">Update Manager</span>
          </Link>
        </div>

        <div className="flex items-center gap-4">
          {/* Search bar placeholder */}
          <div className="hidden md:block">
            <input
              type="text"
              placeholder="Search..."
              className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* Notification bell */}
          <div className="relative">
            <button
              onClick={() => setNotificationCenterOpen(!notificationCenterOpen)}
              className="relative p-2 text-gray-600 hover:text-gray-900 transition-colors"
              aria-label="Notifications"
            >
              <svg
                className="w-6 h-6"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
                />
              </svg>
              <NotificationBadge />
            </button>
            <NotificationCenter
              isOpen={notificationCenterOpen}
              onClose={() => setNotificationCenterOpen(false)}
            />
          </div>

          {/* User menu placeholder */}
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-gray-300 rounded-full"></div>
            <span className="hidden md:block text-sm font-medium text-gray-700">User</span>
          </div>
        </div>
      </div>
    </header>
  );
};

const NotificationBadge: React.FC = () => {
  const { unreadNotificationCount } = useNotifications();

  if (unreadNotificationCount === 0) return null;

  return (
    <span className="absolute top-0 right-0 block h-5 w-5 rounded-full bg-red-500 text-white text-xs font-bold flex items-center justify-center">
      {unreadNotificationCount > 99 ? '99+' : unreadNotificationCount}
    </span>
  );
};

