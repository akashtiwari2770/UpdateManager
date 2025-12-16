import React from 'react';
import { NotificationsList } from '@/components/notifications';
import { Button } from '@/components/ui';
import { useNavigate } from 'react-router-dom';

export const Notifications: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Notifications</h1>
        <Button onClick={() => navigate('/notifications/new')}>
          Create Notification
        </Button>
      </div>
      <NotificationsList showFilters={true} showMarkAllAsRead={true} />
    </div>
  );
};

