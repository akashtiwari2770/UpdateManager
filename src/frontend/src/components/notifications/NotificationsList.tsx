import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { notificationsApi } from '@/services/api/notifications';
import {
  Notification,
  NotificationType,
  NotificationPriority,
  ListNotificationsQuery,
} from '@/types';
import { useAppStore } from '@/store/useAppStore';
import { Button, Badge, Spinner, Alert, Select } from '@/components/ui';
import { NotificationIcon } from './NotificationIcon';

interface NotificationsListProps {
  showFilters?: boolean;
  showMarkAllAsRead?: boolean;
}

export const NotificationsList: React.FC<NotificationsListProps> = ({
  showFilters = true,
  showMarkAllAsRead = true,
}) => {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [filters, setFilters] = useState<ListNotificationsQuery>({
    page: 1,
    limit: 20,
  });
  const navigate = useNavigate();
  const { setUnreadNotificationCount } = useAppStore();

  useEffect(() => {
    loadNotifications();
  }, [filters]);

  const loadNotifications = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await notificationsApi.getAll(filters);
      setNotifications(response.data);
      setTotalPages(response.pagination?.total_pages || 1);
      setPage(response.pagination?.current_page || 1);
    } catch (err) {
      setError('Failed to load notifications');
      console.error('Failed to load notifications:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleMarkAsRead = async (id: string) => {
    try {
      await notificationsApi.markAsRead(id);
      setNotifications((prev) =>
        prev.map((n) => (n.id === id ? { ...n, is_read: true } : n))
      );
      // Update unread count
      const count = await notificationsApi.getUnreadCount();
      setUnreadNotificationCount(count);
    } catch (err) {
      console.error('Failed to mark notification as read:', err);
    }
  };

  const handleMarkAllAsRead = async () => {
    try {
      await notificationsApi.markAllAsRead();
      setNotifications((prev) =>
        prev.map((n) => ({ ...n, is_read: true }))
      );
      setUnreadNotificationCount(0);
    } catch (err) {
      console.error('Failed to mark all as read:', err);
    }
  };

  const handleNotificationClick = (notification: Notification) => {
    if (!notification.is_read) {
      handleMarkAsRead(notification.id);
    }

    // Navigate based on notification type
    if (notification.version_id) {
      navigate(`/versions/${notification.version_id}`);
    } else if (notification.product_id) {
      navigate(`/products/${notification.product_id}`);
    }
  };

  const handleFilterChange = (key: keyof ListNotificationsQuery, value: any) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value,
      page: 1, // Reset to first page on filter change
    }));
  };

  const handlePageChange = (newPage: number) => {
    setFilters((prev) => ({ ...prev, page: newPage }));
  };

  const getPriorityBadgeColor = (priority: NotificationPriority) => {
    switch (priority) {
      case NotificationPriority.CRITICAL:
        return 'bg-red-100 text-red-800';
      case NotificationPriority.HIGH:
        return 'bg-orange-100 text-orange-800';
      case NotificationPriority.NORMAL:
        return 'bg-blue-100 text-blue-800';
      case NotificationPriority.LOW:
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="space-y-4">
      {/* Filters and Actions */}
      {(showFilters || showMarkAllAsRead) && (
        <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
          {showFilters && (
            <div className="flex flex-wrap gap-4 items-center">
              <Select
                value={filters.is_read === undefined ? 'all' : filters.is_read ? 'read' : 'unread'}
                onChange={(e) => {
                  const value = e.target.value;
                  handleFilterChange(
                    'is_read',
                    value === 'all' ? undefined : value === 'read'
                  );
                }}
                className="min-w-[150px]"
              >
                <option value="all">All Status</option>
                <option value="unread">Unread</option>
                <option value="read">Read</option>
              </Select>

              <Select
                value={filters.type || 'all'}
                onChange={(e) => {
                  handleFilterChange(
                    'type',
                    e.target.value === 'all' ? undefined : (e.target.value as NotificationType)
                  );
                }}
                className="min-w-[150px]"
              >
                <option value="all">All Types</option>
                <option value={NotificationType.UPDATE_AVAILABLE}>Update Available</option>
                <option value={NotificationType.NEW_VERSION}>New Version</option>
                <option value={NotificationType.SECURITY_RELEASE}>Security Release</option>
                <option value={NotificationType.EOL_WARNING}>EOL Warning</option>
              </Select>

              <Select
                value={filters.priority || 'all'}
                onChange={(e) => {
                  handleFilterChange(
                    'priority',
                    e.target.value === 'all' ? undefined : (e.target.value as NotificationPriority)
                  );
                }}
                className="min-w-[150px]"
              >
                <option value="all">All Priorities</option>
                <option value={NotificationPriority.CRITICAL}>Critical</option>
                <option value={NotificationPriority.HIGH}>High</option>
                <option value={NotificationPriority.NORMAL}>Normal</option>
                <option value={NotificationPriority.LOW}>Low</option>
              </Select>
            </div>
          )}

          {showMarkAllAsRead && (
            <Button
              variant="secondary"
              size="sm"
              onClick={handleMarkAllAsRead}
              disabled={notifications.filter((n) => !n.is_read).length === 0}
            >
              Mark All as Read
            </Button>
          )}
        </div>
      )}

      {error && <Alert variant="error">{error}</Alert>}

      {/* Notifications List */}
      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Spinner />
        </div>
      ) : notifications.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          <p>No notifications found</p>
        </div>
      ) : (
        <>
          <div className="bg-white rounded-lg border border-gray-200 divide-y divide-gray-200">
            {notifications.map((notification) => (
              <div
                key={notification.id}
                onClick={() => handleNotificationClick(notification)}
                className={`p-4 hover:bg-gray-50 cursor-pointer transition-colors ${
                  !notification.is_read ? 'bg-blue-50' : ''
                }`}
              >
                <div className="flex items-start gap-4">
                  <div className="flex-shrink-0 mt-1">
                    <NotificationIcon type={notification.type} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-4">
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <h3
                            className={`text-base font-medium ${
                              !notification.is_read
                                ? 'text-gray-900'
                                : 'text-gray-700'
                            }`}
                          >
                            {notification.title}
                          </h3>
                          {!notification.is_read && (
                            <div className="w-2 h-2 bg-blue-600 rounded-full" />
                          )}
                        </div>
                        <p className="text-sm text-gray-600 mt-1">
                          {notification.message}
                        </p>
                        <div className="flex items-center gap-3 mt-2">
                          <Badge className={getPriorityBadgeColor(notification.priority)}>
                            {notification.priority}
                          </Badge>
                          <span className="text-xs text-gray-400">
                            {new Date(notification.created_at).toLocaleString()}
                          </span>
                        </div>
                      </div>
                      {!notification.is_read && (
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleMarkAsRead(notification.id);
                          }}
                        >
                          Mark as Read
                        </Button>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between">
              <div className="text-sm text-gray-600">
                Page {page} of {totalPages}
              </div>
              <div className="flex gap-2">
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => handlePageChange(page - 1)}
                  disabled={page === 1}
                >
                  Previous
                </Button>
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => handlePageChange(page + 1)}
                  disabled={page === totalPages}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
};

