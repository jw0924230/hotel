import fs from "node:fs";
import path from "node:path";

const dist = path.resolve("dist");
const indexFile = path.join(dist, "sitemap_index.xml");
const compatibilityFile = path.join(dist, "sitemap.xml");

if (!fs.existsSync(indexFile) || !fs.statSync(indexFile).isFile()) {
  throw new Error("sitemap_index.xml was not generated");
}

if (fs.existsSync(compatibilityFile) && fs.statSync(compatibilityFile).isDirectory()) {
  fs.rmSync(compatibilityFile, { recursive: true, force: true });
}

fs.copyFileSync(indexFile, compatibilityFile);
console.log("Created sitemap.xml compatibility index");
