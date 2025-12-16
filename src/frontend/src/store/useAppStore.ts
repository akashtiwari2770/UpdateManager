import { create } from 'zustand';

interface AppState {
  // User state
  userId: string | null;
  setUserId: (userId: string | null) => void;

  // UI state
  sidebarCollapsed: boolean;
  setSidebarCollapsed: (collapsed: boolean) => void;

  // Notifications
  unreadNotificationCount: number;
  setUnreadNotificationCount: (count: number) => void;
}

export const useAppStore = create<AppState>((set) => ({
  // User state
  userId: localStorage.getItem('user_id'),
  setUserId: (userId) => {
    if (userId) {
      localStorage.setItem('user_id', userId);
    } else {
      localStorage.removeItem('user_id');
    }
    set({ userId });
  },

  // UI state
  sidebarCollapsed: false,
  setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),

  // Notifications
  unreadNotificationCount: 0,
  setUnreadNotificationCount: (count) => set({ unreadNotificationCount: count }),
}));

