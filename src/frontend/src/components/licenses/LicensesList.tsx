import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { licensesApi } from '@/services/api/licenses';
import { License, LicenseType, LicenseStatus, ListLicensesQuery } from '@/types';
import { Button, Card, Spinner, Select } from '@/components/ui';
import { LicenseTypeBadge } from './LicenseTypeBadge';
import { LicenseStatusBadge } from './LicenseStatusBadge';

interface LicensesListProps {
  customerId: string;
  subscriptionId: string;
}

export const LicensesList: React.FC<LicensesListProps> = ({ customerId, subscriptionId }) => {
  const navigate = useNavigate();
  const [licenses, setLicenses] = useState<License[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<ListLicensesQuery>({
    page: 1,
    limit: 20,
  });
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    loadLicenses();
  }, [customerId, subscriptionId, filters]);

  const loadLicenses = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await licensesApi.getAll(customerId, subscriptionId, filters);
      setLicenses(response?.data || []);
      setTotalPages(response?.pagination?.total_pages || 1);
      setTotal(response?.pagination?.total || 0);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load licenses');
      console.error('Error loading licenses:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (key: keyof ListLicensesQuery, value: any) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value,
      page: 1,
    }));
  };

  const handlePageChange = (newPage: number) => {
    setFilters((prev) => ({ ...prev, page: newPage }));
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  if (loading && licenses.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-900">Licenses</h2>
        <Button
          onClick={() => navigate(`/customers/${customerId}/subscriptions/${subscriptionId}/licenses/new`)}
          variant="primary"
        >
          Assign License
        </Button>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      <Card>
        <div className="space-y-4">
          {/* Filters */}
          <div className="flex items-center gap-4">
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-700 mb-1">License Type</label>
              <Select
                value={filters.license_type || ''}
                onChange={(e) =>
                  handleFilterChange('license_type', e.target.value || undefined)
                }
              >
                <option value="">All Types</option>
                <option value={LicenseType.PERPETUAL}>Perpetual</option>
                <option value={LicenseType.TIME_BASED}>Time-based</option>
              </Select>
            </div>
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
              <Select
                value={filters.status || ''}
                onChange={(e) => handleFilterChange('status', e.target.value || undefined)}
              >
                <option value="">All Statuses</option>
                <option value={LicenseStatus.ACTIVE}>Active</option>
                <option value={LicenseStatus.INACTIVE}>Inactive</option>
                <option value={LicenseStatus.EXPIRED}>Expired</option>
                <option value={LicenseStatus.REVOKED}>Revoked</option>
              </Select>
            </div>
          </div>

          {/* Licenses Table */}
          {licenses.length === 0 ? (
            <div className="text-center py-8 text-gray-500">No licenses found</div>
          ) : (
            <>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        License ID
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Product
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Type
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Seats
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        End Date
                      </th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {licenses.map((license) => (
                      <tr key={license.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {license.license_id}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {license.product_id}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <LicenseTypeBadge type={license.license_type} />
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {license.number_of_seats}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <LicenseStatusBadge status={license.status} />
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {license.end_date ? formatDate(license.end_date) : '-'}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <Button
                            variant="secondary"
                            size="sm"
                            onClick={() =>
                              navigate(
                                `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${license.license_id}`
                              )
                            }
                          >
                            View
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {/* Pagination */}
              {totalPages > 1 && (
                <div className="flex items-center justify-between border-t border-gray-200 px-4 py-3">
                  <div className="text-sm text-gray-700">
                    Showing page {filters.page} of {totalPages} ({total} total)
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="secondary"
                      size="sm"
                      onClick={() => handlePageChange((filters.page || 1) - 1)}
                      disabled={filters.page === 1}
                    >
                      Previous
                    </Button>
                    <Button
                      variant="secondary"
                      size="sm"
                      onClick={() => handlePageChange((filters.page || 1) + 1)}
                      disabled={filters.page === totalPages}
                    >
                      Next
                    </Button>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </Card>
    </div>
  );
};

