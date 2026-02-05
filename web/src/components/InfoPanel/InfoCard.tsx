import { createMemo, Show } from "solid-js";
import type { Song } from "types/song";
import { Clipboard, Download } from "components/Icons";
import { DescriptionCard } from "components/InfoPanel/Description";
import Duration from "components/InfoPanel/Duration";
import InfoGridRow from "components/InfoPanel/InfoGridRow";
import MediaTypeIcon from "components/InfoPanel/MediaTypeIcon";
import SongSingers from "components/InfoPanel/Singers";
import Tags from "components/InfoPanel/Tags";
import Tooltip from "components/Tooltip";
import { formatSampleRate } from "utils/audio";
import { toUpper } from "utils/string";

// TODO: Update component to handle Audio & Video
export default function InfoCard(props: {
  song: Song;
  onDownload: () => void;
}) {
  const media = createMemo(() =>
    props.song.media_sources.find((m) => m.file_type === "audio"),
  );

  const handleCopy = (text?: string) => {
    if (!text) return;
    navigator.clipboard.writeText(text);
  };

  return (
    <aside class="mb-4 min-w-0 flex-1 space-y-4 md:mx-4">
      <div class="mb-5">
        <div class="mb-3 flex gap-2">
          <span class="text-xs text-gray-500 hover:underline">
            <a href="/">back to selection</a>
          </span>
        </div>
        <h2 class="text-2xl font-semibold text-gray-800">{props.song.title}</h2>
        <SongSingers singers={props.song.singers} />
        <MediaTypeIcon type={media()?.file_type} />
        <Show when={props.song.description}>
          <DescriptionCard description={props.song.description} />
        </Show>
        <Show when={props.song.tags?.length}>
          <Tags tags={props.song.tags} />
        </Show>
      </div>
      <div class="grid grid-cols-[auto_1fr] gap-x-2 gap-y-2 text-sm text-gray-700">
        <InfoGridRow
          label="Duration"
          content={
            <Duration durationInSec={media()?.audio_metadata?.duration ?? 0} />
          }
        />
        <InfoGridRow
          label="Bitrate"
          content={formatSampleRate(media()?.audio_metadata?.sample_rate)}
        />
        <InfoGridRow
          label="Format"
          content={
            <div class="flex items-center gap-1">
              <div>{toUpper(media()?.format_type)}</div>{" "}
              <div class="cursor-pointer">
                <Tooltip text="Download file">
                  <div
                    class="rounded p-1 hover:bg-gray-100 active:bg-gray-200"
                    onClick={() => props.onDownload()}
                  >
                    <Download size={14} />
                  </div>
                </Tooltip>
              </div>
            </div>
          }
        />
        <InfoGridRow
          label="Checksum"
          content={
            <div class="flex items-center gap-1">
              <div>SHA-256</div>{" "}
              <div class="cursor-pointer">
                <Tooltip text="Copy to clipboard">
                  <div
                    class="rounded p-1 hover:bg-gray-100 active:bg-gray-200"
                    onClick={() => handleCopy(media()?.checksum)}
                  >
                    <Clipboard size={14} />
                  </div>
                </Tooltip>
              </div>
            </div>
          }
        />
      </div>
      <hr class="my-4 border-gray-200" />
    </aside>
  );
}
