import { createSignal, createContext, useContext } from "solid-js";

export type ToastType = "info" | "success" | "error";

export interface Toast {
  id: number;
  message: string;
  type?: ToastType;
  durationMsec?: number;
}

interface ToastContextValue {
  toasts: Toast[];
  addToast: (message: string, type?: ToastType, duration?: number) => void;
  removeToast: (id: number) => void;
}

const ToastContext = createContext<ToastContextValue>();

export function ToastProvider(props: { children: any }) {
  const [toasts, setToasts] = createSignal<Toast[]>([]);

  const addToast = (
    message: string,
    type: ToastType = "info",
    durationMsec = 3000,
  ) => {
    const id = Date.now() + Math.random();
    const toast: Toast = { id, message, type, durationMsec };
    setToasts((prev) => [...prev, toast]);

    setTimeout(() => removeToast(id), durationMsec);
  };

  const removeToast = (id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  };

  return (
    <ToastContext.Provider value={{ toasts: toasts(), addToast, removeToast }}>
      {props.children}
    </ToastContext.Provider>
  );
}

export const useToast = () => useContext(ToastContext)!;
