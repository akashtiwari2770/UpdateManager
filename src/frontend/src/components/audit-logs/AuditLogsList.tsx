import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { auditLogsApi } from '@/services/api/audit-logs';
import { AuditLog, ListAuditLogsQuery } from '@/types';
import { Button, Spinner, Alert, Badge } from '@/components/ui';
import { ActionBadge } from './ActionBadge';
import { AuditLogDetails } from './AuditLogDetails';

interface AuditLogsListProps {
  filters?: ListAuditLogsQuery;
  showFilters?: boolean;
  onFiltersChange?: (filters: ListAuditLogsQuery) => void;
}

export const AuditLogsList: React.FC<AuditLogsListProps> = ({
  filters: initialFilters,
  showFilters = false,
  onFiltersChange,
}) => {
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());
  const navigate = useNavigate();

  const filters = initialFilters || { page: 1, limit: 20 };

  useEffect(() => {
    loadAuditLogs();
  }, [filters, sortOrder]);

  const loadAuditLogs = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await auditLogsApi.getAll({
        ...filters,
        page: filters.page || page,
      });
      setAuditLogs(response.data);
      setTotalPages(response.pagination?.total_pages || 1);
      setPage(response.pagination?.current_page || 1);
    } catch (err) {
      setError('Failed to load audit logs');
      console.error('Failed to load audit logs:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSort = () => {
    setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    // In a real implementation, you'd sort the data or make a new API call
  };

  const toggleRowExpansion = (id: string) => {
    setExpandedRows((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  };

  const handleResourceClick = (resourceType: string, resourceId: string) => {
    switch (resourceType.toLowerCase()) {
      case 'product':
        navigate(`/products/${resourceId}`);
        break;
      case 'version':
        navigate(`/versions/${resourceId}`);
        break;
      default:
        // For other types, we might not have a page yet
        break;
    }
  };

  const formatTimestamp = (timestamp: string) => {
    return new Date(timestamp).toLocaleString();
  };

  const getUserInitials = (email: string) => {
    return email
      .split('@')[0]
      .split('.')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .substring(0, 2);
  };

  return (
    <div className="space-y-4">
      {error && <Alert variant="error">{error}</Alert>}

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Spinner />
        </div>
      ) : auditLogs.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          <p>No audit logs found</p>
        </div>
      ) : (
        <>
          <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      <button
                        onClick={handleSort}
                        className="flex items-center gap-1 hover:text-gray-700"
                      >
                        Timestamp
                        {sortOrder === 'asc' ? (
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                          </svg>
                        ) : (
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                          </svg>
                        )}
                      </button>
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      User
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Action
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Resource Type
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Resource ID
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Details
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {auditLogs.map((log) => (
                    <React.Fragment key={log.id}>
                      <tr className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {formatTimestamp(log.timestamp)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center">
                            <div className="flex-shrink-0 h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center">
                              <span className="text-xs font-medium text-blue-800">
                                {getUserInitials(log.user_email)}
                              </span>
                            </div>
                            <div className="ml-3">
                              <div className="text-sm font-medium text-gray-900">{log.user_email}</div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <ActionBadge action={log.action} />
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          <Badge className="bg-gray-100 text-gray-800">{log.resource_type}</Badge>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          <button
                            onClick={() => handleResourceClick(log.resource_type, log.resource_id)}
                            className="text-blue-600 hover:text-blue-800 hover:underline"
                          >
                            {log.resource_id}
                          </button>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          <button
                            onClick={() => toggleRowExpansion(log.id)}
                            className="text-blue-600 hover:text-blue-800"
                          >
                            {expandedRows.has(log.id) ? 'Hide' : 'Show'}
                          </button>
                        </td>
                      </tr>
                      {expandedRows.has(log.id) && (
                        <tr>
                          <td colSpan={6} className="px-6 py-4">
                            <AuditLogDetails auditLog={log} />
                          </td>
                        </tr>
                      )}
                    </React.Fragment>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between">
              <div className="text-sm text-gray-600">
                Page {page} of {totalPages}
              </div>
              <div className="flex gap-2">
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => {
                    const newFilters = { ...filters, page: page - 1 };
                    if (onFiltersChange) {
                      onFiltersChange(newFilters);
                    }
                  }}
                  disabled={page === 1}
                >
                  Previous
                </Button>
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => {
                    const newFilters = { ...filters, page: page + 1 };
                    if (onFiltersChange) {
                      onFiltersChange(newFilters);
                    }
                  }}
                  disabled={page === totalPages}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
};

