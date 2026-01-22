#!/usr/bin/env node
"use strict";

const path = require("path");
const { spawn } = require("child_process");

const BIN_NAME = process.platform === "win32" ? "agentx.exe" : "agentx";
const BIN_PATH = path.resolve(__dirname, "..", "vendor", BIN_NAME);

const child = spawn(BIN_PATH, process.argv.slice(2), {
  stdio: "inherit",
});

child.on("exit", (code) => {
  process.exit(code ?? 1);
});

child.on("error", (err) => {
  console.error(`Failed to run AgentX: ${err.message}`);
  process.exit(1);
});
