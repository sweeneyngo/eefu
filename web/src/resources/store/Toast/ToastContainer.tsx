import { For } from "solid-js";
import { useToast } from "resources/store/Toast/ToastStore";

export function ToastContainer() {
  const { toasts, removeToast } = useToast();

  return (
    <div class="fixed top-4 right-4 z-50 flex flex-col gap-2">
      <For each={toasts}>
        {(toast) => (
          <div
            class={`rounded px-4 py-2 text-white shadow ${
              toast.type === "success"
                ? "bg-green-500"
                : toast.type === "error"
                  ? "bg-red-500"
                  : "bg-gray-700"
            }`}
          >
            <div class="flex items-center justify-between gap-2">
              <span>{toast.message}</span>
              <button onClick={() => removeToast(toast.id)} class="ml-2">
                x
              </button>
            </div>
          </div>
        )}
      </For>
    </div>
  );
}
