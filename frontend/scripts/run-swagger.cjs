const { spawnSync } = require('node:child_process');
const path = require('node:path');

const backendDir = path.resolve(__dirname, '../../backend');
const args = [
  'run',
  'github.com/swaggo/swag/cmd/swag@v1.16.6',
  'init',
  '-g',
  'cmd/server/main.go',
  '-o',
  'docs',
];

const result = spawnSync('go', args, {
  cwd: backendDir,
  stdio: 'inherit',
  shell: false,
});

if (result.error) {
  console.error('[swagger] Failed to start Go:', result.error.message);
  process.exit(1);
}

process.exit(result.status ?? 1);
