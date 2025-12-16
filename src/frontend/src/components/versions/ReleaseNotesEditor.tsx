import React, { useState } from 'react';
import { ReleaseNotes } from '@/types';
import { Button, Card, Input, Alert } from '@/components/ui';

export interface ReleaseNotesEditorProps {
  releaseNotes?: ReleaseNotes;
  versionNumber?: string;
  releaseDate?: string;
  releaseType?: string;
  onSave: (releaseNotes: ReleaseNotes) => Promise<void>;
  onCancel: () => void;
  loading?: boolean;
}

export const ReleaseNotesEditor: React.FC<ReleaseNotesEditorProps> = ({
  releaseNotes,
  versionNumber = '',
  releaseDate = '',
  releaseType = 'feature',
  onSave,
  onCancel,
  loading = false,
}) => {
  const [whatsNew, setWhatsNew] = useState<string[]>(
    releaseNotes?.whats_new || []
  );
  const [bugFixes, setBugFixes] = useState(
    releaseNotes?.bug_fixes || []
  );
  const [breakingChanges, setBreakingChanges] = useState(
    releaseNotes?.breaking_changes || []
  );
  const [knownIssues, setKnownIssues] = useState(
    releaseNotes?.known_issues || []
  );
  const [upgradeInstructions, setUpgradeInstructions] = useState(
    releaseNotes?.upgrade_instructions || ''
  );
  const [error, setError] = useState<string | null>(null);

  const addWhatsNewItem = () => {
    setWhatsNew([...whatsNew, '']);
  };

  const removeWhatsNewItem = (index: number) => {
    setWhatsNew(whatsNew.filter((_, i) => i !== index));
  };

  const updateWhatsNewItem = (index: number, value: string) => {
    const updated = [...whatsNew];
    updated[index] = value;
    setWhatsNew(updated);
  };

  const addBugFix = () => {
    setBugFixes([
      ...bugFixes,
      { id: '', description: '', issue_number: '' },
    ]);
  };

  const removeBugFix = (index: number) => {
    setBugFixes(bugFixes.filter((_, i) => i !== index));
  };

  const updateBugFix = (index: number, field: 'id' | 'description' | 'issue_number', value: string) => {
    const updated = [...bugFixes];
    updated[index] = { ...updated[index], [field]: value };
    setBugFixes(updated);
  };

  const addBreakingChange = () => {
    setBreakingChanges([
      ...breakingChanges,
      { description: '', migration_steps: '', configuration_changes: '' },
    ]);
  };

  const removeBreakingChange = (index: number) => {
    setBreakingChanges(breakingChanges.filter((_, i) => i !== index));
  };

  const updateBreakingChange = (
    index: number,
    field: 'description' | 'migration_steps' | 'configuration_changes',
    value: string
  ) => {
    const updated = [...breakingChanges];
    updated[index] = { ...updated[index], [field]: value };
    setBreakingChanges(updated);
  };

  const addKnownIssue = () => {
    setKnownIssues([
      ...knownIssues,
      { id: '', description: '', workaround: '', planned_fix: '' },
    ]);
  };

  const removeKnownIssue = (index: number) => {
    setKnownIssues(knownIssues.filter((_, i) => i !== index));
  };

  const updateKnownIssue = (
    index: number,
    field: 'id' | 'description' | 'workaround' | 'planned_fix',
    value: string
  ) => {
    const updated = [...knownIssues];
    updated[index] = { ...updated[index], [field]: value };
    setKnownIssues(updated);
  };

  const handleSave = async () => {
    try {
      setError(null);
      
      // Use existing version_info if available, otherwise use provided props
      const existingVersionInfo = releaseNotes?.version_info;
      const notes: ReleaseNotes = {
        version_info: {
          version_number: existingVersionInfo?.version_number || versionNumber || '',
          release_date: existingVersionInfo?.release_date || releaseDate || new Date().toISOString(),
          release_type: (existingVersionInfo?.release_type || releaseType || 'feature') as any,
        },
        whats_new: whatsNew.filter((item) => item.trim() !== ''),
        bug_fixes: bugFixes.filter(
          (fix) => fix.description.trim() !== ''
        ),
        breaking_changes: breakingChanges.filter(
          (change) => change.description.trim() !== ''
        ),
        compatibility: releaseNotes?.compatibility || {
          server_version_requirements: '',
          client_version_requirements: '',
          os_requirements: [],
        },
        upgrade_instructions: upgradeInstructions,
        known_issues: knownIssues.filter(
          (issue) => issue.description.trim() !== ''
        ),
      };
      await onSave(notes);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to save release notes');
    }
  };

  return (
    <div className="space-y-6">
      {error && (
        <Alert variant="error" title="Error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* What's New */}
      <Card title="What's New">
        <div className="space-y-3">
          {whatsNew.map((item, index) => (
            <div key={index} className="flex gap-2">
              <Input
                value={item}
                onChange={(e) => updateWhatsNewItem(index, e.target.value)}
                placeholder="Enter new feature or improvement..."
                className="flex-1"
              />
              <Button
                variant="ghost"
                size="sm"
                onClick={() => removeWhatsNewItem(index)}
              >
                Remove
              </Button>
            </div>
          ))}
          <Button variant="secondary" onClick={addWhatsNewItem}>
            + Add Item
          </Button>
        </div>
      </Card>

      {/* Bug Fixes */}
      <Card title="Bug Fixes">
        <div className="space-y-4">
          {bugFixes.map((fix, index) => (
            <div key={index} className="border border-gray-200 rounded-lg p-4 space-y-3">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                <Input
                  label="Fix ID"
                  value={fix.id}
                  onChange={(e) => updateBugFix(index, 'id', e.target.value)}
                  placeholder="e.g., BUG-123"
                />
                <Input
                  label="Issue Number"
                  value={fix.issue_number}
                  onChange={(e) => updateBugFix(index, 'issue_number', e.target.value)}
                  placeholder="e.g., #456"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <textarea
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  rows={2}
                  value={fix.description}
                  onChange={(e) => updateBugFix(index, 'description', e.target.value)}
                  placeholder="Describe the bug fix..."
                />
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => removeBugFix(index)}
              >
                Remove
              </Button>
            </div>
          ))}
          <Button variant="secondary" onClick={addBugFix}>
            + Add Bug Fix
          </Button>
        </div>
      </Card>

      {/* Breaking Changes */}
      <Card title="Breaking Changes">
        <div className="space-y-4">
          {breakingChanges.map((change, index) => (
            <div key={index} className="border border-gray-200 rounded-lg p-4 space-y-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <textarea
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  rows={2}
                  value={change.description}
                  onChange={(e) => updateBreakingChange(index, 'description', e.target.value)}
                  placeholder="Describe the breaking change..."
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Migration Steps
                </label>
                <textarea
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  rows={3}
                  value={change.migration_steps}
                  onChange={(e) => updateBreakingChange(index, 'migration_steps', e.target.value)}
                  placeholder="Steps to migrate..."
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Configuration Changes
                </label>
                <textarea
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  rows={2}
                  value={change.configuration_changes}
                  onChange={(e) => updateBreakingChange(index, 'configuration_changes', e.target.value)}
                  placeholder="Configuration changes required..."
                />
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => removeBreakingChange(index)}
              >
                Remove
              </Button>
            </div>
          ))}
          <Button variant="secondary" onClick={addBreakingChange}>
            + Add Breaking Change
          </Button>
        </div>
      </Card>

      {/* Known Issues */}
      <Card title="Known Issues">
        <div className="space-y-4">
          {knownIssues.map((issue, index) => (
            <div key={index} className="border border-gray-200 rounded-lg p-4 space-y-3">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                <Input
                  label="Issue ID"
                  value={issue.id}
                  onChange={(e) => updateKnownIssue(index, 'id', e.target.value)}
                  placeholder="e.g., ISSUE-789"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <textarea
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  rows={2}
                  value={issue.description}
                  onChange={(e) => updateKnownIssue(index, 'description', e.target.value)}
                  placeholder="Describe the known issue..."
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Workaround (Optional)
                </label>
                <textarea
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  rows={2}
                  value={issue.workaround || ''}
                  onChange={(e) => updateKnownIssue(index, 'workaround', e.target.value)}
                  placeholder="Workaround if available..."
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Planned Fix (Optional)
                </label>
                <textarea
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  rows={2}
                  value={issue.planned_fix || ''}
                  onChange={(e) => updateKnownIssue(index, 'planned_fix', e.target.value)}
                  placeholder="Planned fix or timeline..."
                />
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => removeKnownIssue(index)}
              >
                Remove
              </Button>
            </div>
          ))}
          <Button variant="secondary" onClick={addKnownIssue}>
            + Add Known Issue
          </Button>
        </div>
      </Card>

      {/* Upgrade Instructions */}
      <Card title="Upgrade Instructions">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Instructions
          </label>
          <textarea
            className="w-full px-3 py-2 border border-gray-300 rounded-lg"
            rows={6}
            value={upgradeInstructions}
            onChange={(e) => setUpgradeInstructions(e.target.value)}
            placeholder="Provide upgrade instructions, prerequisites, and any important notes..."
          />
          <p className="mt-2 text-sm text-gray-500">
            You can use markdown formatting for better readability.
          </p>
        </div>
      </Card>

      {/* Actions */}
      <div className="flex items-center justify-end gap-4 pt-4 border-t">
        <Button variant="secondary" onClick={onCancel} disabled={loading}>
          Cancel
        </Button>
        <Button variant="primary" onClick={handleSave} isLoading={loading}>
          Save Release Notes
        </Button>
      </div>
    </div>
  );
};

