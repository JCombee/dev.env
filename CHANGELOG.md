# Changelog
## [0.8.0] - 2026-09-05

### New Features

- Add global apps with LiteLLM ([8c14f9d](https://github.com/JCombee/dev.env/commit/8c14f9d6961564e18005036fce52f9b434c453c9))

## [0.7.0] - 2026-05-21

### New Features

- Provision project database after container start ([a1752e6](https://github.com/JCombee/dev.env/commit/a1752e62a21bf7ae6784b0c868ffbf07f40308ba))
- Add dev db import command ([ce7fcb6](https://github.com/JCombee/dev.env/commit/ce7fcb65ea870bd1982df6aac366392239d3ac98))

## [0.6.5] - 2026-05-13

### Bug Fixes

- Check error return from fmt.Scanln in confirmPrompt ([5b1a498](https://github.com/JCombee/dev.env/commit/5b1a4983cb5afe15a5ada03f501b04cf202aca20))

## [0.6.4] - 2026-05-13

### Bug Fixes

- Use github native changelog in goreleaser ([aa241d9](https://github.com/JCombee/dev.env/commit/aa241d97c6d2d900f22145f11c3088bdc7ead4e5))

## [0.6.3] - 2026-05-13

### Bug Fixes

- Skip goreleaser changelog, use git-cliff output only ([8ef43c1](https://github.com/JCombee/dev.env/commit/8ef43c1c0e9967e895068fd13ffcab93be053ca6))

## [0.6.2] - 2026-05-13

### Bug Fixes

- Remove unsupported git-cliff sub-config from goreleaser changelog ([de09bfc](https://github.com/JCombee/dev.env/commit/de09bfc680ebdbd547f6bae040a90be7fb67cf0d))

## [0.6.1] - 2026-05-13

### Bug Fixes

- Use git-cliff changelog provider in goreleaser v2 ([ea906a9](https://github.com/JCombee/dev.env/commit/ea906a9bb30cf8ac4224c7145a4755d7477387e7))

## [0.6.0] - 2026-05-13

### Bug Fixes

- Run goreleaser before changelog commit to avoid HEAD/tag mismatch ([6da18f5](https://github.com/JCombee/dev.env/commit/6da18f54df4c7fb580a22da6bffe7a47d5c88259))

### New Features

- Add dev exec command for interactive container sessions ([12739cc](https://github.com/JCombee/dev.env/commit/12739cca73dd906caf66b4e4e2038a854c44dcd4))
- Pass extra flags through dev exec to inner command ([0fb5968](https://github.com/JCombee/dev.env/commit/0fb5968a15c2068a4dffe04100ed0b527dd6505d))

## [0.5.0] - 2026-05-12

### Bug Fixes

- Remove invalid abbrev_length field from goreleaser changelog config ([4d58f26](https://github.com/JCombee/dev.env/commit/4d58f268a5a13ed81e8301c4434d3e54c4e230c2))
- Rebase CHANGELOG commit onto main before push to avoid non-fast-forward ([3dd95db](https://github.com/JCombee/dev.env/commit/3dd95db536be11165da7f506ab67ab9bd0a56375))

### New Features

- Implement port allocation for containers ([bee1d43](https://github.com/JCombee/dev.env/commit/bee1d438eebc54976939c089fb81d8db6b13e73d))
- Implement .env file generation from running services ([c3df673](https://github.com/JCombee/dev.env/commit/c3df6739bc1878379207e59a7da40019c13f832b))
- Add .env wizard with diff and confirmation for existing files ([46e65b0](https://github.com/JCombee/dev.env/commit/46e65b0c74025726875f00fc50af16cce0da7a1e))

## [0.4.0] - 2026-04-27

### Bug Fixes

- Handle old list-format services.yaml during unmarshal ([f37c55b](https://github.com/JCombee/dev.env/commit/f37c55bc77cf7905c5b0e0955c2abd6a6c70794d))
- Use errors.As for ErrNameCollision type check (errorlint) ([3badf6e](https://github.com/JCombee/dev.env/commit/3badf6ed222d10872bc08c869d242944989c4b0e))
- Add pull-requests: read permission to release workflow ([a9a1886](https://github.com/JCombee/dev.env/commit/a9a18867a23b0165c3f902d2ed61f6891691891e))
- Prevent CHANGELOG commit from rewinding main to old tag ([ab903b5](https://github.com/JCombee/dev.env/commit/ab903b5464b0cfe70593dc968b1d843a1b360225))
- Pass GITHUB_TOKEN to git-cliff for private repo PR access ([9a67451](https://github.com/JCombee/dev.env/commit/9a67451f94e5d90a14296257cccfea1fbc4af6c1))

### New Features

- Implement docker compose start/stop with reference counting (phase 4) ([7d6d32c](https://github.com/JCombee/dev.env/commit/7d6d32ce100260997873129e00ec71505fcefaf1))
- Add workflow_dispatch to release workflow ([e07954b](https://github.com/JCombee/dev.env/commit/e07954b41983133aee9c3db09326df5178220e2c))

### Refactoring

- Simplify release workflow, drop workflow_dispatch ([22472b5](https://github.com/JCombee/dev.env/commit/22472b546b36487057a09ff5a4171d93d01e0d9d))

## [0.3.0] - 2026-04-26

### New Features

- Phase 3 — global state, project registration, secrets, dev status ([babf2cc](https://github.com/JCombee/dev.env/commit/babf2cc3d9e33782ccb837094784b5247e06f453))

## [0.2.0] - 2026-04-26

### New Features

- Generate CHANGELOG.md via git-cliff on release ([61b75e1](https://github.com/JCombee/dev.env/commit/61b75e151fcf358720f0c9e7d66714519cc98d05))
- Phase 2 — config model, service registry, dev init, dev config ([1f1846e](https://github.com/JCombee/dev.env/commit/1f1846e5fcc37177681f0b5d740cd76a4116dee3))

## [0.1.0] - 2026-04-26

### Bug Fixes

- Lower go directive to 1.24 for golangci-lint compatibility ([e839c5a](https://github.com/JCombee/dev.env/commit/e839c5abca5ee70a265cfaae158dfa3c700a243f))
- Remove v2 version field from golangci config ([f945958](https://github.com/JCombee/dev.env/commit/f945958f82a79c46af332f31018c4a6563c885fe))
- Write proper YAML keys in stub settings/services files ([ece2718](https://github.com/JCombee/dev.env/commit/ece27185ba37bbe1fbcafadc9f450b0ae9ba982c))

### New Features

- Phase 1 skeleton, CI pipeline, and e2e test harness ([3f02812](https://github.com/JCombee/dev.env/commit/3f0281267c0c45f288e060d660d3e6e8469246d0))
- Automated draft releases via GoReleaser ([fca3eb3](https://github.com/JCombee/dev.env/commit/fca3eb3226f72bf3c6f49d61c2e3a12660c8eace))


