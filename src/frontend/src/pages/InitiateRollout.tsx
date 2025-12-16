import React from 'react';
import { useNavigate } from 'react-router-dom';
import { InitiateRolloutForm } from '@/components/rollouts';
import { UpdateRollout } from '@/types';

export const InitiateRollout: React.FC = () => {
  const navigate = useNavigate();

  const handleSuccess = (rollout: UpdateRollout) => {
    navigate(`/updates/rollouts/${rollout.id}`);
  };

  const handleCancel = () => {
    navigate('/updates');
  };

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-900">Initiate Update Rollout</h1>
      <InitiateRolloutForm onSuccess={handleSuccess} onCancel={handleCancel} />
    </div>
  );
};

