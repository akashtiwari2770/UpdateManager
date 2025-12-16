import React from 'react';
import { Badge } from '@/components/ui';
import { SubscriptionStatus } from '@/types';

interface SubscriptionStatusBadgeProps {
  status: SubscriptionStatus;
}

export const SubscriptionStatusBadge: React.FC<SubscriptionStatusBadgeProps> = ({ status }) => {
  const getStatusConfig = (status: SubscriptionStatus) => {
    switch (status) {
      case SubscriptionStatus.ACTIVE:
        return { label: 'Active', className: 'bg-green-100 text-green-800' };
      case SubscriptionStatus.INACTIVE:
        return { label: 'Inactive', className: 'bg-gray-100 text-gray-800' };
      case SubscriptionStatus.EXPIRED:
        return { label: 'Expired', className: 'bg-red-100 text-red-800' };
      case SubscriptionStatus.SUSPENDED:
        return { label: 'Suspended', className: 'bg-yellow-100 text-yellow-800' };
      default:
        return { label: status, className: 'bg-gray-100 text-gray-800' };
    }
  };

  const config = getStatusConfig(status);
  return <Badge className={config.className}>{config.label}</Badge>;
};

