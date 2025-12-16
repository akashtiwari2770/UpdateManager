import React from 'react';
import { ListAuditLogsQuery, AuditAction } from '@/types';
import { Input, Select, Button } from '@/components/ui';

interface AuditLogFiltersProps {
  filters: ListAuditLogsQuery;
  onFilterChange: (filters: ListAuditLogsQuery) => void;
  onClear: () => void;
  users?: Array<{ id: string; email: string }>;
}

export const AuditLogFilters: React.FC<AuditLogFiltersProps> = ({
  filters,
  onFilterChange,
  onClear,
  users = [],
}) => {
  const handleChange = (key: keyof ListAuditLogsQuery, value: any) => {
    onFilterChange({
      ...filters,
      [key]: value || undefined,
      page: 1, // Reset to first page on filter change
    });
  };

  const hasActiveFilters = Boolean(
    filters.user_id ||
    filters.action ||
    filters.resource_type ||
    filters.resource_id ||
    filters.start_date ||
    filters.end_date
  );

  return (
    <div className="bg-white p-4 rounded-lg border border-gray-200 space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Filters</h3>
        {hasActiveFilters && (
          <Button variant="ghost" size="sm" onClick={onClear}>
            Clear Filters
          </Button>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {/* User Filter */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">User</label>
          <Select
            value={filters.user_id || ''}
            onChange={(e) => handleChange('user_id', e.target.value)}
          >
            <option value="">All Users</option>
            {users.map((user) => (
              <option key={user.id} value={user.id}>
                {user.email}
              </option>
            ))}
          </Select>
        </div>

        {/* Action Filter */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Action</label>
          <Select
            value={filters.action || ''}
            onChange={(e) => handleChange('action', e.target.value as AuditAction)}
          >
            <option value="">All Actions</option>
            <option value={AuditAction.CREATE}>Create</option>
            <option value={AuditAction.UPDATE}>Update</option>
            <option value={AuditAction.DELETE}>Delete</option>
            <option value={AuditAction.APPROVE}>Approve</option>
            <option value={AuditAction.REJECT}>Reject</option>
            <option value={AuditAction.RELEASE}>Release</option>
            <option value={AuditAction.UPLOAD}>Upload</option>
            <option value={AuditAction.DOWNLOAD}>Download</option>
          </Select>
        </div>

        {/* Resource Type Filter */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Resource Type</label>
          <Select
            value={filters.resource_type || ''}
            onChange={(e) => handleChange('resource_type', e.target.value)}
          >
            <option value="">All Types</option>
            <option value="product">Product</option>
            <option value="version">Version</option>
            <option value="package">Package</option>
            <option value="rollout">Rollout</option>
            <option value="notification">Notification</option>
          </Select>
        </div>

        {/* Resource ID Search */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Resource ID</label>
          <Input
            type="text"
            placeholder="Search by resource ID..."
            value={filters.resource_id || ''}
            onChange={(e) => handleChange('resource_id', e.target.value)}
          />
        </div>

        {/* Start Date */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Start Date</label>
          <Input
            type="date"
            value={filters.start_date || ''}
            onChange={(e) => handleChange('start_date', e.target.value)}
          />
        </div>

        {/* End Date */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">End Date</label>
          <Input
            type="date"
            value={filters.end_date || ''}
            onChange={(e) => handleChange('end_date', e.target.value)}
          />
        </div>
      </div>
    </div>
  );
};

