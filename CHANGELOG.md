# Changelog
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


