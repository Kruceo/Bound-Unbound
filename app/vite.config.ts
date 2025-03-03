import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import vike from "vike/plugin";
import { telefunc } from "telefunc/vite";

export default defineConfig({
  plugins: [vike({}), react({}),telefunc()],
  build: {
    target: "es2022",
  },
});
