#!/usr/bin/env python3

# Copyright 2023 Flant JSC
# Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE

from shell_operator import hook
from yaml import dump


# Converts FalcoAuditRules CRD format to the native Falco rules
def convert_spec(spec: dict) -> list:
    result = []

    required_engine_version = spec.get("requiredEngineVersion")
    if required_engine_version is not None:
        result.append({
            "required_engine_version": required_engine_version,
        })

    required_k8saudit_plugin_version = spec.get("requiredK8sAuditPluginVersion")
    if required_k8saudit_plugin_version is not None:
        result.append({
            "required_plugin_versions": [
                {
                    "name": "k8saudit",
                    "version": required_k8saudit_plugin_version,
                },
            ],
        })

    for item in spec["rules"]:
        if item.get("rule") is not None:
            converted_item = {**item["rule"]}
            converted_item["rule"] = converted_item.pop("name")

            source = item["rule"].get("source")
            if source is not None:
                converted_item["source"] = source.lower()

            result.append(converted_item)
            continue
        if item.get("macro") is not None:
            result.append({
                "macro": item["macro"]["name"],
                "condition": item["macro"]["condition"],
            })
            continue
        if item.get("list") is not None:
            result.append({
                "list": item["list"]["name"],
                "items": item["list"]["items"],
            })
            continue

    return result


def main(ctx: hook.Context):
    for s in ctx.snapshots["rules"]:
        filtered = s["filterResult"]
        if filtered is None:
            # Should not happen
            continue

        filename = f'{filtered["name"]}.yaml'

        with open(f'/etc/falco/rules.d/{filename}', "w") as file:
            spec = convert_spec(filtered["spec"])
            file.write(dump(spec))

    with open('/tmp/ready', "w") as file:
        file.write("ok")


if __name__ == "__main__":
    hook.run(main, configpath="rules/config.yaml")
