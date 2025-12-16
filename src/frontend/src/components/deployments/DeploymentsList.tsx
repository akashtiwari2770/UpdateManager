import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { deploymentsApi } from '@/services/api/deployments';
import { tenantsApi } from '@/services/api/tenants';
import { Deployment, DeploymentType, DeploymentStatus, ListDeploymentsQuery } from '@/types';
import { Button, Card, Badge, Spinner, Select } from '@/components/ui';
import { DeploymentTypeBadge } from './DeploymentTypeBadge';
import { UpdateBadge } from '@/components/ui/UpdateBadge';
import { pendingUpdatesApi } from '@/services/api/pending-updates';
import { PendingUpdatesResponse } from '@/types';

interface DeploymentsListProps {
  customerId: string;
  tenantId: string;
}

export const DeploymentsList: React.FC<DeploymentsListProps> = ({
  customerId,
  tenantId,
}) => {
  const navigate = useNavigate();
  const [deployments, setDeployments] = useState<Deployment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<ListDeploymentsQuery>({
    page: 1,
    limit: 20,
  });
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);
  const [pendingUpdatesMap, setPendingUpdatesMap] = useState<Record<string, PendingUpdatesResponse>>({});

  useEffect(() => {
    loadDeployments();
  }, [customerId, tenantId, filters]);

  useEffect(() => {
    if (deployments.length > 0) {
      loadPendingUpdates();
    }
  }, [deployments, customerId, tenantId]);

  const loadDeployments = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await tenantsApi.getDeployments(customerId, tenantId, filters);
      setDeployments(response?.data || []);
      setTotalPages(response?.pagination?.total_pages || 1);
      setTotal(response?.pagination?.total || 0);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load deployments');
      console.error('Error loading deployments:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadPendingUpdates = async () => {
    const updatesMap: Record<string, PendingUpdatesResponse> = {};
    const promises = deployments.map(async (deployment) => {
      try {
        const updates = await pendingUpdatesApi.getDeploymentPendingUpdates(
          customerId,
          tenantId,
          deployment.id
        );
        updatesMap[deployment.id] = updates;
      } catch (err) {
        // Silently fail for individual deployments
        console.error(`Failed to load pending updates for deployment ${deployment.id}:`, err);
      }
    });
    await Promise.all(promises);
    setPendingUpdatesMap(updatesMap);
  };

  const handleFilterChange = (key: keyof ListDeploymentsQuery, value: any) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value,
      page: 1,
    }));
  };

  const handlePageChange = (newPage: number) => {
    setFilters((prev) => ({ ...prev, page: newPage }));
  };

  const handleDelete = async (deploymentId: string) => {
    if (!confirm('Are you sure you want to delete this deployment?')) {
      return;
    }

    try {
      await deploymentsApi.delete(customerId, tenantId, deploymentId);
      loadDeployments();
    } catch (err: any) {
      alert(err.response?.data?.error?.message || 'Failed to delete deployment');
    }
  };

  const getStatusBadge = (status: DeploymentStatus) => {
    const statusConfig = {
      [DeploymentStatus.ACTIVE]: { label: 'Active', className: 'bg-green-100 text-green-800' },
      [DeploymentStatus.INACTIVE]: { label: 'Inactive', className: 'bg-gray-100 text-gray-800' },
    };
    const config = statusConfig[status] || statusConfig[DeploymentStatus.INACTIVE];
    return <Badge className={config.className}>{config.label}</Badge>;
  };

  if (loading && deployments.length === 0) {
    return (
      <div className="flex items-center justify-center h-32">
        <Spinner size="md" />
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
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Select
          label="Deployment Type"
          value={filters.deployment_type || ''}
          onChange={(e) =>
            handleFilterChange('deployment_type', e.target.value || undefined)
          }
          options={[
            { value: '', label: 'All Types' },
            { value: DeploymentType.UAT, label: 'UAT' },
            { value: DeploymentType.TESTING, label: 'Testing' },
            { value: DeploymentType.PRODUCTION, label: 'Production' },
          ]}
        />
        <Select
          label="Status"
          value={filters.status || ''}
          onChange={(e) => handleFilterChange('status', e.target.value || undefined)}
          options={[
            { value: '', label: 'All Statuses' },
            { value: DeploymentStatus.ACTIVE, label: 'Active' },
            { value: DeploymentStatus.INACTIVE, label: 'Inactive' },
          ]}
        />
      </div>

      {/* Deployments Table */}
      {deployments.length === 0 ? (
        <Card>
          <div className="text-center py-12">
            <p className="text-gray-500 mb-4">No deployments found</p>
            <Button
              variant="primary"
              onClick={() =>
                navigate(`/customers/${customerId}/tenants/${tenantId}/deployments/new`)
              }
            >
              Create First Deployment
            </Button>
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
                      Product
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Type
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Version
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Users
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Updates
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {deployments.map((deployment) => (
                    <tr key={deployment.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {deployment.product_id}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <DeploymentTypeBadge type={deployment.deployment_type} />
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {deployment.installed_version}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {deployment.number_of_users || '-'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">{getStatusBadge(deployment.status)}</td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {pendingUpdatesMap[deployment.id] ? (
                          <UpdateBadge
                            count={pendingUpdatesMap[deployment.id].update_count}
                            priority={pendingUpdatesMap[deployment.id].priority}
                          />
                        ) : (
                          <span className="text-gray-400 text-sm">Loading...</span>
                        )}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        <div className="flex space-x-2">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() =>
                              navigate(
                                `/customers/${customerId}/tenants/${tenantId}/deployments/${deployment.id}`
                              )
                            }
                          >
                            View
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() =>
                              navigate(
                                `/customers/${customerId}/tenants/${tenantId}/deployments/${deployment.id}/edit`
                              )
                            }
                          >
                            Edit
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDelete(deployment.id)}
                          >
                            Delete
                          </Button>
                        </div>
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

