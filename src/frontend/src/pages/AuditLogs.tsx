import React, { useState } from 'react';
import { AuditLogsList, AuditLogFilters, ExportAuditLogs } from '@/components/audit-logs';
import { ListAuditLogsQuery } from '@/types';

export const AuditLogs: React.FC = () => {
  const [filters, setFilters] = useState<ListAuditLogsQuery>({
    page: 1,
    limit: 20,
  });

  const handleFilterChange = (newFilters: ListAuditLogsQuery) => {
    setFilters(newFilters);
  };

  const handleClearFilters = () => {
    setFilters({ page: 1, limit: 20 });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Audit Logs</h1>
        <ExportAuditLogs filters={filters} />
      </div>

      <AuditLogFilters
        filters={filters}
        onFilterChange={handleFilterChange}
        onClear={handleClearFilters}
      />

      <AuditLogsList filters={filters} onFiltersChange={handleFilterChange} />
    </div>
  );
};

