import React from 'react';
import { useParams } from 'react-router-dom';
import { RolloutStatus } from '@/components/rollouts';

export const RolloutDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();

  if (!id) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">Invalid rollout ID.</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-900">Rollout Details</h1>
      <RolloutStatus rolloutId={id} />
    </div>
  );
};

