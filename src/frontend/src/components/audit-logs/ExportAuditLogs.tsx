import React, { useState } from 'react';
import { auditLogsApi } from '@/services/api/audit-logs';
import { ListAuditLogsQuery, AuditLog } from '@/types';
import { Button, Modal, Select, Alert } from '@/components/ui';

interface ExportAuditLogsProps {
  filters: ListAuditLogsQuery;
}

export const ExportAuditLogs: React.FC<ExportAuditLogsProps> = ({ filters }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [format, setFormat] = useState<'csv' | 'json'>('csv');
  const [exporting, setExporting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const exportToCSV = (logs: AuditLog[]) => {
    const headers = ['Timestamp', 'User', 'Action', 'Resource Type', 'Resource ID', 'IP Address', 'User Agent'];
    const rows = logs.map((log) => [
      log.timestamp,
      log.user_email,
      log.action,
      log.resource_type,
      log.resource_id,
      log.ip_address,
      log.user_agent,
    ]);

    const csvContent = [
      headers.join(','),
      ...rows.map((row) => row.map((cell) => `"${String(cell).replace(/"/g, '""')}"`).join(',')),
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', `audit-logs-${new Date().toISOString().split('T')[0]}.csv`);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const exportToJSON = (logs: AuditLog[]) => {
    const jsonContent = JSON.stringify(logs, null, 2);
    const blob = new Blob([jsonContent], { type: 'application/json' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', `audit-logs-${new Date().toISOString().split('T')[0]}.json`);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const handleExport = async () => {
    setExporting(true);
    setError(null);

    try {
      // Fetch all logs matching the filters
      // Note: In a real implementation, you might need to paginate through all results
      const response = await auditLogsApi.getAll({
        ...filters,
        limit: 1000, // Adjust based on your needs
      });

      if (format === 'csv') {
        exportToCSV(response.data);
      } else {
        exportToJSON(response.data);
      }

      setIsOpen(false);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to export audit logs');
    } finally {
      setExporting(false);
    }
  };

  return (
    <>
      <Button variant="secondary" onClick={() => setIsOpen(true)}>
        Export Audit Logs
      </Button>

      <Modal
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="Export Audit Logs"
        size="md"
        footer={
          <div className="flex gap-2 justify-end">
            <Button variant="secondary" onClick={() => setIsOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleExport} isLoading={exporting}>
              Export
            </Button>
          </div>
        }
      >
        <div className="space-y-4">
          {error && <Alert variant="error">{error}</Alert>}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Export Format
            </label>
            <Select value={format} onChange={(e) => setFormat(e.target.value as 'csv' | 'json')}>
              <option value="csv">CSV</option>
              <option value="json">JSON</option>
            </Select>
          </div>

          <div className="text-sm text-gray-600">
            <p>
              The export will include all audit logs matching your current filters.
              {filters.limit && filters.limit < 1000 && (
                <span className="text-yellow-600 block mt-1">
                  Note: Only the first {filters.limit} results will be exported. Consider adjusting filters for a complete export.
                </span>
              )}
            </p>
          </div>
        </div>
      </Modal>
    </>
  );
};

