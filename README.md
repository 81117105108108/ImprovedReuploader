<div align="center">
        <img src="https://github.com/user-attachments/assets/4e753b87-3069-4a8f-824b-69fc910584c5" alt="Icon" />
    <br><br>
    <a href="https://github.com/kartFr/Asset-Reuploader/releases/latest"><img src="https://img.shields.io/github/downloads/kartfr/Asset-Reuploader/total?color=yellow" alt="Latest download" /></a>
    <a href="https://github.com/kartFr/Asset-Reuploader/releases/latest"><img src="https://img.shields.io/github/v/release/kartfr/Asset-Reuploader?color=yellow" alt="Latest release" /></a>
    <a href="https://discord.gg/XTEtUqPTat"><img src="https://img.shields.io/discord/1238572493925646347?label=discord&logo=discord&logoColor=white&color=yellow" alt="Discord" /></a>
    <a href="https://github.com/kartFr/Asset-Reuploader?tab=GPL-3.0-1-ov-file"><img src="https://img.shields.io/github/license/kartFr/Asset-Reuploader?color=yellow" alt="License" /></a>
</div>

#

Roblox asset reuploader using a [roblox plugin](https://create.roblox.com/store/asset/89096096219225/Asset-Reuploader) and go.

Intended to reupload assets since you are unable to use assets if another account owns them.

### Community

Join the [discord](https://discord.gg/XTEtUqPTat) and interact with the community. We have:

- Reuploader themes
- Reuploader support
- Coding support
- Early access to updates when supporting me
- Advertise your game/discord server

### Currently supported assets

- Animations
- Audio (paid, join discord)
- Meshes (paid)

## This fork: reliability + concurrency hardening

Upstream fork maintained at `81117105108108/ImprovedReuploader` (`main`). Rounds 1–6 applied, all verified (`go build`/`go vet`/`go test`, `selene` clean, graph 719n/2436e):

- **Race fixes:** HTTP router state mutex-guarded (was data race + `log.Fatal` in handlers), `pauseController.IsPausedNow()`, `logger` mutex + `HistoryString()/Historyf()`, `shardedmap.GetOrCreateShard` (was TOCTOU + mutex leak), `AtomicArray` copy-on-write, cookie singleflight (was thundering herd), permission-limiter mixup fixed.
- **Dedup:** shared `NewUploadHandler` (Animation/Mesh), `newGamesHandler`, `assetutils.MoveValueToTop`, `pipeline` helpers (`RatePolicy`, `GetDefaultPlaceIDs`, `GroupByCreator`, `CalcBatchSize`, `QueuedDo`), sound/mesh use `QueuedDo` (animation keeps drain-in-slot intentionally).
- **Correctness:** `retry` exponential `Factor` + `Tries(-1)` infinite, `fixedWindow` race fix, poll 5m deadline + cancellable `sleepCtx`/`*Ctx` variants, `parseCookieFile` first-wins, dynamic `config.Get`, mesh/sound nil/bounds guards, missing `return` fixed.
- **Config:** `starts_per_minute`/`max_concurrent`/`max_parallel_chunks`/`upload_tries` + per-type `animation_/sound_/mesh_` overrides (empty = fallback). Chunk caps read config live.
- **Tests:** `retry`, `taskqueue`, `atomicarray`, `shardedmap`, `pipeline`, `config` (`go test ./internal/...`).
- **Luau:** `_G.getPlaceList` replaced with `App/GetPlaceList.luau` (set/wait-get, no race), `Theme.set` fallback to Studio, `Retry` parser-compat, `selene` 0/0 (see `selene.toml` notes on roblox string-`require`/outdated std).

## Dev

```powershell
go build ./...
go vet ./...
go test ./internal/retry/... ./internal/taskqueue/... ./internal/atomicarray/... ./internal/shardedmap/... ./internal/app/assets/pipeline/... ./internal/app/config/...
$env:PATH = "$env:USERPROFILE/.aftman/bin;$env:PATH"
aftman trust Kampfkarren/selene; aftman install
selene plugin/src
```

Rate tuning (`config.ini`, all optional — defaults shown):
`starts_per_minute=420`, `max_concurrent=24`, `max_parallel_chunks=6`, `upload_tries=3`, plus `animation_/sound_/mesh_starts_per_minute|max_concurrent|upload_tries` overrides.

Known indexer gaps (external): 12 Luau `parse_partial` from `typeof(LocalFunc)` forward-refs — valid new Luau, old tree-sitter grammar; grep fallback in flagged ranges.

## Contributing

ALL contributions are welcome. Feel free to make a pull/fork at any time. Please read the [contribution guide](https://github.com/kartFr/Asset-Reuploader/blob/main/CONTRIBUTING.md).

The [discord](https://discord.gg/XTEtUqPTat) does have a dedicated suggestions channel. If you do have a feature request please do put it in there. But, you do have free will, so choose whichever.

If you do need help with anything from stupid questions(pls no) to setting up a fork feel free to DM me on discord, my tag is alekfart.

## License

Copyright (c) 2024, kartFr

Licensed under the [GPL-3.0](https://github.com/kartFr/Asset-Reuploader/blob/main/LICENSE.txt) license.
