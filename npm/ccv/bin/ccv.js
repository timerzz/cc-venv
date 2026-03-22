#!/usr/bin/env node

const { spawn } = require("node:child_process");
const path = require("node:path");

function resolvePackageName(platform, arch) {
  if (platform === "linux" && arch === "x64") return "@timerzz/ccv-linux-x64";
  if (platform === "linux" && arch === "arm64") return "@timerzz/ccv-linux-arm64";
  if (platform === "darwin" && arch === "x64") return "@timerzz/ccv-darwin-x64";
  if (platform === "darwin" && arch === "arm64") return "@timerzz/ccv-darwin-arm64";
  if (platform === "win32" && arch === "x64") return "@timerzz/ccv-win32-x64";
  if (platform === "win32" && arch === "arm64") return "@timerzz/ccv-win32-arm64";

  return null;
}

function resolveBinaryPath() {
  const pkgName = resolvePackageName(process.platform, process.arch);
  if (!pkgName) {
    console.error(`Unsupported platform: ${process.platform}/${process.arch}`);
    process.exit(1);
  }

  let packageJsonPath;
  try {
    packageJsonPath = require.resolve(`${pkgName}/package.json`);
  } catch (error) {
    console.error(`Missing platform package ${pkgName}. Try reinstalling ccv.`);
    process.exit(1);
  }

  const binName = process.platform === "win32" ? "ccv.exe" : "ccv";
  return path.join(path.dirname(packageJsonPath), "bin", binName);
}

const binaryPath = resolveBinaryPath();
const child = spawn(binaryPath, process.argv.slice(2), { stdio: "inherit" });

child.on("exit", (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 1);
});

child.on("error", (error) => {
  console.error(error.message);
  process.exit(1);
});
