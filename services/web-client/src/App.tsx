import { startTransition, useActionState, useState } from "react";

import "./App.css";
import reactLogo from "./assets/react.svg";
import { type ApiErrorData, type TokenResponse, authApi } from "./utils/client";
import viteLogo from "/vite.svg";

function App() {
  const [count, setCount] = useState(0);
  const [token, setToken] = useState("");
  const [state, formAction, isPending] = useActionState<
    TokenResponse | ApiErrorData | null,
    string
  >(async (_, loginToken) => await authApi.create(loginToken), null);

  return (
    <>
      <div>
        <label>
          Token
          <input onChange={(e) => setToken(e.target.value)} />
        </label>
        <button onClick={() => startTransition(() => formAction(token))}>
          {isPending ? "Submit..." : "Submit"}
        </button>
        <pre
          style={{
            maxWidth: "300px",
            overflowWrap: "anywhere",
            overflow: "hidden",
            whiteSpace: "normal",
          }}
        >
          {state && JSON.stringify(state, null, 2)}
        </pre>
      </div>
      <div>
        <a href="https://vite.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  );
}

export default App;
