import fs from "fs-extra"
import commonjs from "@rollup/plugin-commonjs"
import resolve from "@rollup/plugin-node-resolve"
import livereload from "rollup-plugin-livereload"
import svelte from "rollup-plugin-svelte"
import { terser } from "rollup-plugin-terser"

const production = !process.env.ROLLUP_WATCH

export default {
  input: "src/main.js",
  output: {
    sourcemap: true,
    format: "iife",
    name: "app",
    file: "public/bundle.js",
  },
  plugins: [
    svelte({
      dev: !production,
      css: (css) => css.write("public/bundle.css"),
    }),

    resolve({
      browser: true,
      dedupe: ["svelte"],
    }),
    commonjs(),

    !production && serve(),
    !production && livereload("public"),
    production && terser(),
    production && copyToOut(),
  ],
  watch: {
    clearScreen: false,
  },
}

function copyToOut() {
  return {
    writeBundle() {
      fs.copySync("public", "../out/public")
    },
  }
}

function serve() {
  let server

  function toExit() {
    if (server) server.kill(0)
  }

  return {
    writeBundle() {
      if (server) return
      server = require("child_process").spawn(
        "npm",
        ["run", "start", "--", "--dev"],
        {
          stdio: ["ignore", "inherit", "inherit"],
          shell: true,
        }
      )

      process.on("SIGTERM", toExit)
      process.on("exit", toExit)
    },
  }
}
