/* @refresh reload */
import { render } from "solid-js/web";
import "./index.css";
import App from "./App.tsx";

import { ToastProvider } from "resources/store/Toast/ToastStore";
import { ToastContainer } from "resources/store/Toast/ToastContainer";

const root = document.getElementById("root");

render(
  () => (
    <ToastProvider>
      <App />
      <ToastContainer />
    </ToastProvider>
  ),
  root!,
);
