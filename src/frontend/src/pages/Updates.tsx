import React, { useState } from 'react';
import { UpdatesDashboard } from '@/components/updates';
import { UpdateDetectionList, UpdateDetectionForm } from '@/components/update-detections';
import { RolloutList } from '@/components/rollouts';
import { PendingUpdatesList } from '@/components/updates/PendingUpdatesList';
import { Button, Modal } from '@/components/ui';

type TabType = 'dashboard' | 'detections' | 'rollouts' | 'deployments';

export const Updates: React.FC = () => {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard');
  const [showDetectionForm, setShowDetectionForm] = useState(false);
  const [detectionKey, setDetectionKey] = useState(0); // Force re-render of detection list

  const tabs: { id: TabType; label: string }[] = [
    { id: 'dashboard', label: 'Available Updates' },
    { id: 'deployments', label: 'Deployment Updates' },
    { id: 'detections', label: 'Update Detections' },
    { id: 'rollouts', label: 'Rollouts' },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Updates</h1>
        {activeTab === 'detections' && (
          <Button variant="primary" onClick={() => setShowDetectionForm(true)}>
            Register Detection
          </Button>
        )}
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200">
        <nav className="-mb-px flex space-x-8">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`
                py-4 px-1 border-b-2 font-medium text-sm
                ${
                  activeTab === tab.id
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }
              `}
            >
              {tab.label}
            </button>
          ))}
        </nav>
      </div>

      {/* Tab Content */}
      <div className="mt-6">
        {activeTab === 'dashboard' && <UpdatesDashboard />}
        {activeTab === 'deployments' && <PendingUpdatesList view="all" />}
        {activeTab === 'detections' && <UpdateDetectionList key={detectionKey} />}
        {activeTab === 'rollouts' && <RolloutList />}
      </div>

      {/* Register Detection Modal */}
      <Modal
        isOpen={showDetectionForm}
        onClose={() => setShowDetectionForm(false)}
        title="Register Update Detection"
        size="lg"
      >
        <UpdateDetectionForm
          onSuccess={() => {
            setShowDetectionForm(false);
            // Force reload of detection list by changing key
            setDetectionKey(prev => prev + 1);
            // Switch to detections tab to show the new detection
            setActiveTab('detections');
          }}
          onCancel={() => setShowDetectionForm(false)}
        />
      </Modal>
    </div>
  );
};

