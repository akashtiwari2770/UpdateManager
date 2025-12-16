import React, { useState, useEffect } from 'react';
import { UpgradePath, UpgradePathType } from '@/types';
import { Card, Badge, Alert, Spinner, Button } from '@/components/ui';
import { upgradePathsApi } from '@/services/api/upgrade-paths';

export interface UpgradePathViewerProps {
  productId: string;
  fromVersion: string;
  toVersion: string;
  onBlock?: () => void;
}

export const UpgradePathViewer: React.FC<UpgradePathViewerProps> = ({
  productId,
  fromVersion,
  toVersion,
  onBlock,
}) => {
  const [path, setPath] = useState<UpgradePath | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadPath();
  }, [productId, fromVersion, toVersion]);

  const loadPath = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await upgradePathsApi.get(productId, fromVersion, toVersion);
      setPath(data);
    } catch (err: any) {
      if (err.response?.status !== 404) {
        setError(err.response?.data?.error?.message || 'Failed to load upgrade path');
      }
    } finally {
      setLoading(false);
    }
  };

  const getPathTypeColor = (type: UpgradePathType): string => {
    switch (type) {
      case UpgradePathType.DIRECT:
        return 'green';
      case UpgradePathType.MULTI_STEP:
        return 'blue';
      case UpgradePathType.BLOCKED:
        return 'red';
      default:
        return 'gray';
    }
  };

  const getPathTypeLabel = (type: UpgradePathType): string => {
    return type.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase());
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

  if (!path) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">No upgrade path found from {fromVersion} to {toVersion}.</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Path Overview */}
      <Card>
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">Upgrade Path</h3>
            <div className="flex items-center gap-4">
              <Badge color={getPathTypeColor(path.path_type)} size="lg">
                {getPathTypeLabel(path.path_type)}
              </Badge>
              {path.is_blocked && (
                <Badge color="red" size="lg">
                  Blocked
                </Badge>
              )}
            </div>
          </div>
          {onBlock && !path.is_blocked && (
            <Button variant="secondary" onClick={onBlock}>
              Block Path
            </Button>
          )}
        </div>
      </Card>

      {/* Path Visualization */}
      <Card title="Path Steps">
        <div className="space-y-4">
          {/* From Version */}
          <div className="flex items-center">
            <div className="flex-1">
              <div className="inline-flex items-center px-4 py-2 bg-blue-100 text-blue-800 rounded-lg">
                <span className="font-medium">{path.from_version}</span>
                <span className="ml-2 text-sm">(From)</span>
              </div>
            </div>
            <div className="mx-4">
              <svg className="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
              </svg>
            </div>
            <div className="flex-1"></div>
          </div>

          {/* Intermediate Versions */}
          {path.intermediate_versions && path.intermediate_versions.length > 0 && (
            <>
              {path.intermediate_versions.map((version, index) => (
                <div key={index} className="flex items-center">
                  <div className="flex-1"></div>
                  <div className="mx-4">
                    <div className="inline-flex items-center px-4 py-2 bg-gray-100 text-gray-800 rounded-lg">
                      <span className="font-medium">{version}</span>
                      <span className="ml-2 text-sm">(Step {index + 1})</span>
                    </div>
                  </div>
                  <div className="mx-4">
                    <svg className="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                    </svg>
                  </div>
                  <div className="flex-1"></div>
                </div>
              ))}
            </>
          )}

          {/* To Version */}
          <div className="flex items-center">
            <div className="flex-1"></div>
            <div className="mx-4">
              {path.intermediate_versions && path.intermediate_versions.length > 0 && (
                <svg className="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                </svg>
              )}
            </div>
            <div className="flex-1">
              <div className="inline-flex items-center px-4 py-2 bg-green-100 text-green-800 rounded-lg">
                <span className="font-medium">{path.to_version}</span>
                <span className="ml-2 text-sm">(To)</span>
              </div>
            </div>
          </div>
        </div>
      </Card>

      {/* Blocked Information */}
      {path.is_blocked && path.block_reason && (
        <Card title="Block Information">
          <div className="bg-red-50 border border-red-200 rounded-lg p-4">
            <p className="text-sm font-medium text-red-900 mb-2">This upgrade path is blocked</p>
            <p className="text-sm text-red-700">{path.block_reason}</p>
          </div>
        </Card>
      )}

      {/* Metadata */}
      <Card title="Path Information">
        <dl className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <dt className="text-sm font-medium text-gray-500">Created At</dt>
            <dd className="mt-1 text-sm text-gray-900">
              {path.created_at
                ? new Date(path.created_at).toLocaleString()
                : 'N/A'}
            </dd>
          </div>
        </dl>
      </Card>
    </div>
  );
};

