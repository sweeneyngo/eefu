import { For } from "solid-js";
import type { Tag } from "types/song";

export default function Tags(props: { tags: Tag[] }) {
  return (
    <div class="mt-5 mb-4 flex flex-wrap gap-2">
      <For each={props.tags}>
        {(tag) => (
          <span class="rounded-md bg-gray-200 px-2 py-1 text-xs font-medium text-gray-600">
            {tag.name}
          </span>
        )}
      </For>
    </div>
  );
}
