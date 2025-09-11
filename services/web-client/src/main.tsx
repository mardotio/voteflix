// sort-imports-ignore
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import { App } from "./App";
import "./syles/reset.scss";
import "./syles/base.scss";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
