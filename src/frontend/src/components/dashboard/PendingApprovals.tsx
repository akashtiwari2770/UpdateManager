import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Version, VersionState } from '@/types';
import { Card, Button, Badge, Spinner, Alert } from '@/components/ui';
import { versionsApi } from '@/services/api/versions';

interface PendingApprovalsProps {
  approvals: Version[];
  loading?: boolean;
  onApprove?: (versionId: string) => void;
}

export const PendingApprovals: React.FC<PendingApprovalsProps> = ({
  approvals,
  loading,
  onApprove,
}) => {
  const navigate = useNavigate();
  const [processing, setProcessing] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleApprove = async (versionId: string) => {
    setProcessing(versionId);
    setError(null);
    try {
      // Get current user ID from store or localStorage
      const userId = localStorage.getItem('user_id') || 'admin';
      await versionsApi.approve(versionId, userId);
      if (onApprove) {
        onApprove(versionId);
      }
      // Refresh the page or update state
      window.location.reload();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to approve version');
    } finally {
      setProcessing(null);
    }
  };

  const isPendingMoreThan7Days = (createdAt: string) => {
    const daysDiff = (Date.now() - new Date(createdAt).getTime()) / (1000 * 60 * 60 * 24);
    return daysDiff > 7;
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const pendingVersions = approvals.filter(
    (v) => v.state === VersionState.PENDING_REVIEW
  );

  return (
    <Card title="Pending Approvals">
      {error && <Alert variant="error" className="mb-4">{error}</Alert>}
      {loading ? (
        <div className="flex items-center justify-center py-8">
          <Spinner />
        </div>
      ) : pendingVersions.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <p>No pending approvals</p>
        </div>
      ) : (
        <>
          <div className="space-y-4">
            {pendingVersions.map((version) => {
              const isOverdue = isPendingMoreThan7Days(version.created_at);
              return (
                <div
                  key={version.id}
                  className={`p-4 border rounded-lg ${
                    isOverdue
                      ? 'border-red-300 bg-red-50'
                      : 'border-gray-200 bg-white'
                  }`}
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-2">
                        <h4 className="text-sm font-semibold text-gray-900">
                          {version.product_id} - {version.version_number}
                        </h4>
                        {isOverdue && (
                          <Badge className="bg-red-100 text-red-800">Overdue</Badge>
                        )}
                      </div>
                      <p className="text-xs text-gray-500 mb-1">
                        Submitted: {formatDate(version.created_at)}
                      </p>
                      <p className="text-xs text-gray-500">
                        By: {version.created_by || 'Unknown'}
                      </p>
                    </div>
                    <div className="flex gap-2 ml-4">
                      <Button
                        size="sm"
                        onClick={() => handleApprove(version.id)}
                        isLoading={processing === version.id}
                        disabled={processing !== null}
                      >
                        Approve
                      </Button>
                      <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => navigate(`/versions/${version.id}`)}
                      >
                        View
                      </Button>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
          {pendingVersions.length > 0 && (
            <div className="mt-4 text-center">
              <button
                onClick={() => navigate('/versions?state=pending_review')}
                className="text-sm text-blue-600 hover:text-blue-700 font-medium"
              >
                View All Pending →
              </button>
            </div>
          )}
        </>
      )}
    </Card>
  );
};

