import React, { useState, useEffect } from 'react';
import { UpdateDetection } from '@/types';
import { Card, Input, Select, Spinner, Alert, Button } from '@/components/ui';
import { updateDetectionsApi, ListUpdateDetectionsQuery } from '@/services/api/update-detections';
import { productsApi } from '@/services/api/products';
import { Product } from '@/types';

export const UpdateDetectionList: React.FC = () => {
  const [detections, setDetections] = useState<UpdateDetection[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Filters and pagination
  const [filters, setFilters] = useState<ListUpdateDetectionsQuery>({
    page: 1,
    limit: 25,
  });
  const [endpointFilter, setEndpointFilter] = useState('');
  const [productFilter, setProductFilter] = useState('');
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    loadProducts();
  }, []);

  useEffect(() => {
    loadDetections();
  }, [filters, endpointFilter, productFilter]);

  const loadProducts = async () => {
    try {
      const response = await productsApi.getAll({ page: 1, limit: 1000, is_active: true });
      setProducts(response.data || []);
    } catch (err) {
      console.error('Error loading products:', err);
    }
  };

  const loadDetections = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const query: ListUpdateDetectionsQuery = {
        page: filters.page,
        limit: filters.limit,
      };
      
      if (endpointFilter) {
        query.endpoint_id = endpointFilter;
      }
      
      if (productFilter) {
        query.product_id = productFilter;
      }
      
      const response = await updateDetectionsApi.list(query);
      setDetections(response.data || []);
      
      if (response.meta) {
        setTotalPages(response.meta.total_pages || 1);
        setTotal(response.meta.total || 0);
      }
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load update detections');
    } finally {
      setLoading(false);
    }
  };

  const handlePageChange = (newPage: number) => {
    setFilters(prev => ({ ...prev, page: newPage }));
  };

  const handleFilterChange = () => {
    setFilters(prev => ({ ...prev, page: 1 }));
  };

  const clearFilters = () => {
    setEndpointFilter('');
    setProductFilter('');
    setFilters({ page: 1, limit: 25 });
  };

  const getProductName = (productId: string): string => {
    const product = products.find(p => p.product_id === productId);
    return product ? product.name : productId;
  };

  const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleString();
  };

  if (loading && detections.length === 0) {
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
              Endpoint ID
            </label>
            <Input
              value={endpointFilter}
              onChange={(e) => {
                setEndpointFilter(e.target.value);
                handleFilterChange();
              }}
              placeholder="Filter by endpoint ID..."
            />
          </div>

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

      {/* Detections Table */}
      <Card title={`Update Detections (${total})`}>
        {loading ? (
          <div className="flex justify-center items-center h-64">
            <Spinner />
          </div>
        ) : detections.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500">No update detections found.</p>
            {(endpointFilter || productFilter) && (
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
                      Endpoint ID
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Product
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Current Version
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Available Version
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Detected At
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Last Checked
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {detections.map((detection) => (
                    <tr key={detection.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {detection.endpoint_id}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {getProductName(detection.product_id)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {detection.current_version}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-600">
                        {detection.available_version}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {formatDate(detection.detected_at)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {formatDate(detection.last_checked_at)}
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

