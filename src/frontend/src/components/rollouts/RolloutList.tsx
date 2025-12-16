import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { UpdateRollout, RolloutStatus, Product } from '@/types';
import { Card, Badge, Select, Spinner, Alert, Button } from '@/components/ui';
import { updateRolloutsApi, ListUpdateRolloutsQuery } from '@/services/api/update-rollouts';
import { productsApi } from '@/services/api/products';

export const RolloutList: React.FC = () => {
  const navigate = useNavigate();
  const [rollouts, setRollouts] = useState<UpdateRollout[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Filters and pagination
  const [filters, setFilters] = useState<ListUpdateRolloutsQuery>({
    page: 1,
    limit: 25,
  });
  const [productFilter, setProductFilter] = useState('');
  const [statusFilter, setStatusFilter] = useState<RolloutStatus | ''>('');
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    loadProducts();
  }, []);

  useEffect(() => {
    loadRollouts();
  }, [filters, productFilter, statusFilter]);

  const loadProducts = async () => {
    try {
      const response = await productsApi.getAll({ page: 1, limit: 1000, is_active: true });
      setProducts(response.data || []);
    } catch (err) {
      console.error('Error loading products:', err);
    }
  };

  const loadRollouts = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const query: ListUpdateRolloutsQuery = {
        page: filters.page,
        limit: filters.limit,
      };
      
      if (productFilter) {
        query.product_id = productFilter;
      }
      
      if (statusFilter) {
        query.status = statusFilter as RolloutStatus;
      }
      
      const response = await updateRolloutsApi.list(query);
      setRollouts(response.data || []);
      
      if (response.meta) {
        setTotalPages(response.meta.total_pages || 1);
        setTotal(response.meta.total || 0);
      }
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load rollouts');
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: RolloutStatus): string => {
    switch (status) {
      case RolloutStatus.PENDING:
        return 'yellow';
      case RolloutStatus.IN_PROGRESS:
        return 'blue';
      case RolloutStatus.COMPLETED:
        return 'green';
      case RolloutStatus.FAILED:
        return 'red';
      case RolloutStatus.CANCELLED:
        return 'gray';
      default:
        return 'gray';
    }
  };

  const getStatusLabel = (status: RolloutStatus): string => {
    return status.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase());
  };

  const handlePageChange = (newPage: number) => {
    setFilters(prev => ({ ...prev, page: newPage }));
  };

  const handleFilterChange = () => {
    setFilters(prev => ({ ...prev, page: 1 }));
  };

  const clearFilters = () => {
    setProductFilter('');
    setStatusFilter('');
    setFilters({ page: 1, limit: 25 });
  };

  const getProductName = (productId: string): string => {
    const product = products.find(p => p.product_id === productId);
    return product ? product.name : productId;
  };

  const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleString();
  };

  if (loading && rollouts.length === 0) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spinner />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Filters */}
      <Card title="Filters">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Product
            </label>
            <Select
              value={productFilter}
              onChange={(e) => {
                setProductFilter(e.target.value);
                handleFilterChange();
              }}
              options={[
                { value: '', label: 'All Products' },
                ...products.map(p => ({ value: p.product_id, label: p.name || p.product_id })),
              ]}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Status
            </label>
            <Select
              value={statusFilter}
              onChange={(e) => {
                setStatusFilter(e.target.value as RolloutStatus | '');
                handleFilterChange();
              }}
              options={[
                { value: '', label: 'All Statuses' },
                { value: RolloutStatus.PENDING, label: 'Pending' },
                { value: RolloutStatus.IN_PROGRESS, label: 'In Progress' },
                { value: RolloutStatus.COMPLETED, label: 'Completed' },
                { value: RolloutStatus.FAILED, label: 'Failed' },
                { value: RolloutStatus.CANCELLED, label: 'Cancelled' },
              ]}
            />
          </div>

          <div className="flex items-end">
            <Button variant="secondary" onClick={clearFilters}>
              Clear Filters
            </Button>
          </div>
        </div>
      </Card>

      {error && (
        <Alert variant="error" title="Error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* Rollouts Table */}
      <Card title={`Update Rollouts (${total})`}>
        {loading ? (
          <div className="flex justify-center items-center h-64">
            <Spinner />
          </div>
        ) : rollouts.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500">No rollouts found.</p>
            {(productFilter || statusFilter) && (
              <Button variant="secondary" onClick={clearFilters} className="mt-4">
                Clear Filters
              </Button>
            )}
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Product
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      From → To
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Progress
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Initiated At
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {rollouts.map((rollout) => (
                    <tr key={rollout.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {getProductName(rollout.product_id)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        <span className="text-gray-600">{rollout.from_version}</span>
                        <span className="mx-2 text-gray-400">→</span>
                        <span className="font-semibold text-blue-600">{rollout.to_version}</span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <Badge color={getStatusColor(rollout.status)}>
                          {getStatusLabel(rollout.status)}
                        </Badge>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center gap-2">
                          <div className="w-16 bg-gray-200 rounded-full h-2">
                            <div
                              className="bg-blue-600 h-2 rounded-full"
                              style={{ width: `${rollout.progress}%` }}
                            />
                          </div>
                          <span className="text-sm text-gray-600">{rollout.progress}%</span>
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {formatDate(rollout.initiated_at)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm">
                        <Button
                          variant="secondary"
                          size="sm"
                          onClick={() => navigate(`/updates/rollouts/${rollout.id}`)}
                        >
                          View Details
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="mt-4 flex items-center justify-between border-t border-gray-200 pt-4">
                <div className="text-sm text-gray-700">
                  Showing page {filters.page} of {totalPages} ({total} total)
                </div>
                <div className="flex gap-2">
                  <Button
                    variant="secondary"
                    onClick={() => handlePageChange(filters.page - 1)}
                    disabled={filters.page <= 1}
                  >
                    Previous
                  </Button>
                  <Button
                    variant="secondary"
                    onClick={() => handlePageChange(filters.page + 1)}
                    disabled={filters.page >= totalPages}
                  >
                    Next
                  </Button>
                </div>
              </div>
            )}
          </>
        )}
      </Card>
    </div>
  );
};

