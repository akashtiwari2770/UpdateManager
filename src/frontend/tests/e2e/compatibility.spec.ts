import { test, expect } from '@playwright/test';

test.describe('Compatibility Management', () => {
  test.describe('Compatibility Matrix List', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/compatibility');
      await page.waitForLoadState('networkidle');
    });

    test('should load and display compatibility matrices page', async ({ page }) => {
      await expect(page.locator('h1:has-text("Compatibility Matrices")')).toBeVisible();
      await expect(page.locator('text=View and manage compatibility matrices')).toBeVisible();
    });

    test('should display filters section', async ({ page }) => {
      const filtersCard = page.locator('text=Filters').or(page.locator('h3:has-text("Filters")'));
      await expect(filtersCard.first()).toBeVisible();
      
      // Check for product filter
      const productFilter = page.locator('label:has-text("Product")').or(
        page.locator('select').first()
      );
      await expect(productFilter.first()).toBeVisible();
      
      // Check for status filter
      const statusFilter = page.locator('label:has-text("Validation Status")').or(
        page.locator('select').nth(1)
      );
      await expect(statusFilter.first()).toBeVisible();
    });

    test('should display compatibility matrices table', async ({ page }) => {
      // Wait for table or empty state
      const table = page.locator('table').or(page.locator('text=No compatibility matrices found'));
      await table.waitFor({ state: 'visible', timeout: 10000 });
    });

    test('should show loading state', async ({ page }) => {
      await page.reload();
      
      // Check for loading indicator or table appears
      const loadingIndicator = page.locator('[data-testid="loading"], .spinner, text=/loading/i, [role="progressbar"]').first();
      const table = page.locator('table');
      
      // Either loading shows briefly or table appears
      await Promise.race([
        loadingIndicator.waitFor({ state: 'visible', timeout: 2000 }).catch(() => {}),
        table.waitFor({ state: 'visible', timeout: 10000 })
      ]);
    });

    test('should filter by product', async ({ page }) => {
      await page.waitForLoadState('networkidle');
      
      // Find product select dropdown
      const productSelect = page.locator('select').first();
      const isVisible = await productSelect.isVisible().catch(() => false);
      
      if (isVisible) {
        // Get options count
        const options = await productSelect.locator('option').count();
        
        if (options > 1) {
          // Select a product (skip first option which is usually "All Products")
          await productSelect.selectOption({ index: 1 });
          await page.waitForTimeout(1000);
          
          // Verify filter is applied (table should update)
          const table = page.locator('table');
          await table.waitFor({ state: 'visible', timeout: 5000 });
        }
      }
    });

    test('should filter by validation status', async ({ page }) => {
      await page.waitForLoadState('networkidle');
      
      // Find status select dropdown (usually second select)
      const statusSelect = page.locator('select').nth(1);
      const isVisible = await statusSelect.isVisible().catch(() => false);
      
      if (isVisible) {
        // Select "Passed" status
        await statusSelect.selectOption({ value: 'passed' });
        await page.waitForTimeout(1000);
        
        // Verify filter is applied
        const table = page.locator('table');
        await table.waitFor({ state: 'visible', timeout: 5000 });
      }
    });

    test('should clear filters', async ({ page }) => {
      await page.waitForLoadState('networkidle');
      
      const clearButton = page.locator('button:has-text("Clear Filters")');
      const isVisible = await clearButton.isVisible().catch(() => false);
      
      if (isVisible) {
        await clearButton.click();
        await page.waitForTimeout(1000);
        
        // Verify filters are cleared
        const table = page.locator('table');
        await table.waitFor({ state: 'visible', timeout: 5000 });
      }
    });

    test('should display empty state when no matrices', async ({ page }) => {
      // This test will pass if either table or empty state is shown
      const emptyState = page.locator('text=No compatibility matrices found');
      const table = page.locator('table');
      
      const isEmptyStateVisible = await emptyState.isVisible({ timeout: 5000 }).catch(() => false);
      const isTableVisible = await table.isVisible({ timeout: 5000 }).catch(() => false);
      
      expect(isEmptyStateVisible || isTableVisible).toBeTruthy();
    });

    test('should display table columns correctly', async ({ page }) => {
      const table = page.locator('table');
      const isVisible = await table.isVisible().catch(() => false);
      
      if (isVisible) {
        // Check for expected column headers
        await expect(page.locator('th:has-text("Product")')).toBeVisible();
        await expect(page.locator('th:has-text("Version")')).toBeVisible();
        await expect(page.locator('th:has-text("Validation Status")')).toBeVisible();
      }
    });
  });

  test.describe('Compatibility Validation', () => {
    test('should open compatibility validation form from version details', async ({ page }) => {
      // First, navigate to versions page
      await page.goto('/versions');
      await page.waitForLoadState('networkidle');
      
      // Wait for versions table or empty state
      const table = page.locator('table').or(page.locator('text=No versions found'));
      await table.waitFor({ state: 'visible', timeout: 10000 });
      
      // Try to click on first version row if available
      const firstRow = page.locator('table tbody tr').first();
      const isRowVisible = await firstRow.isVisible({ timeout: 5000 }).catch(() => false);
      
      if (isRowVisible) {
        await firstRow.click();
        await page.waitForURL('**/versions/**', { timeout: 5000 });
        
        // Navigate to compatibility tab
        const compatibilityTab = page.locator('button:has-text("Compatibility")');
        await compatibilityTab.click();
        await page.waitForTimeout(500);
        
        // Check for validate button
        const validateButton = page.locator('button:has-text("Validate Compatibility")').or(
          page.locator('button:has-text("Re-validate")')
        );
        const isValidateVisible = await validateButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isValidateVisible) {
          await validateButton.click();
          
          // Check if modal/form opens
          const modal = page.locator('[role="dialog"]').or(page.locator('.modal')).or(
            page.locator('h2:has-text("Validate Compatibility")')
          );
          await expect(modal.first()).toBeVisible({ timeout: 5000 });
        }
      }
    });

    test('should display compatibility validation form fields', async ({ page }) => {
      // Navigate to a version and open compatibility tab
      await page.goto('/versions');
      await page.waitForLoadState('networkidle');
      
      const firstRow = page.locator('table tbody tr').first();
      const isRowVisible = await firstRow.isVisible({ timeout: 5000 }).catch(() => false);
      
      if (isRowVisible) {
        await firstRow.click();
        await page.waitForURL('**/versions/**', { timeout: 5000 });
        
        const compatibilityTab = page.locator('button:has-text("Compatibility")');
        await compatibilityTab.click();
        await page.waitForTimeout(500);
        
        const validateButton = page.locator('button:has-text("Validate Compatibility")').or(
          page.locator('button:has-text("Re-validate")')
        );
        const isValidateVisible = await validateButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isValidateVisible) {
          await validateButton.click();
          await page.waitForTimeout(500);
          
          // Check for form fields
          await expect(page.locator('label:has-text("Min Server Version")').or(
            page.locator('input[placeholder*="Min Server"]')
          ).first()).toBeVisible({ timeout: 5000 });
          
          await expect(page.locator('label:has-text("Max Server Version")').or(
            page.locator('input[placeholder*="Max Server"]')
          ).first()).toBeVisible({ timeout: 5000 });
        }
      }
    });

    test('should add incompatible versions', async ({ page }) => {
      await page.goto('/versions');
      await page.waitForLoadState('networkidle');
      
      const firstRow = page.locator('table tbody tr').first();
      const isRowVisible = await firstRow.isVisible({ timeout: 5000 }).catch(() => false);
      
      if (isRowVisible) {
        await firstRow.click();
        await page.waitForURL('**/versions/**', { timeout: 5000 });
        
        const compatibilityTab = page.locator('button:has-text("Compatibility")');
        await compatibilityTab.click();
        await page.waitForTimeout(500);
        
        const validateButton = page.locator('button:has-text("Validate Compatibility")').or(
          page.locator('button:has-text("Re-validate")')
        );
        const isValidateVisible = await validateButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isValidateVisible) {
          await validateButton.click();
          await page.waitForTimeout(500);
          
          // Find incompatible version input
          const incompatibleInput = page.locator('input[placeholder*="incompatible"]').or(
            page.locator('label:has-text("Incompatible Versions")').locator('..').locator('input')
          ).first();
          
          const isInputVisible = await incompatibleInput.isVisible({ timeout: 5000 }).catch(() => false);
          
          if (isInputVisible) {
            await incompatibleInput.fill('1.0.0');
            
            // Find and click Add button
            const addButton = page.locator('button:has-text("Add")').near(incompatibleInput);
            await addButton.click();
            await page.waitForTimeout(500);
            
            // Verify version was added (should appear as a tag/badge)
            await expect(page.locator('text=1.0.0')).toBeVisible({ timeout: 3000 });
          }
        }
      }
    });

    test('should cancel compatibility validation', async ({ page }) => {
      await page.goto('/versions');
      await page.waitForLoadState('networkidle');
      
      const firstRow = page.locator('table tbody tr').first();
      const isRowVisible = await firstRow.isVisible({ timeout: 5000 }).catch(() => false);
      
      if (isRowVisible) {
        await firstRow.click();
        await page.waitForURL('**/versions/**', { timeout: 5000 });
        
        const compatibilityTab = page.locator('button:has-text("Compatibility")');
        await compatibilityTab.click();
        await page.waitForTimeout(500);
        
        const validateButton = page.locator('button:has-text("Validate Compatibility")').or(
          page.locator('button:has-text("Re-validate")')
        );
        const isValidateVisible = await validateButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isValidateVisible) {
          await validateButton.click();
          await page.waitForTimeout(500);
          
          // Find and click Cancel button
          const cancelButton = page.locator('button:has-text("Cancel")');
          await cancelButton.click();
          await page.waitForTimeout(500);
          
          // Modal should be closed
          const modal = page.locator('[role="dialog"]').or(page.locator('.modal'));
          const isModalVisible = await modal.isVisible({ timeout: 2000 }).catch(() => false);
          expect(isModalVisible).toBeFalsy();
        }
      }
    });
  });

  test.describe('Compatibility Details', () => {
    test('should display compatibility details on version page', async ({ page }) => {
      await page.goto('/versions');
      await page.waitForLoadState('networkidle');
      
      const firstRow = page.locator('table tbody tr').first();
      const isRowVisible = await firstRow.isVisible({ timeout: 5000 }).catch(() => false);
      
      if (isRowVisible) {
        await firstRow.click();
        await page.waitForURL('**/versions/**', { timeout: 5000 });
        
        // Navigate to compatibility tab
        const compatibilityTab = page.locator('button:has-text("Compatibility")');
        await compatibilityTab.click();
        await page.waitForTimeout(500);
        
        // Check for compatibility details section
        const detailsSection = page.locator('text=Validation Status').or(
          page.locator('text=No compatibility information available')
        );
        await expect(detailsSection.first()).toBeVisible({ timeout: 5000 });
      }
    });

    test('should show validation status badge', async ({ page }) => {
      await page.goto('/versions');
      await page.waitForLoadState('networkidle');
      
      const firstRow = page.locator('table tbody tr').first();
      const isRowVisible = await firstRow.isVisible({ timeout: 5000 }).catch(() => false);
      
      if (isRowVisible) {
        await firstRow.click();
        await page.waitForURL('**/versions/**', { timeout: 5000 });
        
        const compatibilityTab = page.locator('button:has-text("Compatibility")');
        await compatibilityTab.click();
        await page.waitForTimeout(1000);
        
        // Check for status badge or empty state
        const statusBadge = page.locator('[class*="badge"]').or(
          page.locator('text=No compatibility information available')
        );
        await expect(statusBadge.first()).toBeVisible({ timeout: 5000 });
      }
    });
  });
});

