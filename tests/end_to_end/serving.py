#  This Source Code Form is subject to the terms of the Mozilla Public
#  License, v. 2.0. If a copy of the MPL was not distributed with this
#  file, You can obtain one at http://mozilla.org/MPL/2.0/.
#
#  Copyright 2024 FeatureForm Inc.
#

import os
import time

from dotenv import load_dotenv

import featureform as ff


SLEEP_DURATION = 30
NUMBER_OF_SLEEPS = 20

FILE_DIRECTORY = os.getenv("FEATUREFORM_TEST_PATH", "")
featureform_location = os.path.dirname(os.path.dirname(FILE_DIRECTORY))
env_file_path = os.path.join(featureform_location, ".env")
load_dotenv(env_file_path)


def read_file(filename):
    global FILE_DIRECTORY
    try:
        print(f"Current Directory: {os.getcwd()}:{os.listdir('.')}")
        print(f"Opening File: {FILE_DIRECTORY}/{filename}")
        print(f"File Directory: {FILE_DIRECTORY}")
        with open(f"{FILE_DIRECTORY}/{filename}", "r") as f:
            data = f.read().strip()
            return data
    except:
        return None


def parse_versions(data):
    versions = data.split(":")
    feature_name, feature_variant = versions[0].split(",")
    training_name, training_variant = versions[1].split(",")
    return feature_name, feature_variant, training_name, training_variant


def parse_feature(data):
    return data.split(":")


FEATURE_NAME, FEATURE_VARIANT, TRAININGSET_NAME, TRAININGSET_VARIANT = parse_versions(
    read_file("version.txt")
)
FEATURE_ENTITY, FEATURE_VALUE = parse_feature(read_file("feature.txt"))
if (
    FEATURE_NAME == None
    or FEATURE_VARIANT == None
    or TRAININGSET_NAME == None
    or TRAININGSET_VARIANT == None
):
    raise TypeError("VERSION is set to None.")

client = ff.ServingClient()


def serve_data():
    for _ in range(NUMBER_OF_SLEEPS):
        try:
            dataset = client.training_set(TRAININGSET_NAME, TRAININGSET_VARIANT)
            training_dataset = dataset.repeat(10).shuffle(1000).batch(8)
            for i, feature_batch in enumerate(training_dataset):
                if i >= 1:
                    return
                print(feature_batch.to_list())

        except Exception as e:
            print(f"\twaiting for {SLEEP_DURATION} seconds")
            print(e)
            time.sleep(SLEEP_DURATION)

    raise Exception(
        f"Serving for {TRAININGSET_NAME}:{TRAININGSET_VARIANT} could not be completed."
    )


def serve_feature():
    for _ in range(NUMBER_OF_SLEEPS):
        try:
            fpf = client.features(
                [(FEATURE_NAME, FEATURE_VARIANT)], {FEATURE_ENTITY: FEATURE_VALUE}
            )
            print(fpf)
            return
        except Exception as e:
            print(f"\twaiting for {SLEEP_DURATION} seconds")
            print(e)
            time.sleep(SLEEP_DURATION)
    raise Exception(
        f"Serving for {FEATURE_NAME}:{FEATURE_VARIANT}-{FEATURE_ENTITY}:{FEATURE_VALUE} could not be completed."
    )


print(f"Serving the training set ({TRAININGSET_NAME}:{TRAININGSET_VARIANT})")
serve_data()

print("\n")

print(f"Serving the feature for ({FEATURE_NAME}:{FEATURE_VARIANT})")
serve_feature()
