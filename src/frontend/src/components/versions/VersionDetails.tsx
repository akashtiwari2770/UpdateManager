import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { versionsApi } from '@/services/api/versions';
import { Version, VersionState, ReleaseType, ReleaseNotes, PackageInfo, PackageType } from '@/types';
import { Button, Card, Badge, Spinner, Alert, Modal } from '@/components/ui';
import { ReleaseNotesEditor, ReleaseNotesViewer, PackageUpload, PackagesList } from './index';
import { packagesApi } from '@/services/api/packages';
import { CompatibilityDetails, CompatibilityValidationForm } from '@/components/compatibility';
import { CreateUpgradePathForm, UpgradePathViewer, BlockUpgradePathDialog } from '@/components/upgrade-paths';

type TabType = 'overview' | 'release-notes' | 'packages' | 'compatibility';

export const VersionDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [version, setVersion] = useState<Version | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<TabType>('overview');
  const [showSubmitModal, setShowSubmitModal] = useState(false);
  const [showApproveModal, setShowApproveModal] = useState(false);
  const [showReleaseModal, setShowReleaseModal] = useState(false);
  const [showReleaseNotesEditor, setShowReleaseNotesEditor] = useState(false);
  const [showPackageUpload, setShowPackageUpload] = useState(false);
  const [showCompatibilityValidation, setShowCompatibilityValidation] = useState(false);
  const [showCreateUpgradePath, setShowCreateUpgradePath] = useState(false);
  const [showBlockUpgradePath, setShowBlockUpgradePath] = useState(false);
  const [selectedFromVersion, setSelectedFromVersion] = useState('');
  const [selectedToVersion, setSelectedToVersion] = useState('');
  const [processing, setProcessing] = useState(false);
  const [savingReleaseNotes, setSavingReleaseNotes] = useState(false);
  const [uploadingPackage, setUploadingPackage] = useState(false);
  const [approvalComment, setApprovalComment] = useState('');
  const [compatibilityRefreshTrigger, setCompatibilityRefreshTrigger] = useState(0);

  useEffect(() => {
    if (id) {
      loadVersion();
    }
  }, [id]);

  const loadVersion = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await versionsApi.getById(id!);
      setVersion(data);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load version');
      console.error('Error loading version:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmitForReview = async () => {
    if (!version) return;
    
    try {
      setProcessing(true);
      const updated = await versionsApi.submitForReview(version.id);
      setVersion(updated);
      setShowSubmitModal(false);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to submit for review');
    } finally {
      setProcessing(false);
    }
  };

  const handleApprove = async () => {
    if (!version) return;
    
    try {
      setProcessing(true);
      // Get current user from context or localStorage
      const approvedBy = localStorage.getItem('user_id') || 'current-user';
      const updated = await versionsApi.approve(version.id, approvedBy);
      setVersion(updated);
      setShowApproveModal(false);
      setApprovalComment('');
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to approve version');
    } finally {
      setProcessing(false);
    }
  };

  const handleRelease = async () => {
    if (!version) return;
    
    try {
      setProcessing(true);
      const updated = await versionsApi.release(version.id);
      setVersion(updated);
      setShowReleaseModal(false);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to release version');
    } finally {
      setProcessing(false);
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleString();
  };

  const getStateBadgeColor = (state: VersionState) => {
    switch (state) {
      case VersionState.DRAFT:
        return 'gray';
      case VersionState.PENDING_REVIEW:
        return 'yellow';
      case VersionState.APPROVED:
        return 'blue';
      case VersionState.RELEASED:
        return 'green';
      case VersionState.DEPRECATED:
        return 'orange';
      case VersionState.EOL:
        return 'red';
      default:
        return 'gray';
    }
  };

  const getReleaseTypeLabel = (type: ReleaseType) => {
    switch (type) {
      case ReleaseType.SECURITY:
        return 'Security';
      case ReleaseType.FEATURE:
        return 'Feature';
      case ReleaseType.MAINTENANCE:
        return 'Maintenance';
      case ReleaseType.MAJOR:
        return 'Major';
      default:
        return type;
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spinner />
      </div>
    );
  }

  if (error && !version) {
    return (
      <div className="space-y-4">
        <Button variant="ghost" onClick={() => navigate('/versions')}>
          ← Back to Versions
        </Button>
        <Alert variant="error" title="Error">
          {error}
        </Alert>
      </div>
    );
  }

  if (!version) {
    return (
      <div className="space-y-4">
        <Button variant="ghost" onClick={() => navigate('/versions')}>
          ← Back to Versions
        </Button>
        <Alert variant="info" title="Not Found">
          Version not found
        </Alert>
      </div>
    );
  }

  const canEdit = version.state === VersionState.DRAFT;
  const canSubmit = version.state === VersionState.DRAFT;
  const canApprove = version.state === VersionState.PENDING_REVIEW;
  const canRelease = version.state === VersionState.APPROVED;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" onClick={() => navigate('/versions')}>
            ← Back to Versions
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-gray-900">{version.version_number}</h1>
            <p className="text-gray-500 mt-1">Product: {version.product_id}</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Badge color={getStateBadgeColor(version.state)}>
            {version.state.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase())}
          </Badge>
          {canEdit && (
            <Button
              variant="secondary"
              onClick={() => navigate(`/versions/${version.id}/edit`)}
            >
              Edit Version
            </Button>
          )}
        </div>
      </div>

      {error && (
        <Alert variant="error" title="Error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* State-based Action Buttons */}
      <Card>
        <div className="flex items-center gap-2">
          {canSubmit && (
            <Button
              variant="primary"
              onClick={() => setShowSubmitModal(true)}
            >
              Submit for Review
            </Button>
          )}
          {canApprove && (
            <Button
              variant="primary"
              onClick={() => setShowApproveModal(true)}
            >
              Approve Version
            </Button>
          )}
          {canRelease && (
            <Button
              variant="primary"
              onClick={() => setShowReleaseModal(true)}
            >
              Release Version
            </Button>
          )}
        </div>
      </Card>

      {/* Tabs */}
      <Card>
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex space-x-8">
            {(['overview', 'release-notes', 'packages', 'compatibility'] as TabType[]).map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`
                  py-4 px-1 border-b-2 font-medium text-sm
                  ${
                    activeTab === tab
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }
                `}
              >
                {tab.replace('-', ' ').replace(/\b\w/g, (l) => l.toUpperCase())}
              </button>
            ))}
          </nav>
        </div>

        <div className="mt-6">
          {/* Overview Tab */}
          {activeTab === 'overview' && (
            <div className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h3 className="text-sm font-medium text-gray-500 mb-2">Version Information</h3>
                  <dl className="space-y-2">
                    <div>
                      <dt className="text-sm text-gray-500">Version Number</dt>
                      <dd className="text-sm font-medium text-gray-900">{version.version_number}</dd>
                    </div>
                    <div>
                      <dt className="text-sm text-gray-500">Release Type</dt>
                      <dd className="text-sm font-medium text-gray-900">{getReleaseTypeLabel(version.release_type)}</dd>
                    </div>
                    <div>
                      <dt className="text-sm text-gray-500">Release Date</dt>
                      <dd className="text-sm font-medium text-gray-900">{formatDate(version.release_date)}</dd>
                    </div>
                    {version.eol_date && (
                      <div>
                        <dt className="text-sm text-gray-500">End of Life Date</dt>
                        <dd className="text-sm font-medium text-gray-900">{formatDate(version.eol_date)}</dd>
                      </div>
                    )}
                  </dl>
                </div>

                <div>
                  <h3 className="text-sm font-medium text-gray-500 mb-2">State Information</h3>
                  <dl className="space-y-2">
                    <div>
                      <dt className="text-sm text-gray-500">Current State</dt>
                      <dd>
                        <Badge color={getStateBadgeColor(version.state)}>
                          {version.state.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase())}
                        </Badge>
                      </dd>
                    </div>
                    <div>
                      <dt className="text-sm text-gray-500">Created By</dt>
                      <dd className="text-sm font-medium text-gray-900">{version.created_by}</dd>
                    </div>
                    <div>
                      <dt className="text-sm text-gray-500">Created At</dt>
                      <dd className="text-sm font-medium text-gray-900">{formatDate(version.created_at)}</dd>
                    </div>
                    <div>
                      <dt className="text-sm text-gray-500">Last Updated</dt>
                      <dd className="text-sm font-medium text-gray-900">{formatDate(version.updated_at)}</dd>
                    </div>
                  </dl>
                </div>
              </div>

              {version.approved_by && (
                <div>
                  <h3 className="text-sm font-medium text-gray-500 mb-2">Approval History</h3>
                  <dl className="space-y-2">
                    <div>
                      <dt className="text-sm text-gray-500">Approved By</dt>
                      <dd className="text-sm font-medium text-gray-900">{version.approved_by}</dd>
                    </div>
                    {version.approved_at && (
                      <div>
                        <dt className="text-sm text-gray-500">Approved At</dt>
                        <dd className="text-sm font-medium text-gray-900">{formatDate(version.approved_at)}</dd>
                      </div>
                    )}
                  </dl>
                </div>
              )}

              {(version.min_server_version || version.max_server_version || version.recommended_server_version) && (
                <div>
                  <h3 className="text-sm font-medium text-gray-500 mb-2">Server Version Requirements</h3>
                  <dl className="space-y-2">
                    {version.min_server_version && (
                      <div>
                        <dt className="text-sm text-gray-500">Minimum Server Version</dt>
                        <dd className="text-sm font-medium text-gray-900">{version.min_server_version}</dd>
                      </div>
                    )}
                    {version.max_server_version && (
                      <div>
                        <dt className="text-sm text-gray-500">Maximum Server Version</dt>
                        <dd className="text-sm font-medium text-gray-900">{version.max_server_version}</dd>
                      </div>
                    )}
                    {version.recommended_server_version && (
                      <div>
                        <dt className="text-sm text-gray-500">Recommended Server Version</dt>
                        <dd className="text-sm font-medium text-gray-900">{version.recommended_server_version}</dd>
                      </div>
                    )}
                  </dl>
                </div>
              )}
            </div>
          )}

          {/* Release Notes Tab */}
          {activeTab === 'release-notes' && (
            <ReleaseNotesViewer
              version={version}
              releaseNotes={version.release_notes}
              onEdit={() => setShowReleaseNotesEditor(true)}
              canEdit={version.state === VersionState.DRAFT}
            />
          )}

          {/* Packages Tab */}
          {activeTab === 'packages' && (
            <PackagesList
              packages={version.packages || []}
              onDownload={async (packageId: string, fileName: string) => {
                // Find the package
                const pkg = version.packages?.find((p) => p.id === packageId);
                if (pkg?.download_url) {
                  window.open(pkg.download_url, '_blank');
                } else {
                  // Fallback: construct download URL
                  window.open(`/api/v1/versions/${version.id}/packages/${packageId}/download`, '_blank');
                }
              }}
              onUpload={() => setShowPackageUpload(true)}
              canUpload={version.state === VersionState.DRAFT}
            />
          )}

          {/* Compatibility Tab */}
          {activeTab === 'compatibility' && (
            <CompatibilityDetails
              productId={version.product_id}
              versionNumber={version.version_number}
              onValidate={() => setShowCompatibilityValidation(true)}
              refreshTrigger={compatibilityRefreshTrigger}
            />
          )}
        </div>
      </Card>

      {/* Submit for Review Modal */}
      <Modal
        isOpen={showSubmitModal}
        onClose={() => setShowSubmitModal(false)}
        title="Submit for Review"
      >
        <div className="space-y-4">
          <p>Are you sure you want to submit this version for review? Once submitted, you won't be able to edit it.</p>
          <div className="flex justify-end gap-2">
            <Button variant="secondary" onClick={() => setShowSubmitModal(false)} disabled={processing}>
              Cancel
            </Button>
            <Button onClick={handleSubmitForReview} disabled={processing}>
              {processing ? <Spinner size="sm" /> : 'Submit for Review'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* Approve Modal */}
      <Modal
        isOpen={showApproveModal}
        onClose={() => setShowApproveModal(false)}
        title="Approve Version"
      >
        <div className="space-y-4">
          <p>Are you sure you want to approve this version?</p>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Approval Comment (Optional)
            </label>
            <textarea
              className="w-full border border-gray-300 rounded-md px-3 py-2"
              rows={3}
              value={approvalComment}
              onChange={(e) => setApprovalComment(e.target.value)}
              placeholder="Add any comments about this approval..."
            />
          </div>
          <div className="flex justify-end gap-2">
            <Button variant="secondary" onClick={() => setShowApproveModal(false)} disabled={processing}>
              Cancel
            </Button>
            <Button onClick={handleApprove} disabled={processing}>
              {processing ? <Spinner size="sm" /> : 'Approve'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* Release Modal */}
      <Modal
        isOpen={showReleaseModal}
        onClose={() => setShowReleaseModal(false)}
        title="Release Version"
      >
        <div className="space-y-4">
          <p>Are you sure you want to release version <strong>{version.version_number}</strong>?</p>
          <div className="bg-yellow-50 border border-yellow-200 rounded-md p-3">
            <p className="text-sm text-yellow-800">
              This will make the version available for deployment. This action cannot be undone.
            </p>
          </div>
          <div className="flex justify-end gap-2">
            <Button variant="secondary" onClick={() => setShowReleaseModal(false)} disabled={processing}>
              Cancel
            </Button>
            <Button onClick={handleRelease} disabled={processing}>
              {processing ? <Spinner size="sm" /> : 'Release Version'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* Release Notes Editor Modal */}
      <Modal
        isOpen={showReleaseNotesEditor}
        onClose={() => setShowReleaseNotesEditor(false)}
        title="Edit Release Notes"
        size="lg"
      >
        <ReleaseNotesEditor
          releaseNotes={version.release_notes}
          versionNumber={version.version_number}
          releaseDate={version.release_date}
          releaseType={version.release_type}
          onSave={async (notes: ReleaseNotes) => {
            try {
              setSavingReleaseNotes(true);
              const updated = await versionsApi.update(version.id, { release_notes: notes });
              setVersion(updated);
              setShowReleaseNotesEditor(false);
            } catch (err: any) {
              throw err;
            } finally {
              setSavingReleaseNotes(false);
            }
          }}
          onCancel={() => setShowReleaseNotesEditor(false)}
          loading={savingReleaseNotes}
        />
      </Modal>

      {/* Package Upload Modal */}
      <Modal
        isOpen={showPackageUpload}
        onClose={() => setShowPackageUpload(false)}
        title="Upload Package"
        size="lg"
      >
        <PackageUpload
          versionId={version.id}
          onUpload={async (file: File, metadata, onProgress) => {
            try {
              setUploadingPackage(true);
              await packagesApi.upload(version.id, {
                file,
                package_type: metadata.package_type,
                os: metadata.os,
                architecture: metadata.architecture,
              }, onProgress);
              // Reload version to get updated packages list
              await loadVersion();
              setShowPackageUpload(false);
            } catch (err: any) {
              throw err;
            } finally {
              setUploadingPackage(false);
            }
          }}
          onCancel={() => setShowPackageUpload(false)}
          loading={uploadingPackage}
        />
      </Modal>

      {/* Compatibility Validation Modal */}
      <Modal
        isOpen={showCompatibilityValidation}
        onClose={() => setShowCompatibilityValidation(false)}
        title="Validate Compatibility"
        size="lg"
      >
        <CompatibilityValidationForm
          productId={version.product_id}
          versionNumber={version.version_number}
          onValidate={async (result) => {
            // Trigger refresh of compatibility details
            setCompatibilityRefreshTrigger(prev => prev + 1);
            setShowCompatibilityValidation(false);
          }}
          onCancel={() => setShowCompatibilityValidation(false)}
        />
      </Modal>

      {/* Create Upgrade Path Modal */}
      <Modal
        isOpen={showCreateUpgradePath}
        onClose={() => setShowCreateUpgradePath(false)}
        title="Create Upgrade Path"
        size="lg"
      >
        <CreateUpgradePathForm
          productId={version.product_id}
          onSuccess={async () => {
            await loadVersion();
            setShowCreateUpgradePath(false);
          }}
          onCancel={() => setShowCreateUpgradePath(false)}
        />
      </Modal>

      {/* Block Upgrade Path Dialog */}
      <BlockUpgradePathDialog
        isOpen={showBlockUpgradePath}
        productId={version.product_id}
        fromVersion={selectedFromVersion}
        toVersion={selectedToVersion}
        onClose={() => {
          setShowBlockUpgradePath(false);
          setSelectedFromVersion('');
          setSelectedToVersion('');
        }}
        onSuccess={async () => {
          await loadVersion();
        }}
      />
    </div>
  );
};

