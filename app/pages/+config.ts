import vikeReact from "vike-react/config";
import type { Config } from "vike/types";
import Layout from "../layouts/LayoutDefault.jsx";

// Default config (can be overridden by pages)
// https://vike.dev/config

export default {
  // https://vike.dev/Layout
  // Layout,
  // https://vike.dev/head-tags
  title: "Bound Unbound",
  description: "Multiple Unbound instances manager.",
  extends: vikeReact,
} satisfies Config;
