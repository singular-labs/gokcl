#!/usr/bin/env python

import json
import os
import uuid

import click
import jinja2

from kcl_utils import run_kcl_process

APPLICATION_NAME = os.environ.get("APPLICATION_NAME", "KCLWorker")
EXECUTABLE_PATH = os.environ.get("EXECUTABLE_PATH")
RETRIEVAL_MODE = os.environ.get("RETRIEVAL_MODE", "FANOUT")
ENV_ID = os.environ.get("ENV_ID").lower()

TEMPLATE_CONFIG_FILE_PATH = os.path.abspath(os.path.join(__file__, "../template.properties"))
CONFIG_FILE_PATH_PATTERN = '/tmp/{stream_name}.properties'

JAVA_PATH = '/usr/bin/java'


def create_config_file(stream_name, region, extra=''):
    """
    Create a KCL config file from the template configuration

    :param stream_name: The Kinesis stream name
     :type stream_name: str
    :param region: The region the worker should run in
     :type region: str
    :param extra: Extra values for the config
     :type extra: str
    :return: The path of the config file
     :rtype: str
    """

    if EXECUTABLE_PATH is None:
        raise Exception("Executable path must be set!")

    with open(TEMPLATE_CONFIG_FILE_PATH) as config_file:
        config_template_data = config_file.read()

    env = jinja2.Environment(undefined=jinja2.StrictUndefined)
    template = env.from_string(config_template_data)

    max_shards = os.environ.get('MAX_SHARDS_PER_CONTAINER', '1024')
    container_id = str(uuid.uuid4())

    application_name = APPLICATION_NAME if ENV_ID == "prod" else ENV_ID + '-' + APPLICATION_NAME

    template_vars = dict(
        REGION=region,
        STREAM_NAME=stream_name,
        APPLICATION_NAME=application_name,
        EXECUTABLE_PATH=EXECUTABLE_PATH,
        MAX_SHARDS=max_shards,
        WORKER_ID=container_id,
        RETRIEVAL_MODE=RETRIEVAL_MODE
    )

    config_data = json.loads(extra) if len(extra) > 0 else {}

    config_file_path = CONFIG_FILE_PATH_PATTERN.format(stream_name=stream_name)
    with open(config_file_path, "w") as config_file:
        config_file.write(template.render(template_vars))
        for key, value in config_data.items():
            config_file.write('\n%(key)s = %(value)s\n' % locals())

    return config_file_path


@click.command()
@click.argument('stream_name')
@click.argument('region')
@click.option("-extra", type=str, default='',
              help="Optional mapping added into the config.")
def main(stream_name, region, extra):
    config_file_path = create_config_file(stream_name, region, extra)
    run_kcl_process(JAVA_PATH, config_file_path)


if __name__ == "__main__":
    main()
