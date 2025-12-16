import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { Button, Card, Input, Select, Alert, Spinner } from '@/components/ui';
import { updateRolloutsApi, InitiateRolloutRequest } from '@/services/api/update-rollouts';
import { upgradePathsApi } from '@/services/api/upgrade-paths';
import { productsApi } from '@/services/api/products';
import { versionsApi } from '@/services/api/versions';
import { UpdateRollout, Product, Version, VersionState, UpgradePathType } from '@/types';

export interface InitiateRolloutFormProps {
  onSuccess: (rollout: UpdateRollout) => void;
  onCancel?: () => void;
}

export const InitiateRolloutForm: React.FC<InitiateRolloutFormProps> = ({
  onSuccess,
  onCancel,
}) => {
  const [searchParams] = useSearchParams();
  const [endpointId, setEndpointId] = useState(searchParams.get('endpoint_id') || '');
  const [productId, setProductId] = useState(searchParams.get('product_id') || '');
  const [fromVersion, setFromVersion] = useState(searchParams.get('from_version') || '');
  const [toVersion, setToVersion] = useState(searchParams.get('to_version') || '');
  const [loading, setLoading] = useState(false);
  const [checkingPath, setCheckingPath] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});
  const [upgradePathValid, setUpgradePathValid] = useState<boolean | null>(null);
  const [upgradePathBlocked, setUpgradePathBlocked] = useState(false);
  
  // Data for dropdowns
  const [products, setProducts] = useState<Product[]>([]);
  const [versions, setVersions] = useState<Version[]>([]);
  const [loadingProducts, setLoadingProducts] = useState(true);
  const [loadingVersions, setLoadingVersions] = useState(false);

  useEffect(() => {
    loadProducts();
  }, []);

  useEffect(() => {
    if (productId) {
      loadVersions(productId);
    } else {
      setVersions([]);
      setFromVersion('');
      setToVersion('');
    }
  }, [productId]);

  useEffect(() => {
    // Check upgrade path when versions are provided
    if (productId && fromVersion && toVersion) {
      checkUpgradePath();
    }
  }, [productId, fromVersion, toVersion]);

  const loadProducts = async () => {
    try {
      setLoadingProducts(true);
      const response = await productsApi.getAll({ page: 1, limit: 1000, is_active: true });
      setProducts(response.data || []);
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
    // Only show released versions for to version
    return versions.filter(v => v.state === VersionState.RELEASED);
  };

  const getFromVersions = (): Version[] => {
    // Show all versions for from version (any state)
    return versions;
  };

  const checkUpgradePath = async () => {
    if (!productId || !fromVersion || !toVersion) return;

    try {
      setCheckingPath(true);
      setUpgradePathValid(null);
      setUpgradePathBlocked(false);
      setError(null);

      const path = await upgradePathsApi.get(productId, fromVersion, toVersion);
      
      if (path.is_blocked) {
        setUpgradePathBlocked(true);
        setUpgradePathValid(false);
        setError(`Upgrade path is blocked: ${path.block_reason || 'No reason provided'}`);
      } else {
        setUpgradePathValid(true);
        setError(null);
      }
    } catch (err: any) {
      if (err.response?.status === 404) {
        // Upgrade path doesn't exist - try to create it automatically
        try {
          await upgradePathsApi.create(productId, {
            from_version: fromVersion,
            to_version: toVersion,
            path_type: UpgradePathType.DIRECT,
            is_blocked: false,
          });
          // Successfully created, verify it
          setUpgradePathValid(true);
          setError(null);
        } catch (createErr: any) {
          // Failed to create - show warning but allow proceeding
          setUpgradePathValid(null);
          setUpgradePathBlocked(false);
          setError(null); // Clear error, we'll show a warning instead
        }
      } else {
        setUpgradePathValid(null);
        setUpgradePathBlocked(false);
        setError('Failed to check upgrade path. You can still proceed, but please verify manually.');
      }
    } finally {
      setCheckingPath(false);
    }
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!endpointId.trim()) {
      errors.endpointId = 'Endpoint ID is required';
    }
    if (!productId.trim()) {
      errors.productId = 'Product is required';
    } else {
      // Validate product exists
      const productExists = products.some(p => p.product_id === productId);
      if (!productExists) {
        errors.productId = 'Selected product does not exist';
      }
    }
    if (!fromVersion.trim()) {
      errors.fromVersion = 'From version is required';
    } else {
      // Validate from version exists for the product
      const versionExists = versions.some(v => v.version_number === fromVersion);
      if (!versionExists) {
        errors.fromVersion = 'Selected version does not exist for this product';
      }
    }
    if (!toVersion.trim()) {
      errors.toVersion = 'To version is required';
    } else {
      // Validate to version exists and is released
      const toVersionObj = versions.find(v => v.version_number === toVersion);
      if (!toVersionObj) {
        errors.toVersion = 'Selected version does not exist for this product';
      } else if (toVersionObj.state !== VersionState.RELEASED) {
        errors.toVersion = 'To version must be in Released state';
      }
    }
    if (fromVersion && toVersion && fromVersion === toVersion) {
      errors.toVersion = 'To version must be different from from version';
    }

    // Note: upgrade path validation is optional - we don't block if it's not found
    // Only block if it's explicitly blocked
    if (upgradePathBlocked) {
      errors.upgradePath = 'Upgrade path is blocked';
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) {
      return;
    }

    // Only block if upgrade path is explicitly blocked
    if (upgradePathBlocked) {
      setError('Cannot proceed: Upgrade path is blocked.');
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const request: InitiateRolloutRequest = {
        endpoint_id: endpointId.trim(),
        product_id: productId.trim(),
        from_version: fromVersion.trim(),
        to_version: toVersion.trim(),
      };

      const rollout = await updateRolloutsApi.initiate(request);
      onSuccess(rollout);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to initiate rollout');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card title="Initiate Update Rollout">
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <Alert 
            variant={upgradePathBlocked ? "error" : "error"} 
            title="Error" 
            onClose={() => setError(null)}
          >
            {error}
          </Alert>
        )}

        {upgradePathValid === true && (
          <Alert variant="success" title="Upgrade Path Valid">
            Upgrade path verified. You can proceed with the rollout.
          </Alert>
        )}
        {upgradePathValid === null && fromVersion && toVersion && !checkingPath && !upgradePathBlocked && (
          <Alert variant="warning" title="Upgrade Path Not Found">
            No upgrade path found for {fromVersion} → {toVersion}. A direct upgrade path was attempted to be created automatically. You can still proceed with the rollout.
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
              value={productId}
              onChange={(e) => setProductId(e.target.value)}
              error={validationErrors.productId}
              required
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
              From Version <span className="text-red-500">*</span>
            </label>
            {loadingVersions ? (
              <div className="flex items-center gap-2">
                <Spinner size="sm" />
                <span className="text-sm text-gray-500">Loading versions...</span>
              </div>
            ) : (
              <Select
                value={fromVersion}
                onChange={(e) => setFromVersion(e.target.value)}
                error={validationErrors.fromVersion}
                required
                disabled={!productId}
                options={[
                  { value: '', label: 'Select from version...' },
                  ...getFromVersions().map(v => ({ 
                    value: v.version_number, 
                    label: `${v.version_number} (${v.state})` 
                  })),
                ]}
              />
            )}
            {validationErrors.fromVersion && (
              <p className="mt-1 text-sm text-red-600">{validationErrors.fromVersion}</p>
            )}
            {!productId && (
              <p className="mt-1 text-xs text-gray-500">Please select a product first</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              To Version <span className="text-red-500">*</span>
            </label>
            {loadingVersions ? (
              <div className="flex items-center gap-2">
                <Spinner size="sm" />
                <span className="text-sm text-gray-500">Loading versions...</span>
              </div>
            ) : (
              <>
                <Select
                  value={toVersion}
                  onChange={(e) => setToVersion(e.target.value)}
                  error={validationErrors.toVersion}
                  required
                  disabled={!productId}
                  options={[
                    { value: '', label: 'Select to version...' },
                    ...getAvailableVersions().map(v => ({ 
                      value: v.version_number, 
                      label: `${v.version_number} (Released)` 
                    })),
                  ]}
                />
                {fromVersion && toVersion && (
                  <Button
                    type="button"
                    variant="secondary"
                    size="sm"
                    onClick={checkUpgradePath}
                    isLoading={checkingPath}
                    className="mt-2"
                  >
                    Verify Upgrade Path
                  </Button>
                )}
              </>
            )}
            {validationErrors.toVersion && (
              <p className="mt-1 text-sm text-red-600">{validationErrors.toVersion}</p>
            )}
            {!productId && (
              <p className="mt-1 text-xs text-gray-500">Please select a product first</p>
            )}
            {productId && getAvailableVersions().length === 0 && (
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
          <Button 
            variant="primary" 
            type="submit" 
            isLoading={loading}
            disabled={upgradePathBlocked || checkingPath}
          >
            Initiate Rollout
          </Button>
        </div>
      </form>
    </Card>
  );
};

