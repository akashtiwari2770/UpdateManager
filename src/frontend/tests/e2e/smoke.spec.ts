import { test, expect } from '@playwright/test';

test.describe('Smoke Tests', () => {
  test('app loads successfully', async ({ page }) => {
    await page.goto('/', { waitUntil: 'domcontentloaded' });
    // Wait for title to be set
    await page.waitForFunction(() => document.title.includes('Update Manager'), { timeout: 10000 });
    await expect(page).toHaveTitle(/Update Manager/, { timeout: 5000 });
  });

  test('navigation works', async ({ page }) => {
    await page.goto('/', { waitUntil: 'domcontentloaded' });
    
    // Wait for sidebar to be ready
    await page.waitForSelector('aside', { state: 'visible', timeout: 10000 });
    await page.waitForSelector('a[href="/products"]', { state: 'visible', timeout: 10000 });

    // Check that dashboard is visible
    const heading = page.locator('h1').first();
    await heading.waitFor({ state: 'visible', timeout: 10000 });
    await expect(heading).toContainText('Dashboard', { timeout: 5000 });

    // Navigate to Products
    const productsLink = page.locator('a:has-text("Products")').first();
    await productsLink.waitFor({ state: 'visible', timeout: 10000 });
    
    await Promise.all([
      page.waitForURL(/.*products/, { timeout: 15000 }),
      productsLink.click(),
    ]);
    
    await page.waitForLoadState('domcontentloaded');
    const productsHeading = page.locator('h1').first();
    await productsHeading.waitFor({ state: 'visible', timeout: 10000 });
    await expect(productsHeading).toContainText('Products', { timeout: 5000 });

    // Navigate to Versions
    const versionsLink = page.locator('a:has-text("Versions")').first();
    await versionsLink.waitFor({ state: 'visible', timeout: 10000 });
    
    await Promise.all([
      page.waitForURL(/.*versions/, { timeout: 15000 }),
      versionsLink.click(),
    ]);
    
    await page.waitForLoadState('domcontentloaded');
    const versionsHeading = page.locator('h1').first();
    await versionsHeading.waitFor({ state: 'visible', timeout: 10000 });
    await expect(versionsHeading).toContainText('Versions', { timeout: 5000 });
  });

  test('health check endpoint accessible', async ({ request }) => {
    const response = await request.get('http://localhost:8080/health');
    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body).toHaveProperty('status');
  });
});

