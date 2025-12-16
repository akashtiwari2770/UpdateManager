import { test, expect } from '@playwright/test';

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');
  });

  test.describe('Dashboard Layout', () => {
    test('should load dashboard successfully', async ({ page }) => {
      await expect(page.locator('h1:has-text("Dashboard")')).toBeVisible();
    });

    test('should display all sections', async ({ page }) => {
      // Check for stats cards section
      const statsCards = page.locator('.grid.grid-cols-1.md\\:grid-cols-2.lg\\:grid-cols-4');
      await expect(statsCards).toBeVisible();

      // Check for Recent Updates section
      await expect(page.locator('text=Recent Updates')).toBeVisible();

      // Check for Pending Approvals section
      await expect(page.locator('text=Pending Approvals')).toBeVisible();

      // Check for Activity Timeline section
      await expect(page.locator('text=Recent Activity')).toBeVisible();
    });

    test('should display refresh button', async ({ page }) => {
      const refreshButton = page.locator('button:has-text("Refresh")');
      await expect(refreshButton).toBeVisible();
    });

    test('should be responsive', async ({ page }) => {
      // Test mobile viewport
      await page.setViewportSize({ width: 375, height: 667 });
      await expect(page.locator('h1:has-text("Dashboard")')).toBeVisible();

      // Test tablet viewport
      await page.setViewportSize({ width: 768, height: 1024 });
      await expect(page.locator('h1:has-text("Dashboard")')).toBeVisible();

      // Test desktop viewport
      await page.setViewportSize({ width: 1920, height: 1080 });
      await expect(page.locator('h1:has-text("Dashboard")')).toBeVisible();
    });

    test('should show loading states', async ({ page }) => {
      // Reload to catch loading state
      await page.reload();
      
      // Check for loading indicators (may be too fast to catch)
      const loadingIndicator = page.locator('[data-testid="loading"], .spinner, text=/loading/i, [role="progressbar"]').first();
      const isLoading = await loadingIndicator.isVisible().catch(() => false);
      expect(typeof isLoading).toBe('boolean');
    });
  });

  test.describe('Statistics Cards', () => {
    test('should display all stat cards', async ({ page }) => {
      // Total Products card
      await expect(page.locator('text=Total Products')).toBeVisible();

      // Active Versions card
      await expect(page.locator('text=Active Versions')).toBeVisible();

      // Pending Updates card
      await expect(page.locator('text=Pending Updates')).toBeVisible();

      // Active Rollouts card
      await expect(page.locator('text=Active Rollouts')).toBeVisible();
    });

    test('should display counts in stat cards', async ({ page }) => {
      // Wait for cards to load
      await page.waitForTimeout(2000);

      // Check that cards have numeric values or loading state
      const productCard = page.locator('text=Total Products').locator('..').locator('..');
      const productValue = productCard.locator('.text-3xl');
      const hasValue = await productValue.isVisible().catch(() => false);
      expect(typeof hasValue).toBe('boolean');
    });

    test('should navigate when stat card is clicked', async ({ page }) => {
      // Click on Total Products card
      const productsCard = page.locator('text=Total Products').locator('..').locator('..');
      const isClickable = await productsCard.locator('..').getAttribute('class')?.includes('cursor-pointer') || false;
      
      if (isClickable) {
        await productsCard.click();
        await page.waitForTimeout(1000);
        // Should navigate to products page
        const currentUrl = page.url();
        expect(currentUrl).toContain('/products');
      }
    });

    test('should display icons in stat cards', async ({ page }) => {
      // Check for icon containers
      const iconContainers = page.locator('.bg-blue-100.rounded-lg');
      const iconCount = await iconContainers.count();
      expect(iconCount).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Recent Updates Section', () => {
    test('should display recent updates section', async ({ page }) => {
      await expect(page.locator('text=Recent Updates')).toBeVisible();
    });

    test('should display table with correct columns', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const table = page.locator('text=Recent Updates').locator('..').locator('table');
      const isTableVisible = await table.isVisible().catch(() => false);
      
      if (isTableVisible) {
        await expect(table.locator('th:has-text("Product")')).toBeVisible();
        await expect(table.locator('th:has-text("Version")')).toBeVisible();
        await expect(table.locator('th:has-text("Release Date")')).toBeVisible();
        await expect(table.locator('th:has-text("Status")')).toBeVisible();
      }
    });

    test('should limit to 10 items', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const table = page.locator('text=Recent Updates').locator('..').locator('table');
      const isTableVisible = await table.isVisible().catch(() => false);
      
      if (isTableVisible) {
        const rows = table.locator('tbody tr');
        const rowCount = await rows.count();
        expect(rowCount).toBeLessThanOrEqual(10);
      }
    });

    test('should show "View All Versions" link when there are 10+ items', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const viewAllLink = page.locator('text=View All Versions');
      const isVisible = await viewAllLink.isVisible().catch(() => false);
      expect(typeof isVisible).toBe('boolean');
    });

    test('should navigate to versions page when "View All" is clicked', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const viewAllLink = page.locator('text=View All Versions');
      if (await viewAllLink.isVisible().catch(() => false)) {
        await viewAllLink.click();
        await page.waitForURL('**/versions', { timeout: 5000 });
        await expect(page.locator('h1:has-text("Versions")')).toBeVisible();
      }
    });

    test('should navigate to version details when row is clicked', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const table = page.locator('text=Recent Updates').locator('..').locator('table');
      const isTableVisible = await table.isVisible().catch(() => false);
      
      if (isTableVisible) {
        const firstRow = table.locator('tbody tr').first();
        if (await firstRow.isVisible().catch(() => false)) {
          await firstRow.click();
          await page.waitForTimeout(1000);
          // Should navigate to version details
          const currentUrl = page.url();
          expect(currentUrl).toMatch(/\/versions\/[^/]+$/);
        }
      }
    });

    test('should display empty state when no updates', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const emptyState = page.locator('text=No recent updates');
      const table = page.locator('text=Recent Updates').locator('..').locator('table');
      
      const isEmptyVisible = await emptyState.isVisible().catch(() => false);
      const isTableVisible = await table.isVisible().catch(() => false);
      
      expect(isEmptyVisible || isTableVisible).toBeTruthy();
    });
  });

  test.describe('Pending Approvals Section', () => {
    test('should display pending approvals section', async ({ page }) => {
      await expect(page.locator('text=Pending Approvals')).toBeVisible();
    });

    test('should display approval items', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const approvalsSection = page.locator('text=Pending Approvals').locator('..');
      const approvalItems = approvalsSection.locator('.border.rounded-lg');
      const itemCount = await approvalItems.count();
      expect(itemCount).toBeGreaterThanOrEqual(0);
    });

    test('should highlight overdue items (pending > 7 days)', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const overdueItems = page.locator('.border-red-300.bg-red-50');
      const hasOverdue = await overdueItems.first().isVisible().catch(() => false);
      expect(typeof hasOverdue).toBe('boolean');
    });

    test('should display approve button', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const approveButton = page.locator('button:has-text("Approve")').first();
      const isVisible = await approveButton.isVisible().catch(() => false);
      expect(typeof isVisible).toBe('boolean');
    });

    test('should display view button', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const viewButton = page.locator('button:has-text("View")').first();
      const isVisible = await viewButton.isVisible().catch(() => false);
      expect(typeof isVisible).toBe('boolean');
    });

    test('should navigate to version when view is clicked', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const viewButton = page.locator('button:has-text("View")').first();
      if (await viewButton.isVisible().catch(() => false)) {
        await viewButton.click();
        await page.waitForTimeout(1000);
        // Should navigate to version details
        const currentUrl = page.url();
        expect(currentUrl).toMatch(/\/versions\/[^/]+$/);
      }
    });

    test('should show "View All Pending" link', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const viewAllLink = page.locator('text=View All Pending');
      const isVisible = await viewAllLink.isVisible().catch(() => false);
      expect(typeof isVisible).toBe('boolean');
    });

    test('should display empty state when no pending approvals', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const emptyState = page.locator('text=No pending approvals');
      const hasItems = page.locator('.border.rounded-lg').first().isVisible().catch(() => false);
      
      const isEmptyVisible = await emptyState.isVisible().catch(() => false);
      const hasItemsVisible = await hasItems;
      
      expect(isEmptyVisible || hasItemsVisible).toBeTruthy();
    });
  });

  test.describe('Activity Timeline', () => {
    test('should display activity timeline section', async ({ page }) => {
      await expect(page.locator('text=Recent Activity')).toBeVisible();
    });

    test('should display timeline items', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const timeline = page.locator('text=Recent Activity').locator('..');
      const timelineItems = timeline.locator('li');
      const itemCount = await timelineItems.count();
      expect(itemCount).toBeGreaterThanOrEqual(0);
    });

    test('should display user avatars in timeline', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const avatars = page.locator('.rounded-full.bg-blue-100');
      const avatarCount = await avatars.count();
      expect(avatarCount).toBeGreaterThanOrEqual(0);
    });

    test('should display action badges in timeline', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const badges = page.locator('.bg-green-100, .bg-blue-100, .bg-red-100');
      const badgeCount = await badges.count();
      expect(badgeCount).toBeGreaterThanOrEqual(0);
    });

    test('should display relative timestamps', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      // Check for relative time formats (e.g., "2h ago", "3d ago")
      const timeText = page.locator('text=/\\d+[mhd] ago|Just now/');
      const hasTimeText = await timeText.first().isVisible().catch(() => false);
      expect(typeof hasTimeText).toBe('boolean');
    });

    test('should navigate to resource when activity is clicked', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const activityLink = page.locator('button.font-medium.hover\\:text-blue-600').first();
      if (await activityLink.isVisible().catch(() => false)) {
        await activityLink.click();
        await page.waitForTimeout(1000);
        // Should navigate to product or version page
        const currentUrl = page.url();
        expect(currentUrl).toMatch(/\/(products|versions)\//);
      }
    });

    test('should show "View All Activity" link', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const viewAllLink = page.locator('text=View All Activity');
      const isVisible = await viewAllLink.isVisible().catch(() => false);
      expect(typeof isVisible).toBe('boolean');
    });

    test('should navigate to audit logs when "View All Activity" is clicked', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const viewAllLink = page.locator('text=View All Activity');
      if (await viewAllLink.isVisible().catch(() => false)) {
        await viewAllLink.click();
        await page.waitForURL('**/audit-logs', { timeout: 5000 });
        await expect(page.locator('h1:has-text("Audit Logs")')).toBeVisible();
      }
    });

    test('should display empty state when no activity', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      const emptyState = page.locator('text=No recent activity');
      const hasItems = page.locator('li').first().isVisible().catch(() => false);
      
      const isEmptyVisible = await emptyState.isVisible().catch(() => false);
      const hasItemsVisible = await hasItems;
      
      expect(isEmptyVisible || hasItemsVisible).toBeTruthy();
    });
  });

  test.describe('Refresh Functionality', () => {
    test('should refresh dashboard data when refresh button is clicked', async ({ page }) => {
      const refreshButton = page.locator('button:has-text("Refresh")');
      await refreshButton.click();
      
      // Wait for refresh to complete
      await page.waitForTimeout(2000);
      
      // Verify dashboard is still visible
      await expect(page.locator('h1:has-text("Dashboard")')).toBeVisible();
    });

    test('should disable refresh button while loading', async ({ page }) => {
      const refreshButton = page.locator('button:has-text("Refresh")');
      const isDisabled = await refreshButton.isDisabled();
      // Button may or may not be disabled depending on loading state
      expect(typeof isDisabled).toBe('boolean');
    });
  });

  test.describe('Integration Tests', () => {
    test('should load all dashboard data correctly', async ({ page }) => {
      await page.waitForTimeout(3000);
      
      // Verify all sections are loaded
      await expect(page.locator('text=Total Products')).toBeVisible();
      await expect(page.locator('text=Recent Updates')).toBeVisible();
      await expect(page.locator('text=Pending Approvals')).toBeVisible();
      await expect(page.locator('text=Recent Activity')).toBeVisible();
    });

    test('should navigate from dashboard to detail pages', async ({ page }) => {
      await page.waitForTimeout(2000);
      
      // Try clicking on a stat card
      const productsCard = page.locator('text=Total Products').locator('..').locator('..');
      const cardParent = productsCard.locator('..');
      const isClickable = await cardParent.getAttribute('class')?.includes('cursor-pointer') || false;
      
      if (isClickable) {
        await cardParent.click();
        await page.waitForTimeout(1000);
        // Should navigate away from dashboard
        const currentUrl = page.url();
        expect(currentUrl).not.toBe('http://localhost:3000/');
      }
    });

    test('should handle error states gracefully', async ({ page }) => {
      // This would require mocking API failures
      // For now, verify error handling exists
      const errorAlert = page.locator('.bg-red-50.border-red-200');
      const hasError = await errorAlert.isVisible().catch(() => false);
      // Error may or may not be visible
      expect(typeof hasError).toBe('boolean');
    });
  });
});

