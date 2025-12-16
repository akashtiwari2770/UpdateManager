import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { deploymentsApi } from '@/services/api/deployments';
import { productsApi } from '@/services/api/products';
import { versionsApi } from '@/services/api/versions';
import {
  CreateDeploymentRequest,
  UpdateDeploymentRequest,
  DeploymentType,
  DeploymentStatus,
  Product,
  Version,
} from '@/types';
import { Button, Card, Input, Select, Spinner } from '@/components/ui';

export const DeploymentForm: React.FC = () => {
  const navigate = useNavigate();
  const params = useParams<{
    customerId: string;
    tenantId: string;
    deploymentId?: string;
  }>();
  const { customerId, tenantId, deploymentId } = params;
  const isEdit = !!deploymentId;

  const [loading, setLoading] = useState(false);
  const [loadingData, setLoadingData] = useState(isEdit);
  const [error, setError] = useState<string | null>(null);
  const [products, setProducts] = useState<Product[]>([]);
  const [versions, setVersions] = useState<Version[]>([]);
  const [formData, setFormData] = useState<CreateDeploymentRequest>({
    product_id: '',
    deployment_type: DeploymentType.UAT,
    installed_version: '',
    status: DeploymentStatus.ACTIVE,
  });
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    loadProducts();
    if (isEdit && deploymentId && customerId && tenantId) {
      loadDeployment();
    }
  }, [isEdit, deploymentId, customerId, tenantId]);

  useEffect(() => {
    if (formData.product_id) {
      loadVersions();
    } else {
      setVersions([]);
      setFormData((prev) => ({ ...prev, installed_version: '' }));
    }
  }, [formData.product_id]);

  const loadProducts = async () => {
    try {
      const response = await productsApi.getAll({ is_active: true, limit: 1000 });
      setProducts(response?.data || []);
    } catch (err: any) {
      console.error('Error loading products:', err);
    }
  };

  const loadVersions = async () => {
    if (!formData.product_id) return;
    try {
      const versionsData = await versionsApi.getByProduct(formData.product_id);
      setVersions(Array.isArray(versionsData) ? versionsData : []);
    } catch (err: any) {
      console.error('Error loading versions:', err);
      setVersions([]);
    }
  };

  const loadDeployment = async () => {
    if (!customerId || !tenantId || !deploymentId) return;
    try {
      setLoadingData(true);
      const deployment = await deploymentsApi.getById(customerId, tenantId, deploymentId);
      setFormData({
        product_id: deployment.product_id,
        deployment_type: deployment.deployment_type,
        installed_version: deployment.installed_version,
        number_of_users: deployment.number_of_users,
        license_info: deployment.license_info || '',
        server_hostname: deployment.server_hostname || '',
        environment_details: deployment.environment_details || '',
        status: deployment.status,
      });
      // Load versions for the product
      await loadVersions();
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load deployment');
    } finally {
      setLoadingData(false);
    }
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.product_id.trim()) {
      errors.product_id = 'Product is required';
    }
    if (!formData.installed_version.trim()) {
      errors.installed_version = 'Installed version is required';
    }
    if (!formData.deployment_type) {
      errors.deployment_type = 'Deployment type is required';
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate() || !customerId || !tenantId) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      if (isEdit && deploymentId) {
        const updateData: UpdateDeploymentRequest = {
          deployment_type: formData.deployment_type,
          installed_version: formData.installed_version,
          number_of_users: formData.number_of_users,
          license_info: formData.license_info,
          server_hostname: formData.server_hostname,
          environment_details: formData.environment_details,
          status: formData.status,
        };
        await deploymentsApi.update(customerId, tenantId, deploymentId, updateData);
      } else {
        await deploymentsApi.create(customerId, tenantId, formData);
      }

      navigate(`/customers/${customerId}/tenants/${tenantId}`);
    } catch (err: any) {
      console.error('Deployment form error:', err);
      let errorMessage = isEdit ? 'Failed to update deployment' : 'Failed to create deployment';

      if (err.response?.data) {
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

      if (errorMessage.toLowerCase().includes('duplicate')) {
        setValidationErrors({
          deployment_type: 'A deployment with this product and type already exists',
        });
      }
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field: keyof CreateDeploymentRequest, value: any) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    if (validationErrors[field]) {
      setValidationErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  if (loadingData) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  if (!customerId || !tenantId) {
    return (
      <div className="space-y-6">
        <Card>
          <div className="text-center py-12">
            <p className="text-red-600 mb-4">Customer ID and Tenant ID are required</p>
            <Button variant="primary" onClick={() => navigate('/customers')}>
              Back to Customers
            </Button>
          </div>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">
          {isEdit ? 'Edit Deployment' : 'Create Deployment'}
        </h1>
        <Button
          variant="ghost"
          onClick={() => navigate(`/customers/${customerId}/tenants/${tenantId}`)}
        >
          Cancel
        </Button>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      <Card>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Select
              label="Product *"
              value={formData.product_id}
              onChange={(e) => handleChange('product_id', e.target.value)}
              error={validationErrors.product_id}
              required
              disabled={isEdit}
              options={[
                { value: '', label: 'Select a product' },
                ...products.map((p) => ({ value: p.product_id, label: p.name })),
              ]}
            />

            <Select
              label="Deployment Type *"
              value={formData.deployment_type}
              onChange={(e) => handleChange('deployment_type', e.target.value as DeploymentType)}
              error={validationErrors.deployment_type}
              required
              options={[
                { value: DeploymentType.UAT, label: 'UAT' },
                { value: DeploymentType.TESTING, label: 'Testing' },
                { value: DeploymentType.PRODUCTION, label: 'Production' },
              ]}
            />

            <Select
              label="Installed Version *"
              value={formData.installed_version}
              onChange={(e) => handleChange('installed_version', e.target.value)}
              error={validationErrors.installed_version}
              required
              disabled={!formData.product_id || versions.length === 0}
              options={[
                { value: '', label: versions.length === 0 ? 'No versions available' : 'Select version' },
                ...versions.map((v) => ({
                  value: v.version_number,
                  label: `${v.version_number} (${v.state})`,
                })),
              ]}
            />

            <Input
              label="Number of Users"
              type="number"
              value={formData.number_of_users?.toString() || ''}
              onChange={(e) =>
                handleChange(
                  'number_of_users',
                  e.target.value ? parseInt(e.target.value, 10) : undefined
                )
              }
              placeholder="Optional"
            />

            <Select
              label="Status *"
              value={formData.status}
              onChange={(e) => handleChange('status', e.target.value as DeploymentStatus)}
              required
              options={[
                { value: DeploymentStatus.ACTIVE, label: 'Active' },
                { value: DeploymentStatus.INACTIVE, label: 'Inactive' },
              ]}
            />

            <Input
              label="Server Hostname"
              value={formData.server_hostname || ''}
              onChange={(e) => handleChange('server_hostname', e.target.value)}
              placeholder="Optional"
            />

            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                License Information
              </label>
              <textarea
                value={formData.license_info || ''}
                onChange={(e) => handleChange('license_info', e.target.value)}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Optional license information"
              />
            </div>

            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Environment Details
              </label>
              <textarea
                value={formData.environment_details || ''}
                onChange={(e) => handleChange('environment_details', e.target.value)}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Optional environment details"
              />
            </div>
          </div>

          <div className="flex justify-end space-x-4 pt-4 border-t">
            <Button
              type="button"
              variant="ghost"
              onClick={() => navigate(`/customers/${customerId}/tenants/${tenantId}`)}
              disabled={loading}
            >
              Cancel
            </Button>
            <Button type="submit" variant="primary" disabled={loading} isLoading={loading}>
              {isEdit ? 'Update Deployment' : 'Create Deployment'}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
};

