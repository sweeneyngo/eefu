import { createSignal } from "solid-js";

export default function CoverArt(props: { src: string }) {
  const [loaded, setLoaded] = createSignal(false);

  return (
    <div class="relative mb-2 aspect-square w-full overflow-hidden rounded-md">
      {!loaded() && <div class="absolute inset-0 animate-pulse bg-gray-200" />}
      <img
        src={props.src}
        class={`h-full w-full rounded-md object-cover transition-opacity duration-200 ${
          loaded() ? "opacity-100" : "opacity-0"
        }`}
        onLoad={() => setLoaded(true)}
      />
    </div>
  );
}
