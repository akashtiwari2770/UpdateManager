import React, { useEffect } from 'react';
import { Alert } from './Alert';

export interface ToastProps {
  message: string;
  type?: 'success' | 'error' | 'info' | 'warning';
  duration?: number;
  onClose: () => void;
}

export const Toast: React.FC<ToastProps> = ({
  message,
  type = 'info',
  duration = 5000,
  onClose,
}) => {
  useEffect(() => {
    const timer = setTimeout(() => {
      onClose();
    }, duration);

    return () => clearTimeout(timer);
  }, [duration, onClose]);

  const variantMap = {
    success: 'success' as const,
    error: 'error' as const,
    info: 'info' as const,
    warning: 'warning' as const,
  };

  return (
    <div className="animate-in slide-in-from-top-5">
      <Alert variant={variantMap[type]}>{message}</Alert>
    </div>
  );
};

