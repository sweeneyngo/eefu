import type { Song } from "types/song";

interface TagFilter {
  tag: string;
  op?: "=" | ">=" | "<=" | ">" | "<";
  value: string;
  incomplete?: boolean; // optional, for internal use
}

export function parseSearch(query: string): {
  keywords: string[];
  tags: TagFilter[];
} {
  const keywords: string[] = [];
  const tags: TagFilter[] = [];

  const parts = query.split(" ");

  parts.forEach((part) => {
    if (part.startsWith("/")) {
      const match = part.match(/\/(\w+):(>=|<=|=|>|<)?(.+)/);
      if (match) {
        const [, tag, op, value] = match;
        console.log({ tag, op, value });
        tags.push({ tag, op: (op as any) || "=", value, incomplete: false });
      }
    } else if (part.trim()) {
      keywords.push(part.trim());
    }
  });

  return { keywords, tags };
}

export function filterSongs(songs: Song[], query: string) {
  const { keywords, tags } = parseSearch(query);

  return songs.filter((song) => {
    const matchesText = keywords.length
      ? keywords.every((kw) =>
          song.title.toLowerCase().includes(kw.toLowerCase()),
        )
      : true;

    const matchesTags = tags.every((tag) => {
      switch (tag.tag) {
        case "year":
          if (!song.released_at) return false;
          const songYear = new Date(song.released_at).getFullYear();
          const val = parseInt(tag.value, 10);
          switch (tag.op) {
            case "=":
              return songYear === val;
            case ">=":
              return songYear >= val;
            case "<=":
              return songYear <= val;
            case ">":
              return songYear > val;
            case "<":
              return songYear < val;
          }
          break;

        default:
          return true;
      }
    });

    return matchesText && matchesTags;
  });
}
