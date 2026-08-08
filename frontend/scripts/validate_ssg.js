import fs from "node:fs";
import path from "node:path";

const backendUrl = process.env.BACKEND_API_URL;
if (!backendUrl) {
  throw new Error("BACKEND_API_URL is required for SSG validation");
}

const headers = {};
if (process.env.NEXT_PUBLIC_SSG_BUILD_TOKEN) {
  headers["x-github-build-token"] = process.env.NEXT_PUBLIC_SSG_BUILD_TOKEN;
}

const response = await fetch(`${backendUrl}/api/regions/combined`, { headers });
if (!response.ok) throw new Error(`Locations API returned ${response.status}`);
const locations = await response.json();
const townships = (locations.cities || []).flatMap((city) =>
  (city.townships || []).map((township) => ({ ...township, city })),
);
if (townships.length !== 368) {
  throw new Error(`Expected 368 township categories, received ${townships.length}`);
}

const dist = path.resolve("dist");
const missing = [];
const invalid = [];
const checkPage = (relativePath, expectedText) => {
  const file = path.join(dist, relativePath, "index.html");
  if (!fs.existsSync(file)) {
    missing.push(relativePath);
    return;
  }
  const html = fs.readFileSync(file, "utf8");
  if (!html.includes(expectedText)) invalid.push(relativePath);
};

for (const city of locations.cities || []) {
  checkPage(path.join("area", String(city.id), "1"), `${city.name}住宿與休息推薦`);
}
for (const township of townships) {
  const pages = Math.max(1, Math.ceil(Number(township.hotel_count || 0) / 20));
  for (let page = 1; page <= pages; page++) {
    checkPage(
      path.join("area", String(township.city.id), String(township.id), String(page)),
      `${township.city.name}${township.name}住宿與休息推薦`,
    );
  }
}

const sitemapFile = path.join(dist, "sitemap.xml");
if (!fs.existsSync(sitemapFile)) {
  missing.push("sitemap.xml");
} else {
  const sitemap = fs.readFileSync(sitemapFile, "utf8");
  for (const township of townships) {
    const route = `/area/${township.city.id}/${township.id}/1`;
    if (!sitemap.includes(route)) invalid.push(`sitemap:${route}`);
  }
}

if (missing.length || invalid.length) {
  throw new Error(
    `SSG validation failed. Missing: ${missing.slice(0, 20).join(", ") || "none"}; ` +
    `invalid: ${invalid.slice(0, 20).join(", ") || "none"}`,
  );
}

console.log(`SSG validation passed: ${townships.length} township categories checked.`);
