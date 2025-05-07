import vikeReact from "vike-react/config";
import type { Config } from "vike/types";

// Default config (can be overridden by pages)
// https://vike.dev/config

export default {
  title: "Bound Unbound",
  description: "Multiple Unbound instances manager.",
  extends: vikeReact,
  passToClient: [
    'user'
  ],
  meta: {
    onBeforeRender: { env: { client: false, server: true } }
  }
} satisfies Config;
