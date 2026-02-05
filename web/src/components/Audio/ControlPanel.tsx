import { Pause, Play } from "components/Icons";
import { round, toDuration } from "utils/math";

export default function ControlPanel(props: {
  isPlaying: boolean;
  togglePlay: () => void;
  progress: number;
  currentTime: number;
  duration: number;
  title?: string;
  disabled?: boolean;
}) {

  return (
    <div class="flex items-center space-x-4">
      <button
        onClick={() => props.togglePlay()}
        disabled={props.disabled}
        class={`flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-md border border-gray-300 text-gray-800 ${props.disabled ? "cursor-not-allowed opacity-50" : "cursor-pointer hover:bg-gray-200"}`}
      >
        {props.isPlaying ? <Pause /> : <Play />}
      </button>
      <div class={`flex-1 ${props.disabled ? "opacity-50" : ""}`}>
        <div class="relative h-2 overflow-hidden rounded-full bg-gray-200">
          <div
            class="h-2 rounded-full bg-gray-500"
            style={{ width: `${props.progress * 100 || 0}%` }}
          />
        </div>
        <div class="mt-1 flex w-full items-center justify-between px-1 text-xs text-gray-500">
          <span class="max-w-[12rem] truncate sm:max-w-none md:max-w-[20rem]">
            {props.title}
          </span>
          <span class="ml-2 flex-shrink-0">
            {toDuration(round(props.currentTime))}/
            {toDuration(round(props.duration))}
          </span>
        </div>
      </div>
    </div>
  );
}
