import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { deploymentsApi } from '@/services/api/deployments';
import { Deployment, DeploymentStatus } from '@/types';
import { Button, Card, Badge, Spinner, Alert } from '@/components/ui';
import { DeploymentTypeBadge } from './DeploymentTypeBadge';
import { DeploymentPendingUpdates } from './DeploymentPendingUpdates';

type TabType = 'overview' | 'updates';

export const DeploymentDetails: React.FC = () => {
  const { customerId, tenantId, deploymentId } = useParams<{
    customerId: string;
    tenantId: string;
    deploymentId: string;
  }>();
  const navigate = useNavigate();
  const [deployment, setDeployment] = useState<Deployment | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<TabType>('overview');

  useEffect(() => {
    if (customerId && tenantId && deploymentId) {
      loadDeployment();
    }
  }, [customerId, tenantId, deploymentId]);

  const loadDeployment = async () => {
    if (!customerId || !tenantId || !deploymentId) return;
    try {
      setLoading(true);
      setError(null);
      const data = await deploymentsApi.getById(customerId, tenantId, deploymentId);
      setDeployment(data);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load deployment');
      console.error('Error loading deployment:', err);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: DeploymentStatus) => {
    const statusConfig = {
      [DeploymentStatus.ACTIVE]: { label: 'Active', className: 'bg-green-100 text-green-800' },
      [DeploymentStatus.INACTIVE]: { label: 'Inactive', className: 'bg-gray-100 text-gray-800' },
    };

    const config = statusConfig[status] || statusConfig[DeploymentStatus.INACTIVE];
    return <Badge className={config.className}>{config.label}</Badge>;
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  if (error || !deployment) {
    return (
      <div className="space-y-4">
        {error && <Alert variant="error">{error}</Alert>}
        <Button variant="secondary" onClick={() => navigate(-1)}>
          Go Back
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{deployment.deployment_id}</h1>
          <p className="text-sm text-gray-500 mt-1">
            {deployment.product_id} • {deployment.deployment_type}
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            variant="secondary"
            onClick={() =>
              navigate(`/customers/${customerId}/tenants/${tenantId}/deployments/${deploymentId}/edit`)
            }
          >
            Edit
          </Button>
          <Button variant="secondary" onClick={() => navigate(-1)}>
            Back
          </Button>
        </div>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200">
        <nav className="-mb-px flex space-x-8">
          <button
            onClick={() => setActiveTab('overview')}
            className={`
              py-4 px-1 border-b-2 font-medium text-sm
              ${
                activeTab === 'overview'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }
            `}
          >
            Overview
          </button>
          <button
            onClick={() => setActiveTab('updates')}
            className={`
              py-4 px-1 border-b-2 font-medium text-sm
              ${
                activeTab === 'updates'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }
            `}
          >
            Pending Updates
          </button>
        </nav>
      </div>

      {/* Tab Content */}
      <div className="mt-6">
        {activeTab === 'overview' && (
          <div className="space-y-6">
            <Card className="p-6">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">Deployment Information</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="text-sm font-medium text-gray-500">Deployment ID</label>
                  <p className="mt-1 text-sm font-medium text-gray-900">{deployment.deployment_id}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Product</label>
                  <p className="mt-1 text-sm font-medium text-gray-900">{deployment.product_id}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Deployment Type</label>
                  <div className="mt-1">
                    <DeploymentTypeBadge type={deployment.deployment_type} />
                  </div>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Status</label>
                  <div className="mt-1">{getStatusBadge(deployment.status)}</div>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Installed Version</label>
                  <p className="mt-1 text-sm font-medium text-gray-900">{deployment.installed_version}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Number of Users</label>
                  <p className="mt-1 text-sm font-medium text-gray-900">
                    {deployment.number_of_users || 'Not specified'}
                  </p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">License Information</label>
                  <p className="mt-1 text-sm text-gray-900">
                    {deployment.license_info || 'Not specified'}
                  </p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Server Hostname</label>
                  <p className="mt-1 text-sm text-gray-900">
                    {deployment.server_hostname || 'Not specified'}
                  </p>
                </div>
                <div className="md:col-span-2">
                  <label className="text-sm font-medium text-gray-500">Environment Details</label>
                  <p className="mt-1 text-sm text-gray-900">
                    {deployment.environment_details || 'Not specified'}
                  </p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Deployment Date</label>
                  <p className="mt-1 text-sm text-gray-900">
                    {new Date(deployment.deployment_date).toLocaleDateString()}
                  </p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Last Updated</label>
                  <p className="mt-1 text-sm text-gray-900">
                    {new Date(deployment.last_updated_date).toLocaleDateString()}
                  </p>
                </div>
              </div>
            </Card>
          </div>
        )}

        {activeTab === 'updates' && customerId && tenantId && deploymentId && (
          <DeploymentPendingUpdates
            customerId={customerId}
            tenantId={tenantId}
            deploymentId={deploymentId}
          />
        )}
      </div>
    </div>
  );
};

