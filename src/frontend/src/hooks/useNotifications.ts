import { useEffect, useRef } from 'react';
import { notificationsApi } from '@/services/api/notifications';
import { useAppStore } from '@/store/useAppStore';
import { Notification } from '@/types';

interface UseNotificationsOptions {
  pollInterval?: number; // in milliseconds
  onNewNotification?: (notification: Notification) => void;
  enabled?: boolean;
}

export const useNotifications = (options: UseNotificationsOptions = {}) => {
  const {
    pollInterval = 30000, // 30 seconds default
    onNewNotification,
    enabled = true,
  } = options;

  const { unreadNotificationCount, setUnreadNotificationCount } = useAppStore();
  const previousCountRef = useRef<number>(0);
  const lastNotificationIdRef = useRef<string | null>(null);

  useEffect(() => {
    if (!enabled) return;

    // Load initial count
    const loadUnreadCount = async () => {
      try {
        const count = await notificationsApi.getUnreadCount();
        setUnreadNotificationCount(count);
        previousCountRef.current = count;
      } catch (error) {
        console.error('Failed to load unread notification count:', error);
      }
    };

    loadUnreadCount();

    // Poll for new notifications
    const interval = setInterval(async () => {
      try {
        const count = await notificationsApi.getUnreadCount();
        setUnreadNotificationCount(count);

        // Check if there are new notifications
        if (count > previousCountRef.current && onNewNotification) {
          // Fetch the latest notification
          const response = await notificationsApi.getAll({
            is_read: false,
            limit: 1,
            page: 1,
          });

          if (response.data.length > 0) {
            const latestNotification = response.data[0];
            // Only trigger if it's a new notification (different ID)
            if (latestNotification.id !== lastNotificationIdRef.current) {
              lastNotificationIdRef.current = latestNotification.id;
              onNewNotification(latestNotification);
            }
          }
        }

        previousCountRef.current = count;
      } catch (error) {
        console.error('Failed to poll notifications:', error);
      }
    }, pollInterval);

    return () => clearInterval(interval);
  }, [pollInterval, onNewNotification, enabled, setUnreadNotificationCount]);

  return {
    unreadCount: unreadNotificationCount,
    refreshCount: async () => {
      try {
        const count = await notificationsApi.getUnreadCount();
        setUnreadNotificationCount(count);
        previousCountRef.current = count;
      } catch (error) {
        console.error('Failed to refresh notification count:', error);
      }
    },
  };
};

