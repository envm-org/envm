# Testing Environment File Tracking and Preservation

This directory `examples/` has been set up with the following environment files to help you test the new tracking capabilities of the `envm` CLI:

*   `.env`
*   `.env.local`
*   `.env.prod`
*   `.env.example`

## Test 1: Tracking Everything Except .example

You can test that the `init` or `load` commands track `.env`, `.env.local`, and `.env.prod`, but safely ignore `.env.example`.

1. Ensure your backend server is running (`docker compose up -d` in the root).
2. Authenticate the CLI with `envm login`.
3. In this `examples` directory, run:
   ```bash
   envm init
   ```
4. **Expected Result:** The CLI should output: `Found 3 .env files` and will prompt you to create the project. It should explicitly *not* track `.env.example`.

## Test 2: Preserving Extensions on Pull

Once the files are pushed up, you can test if `envm pull` properly reconstructs them.

1. First, push the files up after initializing the project:
   ```bash
   envm push
   ```
2. Delete the local environment files from the directory:
   ```bash
   rm .env .env.local .env.prod
   ```
3. Run the pull command:
   ```bash
   envm pull
   ```
## Test 3: Recursive Scanning (Monorepo/Workspace)

The `envm` CLI is designed to recursively scan subdirectories for environment files (e.g., in a monorepo).
I have created a `workspace/` directory inside `examples/` with three mock applications:
*   `workspace/app1/.env` and `.env.staging`
*   `workspace/app2/.env.local` and `.env.example`
*   `workspace/app3/.env`

Run `envm init` from this `examples/` directory or directly inside `examples/workspace/`. 
**Expected Result:** `envm` will deep-scan and find all valid `.env` variants across all nested applications, while successfully ignoring `workspace/app2/.env.example`.
