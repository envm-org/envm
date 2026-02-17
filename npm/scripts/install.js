const axios = require('axios');
const fs = require('fs');
const path = require('path');
const tar = require('tar');
const unzipper = require('unzipper');
const packageJson = require('../package.json');

const BIN_DIR = path.join(__dirname, '../bin');
const BIN_NAME = process.platform === 'win32' ? 'envm.exe' : 'envm';
const BIN_PATH = path.join(BIN_DIR, BIN_NAME);

if (!fs.existsSync(BIN_DIR)) {
  fs.mkdirSync(BIN_DIR);
}

function getArtifactName() {
  const platform = process.platform;
  const arch = process.arch;

  let os = '';
  let architecture = '';
  let ext = 'tar.gz';

  switch (platform) {
    case 'linux':
      os = 'Linux';
      break;
    case 'darwin':
      os = 'Darwin';
      break;
    case 'win32':
      os = 'Windows';
      ext = 'zip';
      break;
    default:
      console.error(`Unsupported OS: ${platform}`);
      process.exit(1);
  }

  switch (arch) {
    case 'x64':
      architecture = 'x86_64';
      break;
    case 'arm64':
      architecture = 'arm64';
      break;
    case 'ia32':
      architecture = 'i386';
      break;
    default:
      console.error(`Unsupported architecture: ${arch}`);
      process.exit(1);
  }

  const version = packageJson.version;
  
  return `envm_${version}_${os}_${architecture}.${ext}`;
}

async function downloadAndExtract() {
  const artifactName = getArtifactName();
  const version = packageJson.version;
  const url = `https://github.com/envm-org/envm/releases/download/v${version}/${artifactName}`;

  console.log(`Downloading ${artifactName} from ${url}...`);

  try {
    const response = await axios({
      url,
      method: 'GET',
      responseType: 'stream',
    });

    if (artifactName.endsWith('.zip')) {
      response.data
        .pipe(unzipper.Parse())
        .on('entry', function (entry) {
          if (entry.path === BIN_NAME) {
            entry.pipe(fs.createWriteStream(BIN_PATH, { mode: 0o755 }));
          } else {
            entry.autodrain();
          }
        })
        .on('close', () => {
          console.log('Successfully installed envm.');
        });
    } else {
      const extract = tar.x({
        cwd: BIN_DIR,
        filter: (path) => path === BIN_NAME,
      });
      
      response.data.pipe(extract);
      
      extract.on('close', () => {
         if(fs.existsSync(BIN_PATH)) {
            fs.chmodSync(BIN_PATH, 0o755);
            console.log('Successfully installed envm.');
         } else {
             console.error("Binary not found after extraction.");
             process.exit(1);
         }
      });
    }
  } catch (error) {
    console.error(`Failed to download or install envm: ${error.message}`);
    if (error.response && error.response.status === 404) {
      console.error("Release not found. Please ensure the version in package.json matches a GitHub release.");
    }
    process.exit(1);
  }
}

downloadAndExtract();
