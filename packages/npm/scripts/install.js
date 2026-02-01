#!/usr/bin/env node
const https = require("https");
const fs = require("fs");
const path = require("path");
const os = require("os");

const REPO = "Europroiect-Estate/Codeez-AI";
const VERSION = process.env.CODEEZ_VERSION || "latest";

function get(url) {
  return new Promise((resolve, reject) => {
    https.get(url, { headers: { "User-Agent": "codeez-npm" } }, (res) => {
      let data = "";
      res.on("data", (ch) => (data += ch));
      res.on("end", () => resolve(data));
      res.on("error", reject);
    }).on("error", reject);
  });
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https.get(url, { headers: { "User-Agent": "codeez-npm" } }, (res) => {
      res.pipe(file);
      file.on("finish", () => {
        file.close();
        fs.chmodSync(dest, 0o755);
        resolve();
      });
    }).on("error", (err) => {
      fs.unlink(dest, () => {});
      reject(err);
    });
  });
}

async function main() {
  const platform = os.platform();
  const arch = os.arch();
  const goos = platform === "darwin" ? "Darwin" : platform === "win32" ? "Windows" : "Linux";
  const goarch = arch === "x64" ? "amd64" : arch === "arm64" ? "arm64" : "amd64";
  let version = VERSION;
  if (version === "latest") {
    const data = JSON.parse(await get(`https://api.github.com/repos/${REPO}/releases/latest`));
    version = data.tag_name;
  }
  const name = `codeez_${version.replace(/^v/, "")}_${goos}_${goarch}`;
  const ext = goos === "Windows" ? "zip" : "tar.gz";
  const archive = `${name}.${ext}`;
  const url = `https://github.com/${REPO}/releases/download/${version}/${archive}`;
  const binDir = path.join(__dirname, "..", "bin");
  const binName = goos === "Windows" ? "codeez.exe" : "codeez";
  const binPath = path.join(binDir, binName);
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }
  try {
    await download(url, binPath);
    console.log("codeez binary installed to", binPath);
  } catch (e) {
    console.error("Download failed. Build from source: go build -o codeez ./cmd/codeez");
    console.error(e.message);
  }
}

main();
