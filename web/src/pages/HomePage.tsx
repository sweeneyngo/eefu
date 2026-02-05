import { createEffect, createSignal, For, Show } from "solid-js";
import { fetchSongs } from "api/Song";
import type { Song, Tag } from "types/song";
import { formatDate } from "utils/date";

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

  createEffect(async () => {
    setLoaded(false);

    try {
      const data = await fetchSongs();
      if (!data || data.length === 0) return;
      const { nonTentative, tentative } = groupSongsByLatest(data);
      setNonTentativeSongs(nonTentative);
      setTentativeSongs(tentative);
    } catch (err) {
      console.error(err);
    } finally {
      setLoaded(true);
    }
  });

  return (
    <div class="mt-4">
      <Show when={loaded()}>
        <div class="grid grid-cols-1 gap-2">
          <div>
            <For each={nonTentativeSongs()}>
              {(song) => {
                return (
                  <div class="bg-white pl-1 transition">
                    <h2 class={"flex items-center text-sm text-gray-700"}>
                      <a
                        class="hover:underline"
                        href={`/song/${song.song_group_hash_id}`}
                      >
                        {song.title}
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
            <For each={tentativeSongs()}>
              {(song) => {
                return (
                  <div class="bg-white pl-1 transition">
                    <h2 class={"flex items-center text-sm text-gray-500"}>
                      <a
                        class="hover:underline"
                        href={`/song/${song.song_group_hash_id}`}
                      >
                        {song.title}
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
