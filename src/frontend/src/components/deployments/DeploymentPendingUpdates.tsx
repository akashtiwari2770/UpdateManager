import React, { useState, useEffect } from 'react';
import { pendingUpdatesApi } from '@/services/api/pending-updates';
import { PendingUpdatesResponse, AvailableUpdate } from '@/types';
import { Card, Spinner, Badge } from '@/components/ui';
import { UpdateBadge } from '@/components/ui/UpdateBadge';
import { PriorityBadge } from '@/components/ui/PriorityBadge';
import { VersionGapBadge } from '@/components/ui/VersionGapBadge';

interface DeploymentPendingUpdatesProps {
  customerId: string;
  tenantId: string;
  deploymentId: string;
}

export const DeploymentPendingUpdates: React.FC<DeploymentPendingUpdatesProps> = ({
  customerId,
  tenantId,
  deploymentId,
}) => {
  const [pendingUpdates, setPendingUpdates] = useState<PendingUpdatesResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadPendingUpdates();
  }, [customerId, tenantId, deploymentId]);

  const loadPendingUpdates = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await pendingUpdatesApi.getDeploymentPendingUpdates(
        customerId,
        tenantId,
        deploymentId
      );
      setPendingUpdates(data);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load pending updates');
      console.error('Error loading pending updates:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-32">
        <Spinner size="md" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
        {error}
      </div>
    );
  }

  if (!pendingUpdates || pendingUpdates.update_count === 0) {
    return (
      <Card className="p-4">
        <div className="text-center py-8">
          <p className="text-gray-500">No pending updates available</p>
          <p className="text-sm text-gray-400 mt-2">
            Current version: <span className="font-medium">{pendingUpdates?.current_version}</span>
          </p>
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      {/* Summary Card */}
      <Card className="p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-medium text-gray-900">Pending Updates</h3>
          <UpdateBadge count={pendingUpdates.update_count} priority={pendingUpdates.priority} />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
          <div>
            <label className="text-sm font-medium text-gray-500">Current Version</label>
            <p className="mt-1 text-sm font-medium text-gray-900">
              {pendingUpdates.current_version}
            </p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500">Latest Version</label>
            <p className="mt-1 text-sm font-medium text-gray-900">{pendingUpdates.latest_version}</p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500">Update Type</label>
            <div className="mt-1">
              <VersionGapBadge gapType={pendingUpdates.version_gap_type} />
            </div>
          </div>
        </div>

        <div className="flex items-center gap-4">
          <div>
            <label className="text-sm font-medium text-gray-500">Priority</label>
            <div className="mt-1">
              <PriorityBadge priority={pendingUpdates.priority} />
            </div>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500">Total Updates</label>
            <p className="mt-1 text-sm font-medium text-gray-900">
              {pendingUpdates.update_count} available
            </p>
          </div>
        </div>
      </Card>

      {/* Available Updates List */}
      <Card className="p-6">
        <h4 className="text-md font-medium text-gray-900 mb-4">Available Updates</h4>
        <div className="space-y-3">
          {pendingUpdates.available_updates.map((update: AvailableUpdate, index: number) => (
            <div
              key={index}
              className="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors"
            >
              <div className="flex items-center justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="text-lg font-semibold text-gray-900">
                      {update.version_number}
                    </span>
                    {update.is_security_update && (
                      <Badge className="bg-red-100 text-red-800">Security</Badge>
                    )}
                    <Badge className="bg-gray-100 text-gray-800">{update.release_type}</Badge>
                  </div>
                  <p className="text-sm text-gray-500">
                    Released: {new Date(update.release_date).toLocaleDateString()}
                  </p>
                  {update.upgrade_path && update.upgrade_path.length > 1 && (
                    <p className="text-xs text-gray-400 mt-1">
                      Upgrade path: {update.upgrade_path.join(' → ')}
                    </p>
                  )}
                </div>
                <div className="text-right">
                  <Badge
                    className={
                      update.compatibility_status === 'compatible'
                        ? 'bg-green-100 text-green-800'
                        : 'bg-yellow-100 text-yellow-800'
                    }
                  >
                    {update.compatibility_status}
                  </Badge>
                </div>
              </div>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
};

