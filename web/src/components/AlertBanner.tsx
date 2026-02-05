import { Show, createUniqueId } from "solid-js";
import type { JSX } from "solid-js";
import { Alert } from "components/Icons";

type Variant = "info" | "success" | "warning" | "error";

const variantStyles: Record<
  Variant,
  {
    bg: string;
    border: string;
    color: string;
    defaultIcon: JSX.Element;
  }
> = {
  info: {
    bg: "#eef6ff",
    border: "#60a5fa",
    color: "#084298",
    defaultIcon: <Alert size={20} color="#084298" />,
  },

  success: {
    bg: "#ecfdf5",
    border: "#34d399",
    color: "#065f46",
    defaultIcon: <Alert size={20} color="#065f46" />,
  },

  warning: {
    bg: "#fff7ed",
    border: "#f59e0b",
    color: "#92400e",
    defaultIcon: <Alert size={20} color="#92400e" />,
  },

  error: {
    bg: "#fff1f2",
    border: "#ef4444",
    color: "#7f1d1d",
    defaultIcon: <Alert size={20} color="#7f1d1d" />,
  },
};

export default function AlertBanner(props: {
  variant?: Variant;
  title: JSX.Element;
  description?: JSX.Element;
  icon?: JSX.Element;
  onClose?: () => void;
  dismissible?: boolean;
  class?: string;
}) {
  const styles = () => variantStyles[props.variant ?? "info"];

  const titleId = createUniqueId();
  const descId = props.description ? createUniqueId() : undefined;

  return (
    <div
      role="alert"
      aria-labelledby={titleId}
      aria-describedby={descId}
      class={props.class}
      style={{
        display: "flex",
        gap: "12px",
        "align-items": "flex-start",
        background: styles().bg,
        color: styles().color,
        "border-left": `4px solid ${styles().border}`,
        padding: "12px 16px",
        "border-radius": "6px",
      }}
    >
      <div style={{ "margin-top": "2px", flex: "0 0 auto" }}>
        {props.icon ?? styles().defaultIcon}
      </div>

      <div style={{ flex: "1 1 auto", "min-width": "0" }}>
        <div
          id={titleId}
          style={{
            "font-weight": 600,
            "font-size": "14px",
            "line-height": 1.2,
          }}
        >
          {props.title}
        </div>

        <Show when={props.description}>
          <div
            id={descId}
            style={{
              "margin-top": "4px",
              "font-size": "13px",
              "line-height": 1.3,
              opacity: 0.95,
            }}
          >
            {props.description}
          </div>
        </Show>
      </div>

      <Show when={props.dismissible}>
        <button
          aria-label="Close alert"
          onClick={props.onClose}
          style={{
            "margin-left": "12px",
            background: "transparent",
            border: "none",
            color: styles().color,
            cursor: "pointer",
            padding: "6px",
            display: "inline-flex",
            "align-items": "center",
            "justify-content": "center",
            "border-radius": "4px",
          }}
        >
          <svg
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            aria-hidden
          >
            <path
              d="M18 6L6 18"
              stroke="currentColor"
              stroke-width="1.6"
              stroke-linecap="round"
            />
            <path
              d="M6 6l12 12"
              stroke="currentColor"
              stroke-width="1.6"
              stroke-linecap="round"
            />
          </svg>
        </button>
      </Show>
    </div>
  );
}
