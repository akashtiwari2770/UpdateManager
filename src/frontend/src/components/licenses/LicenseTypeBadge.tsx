import React from 'react';
import { Badge } from '@/components/ui';
import { LicenseType } from '@/types';

interface LicenseTypeBadgeProps {
  type: LicenseType;
}

export const LicenseTypeBadge: React.FC<LicenseTypeBadgeProps> = ({ type }) => {
  const getTypeConfig = (type: LicenseType) => {
    switch (type) {
      case LicenseType.PERPETUAL:
        return { label: 'Perpetual', className: 'bg-blue-100 text-blue-800' };
      case LicenseType.TIME_BASED:
        return { label: 'Time-based', className: 'bg-purple-100 text-purple-800' };
      default:
        return { label: type, className: 'bg-gray-100 text-gray-800' };
    }
  };

  const config = getTypeConfig(type);
  return <Badge className={config.className}>{config.label}</Badge>;
};

