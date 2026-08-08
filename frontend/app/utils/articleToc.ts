export interface ArticleTocItem {
  id: string
  title: string
}

export interface ArticleWithToc {
  html: string
  items: ArticleTocItem[]
}

const decodeCodePoint = (code: number, fallback: string) => (
  Number.isInteger(code) && code >= 0 && code <= 0x10ffff
    ? String.fromCodePoint(code)
    : fallback
)

const decodeHtmlEntities = (value: string) => value
  .replace(/&#(\d+);/g, (match, code) => decodeCodePoint(Number(code), match))
  .replace(/&#x([\da-f]+);/gi, (match, code) => decodeCodePoint(Number.parseInt(code, 16), match))
  .replace(/&nbsp;/gi, ' ')
  .replace(/&amp;/gi, '&')
  .replace(/&lt;/gi, '<')
  .replace(/&gt;/gi, '>')
  .replace(/&quot;/gi, '"')
  .replace(/&#39;|&apos;/gi, "'")

const headingText = (html: string) => {
  const firstLine = html.split(/<br\s*\/?>/i)[0] || ''
  return decodeHtmlEntities(
    firstLine.replace(/<[^>]*>/g, ' ').replace(/\s+/g, ' ').trim()
  )
}

const headingSlug = (title: string, index: number) => {
  const slug = title
    .normalize('NFKC')
    .toLocaleLowerCase('zh-TW')
    .replace(/[^\p{Letter}\p{Number}]+/gu, '-')
    .replace(/^-+|-+$/g, '')

  return `article-${slug || `section-${index + 1}`}`
}

export const buildArticleToc = (html: string): ArticleWithToc => {
  const items: ArticleTocItem[] = []
  const usedIds = new Set<string>()

  const renderedHtml = html.replace(
    /<h2\b([^>]*)>([\s\S]*?)<\/h2>/gi,
    (_match, rawAttributes: string, innerHtml: string) => {
      const title = headingText(innerHtml)
      if (!title) return _match

      const baseId = headingSlug(title, items.length)
      let id = baseId
      let suffix = 2
      while (usedIds.has(id)) {
        id = `${baseId}-${suffix}`
        suffix += 1
      }
      usedIds.add(id)

      const attributes = rawAttributes.replace(
        /\s+id\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]+)/i,
        ''
      )
      items.push({ id, title })

      return `<h2${attributes} id="${id}" data-article-toc-id="${id}">${innerHtml}</h2>`
    }
  )

  return { html: renderedHtml, items }
}
