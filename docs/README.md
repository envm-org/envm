# ENVM Documentation

This is the official documentation site for [ENVM](https://github.com/envm-org/envm) - a secure environment variable management and sync tool.

Built with [Docusaurus](https://docusaurus.io/).

## Development

```bash
# Install dependencies
bun install

# Start development server
bun run start

# Build for production
bun run build
```

## Deployment

The documentation is automatically deployed to GitHub Pages when changes are pushed to the `main` branch.

**Live site:** https://envm-org.github.io/envm/

### Manual Deployment

```bash
bun run build
```

The `build` folder contains the static files ready for deployment.

## GitHub Pages Setup

1. Go to your repository Settings → Pages
2. Under "Build and deployment", select **GitHub Actions** as the source
3. Push to `main` branch to trigger deployment

## License

MIT License - see [LICENSE](LICENSE) for details.