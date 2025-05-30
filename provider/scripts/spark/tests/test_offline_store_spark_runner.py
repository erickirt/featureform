#  This Source Code Form is subject to the terms of the Mozilla Public
#  License, v. 2.0. If a copy of the MPL was not distributed with this
#  file, You can obtain one at http://mozilla.org/MPL/2.0/.
#
#  Copyright 2024 FeatureForm Inc.
#

import os
import sys

sys.path.insert(0, "provider/scripts/spark")

import pytest
from unittest.mock import patch
from deepdiff import DeepDiff

from offline_store_spark_runner import (
    main,
    parse_args,
    execute_sql_query,
    execute_df_job,
    set_spark_configs,
    download_blobs_to_local,
    split_key_value,
    get_credentials_dict,
    delete_file,
    check_dill_exception,
    get_s3_object,
)


real_path = os.path.realpath(__file__)
dir_path = os.path.dirname(real_path)


@pytest.mark.skipif(sys.platform.startswith("win"), reason="should not run on windows")
@pytest.mark.parametrize(
    "arguments",
    [
        "sql_local_all_arguments",
        "df_local_all_arguments",
        pytest.param("invalid_arguments", marks=pytest.mark.xfail),
        pytest.param("sql_invalid_local_arguments", marks=pytest.mark.xfail),
    ],
)
def test_main(arguments, request):
    expected_args = request.getfixturevalue(arguments)
    main(expected_args)


@pytest.mark.skipif(sys.platform.startswith("win"), reason="should not run on windows")
@pytest.mark.parametrize(
    "arguments",
    [
        "sql_all_arguments",
        "sql_partial_arguments",
        "df_all_arguments",
        "sql_databricks_all_arguments",
        pytest.param("df_partial_arguments", marks=pytest.mark.xfail),
        pytest.param("sql_invalid_arguments", marks=pytest.mark.xfail),
        pytest.param("df_invalid_arguments", marks=pytest.mark.xfail),
        pytest.param("invalid_arguments", marks=pytest.mark.xfail),
    ],
)
def test_parse_args(arguments, request):
    input_args, expected_args = request.getfixturevalue(arguments)
    args = parse_args(input_args)
    expected = vars(expected_args)
    found = vars(args)
    diff = DeepDiff(expected, found, verbose_level=2, ignore_order=True)
    assert not diff

@pytest.mark.skipif(sys.platform.startswith("win"), reason="should not run on windows")
@pytest.mark.parametrize(
    "arguments,expected_output",
    [
        ("sql_local_all_arguments", f"{dir_path}/test_files/input/transaction.parquet"),
        pytest.param(
            "sql_invalid_local_arguments",
            f"{dir_path}/test_files/expected/test_execute_sql_job_success",
            marks=pytest.mark.xfail,
        ),
    ],
)
def test_execute_sql_query(arguments, expected_output, spark, request):
    args = request.getfixturevalue(arguments)
    output_file = execute_sql_query(
        args.job_type,
        args.output,
        args.sql_query,
        args.spark_config,
        args.sources,
        args.output_format,
        args.headers,
        args.credential,
    )

    expected_df = spark.read.parquet(expected_output)
    output_df = spark.read.parquet(output_file)

    assert expected_df.count() == output_df.count()
    assert expected_df.schema == output_df.schema


@pytest.mark.skipif(sys.platform.startswith("win"), reason="should not run on windows")
@pytest.mark.parametrize(
    "arguments,expected_output",
    [
        ("df_local_all_arguments", f"{dir_path}/test_files/input/transaction.parquet"),
        pytest.param(
            "df_local_pass_none_code_failure",
            f"{dir_path}/test_files/expected/test_execute_df_job_success",
            marks=pytest.mark.xfail,
        ),
    ],
)
def test_execute_df_job(arguments, expected_output, spark, request):
    args = request.getfixturevalue(arguments)
    output_file = execute_df_job(
        args.output,
        args.code,
        args.store_type,
        args.spark_config,
        args.headers,
        args.credential,
        args.sources,
    )

    expected_df = spark.read.parquet(expected_output)
    output_df = spark.read.parquet(output_file)

    assert expected_df.count() == output_df.count()
    assert expected_df.schema == output_df.schema


def test_set_spark_config(sparkbuilder):
    config = {
        "fs.azure.account.key.account_name.dfs.core.windows.net": "adfjaidfasdklciadsj=="
    }

    set_spark_configs(sparkbuilder, config)
    spark = sparkbuilder.getOrCreate()

    for key, value in config.items():
        assert spark.conf.get(key) == value


@pytest.mark.skipif(sys.platform.startswith("win"), reason="should not run on windows")
@pytest.mark.hosted
def test_download_blobs_to_local(container_client):
    blob = "scripts/spark/python_packages.sh"
    local_filename = "python_packages.sh"
    output_file = download_blobs_to_local(container_client, blob, local_filename)

    assert os.path.exists(output_file)


def test_get_credentials_dict():
    input_base64_creds = "eyJ0ZXN0aW5nIjogImZpbGUifQ=="
    expected_output = {"testing": "file"}

    creds = get_credentials_dict(input_base64_creds)

    assert creds == expected_output


def test_delete_file(tmp_path):
    file_path = f"{tmp_path}/test.txt"
    print(tmp_path, file_path)
    with open(file_path, "w") as f:
        f.write("hi.world\n")

    assert os.path.isfile(file_path)
    delete_file(file_path)
    assert not os.path.isfile(file_path)
    delete_file(file_path)
    assert not os.path.isfile(file_path)


def test_split_key_value():
    key_values = ["a=b", "b=c", "c=b", "d=e=="]
    expected_output = {"a": "b", "b": "c", "c": "b", "d": "e=="}

    output = split_key_value(key_values)
    assert output == expected_output


@pytest.mark.skipif(sys.platform.startswith("win"), reason="should not run on windows")
@pytest.mark.parametrize(
    "exception_message, error",
    [
        (
            Exception("TypeError: code() takes at most 16 arguments (19 given)"),
            "dill_python_version_error",
        ),
        (
            Exception("unknown opcode"),
            "dill_python_version_error",
        ),
        (Exception("generic error"), "generic_error"),
    ],
)
def test_check_dill_exception(exception_message, error, request):
    expected_error = request.getfixturevalue(error)
    error = check_dill_exception(exception_message)
    assert str(error) == str(expected_error)


@pytest.mark.parametrize(
    "credentials",
    [
        {
            "aws_region": "us-west-2",
            "aws_access_key_id": "key",
            "aws_secret_access_key": "secret",
            "aws_bucket_name": "bucket",
        },
        {
            "aws_region": "us-west-2",
            "aws_access_key_id": "key",
            "aws_secret_access_key": "secret",
            "aws_bucket_name": "bucket",
            "use_service_account": "true",
        },
        # Invalid credentials, expected failure
        pytest.param(
            {},
            marks=pytest.mark.xfail,
        ),
    ],
)
def test_get_s3_object(credentials):
    class MockS3:
        @staticmethod
        def Object(bucket, filepath):
            return f"{bucket}/{filepath}"

    class MockSession:
        @staticmethod
        def resource(resource_name, region_name="us-east-1"):
            return MockS3()

    with patch("boto3.Session") as mock_session:
        mock_session.return_value = MockSession
        output = get_s3_object("file_path", credentials)
        assert output == "bucket/file_path"
