import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { pendingUpdatesApi } from '@/services/api/pending-updates';
import { PendingUpdatesResponse, PendingUpdatesQuery } from '@/types';
import { Button, Card, Spinner, Select, Input } from '@/components/ui';
import { UpdateBadge } from '@/components/ui/UpdateBadge';
import { PriorityBadge } from '@/components/ui/PriorityBadge';
import { DeploymentTypeBadge } from '@/components/deployments/DeploymentTypeBadge';

interface PendingUpdatesListProps {
  customerId?: string;
  view?: 'all' | 'customer' | 'tenant';
}

export const PendingUpdatesList: React.FC<PendingUpdatesListProps> = ({
  customerId,
  view = 'all',
}) => {
  const navigate = useNavigate();
  const [pendingUpdates, setPendingUpdates] = useState<PendingUpdatesResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<PendingUpdatesQuery>({
    page: 1,
    limit: 20,
  });
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    loadPendingUpdates();
  }, [customerId, filters, view]);

  const loadPendingUpdates = async () => {
    try {
      setLoading(true);
      setError(null);

      if (view === 'all') {
        const result = await pendingUpdatesApi.getAllPendingUpdates(filters);
        setPendingUpdates(result.data || []);
        setTotalPages(result.pagination?.total_pages || 1);
        setTotal(result.pagination?.total || 0);
      } else {
        // For customer/tenant views, we'd need additional API methods
        // For now, use getAllPendingUpdates with customer filter
        const query = customerId ? { ...filters, customer_id: customerId } : filters;
        const result = await pendingUpdatesApi.getAllPendingUpdates(query);
        setPendingUpdates(result.data || []);
        setTotalPages(result.pagination?.total_pages || 1);
        setTotal(result.pagination?.total || 0);
      }
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load pending updates');
      console.error('Error loading pending updates:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (key: keyof PendingUpdatesQuery, value: any) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value,
      page: 1,
    }));
  };

  const handlePageChange = (newPage: number) => {
    setFilters((prev) => ({ ...prev, page: newPage }));
  };

  const handleViewDeployment = (deployment: PendingUpdatesResponse) => {
    if (deployment.customer_id && deployment.tenant_id && deployment.deployment_id) {
      navigate(
        `/customers/${deployment.customer_id}/tenants/${deployment.tenant_id}/deployments/${deployment.deployment_id}`
      );
    } else {
      console.warn('Missing required IDs for navigation:', {
        customer_id: deployment.customer_id,
        tenant_id: deployment.tenant_id,
        deployment_id: deployment.deployment_id,
      });
    }
  };

  if (loading && pendingUpdates.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  if (error && pendingUpdates.length === 0) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
        {error}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {/* Filters */}
      <Card className="p-4">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Select
            label="Priority"
            value={filters.priority || ''}
            onChange={(e) => handleFilterChange('priority', e.target.value || undefined)}
            options={[
              { value: '', label: 'All Priorities' },
              { value: 'critical', label: 'Critical' },
              { value: 'high', label: 'High' },
              { value: 'normal', label: 'Normal' },
            ]}
          />
          <Input
            label="Product ID"
            value={filters.product_id || ''}
            onChange={(e) => handleFilterChange('product_id', e.target.value || undefined)}
            placeholder="Filter by product"
          />
          <Select
            label="Deployment Type"
            value={filters.deployment_type || ''}
            onChange={(e) =>
              handleFilterChange('deployment_type', (e.target.value as any) || undefined)
            }
            options={[
              { value: '', label: 'All Types' },
              { value: 'uat', label: 'UAT' },
              { value: 'testing', label: 'Testing' },
              { value: 'production', label: 'Production' },
            ]}
          />
        </div>
      </Card>

      {/* Pending Updates Table */}
      {pendingUpdates.length === 0 ? (
        <Card>
          <div className="text-center py-12">
            <p className="text-gray-500">No deployments with pending updates found</p>
          </div>
        </Card>
      ) : (
        <>
          <Card>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Customer
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Tenant
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Product
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Type
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Current
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Latest
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Updates
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Priority
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {pendingUpdates.map((update) => (
                    <tr key={update.deployment_id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {update.customer_name || update.customer_id || '-'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {update.tenant_name || update.tenant_id || '-'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {update.product_id}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {update.deployment_type && (
                          <DeploymentTypeBadge type={update.deployment_type} />
                        )}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {update.current_version}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {update.latest_version}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <UpdateBadge count={update.update_count} priority={update.priority} />
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <PriorityBadge priority={update.priority} />
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        {update.customer_id && update.tenant_id && update.deployment_id ? (
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleViewDeployment(update)}
                          >
                            View
                          </Button>
                        ) : (
                          <span className="text-gray-400 text-xs">N/A</span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </Card>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6">
              <div className="flex flex-1 justify-between sm:hidden">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handlePageChange(Math.max(1, (filters.page || 1) - 1))}
                  disabled={filters.page === 1}
                >
                  Previous
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handlePageChange(Math.min(totalPages, (filters.page || 1) + 1))}
                  disabled={filters.page === totalPages}
                >
                  Next
                </Button>
              </div>
              <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                <div>
                  <p className="text-sm text-gray-700">
                    Showing{' '}
                    <span className="font-medium">
                      {((filters.page || 1) - 1) * (filters.limit || 20) + 1}
                    </span>{' '}
                    to{' '}
                    <span className="font-medium">
                      {Math.min((filters.page || 1) * (filters.limit || 20), total)}
                    </span>{' '}
                    of <span className="font-medium">{total}</span> results
                  </p>
                </div>
                <div>
                  <nav className="isolate inline-flex -space-x-px rounded-md shadow-sm">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handlePageChange(Math.max(1, (filters.page || 1) - 1))}
                      disabled={filters.page === 1}
                    >
                      Previous
                    </Button>
                    {Array.from({ length: totalPages }, (_, i) => i + 1)
                      .filter(
                        (page) =>
                          page === 1 ||
                          page === totalPages ||
                          (page >= (filters.page || 1) - 2 && page <= (filters.page || 1) + 2)
                      )
                      .map((page, idx, arr) => (
                        <React.Fragment key={page}>
                          {idx > 0 && arr[idx - 1] !== page - 1 && (
                            <span className="px-4 py-2 text-sm text-gray-700">...</span>
                          )}
                          <Button
                            variant={filters.page === page ? 'primary' : 'ghost'}
                            size="sm"
                            onClick={() => handlePageChange(page)}
                          >
                            {page}
                          </Button>
                        </React.Fragment>
                      ))}
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handlePageChange(Math.min(totalPages, (filters.page || 1) + 1))}
                      disabled={filters.page === totalPages}
                    >
                      Next
                    </Button>
                  </nav>
                </div>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
};

