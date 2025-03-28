#  This Source Code Form is subject to the terms of the Mozilla Public
#  License, v. 2.0. If a copy of the MPL was not distributed with this
#  file, You can obtain one at http://mozilla.org/MPL/2.0/.
#
#  Copyright 2024 FeatureForm Inc.
#

import requests
from .format import *


def search(phrase, host):
    response = requests.get(f"http://{host}/data/search?q={phrase}")

    if response.status_code != 200:
        print(
            f"Search request for {phrase} resulted in HTTP status {response.status_code}"
        )
        return

    results = response.json()

    if len(results) == 0:
        print(f"Search phrase {phrase} returned no results.")
    else:
        format_rows("NAME", "VARIANT", "TYPE")
        for r in results:
            format_rows(r["Name"], r["Variant"], r["Type"])

    return results
