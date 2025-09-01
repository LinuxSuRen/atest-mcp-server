#!/usr/bin/env node
import { homedir } from 'os';
import { join, dirname } from 'path';
import { existsSync, mkdirSync, createWriteStream } from 'fs';
import { fileURLToPath } from 'url';
import fetch from 'node-fetch';
import tar from 'tar';
import { spawn } from 'child_process';

const OWNER = 'LinuxSuRen';
const REPO  = 'atest-mcp-server';
const CACHE_DIR = join(homedir(), '.config', 'bin');
const BIN_PATH  = join(CACHE_DIR, 'atest-store-mcp');

(async () => {
  // 1. 如果本地已有缓存的可执行文件，直接启动
  if (existsSync(BIN_PATH)) {
    spawn(BIN_PATH, process.argv.slice(2), { stdio: 'inherit' });
    return;
  }

  mkdirSync(CACHE_DIR, { recursive: true });

  const latestUrl = `https://api.github.com/repos/${OWNER}/${REPO}/releases/latest`;
  const release = await (await fetch(latestUrl)).json();
  const assetName = `atest-mcp-server-${process.platform}-${process.arch}.tar.gz`;
  const asset = release.assets.find(a => a.name === assetName);
  if (!asset) {
    console.error(`❌ 未找到适配的二进制：${assetName}`);
    process.exit(1);
  }

  console.log(`📥 下载 ${asset.browser_download_url}`);
  const tarPath = join(CACHE_DIR, asset.name);
  const res = await fetch(asset.browser_download_url);
  res.body.pipe(createWriteStream(tarPath))
          .on('finish', () => {
            tar.x({ file: tarPath, cwd: CACHE_DIR, strip: 1 })
               .then(() => {
                 // 5. 启动
                 spawn(BIN_PATH, process.argv.slice(2), { stdio: 'inherit' });
               });
          });
})();