import { createMemo, createSignal, onMount } from "solid-js";
import { useNavigate, useParams } from "@solidjs/router";
import type { SongThumbnailWindow } from "types/song";

export default function SongCarousel(props: { items: SongThumbnailWindow }) {
  const navigate = useNavigate();
  const params = useParams();
  const [mounted, setMounted] = createSignal(false);
  const visibleCount = 5;

  const currentSongID = () => params.song_group_hash_id;

  const handleClick = (item: (typeof props.items.songs)[0]) => {
    navigate(`/song/${item.song_group_hash_id}`);
  };

  const startIndex = createMemo(() => {
    const id = params.song_group_hash_id;
    if (!id) return 0;
    const index = props.items.songs.findIndex((s) => s.song_group_hash_id === id);
    if (index === -1) return 0;
    let offset = index - Math.floor(visibleCount / 2);
    if (offset < 0) offset = 0;
    if (offset + visibleCount > props.items.songs.length) {
      offset = Math.max(props.items.songs.length - visibleCount, 0);
    }
    return offset;
  });

  const visibleItems = createMemo(() =>
    props.items.songs.slice(startIndex(), startIndex() + visibleCount)
  );

  onMount(() => {
      // Preload all images in carousel
      visibleItems().forEach((item) => {
        const img = new Image();
        img.src = item.cover_url;
      });
      setMounted(true)
  });


  return (
    <div class="relative inline-flex gap-2 overflow-hidden">
      <div
        class={`flex gap-2 transition-all duration-700 md:mx-4 transform ${
          mounted() ? "opacity-100 translate-x-0" : "opacity-0 -translate-x-10"
        }`}
      >
        {visibleItems().map((item) => {
          const isCurrent = item.song_group_hash_id === currentSongID();
          return (
            <div
              class={`relative aspect-square max-h-[90px] max-w-[90px] cursor-pointer overflow-hidden rounded-md border border-gray-300 transition-transform duration-300 ${
                isCurrent ? "" : "hover:scale-105"
              }`}
              onClick={() => handleClick(item)}
            >
              <img
                src={item.cover_url}
                class={`h-full w-full object-cover scale-200 brightness-70 filter hover:grayscale-0 ${
                  isCurrent ? "grayscale-0" : "grayscale"
                }`}
                loading="lazy"
              />
              {isCurrent && <div class="absolute bottom-0 left-0 h-1 w-full bg-blue-400" />}
            </div>
          );
        })}
      </div>
    </div>
  );
}
