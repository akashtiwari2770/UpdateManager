import React, { useState, useEffect } from 'react';
import { CompatibilityMatrix, ValidationStatus } from '@/types';
import { Card, Badge, Alert, Spinner, Button } from '@/components/ui';
import { compatibilityApi } from '@/services/api/compatibility';

export interface CompatibilityDetailsProps {
  productId: string;
  versionNumber: string;
  onValidate?: () => void;
  refreshTrigger?: number; // Add refresh trigger to force reload
}

export const CompatibilityDetails: React.FC<CompatibilityDetailsProps> = ({
  productId,
  versionNumber,
  onValidate,
  refreshTrigger,
}) => {
  const [compatibility, setCompatibility] = useState<CompatibilityMatrix | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadCompatibility();
  }, [productId, versionNumber, refreshTrigger]);

  const loadCompatibility = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await compatibilityApi.get(productId, versionNumber);
      setCompatibility(data);
    } catch (err: any) {
      // If compatibility doesn't exist, that's okay - just show empty state
      if (err.response?.status !== 404) {
        setError(err.response?.data?.error?.message || 'Failed to load compatibility information');
      }
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
      default:
        return 'gray';
    }
  };

  const getStatusLabel = (status: ValidationStatus): string => {
    return status.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase());
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spinner />
      </div>
    );
  }

  if (error) {
    return (
      <Alert variant="error" title="Error" onClose={() => setError(null)}>
        {error}
      </Alert>
    );
  }

  if (!compatibility) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 mb-4">No compatibility information available for this version.</p>
        {onValidate && (
          <Button variant="primary" onClick={onValidate}>
            Validate Compatibility
          </Button>
        )}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Status Card */}
      <Card>
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">Validation Status</h3>
            <Badge color={getStatusColor(compatibility.validation_status)} size="lg">
              {getStatusLabel(compatibility.validation_status)}
            </Badge>
          </div>
          {onValidate && (
            <Button variant="secondary" onClick={onValidate}>
              Re-validate
            </Button>
          )}
        </div>
      </Card>

      {/* Server Version Requirements */}
      <Card title="Server Version Requirements">
        <dl className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <dt className="text-sm font-medium text-gray-500">Minimum Server Version</dt>
            <dd className="mt-1 text-sm text-gray-900">
              {compatibility.min_server_version || <span className="text-gray-400">Not specified</span>}
            </dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Maximum Server Version</dt>
            <dd className="mt-1 text-sm text-gray-900">
              {compatibility.max_server_version || <span className="text-gray-400">Not specified</span>}
            </dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Recommended Server Version</dt>
            <dd className="mt-1 text-sm text-gray-900">
              {compatibility.recommended_server_version || <span className="text-gray-400">Not specified</span>}
            </dd>
          </div>
        </dl>
      </Card>


      {/* Incompatible Versions */}
      {compatibility.incompatible_versions && compatibility.incompatible_versions.length > 0 && (
        <Card title="Incompatible Versions">
          <div className="bg-red-50 border border-red-200 rounded-lg p-4">
            <p className="text-sm font-medium text-red-900 mb-2">
              The following versions are incompatible:
            </p>
            <ul className="list-disc list-inside space-y-1">
              {compatibility.incompatible_versions.map((version, index) => (
                <li key={index} className="text-sm text-red-700">{version}</li>
              ))}
            </ul>
          </div>
        </Card>
      )}

      {/* Validation Errors */}
      {compatibility.validation_errors && compatibility.validation_errors.length > 0 && (
        <Card title="Validation Errors">
          <div className="bg-red-50 border border-red-200 rounded-lg p-4">
            <ul className="list-disc list-inside space-y-1">
              {compatibility.validation_errors.map((error, index) => (
                <li key={index} className="text-sm text-red-700">{error}</li>
              ))}
            </ul>
          </div>
        </Card>
      )}

      {/* Validation Metadata */}
      <Card title="Validation Information">
        <dl className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <dt className="text-sm font-medium text-gray-500">Validated At</dt>
            <dd className="mt-1 text-sm text-gray-900">
              {compatibility.validated_at
                ? new Date(compatibility.validated_at).toLocaleString()
                : 'Not validated'}
            </dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Validated By</dt>
            <dd className="mt-1 text-sm text-gray-900">
              {compatibility.validated_by || 'N/A'}
            </dd>
          </div>
        </dl>
      </Card>
    </div>
  );
};

