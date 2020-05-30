from __future__ import annotations
from egoist.app import create_app, SettingsDict, parse_args
from egoist.go.types import GoError

settings: SettingsDict = {"rootdir": "../cmd", "here": __file__}
app = create_app(settings)


app.include("egoist.directives.define_cli")


@app.define_cli("egoist.generators.clikit:walk")
def selfish(*, alias: str, delete: bool, silent: bool) -> GoError:
    """individual gist client """
    from egoist.generators.clikit import runtime, clikit
    from egoist.go import di
    from egoist.go.types import GoError
    import components  # ./components.py

    options = runtime.get_cli_options()
    options.alias.help = "alias name of uploaded gists"
    options.delete.help = "delete uploaded gists"
    options.silent.help = "don't open gist pages, after uploaded"

    with runtime.generate(clikit) as m:
        b = di.Builder()
        b.add_provider(components.NewCommitHistory)
        b.add_provider(components.LoadConfig)
        b.add_provider(components.NewClient)
        b.add_provider(components.NewApp)

        injector = b.build(variables={**locals(), "name": m.symbol('"selfish"')})
        app = injector.inject(m)

        context_pkg = m.import_("context")
        args = runtime.get_cli_rest_args()
        m.return_(app.Run(context_pkg.Background(), args))


if __name__ == "__main__":
    for argv in parse_args(sep="-"):
        app.run(argv)
