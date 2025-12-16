import React from 'react';
import { Badge } from './Badge';

interface VersionGapBadgeProps {
  gapType: 'patch' | 'minor' | 'major';
}

export const VersionGapBadge: React.FC<VersionGapBadgeProps> = ({ gapType }) => {
  const gapConfig = {
    patch: {
      className: 'bg-green-100 text-green-800',
      label: 'Patch',
    },
    minor: {
      className: 'bg-yellow-100 text-yellow-800',
      label: 'Minor',
    },
    major: {
      className: 'bg-purple-100 text-purple-800',
      label: 'Major',
    },
  };

  const config = gapConfig[gapType] || gapConfig.patch;

  return <Badge className={config.className}>{config.label}</Badge>;
};

