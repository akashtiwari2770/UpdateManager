import React, { useState, useEffect } from 'react';
import { UpdateDetection, Product, ReleaseType } from '@/types';
import { Card, Badge, Button, Spinner, Alert, Select, Input } from '@/components/ui';
import { updateDetectionsApi } from '@/services/api/update-detections';
import { productsApi } from '@/services/api/products';
import { useNavigate } from 'react-router-dom';

export const UpdatesDashboard: React.FC = () => {
  const navigate = useNavigate();
  const [detections, setDetections] = useState<UpdateDetection[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [productFilter, setProductFilter] = useState('');
  const [lastRefresh, setLastRefresh] = useState<Date>(new Date());

  useEffect(() => {
    loadProducts();
    loadDetections();
  }, [productFilter]);

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
      
      const query: any = {
        page: 1,
        limit: 100,
      };
      
      if (productFilter) {
        query.product_id = productFilter;
      }
      
      const response = await updateDetectionsApi.list(query);
      setDetections(response.data || []);
      setLastRefresh(new Date());
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load available updates');
    } finally {
      setLoading(false);
    }
  };

  const getProductName = (productId: string): string => {
    const product = products.find(p => p.product_id === productId);
    return product ? product.name : productId;
  };

  const getProduct = (productId: string): Product | undefined => {
    return products.find(p => p.product_id === productId);
  };

  const handleStartUpdate = (detection: UpdateDetection) => {
    // Navigate to initiate rollout with pre-filled data
    const params = new URLSearchParams({
      endpoint_id: detection.endpoint_id,
      product_id: detection.product_id,
      from_version: detection.current_version,
      to_version: detection.available_version,
    });
    navigate(`/updates/rollout/new?${params.toString()}`);
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
      {/* Header with filters and refresh */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold text-gray-900">Available Updates</h2>
          <p className="text-sm text-gray-500 mt-1">
            Last refreshed: {lastRefresh.toLocaleTimeString()}
          </p>
        </div>
        <div className="flex items-center gap-4">
          <div className="w-48">
            <Select
              value={productFilter}
              onChange={(e) => setProductFilter(e.target.value)}
              options={[
                { value: '', label: 'All Products' },
                ...products.map(p => ({ value: p.product_id, label: p.name || p.product_id })),
              ]}
            />
          </div>
          <Button variant="secondary" onClick={loadDetections}>
            Refresh
          </Button>
        </div>
      </div>

      {error && (
        <Alert variant="error" title="Error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* Update Cards */}
      {detections.length === 0 ? (
        <Card>
          <div className="text-center py-12">
            <p className="text-gray-500">No available updates found.</p>
            {productFilter && (
              <Button variant="secondary" onClick={() => setProductFilter('')} className="mt-4">
                Clear Filter
              </Button>
            )}
          </div>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {detections.map((detection) => {
            const product = getProduct(detection.product_id);
            return (
              <Card key={detection.id} className="hover:shadow-lg transition-shadow">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex-1">
                    <h3 className="text-lg font-semibold text-gray-900">
                      {getProductName(detection.product_id)}
                    </h3>
                    <p className="text-sm text-gray-500 mt-1">
                      Endpoint: {detection.endpoint_id}
                    </p>
                  </div>
                  <div className="w-3 h-3 bg-green-500 rounded-full flex-shrink-0 mt-1" title="Update Available" />
                </div>

                <div className="space-y-3 mb-4">
                  <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div>
                      <p className="text-xs text-gray-500">Current Version</p>
                      <p className="text-sm font-medium text-gray-900">{detection.current_version}</p>
                    </div>
                    <div className="text-gray-400">→</div>
                    <div className="text-right">
                      <p className="text-xs text-gray-500">Available Version</p>
                      <p className="text-sm font-semibold text-blue-600">{detection.available_version}</p>
                    </div>
                  </div>

                  <div className="text-xs text-gray-500">
                    Detected: {formatDate(detection.detected_at)}
                  </div>
                </div>

                <div className="flex gap-2 pt-4 border-t">
                  <Button
                    variant="primary"
                    onClick={() => handleStartUpdate(detection)}
                    className="flex-1"
                  >
                    Start Update
                  </Button>
                  <Button
                    variant="secondary"
                    onClick={() => navigate(`/products/${detection.product_id}`)}
                  >
                    View Product
                  </Button>
                </div>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
};

