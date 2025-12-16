import React from 'react';
import { DeploymentType } from '@/types';
import { Badge } from '@/components/ui';

interface DeploymentTypeBadgeProps {
  type: DeploymentType;
}

export const DeploymentTypeBadge: React.FC<DeploymentTypeBadgeProps> = ({ type }) => {
  const typeConfig = {
    [DeploymentType.UAT]: {
      label: 'UAT',
      className: 'bg-yellow-100 text-yellow-800',
    },
    [DeploymentType.TESTING]: {
      label: 'Testing',
      className: 'bg-blue-100 text-blue-800',
    },
    [DeploymentType.PRODUCTION]: {
      label: 'Production',
      className: 'bg-green-100 text-green-800',
    },
  };

  const config = typeConfig[type] || typeConfig[DeploymentType.UAT];

  return <Badge className={config.className}>{config.label}</Badge>;
};

