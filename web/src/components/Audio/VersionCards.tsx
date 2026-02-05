import { For, Show } from "solid-js";
import type { Song } from "types/song";
import { formatDate } from "utils/date";

export function VersionCards(props: {
  song: Song[];
  version: number;
  onSelectVersion: (version: number) => void;
}) {
  return (
    <div class="-mt-1 flex flex-col">
      <Show when={props.song}>
        <For each={props.song}>
          {(song, i) => (
            <div
              class={`cursor-pointer border border-gray-300 p-2 transition ${props.version === i() ? "z-10 border-blue-400 bg-gray-100" : "bg-white hover:bg-gray-50"} ${i() !== 0 ? "-mt-[1px]" : ""} ${i() === props.song.length - 1 ? "rounded-b-md" : "rounded-b-none"}`}
              onClick={() => props.onSelectVersion(i())}
              style={{ "z-index": props.version === i() ? 20 : 0 }}
            >
              <div class="mb-1 flex items-center justify-between px-3">
                <span class="text-sm">{`v${song.version}`}</span>
                <span class="text-xs text-gray-500">
                  {formatDate(song.released_at)}
                </span>
              </div>
            </div>
          )}
        </For>
      </Show>
    </div>
  );
}
