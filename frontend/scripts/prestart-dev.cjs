const { spawnSync } = require('node:child_process');

function runCommand(command, args, label, options = {}) {
  const result = spawnSync(command, args, {
    stdio: 'inherit',
    shell: process.platform === 'win32',
    ...options,
  });

  if (result.error) {
    throw new Error(`[${label}] ${result.error.message}`);
  }

  return result.status ?? 1;
}

const npmCommand = process.platform === 'win32' ? 'npm.cmd' : 'npm';

try {
  const swaggerStatus = runCommand(npmCommand, ['run', 'swagger'], 'swagger');

  if (swaggerStatus !== 0) {
    console.warn('[prestart:dev] Swagger generation failed, continuing with existing docs.');
  }
} catch (error) {
  console.warn(`${error.message}\n[prestart:dev] Swagger generation failed, continuing with existing docs.`);
}

const openapiStatus = runCommand(npmCommand, ['run', 'openapi'], 'openapi');
process.exit(openapiStatus);
