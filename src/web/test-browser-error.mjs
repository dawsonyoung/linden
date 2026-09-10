import { chromium } from '@playwright/test';

(async () => {
  const targetUrl = process.argv[2] || 'http://localhost:8080/';
  let hasErrors = false;

  const browser = await chromium.launch();
  const page = await browser.newPage();

  try {
    await page.goto(targetUrl, { waitUntil: 'networkidle' });

    // Ensure the app mounted
    const chatInput = await page.$('input[placeholder="Type a message..."]');
    if (!chatInput) {
      console.error('Chat input not found!');
      process.exit(1);
    }

    // Type a message and send it
    await chatInput.fill('Force error');
    const sendButton = await page.$('button[type="submit"]');
    await sendButton.click();

    // Wait for the alert-error banner to appear
    try {
      await page.waitForSelector('.alert-error', { timeout: 3000 });
      const errorText = await page.$eval('.alert-error span', el => el.textContent);
      console.log(`Successfully caught error in UI: ${errorText}`);
      
      // We expect it to NOT be raw JSON.
      if (errorText.includes('"error":')) {
        console.error('Error banner contains raw JSON instead of parsed text!');
        hasErrors = true;
      }
    } catch (err) {
      console.error('Error banner did not appear within timeout!');
      hasErrors = true;
    }

  } catch (err) {
    console.error('Failed test:', err);
    hasErrors = true;
  } finally {
    await browser.close();
  }

  if (hasErrors) {
    process.exit(1);
  }
})();
