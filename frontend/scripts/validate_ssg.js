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

const tagsResponse = await fetch(`${backendUrl}/api/hotel-tags`, { headers });
if (!tagsResponse.ok) throw new Error(`Hotel tags API returned ${tagsResponse.status}`);
const tags = await tagsResponse.json();
const articleTagsResponse = await fetch(`${backendUrl}/api/article-tags`, { headers });
if (!articleTagsResponse.ok) throw new Error(`Article tags API returned ${articleTagsResponse.status}`);
const articleTags = await articleTagsResponse.json();
const hotelsResponse = await fetch(`${backendUrl}/api/hotels?limit=10000`, { headers });
if (!hotelsResponse.ok) throw new Error(`Hotels API returned ${hotelsResponse.status}`);
const hotelsResult = await hotelsResponse.json();
const postsResponse = await fetch(`${backendUrl}/api/posts?limit=10000`, { headers });
if (!postsResponse.ok) throw new Error(`Posts API returned ${postsResponse.status}`);
const postsResult = await postsResponse.json();

const dist = path.resolve("dist");
const missing = [];
const invalid = [];
const canonicalOrigin = "https://www.qk3houronline.com";
const expectedSitemapRoutes = new Set(["/", "/blog"]);
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
  const cityResponse = await fetch(
    `${backendUrl}/api/hotels?limit=1&area=${encodeURIComponent(city.name)}`,
    { headers },
  );
  if (!cityResponse.ok) throw new Error(`Hotels API failed for city ${city.name}`);
  const cityResult = await cityResponse.json();
  const cityPages = Math.max(1, Math.ceil(Number(cityResult.total || 0) / 20));
  for (let page = 1; page <= cityPages; page++) {
    const relativePath = path.join("area", String(city.id), String(page));
    expectedSitemapRoutes.add(`/area/${city.id}/${page}`);
    checkPage(relativePath, `${city.name}住宿與休息推薦`);
    const file = path.join(dist, relativePath, "index.html");
    if (fs.existsSync(file)) {
      const html = fs.readFileSync(file, "utf8");
      const route = `/area/${city.id}/${page}`;
      if (!html.includes(`rel="canonical" href="${canonicalOrigin}${route}"`)) {
        invalid.push(`${relativePath}:canonical`);
      }
      if (html.includes('href="/search/')) invalid.push(`${relativePath}:search-link`);
      if (page === 1) {
        for (const township of city.townships || []) {
          if (!html.includes(`href="/area/${city.id}/${township.id}/1"`)) {
            invalid.push(`${relativePath}:missing-township-${township.id}`);
          }
        }
      }
    }
  }
}
for (const township of townships) {
  const pages = Math.max(1, Math.ceil(Number(township.hotel_count || 0) / 20));
  for (let page = 1; page <= pages; page++) {
    expectedSitemapRoutes.add(`/area/${township.city.id}/${township.id}/${page}`);
    checkPage(
      path.join("area", String(township.city.id), String(township.id), String(page)),
      `${township.city.name}${township.name}住宿與休息推薦`,
    );
    const route = `/area/${township.city.id}/${township.id}/${page}`;
    const file = path.join(dist, route, "index.html");
    if (fs.existsSync(file)) {
      const html = fs.readFileSync(file, "utf8");
      if (!html.includes(`rel="canonical" href="${canonicalOrigin}${route}"`)) {
        invalid.push(`${route}:canonical`);
      }
      if (html.includes('href="/search/')) invalid.push(`${route}:search-link`);
    }
  }
}
for (const tag of tags) {
  const pages = Math.max(1, Math.ceil(Number(tag.enabled_hotel_count || 0) / 20));
  for (let page = 1; page <= pages; page++) {
    expectedSitemapRoutes.add(`/tag/${tag.id}/${page}`);
    const relativePath = path.join("tag", String(tag.id), String(page));
    checkPage(relativePath, `${tag.name}旅館推薦`);
    const file = path.join(dist, relativePath, "index.html");
    if (fs.existsSync(file)) {
      const html = fs.readFileSync(file, "utf8");
      const route = `/tag/${tag.id}/${page}`;
      if (!html.includes(`rel="canonical" href="${canonicalOrigin}${route}"`)) {
        invalid.push(`${relativePath}:canonical`);
      }
    }
  }

  const firstPageFile = path.join(dist, "tag", String(tag.id), "1", "index.html");
  if (fs.existsSync(firstPageFile)) {
    const html = fs.readFileSync(firstPageFile, "utf8");
    if (!html.includes("特色分類")) invalid.push(`tag/${tag.id}/1:missing-feature-category-heading`);
    for (const option of tags) {
      if (!html.includes(`/tag/${option.id}/1`)) {
        invalid.push(`tag/${tag.id}/1:missing-tag-option-${option.id}`);
      }
    }
  }
}
for (const tag of articleTags) {
  const pages = Math.max(1, Math.ceil(Number(tag.post_count || 0) / 12));
  for (let page = 1; page <= pages; page++) {
    expectedSitemapRoutes.add(`/blog/tag/${tag.id}/${page}`);
    const relativePath = path.join("blog", "tag", String(tag.id), String(page));
    checkPage(relativePath, `${tag.name}文章推薦`);
    const file = path.join(dist, relativePath, "index.html");
    if (fs.existsSync(file)) {
      const html = fs.readFileSync(file, "utf8");
      const route = `/blog/tag/${tag.id}/${page}`;
      if (!html.includes(`rel="canonical" href="${canonicalOrigin}${route}"`)) invalid.push(`${relativePath}:canonical`);
      if (page === 1) for (const option of articleTags) if (!html.includes(`/blog/tag/${option.id}/1`)) invalid.push(`${relativePath}:missing-article-tag-${option.id}`);
    }
  }
}
for (const hotel of hotelsResult.data || []) expectedSitemapRoutes.add(`/detail/${hotel.id}`);
for (const post of postsResult.data || []) expectedSitemapRoutes.add(`/blog/${post.id}`);

