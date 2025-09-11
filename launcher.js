#!/usr/bin/env node
import { homedir } from 'os';
import { join } from 'path';
import { existsSync, mkdirSync, createWriteStream } from 'fs';
import fetch from 'node-fetch';
import tar from 'tar';
import { createReadStream } from 'fs';
import { pipeline } from 'stream/promises';
import { Extract } from 'unzipper';
import { spawn } from 'child_process';
import { execSync } from 'child_process';

const OWNER = 'LinuxSuRen';
const REPO  = 'atest-mcp-server';
const CACHE_DIR = join(homedir(), '.config', 'atest', 'bin');
const BIN_PATH  = join(CACHE_DIR, 'atest-store-mcp');

// Function to find executable in system PATH
function findExecutableInPath(execName) {
  try {
    const command = process.platform === 'win32' ? 'where' : 'which';
    const result = execSync(`${command} ${execName}`, { encoding: 'utf8', stdio: 'pipe' });
    return result.trim().split('\n')[0]; // Return first match
  } catch (error) {
    return null; // Not found in PATH
  }
}

(async () => {
  // 1. Check if cached executable exists locally
  if (existsSync(BIN_PATH)) {
    spawn(BIN_PATH, process.argv.slice(2), { stdio: 'inherit' });
    return;
  }

  // 2. Check if executable exists in system PATH
  const executableName = process.platform === 'win32' ? 'atest-store-mcp.exe' : 'atest-store-mcp';
  const systemExecutable = findExecutableInPath(executableName);
  if (systemExecutable) {
    spawn(systemExecutable, process.argv.slice(2), { stdio: 'inherit' });
    return;
  }

  mkdirSync(CACHE_DIR, { recursive: true });

  const latestUrl = `https://api.github.com/repos/${OWNER}/${REPO}/releases/latest`;
  const release = await (await fetch(latestUrl)).json();
  
  let platform
  let suffix
  switch (process.platform) {
    case 'win32':
      platform = 'windows';
      suffix = 'zip'
      break;
    case 'darwin':
    case 'linux':
      platform = process.platform;
      suffix = 'tar.gz'
      break;
    default:
      console.error(`❌ Unsupported platform: ${process.platform}`);
      process.exit(1);
  }

  let arch
  switch (process.arch) {
    case 'x64':
      arch = 'amd64';
      break;
    case 'arm64':
      arch = 'arm64';
      break;
    default:
      console.error(`❌ Unsupported architecture: ${process.arch}`);
      process.exit(1);
  }

  const assetName = `atest-store-mcp-${platform}-${arch}.${suffix}`;
  const asset = release.assets.find(a => a.name === assetName);
  if (!asset) {
    console.error(`❌ Compatible binary not found: ${assetName}`);
    process.exit(1);
  }

  console.log(`📥 Downloading ${asset.browser_download_url}`);
  const archivePath = join(CACHE_DIR, asset.name);
  const res = await fetch(asset.browser_download_url);
  res.body.pipe(createWriteStream(archivePath))
          .on('finish', async () => {
            try {
              if (suffix === 'zip') {
                // Extract zip for Windows
                await pipeline(
                  createReadStream(archivePath),
                  Extract({ path: CACHE_DIR })
                );
              } else {
                // Extract tar.gz for Linux/macOS
                await tar.x({ file: archivePath, cwd: CACHE_DIR, strip: 1 });
              }
              // 3. Launch
              spawn(BIN_PATH, process.argv.slice(2), { stdio: 'inherit' });
            } catch (error) {
              console.error(`❌ Extraction failed: ${error.message}`);
              process.exit(1);
            }
          });
})();