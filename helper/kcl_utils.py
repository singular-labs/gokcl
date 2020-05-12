"""
Contains Helper functions for initializing the KCL process.
Mostly taken from the amazon_kclpy_helper.py file provided in the KCL python library, with some changes
to make it cleaner to use.
"""

import os
import subprocess
from glob import glob

from amazon_kclpy import kcl


MULTI_LANG_DAEMON_CLASS = 'software.amazon.kinesis.multilang.MultiLangDaemon'
DEBUG_CLIENT_CLASS = "DebugClient"


def get_dir_of_file(f):
    """
    Returns the absolute path to the directory containing the specified file.

    :param f: A path to a file, either absolute or relative
     :type f: str
    :return: The absolute path of the directory represented by the relative path provided.
     :rtype:  str
    """
    return os.path.dirname(os.path.abspath(f))


def get_kcl_dir():
    """
    Returns the absolute path to the dir containing the amazon_kclpy.kcl module.

    :return: The absolute path of the KCL package.
     :rtype: str
    """
    return get_dir_of_file(kcl.__file__)


def get_kcl_jar_path():
    """
    Returns the absolute path to the KCL jars needed to run an Amazon KCLpy app.

    :return: The absolute path of the KCL jar files needed to run the MultiLangDaemon.
     :rtype: str
    """
    return ':'.join(glob(os.path.join(get_kcl_dir(), 'jars', '*jar')))


def debug_main_file_path():
    return os.path.abspath(os.path.join(os.path.dirname(__file__), "../", "./src/main/java/singular"))


def get_kcl_classpath(properties=None, paths=None, debug=False):
    """
    Generates a classpath that includes the location of the kcl jars, the
    properties file and the optional paths.

    :param properties: Path to properties file.
     :type properties: str
    :param paths: List of strings. The paths that will be prepended to the classpath.
     :type paths: list
    :return: A java class path that will allow your properties to be found and the MultiLangDaemon and its deps and
             any custom paths you provided.
     :rtype: str
    """

    # First make all the user provided paths absolute
    paths = paths or []
    paths = [os.path.abspath(p) for p in paths]
    # We add our paths after the user provided paths because this permits users to
    # potentially inject stuff before our paths (otherwise our stuff would always
    # take precedence).

    if debug:
        debug_file_path = debug_main_file_path()
        paths.append(debug_file_path)

    paths.append(get_kcl_jar_path())
    if properties:
        # Add the dir that the props file is in
        dir_of_file = get_dir_of_file(properties)
        paths.append(dir_of_file)

    return ":".join([p for p in paths if p != ''])


def get_kcl_app_command(java, multi_lang_daemon_class, properties, paths=None, debug=False):
    """
    Generates a command to run the MultiLangDaemon.

    :param java: Path to java
     :type java: str
    :param multi_lang_daemon_class: Name of multi language daemon class e.g. com.amazonaws.services.kinesis.multilang.MultiLangDaemon
     :type multi_lang_daemon_class: str
    :param properties: Optional properties file to be included in the classpath.
     :type properties: str
    :param paths: List of strings. Additional paths to prepend to the classpath.
     :type paths: list
    :return: A command that will run the MultiLangDaemon with your properties and custom paths and java.
     :rtype: list
    """

    command = java
    args = [
        java,
    ]

    classpath = get_kcl_classpath(properties, paths, debug=debug)

    # This is to disable SSL verification in local stack kinesis :-(
    if debug:
        multi_lang_daemon_class = DEBUG_CLIENT_CLASS

    args.extend([
        "-cp",
        classpath,
        multi_lang_daemon_class,
        properties
    ])

    return command, args


def run_kcl_process(java_path, config_file_path, debug=False):
    """
    Runs the KCL process according to the given parameters (does not return).

    :param java_path: The path of the JAVA binary file
     :type java_path: str
    :param config_file_path: The path of the KCL properties file
     :type config_file_path: str
    """

    command, args = get_kcl_app_command(
        java_path,
        MULTI_LANG_DAEMON_CLASS,
        config_file_path,
        debug=debug)

    # Need to compile with the same args and env
    if debug:
        # Use javac
        build_args = [args[0] + "c"]

        # Need classpath, can omit properties
        build_args.extend(args[1:-1])

        # Need to specify full path of file to build
        build_args[-1] = os.path.join(debug_main_file_path(), "%s.java" % DEBUG_CLIENT_CLASS)

        subprocess.check_call(' '.join(build_args), shell=True, stdout=subprocess.PIPE)

    os.execv(command, args)
