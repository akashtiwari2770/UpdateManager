import React from 'react';
import { useNavigate } from 'react-router-dom';
import { AuditLog, AuditAction } from '@/types';
import { Card, Badge, Spinner } from '@/components/ui';
import { ActionBadge } from '@/components/audit-logs/ActionBadge';

interface ActivityTimelineProps {
  activities: AuditLog[];
  loading?: boolean;
  limit?: number;
}

export const ActivityTimeline: React.FC<ActivityTimelineProps> = ({
  activities,
  loading,
  limit = 20,
}) => {
  const navigate = useNavigate();

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  };

  const getUserInitials = (email: string) => {
    return email
      .split('@')[0]
      .split('.')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .substring(0, 2);
  };

  const handleActivityClick = (activity: AuditLog) => {
    switch (activity.resource_type.toLowerCase()) {
      case 'product':
        navigate(`/products/${activity.resource_id}`);
        break;
      case 'version':
        navigate(`/versions/${activity.resource_id}`);
        break;
      default:
        break;
    }
  };

  return (
    <Card title="Recent Activity">
      {loading ? (
        <div className="flex items-center justify-center py-8">
          <Spinner />
        </div>
      ) : activities.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <p>No recent activity</p>
        </div>
      ) : (
        <div className="flow-root">
          <ul className="-mb-8">
            {activities.slice(0, limit).map((activity, index) => (
              <li key={activity.id}>
                <div className="relative pb-8">
                  {index !== activities.length - 1 && (
                    <span
                      className="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200"
                      aria-hidden="true"
                    />
                  )}
                  <div className="relative flex space-x-3">
                    <div>
                      <div className="flex h-8 w-8 items-center justify-center rounded-full bg-blue-100 ring-8 ring-white">
                        <span className="text-xs font-medium text-blue-800">
                          {getUserInitials(activity.user_email)}
                        </span>
                      </div>
                    </div>
                    <div className="flex min-w-0 flex-1 justify-between space-x-4 pt-1.5">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-1">
                          <ActionBadge action={activity.action} />
                          <span className="text-sm text-gray-500">
                            {activity.user_email}
                          </span>
                        </div>
                        <p className="text-sm text-gray-900">
                          <button
                            onClick={() => handleActivityClick(activity)}
                            className="font-medium hover:text-blue-600"
                          >
                            {activity.resource_type} {activity.resource_id}
                          </button>
                        </p>
                      </div>
                      <div className="whitespace-nowrap text-right text-sm text-gray-500">
                        {formatTimestamp(activity.timestamp)}
                      </div>
                    </div>
                  </div>
                </div>
              </li>
            ))}
          </ul>
          {activities.length >= limit && (
            <div className="mt-4 text-center">
              <button
                onClick={() => navigate('/audit-logs')}
                className="text-sm text-blue-600 hover:text-blue-700 font-medium"
              >
                View All Activity →
              </button>
            </div>
          )}
        </div>
      )}
    </Card>
  );
};

