from __future__ import annotations
import typing as t
from egoist.go.types import GoError, gopackage, rename


class CommitHistory:
    pass


@gopackage("github.com/podhmo/selfish/internal/commithistory")
@rename("New")
def NewCommitHistory(name: str) -> CommitHistory:
    pass


class Config:
    pass


@gopackage("github.com/podhmo/selfish")
def LoadConfig(commit_history: CommitHistory) -> t.Tuple[Config, GoError]:
    pass


class Client:
    pass


@gopackage("github.com/podhmo/selfish")
def NewClient(config: Config) -> Client:
    pass


class App:
    pass


@gopackage("github.com/podhmo/selfish/cmd/selfish/internal")
def NewApp(
    c: CommitHistory,
    client: Client,
    config: Config,
    silent: bool,
    delete: bool,
    alias: str,
) -> App:
    pass
