import React, { useState, useEffect } from 'react';
import { Button, Card, Input, Select, Alert, Spinner } from '@/components/ui';
import { updateDetectionsApi } from '@/services/api/update-detections';
import { productsApi } from '@/services/api/products';
import { versionsApi } from '@/services/api/versions';
import { UpdateDetection, Product, Version, VersionState } from '@/types';

export interface UpdateDetectionFormProps {
  productId?: string;
  onSuccess: (detection: UpdateDetection) => void;
  onCancel?: () => void;
}

export const UpdateDetectionForm: React.FC<UpdateDetectionFormProps> = ({
  productId,
  onSuccess,
  onCancel,
}) => {
  const [endpointId, setEndpointId] = useState('');
  const [selectedProductId, setSelectedProductId] = useState(productId || '');
  const [currentVersion, setCurrentVersion] = useState('');
  const [availableVersion, setAvailableVersion] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});
  
  // Data for dropdowns
  const [products, setProducts] = useState<Product[]>([]);
  const [versions, setVersions] = useState<Version[]>([]);
  const [loadingProducts, setLoadingProducts] = useState(true);
  const [loadingVersions, setLoadingVersions] = useState(false);

  useEffect(() => {
    loadProducts();
  }, []);

  useEffect(() => {
    if (selectedProductId) {
      loadVersions(selectedProductId);
    } else {
      setVersions([]);
      setCurrentVersion('');
      setAvailableVersion('');
    }
  }, [selectedProductId]);

  const loadProducts = async () => {
    try {
      setLoadingProducts(true);
      const response = await productsApi.getAll({ page: 1, limit: 1000, is_active: true });
      setProducts(response.data || []);
      
      // If productId prop is provided, ensure it's in the list
      if (productId && !response.data?.find(p => p.product_id === productId)) {
        // Try to get the product directly
        try {
          const product = await productsApi.getById(productId);
          setProducts(prev => [...prev, product]);
        } catch (err) {
          console.error('Failed to load product:', err);
        }
      }
    } catch (err) {
      setError('Failed to load products');
    } finally {
      setLoadingProducts(false);
    }
  };

  const loadVersions = async (productId: string) => {
    try {
      setLoadingVersions(true);
      const productVersions = await versionsApi.getByProduct(productId);
      setVersions(productVersions || []);
    } catch (err) {
      setError('Failed to load versions');
      setVersions([]);
    } finally {
      setLoadingVersions(false);
    }
  };

  const getAvailableVersions = (): Version[] => {
    // Only show released versions for available version
    return versions.filter(v => v.state === VersionState.RELEASED);
  };

  const getCurrentVersions = (): Version[] => {
    // Show all versions for current version (any state)
    return versions;
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!endpointId.trim()) {
      errors.endpointId = 'Endpoint ID is required';
    }
    if (!selectedProductId.trim()) {
      errors.productId = 'Product is required';
    } else {
      // Validate product exists
      const productExists = products.some(p => p.product_id === selectedProductId);
      if (!productExists) {
        errors.productId = 'Selected product does not exist';
      }
    }
    if (!currentVersion.trim()) {
      errors.currentVersion = 'Current version is required';
    } else {
      // Validate current version exists for the product
      const versionExists = versions.some(v => v.version_number === currentVersion);
      if (!versionExists) {
        errors.currentVersion = 'Selected version does not exist for this product';
      }
    }
    if (!availableVersion.trim()) {
      errors.availableVersion = 'Available version is required';
    } else {
      // Validate available version exists and is released
      const availableVersionObj = versions.find(v => v.version_number === availableVersion);
      if (!availableVersionObj) {
        errors.availableVersion = 'Selected version does not exist for this product';
      } else if (availableVersionObj.state !== VersionState.RELEASED) {
        errors.availableVersion = 'Available version must be in Released state';
      }
    }
    if (currentVersion && availableVersion && currentVersion === availableVersion) {
      errors.availableVersion = 'Available version must be different from current version';
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

      const detection = await updateDetectionsApi.detect({
        endpoint_id: endpointId.trim(),
        product_id: selectedProductId.trim(),
        current_version: currentVersion.trim(),
        available_version: availableVersion.trim(),
      });

      onSuccess(detection);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to register update detection');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card title="Register Update Detection">
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <Alert variant="error" title="Error" onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        <Input
          label="Endpoint ID"
          value={endpointId}
          onChange={(e) => setEndpointId(e.target.value)}
          placeholder="e.g., endpoint-123"
          error={validationErrors.endpointId}
          required
        />

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Product <span className="text-red-500">*</span>
          </label>
          {loadingProducts ? (
            <div className="flex items-center gap-2">
              <Spinner size="sm" />
              <span className="text-sm text-gray-500">Loading products...</span>
            </div>
          ) : (
            <Select
              value={selectedProductId}
              onChange={(e) => setSelectedProductId(e.target.value)}
              error={validationErrors.productId}
              required
              disabled={!!productId}
              options={[
                { value: '', label: 'Select a product...' },
                ...products.map(p => ({ 
                  value: p.product_id, 
                  label: `${p.name || p.product_id} (${p.product_id})` 
                })),
              ]}
            />
          )}
          {validationErrors.productId && (
            <p className="mt-1 text-sm text-red-600">{validationErrors.productId}</p>
          )}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Current Version <span className="text-red-500">*</span>
            </label>
            {loadingVersions ? (
              <div className="flex items-center gap-2">
                <Spinner size="sm" />
                <span className="text-sm text-gray-500">Loading versions...</span>
              </div>
            ) : (
              <Select
                value={currentVersion}
                onChange={(e) => setCurrentVersion(e.target.value)}
                error={validationErrors.currentVersion}
                required
                disabled={!selectedProductId}
                options={[
                  { value: '', label: 'Select current version...' },
                  ...getCurrentVersions().map(v => ({ 
                    value: v.version_number, 
                    label: `${v.version_number} (${v.state})` 
                  })),
                ]}
              />
            )}
            {validationErrors.currentVersion && (
              <p className="mt-1 text-sm text-red-600">{validationErrors.currentVersion}</p>
            )}
            {!selectedProductId && (
              <p className="mt-1 text-xs text-gray-500">Please select a product first</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Available Version <span className="text-red-500">*</span>
            </label>
            {loadingVersions ? (
              <div className="flex items-center gap-2">
                <Spinner size="sm" />
                <span className="text-sm text-gray-500">Loading versions...</span>
              </div>
            ) : (
              <Select
                value={availableVersion}
                onChange={(e) => setAvailableVersion(e.target.value)}
                error={validationErrors.availableVersion}
                required
                disabled={!selectedProductId}
                options={[
                  { value: '', label: 'Select available version...' },
                  ...getAvailableVersions().map(v => ({ 
                    value: v.version_number, 
                    label: `${v.version_number} (Released)` 
                  })),
                ]}
              />
            )}
            {validationErrors.availableVersion && (
              <p className="mt-1 text-sm text-red-600">{validationErrors.availableVersion}</p>
            )}
            {!selectedProductId && (
              <p className="mt-1 text-xs text-gray-500">Please select a product first</p>
            )}
            {selectedProductId && getAvailableVersions().length === 0 && (
              <p className="mt-1 text-xs text-yellow-600">No released versions available for this product</p>
            )}
          </div>
        </div>

        <div className="flex items-center justify-end gap-4 pt-4 border-t">
          {onCancel && (
            <Button variant="secondary" onClick={onCancel} disabled={loading}>
              Cancel
            </Button>
          )}
          <Button variant="primary" type="submit" isLoading={loading}>
            Register Detection
          </Button>
        </div>
      </form>
    </Card>
  );
};

