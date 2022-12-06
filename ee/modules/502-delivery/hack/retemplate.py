#!/usr/bin/env python3

import os


def map_image(image):
    """
    Match image name by substring and return Helm template function call.
    """
    if image.find("updater") > -1:
        return '{{ include "helm_lib_module_image" (list . "argocdImageUpdater") }}'
    if image.find("argo") > -1:
        return '{{ include "helm_lib_module_image" (list . "argocd") }}'
    if image.find("redis") > -1:
        return '{{ include "helm_lib_module_image" (list . "redis") }}'
    return "UNKNOWN"


def include_container_security_context(nindent):
    """
    Get common security context for containers and Pods
    """
    return " ".join(
        (
            " " * (nindent - 1),  # Helm template indentation
            "{{-",
            'include "helm_lib_module_container_security_context_read_only_root_filesystem_capabilities_drop_all"',
            f"| nindent {nindent}",
            "}}\n",
        )
    )


def include_module_labels(labels: dict, nindent: int):
    """
    Get module labels Helm template function call.
    """
    kv_str = " ".join([f'"{k}" "{v}"' for k, v in labels.items()])
    return " ".join(
        (
            " " * (nindent - 1),  # Helm template indentation
            "{{-",
            f'include "helm_lib_module_labels" (list . (dict {kv_str}))',
            f"| nindent {nindent}",
            "}}\n",
        )
    )


def re_template(filepath: str):
    """
    Walk through ArgoCD templates and substitute Helm template functions.
    """

    print(f"File {filepath}")
    lines = list(open(filepath).readlines())

    overrides, erasures, insertions = [], [], []

    in_labels = False
    labels = {}

    in_seccontext = False

    for i, l in enumerate(lines):
        if l.strip() == "labels:":
            l_indent = l.count(" ", 0, l.find("l"))
            in_labels = True
            i_labels_start = i
            print(f"Found labels at {i}")
            print(f"Labels indent is {l_indent}")
            continue
        if in_labels:
            if l.startswith((2 + l_indent) * " "):
                k, v = l.strip().split(":")
                print(f"Found label {k}={v}")
                labels[k.strip()] = v.strip()
            else:
                # After labels
                if l_indent > 2:
                    # Pods should have only additional 'app' label
                    indent = (l_indent + 2) * " "
                    app_label = (
                        f'{indent}app: {labels.get("app.kubernetes.io/name", "")}\n'
                    )
                    insertions.append((i, app_label))  # (index, text)
                else:
                    # Inject custom labels
                    app_name = labels.get("app.kubernetes.io/name", "")
                    if app_name != "":
                        labels["app"] = app_name

                    print(f"Labels end at {i}")
                    overrides.append(
                        (i_labels_start, include_module_labels(labels, l_indent))
                    )
                    erasures += range(i_labels_start + 1, i)
                labels = {}
                in_labels = False
            continue

        # Substitute securityContext with *drop_all
        if l.strip() == "securityContext:":
            print(f"Found securityContext at {i}")
            in_seccontext = True
            sc_indent = l.count(" ", 0, l.find("s"))
            # overwrite currentstarting line
            overrides.append((i, include_container_security_context(sc_indent)))
            continue
        if in_seccontext:
            # drop all inner lines
            if l.startswith((2 + sc_indent) * " "):
                erasures.append(i)
            else:
                in_seccontext = False
            continue

        # Subtitute images
        if l.strip().startswith("image: "):
            image_parts = lines[i].split(": ")
            print(f"Found image '{image_parts[1]}'")
            lines[i] = image_parts[0] + ": " + map_image(image_parts[1].strip()) + "\n"

    # Apply changes
    for i, l in overrides:
        lines[i] = l

    for i in erasures:
        lines[i] = ""

    shift = 0
    for i, l in insertions:
        lines.insert(i + shift, l)
        shift += 1

    with open(filepath, "w") as f:
        f.writelines([l for l in lines if l.strip() != ""])


if __name__ == "__main__":
    templates_root = os.path.join(os.getcwd(), "templates", "argocd")
    for root, dirs, files in os.walk(templates_root):
        print(root)
        # print("\n".join(["  ./" + d for d in dirs]))

        argo_files = [f for f in files if f.endswith("yaml")]
        print("\n".join(["    " + f for f in argo_files]))

        for f in argo_files:
            filepath = os.path.join(root, f)
            re_template(filepath)
