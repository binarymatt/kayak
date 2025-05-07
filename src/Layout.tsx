import { Component, For, PropsWithChildren } from "solid-js";
import { A } from "@solidjs/router";
import { clusters } from "./stores";

const Layout: Component = (props: PropsWithChildren) => {
  return (
    <div class="drawer">
      <input id="my-drawer" type="checkbox" class="drawer-toggle" />
      <div class="drawer-content">
        <div class="flex flex-col h-screen justify-between">
          <header class="">
            <div class="navbar bg-base-300">
              <div class="flex-none">
                <label
                  for="my-drawer"
                  aria-label="open sidebar"
                  class="btn btn-square btn-ghost"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    class="inline-block h-5 w-5 stroke-current"
                  >
                    {" "}
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M4 6h16M4 12h16M4 18h16"
                    ></path>{" "}
                  </svg>
                </label>
              </div>
              <div class="flex-1">
                <a href="/" class="btn btn-ghost normal-case text-xl">
                  Kayak UI
                </a>
              </div>
              <div class="flex-2"></div>
              <div class="flex-none gap-2"></div>
            </div>
          </header>
          <main class="m-4 mb-auto">
            <div class="object-center">{props.children}</div>
          </main>
          <footer class="footer bg-neutral text-neutral-content p-4">
            <input
              type="checkbox"
              class="toggle"
              data-toggle-theme="dark,light"
              data-act-class="ACTIVECLASS"
            />
            <div class="items-right grid-flow-col">
              <p>Copyright © 2025 - All right reserved</p>
            </div>
          </footer>
        </div>
      </div>
      <div class="drawer-side">
        <label
          for="my-drawer"
          aria-label="close sidebar"
          class="drawer-overlay"
        ></label>
        <ul class="menu bg-base-200 text-base-content min-h-full w-80 p-4">
          <For each={clusters}>
            {(cluster, i) => (
              <li>
                <h2 class="menu-title">{cluster.name}</h2>
                <ul>
                  <li>
                    <A
                      href={"/nodes/" + i()}
                      onClick={() => {
                        document.getElementById("my-drawer").click();
                      }}
                    >
                      Nodes
                    </A>
                  </li>
                  <li>
                    <A
                      href={"/streams/" + i()}
                      onClick={() => {
                        document.getElementById("my-drawer").click();
                      }}
                    >
                      Streams
                    </A>
                  </li>
                </ul>
              </li>
            )}
          </For>
        </ul>
      </div>
    </div>
  );
};
export default Layout;
