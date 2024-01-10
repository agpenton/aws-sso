gt_token=${GITHUB_TOKEN}


#@PHONY: help
@PHONY: run

run:
	gh act -s github.token=gt_token -W ./.github/workflows/goreleaser.yaml
