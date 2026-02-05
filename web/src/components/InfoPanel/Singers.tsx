import type { SongSinger } from "types/song";
import { Show, createMemo } from "solid-js";

export default function SongSingers(props: { singers: SongSinger[] | null }) {
  const mainSinger = createMemo(
    () => props.singers?.find((s) => s.role === "main")?.singer.name ?? "",
  );

  const supportingSingers = createMemo(
    () =>
      props.singers
        ?.filter((s) => s.role !== "main")
        .map((s) => s.singer.name) ?? [],
  );

  return (
    <Show when={props.singers && props.singers.length > 0}>
      <p>
        <span class="text-gray-500">
          {mainSinger()}
          <Show when={supportingSingers().length > 0}>
            <> (ft. {supportingSingers().join(", ")})</>
          </Show>
        </span>
      </p>
    </Show>
  );
}
