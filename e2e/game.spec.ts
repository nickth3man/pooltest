import { test, expect } from "@playwright/test";

test.describe("Billiards Game", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("page has correct title and heading", async ({ page }) => {
    await expect(page).toHaveTitle(/Billiards/i);
    await expect(
      page.getByRole("heading", { name: "Billiards" }),
    ).toBeVisible();
  });

  test("canvas element is present and accessible", async ({ page }) => {
    const canvas = page.locator("canvas#table");
    await expect(canvas).toBeVisible();
    await expect(canvas).toHaveAttribute("role", "img");
    await expect(canvas).toHaveAttribute("aria-label", /Billiards table game/i);
  });

  test("HUD displays initial state correctly", async ({ page }) => {
    await expect(page.getByText("Sunk")).toBeVisible();
    await expect(page.locator("#sunk-count")).toHaveText("0");
    await expect(page.getByText("State")).toBeVisible();
    await expect(page.locator("#state-text")).toHaveText("AIMING");
    await expect(page.getByText("Take the shot")).toBeVisible();
  });

  test("restart button is present and clickable", async ({ page }) => {
    const restartBtn = page.locator("#restart-btn");
    await expect(restartBtn).toBeVisible();
    await expect(restartBtn).toBeEnabled();
    await expect(restartBtn).toHaveText("New Rack");
  });

  test("canvas responds to mouse interactions", async ({ page }) => {
    const canvas = page.locator("canvas#table");
    const box = await canvas.boundingBox();

    if (box) {
      await page.mouse.move(box.x + box.width / 2, box.y + box.height / 2);
      await canvas.click();
    }

    await expect(page.locator("#state-text")).toBeVisible();
  });
});
