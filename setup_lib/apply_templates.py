"""Apply project-specific source/Dockerfile/handler templates."""
from __future__ import print_function

import os
import shutil
import subprocess
import sys

from setup_lib import constants

def apply_go_mod_module(app_name):
    if not app_name or not os.path.isfile(constants.GO_MOD_PATH):
        return
    with open(constants.GO_MOD_PATH, "r", encoding="utf-8") as f:
        lines = f.readlines()
    for i, line in enumerate(lines):
        if line.startswith("module "):
            lines[i] = "module {}\n".format(app_name.strip())
            break
    with open(constants.GO_MOD_PATH, "w", encoding="utf-8") as f:
        f.writelines(lines)
    print("Updated go.mod module to {}".format(app_name))


def apply_handler_template(app_type):
    tmpl_name = constants.TEMPLATE_FILE_BY_TYPE.get(app_type)
    if not tmpl_name:
        return
    src = os.path.join(constants.LAMBDA_CMD_DIR, tmpl_name)
    if not os.path.isfile(src):
        print("Warning: {} not found; skipping handler setup".format(tmpl_name))
        return
    shutil.copy2(src, constants.MAIN_GO_PATH)
    print("Created cmd/lambda/main.go from {}".format(tmpl_name))
    _sync_go_module()


def _sync_go_module():
    """Resolve imports in go.mod/go.sum after handler template changes."""
    try:
        subprocess.run(
            ["go", "mod", "tidy"],
            cwd=constants.SCRIPT_DIR,
            check=True,
            capture_output=True,
            text=True,
        )
        print("Updated go.mod and go.sum")
    except FileNotFoundError:
        print("Warning: go not found; run 'go mod tidy' manually.", file=sys.stderr)
    except subprocess.CalledProcessError as e:
        print("Warning: go mod tidy failed: {}".format(e.stderr or e), file=sys.stderr)
