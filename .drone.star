def main(ctx):
    if ctx.build.event == "push":
        return [
            {
                "kind": "pipeline",
                "type": "docker",
                "name": "cov-tests",
                "steps": [
                    {
                        "name": "tests",
                        "image": "golang:latest",
                        "commands": [
                            "go test -coverprofile=coverage.out && go tool cover -func=coverage.out"
                        ],
                    },
                ],
                "node": {"docker": "slow"},
            },
        ]
