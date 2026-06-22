import { defineConfig } from "vitest/config";
import solid from "vite-plugin-solid";

export default defineConfig({
  plugins: [solid()],
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./src/test/setup.ts"],
    // Vitest 4: optimizer key is "ssr" (not "web"); transformAssets/CSS are
    // now first-class via deps.web.* flags. solid-js must be pre-bundled.
    deps: {
      optimizer: {
        ssr: {
          include: ["solid-js", "@solidjs/testing-library"],
        },
      },
    },
    server: {
      deps: {
        // Inline solid-js so Vite processes its ESM instead of passing raw to Node
        inline: ["solid-js", "@solidjs/testing-library"],
      },
    },
  },
  resolve: {
    conditions: ["development", "browser"],
  },
});