const decodeXml = (value) => value
  .replaceAll("&amp;", "&")
  .replaceAll("&lt;", "<")
  .replaceAll("&gt;", ">")
  .replaceAll("&quot;", "\"")
  .replaceAll("&apos;", "'");
const extractLocs = (xml) => [...xml.matchAll(/<loc>(.*?)<\/loc>/gs)].map((match) => decodeXml(match[1].trim()));
const readSitemapIndex = (name) => {
  const file = path.join(dist, name);
  if (!fs.existsSync(file) || !fs.statSync(file).isFile()) {
    missing.push(name);
    return [];
  }
  const xml = fs.readFileSync(file, "utf8");
  if (!xml.includes("<sitemapindex")) invalid.push(`${name}:not-sitemap-index`);
  return extractLocs(xml);
};

const compatibilityChildren = readSitemapIndex("sitemap.xml");
const canonicalChildren = readSitemapIndex("sitemap_index.xml");
if (JSON.stringify(compatibilityChildren) !== JSON.stringify(canonicalChildren)) {
  invalid.push("sitemap.xml:index-alias-mismatch");
}

const actualSitemapRoutes = [];
for (const childLocation of canonicalChildren) {
  let childUrl;
  try {
    childUrl = new URL(childLocation);
  } catch {
    invalid.push(`sitemap-index:invalid-child-url-${childLocation}`);
    continue;
  }
  if (childUrl.origin !== canonicalOrigin || !childUrl.pathname.startsWith("/__sitemap__/")) {
    invalid.push(`sitemap-index:invalid-child-location-${childLocation}`);
    continue;
  }
  const childFile = path.resolve(dist, `.${childUrl.pathname}`);
  if (!childFile.startsWith(`${dist}${path.sep}`) || !fs.existsSync(childFile)) {
    missing.push(childUrl.pathname);
    continue;
  }
  if (fs.statSync(childFile).size > 50 * 1024 * 1024) invalid.push(`${childUrl.pathname}:over-50mb`);
  const xml = fs.readFileSync(childFile, "utf8");
  if (!xml.includes("<urlset")) invalid.push(`${childUrl.pathname}:not-urlset`);
  const urls = extractLocs(xml);
  if (urls.length > 10000) invalid.push(`${childUrl.pathname}:over-10000-urls`);
  for (const location of urls) {
    let url;
    try {
      url = new URL(location);
    } catch {
      invalid.push(`${childUrl.pathname}:invalid-url-${location}`);
      continue;
    }
    if (url.origin !== canonicalOrigin || url.protocol !== "https:") {
      invalid.push(`${childUrl.pathname}:non-canonical-url-${location}`);
      continue;
    }
    actualSitemapRoutes.push(url.pathname);
  }
}

const duplicateSitemapRoutes = actualSitemapRoutes.filter((route, index) => actualSitemapRoutes.indexOf(route) !== index);
if (duplicateSitemapRoutes.length) invalid.push(`sitemap:duplicate-${duplicateSitemapRoutes[0]}`);
const actualSitemapRouteSet = new Set(actualSitemapRoutes);
for (const route of expectedSitemapRoutes) {
  if (!actualSitemapRouteSet.has(route)) invalid.push(`sitemap:missing-${route}`);
}
for (const route of actualSitemapRouteSet) {
  if (!expectedSitemapRoutes.has(route)) invalid.push(`sitemap:unexpected-${route}`);
  if (/^\/(search|cms|api)(\/|$)/.test(route) || route.includes("/_payload")) {
    invalid.push(`sitemap:excluded-route-${route}`);
  }
}

const robotsFile = path.join(dist, "robots.txt");
if (!fs.existsSync(robotsFile)) {
  missing.push("robots.txt");
} else if (!fs.readFileSync(robotsFile, "utf8").includes(`Sitemap: ${canonicalOrigin}/sitemap.xml`)) {
  invalid.push("robots.txt:missing-sitemap-index");
}

const latestTag = articleTags.find((tag) => tag.is_system);
const homeFile = path.join(dist, "index.html");
if (latestTag && fs.existsSync(homeFile)) {
  const home = fs.readFileSync(homeFile, "utf8");
  if (!home.includes("最新文章") || !home.includes(`/blog/tag/${latestTag.id}/1`)) invalid.push("home:latest-articles");
}

if (missing.length || invalid.length) {
  throw new Error(
    `SSG validation failed. Missing: ${missing.slice(0, 20).join(", ") || "none"}; ` +
    `invalid: ${invalid.slice(0, 20).join(", ") || "none"}`,
  );
}

console.log(`SSG validation passed: ${actualSitemapRoutes.length} URLs across ${canonicalChildren.length} sitemap files; ${townships.length} townships, ${tags.length} hotel tags, and ${articleTags.length} article tags checked.`);
