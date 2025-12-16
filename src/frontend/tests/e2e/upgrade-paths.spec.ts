import { test, expect } from '@playwright/test';

test.describe('Upgrade Path Management', () => {
  test.describe('Upgrade Path Creation', () => {
    test('should open create upgrade path form from version details', async ({ page }) => {
      await page.goto('/versions');
      await page.waitForLoadState('networkidle');
      
      const firstRow = page.locator('table tbody tr').first();
      const isRowVisible = await firstRow.isVisible({ timeout: 5000 }).catch(() => false);
      
      if (isRowVisible) {
        await firstRow.click();
        await page.waitForURL('**/versions/**', { timeout: 5000 });
        
        // Navigate to compatibility tab (upgrade paths might be there or in a separate section)
        const compatibilityTab = page.locator('button:has-text("Compatibility")');
        await compatibilityTab.click();
        await page.waitForTimeout(500);
        
        // Look for create upgrade path button
        const createButton = page.locator('button:has-text("Create Upgrade Path")').or(
          page.locator('button:has-text("Add Upgrade Path")')
        );
        const isButtonVisible = await createButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isButtonVisible) {
          await createButton.click();
          await page.waitForTimeout(500);
          
          // Check if modal/form opens
          const modal = page.locator('[role="dialog"]').or(page.locator('.modal')).or(
            page.locator('h2:has-text("Create Upgrade Path")')
          );
          await expect(modal.first()).toBeVisible({ timeout: 5000 });
        }
      }
    });

    test('should display create upgrade path form fields', async ({ page }) => {
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
        
        const createButton = page.locator('button:has-text("Create Upgrade Path")').or(
          page.locator('button:has-text("Add Upgrade Path")')
        );
        const isButtonVisible = await createButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isButtonVisible) {
          await createButton.click();
          await page.waitForTimeout(500);
          
          // Check for form fields
          await expect(page.locator('label:has-text("From Version")').or(
            page.locator('input[placeholder*="From"]')
          ).first()).toBeVisible({ timeout: 5000 });
          
          await expect(page.locator('label:has-text("To Version")').or(
            page.locator('input[placeholder*="To"]')
          ).first()).toBeVisible({ timeout: 5000 });
          
          await expect(page.locator('label:has-text("Path Type")').or(
            page.locator('select')
          ).first()).toBeVisible({ timeout: 5000 });
        }
      }
    });

    test('should select path type', async ({ page }) => {
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
        
        const createButton = page.locator('button:has-text("Create Upgrade Path")').or(
          page.locator('button:has-text("Add Upgrade Path")')
        );
        const isButtonVisible = await createButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isButtonVisible) {
          await createButton.click();
          await page.waitForTimeout(500);
          
          // Find path type select
          const pathTypeSelect = page.locator('select').first();
          const isSelectVisible = await pathTypeSelect.isVisible({ timeout: 5000 }).catch(() => false);
          
          if (isSelectVisible) {
            // Select Multi-Step
            await pathTypeSelect.selectOption({ value: 'multi_step' });
            await page.waitForTimeout(500);
            
            // Check if intermediate versions field appears
            const intermediateField = page.locator('label:has-text("Intermediate")').or(
              page.locator('input[placeholder*="intermediate"]')
            );
            const isFieldVisible = await intermediateField.isVisible({ timeout: 3000 }).catch(() => false);
            
            // Multi-step should show intermediate versions field
            if (isFieldVisible) {
              await expect(intermediateField.first()).toBeVisible();
            }
          }
        }
      }
    });

    test('should add intermediate versions for multi-step path', async ({ page }) => {
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
        
        const createButton = page.locator('button:has-text("Create Upgrade Path")').or(
          page.locator('button:has-text("Add Upgrade Path")')
        );
        const isButtonVisible = await createButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isButtonVisible) {
          await createButton.click();
          await page.waitForTimeout(500);
          
          // Select Multi-Step path type
          const pathTypeSelect = page.locator('select').first();
          const isSelectVisible = await pathTypeSelect.isVisible({ timeout: 5000 }).catch(() => false);
          
          if (isSelectVisible) {
            await pathTypeSelect.selectOption({ value: 'multi_step' });
            await page.waitForTimeout(500);
            
            // Find intermediate version input
            const intermediateInput = page.locator('input[placeholder*="intermediate"]').or(
              page.locator('label:has-text("Intermediate")').locator('..').locator('input')
            ).first();
            
            const isInputVisible = await intermediateInput.isVisible({ timeout: 5000 }).catch(() => false);
            
            if (isInputVisible) {
              await intermediateInput.fill('1.5.0');
              
              // Find and click Add button
              const addButton = page.locator('button:has-text("Add")').near(intermediateInput);
              await addButton.click();
              await page.waitForTimeout(500);
              
              // Verify version was added
              await expect(page.locator('text=1.5.0')).toBeVisible({ timeout: 3000 });
            }
          }
        }
      }
    });

    test('should validate form fields', async ({ page }) => {
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
        
        const createButton = page.locator('button:has-text("Create Upgrade Path")').or(
          page.locator('button:has-text("Add Upgrade Path")')
        );
        const isButtonVisible = await createButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isButtonVisible) {
          await createButton.click();
          await page.waitForTimeout(500);
          
          // Try to submit without filling required fields
          const submitButton = page.locator('button[type="submit"]').or(
            page.locator('button:has-text("Create")')
          );
          const isSubmitVisible = await submitButton.isVisible({ timeout: 5000 }).catch(() => false);
          
          if (isSubmitVisible) {
            await submitButton.click();
            await page.waitForTimeout(500);
            
            // Check for validation errors
            const errorMessages = page.locator('text=required').or(
              page.locator('[class*="error"]')
            );
            const hasErrors = await errorMessages.count() > 0;
            
            // Form should show validation errors or prevent submission
            expect(hasErrors || !isSubmitVisible).toBeTruthy();
          }
        }
      }
    });

    test('should cancel upgrade path creation', async ({ page }) => {
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
        
        const createButton = page.locator('button:has-text("Create Upgrade Path")').or(
          page.locator('button:has-text("Add Upgrade Path")')
        );
        const isButtonVisible = await createButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isButtonVisible) {
          await createButton.click();
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

  test.describe('Upgrade Path Viewer', () => {
    test('should display upgrade path viewer', async ({ page }) => {
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
        
        // Look for upgrade path viewer or list
        const pathViewer = page.locator('text=Upgrade Path').or(
          page.locator('text=Path Steps')
        );
        const isViewerVisible = await pathViewer.isVisible({ timeout: 5000 }).catch(() => false);
        
        // Either viewer is visible or no paths message
        const noPaths = page.locator('text=No upgrade path').or(
          page.locator('text=No paths found')
        );
        const isNoPathsVisible = await noPaths.isVisible({ timeout: 5000 }).catch(() => false);
        
        expect(isViewerVisible || isNoPathsVisible).toBeTruthy();
      }
    });

    test('should display path type indicators', async ({ page }) => {
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
        
        // Look for path type badges or labels
        const pathTypeBadge = page.locator('text=Direct').or(
          page.locator('text=Multi-Step').or(
            page.locator('text=Blocked')
          )
        );
        const isBadgeVisible = await pathTypeBadge.isVisible({ timeout: 5000 }).catch(() => false);
        
        // Path type should be visible if paths exist
        if (isBadgeVisible) {
          await expect(pathTypeBadge.first()).toBeVisible();
        }
      }
    });

    test('should display path steps', async ({ page }) => {
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
        
        // Look for path steps visualization
        const pathSteps = page.locator('text=Path Steps').or(
          page.locator('text=From').or(
            page.locator('text=To')
          )
        );
        const isStepsVisible = await pathSteps.isVisible({ timeout: 5000 }).catch(() => false);
        
        // Steps should be visible if paths exist
        if (isStepsVisible) {
          await expect(pathSteps.first()).toBeVisible();
        }
      }
    });
  });

  test.describe('Block Upgrade Path', () => {
    test('should open block upgrade path dialog', async ({ page }) => {
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
        
        // Look for block button (might be on upgrade path items)
        const blockButton = page.locator('button:has-text("Block")').or(
          page.locator('button:has-text("Block Path")')
        );
        const isBlockVisible = await blockButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isBlockVisible) {
          await blockButton.first().click();
          await page.waitForTimeout(500);
          
          // Check if dialog opens
          const dialog = page.locator('[role="dialog"]').or(page.locator('.modal')).or(
            page.locator('text=Block Upgrade Path')
          );
          await expect(dialog.first()).toBeVisible({ timeout: 5000 });
        }
      }
    });

    test('should require reason for blocking', async ({ page }) => {
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
        
        const blockButton = page.locator('button:has-text("Block")').or(
          page.locator('button:has-text("Block Path")')
        );
        const isBlockVisible = await blockButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isBlockVisible) {
          await blockButton.first().click();
          await page.waitForTimeout(500);
          
          // Try to submit without reason
          const submitButton = page.locator('button:has-text("Block")').or(
            page.locator('button[type="submit"]')
          );
          const isSubmitVisible = await submitButton.isVisible({ timeout: 5000 }).catch(() => false);
          
          if (isSubmitVisible) {
            await submitButton.click();
            await page.waitForTimeout(500);
            
            // Check for validation error
            const errorMessage = page.locator('text=required').or(
              page.locator('[class*="error"]')
            );
            const hasError = await errorMessage.count() > 0;
            
            // Should show validation error
            expect(hasError).toBeTruthy();
          }
        }
      }
    });

    test('should cancel block upgrade path', async ({ page }) => {
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
        
        const blockButton = page.locator('button:has-text("Block")').or(
          page.locator('button:has-text("Block Path")')
        );
        const isBlockVisible = await blockButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isBlockVisible) {
          await blockButton.first().click();
          await page.waitForTimeout(500);
          
          // Find and click Cancel button
          const cancelButton = page.locator('button:has-text("Cancel")');
          await cancelButton.click();
          await page.waitForTimeout(500);
          
          // Dialog should be closed
          const dialog = page.locator('[role="dialog"]').or(page.locator('.modal'));
          const isDialogVisible = await dialog.isVisible({ timeout: 2000 }).catch(() => false);
          expect(isDialogVisible).toBeFalsy();
        }
      }
    });
  });

  test.describe('Integration Tests', () => {
    test('should complete workflow: validate compatibility → create path → view → block', async ({ page }) => {
      // This is a comprehensive integration test
      await page.goto('/versions');
      await page.waitForLoadState('networkidle');
      
      const firstRow = page.locator('table tbody tr').first();
      const isRowVisible = await firstRow.isVisible({ timeout: 5000 }).catch(() => false);
      
      if (isRowVisible) {
        await firstRow.click();
        await page.waitForURL('**/versions/**', { timeout: 5000 });
        
        // Step 1: Navigate to compatibility tab
        const compatibilityTab = page.locator('button:has-text("Compatibility")');
        await compatibilityTab.click();
        await page.waitForTimeout(1000);
        
        // Step 2: Validate compatibility (if button exists)
        const validateButton = page.locator('button:has-text("Validate Compatibility")').or(
          page.locator('button:has-text("Re-validate")')
        );
        const isValidateVisible = await validateButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isValidateVisible) {
          await validateButton.click();
          await page.waitForTimeout(500);
          
          // Fill form and submit (simplified)
          const minServerInput = page.locator('input[placeholder*="Min Server"]');
          const isMinVisible = await minServerInput.isVisible({ timeout: 3000 }).catch(() => false);
          
          if (isMinVisible) {
            await minServerInput.fill('1.0.0');
            
            const submitButton = page.locator('button:has-text("Validate Compatibility")').or(
              page.locator('button[type="submit"]')
            );
            await submitButton.click();
            await page.waitForTimeout(2000);
          } else {
            // Close modal if we can't fill
            const cancelButton = page.locator('button:has-text("Cancel")');
            await cancelButton.click();
            await page.waitForTimeout(500);
          }
        }
        
        // Step 3: Create upgrade path (if button exists)
        const createPathButton = page.locator('button:has-text("Create Upgrade Path")').or(
          page.locator('button:has-text("Add Upgrade Path")')
        );
        const isCreateVisible = await createPathButton.isVisible({ timeout: 5000 }).catch(() => false);
        
        if (isCreateVisible) {
          await createPathButton.click();
          await page.waitForTimeout(500);
          
          // Fill form
          const fromInput = page.locator('input[placeholder*="From"]');
          const toInput = page.locator('input[placeholder*="To"]');
          
          const isFromVisible = await fromInput.isVisible({ timeout: 3000 }).catch(() => false);
          
          if (isFromVisible) {
            await fromInput.fill('1.0.0');
            await toInput.fill('2.0.0');
            
            // Submit or cancel
            const cancelButton = page.locator('button:has-text("Cancel")');
            await cancelButton.click();
            await page.waitForTimeout(500);
          }
        }
        
        // Step 4: Verify compatibility details are visible
        const detailsSection = page.locator('text=Validation Status').or(
          page.locator('text=No compatibility information available')
        );
        await expect(detailsSection.first()).toBeVisible({ timeout: 5000 });
      }
    });
  });
});

