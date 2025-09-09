import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import { App } from "./App";
import "./syles/base.scss";
import "./syles/reset.scss";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
