import { test, expect } from '@playwright/test';

test.describe('Navigation Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/', { waitUntil: 'domcontentloaded' });
    // Wait for React to hydrate and sidebar to be ready
    await page.waitForSelector('aside', { state: 'visible', timeout: 10000 });
    // Wait for navigation links to be ready
    await page.waitForSelector('a[href="/products"]', { state: 'visible', timeout: 10000 });
    // Small delay to ensure React has finished rendering
    await page.waitForLoadState('networkidle', { timeout: 10000 }).catch(() => {
      // Ignore if networkidle times out, page might be already loaded
    });
  });

  test('navigation between all pages', async ({ page }) => {
    const pages = [
      { name: 'Dashboard', path: '/', heading: 'Dashboard' },
      { name: 'Products', path: '/products', heading: 'Products' },
      { name: 'Versions', path: '/versions', heading: 'Versions' },
      { name: 'Updates', path: '/updates', heading: 'Updates' },
      { name: 'Notifications', path: '/notifications', heading: 'Notifications' },
      { name: 'Audit Logs', path: '/audit-logs', heading: 'Audit Logs' },
    ];

    for (const pageInfo of pages) {
      // Wait for navigation link to be visible and clickable
      const link = page.locator(`a:has-text("${pageInfo.name}")`).first();
      await link.waitFor({ state: 'visible', timeout: 10000 });
      await link.waitFor({ state: 'attached' });
      
      // Click and wait for navigation
      await Promise.all([
        page.waitForURL(new RegExp(`.*${pageInfo.path.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}`), { timeout: 15000 }),
        link.click(),
      ]);
      
      // Wait for page content to load
      await page.waitForLoadState('domcontentloaded');
      
      // Verify heading exists and has correct text
      const heading = page.locator('h1').first();
      await heading.waitFor({ state: 'visible', timeout: 10000 });
      await expect(heading).toContainText(pageInfo.heading, { timeout: 5000 });
    }
  });

  test('active route highlighting', async ({ page }) => {
    // Navigate to Products
    const productsLink = page.locator('a:has-text("Products")').first();
    await productsLink.waitFor({ state: 'visible', timeout: 10000 });
    
    await Promise.all([
      page.waitForURL(/.*products/, { timeout: 15000 }),
      productsLink.click(),
    ]);
    
    // Wait for React to update the active state using waitForFunction
    await page.waitForFunction(
      () => {
        const link = document.querySelector('a[href="/products"]');
        return link && link.className.includes('bg-blue-50');
      },
      { timeout: 5000 }
    );
    
    // Verify active class
    const activeProductsLink = page.locator('a[href="/products"]').first();
    await activeProductsLink.waitFor({ state: 'visible' });
    const classList = await activeProductsLink.evaluate((el) => el.className);
    expect(classList).toContain('bg-blue-50');

    // Navigate to Versions
    const versionsLink = page.locator('a:has-text("Versions")').first();
    await versionsLink.waitFor({ state: 'visible', timeout: 10000 });
    
    await Promise.all([
      page.waitForURL(/.*versions/, { timeout: 15000 }),
      versionsLink.click(),
    ]);
    
    // Wait for React to update the active state
    await page.waitForFunction(
      () => {
        const link = document.querySelector('a[href="/versions"]');
        return link && link.className.includes('bg-blue-50');
      },
      { timeout: 5000 }
    );
    
    // Verify active class
    const activeVersionsLink = page.locator('a[href="/versions"]').first();
    await activeVersionsLink.waitFor({ state: 'visible' });
    const versionsClassList = await activeVersionsLink.evaluate((el) => el.className);
    expect(versionsClassList).toContain('bg-blue-50');
  });

  test('sidebar collapse/expand', async ({ page }) => {
    // Wait for sidebar to be visible
    const sidebar = page.locator('aside').first();
    await sidebar.waitFor({ state: 'visible', timeout: 10000 });
    
    // Find and click the collapse button
    const collapseButton = page.locator('button[aria-label*="sidebar"]').first();
    await collapseButton.waitFor({ state: 'visible', timeout: 10000 });
    
    // Verify sidebar is expanded initially (check width and text visibility)
    const productsLink = page.locator('a[href="/products"]').first();
    await productsLink.waitFor({ state: 'visible', timeout: 10000 });
    
    // Check sidebar width is expanded (w-64 = 256px)
    await page.waitForFunction(
      () => {
        const aside = document.querySelector('aside');
        return aside && aside.offsetWidth > 200;
      },
      { timeout: 5000 }
    );
    let sidebarWidth = await sidebar.evaluate((el) => el.offsetWidth);
    expect(sidebarWidth).toBeGreaterThan(200); // Should be around 256px (w-64)
    
    // Check that text is visible
    const initialText = productsLink.locator('span:has-text("Products")');
    await expect(initialText).toBeVisible({ timeout: 5000 });
    
    // Click to collapse
    await collapseButton.click();

    // Wait for collapse animation and React state update using waitForFunction
    await page.waitForFunction(
      () => {
        const aside = document.querySelector('aside');
        const link = document.querySelector('a[href="/products"]');
        if (!aside || !link) return false;
        // Check if text span exists (when collapsed, it's conditionally rendered out)
        const spans = link.querySelectorAll('span');
        const hasTextSpan = Array.from(spans).some(span => span.textContent?.trim() === 'Products');
        return aside.offsetWidth < 100 && !hasTextSpan;
      },
      { timeout: 5000 }
    );
    
    // Verify collapsed state
    sidebarWidth = await sidebar.evaluate((el) => el.offsetWidth);
    expect(sidebarWidth).toBeLessThan(100); // Should be around 64px (w-16)
    
    // Check that text span doesn't exist (conditionally rendered)
    const collapsedText = productsLink.locator('span:has-text("Products")');
    await expect(collapsedText).toHaveCount(0, { timeout: 2000 });

    // Expand sidebar
    await collapseButton.click();
    
    // Wait for expand animation and React state update
    await page.waitForFunction(
      () => {
        const aside = document.querySelector('aside');
        const link = document.querySelector('a[href="/products"]');
        if (!aside || !link) return false;
        // Check if text span exists (when expanded, it's rendered)
        const spans = link.querySelectorAll('span');
        const hasTextSpan = Array.from(spans).some(span => span.textContent?.trim() === 'Products');
        return aside.offsetWidth > 200 && hasTextSpan;
      },
      { timeout: 5000 }
    );
    
    // Verify expanded state
    sidebarWidth = await sidebar.evaluate((el) => el.offsetWidth);
    expect(sidebarWidth).toBeGreaterThan(200);
    
    // Check that text is visible again
    const expandedText = productsLink.locator('span:has-text("Products")');
    await expect(expandedText).toBeVisible({ timeout: 5000 });
  });
});

