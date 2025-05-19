/* @refresh reload */
import "./index.css";
import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";

import Layout from "./Layout";
import { Dashboard, Nodes, Streams, Stream } from "./views";

const root = document.getElementById("root");

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    "Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got misspelled?",
  );
}

render(
  () => (
    <Router root={Layout}>
      <Route path="/" component={Dashboard} />
      <Route path="/nodes/:index" component={Nodes} />
      <Route path="/streams/:index" component={Streams} />
      <Route path="/stream/:name" component={Stream} />
    </Router>
  ),
  root!,
);
