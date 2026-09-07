import { chromium } from 'playwright';
import fs from 'fs';

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();
  
  try {
    await page.goto('http://localhost:8080');
    await page.waitForTimeout(2000);
    
    const html = await page.content();
    fs.writeFileSync('dom-dump.html', html);
    
    await page.screenshot({ path: 'screenshot.png' });
    console.log('Successfully captured DOM and screenshot.');
  } catch (err) {
    console.error('Failed to reach local server:', err);
  } finally {
    await browser.close();
  }
})();
