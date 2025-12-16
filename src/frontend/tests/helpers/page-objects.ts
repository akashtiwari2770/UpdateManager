import { Page, Locator } from '@playwright/test';

export class BasePage {
  protected page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  async goto(path: string) {
    await this.page.goto(path);
  }

  async waitForLoad() {
    await this.page.waitForLoadState('networkidle');
  }
}

export class DashboardPage extends BasePage {
  get heading(): Locator {
    return this.page.locator('h1');
  }

  async navigateToProducts() {
    await this.page.click('text=Products');
  }
}

export class ProductsPage extends BasePage {
  get heading(): Locator {
    return this.page.locator('h1');
  }
}

