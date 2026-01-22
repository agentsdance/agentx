#!/usr/bin/env node
/* eslint-disable no-console */
"use strict";

const fs = require("fs");
const os = require("os");
const path = require("path");
const https = require("https");
const { pipeline } = require("stream");
const { promisify } = require("util");
const { spawn } = require("child_process");

const streamPipeline = promisify(pipeline);

const PKG_ROOT = path.resolve(__dirname, "..");
const BIN_DIR = path.join(PKG_ROOT, "vendor");
const BIN_NAME = os.platform() === "win32" ? "agentx.exe" : "agentx";
const BIN_PATH = path.join(BIN_DIR, BIN_NAME);

const VERSION = require(path.join(PKG_ROOT, "package.json")).version;
const REPO = "agentsdance/agentx";

function getPlatform() {
  const platform = os.platform();
  if (platform === "darwin" || platform === "linux" || platform === "win32") {
    return platform;
  }
  throw new Error(`Unsupported platform: ${platform}`);
}

function getArch() {
  const arch = os.arch();
  if (arch === "x64") return "amd64";
  if (arch === "arm64") return "arm64";
  throw new Error(`Unsupported architecture: ${arch}`);
}

function getAssetInfo() {
  const platform = getPlatform();
  const arch = getArch();
  const ext = platform === "win32" ? "zip" : "tar.gz";
  const osName = platform === "win32" ? "windows" : platform;
  const filename = `agentx_${VERSION}_${osName}_${arch}.${ext}`;
  return { platform, arch, ext, filename, osName };
}

async function ensureDir(dir) {
  await fs.promises.mkdir(dir, { recursive: true });
}

async function downloadTo(url, dest, redirects = 0) {
  await new Promise((resolve, reject) => {
    https
      .get(url, (res) => {
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          if (redirects >= 5) {
            reject(new Error("Too many redirects while downloading."));
            res.resume();
            return;
          }
          res.resume();
          const nextUrl = new URL(res.headers.location, url).toString();
          downloadTo(nextUrl, dest, redirects + 1).then(resolve).catch(reject);
          return;
        }
        if (res.statusCode !== 200) {
          const err = new Error(
            `Failed to download ${url} (status ${res.statusCode})`
          );
          err.statusCode = res.statusCode;
          reject(err);
          res.resume();
          return;
        }
        const file = fs.createWriteStream(dest);
        streamPipeline(res, file).then(resolve).catch(reject);
      })
      .on("error", reject);
  });
}

async function extractArchive(archivePath, ext) {
  if (ext === "zip") {
    if (os.platform() === "win32") {
      const escapedArchive = archivePath.replace(/'/g, "''");
      const escapedDest = BIN_DIR.replace(/'/g, "''");
      await runCommand("powershell", [
        "-NoProfile",
        "-Command",
        `Expand-Archive -LiteralPath '${escapedArchive}' -DestinationPath '${escapedDest}' -Force`,
      ]);
      return;
    }
    await runCommand("unzip", ["-o", archivePath, "-d", BIN_DIR]);
    return;
  }
  await runCommand("tar", ["-xzf", archivePath, "-C", BIN_DIR]);
}

async function findInstalledBinary() {
  const directPath = path.join(BIN_DIR, BIN_NAME);
  try {
    await fs.promises.access(directPath, fs.constants.X_OK);
    return directPath;
  } catch (_) {
    // continue
  }
  // Look for nested path like vendor/agentx_<ver>_<os>_<arch>/agentx
  const entries = await fs.promises.readdir(BIN_DIR, { withFileTypes: true });
  for (const entry of entries) {
    if (!entry.isDirectory()) continue;
    const candidate = path.join(BIN_DIR, entry.name, BIN_NAME);
    try {
      await fs.promises.access(candidate, fs.constants.X_OK);
      return candidate;
    } catch (_) {
      // continue
    }
  }
  return null;
}

async function promoteBinary(foundPath) {
  if (foundPath === BIN_PATH) return;
  await fs.promises.copyFile(foundPath, BIN_PATH);
}

async function makeExecutable(filePath) {
  if (os.platform() === "win32") return;
  await fs.promises.chmod(filePath, 0o755);
}

function runCommand(cmd, args) {
  return new Promise((resolve, reject) => {
    const child = spawn(cmd, args, { stdio: "inherit" });
    child.on("error", reject);
    child.on("exit", (code) => {
      if (code === 0) {
        resolve();
      } else {
        reject(new Error(`${cmd} exited with code ${code}`));
      }
    });
  });
}

async function main() {
  const { filename, ext, osName, arch } = getAssetInfo();
  await ensureDir(BIN_DIR);

  if (fs.existsSync(BIN_PATH)) {
    return;
  }

  const archivePath = path.join(BIN_DIR, filename);
  const release = await resolveRelease(filename, osName, arch, ext);
  console.log(
    `Downloading AgentX ${release.tag} for ${os.platform()} ${os.arch()}...`
  );
  await downloadTo(release.url, archivePath);
  await extractArchive(archivePath, ext);

  const found = await findInstalledBinary();
  if (!found) {
    throw new Error("Downloaded archive but could not find agentx binary.");
  }
  await promoteBinary(found);
  await makeExecutable(BIN_PATH);
}

async function resolveRelease(expectedFilename, osName, arch, ext) {
  const tag = `v${VERSION}`;
  const byTag = await fetchRelease(`https://api.github.com/repos/${REPO}/releases/tags/${tag}`);
  if (byTag) {
    const asset = pickAsset(byTag, expectedFilename);
    if (asset) return { url: asset.browser_download_url, tag };
  }

  const latest = await fetchRelease(
    `https://api.github.com/repos/${REPO}/releases/latest`
  );
  if (latest) {
    const fallbackName = `agentx_${latest.tag_name.replace(/^v/, "")}_${osName}_${arch}.${ext}`;
    const asset = pickAsset(latest, fallbackName);
    if (asset) return { url: asset.browser_download_url, tag: latest.tag_name };
  }
  throw new Error("Could not locate a matching release asset on GitHub.");
}

function pickAsset(release, filename) {
  if (!release || !Array.isArray(release.assets)) return null;
  return release.assets.find((asset) => asset.name === filename) || null;
}

async function fetchRelease(url) {
  try {
    return await fetchJson(url);
  } catch (err) {
    if (err.statusCode === 404) return null;
    throw err;
  }
}

async function fetchJson(url, redirects = 0) {
  return await new Promise((resolve, reject) => {
    const options = {
      headers: {
        "User-Agent": "agentx-npm-installer",
        Accept: "application/vnd.github+json",
      },
    };
    https
      .get(url, options, (res) => {
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          if (redirects >= 5) {
            const err = new Error("Too many redirects while fetching JSON.");
            err.statusCode = res.statusCode;
            reject(err);
            res.resume();
            return;
          }
          res.resume();
          const nextUrl = new URL(res.headers.location, url).toString();
          fetchJson(nextUrl, redirects + 1).then(resolve).catch(reject);
          return;
        }
        if (res.statusCode !== 200) {
          const err = new Error(`Failed to fetch ${url} (status ${res.statusCode})`);
          err.statusCode = res.statusCode;
          reject(err);
          res.resume();
          return;
        }
        let body = "";
        res.setEncoding("utf8");
        res.on("data", (chunk) => {
          body += chunk;
        });
        res.on("end", () => {
          try {
            resolve(JSON.parse(body));
          } catch (parseErr) {
            reject(parseErr);
          }
        });
      })
      .on("error", reject);
  });
}

main().catch((err) => {
  console.error(`AgentX install failed: ${err.message}`);
  process.exit(1);
});
