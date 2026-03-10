#!/usr/bin/env node
'use strict'

const { execFileSync, execSync } = require('child_process')
const { existsSync, mkdtempSync, createWriteStream, chmodSync } = require('fs')
const { join } = require('path')
const { tmpdir, platform, arch } = require('os')
const https = require('https')
const { createGunzip } = require('zlib')
const tar = require('tar')
const path = require('path')

const REPO_OWNER = 'envm-org'
const REPO_NAME = 'envm'
const VERSION = require('../package.json').version

function getPlatformInfo() {
  const p = platform()
  const a = arch()

  const osMap = { linux: 'Linux', darwin: 'Darwin', win32: 'Windows' }
  const archMap = { x64: 'x86_64', arm64: 'arm64', ia32: 'i386' }

  const osName = osMap[p]
  const archName = archMap[a]

  if (!osName) throw new Error(`Unsupported platform: ${p}`)
  if (!archName) throw new Error(`Unsupported architecture: ${a}`)

  return { osName, archName }
}

async function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    const file = createWriteStream(dest)
    https
      .get(url, (res) => {
        if (res.statusCode === 302 || res.statusCode === 301) {
          return downloadFile(res.headers.location, dest).then(resolve).catch(reject)
        }
        if (res.statusCode !== 200) {
          return reject(new Error(`HTTP ${res.statusCode} for ${url}`))
        }
        res.pipe(file)
        file.on('finish', () => file.close(resolve))
      })
      .on('error', reject)
  })
}

async function installBinary() {
  const { osName, archName } = getPlatformInfo()
  const ext = osName === 'Windows' ? 'zip' : 'tar.gz'
  const archiveName = `${REPO_NAME}_${osName}_${archName}.${ext}`
  const tag = `v${VERSION}`
  const url = `https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${tag}/${archiveName}`

  const tmpDir = mkdtempSync(join(tmpdir(), 'envm-'))
  const archivePath = join(tmpDir, archiveName)

  process.stderr.write(`Downloading envm ${tag} for ${osName}/${archName}...\n`)

  await downloadFile(url, archivePath)

  process.stderr.write('Extracting...\n')
  await tar.x({ file: archivePath, cwd: tmpDir })

  const binaryName = osName === 'Windows' ? 'envm.exe' : 'envm'
  const binaryPath = join(tmpDir, binaryName)

  if (!existsSync(binaryPath)) {
    throw new Error(`Binary not found in archive at ${binaryPath}`)
  }

  chmodSync(binaryPath, 0o755)

  // Re-exec with the downloaded binary passing all CLI args
  const args = process.argv.slice(2)
  try {
    execFileSync(binaryPath, args, { stdio: 'inherit' })
  } catch (err) {
    process.exit(err.status ?? 1)
  }
}

installBinary().catch((err) => {
  process.stderr.write(`Error: ${err.message}\n`)
  process.exit(1)
})
