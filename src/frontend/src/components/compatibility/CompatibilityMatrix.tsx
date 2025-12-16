import React, { useState, useEffect } from 'react';
import { CompatibilityMatrix as CompatibilityMatrixType, ValidationStatus } from '@/types';
import { Card, Badge, Select, Input, Spinner, Alert } from '@/components/ui';
import { compatibilityApi } from '@/services/api/compatibility';

export interface CompatibilityMatrixProps {
  productId?: string;
  versionNumber?: string;
}

export const CompatibilityMatrix: React.FC<CompatibilityMatrixProps> = ({
  productId,
  versionNumber,
}) => {
  const [matrix, setMatrix] = useState<CompatibilityMatrixType | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filterProductId, setFilterProductId] = useState(productId || '');
  const [filterVersion, setFilterVersion] = useState(versionNumber || '');

  useEffect(() => {
    if (productId && versionNumber) {
      loadCompatibility();
    }
  }, [productId, versionNumber]);

  const loadCompatibility = async () => {
    if (!productId || !versionNumber) return;

    try {
      setLoading(true);
      setError(null);
      const data = await compatibilityApi.get(productId, versionNumber);
      setMatrix(data);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load compatibility matrix');
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

  if (!matrix) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">No compatibility matrix available.</p>
        {productId && versionNumber && (
          <button
            onClick={loadCompatibility}
            className="mt-4 text-blue-600 hover:text-blue-700"
          >
            Load Compatibility Matrix
          </button>
        )}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Filters */}
      <Card>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Input
            label="Filter by Product ID"
            value={filterProductId}
            onChange={(e) => setFilterProductId(e.target.value)}
            placeholder="Enter product ID..."
          />
          <Input
            label="Filter by Version"
            value={filterVersion}
            onChange={(e) => setFilterVersion(e.target.value)}
            placeholder="Enter version number..."
          />
        </div>
      </Card>

      {/* Compatibility Matrix Table */}
      <Card title="Compatibility Matrix">
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
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              <tr>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {matrix.product_id || 'N/A'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {matrix.version_number || 'N/A'}
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
              </tr>
            </tbody>
          </table>
        </div>

        {/* Validation Errors */}
        {matrix.validation_errors && matrix.validation_errors.length > 0 && (
          <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
            <h4 className="text-sm font-medium text-red-900 mb-2">Validation Errors:</h4>
            <ul className="list-disc list-inside space-y-1">
              {matrix.validation_errors.map((error, index) => (
                <li key={index} className="text-sm text-red-700">{error}</li>
              ))}
            </ul>
          </div>
        )}

        {/* Validation Info */}
        <div className="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <p className="text-sm text-gray-500">Validated At</p>
            <p className="text-sm font-medium text-gray-900">
              {matrix.validated_at
                ? new Date(matrix.validated_at).toLocaleString()
                : 'Not validated'}
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-500">Validated By</p>
            <p className="text-sm font-medium text-gray-900">
              {matrix.validated_by || 'N/A'}
            </p>
          </div>
        </div>
      </Card>
    </div>
  );
};

