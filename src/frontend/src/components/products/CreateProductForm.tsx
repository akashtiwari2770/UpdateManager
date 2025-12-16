import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { productsApi } from '@/services/api/products';
import { CreateProductRequest, ProductType } from '@/types';
import { Button, Card, Input, Select, Alert, Spinner } from '@/components/ui';

export const CreateProductForm: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState<CreateProductRequest>({
    product_id: '',
    name: '',
    type: ProductType.SERVER,
    description: '',
    vendor: '',
  });
  const [validationErrors, setValidationErrors] = useState<
    Record<string, string>
  >({});

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.product_id.trim()) {
      errors.product_id = 'Product ID is required';
    } else if (!/^[a-zA-Z0-9_-]+$/.test(formData.product_id)) {
      errors.product_id =
        'Product ID can only contain letters, numbers, underscores, and hyphens';
    }

    if (!formData.name.trim()) {
      errors.name = 'Product Name is required';
    }

    if (!formData.type) {
      errors.type = 'Product Type is required';
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) {
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const product = await productsApi.create(formData);
      
      // Validate product was created successfully
      if (!product || !product.id) {
        throw new Error('Product was created but ID is missing. Please refresh the products list.');
      }
      
      navigate(`/products/${product.id}`);
    } catch (err: any) {
      console.error('Create product error:', err);
      let errorMessage = 'Failed to create product';
      
      if (err.response?.data) {
        // Check if error is in the unwrapped format
        if (err.response.data.error) {
          errorMessage = err.response.data.error.message || err.response.data.error;
        } else if (typeof err.response.data === 'string') {
          errorMessage = err.response.data;
        } else if (err.response.data.message) {
          errorMessage = err.response.data.message;
        }
      } else if (err.message) {
        errorMessage = err.message;
      }
      
      setError(errorMessage);

      // Handle duplicate product ID error
      if (errorMessage.toLowerCase().includes('already exists') || 
          errorMessage.toLowerCase().includes('duplicate')) {
        setValidationErrors({
          product_id: 'A product with this ID already exists',
        });
      }
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (
    field: keyof CreateProductRequest,
    value: string | ProductType
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    // Clear validation error for this field
    if (validationErrors[field]) {
      setValidationErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Create Product</h1>
        <Button variant="ghost" onClick={() => navigate('/products')}>
          Cancel
        </Button>
      </div>

      <Card title="Product Details">
        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <Alert variant="error" title="Error" onClose={() => setError(null)}>
              {error}
            </Alert>
          )}

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Input
              label="Product ID"
              required
              value={formData.product_id}
              onChange={(e) => handleChange('product_id', e.target.value)}
              error={validationErrors.product_id}
              placeholder="e.g., my-product-v1"
              helperText="Unique identifier for the product. Only letters, numbers, underscores, and hyphens allowed."
            />

            <Input
              label="Product Name"
              required
              value={formData.name}
              onChange={(e) => handleChange('name', e.target.value)}
              error={validationErrors.name}
              placeholder="e.g., My Awesome Product"
            />
          </div>

          <div>
            <Select
              label="Product Type"
              required
              value={formData.type}
              onChange={(e) => handleChange('type', e.target.value as ProductType)}
              error={validationErrors.type}
              options={[
                { value: ProductType.SERVER, label: 'Server' },
                { value: ProductType.CLIENT, label: 'Client' },
              ]}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Description
            </label>
            <textarea
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              rows={4}
              value={formData.description || ''}
              onChange={(e) => handleChange('description', e.target.value)}
              placeholder="Optional description of the product..."
            />
          </div>

          <div className="flex items-center justify-end gap-4 pt-4 border-t">
            <Button
              type="button"
              variant="secondary"
              onClick={() => navigate('/products')}
              disabled={loading}
            >
              Cancel
            </Button>
            <Button type="submit" variant="primary" isLoading={loading}>
              Create Product
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
};

