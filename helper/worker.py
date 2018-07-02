#!/usr/bin/env python

import click
import jinja2
import os
import uuid


from kcl_utils import run_kcl_process


APPLICATION_NAME = os.environ.get("APPLICATION_NAME", "KCLWorker")
EXECUTABLE_PATH = os.environ.get("EXECUTABLE_PATH")

TEMPLATE_CONFIG_FILE_PATH = os.path.abspath(os.path.join(__file__, "../template.properties"))
CONFIG_FILE_PATH_PATTERN = '/tmp/{stream_name}.properties'

JAVA_PATH = '/usr/bin/java'


def create_config_file(stream_name, region):
    """
    Create a KCL config file from the template coniguration

    :param stream_name: The Kinesis stream name
     :type stream_name: str
    :param region: The region the worker should run in
     :type region: str
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

    template_vars = dict(
        REGION=region,
        STREAM_NAME=stream_name,
        APPLICATION_NAME=APPLICATION_NAME,
        EXECUTABLE_PATH=EXECUTABLE_PATH,
        MAX_SHARDS=max_shards,
        WORKER_ID=container_id)

    config_file_path = CONFIG_FILE_PATH_PATTERN.format(stream_name=stream_name)
    with open(config_file_path, "w") as config_file:
        config_file.write(template.render(template_vars))

    return config_file_path


@click.command()
@click.argument('stream_name')
@click.argument('region')
def main(stream_name, region):
    config_file_path = create_config_file(stream_name, region)
    run_kcl_process(JAVA_PATH, config_file_path)


if __name__ == "__main__":
    main()
