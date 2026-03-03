import { readFileSync, writeFileSync, existsSync } from "node:fs";
import { resolve } from "node:path";

const target = resolve("node_modules/react-hook-form/dist/index.d.ts");

if (!existsSync(target)) {
  process.exit(0);
}

const original = readFileSync(target, "utf8");
const fixed = original.replaceAll("../src/", "./");

if (fixed !== original) {
  writeFileSync(target, fixed, "utf8");
  console.log("patched react-hook-form dist/index.d.ts exports");
}
