import type { JSX } from "solid-js/h/jsx-runtime";

export default function InfoGridRow(props: {
  label: string;
  content: JSX.Element | any;
}) {
  return (
    <>
      <div class="flex h-[24px] items-center justify-center rounded-xs bg-gray-500 px-2 text-white">
        {props.label}
      </div>
      <div class="flex items-center break-all">{props.content}</div>
    </>
  );
}
