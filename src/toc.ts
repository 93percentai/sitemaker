export interface TocHeading {
  level: number;
  id: string;
  text: string;
}

export function extractHeadings(html: string): TocHeading[] {
  const headings: TocHeading[] = [];
  const re = /<h([2-6])\s+id="([^"]*)"[^>]*>(.*?)<\/h\1>/gi;
  let match;

  while ((match = re.exec(html)) !== null) {
    headings.push({
      level: parseInt(match[1], 10),
      id: match[2],
      text: stripHtml(match[3]),
    });
  }

  return headings;
}

function stripHtml(s: string): string {
  return s.replace(/<[^>]+>/g, "");
}

export function renderTOC(headings: TocHeading[]): string {
  if (headings.length === 0) return "";

  let html = '<nav class="toc">\n';
  let prevLevel = 0;

  for (const h of headings) {
    if (prevLevel === 0) {
      html += "<ul>\n";
      prevLevel = h.level;
    }

    if (h.level > prevLevel) {
      for (let i = prevLevel; i < h.level; i++) {
        html += "<ul>\n";
      }
    } else if (h.level < prevLevel) {
      for (let i = h.level; i < prevLevel; i++) {
        html += "</li>\n</ul>\n";
      }
      html += "</li>\n";
    } else if (prevLevel > 0) {
      html += "</li>\n";
    }

    html += `<li><a href="#${h.id}">${h.text}</a>`;
    prevLevel = h.level;
  }

  for (let i = headings[0].level; i <= prevLevel; i++) {
    html += "</li>\n</ul>\n";
  }

  html += "</nav>\n";
  return html;
}
