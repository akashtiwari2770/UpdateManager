import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { productsApi } from '@/services/api/products';
import { Product, ProductType } from '@/types';
import { Card, Badge, Spinner, Alert } from '@/components/ui';

export const ActiveProducts: React.FC = () => {
  const navigate = useNavigate();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadActiveProducts();
  }, []);

  const loadActiveProducts = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await productsApi.getActive();
      setProducts(data);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load active products');
      console.error('Error loading active products:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  if (error) {
    return (
      <Alert variant="error" title="Error">
        {error}
      </Alert>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Active Products</h1>
      </div>

      <Card>
        {products.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">No active products found</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {products.map((product) => (
              <div
                key={product.id}
                className="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow cursor-pointer"
                onClick={() => navigate(`/products/${product.id}`)}
              >
                <div className="flex items-start justify-between mb-2">
                  <h3 className="text-lg font-semibold text-gray-900">
                    {product.name}
                  </h3>
                  <Badge
                    variant={
                      product.type === ProductType.SERVER ? 'info' : 'default'
                    }
                  >
                    {product.type}
                  </Badge>
                </div>
                <p className="text-sm text-gray-500 font-mono mb-2">
                  {product.product_id}
                </p>
                {product.description && (
                  <p className="text-sm text-gray-600 line-clamp-2">
                    {product.description}
                  </p>
                )}
              </div>
            ))}
          </div>
        )}
      </Card>
    </div>
  );
};

