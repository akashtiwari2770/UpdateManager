import React from 'react';
import { ReleaseNotes, Version } from '@/types';
import { Card, Button } from '@/components/ui';

export interface ReleaseNotesViewerProps {
  version: Version;
  releaseNotes?: ReleaseNotes;
  onEdit?: () => void;
  canEdit?: boolean;
}

export const ReleaseNotesViewer: React.FC<ReleaseNotesViewerProps> = ({
  version,
  releaseNotes,
  onEdit,
  canEdit = false,
}) => {
  const handlePrint = () => {
    window.print();
  };

  const handleShare = () => {
    if (navigator.share) {
      navigator.share({
        title: `Release Notes - ${version.version_number}`,
        text: `Release notes for version ${version.version_number}`,
        url: window.location.href,
      });
    } else {
      // Fallback: copy to clipboard
      navigator.clipboard.writeText(window.location.href);
      alert('Link copied to clipboard!');
    }
  };

  if (!releaseNotes) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 mb-4">No release notes available for this version.</p>
        {canEdit && onEdit && (
          <Button variant="primary" onClick={onEdit}>
            Add Release Notes
          </Button>
        )}
      </div>
    );
  }

  return (
    <div className="space-y-6 print:space-y-4">
      {/* Header Actions */}
      <div className="flex items-center justify-between print:hidden">
        <h2 className="text-2xl font-bold text-gray-900">
          Release Notes - {version.version_number}
        </h2>
        <div className="flex items-center gap-2">
          <Button variant="secondary" onClick={handleShare}>
            Share
          </Button>
          <Button variant="secondary" onClick={handlePrint}>
            Print
          </Button>
          {canEdit && onEdit && (
            <Button variant="primary" onClick={onEdit}>
              Edit
            </Button>
          )}
        </div>
      </div>

      {/* Version Information */}
      <Card>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Version Information</h3>
        <dl className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <dt className="text-sm font-medium text-gray-500">Version Number</dt>
            <dd className="mt-1 text-sm text-gray-900">{version.version_number}</dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Release Type</dt>
            <dd className="mt-1 text-sm text-gray-900">{version.release_type}</dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Release Date</dt>
            <dd className="mt-1 text-sm text-gray-900">
              {new Date(version.release_date).toLocaleDateString()}
            </dd>
          </div>
          {version.eol_date && (
            <div>
              <dt className="text-sm font-medium text-gray-500">End of Life Date</dt>
              <dd className="mt-1 text-sm text-gray-900">
                {new Date(version.eol_date).toLocaleDateString()}
              </dd>
            </div>
          )}
        </dl>
      </Card>

      {/* What's New */}
      {releaseNotes.whats_new && releaseNotes.whats_new.length > 0 && (
        <Card>
          <h3 className="text-lg font-semibold text-gray-900 mb-4">What's New</h3>
          <ul className="list-disc list-inside space-y-2">
            {releaseNotes.whats_new.map((item, index) => (
              <li key={index} className="text-gray-700">{item}</li>
            ))}
          </ul>
        </Card>
      )}

      {/* Bug Fixes */}
      {releaseNotes.bug_fixes && releaseNotes.bug_fixes.length > 0 && (
        <Card>
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Bug Fixes</h3>
          <ul className="space-y-4">
            {releaseNotes.bug_fixes.map((fix, index) => (
              <li key={index} className="border-l-4 border-blue-500 pl-4">
                {fix.id && (
                  <div className="font-medium text-gray-900 mb-1">
                    {fix.id}
                    {fix.issue_number && (
                      <span className="text-gray-500 ml-2">(Issue #{fix.issue_number})</span>
                    )}
                  </div>
                )}
                <div className="text-gray-700">{fix.description}</div>
              </li>
            ))}
          </ul>
        </Card>
      )}

      {/* Breaking Changes */}
      {releaseNotes.breaking_changes && releaseNotes.breaking_changes.length > 0 && (
        <Card>
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Breaking Changes</h3>
          <div className="space-y-4">
            {releaseNotes.breaking_changes.map((change, index) => (
              <div key={index} className="border-l-4 border-red-500 pl-4">
                <div className="text-gray-700 mb-2">{change.description}</div>
                {change.migration_steps && (
                  <div className="mt-3">
                    <h4 className="text-sm font-medium text-gray-900 mb-1">Migration Steps:</h4>
                    <div className="text-sm text-gray-700 whitespace-pre-line">
                      {change.migration_steps}
                    </div>
                  </div>
                )}
                {change.configuration_changes && (
                  <div className="mt-3">
                    <h4 className="text-sm font-medium text-gray-900 mb-1">Configuration Changes:</h4>
                    <div className="text-sm text-gray-700 whitespace-pre-line">
                      {change.configuration_changes}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </Card>
      )}

      {/* Known Issues */}
      {releaseNotes.known_issues && releaseNotes.known_issues.length > 0 && (
        <Card>
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Known Issues</h3>
          <div className="space-y-4">
            {releaseNotes.known_issues.map((issue, index) => (
              <div key={index} className="border-l-4 border-yellow-500 pl-4">
                {issue.id && (
                  <div className="font-medium text-gray-900 mb-1">{issue.id}</div>
                )}
                <div className="text-gray-700 mb-2">{issue.description}</div>
                {issue.workaround && (
                  <div className="mt-2">
                    <h4 className="text-sm font-medium text-gray-900 mb-1">Workaround:</h4>
                    <div className="text-sm text-gray-700">{issue.workaround}</div>
                  </div>
                )}
                {issue.planned_fix && (
                  <div className="mt-2">
                    <h4 className="text-sm font-medium text-gray-900 mb-1">Planned Fix:</h4>
                    <div className="text-sm text-gray-700">{issue.planned_fix}</div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </Card>
      )}

      {/* Upgrade Instructions */}
      {releaseNotes.upgrade_instructions && (
        <Card>
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Upgrade Instructions</h3>
          <div className="text-gray-700 whitespace-pre-line">
            {releaseNotes.upgrade_instructions}
          </div>
        </Card>
      )}

      {/* Empty State */}
      {(!releaseNotes.whats_new || releaseNotes.whats_new.length === 0) &&
        (!releaseNotes.bug_fixes || releaseNotes.bug_fixes.length === 0) &&
        (!releaseNotes.breaking_changes || releaseNotes.breaking_changes.length === 0) &&
        (!releaseNotes.known_issues || releaseNotes.known_issues.length === 0) &&
        !releaseNotes.upgrade_instructions && (
          <div className="text-center py-12">
            <p className="text-gray-500 mb-4">No release notes content available.</p>
            {canEdit && onEdit && (
              <Button variant="primary" onClick={onEdit}>
                Add Release Notes
              </Button>
            )}
          </div>
        )}
    </div>
  );
};

