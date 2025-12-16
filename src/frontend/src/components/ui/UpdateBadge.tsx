import React from 'react';
import { Badge } from './Badge';

interface UpdateBadgeProps {
  count: number;
  priority?: 'critical' | 'high' | 'normal';
}

export const UpdateBadge: React.FC<UpdateBadgeProps> = ({ count, priority = 'normal' }) => {
  if (count === 0) {
    return (
      <Badge className="bg-green-100 text-green-800">Up to date</Badge>
    );
  }

  const priorityConfig = {
    critical: { className: 'bg-red-100 text-red-800', label: `${count} Critical` },
    high: { className: 'bg-orange-100 text-orange-800', label: `${count} Updates` },
    normal: { className: 'bg-blue-100 text-blue-800', label: `${count} Updates` },
  };

  const config = priorityConfig[priority] || priorityConfig.normal;

  return <Badge className={config.className}>{config.label}</Badge>;
};

