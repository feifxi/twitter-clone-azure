// src/theme/applyTheme.ts
import { light, dark } from "./tokens";

type Theme = typeof light;

export function applyTheme(theme: Theme) {
    const root = document.documentElement;
    Object.entries(theme).forEach(([key, value]) => {
        root.style.setProperty(`--color-${key}`, value);
    });
}

// Helper to apply light or dark based on a string
export function setTheme(mode: "light" | "dark") {
    const root = document.documentElement;
    if (mode === "dark") {
        root.classList.add("dark");
        root.classList.remove("light");
    } else {
        root.classList.add("light");
        root.classList.remove("dark");
    }
    applyTheme(mode === "dark" ? dark : light);
}
