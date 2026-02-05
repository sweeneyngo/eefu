import type { JSX } from "solid-js";

export default function Tooltip(props: {
  text: string;
  position?: "top" | "bottom" | "left" | "right";
  children: JSX.Element;
}) {
  const positionClasses: Record<string, string> = {
    top: "bottom-full mb-1 left-1/2 -translate-x-1/2",
    bottom: "top-full mt-1 left-1/2 -translate-x-1/2",
    left: "right-full mr-1 top-1/2 -translate-y-1/2",
    right: "left-full ml-1 top-1/2 -translate-y-1/2",
  };

  return (
    <div class="group relative">
      {props.children}
      <span
        class={`absolute z-50 ${props.position ? positionClasses[props.position] : positionClasses.top} pointer-events-none rounded bg-gray-700 px-2 py-1 text-xs whitespace-nowrap text-white opacity-0 transition-opacity group-hover:opacity-100`}
      >
        {props.text}
      </span>
    </div>
  );
}
