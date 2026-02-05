import { createEffect, createMemo, createSignal, For, Show } from "solid-js";
import { fetchSongs } from "api/Song";
import AlertBanner from "components/AlertBanner";
import SearchBar from "components/Searchbar";
import type { Song, Tag } from "types/song";
import { formatDate } from "utils/date";
import { filterSongs, parseSearch } from "utils/search";

const TENTATIVE = "tentative";
function groupSongsByLatest(songs: Song[]): {
  nonTentative: Song[];
  tentative: Song[];
} {
  const latestByGroup = new Map<string, Song>();
  for (const song of songs) {
    const existing = latestByGroup.get(song.song_group_hash_id);
    if (!existing || song.version > existing.version) {
      latestByGroup.set(song.song_group_hash_id, song);
    }
  }

  const values = Array.from(latestByGroup.values());
  const nonTentative: Song[] = [];
  const tentative: Song[] = [];

  for (const song of values) {
    const hasTentative = Array.isArray((song as Song).tags)
      ? (song as Song).tags.some((t: Tag) => t?.name === TENTATIVE)
      : false;
    if (hasTentative) tentative.push(song);
    else nonTentative.push(song);
  }

  const sortByDateDesc = (a: Song, b: Song) => {
    const dateA = a.released_at ? new Date(a.released_at).getTime() : 0;
    const dateB = b.released_at ? new Date(b.released_at).getTime() : 0;
    return dateB - dateA;
  };

  nonTentative.sort(sortByDateDesc);
  tentative.sort(sortByDateDesc);

  return {
    nonTentative,
    tentative,
  };
}

export default function HomePage() {
  const [nonTentativeSongs, setNonTentativeSongs] = createSignal<Song[]>([]);
  const [tentativeSongs, setTentativeSongs] = createSignal<Song[]>([]);
  const [loaded, setLoaded] = createSignal(false);
  const [error, setError] = createSignal(false);

  const [search, setSearch] = createSignal(""); // Step 1: track search input

  const highlightMatch = (text: string, query: string) => {
    const { keywords, _ } = parseSearch(query);
    const parts: any[] = [];

    let lastIndex = 0;
    keywords.forEach((kw) => {
      const regex = new RegExp(`(${kw})`, "gi");
      let match;
      while ((match = regex.exec(text)) !== null) {
        if (match.index > lastIndex)
          parts.push(<span>{text.slice(lastIndex, match.index)}</span>);
        parts.push(<span class="bg-yellow-200">{match[0]}</span>);
        lastIndex = match.index + match[0].length;
      }
    });

    if (lastIndex < text.length)
      parts.push(<span>{text.slice(lastIndex)}</span>);

    return parts.length ? parts : text;
  };

  createEffect(async () => {
    setLoaded(false);

    try {
      const data = await fetchSongs();
      console.log(data);
      if (!data || data.length === 0) {
        setError(true);
        return;
      }
      const { nonTentative, tentative } = groupSongsByLatest(data);
      setNonTentativeSongs(nonTentative);
      setTentativeSongs(tentative);
    } catch (err) {
      console.error(err);
      setError(true);
    } finally {
      setLoaded(true);
    }
  });

  const filteredNonTentativeSongs = createMemo(() => {
    return filterSongs(nonTentativeSongs(), search());
  });

  const filteredTentativeSongs = createMemo(() => {
    return filterSongs(tentativeSongs(), search());
  });

  return (
    <div class="mt-4">
      {error() && (
        <AlertBanner
          variant="error"
          title="Server failed to respond!"
          description="Please try refreshing the page, or check back later."
        />
      )}

      <Show when={loaded()}>
        <div>
          <SearchBar value={search()} onInput={setSearch} />
        </div>
        <div class="mt-4 grid grid-cols-1 gap-2">
          <div>
            <For each={filteredNonTentativeSongs()}>
              {(song) => {
                return (
                  <div class="bg-white pl-1 transition">
                    <h2 class="flex items-center text-sm text-gray-700">
                      <a
                        class="hover:underline"
                        href={`/song/${song.song_group_hash_id}`}
                      >
                        {highlightMatch(song.title, search())}
                      </a>
                      <span class="ml-2 text-xs text-gray-400">
                        {formatDate((song as any).released_at ?? null)}
                      </span>
                    </h2>
                  </div>
                );
              }}
            </For>
          </div>
          <div>
            <For each={filteredTentativeSongs()}>
              {(song) => {
                return (
                  <div class="bg-white pl-1 transition">
                    <h2 class="flex items-center text-sm text-gray-500">
                      <a
                        class="hover:underline"
                        href={`/song/${song.song_group_hash_id}`}
                      >
                        {highlightMatch(song.title, search())}
                      </a>
                      <span class="ml-2 text-xs text-gray-400">
                        {formatDate((song as any).released_at ?? null)}
                      </span>
                    </h2>
                  </div>
                );
              }}
            </For>
          </div>
        </div>
      </Show>
    </div>
  );
}
