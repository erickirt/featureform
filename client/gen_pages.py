#  This Source Code Form is subject to the terms of the Mozilla Public
#  License, v. 2.0. If a copy of the MPL was not distributed with this
#  file, You can obtain one at http://mozilla.org/MPL/2.0/.
#
#  Copyright 2024 FeatureForm Inc.
#

import mkdocs_gen_files
import os

excluded_files = [
    "__pycache__",
    "proto",
    "__init__.py",
    "__main__.py",
    "cli.py",
    "client.py",
    "dashboard_metadata.py",
    "format.py",
    "get.py",
    "list.py",
    "local.py",
    "resources.py",
    "serving_test.py",
    "serving.py",
    "sqlite_metadata.py",
    "type_objects.py",
    "register.py",
    "dashboard",
    "get_local.py",
    "get_test.py",
    "list_local.py",
    "tls.py",
    "constants.py",
    "exceptions.py",
    "local_utils.py",
    "names_generator.py",
    "parse.py",
    "search.py",
    "search_local.py",
    "status_display.py",
    "local_cache.py",
    "version.py",
    "file_utils.py",
    "metadata.py",
    "providers",
]

for filename in os.listdir("./src/featureform"):
    if filename in excluded_files:
        continue
    if "test" in filename:
        continue

    file = filename.split(".")[0]
    mdFile = f"{file}.md"

    if mdFile == "metadata.md":
        continue
    with mkdocs_gen_files.open(mdFile, "w") as f:
        print(f"::: featureform.{file}", file=f)

    mkdocs_gen_files.set_edit_path(mdFile, "gen_pages.py")
