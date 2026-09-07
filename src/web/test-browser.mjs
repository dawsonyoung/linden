import { chromium } from '@playwright/test';

(async () => {
  const targetUrl = process.argv[2] || 'http://localhost:8080/';
  let hasErrors = false;

  const browser = await chromium.launch();
  const page = await browser.newPage();

  page.on('console', msg => {
    if (msg.type() === 'error') {
      console.error(`BROWSER ERROR: ${msg.text()}`);
      hasErrors = true;
    }
  });

  page.on('pageerror', error => {
    console.error(`PAGE ERROR: ${error.message}`);
    hasErrors = true;
  });

  page.on('requestfailed', request => {
    console.error(`REQUEST FAILED: ${request.url()} - ${request.failure().errorText}`);
    hasErrors = true;
  });

  try {
    await page.goto(targetUrl, { waitUntil: 'networkidle' });

    const content = await page.content();
    if (content.length < 2000) {
      console.error('DOM content looks suspiciously small, UI might not have rendered.');
      hasErrors = true;
    }

    // Specifically check for the Chat UI mounting
    const chatInput = await page.$('input[placeholder="Type a message..."]');
    if (!chatInput) {
      console.error('Chat input not found! App failed to mount.');
      hasErrors = true;
    }
  } catch (err) {
    console.error('Failed to load page:', err);
    hasErrors = true;
  } finally {
    await browser.close();
  }

  if (hasErrors) {
    process.exit(1);
  }
})();
