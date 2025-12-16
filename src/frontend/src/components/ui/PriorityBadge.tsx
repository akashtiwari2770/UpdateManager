import React from 'react';
import { Badge } from './Badge';

interface PriorityBadgeProps {
  priority: 'critical' | 'high' | 'normal';
}

export const PriorityBadge: React.FC<PriorityBadgeProps> = ({ priority }) => {
  const priorityConfig = {
    critical: {
      className: 'bg-red-100 text-red-800',
      label: 'Critical',
      icon: '🔴',
    },
    high: {
      className: 'bg-orange-100 text-orange-800',
      label: 'High',
      icon: '🟠',
    },
    normal: {
      className: 'bg-blue-100 text-blue-800',
      label: 'Normal',
      icon: '🔵',
    },
  };

  const config = priorityConfig[priority] || priorityConfig.normal;

  return (
    <Badge className={config.className}>
      <span className="mr-1">{config.icon}</span>
      {config.label}
    </Badge>
  );
};

