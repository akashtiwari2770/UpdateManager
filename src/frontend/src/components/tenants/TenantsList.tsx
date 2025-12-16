import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { tenantsApi } from '@/services/api/tenants';
import { customersApi } from '@/services/api/customers';
import { CustomerTenant, TenantStatus, ListTenantsQuery } from '@/types';
import { Button, Card, Badge, Spinner, Select } from '@/components/ui';

interface TenantsListProps {
  customerId: string;
}

export const TenantsList: React.FC<TenantsListProps> = ({ customerId }) => {
  const navigate = useNavigate();
  const [tenants, setTenants] = useState<CustomerTenant[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<ListTenantsQuery>({
    page: 1,
    limit: 20,
  });
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    loadTenants();
  }, [customerId, filters]);

  const loadTenants = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await customersApi.getTenants(customerId, filters);
      setTenants(response?.data || []);
      setTotalPages(response?.pagination?.total_pages || 1);
      setTotal(response?.pagination?.total || 0);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load tenants');
      console.error('Error loading tenants:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (key: keyof ListTenantsQuery, value: any) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value,
      page: 1,
    }));
  };

  const handlePageChange = (newPage: number) => {
    setFilters((prev) => ({ ...prev, page: newPage }));
  };

  const handleDelete = async (tenantId: string) => {
    if (!confirm('Are you sure you want to delete this tenant?')) {
      return;
    }

    try {
      await tenantsApi.delete(customerId, tenantId);
      loadTenants();
    } catch (err: any) {
      alert(err.response?.data?.error?.message || 'Failed to delete tenant');
    }
  };

  const getStatusBadge = (status: TenantStatus) => {
    const statusConfig = {
      [TenantStatus.ACTIVE]: { label: 'Active', className: 'bg-green-100 text-green-800' },
      [TenantStatus.INACTIVE]: { label: 'Inactive', className: 'bg-gray-100 text-gray-800' },
    };
    const config = statusConfig[status] || statusConfig[TenantStatus.INACTIVE];
    return <Badge className={config.className}>{config.label}</Badge>;
  };

  if (loading && tenants.length === 0) {
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
      <div className="flex items-center gap-4">
        <Select
          label="Status"
          value={filters.status || ''}
          onChange={(e) => handleFilterChange('status', e.target.value || undefined)}
          options={[
            { value: '', label: 'All Statuses' },
            { value: TenantStatus.ACTIVE, label: 'Active' },
            { value: TenantStatus.INACTIVE, label: 'Inactive' },
          ]}
        />
      </div>

      {/* Tenants Table */}
      {tenants.length === 0 ? (
        <Card>
          <div className="text-center py-12">
            <p className="text-gray-500 mb-4">No tenants found</p>
            <Button
              variant="primary"
              onClick={() => navigate(`/customers/${customerId}/tenants/new`)}
            >
              Create First Tenant
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
                      Tenant ID
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Name
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Description
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {tenants.map((tenant) => (
                    <tr key={tenant.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {tenant.tenant_id}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {tenant.name}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500">
                        {tenant.description || '-'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">{getStatusBadge(tenant.status)}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        <div className="flex space-x-2">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() =>
                              navigate(`/customers/${customerId}/tenants/${tenant.id}`)
                            }
                          >
                            View
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() =>
                              navigate(`/customers/${customerId}/tenants/${tenant.id}/edit`)
                            }
                          >
                            Edit
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDelete(tenant.id)}
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

