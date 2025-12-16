import React from 'react';
import { Badge } from '@/components/ui';
import { LicenseStatus } from '@/types';

interface LicenseStatusBadgeProps {
  status: LicenseStatus;
}

export const LicenseStatusBadge: React.FC<LicenseStatusBadgeProps> = ({ status }) => {
  const getStatusConfig = (status: LicenseStatus) => {
    switch (status) {
      case LicenseStatus.ACTIVE:
        return { label: 'Active', className: 'bg-green-100 text-green-800' };
      case LicenseStatus.INACTIVE:
        return { label: 'Inactive', className: 'bg-gray-100 text-gray-800' };
      case LicenseStatus.EXPIRED:
        return { label: 'Expired', className: 'bg-red-100 text-red-800' };
      case LicenseStatus.REVOKED:
        return { label: 'Revoked', className: 'bg-orange-100 text-orange-800' };
      default:
        return { label: status, className: 'bg-gray-100 text-gray-800' };
    }
  };

  const config = getStatusConfig(status);
  return <Badge className={config.className}>{config.label}</Badge>;
};

