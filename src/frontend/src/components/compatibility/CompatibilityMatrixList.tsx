import React, { useState, useEffect } from 'react';
import { CompatibilityMatrix, ValidationStatus } from '@/types';
import { Card, Badge, Select, Input, Spinner, Alert, Button } from '@/components/ui';
import { compatibilityApi, ListCompatibilityQuery } from '@/services/api/compatibility';
import { productsApi } from '@/services/api/products';
import { Product } from '@/types';

export const CompatibilityMatrixList: React.FC = () => {
  const [matrices, setMatrices] = useState<CompatibilityMatrix[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Filters and pagination
  const [filters, setFilters] = useState<ListCompatibilityQuery>({
    page: 1,
    limit: 25,
  });
  const [productFilter, setProductFilter] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    loadProducts();
  }, []);

  useEffect(() => {
    loadMatrices();
  }, [filters, productFilter, statusFilter]);

  const loadProducts = async () => {
    try {
      const response = await productsApi.getAll({ page: 1, limit: 1000, is_active: true });
      setProducts(response.data || []);
    } catch (err) {
      console.error('Error loading products:', err);
    }
  };

  const loadMatrices = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const query: ListCompatibilityQuery = {
        page: filters.page,
        limit: filters.limit,
      };
      
      if (productFilter) {
        query.product_id = productFilter;
      }
      
      if (statusFilter) {
        query.validation_status = statusFilter;
      }
      
      const response = await compatibilityApi.list(query);
      setMatrices(response.data || []);
      
      if (response.meta) {
        setTotalPages(response.meta.total_pages || 1);
        setTotal(response.meta.total || 0);
      }
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load compatibility matrices');
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: ValidationStatus): string => {
    switch (status) {
      case 'passed':
        return 'green';
      case 'failed':
        return 'red';
      case 'pending':
        return 'yellow';
      case 'skipped':
        return 'gray';
      default:
        return 'gray';
    }
  };

  const getStatusLabel = (status: ValidationStatus): string => {
    return status.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase());
  };

  const handlePageChange = (newPage: number) => {
    setFilters(prev => ({ ...prev, page: newPage }));
  };

  const handleFilterChange = () => {
    setFilters(prev => ({ ...prev, page: 1 })); // Reset to first page when filter changes
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

  if (loading && matrices.length === 0) {
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
              Validation Status
            </label>
            <Select
              value={statusFilter}
              onChange={(e) => {
                setStatusFilter(e.target.value);
                handleFilterChange();
              }}
              options={[
                { value: '', label: 'All Statuses' },
                { value: 'passed', label: 'Passed' },
                { value: 'failed', label: 'Failed' },
                { value: 'pending', label: 'Pending' },
                { value: 'skipped', label: 'Skipped' },
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

      {/* Compatibility Matrix Table */}
      <Card title={`Compatibility Matrices (${total})`}>
        {loading ? (
          <div className="flex justify-center items-center h-64">
            <Spinner />
          </div>
        ) : matrices.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500">No compatibility matrices found.</p>
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
                      Version
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Min Server Version
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Max Server Version
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Recommended Server Version
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Validation Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Validated At
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Validated By
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {matrices.map((matrix) => (
                    <tr key={matrix.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {getProductName(matrix.product_id)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {matrix.version_number}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {matrix.min_server_version || <span className="text-gray-400">—</span>}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {matrix.max_server_version || <span className="text-gray-400">—</span>}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {matrix.recommended_server_version || <span className="text-gray-400">—</span>}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <Badge color={getStatusColor(matrix.validation_status)}>
                          {getStatusLabel(matrix.validation_status)}
                        </Badge>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {matrix.validated_at
                          ? new Date(matrix.validated_at).toLocaleString()
                          : <span className="text-gray-400">—</span>}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {matrix.validated_by || <span className="text-gray-400">—</span>}
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

