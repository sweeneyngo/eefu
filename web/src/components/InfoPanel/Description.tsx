import { createMemo, For } from "solid-js";

function Description(props: { description: string | null }) {
  const urlRegex = /(https?:\/\/[^\s]+)/g;
  const lines = createMemo(() => props.description?.split("\n"));
  return (
    <For each={lines()}>
      {(line) =>
        line === "" ? (
          <div class="h-4" /> // preserves empty line spacing
        ) : (
          <div>
            {line.split(urlRegex).map((part) =>
              urlRegex.test(part) ? (
                <a
                  href={part}
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-gray-600 underline hover:text-blue-500"
                >
                  {part}
                </a>
              ) : (
                part
              )
            )}
          </div>
        )
      }
    </For>
  );
}

export function DescriptionCard(props: { description: string | null }) {
  return (
    <div class="relative mt-5 rounded-md bg-gray-100 px-4">
      <p class="flex flex-col py-3 text-gray-600">
        <Description description={props.description} />
      </p>
    </div>
  );
}
