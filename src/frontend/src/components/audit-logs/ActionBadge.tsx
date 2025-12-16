import React from 'react';
import { AuditAction } from '@/types';
import { Badge } from '@/components/ui';

interface ActionBadgeProps {
  action: AuditAction;
}

export const ActionBadge: React.FC<ActionBadgeProps> = ({ action }) => {
  const getActionConfig = (action: AuditAction) => {
    switch (action) {
      case AuditAction.CREATE:
        return {
          label: 'Create',
          className: 'bg-green-100 text-green-800',
        };
      case AuditAction.UPDATE:
        return {
          label: 'Update',
          className: 'bg-blue-100 text-blue-800',
        };
      case AuditAction.DELETE:
        return {
          label: 'Delete',
          className: 'bg-red-100 text-red-800',
        };
      case AuditAction.APPROVE:
        return {
          label: 'Approve',
          className: 'bg-green-100 text-green-800',
        };
      case AuditAction.REJECT:
        return {
          label: 'Reject',
          className: 'bg-orange-100 text-orange-800',
        };
      case AuditAction.RELEASE:
        return {
          label: 'Release',
          className: 'bg-blue-100 text-blue-800',
        };
      case AuditAction.UPLOAD:
        return {
          label: 'Upload',
          className: 'bg-purple-100 text-purple-800',
        };
      case AuditAction.DOWNLOAD:
        return {
          label: 'Download',
          className: 'bg-indigo-100 text-indigo-800',
        };
      default:
        return {
          label: action,
          className: 'bg-gray-100 text-gray-800',
        };
    }
  };

  const config = getActionConfig(action);

  return <Badge className={config.className}>{config.label}</Badge>;
};

